package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/amane15/ecom_microservice/internal/dbutils"
	"github.com/amane15/ecom_microservice/services/products/internal/data"
	"github.com/amane15/ecom_microservice/services/products/internal/service"
	"github.com/amane15/ecom_microservice/services/products/internal/store/sqlstore"
)

type VariantStore struct {
	DB *sql.DB
	Q  *sqlstore.Queries
}

func NewVariantStore(db *sql.DB) *VariantStore {
	return &VariantStore{
		DB: db,
		Q:  sqlstore.New(db),
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
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
		DeletedAt: dbutils.MapNullTime(row.DeletedAt),
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
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}

	return variant, nil
}

func (vs *VariantStore) Update(ctx context.Context, id int64, input *service.UpdateVariantInput) (*data.Variant, error) {
	if id < 1 {
		return nil, data.ErrRecordNotFound
	}

	params := sqlstore.UpdateVariantParams{
		ID:       id,
		Name:     dbutils.MapStringPtr(input.Name),
		IsActive: dbutils.MapBoolPtr(input.IsActive),
	}

	if input.Price != nil {
		params.Price = sql.NullString{String: input.Price.String(), Valid: true}
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
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
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
			CreatedAt: row.CreatedAt,
			UpdatedAt: row.UpdatedAt,
			DeletedAt: dbutils.MapNullTime(row.DeletedAt),
		}

		variants = append(variants, variant)
	}

	return variants, nil
}
