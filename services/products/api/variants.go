package main

import (
	"errors"
	"net/http"

	"github.com/amane15/ecom_microservice/pkg/validator"
	"github.com/amane15/ecom_microservice/services/products/internal/data"
)

func (app *application) getVariantHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	variant, err := app.variants.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}

		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"variant": variant}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) createVariantHandler(w http.ResponseWriter, r *http.Request) {
	input := &data.CreateVariantInput{}

	err := app.readJSON(w, r, input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	app.logger.Info("variant", "price", input.Price)

	variant := &data.ProductVariant{
		ProductID: input.ProductID,
		Slug:      input.Slug,
		Name:      input.Name,
		Price:     input.Price,
	}
	app.logger.Info("variant og struct", "price", input.Price)

	if input.IsActive != nil {
		variant.IsActive = *input.IsActive
	}

	v := validator.New()

	if data.ValidateVariant(v, variant); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.variants.Insert(variant)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateSlug):
			app.badRequestResponse(w, r, errors.New("variant with this slug already exists"))
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"variant": variant}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateVariantHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	input := &data.UpdateVariantInput{}

	err = app.readJSON(w, r, input)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if input.Name == nil && input.IsActive == nil && input.Price == nil {
		app.badRequestResponse(w, r, errors.New("at least one field must be provided"))
		return
	}

	v := validator.New()

	if data.ValidateUpdateVariantInput(v, input); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	variant, err := app.variants.Update(id, input)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"variant": variant}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteVariantHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.variants.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusNoContent, nil, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
