package database

import (
	"errors"

	"github.com/google/uuid"
	"github.com/rezbow/ecommerce/internal/app/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type IProductRepo interface {
	Get(uuid.UUID) (*models.Product, error)
	GetPaged(*models.Pagination) ([]models.Product, error)
	Create(*models.Product) error
	Update(uuid.UUID, map[string]any) (*models.Product, error)
}

type ProductRepo struct {
	db *gorm.DB
}

func NewProductRepo(db *gorm.DB) *ProductRepo {
	return &ProductRepo{db: db}
}

func (repo *ProductRepo) Get(id uuid.UUID) (*models.Product, error) {
	var product models.Product
	if err := repo.db.First(&product, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRecordNotFound
		}
		return nil, ErrInternal
	}
	return &product, nil
}

func (repo *ProductRepo) Create(product *models.Product) error {
	product.ID = uuid.New()
	if err := repo.db.Create(product).Error; err != nil {
		return ErrInternal
	}
	return nil
}

func (repo *ProductRepo) GetPaged(pagination *models.Pagination) ([]models.Product, error) {
	var products []models.Product
	err := repo.db.Model(models.Product{}).Offset(pagination.Offset).Limit(pagination.Limit).Find(&products).Error
	if err != nil {
		return nil, ErrInternal
	}
	return products, nil
}

func (repo *ProductRepo) Update(id uuid.UUID, updatedColumns map[string]any) (*models.Product, error) {
	product := models.Product{
		ID: id,
	}
	result := repo.db.Model(&product).Clauses(clause.Returning{}).Updates(updatedColumns)
	if result.Error != nil {
		return nil, ErrInternal
	}
	if result.RowsAffected == 0 {
		return nil, ErrRecordNotFound
	}
	return &product, nil
}
