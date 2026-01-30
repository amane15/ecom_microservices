package data

import "errors"

var (
	ErrRecordNotFound    = errors.New("record not found")
	ErrNoFieldsToUpdate  = errors.New("no fields to update")
	ErrItemAlreadyExists = errors.New("item already exists")
)

type ValidationError struct {
	Fields map[string]string
}

func (v ValidationError) Error() string {
	return "validation errors"
}
