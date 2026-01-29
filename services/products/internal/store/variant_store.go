package store

import (
	"context"
	"time"

	"github.com/amane15/ecom_microservice/internal/dbutils"
	"github.com/amane15/ecom_microservice/services/products/internal/data"
	"github.com/amane15/ecom_microservice/services/products/internal/service"
	"github.com/amane15/ecom_microservice/services/products/internal/store/sqlstore"
	"github.com/jackc/pgx/v5/pgxpool"
)

type VariantStore struct {
	DBPool *pgxpool.Pool
	Q      *sqlstore.Queries
}

func NewVariantStore(dbpool *pgxpool.Pool) *VariantStore {
	return &VariantStore{
		DBPool: dbpool,
		Q:      sqlstore.New(dbpool),
	}
}

func (vs *VariantStore) Get(ctx context.Context, id int64) (*data.Variant, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	row, err := vs.Q.GetVariant(ctx, id)
	if err != nil {
		return nil, err
	}

	variant := &data.Variant{
		ID:        row.ID,
		ProductID: row.ProductID,
		Slug:      row.Slug,
		Name:      row.Name,
		Price:     row.Price,
		IsActive:  row.IsActive,
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
		DeletedAt: dbutils.TimeToPtr(row.DeletedAt),
	}

	return variant, nil
}

func (vs *VariantStore) Create(ctx context.Context, input *service.CreateVariantInput) (*data.Variant, error) {
	params := sqlstore.InsertVariantParams{
		ProductID: input.ProductID,
		Slug:      input.Slug,
		Name:      input.Name,
		Price:     input.Price,
	}

	if input.IsActive != nil {
		params.IsActive = *input.IsActive
	}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	row, err := vs.Q.InsertVariant(ctx, params)
	if err != nil {
		return nil, err
	}

	variant := &data.Variant{
		ID:        row.ID,
		ProductID: row.ProductID,
		Name:      row.Name,
		Slug:      row.Slug,
		Price:     row.Price,
		IsActive:  row.IsActive,
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
	}

	return variant, nil
}

func (vs *VariantStore) Update(ctx context.Context, id int64, input *service.UpdateVariantInput) (*data.Variant, error) {
	if id < 1 {
		return nil, data.ErrRecordNotFound
	}

	params := sqlstore.UpdateVariantParams{
		ID:       id,
		Name:     dbutils.PtrToString(input.Name),
		IsActive: dbutils.PtrToBool(input.IsActive),
		Price:    dbutils.DecimalPtrToNumeric(input.Price),
	}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	row, err := vs.Q.UpdateVariant(ctx, params)
	if err != nil {
		return nil, err
	}

	variant := &data.Variant{
		ID:        row.ID,
		ProductID: row.ProductID,
		Name:      row.Name,
		Slug:      row.Slug,
		Price:     row.Price,
		IsActive:  row.IsActive,
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
	}

	return variant, nil
}

func (vs *VariantStore) Delete(ctx context.Context, id int64) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	err := vs.Q.DeleteVariant(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (vs *VariantStore) ListByProduct(ctx context.Context, productID int64) ([]*data.Variant, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	rows, err := vs.Q.ListProductVariants(ctx, productID)
	if err != nil {
		return nil, err
	}

	variants := []*data.Variant{}
	for _, row := range rows {
		variant := &data.Variant{
			ID:        row.ID,
			ProductID: row.ProductID,
			Slug:      row.Slug,
			Name:      row.Name,
			Price:     row.Price,
			IsActive:  row.IsActive,
			CreatedAt: row.CreatedAt.Time,
			UpdatedAt: row.UpdatedAt.Time,
			DeletedAt: dbutils.TimeToPtr(row.DeletedAt),
		}

		variants = append(variants, variant)
	}

	return variants, nil
}
