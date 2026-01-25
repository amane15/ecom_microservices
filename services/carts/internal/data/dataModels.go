package data

import (
	"database/sql"
)

type Models struct {
	Carts CartModel
	Items CartItemModel
}

func NewModels(db *sql.DB) Models {
	queries := New(db)
	return Models{
		Carts: CartModel{queries: queries},
		Items: CartItemModel{queries: queries},
	}
}
