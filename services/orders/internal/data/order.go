package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/shopspring/decimal"
)

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusPaid      OrderStatus = "paid"
	OrderStatusCancelled OrderStatus = "cancelled"
	OrderStatusShipped   OrderStatus = "shipped"
	OrderStatusDelivered OrderStatus = "delivered"
	OrderStatusRefunded  OrderStatus = "refunded"
)

type Order struct {
	ID             int64           `json:"id"`
	UserID         int64           `json:"user_id"`
	Status         OrderStatus     `json:"status"`
	SubtotalAmount decimal.Decimal `json:"subtotal_amount"`
	TaxAmount      decimal.Decimal `json:"tax_amount"`
	ShippingAmount decimal.Decimal `json:"shipping_amount"`
	DiscountAmount decimal.Decimal `json:"discount_amount"`
	TotalAmount    decimal.Decimal `json:"total_amount"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
}

type OrderModel struct {
	DB *sql.DB
}

func (m OrderModel) Get(id int64) (*Order, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
	SELECT id, user_id, status, subtotal_amount, tax_amount, shipping_amount,
		discount_amount, total_amount, created_at, updated_at
	FROM orders 
	WHERE id = $1
	`
	order := &Order{}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&order.ID,
		&order.UserID,
		&order.Status,
		&order.SubtotalAmount,
		&order.TaxAmount,
		&order.ShippingAmount,
		&order.DiscountAmount,
		&order.TotalAmount,
		&order.CreatedAt,
		&order.UpdatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return order, nil
}

func (m OrderModel) Insert(order *Order) error {
	query := `
	INSERT INTO orders (user_id, status, subtotal_amount, tax_amount, shipping_amount,
		discount_amount, total_amount)
	VALUES($1, $2, $3, $4, $5, $6, $7)
	RETURNING id, created_at, updated_at
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{
		order.UserID, order.Status, order.SubtotalAmount, order.TaxAmount,
		order.ShippingAmount, order.DiscountAmount, order.TotalAmount,
	}
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(
		&order.ID,
		&order.CreatedAt,
		&order.UpdatedAt,
	)
	if err != nil {
		return nil
	}

	return nil
}

// Just for status update
// func (m OrderModel) Update(order *Order) error {}
