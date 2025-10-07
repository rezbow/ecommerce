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
	// sync the entire cart with database
	svc.SyncCart(cart)
	// save it to cache
	if err := svc.cartRepo.Save(userId.String(), cart, cacheDuration); err != nil {
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

	if itemCartRequest.Quantity+userCart.ItemQuantity(product.ID) > product.StockQuantity {
		return ErrInsufficientQuantity
	}

	// add to item quantity or insert if item with product.ID
	// doesn't exists in our cart
	userCart.AddQuantityOrInsert(product, itemCartRequest.Quantity)

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

func (svc *CartService) SyncCart(cart *models.Cart) {
	refreshedItems := make([]*models.CartItem, 0)
	for _, item := range cart.Items {
		product, err := svc.productRepo.Get(item.ProductId)
		if err != nil {
			continue
		}
		if product.StockQuantity == 0 {
			continue
		}

		newItem := models.NewCartItem(item.ProductId)
		if product.StockQuantity < item.Quantity {
			newItem.Quantity = product.StockQuantity
		} else {
			newItem.Quantity = item.Quantity
		}
		newItem.Price = product.Price
		newItem.Name = product.Name
		newItem.SubTotal = int64(newItem.Quantity) * newItem.Price

		refreshedItems = append(refreshedItems, newItem)

	}
	cart.SetItems(refreshedItems)
}
