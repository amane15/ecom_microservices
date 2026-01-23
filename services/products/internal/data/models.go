package data

import "database/sql"

type Models struct {
	Products   ProductModel
	Variants   ProductVariantModel
	Categories CategoryModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Products:   ProductModel{DB: db},
		Variants:   ProductVariantModel{DB: db},
		Categories: CategoryModel{DB: db},
	}
}
