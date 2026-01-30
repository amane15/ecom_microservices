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

type CartStore struct {
	DBPool *pgxpool.Pool
	Q      *sqlstore.Queries
}

func NewCartStore(dbpool *pgxpool.Pool) *CartStore {
	return &CartStore{
		DBPool: dbpool,
		Q:      sqlstore.New(dbpool),
	}
}

func (cs *CartStore) Get(ctx context.Context, id int64) (*data.Cart, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	row, err := cs.Q.GetCart(ctx, id)
	if err != nil {
		return nil, checkAndHandlePostgresErrors(err)
	}

	cart := &data.Cart{
		ID:        row.ID,
		UserID:    row.UserID.Int64,
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
	}

	return cart, nil
}

func (cs *CartStore) Create(ctx context.Context, input *service.CreateCartInput) (*data.Cart, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	row, err := cs.Q.InsertCart(ctx, dbutils.PtrToInt8(&input.UserID))
	if err != nil {
		return nil, err
	}

	cart := &data.Cart{
		ID:        row.ID,
		UserID:    row.UserID.Int64,
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
	}

	return cart, nil
}

func (cs *CartStore) Delete(ctx context.Context, id int64) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	err := cs.Q.DeleteCart(ctx, id)
	if err != nil {
		return checkAndHandlePostgresErrors(err)
	}

	return nil
}
