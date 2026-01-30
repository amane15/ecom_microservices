package store

import (
	"context"
	"time"

	"github.com/amane15/ecom_microservice/services/orders/internal/data"
	"github.com/amane15/ecom_microservice/services/orders/internal/service"
	"github.com/amane15/ecom_microservice/services/orders/internal/store/sqlstore"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderStore struct {
	DBPool *pgxpool.Pool
	Q      *sqlstore.Queries
}

func NewOrderStore(dbpool *pgxpool.Pool) *OrderStore {
	return &OrderStore{
		DBPool: dbpool,
		Q:      sqlstore.New(dbpool),
	}
}

func (os *OrderStore) Get(ctx context.Context, id int64) (*data.Order, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	row, err := os.Q.GetOrder(ctx, id)
	if err != nil {
		return nil, checkAndHandlePostgresErrors(err)
	}

	order := &data.Order{
		ID:             row.ID,
		UserID:         row.UserID,
		Status:         data.OrderStatus(row.Status),
		SubtotalAmount: row.SubtotalAmount,
		TaxAmount:      row.TaxAmount,
		ShippingAmount: row.ShippingAmount,
		DiscountAmount: row.DiscountAmount,
		TotalAmount:    row.TaxAmount,
		CreatedAt:      row.CreatedAt.Time,
		UpdatedAt:      row.UpdatedAt.Time,
	}

	return order, nil
}

func (os *OrderStore) Create(ctx context.Context, input *service.CreateOrderInput) (*data.Order, error) {
	return nil, nil
}

func (os *OrderStore) Delete(ctx context.Context, id int64) error {
	return nil
}
