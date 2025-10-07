package services

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/rezbow/ecommerce/internal/app/models"
	"github.com/rezbow/ecommerce/internal/platform/database"
)

const cacheDuration = time.Hour * 24 * 7

var (
	ErrCartNotFound         = errors.New("cart not found")
	ErrInsufficientQuantity = errors.New("not enoguth product qunatity in stock")
	ErrItemNotFound         = errors.New("item not found in cart")
)

type ICartService interface {
	GetUserCart(uuid.UUID) (*models.Cart, error)
	AddToUserCart(uuid.UUID, *models.ItemCartRequest) error
	RemoveItemFromCart(uuid.UUID, uuid.UUID) error
	ClearCart(uuid.UUID) error
}

type CartService struct {
	cartRepo    database.ICartRepo
	productRepo database.IProductRepo
}

func NewCartService(cartRepo database.ICartRepo, productRepo database.IProductRepo) *CartService {
	return &CartService{
		cartRepo:    cartRepo,
		productRepo: productRepo,
	}
}

func (svc *CartService) GetUserCart(userId uuid.UUID) (*models.Cart, error) {
	cart, err := svc.cartRepo.Get(userId.String())
	if err != nil {
		if errors.Is(err, database.ErrRecordNotFound) {
			return nil, ErrCartNotFound
		}
		return nil, ErrInternal
	}
	return cart, nil
}

func (svc *CartService) AddToUserCart(userId uuid.UUID, itemCartRequest *models.ItemCartRequest) error {
	// get users cart
	userCart, err := svc.cartRepo.Get(userId.String())
	if err != nil {
		if errors.Is(err, database.ErrRecordNotFound) {
			userCart = models.NewCart(userId)
		} else {
			return ErrInternal
		}
	}

	// build an item
	// call db for live info
	product, err := svc.productRepo.Get(itemCartRequest.ProductId)
	if err != nil {
		if errors.Is(err, database.ErrRecordNotFound) {
			return ErrProductNotFound
		}
		return ErrInternal
	}

	item := userCart.FindItem(itemCartRequest.ProductId)
	if item == nil {
		item = models.NewCartItem(product.ID)
		userCart.Add(item)
	}

	if itemCartRequest.Quantity+item.Quantity > product.StockQuantity {
		return ErrInsufficientQuantity
	}

	// update item state
	item.AddToQuantity(itemCartRequest.Quantity)
	item.Sync(product)
	// update cart state
	userCart.Update()

	// add it to cache
	if err := svc.cartRepo.Save(userId.String(), userCart, cacheDuration); err != nil {
		return ErrInternal
	}
	return nil
}

func (svc *CartService) RemoveItemFromCart(userId uuid.UUID, productId uuid.UUID) error {
	userCart, err := svc.cartRepo.Get(userId.String())
	if errors.Is(err, database.ErrRecordNotFound) {
		return ErrCartNotFound
	}
	if err != nil {
		return ErrInternal
	}

	if !userCart.Remove(productId) {
		return ErrItemNotFound
	}

	// update cart state
	userCart.Update()

	if err := svc.cartRepo.Save(userId.String(), userCart, cacheDuration); err != nil {
		return ErrInternal
	}
	return nil
}

func (svc *CartService) ClearCart(userId uuid.UUID) error {
	err := svc.cartRepo.Delete(userId.String())
	if errors.Is(err, database.ErrRecordNotFound) {
		return ErrCartNotFound
	}
	if err != nil {
		return ErrInternal
	}
	return nil
}
