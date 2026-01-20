package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
)

type CartItem struct {
	ID        int64 `json:"id"`
	CartID    int64 `json:"cart_id"`
	ProductID int64 `json:"product_id"`
	VariantID int64 `json:"variant_id"`
	Quantity  int   `json:"quantity"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CartItemModel struct {
	DB *sql.DB
}

func (m CartItemModel) Get(id int64) (*CartItem, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
	SELECT id, cart_id, product_id, variant_id, quantity,
		created_at, updated_at
	FROM cart_items
	WHERE id = $1
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	cartItem := &CartItem{}
	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&cartItem.ID,
		&cartItem.CartID,
		&cartItem.ProductID,
		&cartItem.VariantID,
		&cartItem.Quantity,
		&cartItem.CreatedAt,
		&cartItem.UpdatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return cartItem, nil
}

func (m CartItemModel) Insert(item *CartItem) error {
	query := `INSERT INTO cart_items(cart_id, product_id, variant_id, quantity)
	VALUES ($1, $2,$3, $4)
	ON CONFLICT (cart_id, variant_id)
	DO UPDATE
	SET quantity = cart_items.quantity + EXCLUDED.quantity, updated_at = now()
	RETURNING id, quantity, created_at, updated_at`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{&item.CartID, &item.ProductID, &item.VariantID, &item.Quantity}
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(
		&item.ID,
		&item.Quantity,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if err != nil {
		if pqError, ok := err.(*pq.Error); ok {
			switch pqError.Code {
			case "23505":
				return ErrItemAlreadyExists
			}
		}
		return err
	}

	return nil
}

func (m CartItemModel) Update(id int64, item *CartItem) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
	INSERT INTO cart_items (cart_id, product_id, variant_id, quantity)
	VALUES ($1, $2, $3, $4)
	ON CONFLICT (cart_id, variant_id)
	DO UPDATE 
	SET quantity = quantity + EXCLUDED.quantity, updated_at = now()
	RETURNING id, cart_id, product_id, variant_id, quantity, created_at, updated_at
	`
	args := []any{&item.CartID, &item.ProductID, &item.VariantID, &item.Quantity}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(
		&item.ID,
		&item.CartID,
		&item.ProductID,
		&item.VariantID,
		&item.Quantity,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (m CartItemModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}
	query := `DELETE FROM cart_items WHERE id = $1`
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
