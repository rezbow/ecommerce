package models

import "github.com/google/uuid"

type RegisterUser struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type Login struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ProductCreate struct {
	Name          string  `json:"name" binding:"required"`
	Description   *string `json:"description"`
	Price         int64   `json:"price" binding:"required"`
	StockQuantity int     `json:"stock_quantity" binding:"required"`
}

func (p *ProductCreate) Validate() (bool, map[string]string) {
	errs := make(map[string]string)
	if len(p.Name) < 2 {
		errs["name"] = "name must have more than 2 characters"
	}
	if p.Description != nil && len(*p.Description) < 10 {
		errs["description"] = "description must have more than 10 characters"
	}
	if p.Price <= 0 {
		errs["price"] = "price must be greater than 0"
	}

	if p.StockQuantity <= 0 {
		errs["stock_quantity"] = "stock quantity must be greater than 0"
	}
	return len(errs) == 0, errs
}

type ProductUpdateRequest struct {
	Name          *string `json:"name"`
	Description   *string `json:"description"`
	Price         *int64  `json:"price"`
	StockQuantity *int    `json:"stock_quantity"`
}

func (p *ProductUpdateRequest) Validate() (bool, map[string]string) {
	errs := make(map[string]string)
	if p.Name != nil && len(*p.Name) < 2 {
		errs["name"] = "name must have more than 2 characters"
	}
	if p.Description != nil && len(*p.Description) < 10 {
		errs["description"] = "description must have more than 10 characters"
	}
	if p.Price != nil && *p.Price <= 0 {
		errs["price"] = "price must be greater than 0"
	}

	if p.StockQuantity != nil && *p.StockQuantity <= 0 {
		errs["stock_quantity"] = "stock quantity must be greater than 0"
	}
	return len(errs) == 0, errs
}

func (p *ProductUpdateRequest) ToMap() map[string]any {
	result := make(map[string]any)
	if p.Name != nil {
		result["name"] = *p.Name
	}
	if p.Description != nil {
		result["description"] = *p.Description
	}
	if p.Price != nil {
		result["price"] = *p.Price
	}
	if p.StockQuantity != nil {
		result["stock_quantity"] = *p.StockQuantity
	}
	return result
}

type ItemCartRequest struct {
	ProductId uuid.UUID `json:"product_id" binding:"required"`
	Quantity  int       `json:"quantity" binding:"required"`
}

func (i *ItemCartRequest) Validate() (bool, map[string]string) {
	errs := make(map[string]string)
	if i.Quantity <= 0 {
		errs["quantity"] = "quantity should be greater than 0"
	}
	return len(errs) == 0, errs
}

type ItemQuantityUpdate struct {
	NewQuantity int `json:"new_quantity" binding:"required"`
}

func (i *ItemQuantityUpdate) Validate() (bool, map[string]string) {
	errs := make(map[string]string)
	if i.NewQuantity <= 0 {
		errs["new_quantity"] = "new_quantity should be greater than 0"
	}
	return len(errs) == 0, errs
}
