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
		return nil, err
	}

	product := &data.Product{
		ID:               row.ID,
		Name:             row.Name,
		Slug:             row.Slug,
		Status:           data.ProductStatus(row.Status),
		Description:      dbutils.MapNullString(row.Description),
		ShortDescription: dbutils.MapNullString(row.ShortDescription),
		MetaTitle:        dbutils.MapNullString(row.MetaTitle),
		MetaDescription:  dbutils.MapNullString(row.MetaDescription),
		DefaultVariantID: dbutils.MapNullInt64(row.DefaultVariantID),
		CreatedAt:        row.CreatedAt,
		UpdatedAt:        row.UpdatedAt,
		DeletedAt:        dbutils.MapNullTime(row.DeletedAt),
	}

	return product, nil
}

func (ps *ProductStore) Create(ctx context.Context, input *service.CreateProductInput) (*data.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	row, err := ps.Q.InsertProduct(ctx, sqlstore.InsertProductParams{
		Name:             input.Name,
		Slug:             input.Slug,
		Description:      dbutils.MapStringPtr(input.Description),
		ShortDescription: dbutils.MapStringPtr(input.ShortDescription),
		MetaTitle:        dbutils.MapStringPtr(input.MetaDescription),
		MetaDescription:  dbutils.MapStringPtr(input.MetaDescription),
		Status:           sqlstore.ProductStatus(*input.Status),
	})
	if err != nil {
		return nil, err
	}

	product := &data.Product{
		ID:               row.ID,
		Name:             row.Name,
		Slug:             row.Slug,
		Status:           data.ProductStatus(row.Status),
		Description:      dbutils.MapNullString(row.Description),
		ShortDescription: dbutils.MapNullString(row.ShortDescription),
		MetaTitle:        dbutils.MapNullString(row.MetaTitle),
		MetaDescription:  dbutils.MapNullString(row.MetaDescription),
		DefaultVariantID: dbutils.MapNullInt64(row.DefaultVariantID),
		CreatedAt:        row.CreatedAt,
		UpdatedAt:        row.UpdatedAt,
	}
	return product, nil
}

func (ps *ProductStore) Update(ctx context.Context, id int64, input *service.UpdateProductInput) (*data.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	params := sqlstore.UpdateProductParams{
		ID:               id,
		Name:             dbutils.MapStringPtr(input.Name),
		Description:      dbutils.MapStringPtr(input.Description),
		ShortDescription: dbutils.MapStringPtr(input.ShortDescription),
		MetaTitle:        dbutils.MapStringPtr(input.MetaTitle),
		MetaDescription:  dbutils.MapStringPtr(input.MetaDescription),
		DefaultVariantID: dbutils.MapInt64Ptr(input.DefaultVariantID),
	}

	if input.Status != nil {
		params.Status = sqlstore.NullProductStatus{ProductStatus: sqlstore.ProductStatus(*input.Status), Valid: true}
	}

	row, err := ps.Q.UpdateProduct(ctx, params)
	if err != nil {
		return nil, err
	}

	product := &data.Product{
		ID:               row.ID,
		Name:             row.Name,
		Slug:             row.Slug,
		Status:           data.ProductStatus(row.Status),
		Description:      dbutils.MapNullString(row.Description),
		ShortDescription: dbutils.MapNullString(row.ShortDescription),
		MetaTitle:        dbutils.MapNullString(row.MetaTitle),
		MetaDescription:  dbutils.MapNullString(row.MetaDescription),
		DefaultVariantID: dbutils.MapNullInt64(row.DefaultVariantID),
		CreatedAt:        row.CreatedAt,
		UpdatedAt:        row.UpdatedAt,
	}
	return product, nil
}

func (ps *ProductStore) List(ctx context.Context, limit, offset int32) ([]*data.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	rows, err := ps.Q.ListProducts(ctx, sqlstore.ListProductsParams{Limit: limit, Offset: offset})
	if err != nil {
		return nil, err
	}

	products := []*data.Product{}
	for _, row := range rows {
		product := &data.Product{
			ID:               row.ID,
			Name:             row.Name,
			Slug:             row.Slug,
			Description:      dbutils.MapNullString(row.Description),
			ShortDescription: dbutils.MapNullString(row.ShortDescription),
			MetaTitle:        dbutils.MapNullString(row.MetaTitle),
			MetaDescription:  dbutils.MapNullString(row.MetaDescription),
			Status:           data.ProductStatus(row.Status),
			DefaultVariantID: dbutils.MapNullInt64(row.DefaultVariantID),
			CreatedAt:        row.CreatedAt,
			UpdatedAt:        row.UpdatedAt,
			DeletedAt:        dbutils.MapNullTime(row.DeletedAt),
		}

		products = append(products, product)
	}

	return products, nil
}

func (ps *ProductStore) ListByProduct(ctx context.Context, limit, offset int32) ([]*data.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	rows, err := ps.Q.ListProducts(ctx, sqlstore.ListProductsParams{Limit: limit, Offset: offset})
	if err != nil {
		return nil, err
	}

	products := []*data.Product{}

	for _, p := range rows {
		product := &data.Product{
			ID:               p.ID,
			Name:             p.Name,
			Slug:             p.Slug,
			Status:           data.ProductStatus(p.Status),
			Description:      dbutils.MapNullString(p.Description),
			ShortDescription: dbutils.MapNullString(p.ShortDescription),
			MetaTitle:        dbutils.MapNullString(p.MetaTitle),
			MetaDescription:  dbutils.MapNullString(p.MetaDescription),
			DefaultVariantID: dbutils.MapNullInt64(p.DefaultVariantID),
			CreatedAt:        p.CreatedAt,
			UpdatedAt:        p.UpdatedAt,
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
		return err
	}

	return nil
}
