package data

import "database/sql"

type Models struct {
	Orders OrderModel
	Items  OrderItemModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Orders: OrderModel{DB: db},
		Items:  OrderItemModel{DB: db},
	}
}
