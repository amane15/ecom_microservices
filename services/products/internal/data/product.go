package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/lib/pq"
)

var (
	ErrDuplicateSlug    = errors.New("slug already exists")
	ErrDuplicateSku     = errors.New("sku already exists")
	ErrRecordNotFound   = errors.New("record not found")
	ErrNoFieldsToUpdate = errors.New("no fields to update")
)

type ProductStatus string

const (
	ProductStatusDraft     ProductStatus = "draft"
	ProductStatusPublished ProductStatus = "published"
	ProductStatusArchived  ProductStatus = "archived"
)

type Product struct {
	ID               int64         `json:"id"`
	Name             string        `json:"name"`
	Slug             string        `json:"slug"`
	Description      *string       `json:"description,omitempty"`
	ShortDescription *string       `json:"short_description,omitempty"`
	MetaTitle        *string       `json:"meta_title,omitempty"`
	MetaDescription  *string       `json:"meta_description,omitempty"`
	Status           ProductStatus `json:"product_status"`
	DefaultVariantID *int64        `json:"default_variant_id"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type productRow struct {
	ID               int64
	Name             string
	Slug             string
	Description      sql.NullString
	ShortDescription sql.NullString
	MetaTitle        sql.NullString
	MetaDescription  sql.NullString
	Status           ProductStatus
	DefaultVariantID sql.NullInt64

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt sql.NullTime
}

type UpdateProductRow struct {
	Name             *string `json:"name"`
	Description      *string `json:"description"`
	ShortDescription *string `json:"short_description"`
	MetaTitle        *string `json:"meta_title,omitempty"`
	MetaDescription  *string `json:"meta_description,omitempty"`
}

type ProductModel struct {
	DB *sql.DB
}

func (m ProductModel) Get(id int64) (*Product, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
	SELECT id, name, slug, description, short_description, meta_title,
		meta_description, status, default_variant_id,
		created_at, updated_at, deleted_at
	FROM products
	WHERE id = $1
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	productRow := &productRow{}

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&productRow.ID,
		&productRow.Name,
		&productRow.Slug,
		&productRow.Description,
		&productRow.ShortDescription,
		&productRow.MetaTitle,
		&productRow.MetaDescription,
		&productRow.Status,
		&productRow.DefaultVariantID,
		&productRow.CreatedAt,
		&productRow.UpdatedAt,
		&productRow.DeletedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	product := mapProductRow(productRow)
	return product, nil
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

func (m ProductModel) Update(id int64, productInput *UpdateProductRow) (*Product, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	query, args, err := buildPatchQuery(id, productInput)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	product := &Product{}
	err = m.DB.QueryRowContext(ctx, query, args...).Scan(
		&product.ID,
		&product.Name,
		&product.Slug,
		&product.Description,
		&product.ShortDescription,
		&product.MetaTitle,
		&product.MetaDescription,
		&product.Status,
		&product.DefaultVariantID,
		&product.CreatedAt,
		&product.UpdatedAt,
		&product.DeletedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return product, nil
}

// Soft delete only
func (m ProductModel) Delete(id int) error {
	return nil
}

func andleUniqueViolationError(pqError *pq.Error) error {
	switch pqError.Constraint {
	case "products_sku_key":
		return ErrDuplicateSku
	case "products_slug_key":
		return ErrDuplicateSlug
	default:
		return pqError
	}
}

func mapProductRow(pr *productRow) *Product {
	product := &Product{
		ID:        pr.ID,
		Name:      pr.Name,
		Slug:      pr.Slug,
		Status:    pr.Status,
		CreatedAt: pr.CreatedAt,
		UpdatedAt: pr.UpdatedAt,
	}

	if pr.Description.Valid {
		product.Description = &pr.Description.String
	}

	if pr.ShortDescription.Valid {
		product.ShortDescription = &pr.ShortDescription.String
	}

	if pr.MetaTitle.Valid {
		product.MetaTitle = &pr.MetaTitle.String
	}

	if pr.MetaDescription.Valid {
		product.MetaDescription = &pr.MetaDescription.String
	}

	if pr.DefaultVariantID.Valid {
		product.DefaultVariantID = &pr.DefaultVariantID.Int64
	}

	if pr.DeletedAt.Valid {
		product.DeletedAt = &pr.DeletedAt.Time
	}

	return product
}

func buildPatchQuery(id int64, input *UpdateProductRow) (string, []any, error) {
	setClauses := []string{}
	args := []any{}
	argPos := 1

	if input.Name != nil {
		setClauses = append(setClauses, fmt.Sprintf("name = $%d", argPos))
		args = append(args, *input.Name)
		argPos++
	}

	if input.Description != nil {
		setClauses = append(setClauses, fmt.Sprintf("description = $%d", argPos))
		args = append(args, *input.Description)
		argPos++
	}

	if input.ShortDescription != nil {
		setClauses = append(setClauses, fmt.Sprintf("short_description = $%d", argPos))
		args = append(args, *input.ShortDescription)
		argPos++
	}

	if input.MetaTitle != nil {
		setClauses = append(setClauses, fmt.Sprintf("meta_title = $%d", argPos))
		args = append(args, *input.MetaTitle)
		argPos++
	}
	if input.MetaDescription != nil {
		setClauses = append(setClauses, fmt.Sprintf("meta_description = $%d", argPos))
		args = append(args, *input.MetaDescription)
		argPos++
	}

	if len(setClauses) == 0 {
		return "", nil, ErrNoFieldsToUpdate
	}

	query := fmt.Sprintf(`
	UPDATE products
	SET %s, updated_at = now()
	WHERE id = $%d 
	RETURNING id, name, slug, description, short_description, meta_title, 
	meta_description, status, default_variant_id, created_at, updated_at, deleted_at
	`, strings.Join(setClauses, ", "), argPos)

	args = append(args, id)

	return query, args, nil
}
