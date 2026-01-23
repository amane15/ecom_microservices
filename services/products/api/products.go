package main

import (
	"errors"
	"net/http"
	"unicode/utf8"

	"github.com/amane15/ecom_microservice/pkg/validator"
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
	productInput := data.UpdateProductInput{}

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
	app.logger.Info("Change hit")
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	var input struct {
		Status *data.ProductStatus `json:"status"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	if input.Status == nil {
		v.AddError("status", "status must be provided")
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	updateInput := &data.UpdateProductInput{}
	status := *input.Status

	switch status {
	case data.ProductStatusDraft, data.ProductStatusArchived, data.ProductStatusActive:
		if status == data.ProductStatusActive {
			count, err := app.variants.GetVariantCountForProduct(id)
			if err != nil {
				switch {
				case errors.Is(err, data.ErrRecordNotFound):
					count = 0
				default:
					app.serverErrorResponse(w, r, err)
					return
				}
			}

			if status == data.ProductStatusActive && count == 0 {
				app.badRequestResponse(w, r, errors.New("to make a product active you need have at least 1 variant"))
				return
			}
		}
		updateInput.Status = &status
	default:
		v.AddError("status", "must be a valid status")
		app.failedValidationResponse(w, r, v.Errors)
		return

	}

	product, err := app.products.Update(id, updateInput)
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

func (app *application) setDefaultVariantHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	var input struct {
		DefaultVariantID *int64 `json:"default_variant_id"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	if input.DefaultVariantID == nil {
		v.AddError("default_variant_id", "must be provided")
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	variantID := *input.DefaultVariantID
	exists, err := app.variants.IsVariantExists(variantID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	if !exists {
		app.badRequestResponse(w, r, errors.New("variant does not exists"))
		return
	}

	updateInput := &data.UpdateProductInput{}
	updateInput.DefaultVariantID = input.DefaultVariantID

	product, err := app.products.Update(id, updateInput)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.badRequestResponse(w, r, errors.New("product does not exists"))
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

func (app *application) listProductsHandler(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	limit := app.readInt(queryParams, "limit", 10)
	offset := app.readInt(queryParams, "offset", 0)

	products, err := app.products.GetAll(limit, offset)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"products": products}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listProductVariantsHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	variants, err := app.variants.GetVariantsByProduct(id)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"variants": variants}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
