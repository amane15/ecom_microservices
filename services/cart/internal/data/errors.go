package data

import "errors"

var (
	ErrRecordNotFound    = errors.New("record not found")
	ErrNoFieldsToUpdate  = errors.New("no fields to update")
	ErrItemAlreadyExists = errors.New("item already exists")
)
