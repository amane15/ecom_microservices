package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/shopspring/decimal"
)

type OrderItem struct {
	ID          int64           `json:"id"`
	OrderID     int64           `json:"order_id"`
	ProductID   int64           `json:"product_id"`
	VariantID   int64           `json:"variant_id"`
	ProductName string          `json:"product_name"`
	VariantName string          `json:"variant_name"`
	Sku         string          `json:"sku"`
	UnitPrice   decimal.Decimal `json:"unit_price"`
	Quantity    decimal.Decimal `json:"quantity"`
	TotalPrice  decimal.Decimal `json:"total_price"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

type OrderItemModel struct {
	DB *sql.DB
}

func (m OrderItemModel) Get(id int64) (*OrderItem, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
	SELECT id, order_id, product_id, variant_id, product_name,
		variant_name, sku, unit_price, quantity, total_price,
		created_at, updated_at
	FROM order_items
	WHERE id = $1
	`

	item := &OrderItem{}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&item.ID,
		&item.OrderID,
		&item.ProductID,
		&item.VariantID,
		&item.ProductName,
		&item.VariantName,
		&item.Sku,
		&item.UnitPrice,
		&item.Quantity,
		&item.TotalPrice,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return item, nil
}

func (m OrderItemModel) Insert(item *OrderItem) error {
	query := `
	INSERT INTO order_items (order_id, product_id, variant_id
			product_name, variant_name, sku, unit_price, quantity,
			total_price)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	RETURNING id, created_at, updated_at`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{
		item.OrderID, item.ProductID, item.VariantID, item.ProductName,
		item.VariantName, item.Sku, item.UnitPrice, item.Quantity, item.TotalPrice,
	}
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(
		&item.ID,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}

func (m OrderItemModel) Update(item *OrderItem) error {
	return nil
}

func (m OrderItemModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}
	query := `DELETE FROM carts WHERE id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}
