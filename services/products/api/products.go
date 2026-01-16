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
		Name             string  `json:"name"`
		Slug             string  `json:"slug"`
		ShortDescription *string `json:"short_description"`
		Description      *string `json:"description"`
		MetaTitle        *string `json:"meta_title"`
		MetaDescription  *string `json:"meta_description"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	product := &data.Product{
		Name:   input.Name,
		Slug:   input.Slug,
		Status: data.ProductStatusDraft,
	}

	v := validator.New()

	if input.ShortDescription != nil {
		product.ShortDescription = input.ShortDescription
	}
	if input.Description != nil {
		product.Description = input.Description
	}

	if input.MetaTitle != nil {
		product.MetaTitle = input.MetaTitle
		v.Check(utf8.RuneCountInString(*input.MetaTitle) <= 255, "meta_title", "must not be more than 255 characters long")
	}
	if input.MetaDescription != nil {
		product.MetaDescription = input.MetaDescription
	}

	if data.ValidateProduct(v, product); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.products.Create(product)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateSlug):
			app.badRequestResponse(w, r, errors.New("product already exists"))
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

func (app *application) getProductHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	product, err := app.products.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"product": product}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateProductHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	productInput := data.UpdateProductRow{}

	err = app.readJSON(w, r, &productInput)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if productInput.Name == nil && productInput.Description == nil &&
		productInput.ShortDescription == nil && productInput.MetaTitle == nil &&
		productInput.MetaDescription == nil {
		app.badRequestResponse(w, r, errors.New("at least one field must be provided"))
		return
	}

	v := validator.New()
	if data.ValidateUpdateProduct(v, &productInput); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	product, err := app.products.Update(id, &productInput)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrNoFieldsToUpdate):
			app.badRequestResponse(w, r, errors.New("no fields were provided"))
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}

		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"product": product}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) changeProductStatusHandler(w http.ResponseWriter, r *http.Request) {
	_, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
}

func (app *application) setDefaultVariantHandler(w http.ResponseWriter, r *http.Request) {}
