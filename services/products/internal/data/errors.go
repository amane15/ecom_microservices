package data

import (
	"errors"

	"github.com/lib/pq"
)

var (
	ErrDuplicateSlug    = errors.New("slug already exists")
	ErrDuplicateSku     = errors.New("sku already exists")
	ErrRecordNotFound   = errors.New("record not found")
	ErrNoFieldsToUpdate = errors.New("no fields to update")
)

func handleUniqueViolationError(pqError *pq.Error) error {
	switch pqError.Constraint {
	case "products_sku_key":
		return ErrDuplicateSku
	case "products_slug_key":
		return ErrDuplicateSlug
	case "categories_slug_key":
		return ErrDuplicateSlug
	default:
		return pqError
	}
}
