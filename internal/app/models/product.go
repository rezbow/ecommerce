package models

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID            uuid.UUID
	Name          string
	Description   *string
	Price         int64
	StockQuantity int
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type ProductRepo interface {
	Get(uuid.UUID) (*Product, error)
	GetPaged() ([]*Product, error)
	Create(*Product) error
}
