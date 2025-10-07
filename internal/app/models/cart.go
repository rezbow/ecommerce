package models

import "github.com/google/uuid"

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

func (c *Cart) Update() {
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
			return true
		}
	}
	return false
}

func (c *Cart) Add(item *CartItem) {
	c.Items = append(c.Items, item)
}

func (c *Cart) FindItem(id uuid.UUID) *CartItem {
	for _, item := range c.Items {
		if item.ProductId == id {
			return item
		}
	}
	return nil
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

func (item *CartItem) calculateSubTotal() {
	item.SubTotal = item.Price * int64(item.Quantity)
}

func (item *CartItem) Sync(product *Product) {
	item.Name = product.Name
	item.Price = product.Price
	item.calculateSubTotal()
}

func (item *CartItem) AddToQuantity(n int) {
	item.Quantity += n
	item.calculateSubTotal()
}
