package main

import (
	"errors"
	"net/http"
	"unicode/utf8"

	"github.com/amane15/ecom_microservice/pkg/httpx"
	"github.com/amane15/ecom_microservice/pkg/validator"
	"github.com/amane15/ecom_microservice/services/products/internal/data"
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

	err := httpx.ReadJSON(w, r, &input)
	if err != nil {
		app.httpErrRes.BadRequestResponse(w, r, err)
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
		app.httpErrRes.FailedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.products.Create(product)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateSlug):
			app.httpErrRes.BadRequestResponse(w, r, errors.New("product already exists"))
		default:
			app.httpErrRes.ServerErrorResponse(w, r, err)
		}
		return
	}

	err = httpx.WriteJSON(w, http.StatusCreated, httpx.Envelope{"product": product}, nil)
	if err != nil {
		app.httpErrRes.ServerErrorResponse(w, r, err)
	}
}

func (app *application) getProductHandler(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.ReadIDParam(r)
	if err != nil {
		app.httpErrRes.NotFoundResponse(w, r)
		return
	}

	product, err := app.products.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.httpErrRes.NotFoundResponse(w, r)
		default:
			app.httpErrRes.ServerErrorResponse(w, r, err)
		}
		return
	}

	err = httpx.WriteJSON(w, http.StatusOK, httpx.Envelope{"product": product}, nil)
	if err != nil {
		app.httpErrRes.ServerErrorResponse(w, r, err)
	}
}

func (app *application) updateProductHandler(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.ReadIDParam(r)
	if err != nil {
		app.httpErrRes.NotFoundResponse(w, r)
		return
	}
	productInput := data.UpdateProductInput{}

	err = httpx.ReadJSON(w, r, &productInput)
	if err != nil {
		app.httpErrRes.BadRequestResponse(w, r, err)
		return
	}

	if productInput.Name == nil && productInput.Description == nil &&
		productInput.ShortDescription == nil && productInput.MetaTitle == nil &&
		productInput.MetaDescription == nil {
		app.httpErrRes.BadRequestResponse(w, r, errors.New("at least one field must be provided"))
		return
	}

	v := validator.New()
	if data.ValidateUpdateProduct(v, &productInput); !v.Valid() {
		app.httpErrRes.FailedValidationResponse(w, r, v.Errors)
		return
	}

	product, err := app.products.Update(id, &productInput)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrNoFieldsToUpdate):
			app.httpErrRes.BadRequestResponse(w, r, errors.New("no fields were provided"))
		case errors.Is(err, data.ErrRecordNotFound):
			app.httpErrRes.NotFoundResponse(w, r)
		default:
			app.httpErrRes.ServerErrorResponse(w, r, err)
		}

		return
	}

	err = httpx.WriteJSON(w, http.StatusOK, httpx.Envelope{"product": product}, nil)
	if err != nil {
		app.httpErrRes.ServerErrorResponse(w, r, err)
	}
}

func (app *application) changeProductStatusHandler(w http.ResponseWriter, r *http.Request) {
	app.logger.Info("Change hit")
	id, err := httpx.ReadIDParam(r)
	if err != nil {
		app.httpErrRes.NotFoundResponse(w, r)
		return
	}

	var input struct {
		Status *data.ProductStatus `json:"status"`
	}

	err = httpx.ReadJSON(w, r, &input)
	if err != nil {
		app.httpErrRes.BadRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	if input.Status == nil {
		v.AddError("status", "status must be provided")
		app.httpErrRes.FailedValidationResponse(w, r, v.Errors)
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
					app.httpErrRes.ServerErrorResponse(w, r, err)
					return
				}
			}

			if status == data.ProductStatusActive && count == 0 {
				app.httpErrRes.BadRequestResponse(w, r, errors.New("to make a product active you need have at least 1 variant"))
				return
			}
		}
		updateInput.Status = &status
	default:
		v.AddError("status", "must be a valid status")
		app.httpErrRes.FailedValidationResponse(w, r, v.Errors)
		return

	}

	product, err := app.products.Update(id, updateInput)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.httpErrRes.NotFoundResponse(w, r)
		default:
			app.httpErrRes.ServerErrorResponse(w, r, err)
		}
		return

	}

	err = httpx.WriteJSON(w, http.StatusOK, httpx.Envelope{"product": product}, nil)
	if err != nil {
		app.httpErrRes.ServerErrorResponse(w, r, err)
	}
}

func (app *application) setDefaultVariantHandler(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.ReadIDParam(r)
	if err != nil {
		app.httpErrRes.NotFoundResponse(w, r)
		return
	}

	var input struct {
		DefaultVariantID *int64 `json:"default_variant_id"`
	}

	err = httpx.ReadJSON(w, r, &input)
	if err != nil {
		app.httpErrRes.BadRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	if input.DefaultVariantID == nil {
		v.AddError("default_variant_id", "must be provided")
		app.httpErrRes.FailedValidationResponse(w, r, v.Errors)
		return
	}

	variantID := *input.DefaultVariantID
	exists, err := app.variants.IsVariantExists(variantID)
	if err != nil {
		app.httpErrRes.ServerErrorResponse(w, r, err)
		return
	}
	if !exists {
		app.httpErrRes.BadRequestResponse(w, r, errors.New("variant does not exists"))
		return
	}

	updateInput := &data.UpdateProductInput{}
	updateInput.DefaultVariantID = input.DefaultVariantID

	product, err := app.products.Update(id, updateInput)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.httpErrRes.BadRequestResponse(w, r, errors.New("product does not exists"))
		default:
			app.httpErrRes.ServerErrorResponse(w, r, err)
		}
		return
	}

	err = httpx.WriteJSON(w, http.StatusOK, httpx.Envelope{"product": product}, nil)
	if err != nil {
		app.httpErrRes.ServerErrorResponse(w, r, err)
	}
}

func (app *application) listProductsHandler(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	limit := httpx.ReadInt(queryParams, "limit", 10)
	offset := httpx.ReadInt(queryParams, "offset", 0)

	products, err := app.products.GetAll(limit, offset)
	if err != nil {
		app.httpErrRes.ServerErrorResponse(w, r, err)
		return
	}

	err = httpx.WriteJSON(w, http.StatusOK, httpx.Envelope{"products": products}, nil)
	if err != nil {
		app.httpErrRes.ServerErrorResponse(w, r, err)
	}
}

func (app *application) listProductVariantsHandler(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.ReadIDParam(r)
	if err != nil {
		app.httpErrRes.NotFoundResponse(w, r)
		return
	}

	variants, err := app.variants.GetVariantsByProduct(id)
	if err != nil {
		app.httpErrRes.ServerErrorResponse(w, r, err)
		return
	}

	err = httpx.WriteJSON(w, http.StatusOK, httpx.Envelope{"variants": variants}, nil)
	if err != nil {
		app.httpErrRes.ServerErrorResponse(w, r, err)
	}
}
