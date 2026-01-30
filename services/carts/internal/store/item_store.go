package store

import (
	"context"
	"time"

	"github.com/amane15/ecom_microservice/internal/dbutils"
	"github.com/amane15/ecom_microservice/services/carts/internal/data"
	"github.com/amane15/ecom_microservice/services/carts/internal/service"
	"github.com/amane15/ecom_microservice/services/carts/internal/store/sqlstore"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ItemStore struct {
	DBPool *pgxpool.Pool
	Q      *sqlstore.Queries
}

func NewItemStore(dbpool *pgxpool.Pool) *ItemStore {
	return &ItemStore{
		DBPool: dbpool,
		Q:      sqlstore.New(dbpool),
	}
}

func (is *ItemStore) Get(ctx context.Context, id int64) (*data.Item, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	row, err := is.Q.GetCartItem(ctx, id)
	if err != nil {
		return nil, checkAndHandlePostgresErrors(err)
	}

	item := &data.Item{
		ID:        row.ID,
		CartID:    row.CartID,
		ProductID: row.ProductID.Int64,
		VariantID: row.VariantID.Int64,
		Quantity:  int(row.Quantity),
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
	}

	return item, nil
}

func (is *ItemStore) Create(ctx context.Context, input *service.CreateItemInput) (*data.Item, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	row, err := is.Q.InsertCartItem(ctx, sqlstore.InsertCartItemParams{
		CartID:    input.CartID,
		ProductID: dbutils.PtrToInt8(input.ProductID),
		VariantID: dbutils.PtrToInt8(input.VariantID),
		Quantity:  *input.Quantity,
	})
	if err != nil {
		return nil, checkAndHandlePostgresErrors(err)
	}

	item := &data.Item{
		ID:        row.ID,
		CartID:    row.CartID,
		ProductID: row.ProductID.Int64,
		VariantID: row.VariantID.Int64,
		Quantity:  int(row.Quantity),
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
	}

	return item, nil
}

func (is *ItemStore) Delete(ctx context.Context, id int64) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	err := is.Q.DeleteCartItem(ctx, id)
	if err != nil {
		return checkAndHandlePostgresErrors(err)
	}

	return nil
}
