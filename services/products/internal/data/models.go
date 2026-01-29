package data

import (
	"time"

	"github.com/shopspring/decimal"
)

type ProductStatus string

const (
	ProductStatusDraft    ProductStatus = "draft"
	ProductStatusActive   ProductStatus = "active"
	ProductStatusArchived ProductStatus = "archived"
)

type Product struct {
	ID               int64
	Name             string
	Slug             string
	Description      *string
	ShortDescription *string
	MetaTitle        *string
	MetaDescription  *string
	Status           ProductStatus
	DefaultVariantID *int64

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type Category struct {
	ID          int64
	Name        string
	Slug        string
	Description *string
	IsActive    bool

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type Variant struct {
	ID        int64
	ProductID int64
	Slug      string
	Name      string
	Price     decimal.Decimal
	IsActive  bool

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type ValidationErrors struct {
	Fields map[string]string
}

func (v ValidationErrors) Error() string {
	return "failed validation"
}
