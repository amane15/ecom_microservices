package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
)

type CartItemModel struct {
	queries *Queries
}

func (m CartItemModel) Get(id int64) (*CartItem, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	item, err := m.queries.GetCartItem(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &item, nil
}

func (m CartItemModel) Insert(item *CartItem) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	itemRow, err := m.queries.InsertCartItem(ctx, InsertCartItemParams{
		CartID:    item.CartID,
		ProductID: item.ProductID,
		VariantID: item.VariantID,
		Quantity:  item.Quantity,
	})
	if err != nil {
		if pqError, ok := err.(*pq.Error); ok {
			switch pqError.Code {
			case "23505":
				return ErrItemAlreadyExists
			}
		}
		return err
	}

	item.ID = itemRow.ID
	item.Quantity = itemRow.Quantity
	item.CreatedAt = itemRow.CreatedAt
	item.UpdatedAt = itemRow.UpdatedAt

	return nil
}

func (m CartItemModel) Update(id int64, item *CartItem) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	itemRow, err := m.queries.UpdateCartItem(ctx, UpdateCartItemParams{
		CartID:    item.CartID,
		ProductID: item.ProductID,
		VariantID: item.VariantID,
		Quantity:  item.Quantity,
	})
	if err != nil {
		return err
	}

	item = &itemRow

	return nil
}

func (m CartItemModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.queries.DeleteCartItem(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
