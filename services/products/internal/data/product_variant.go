package data

import (
	"database/sql"
	"time"
)

// id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
// product_id BIGINT,
//
// slug VARCHAR(128) NOT NULL UNIQUE,
// name VARCHAR(255),
//
// price DECIMAL(6, 2),
//
// is_active BOOLEAN NOT NULL DEFAULT true,
//
// created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
// updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
// deleted_at TIMESTAMPTZ
type ProductVariant struct {
	ID        int64   `json:"id"`
	ProductID int64   `json:"product_id"`
	Slug      string  `json:"slug"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	IsActive  bool    `json:"is_active"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

type ProductVariantModel struct {
	DB *sql.DB
}
