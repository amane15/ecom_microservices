package store

import (
	"context"
	"errors"
	"time"

	"github.com/amane15/ecom_microservice/internal/dbutils"
	"github.com/amane15/ecom_microservice/services/products/internal/data"
	"github.com/amane15/ecom_microservice/services/products/internal/service"
	"github.com/amane15/ecom_microservice/services/products/internal/store/sqlstore"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CategoryStore struct {
	DBPool *pgxpool.Pool
	Q      *sqlstore.Queries
}

func NewCategoryStore(dbpool *pgxpool.Pool) *CategoryStore {
	return &CategoryStore{
		DBPool: dbpool,
		Q:      sqlstore.New(dbpool),
	}
}

func (cs *CategoryStore) Get(ctx context.Context, id int64) (*data.Category, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	row, err := cs.Q.GetCategory(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, data.ErrRecordNotFound
		default:
			return nil, checkAndHandlePostgresErrors(err)
		}
	}

	category := &data.Category{
		ID:          row.ID,
		Name:        row.Name,
		Slug:        row.Slug,
		Description: dbutils.StringToPtr(row.Description),
		IsActive:    row.IsActive,
		CreatedAt:   row.CreatedAt.Time,
		UpdatedAt:   row.UpdatedAt.Time,
	}

	return category, nil
}

func (cs *CategoryStore) Create(ctx context.Context, input *service.CreateCategoryInput) (*data.Category, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	params := sqlstore.InsertCategoryParams{
		Name:        input.Name,
		Slug:        input.Slug,
		Description: dbutils.PtrToString(input.Description),
	}

	if input.IsActive != nil {
		params.IsActive = *input.IsActive
	}

	row, err := cs.Q.InsertCategory(ctx, params)
	if err != nil {
		return nil, checkAndHandlePostgresErrors(err)
	}

	category := &data.Category{
		ID:          row.ID,
		Name:        row.Name,
		Slug:        row.Slug,
		Description: dbutils.StringToPtr(row.Description),
		IsActive:    row.IsActive,
		CreatedAt:   row.CreatedAt.Time,
		UpdatedAt:   row.UpdatedAt.Time,
	}

	return category, nil
}

func (cs *CategoryStore) Update(ctx context.Context, id int64, input *service.UpdateCategoryInput) (*data.Category, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	params := sqlstore.UpdateCategoryParams{
		ID:          id,
		Name:        dbutils.PtrToString(input.Name),
		Description: dbutils.PtrToString(input.Description),
		IsActive:    dbutils.PtrToBool(input.IsActive),
	}

	row, err := cs.Q.UpdateCategory(ctx, params)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, data.ErrRecordNotFound
		default:
			return nil, checkAndHandlePostgresErrors(err)
		}
	}

	category := &data.Category{
		ID:          row.ID,
		Name:        row.Name,
		Slug:        row.Slug,
		Description: dbutils.StringToPtr(row.Description),
		IsActive:    row.IsActive,
		CreatedAt:   row.CreatedAt.Time,
		UpdatedAt:   row.UpdatedAt.Time,
	}

	return category, nil
}

func (cs *CategoryStore) Delete(ctx context.Context, id int64) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	err := cs.Q.DeleteCategory(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return data.ErrRecordNotFound
		default:
			return checkAndHandlePostgresErrors(err)
		}
	}

	return nil
}

func (cs *CategoryStore) List(ctx context.Context, limit, offset int32) ([]*data.Category, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	rows, err := cs.Q.ListCategories(ctx, sqlstore.ListCategoriesParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, checkAndHandlePostgresErrors(err)
	}

	categories := []*data.Category{}

	for _, row := range rows {
		category := &data.Category{
			ID:          row.ID,
			Name:        row.Name,
			Slug:        row.Slug,
			Description: dbutils.StringToPtr(row.Description),
			IsActive:    row.IsActive,
			CreatedAt:   row.CreatedAt.Time,
			UpdatedAt:   row.UpdatedAt.Time,
		}

		categories = append(categories, category)
	}

	return categories, nil
}
