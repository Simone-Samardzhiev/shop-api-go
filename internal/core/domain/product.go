package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// Product is an entity representing a product.
type Product struct {
	Id          uuid.UUID
	Name        string
	Description string
	Price       decimal.Decimal
	Rating      decimal.Decimal
	Count       int
	ImageUrl    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewProduct creates a new Product instance.
func NewProduct(
	id uuid.UUID,
	name string,
	description string,
	price decimal.Decimal,
	rating decimal.Decimal,
	count int,
	imageUrl string,
	createdAt time.Time,
	updatedAt time.Time,
) *Product {
	return &Product{
		Id:          id,
		Name:        name,
		Description: description,
		Price:       price,
		Rating:      rating,
		Count:       count,
		ImageUrl:    imageUrl,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}
