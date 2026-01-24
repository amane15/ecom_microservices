package http

import (
	"errors"
	"net/http"

	"github.com/amane15/ecom_microservice/pkg/httpx"
	"github.com/amane15/ecom_microservice/pkg/validator"
	"github.com/amane15/ecom_microservice/services/products/internal/data"
)

func (h *Handler) getCategoryHandler(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.ReadIDParam(r)
	if err != nil {
		h.httpErrRes.NotFoundResponse(w, r)
		return
	}

	category, err := h.app.Models.Categories.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			h.httpErrRes.NotFoundResponse(w, r)
		default:
			h.httpErrRes.ServerErrorResponse(w, r, err)
		}
		return
	}

	err = httpx.WriteJSON(w, http.StatusOK, httpx.Envelope{"category": category}, nil)
	if err != nil {
		h.httpErrRes.ServerErrorResponse(w, r, err)
	}
}

func (h *Handler) createCategoryHandler(w http.ResponseWriter, r *http.Request) {
	input := data.CreateCategoryInput{}

	err := httpx.ReadJSON(w, r, &input)
	if err != nil {
		h.httpErrRes.BadRequestResponse(w, r, err)
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
		h.httpErrRes.FailedValidationResponse(w, r, v.Errors)
		return
	}

	err = h.app.Models.Categories.Insert(category)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateSlug):
			h.httpErrRes.BadRequestResponse(w, r, errors.New("category already exists"))
		default:
			h.httpErrRes.ServerErrorResponse(w, r, err)
		}

		return
	}

	err = httpx.WriteJSON(w, http.StatusCreated, httpx.Envelope{"category": category}, nil)
	if err != nil {
		h.httpErrRes.ServerErrorResponse(w, r, err)
	}
}

func (h *Handler) updateCategoryHandler(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.ReadIDParam(r)
	if err != nil {
		h.httpErrRes.NotFoundResponse(w, r)
		return
	}

	input := &data.UpdateCategoryInput{}
	err = httpx.ReadJSON(w, r, &input)
	if err != nil {
		h.httpErrRes.BadRequestResponse(w, r, err)
		return
	}

	if input.Name == nil && input.IsActive == nil && input.Description == nil {
		h.httpErrRes.BadRequestResponse(w, r, errors.New("at least one field must be provided"))
		return
	}

	v := validator.New()

	if data.ValidateCategoryUpdate(v, input); !v.Valid() {
		h.httpErrRes.FailedValidationResponse(w, r, v.Errors)
		return
	}

	category, err := h.app.Models.Categories.Update(id, input)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrNoFieldsToUpdate):
			h.httpErrRes.BadRequestResponse(w, r, errors.New("no fields were provided"))
		case errors.Is(err, data.ErrRecordNotFound):
			h.httpErrRes.NotFoundResponse(w, r)
		default:
			h.httpErrRes.ServerErrorResponse(w, r, err)
		}

		return
	}

	err = httpx.WriteJSON(w, http.StatusOK, httpx.Envelope{"category": category}, nil)
	if err != nil {
		h.httpErrRes.ServerErrorResponse(w, r, err)
	}
}

func (h *Handler) markActiveHandler(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.ReadIDParam(r)
	if err != nil {
		h.httpErrRes.NotFoundResponse(w, r)
		return
	}

	input := &data.UpdateCategoryInput{}

	err = httpx.ReadJSON(w, r, &input)
	if err != nil {
		h.httpErrRes.BadRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	if input.IsActive == nil {
		v.AddError("is_active", "must be provided")
		h.httpErrRes.FailedValidationResponse(w, r, v.Errors)
		return
	}

	category, err := h.app.Models.Categories.Update(id, input)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			h.httpErrRes.NotFoundResponse(w, r)
		default:
			h.httpErrRes.ServerErrorResponse(w, r, err)
		}
		return
	}

	err = httpx.WriteJSON(w, http.StatusOK, httpx.Envelope{"category": category}, nil)
	if err != nil {
		h.httpErrRes.ServerErrorResponse(w, r, err)
	}
}

func (h *Handler) listCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	limit := httpx.ReadInt(queryParams, "limit", 10)
	offset := httpx.ReadInt(queryParams, "offset", 0)

	categories, err := h.app.Models.Categories.GetAll(limit, offset)
	if err != nil {
		h.httpErrRes.ServerErrorResponse(w, r, err)
		return
	}

	err = httpx.WriteJSON(w, http.StatusOK, httpx.Envelope{"categories": categories}, nil)
	if err != nil {
		h.httpErrRes.ServerErrorResponse(w, r, err)
	}
}
