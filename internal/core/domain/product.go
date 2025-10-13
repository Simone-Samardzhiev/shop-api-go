package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// CategorySection is an entity representing a category section.
type CategorySection struct {
	Id   uuid.UUID
	Name string
}

// Category is an entity representing a category.
type Category struct {
	Id        uuid.UUID
	Name      string
	SectionId uuid.UUID
}

// Subcategory is an entity representing a category
type Subcategory struct {
	Id         uuid.UUID
	Name       string
	CategoryID uuid.UUID
}

// Product is an entity representing a product.
type Product struct {
	Id            uuid.UUID
	Name          string
	Description   string
	Price         decimal.Decimal
	Rating        decimal.Decimal
	Count         int
	ImageUrl      string
	Subcategories []Subcategory
	CreatedAt     time.Time
	UpdatedAt     time.Time
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
	subcategories []Subcategory,
	createdAt time.Time,
	updatedAt time.Time,
) *Product {
	return &Product{
		Id:            id,
		Name:          name,
		Description:   description,
		Price:         price,
		Rating:        rating,
		Count:         count,
		ImageUrl:      imageUrl,
		Subcategories: subcategories,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}
}
