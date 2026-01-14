package main

import (
	"errors"
	"net/http"
	"unicode/utf8"

	"github.com/amane15/ecom_microservice/services/proudcts/internal/data"
	"github.com/amane15/ecom_microservice/services/proudcts/internal/validator"
)

func (app *application) createProductHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Sku              string  `json:"sku"`
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

	product := &data.Product{
		Sku:  input.Sku,
		Name: input.Name,
		Slug: input.Slug,
	}
	v := validator.New()

	if input.ShortDescription != nil {
		product.ShortDescription = *input.ShortDescription
	}
	if input.Description != nil {
		product.Description = *input.Description
	}

	if input.MetaTitle != nil {
		product.MetaTitle = *input.MetaTitle
		v.Check(utf8.RuneCountInString(*input.MetaTitle) <= 255, "meta_title", "must not be more than 255 characters long")
	}
	if input.MetaDescription != nil {
		product.MetaDescription = *input.MetaDescription
	}

	if input.IsActive != nil {
		product.IsActive = *input.IsActive
	} else {
		product.IsActive = true
	}

	if data.ValidateProduct(v, product); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.products.Create(product)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateSku):
			app.badRequestResponse(w, r, err)
		case errors.Is(err, data.ErrDuplicateSlug):
			app.badRequestResponse(w, r, err)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"product": product}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getProductHandler(w http.ResponseWriter, r *http.Request) {}

func (app *application) updateProductHandler(w http.ResponseWriter, r *http.Request) {}
