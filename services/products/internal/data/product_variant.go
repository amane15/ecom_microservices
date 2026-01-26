package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/lib/pq"
	"github.com/shopspring/decimal"
)

type ProductVariant struct {
	ID        int64           `json:"id"`
	ProductID int64           `json:"product_id"`
	Slug      string          `json:"slug"`
	Name      string          `json:"name"`
	Price     decimal.Decimal `json:"price"`
	IsActive  bool            `json:"is_active"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type UpdateVariantInput struct {
	Name     *string          `json:"name"`
	Price    *decimal.Decimal `json:"price"`
	IsActive *bool            `json:"is_active"`
}

type CreateVariantInput struct {
	ProductID int64           `json:"product_id"`
	Slug      string          `json:"slug"`
	Name      string          `json:"name"`
	Price     decimal.Decimal `json:"price"`
	IsActive  *bool           `json:"is_active"`
}

type ProductVariantModel struct {
	DB      *sql.DB
	queries *Queries
}

func (m ProductVariantModel) Get(id int64) (*ProductsVariant, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	variant, err := m.queries.GetVariant(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &variant, nil
}

func (m ProductVariantModel) Insert(variant *ProductVariant) error {
	query := `
	INSERT INTO products_variants(product_id, slug, name, price, is_active)
	VALUES($1, $2, $3, $4, $5)
	RETURNING id, created_at, updated_at
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{variant.ProductID, variant.Slug, variant.Name, variant.Price, variant.IsActive}
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(
		&variant.ID,
		&variant.CreatedAt,
		&variant.UpdatedAt,
	)
	if err != nil {
		if pqError, ok := err.(*pq.Error); ok {
			if pqError.Code == "23505" {
				return ErrDuplicateSlug
			}
		}
		return err
	}

	return nil
}

func (m ProductVariantModel) Update(id int64, variantInput *UpdateVariantInput) (*ProductVariant, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query, args, err := buildPatchVariantQuery(id, variantInput)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	variant := &ProductVariant{}
	deleteAt := sql.NullTime{}
	err = m.DB.QueryRowContext(ctx, query, args...).Scan(
		&variant.ID,
		&variant.ProductID,
		&variant.Slug,
		&variant.Name,
		&variant.Price,
		&variant.IsActive,
		&variant.CreatedAt,
		&variant.UpdatedAt,
		&deleteAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	if deleteAt.Valid {
		variant.DeletedAt = &deleteAt.Time
	}

	return variant, nil
}

func (m ProductVariantModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.queries.DeleteVariant(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (m ProductVariantModel) GetVariantCountForProduct(id int64) (int64, error) {
	if id < 1 {
		return 0, ErrRecordNotFound
	}
	query := `SELECT count(*) FROM products_variants WHERE product_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var count int64
	err := m.DB.QueryRowContext(ctx, query, id).Scan(&count)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return 0, ErrRecordNotFound
		default:
			return 0, err
		}
	}

	return count, nil
}

func (m ProductVariantModel) IsVariantExists(id int64) (bool, error) {
	if id < 1 {
		return false, ErrRecordNotFound
	}

	query := `SELECT EXISTS(SELECT 1 FROM products_variants WHERE id = $1) AS "exists"`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	exists := false
	err := m.DB.QueryRowContext(ctx, query, id).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (m ProductVariantModel) GetVariantsByProduct(productID int64) ([]*ProductVariant, error) {
	if productID < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
	SELECT id, product_id, slug, name, price, is_active,
		created_at, updated_at, deleted_at
	FROM products_variants
	WHERE product_id = $1
	ORDER BY id
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	variants := []*ProductVariant{}

	for rows.Next() {
		var variant ProductVariant

		err := rows.Scan(
			&variant.ID,
			&variant.ProductID,
			&variant.Slug,
			&variant.Name,
			&variant.Price,
			&variant.IsActive,
			&variant.CreatedAt,
			&variant.UpdatedAt,
			&variant.DeletedAt,
		)
		if err != nil {
			return nil, err
		}

		variants = append(variants, &variant)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return variants, nil
}

func buildPatchVariantQuery(id int64, input *UpdateVariantInput) (string, []any, error) {
	setClauses := []string{}
	args := []any{}
	argPos := 1

	if input.Name != nil {
		setClauses = append(setClauses, fmt.Sprintf("name = $%d", argPos))
		args = append(args, *input.Name)
		argPos++
	}

	if input.Price != nil {
		setClauses = append(setClauses, fmt.Sprintf("price = $%d::DECIMAL(6, 2)", argPos))
		args = append(args, *input.Price)
		argPos++
	}

	if input.IsActive != nil {
		setClauses = append(setClauses, fmt.Sprintf("is_active = $%d", argPos))
		args = append(args, *input.IsActive)
		argPos++
	}

	if len(setClauses) == 0 {
		return "", nil, ErrNoFieldsToUpdate
	}

	query := fmt.Sprintf(`
	UPDATE products_variants
	SET %s, updated_at = now()
	WHERE id = $%d 
	RETURNING id, product_id, slug, name, price, is_active,
	created_at, updated_at, deleted_at 
	`, strings.Join(setClauses, ", "), argPos)

	args = append(args, id)

	return query, args, nil
}
