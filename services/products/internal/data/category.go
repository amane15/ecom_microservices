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

// id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
//
// name VARCHAR(255) NOT NULL,
// slug VARCHAR(128) NOT NULL UNIQUE,
//
// description TEXT,
//
// is_active BOOLEAN NOT NULL DEFAULT true,
//
// created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
// updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
// deleted_at TIMESTAMPTZ
//

type CategoryModel struct {
	DB *sql.DB
}

type Category struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Slug        string  `json:"slug"`
	Description *string `json:"description,omitempty"`
	IsActive    bool    `json:"is_active"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type CreateCategoryInput struct {
	Name        string  `json:"name"`
	Slug        string  `json:"slug"`
	Description *string `json:"description"`
	IsActive    *bool   `json:"is_active"`
}

type UpdateCategoryInput struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	IsActive    *bool   `json:"is_active"`
}

func (m CategoryModel) Get(id int64) (*Category, error) {
	query := `
	SELECT id, name, slug, description, is_active, 
		created_at, updated_at, deleted_at
	FROM categories
	WHERE id = $1
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	category := &Category{}
	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&category.ID,
		&category.Name,
		&category.Slug,
		&category.Description,
		&category.IsActive,
		&category.CreatedAt,
		&category.UpdatedAt,
		&category.DeletedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return category, nil
}

func (m CategoryModel) Insert(category *Category) error {
	query := `
	INSERT INTO categories(name, slug, description, is_active)
	VALUES ($1, $2, $3, $4)
	RETURNING id, created_at, updated_at
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{category.Name, category.Slug, category.Description, category.IsActive}
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(
		&category.ID,
		&category.CreatedAt,
		&category.UpdatedAt,
	)
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

func (m CategoryModel) Update(id int64, input *UpdateCategoryInput) (*Category, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query, args, err := buildPatchCategoryQuery(id, input)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	category := &Category{}
	deletedAt := sql.NullTime{}
	err = m.DB.QueryRowContext(ctx, query, args...).Scan(
		&category.ID,
		&category.Name,
		&category.Slug,
		&category.Description,
		&category.IsActive,
		&category.CreatedAt,
		&category.UpdatedAt,
		&deletedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	if deletedAt.Valid {
		category.DeletedAt = &deletedAt.Time
	}

	return category, nil
}

func (m CategoryModel) GetAll(limit, offset int) ([]*Category, error) {
	query := `
	SELECT id, name, slug, description, is_active,
		created_at, updated_at, deleted_at
	FROM categories
	WHERE is_active = true
	ORDER BY id
	LIMIT $1 OFFSET $2
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := []*Category{}

	deletedAt := sql.NullTime{}
	for rows.Next() {
		var category Category

		err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.Slug,
			&category.Description,
			&category.IsActive,
			&category.CreatedAt,
			&category.UpdatedAt,
			&deletedAt,
		)
		if err != nil {
			return nil, err
		}

		if deletedAt.Valid {
			category.DeletedAt = &deletedAt.Time
		}

		categories = append(categories, &category)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

func buildPatchCategoryQuery(id int64, input *UpdateCategoryInput) (string, []any, error) {
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

	if input.IsActive != nil {
		setClauses = append(setClauses, fmt.Sprintf("is_active = $%d", argPos))
		args = append(args, *input.IsActive)
		argPos++
	}

	if len(setClauses) == 0 {
		return "", nil, ErrNoFieldsToUpdate
	}

	query := fmt.Sprintf(`
	UPDATE categories
	SET %s, updated_at = now()
	WHERE id = $%d 
	RETURNING id, name, slug, description, is_active,
	created_at, updated_at, deleted_at
	`, strings.Join(setClauses, ", "), argPos)

	args = append(args, id)

	return query, args, nil
}
