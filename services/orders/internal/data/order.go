package data

import (
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
