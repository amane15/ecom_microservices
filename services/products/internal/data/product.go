package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
	"unicode/utf8"

	"github.com/amane15/ecom_microservice/services/proudcts/internal/validator"
	"github.com/lib/pq"
)

var (
	ErrDuplicateSlug = errors.New("slug already exists")
	ErrDuplicateSku  = errors.New("sku already exists")
)

type ProductStatus string

const (
	ProductStatusDraft     ProductStatus = "draft"
	ProductStatusPublished ProductStatus = "published"
	ProductStatusArchived  ProductStatus = "archived"
)

type Product struct {
	ID               int           `json:"id"`
	Name             string        `json:"name"`
	Slug             string        `json:"slug"`
	Description      string        `json:"description,omitempty"`
	ShortDescription string        `json:"short_description,omitempty"`
	MetaTitle        string        `json:"meta_title,omitempty"`
	MetaDescription  string        `json:"meta_description,omitempty"`
	Status           ProductStatus `json:"product_status"`
	DefaultVariantID int           `json:"default_variant_id"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at,omitempty"`
}

type ProductModel struct {
	DB *sql.DB
}

func (m ProductModel) Get(id int) (*Product, error) {
	return nil, nil
}

func (m ProductModel) Create(product *Product) error {
	query := `INSERT INTO products(name, slug, description, short_description,
				meta_title, meta_description, status, default_variant_id)
	 	VALUES($1, $2, $3, $4, $5, $6, $7, $8)
	 	RETURNING id, created_at, updated_at`

	args := []any{
		product.Name, product.Slug, product.Description,
		product.ShortDescription, product.MetaTitle,
		product.MetaDescription, product.Status, product.DefaultVariantID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(
		&product.ID,
		&product.CreatedAt,
		&product.UpdatedAt)
	if err != nil {
		if pqError, ok := err.(*pq.Error); ok {
			switch pqError.Code {
			case "23505":
				return handleUniqueViolationError(pqError)
			default:
				return err
			}
		}
		return err
	}
	return nil
}

func (m ProductModel) Update(product *Product) error {
	return nil
}

// Soft delete only
func (m ProductModel) Delete(id int) error {
	return nil
}

func ValidateProduct(v *validator.Validator, product *Product) {
	v.Check(product.Slug != "", "slug", "must be provided")
	v.Check(utf8.RuneCountInString(product.Slug) <= 128, "slug", "must not be more that 128 characters long")
	v.Check(validator.Matches(product.Slug, validator.HyphenatedRegex), "slug", "must be in a valid format. For e.g. lenovo-r5-16g")

	v.Check(product.Name != "", "name", "must be provided")
	v.Check(utf8.RuneCountInString(product.Name) <= 255, "name", "must not be more than 255 characters long")
}

func handleUniqueViolationError(pqError *pq.Error) error {
	switch pqError.Constraint {
	case "products_sku_key":
		return ErrDuplicateSku
	case "products_slug_key":
		return ErrDuplicateSlug
	default:
		return pqError
	}
}
