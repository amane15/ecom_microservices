package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type CartModel struct {
	queries *Queries
}

func (m CartModel) Get(id int64) (*Cart, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	cart, err := m.queries.GetCart(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &cart, nil
}

func (m CartModel) Insert(userID int64) (*Cart, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	cart, err := m.queries.InsertCart(ctx, sql.NullInt64{Int64: userID, Valid: true})
	if err != nil {
		return nil, err
	}

	return &cart, nil
}

func (m CartModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.queries.DeleteCart(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
