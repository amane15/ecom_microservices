package data

import (
	"errors"
)

var (
	ErrDuplicateSlug       = errors.New("slug already exists")
	ErrRecordNotFound      = errors.New("record not found")
	ErrNoFieldsToUpdate    = errors.New("no fields to update")
	ErrTableNotFound       = errors.New("table not found")
	ErrColumnNotFound      = errors.New("column not found")
	ErrReferenceNotFound   = errors.New("referenced row not found")
	ErrCheckViolation      = errors.New("check constraint failed")
	ErrForeignKeyViolation = errors.New("row still in use")
	ErrNotNullViolation    = errors.New("value cannot be null")
)
