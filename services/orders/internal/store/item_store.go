package store

import (
	"context"

	"github.com/amane15/ecom_microservice/services/orders/internal/data"
	"github.com/amane15/ecom_microservice/services/orders/internal/service"
	"github.com/amane15/ecom_microservice/services/orders/internal/store/sqlstore"
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

func (is *ItemStore) Get(ctx context.Context, id int64) (*data.OrderItem, error) {
	return nil, nil
}

func (is *ItemStore) Create(ctx context.Context, input *service.CreateOrderItemInput) (*data.OrderItem, error) {
	return nil, nil
}

func (is *ItemStore) Delete(ctx context.Context, id int64) error {
	return nil
}
