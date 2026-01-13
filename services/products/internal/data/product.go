package data

import (
	"database/sql"
	"time"
)

type Product struct {
	ID               int    `json:"id"`
	Sku              string `json:"sku"`
	Name             string `json:"name"`
	Slug             string `json:"slug"`
	Description      string `json:"description"`
	ShortDescription string `json:"short_description"`
	MetaTitle        string `json:"meta_title"`
	MetaDescription  string `json:"meta_description"`
	IsActive         bool   `json:"is_active"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

type ProductModel struct {
	DB *sql.DB
}

func (m ProductModel) Get(id int) (*Product, error) {
	return nil, nil
}

func (m ProductModel) Create(product *Product) error {
	return nil
}

func (m ProductModel) Update(product *Product) error {
	return nil
}

// Soft delete only
func (m ProductModel) Delete(id int) error {
	return nil
}
