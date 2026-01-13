package main

import "net/http"

func (app *application) createProductHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		SKU              string  `json:"sku"`
		Name             string  `json:"name"`
		Slug             string  `json:"slug"`
		ShortDescription *string `json:"short_description"`
		Description      *string `json:"description"`
		MetaTitle        *string `json:"meta_title"`
		MetaDescription  *string `json:"meta_description"`
		IsActive         *bool   `json:"is_active"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
}

func (app *application) getProductHandler(w http.ResponseWriter, r *http.Request) {}

func (app *application) updateProductHandler(w http.ResponseWriter, r *http.Request) {}
