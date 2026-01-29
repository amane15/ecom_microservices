package http

import (
	"errors"
	"net/http"

	"github.com/amane15/ecom_microservice/pkg/httpx"
	"github.com/amane15/ecom_microservice/services/products/internal/data"
	"github.com/amane15/ecom_microservice/services/products/internal/service"
)

func (h *Handler) getCategoryHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("get category request received")
	id, err := httpx.ReadIDParam(r)
	if err != nil {
		httpx.NotFoundResponse(w, r)
		return
	}

	category, err := h.categoryService.GetCategory(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			httpx.NotFoundResponse(w, r)
		default:
			h.logger.Error("error while fetching category", "error", err, "id", id)
			httpx.ServerErrorResponse(w, r, err)
		}
		return
	}

	err = httpx.WriteJSON(w, http.StatusOK, httpx.Envelope{"category": category}, nil)
	if err != nil {
		h.logger.Error("error writing json", "error", err)
		httpx.ServerErrorResponse(w, r, err)
	}
}

func (h *Handler) createCategoryHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("create category request received")
	input := &service.CreateCategoryInput{}

	err := httpx.ReadJSON(w, r, &input)
	if err != nil {
		h.logger.Error("error while ready body", "error", err)
		httpx.BadRequestResponse(w, r, err)
		return
	}

	category, err := h.categoryService.CreateCategory(r.Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateSlug):
			httpx.BadRequestResponse(w, r, errors.New("category already exists"))
		default:
			h.logger.Error("error while creating a category", "error", err)
			httpx.ServerErrorResponse(w, r, err)
		}

		return
	}

	err = httpx.WriteJSON(w, http.StatusCreated, httpx.Envelope{"category": category}, nil)
	if err != nil {
		h.logger.Error("error while writing a json", "error", err)
		httpx.ServerErrorResponse(w, r, err)
	}
}

func (h *Handler) updateCategoryHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("update category request received")
	id, err := httpx.ReadIDParam(r)
	if err != nil {
		httpx.NotFoundResponse(w, r)
		return
	}

	input := &service.UpdateCategoryInput{}
	err = httpx.ReadJSON(w, r, &input)
	if err != nil {
		h.logger.Error("error reading body", "error", err)
		httpx.BadRequestResponse(w, r, err)
		return
	}

	category, err := h.categoryService.UpdateCategory(r.Context(), id, input)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrNoFieldsToUpdate):
			httpx.BadRequestResponse(w, r, errors.New("no fields were provided"))
		case errors.Is(err, data.ErrRecordNotFound):
			httpx.NotFoundResponse(w, r)
		default:
			h.logger.Error("error while updating category", "error", err, "id", id)
			httpx.ServerErrorResponse(w, r, err)
		}

		return
	}

	err = httpx.WriteJSON(w, http.StatusOK, httpx.Envelope{"category": category}, nil)
	if err != nil {
		h.logger.Error("error while writing json", "error", err)
		httpx.ServerErrorResponse(w, r, err)
	}
}

func (h *Handler) markActiveHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("mark active category request received")
	id, err := httpx.ReadIDParam(r)
	if err != nil {
		httpx.NotFoundResponse(w, r)
		return
	}

	var input struct {
		IsActive *bool `json:"is_active"`
	}

	err = httpx.ReadJSON(w, r, &input)
	if err != nil {
		h.logger.Error("error reading body", "error", err)
		httpx.BadRequestResponse(w, r, err)
		return
	}

	category, err := h.categoryService.UpdateCategory(r.Context(), id, &service.UpdateCategoryInput{IsActive: input.IsActive})
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			httpx.NotFoundResponse(w, r)
		default:
			h.logger.Error("error marking active category", "error", err, "id", id)
			httpx.ServerErrorResponse(w, r, err)
		}
		return
	}

	err = httpx.WriteJSON(w, http.StatusOK, httpx.Envelope{"category": category}, nil)
	if err != nil {
		h.logger.Error("error writing a json", "error", err)
		httpx.ServerErrorResponse(w, r, err)
	}
}

func (h *Handler) listCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("list category request received")
	queryParams := r.URL.Query()
	limit := httpx.ReadInt(queryParams, "limit", 10)
	offset := httpx.ReadInt(queryParams, "offset", 0)

	categories, err := h.categoryService.ListCategories(r.Context(), int32(limit), int32(offset))
	if err != nil {
		h.logger.Error("error fetching categories", "error", err)
		httpx.ServerErrorResponse(w, r, err)
		return
	}

	err = httpx.WriteJSON(w, http.StatusOK, httpx.Envelope{"categories": categories}, nil)
	if err != nil {
		h.logger.Error("error writing a json", "error", err)
		httpx.ServerErrorResponse(w, r, err)
	}
}
