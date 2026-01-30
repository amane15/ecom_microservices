package data

import "time"

type Cart struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Item struct {
	ID        int64 `json:"id"`
	CartID    int64 `json:"cart_id"`
	ProductID int64 `json:"product_id"`
	VariantID int64 `json:"variant_id"`
	Quantity  int   `json:"quantity"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
