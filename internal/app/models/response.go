package models

import (
	"github.com/google/uuid"
)

type ProductResponse struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	Description   *string   `json:"description,omitempty"`
	Price         int64     `json:"price"`
	StockQuantity int       `json:"stock_quantity"`
}

func ProductToProductResponse(product Product) ProductResponse {
	return ProductResponse{
		ID:            product.ID,
		Name:          product.Name,
		Description:   product.Description,
		Price:         product.Price,
		StockQuantity: product.StockQuantity,
	}
}

func ProductsToProductsResponse(products []Product) []ProductResponse {
	result := make([]ProductResponse, len(products))
	for idx, p := range products {
		result[idx] = ProductToProductResponse(p)
	}
	return result
}
