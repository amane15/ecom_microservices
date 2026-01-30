package data

import (
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
