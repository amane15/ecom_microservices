package data

import "database/sql"

type Models struct {
	Products   ProductModel
	Variants   ProductVariantModel
	Categories CategoryModel
}

func NewModels(db *sql.DB) Models {
	queries := New(db)
	return Models{
		Products:   ProductModel{queries: queries, DB: db},
		Variants:   ProductVariantModel{queries: queries, DB: db},
		Categories: CategoryModel{queries: queries},
	}
}
