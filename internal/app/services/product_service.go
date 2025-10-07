package services

import (
	"errors"

	"github.com/google/uuid"
	"github.com/rezbow/ecommerce/internal/app/models"
	"github.com/rezbow/ecommerce/internal/platform/database"
)

var (
	ErrProductNotFound = errors.New("product not found")
)

type IProductSvc interface {
	ListProducts(*models.Pagination) ([]models.Product, error)
	GetProduct(uuid.UUID) (*models.Product, error)
	CreateProduct(*models.ProductCreate) (*models.Product, error)
	UpdateProduct(uuid.UUID, *models.ProductUpdateRequest) (*models.Product, error)
}

type ProductService struct {
	productRepo database.IProductRepo
}

func NewProductService(repo database.IProductRepo) *ProductService {
	return &ProductService{
		productRepo: repo,
	}
}

func (svc *ProductService) GetProduct(id uuid.UUID) (*models.Product, error) {
	product, err := svc.productRepo.Get(id)
	if err != nil {
		if errors.Is(err, database.ErrRecordNotFound) {
			return nil, ErrProductNotFound
		}
		return nil, ErrInternal
	}
	return product, nil
}

func (svc *ProductService) ListProducts(pagination *models.Pagination) ([]models.Product, error) {
	products, err := svc.productRepo.GetPaged(pagination)
	if err != nil {
		return nil, ErrInternal
	}
	return products, nil
}

func (svc *ProductService) CreateProduct(productCreate *models.ProductCreate) (*models.Product, error) {
	product := &models.Product{
		Name:          productCreate.Name,
		Description:   productCreate.Description,
		Price:         productCreate.Price,
		StockQuantity: productCreate.StockQuantity,
	}

	if err := svc.productRepo.Create(product); err != nil {
		return nil, ErrInternal
	}

	return product, nil
}

func (svc *ProductService) UpdateProduct(productId uuid.UUID, productUpdateRequest *models.ProductUpdateRequest) (*models.Product, error) {
	product, err := svc.productRepo.Update(productId, productUpdateRequest.ToMap())
	if err != nil {
		if errors.Is(err, database.ErrRecordNotFound) {
			return nil, ErrProductNotFound
		}
		return nil, ErrInternal
	}
	return product, nil
}
