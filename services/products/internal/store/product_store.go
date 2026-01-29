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

type ProductStore struct {
	DBPool *pgxpool.Pool
	Q      *sqlstore.Queries
}

func NewProductStore(dbpool *pgxpool.Pool) *ProductStore {
	return &ProductStore{
		DBPool: dbpool,
		Q:      sqlstore.New(dbpool),
	}
}

func (ps *ProductStore) Get(ctx context.Context, id int64) (*data.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	row, err := ps.Q.GetProduct(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, data.ErrRecordNotFound
		default:
			return nil, checkAndHandlePostgresErrors(err)
		}
	}

	product := &data.Product{
		ID:               row.ID,
		Name:             row.Name,
		Slug:             row.Slug,
		Status:           data.ProductStatus(row.Status),
		Description:      dbutils.StringToPtr(row.Description),
		ShortDescription: dbutils.StringToPtr(row.ShortDescription),
		MetaTitle:        dbutils.StringToPtr(row.MetaTitle),
		MetaDescription:  dbutils.StringToPtr(row.MetaDescription),
		DefaultVariantID: dbutils.Int8ToPtr(row.DefaultVariantID),
		CreatedAt:        row.CreatedAt.Time,
		UpdatedAt:        row.UpdatedAt.Time,
		DeletedAt:        dbutils.TimeToPtr(row.DeletedAt),
	}

	return product, nil
}

func (ps *ProductStore) Create(ctx context.Context, input *service.CreateProductInput) (*data.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	row, err := ps.Q.InsertProduct(ctx, sqlstore.InsertProductParams{
		Name:             input.Name,
		Slug:             input.Slug,
		Description:      dbutils.PtrToString(input.Description),
		ShortDescription: dbutils.PtrToString(input.ShortDescription),
		MetaTitle:        dbutils.PtrToString(input.MetaDescription),
		MetaDescription:  dbutils.PtrToString(input.MetaDescription),
		Status:           sqlstore.ProductStatus(*input.Status),
	})
	if err != nil {
		return nil, checkAndHandlePostgresErrors(err)
	}

	product := &data.Product{
		ID:               row.ID,
		Name:             row.Name,
		Slug:             row.Slug,
		Status:           data.ProductStatus(row.Status),
		Description:      dbutils.StringToPtr(row.Description),
		ShortDescription: dbutils.StringToPtr(row.ShortDescription),
		MetaTitle:        dbutils.StringToPtr(row.MetaTitle),
		MetaDescription:  dbutils.StringToPtr(row.MetaDescription),
		DefaultVariantID: dbutils.Int8ToPtr(row.DefaultVariantID),
		CreatedAt:        row.CreatedAt.Time,
		UpdatedAt:        row.UpdatedAt.Time,
	}
	return product, nil
}

func (ps *ProductStore) Update(ctx context.Context, id int64, input *service.UpdateProductInput) (*data.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	params := sqlstore.UpdateProductParams{
		ID:               id,
		Name:             dbutils.PtrToString(input.Name),
		Description:      dbutils.PtrToString(input.Description),
		ShortDescription: dbutils.PtrToString(input.ShortDescription),
		MetaTitle:        dbutils.PtrToString(input.MetaTitle),
		MetaDescription:  dbutils.PtrToString(input.MetaDescription),
		DefaultVariantID: dbutils.PtrToInt8(input.DefaultVariantID),
	}

	if input.Status != nil {
		params.Status = sqlstore.NullProductStatus{ProductStatus: sqlstore.ProductStatus(*input.Status), Valid: true}
	}

	row, err := ps.Q.UpdateProduct(ctx, params)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, data.ErrRecordNotFound
		default:
			return nil, checkAndHandlePostgresErrors(err)
		}
	}

	product := &data.Product{
		ID:               row.ID,
		Name:             row.Name,
		Slug:             row.Slug,
		Status:           data.ProductStatus(row.Status),
		Description:      dbutils.StringToPtr(row.Description),
		ShortDescription: dbutils.StringToPtr(row.ShortDescription),
		MetaTitle:        dbutils.StringToPtr(row.MetaTitle),
		MetaDescription:  dbutils.StringToPtr(row.MetaDescription),
		DefaultVariantID: dbutils.Int8ToPtr(row.DefaultVariantID),
		CreatedAt:        row.CreatedAt.Time,
		UpdatedAt:        row.UpdatedAt.Time,
	}
	return product, nil
}

func (ps *ProductStore) List(ctx context.Context, limit, offset int32) ([]*data.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	rows, err := ps.Q.ListProducts(ctx, sqlstore.ListProductsParams{Limit: limit, Offset: offset})
	if err != nil {
		return nil, checkAndHandlePostgresErrors(err)
	}

	products := []*data.Product{}
	for _, row := range rows {
		product := &data.Product{
			ID:               row.ID,
			Name:             row.Name,
			Slug:             row.Slug,
			Description:      dbutils.StringToPtr(row.Description),
			ShortDescription: dbutils.StringToPtr(row.ShortDescription),
			MetaTitle:        dbutils.StringToPtr(row.MetaTitle),
			MetaDescription:  dbutils.StringToPtr(row.MetaDescription),
			Status:           data.ProductStatus(row.Status),
			DefaultVariantID: dbutils.Int8ToPtr(row.DefaultVariantID),
			CreatedAt:        row.CreatedAt.Time,
			UpdatedAt:        row.UpdatedAt.Time,
			DeletedAt:        dbutils.TimeToPtr(row.DeletedAt),
		}

		products = append(products, product)
	}

	return products, nil
}

func (ps *ProductStore) Delete(ctx context.Context, id int64) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	err := ps.Q.DeleteProduct(ctx, id)
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
