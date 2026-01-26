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

// type ProductStatus string
//
// const (
// 	ProductStatusDraft    ProductStatus = "draft"
// 	ProductStatusActive   ProductStatus = "active"
// 	ProductStatusArchived ProductStatus = "archived"
// )

type _Product struct {
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

type UpdateProductInput struct {
	Name             *string        `json:"name"`
	Description      *string        `json:"description"`
	ShortDescription *string        `json:"short_description"`
	MetaTitle        *string        `json:"meta_title,omitempty"`
	MetaDescription  *string        `json:"meta_description,omitempty"`
	Status           *ProductStatus `json:"status"`
	DefaultVariantID *int64         `json:"default_variant_id"`
}

type ProductModel struct {
	DB      *sql.DB // For compatibility
	queries *Queries
}

func (m ProductModel) Get(id int64) (*_Product, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	productRow, err := m.queries.GetProduct(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	product := mapProductRow(&productRow)
	return product, nil
}

func (m ProductModel) Create(product *Product) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	productRow, err := m.queries.InsertProduct(ctx, InsertProductParams{
		Name:             product.Name,
		Slug:             product.Slug,
		Description:      product.Description,
		ShortDescription: product.ShortDescription,
		MetaTitle:        product.MetaTitle,
		MetaDescription:  product.MetaDescription,
		Status:           product.Status,
		DefaultVariantID: product.DefaultVariantID,
	})
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

	product.ID = productRow.ID
	product.CreatedAt = productRow.CreatedAt
	product.UpdatedAt = productRow.UpdatedAt
	return nil
}

func (m ProductModel) Update(id int64, productInput *UpdateProductInput) (*Product, error) {
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

func mapProductRow(pr *GetProductRow) *_Product {
	product := &_Product{
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

func (m ProductModel) GetAll(limit, offset int) ([]*_Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	productRows, err := m.queries.ListProducts(ctx, ListProductsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, err
	}

	products := []*_Product{}

	for _, product := range productRows {
		p := &_Product{
			ID:               product.ID,
			Name:             product.Name,
			Slug:             product.Slug,
			Description:      &product.Description.String,
			ShortDescription: &product.ShortDescription.String,
			MetaTitle:        &product.MetaTitle.String,
			MetaDescription:  &product.MetaDescription.String,
			Status:           product.Status,
			DefaultVariantID: &product.DefaultVariantID.Int64,
			CreatedAt:        product.CreatedAt,
			UpdatedAt:        product.UpdatedAt,
			DeletedAt:        &product.DeletedAt.Time,
		}

		products = append(products, p)
	}

	return products, nil
}

func buildPatchQuery(id int64, input *UpdateProductInput) (string, []any, error) {
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
	if input.Status != nil {
		setClauses = append(setClauses, fmt.Sprintf("status = $%d", argPos))
		args = append(args, *input.Status)
		argPos++
	}

	if input.DefaultVariantID != nil {
		setClauses = append(setClauses, fmt.Sprintf("default_variant_id = $%d", argPos))
		args = append(args, *input.DefaultVariantID)
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
