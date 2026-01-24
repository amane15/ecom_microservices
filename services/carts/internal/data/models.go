package data

import "database/sql"

type Models struct {
	Carts CartModel
	Items CartItemModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Carts: CartModel{DB: db},
		Items: CartItemModel{DB: db},
	}
}
