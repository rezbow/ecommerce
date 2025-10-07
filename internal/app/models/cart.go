package models

import (
	"github.com/google/uuid"
)

type Cart struct {
	UserId uuid.UUID   `json:"user_id"`
	Items  []*CartItem `json:"items"`
	Total  int64       `json:"total"`
}

func NewCart(userId uuid.UUID) *Cart {
	return &Cart{
		UserId: userId,
		Items:  make([]*CartItem, 0),
	}
}

func (c *Cart) SetItems(items []*CartItem) {
	c.Items = items
	c.update()
}

func (c *Cart) update() {
	var total int64
	for _, item := range c.Items {
		total += item.SubTotal
	}
	c.Total = total
}

func (c *Cart) Remove(productId uuid.UUID) bool {
	for idx, item := range c.Items {
		if item.ProductId == productId {
			c.Items = append(c.Items[:idx], c.Items[idx+1:]...)
			c.update()
			return true
		}
	}
	return false
}

func (c *Cart) findItem(id uuid.UUID) *CartItem {
	for _, item := range c.Items {
		if item.ProductId == id {
			return item
		}
	}
	return nil
}

func (c *Cart) ItemQuantity(id uuid.UUID) int {
	if item := c.findItem(id); item != nil {
		return item.Quantity
	}
	return 0
}

func (c *Cart) AddQuantityOrInsert(product *Product, quantity int) {
	if item := c.findItem(product.ID); item != nil {
		item.Quantity += quantity
		item.SubTotal = int64(item.Quantity) * item.Price
		c.update()
		return
	}
	item := NewCartItem(product.ID)
	item.Name = product.Name
	item.Price = product.Price
	item.Quantity = quantity
	item.SubTotal = int64(item.Quantity) * item.Price

	c.Items = append(c.Items, item)
	c.update()
}

type CartItem struct {
	ProductId uuid.UUID `json:"product_id"`
	Quantity  int       `json:"quantity"`
	Name      string    `json:"name"`
	Price     int64     `json:"price"`
	SubTotal  int64     `json:"subtotal"`
}

func NewCartItem(productId uuid.UUID) *CartItem {
	return &CartItem{
		ProductId: productId,
	}
}
