package entity

import (
	"github.com/google/uuid"
)

type ProductRepository interface {
	Create(product *Product) error
	FindAll() ([]*Product, error)
}

type Product struct {
	ID    string
	Name  string
	Price float32
}

func NewProduct(name string, price float32) *Product {
	return &Product{
		ID:    uuid.New().String(),
		Name:  name,
		Price: price,
	}
}
