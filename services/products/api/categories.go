package main

import (
	"errors"
	"net/http"

	"github.com/amane15/ecom_microservice/pkg/httpx"
	"github.com/amane15/ecom_microservice/pkg/validator"
	"github.com/amane15/ecom_microservice/services/products/internal/data"
)

func (app *application) getCategoryHandler(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.ReadIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	category, err := app.categories.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = httpx.WriteJSON(w, http.StatusOK, httpx.Envelope{"category": category}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) createCategoryHandler(w http.ResponseWriter, r *http.Request) {
	input := data.CreateCategoryInput{}

	err := httpx.ReadJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	category := &data.Category{
		Name: input.Name,
		Slug: input.Slug,
	}

	if input.IsActive == nil {
		category.IsActive = false
	} else {
		category.IsActive = *input.IsActive
	}

	if input.Description != nil {
		category.Description = input.Description
	}

	v := validator.New()

	if data.ValidateCategory(v, category); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.categories.Insert(category)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateSlug):
			app.badRequestResponse(w, r, errors.New("category already exists"))
		default:
			app.serverErrorResponse(w, r, err)
		}

		return
	}

	err = httpx.WriteJSON(w, http.StatusCreated, httpx.Envelope{"category": category}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateCategoryHandler(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.ReadIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	input := &data.UpdateCategoryInput{}
	err = httpx.ReadJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Name == nil && input.IsActive == nil && input.Description == nil {
		app.badRequestResponse(w, r, errors.New("at least one field must be provided"))
		return
	}

	v := validator.New()

	if data.ValidateCategoryUpdate(v, input); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	category, err := app.categories.Update(id, input)
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

	err = httpx.WriteJSON(w, http.StatusOK, httpx.Envelope{"category": category}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) markActiveHandler(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.ReadIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	input := &data.UpdateCategoryInput{}

	err = httpx.ReadJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	if input.IsActive == nil {
		v.AddError("is_active", "must be provided")
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	category, err := app.categories.Update(id, input)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = httpx.WriteJSON(w, http.StatusOK, httpx.Envelope{"category": category}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	limit := httpx.ReadInt(queryParams, "limit", 10)
	offset := httpx.ReadInt(queryParams, "offset", 0)

	categories, err := app.categories.GetAll(limit, offset)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = httpx.WriteJSON(w, http.StatusOK, httpx.Envelope{"categories": categories}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
