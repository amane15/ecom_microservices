package http

import (
	"errors"
	"net/http"

	"github.com/amane15/ecom_microservice/pkg/httpx"
	"github.com/amane15/ecom_microservice/services/products/internal/data"
	"github.com/amane15/ecom_microservice/services/products/internal/service"
)

func (h *Handler) getVariantHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("get variant request received")
	id, err := httpx.ReadIDParam(r)
	if err != nil {
		httpx.NotFoundResponse(w, r)
		return
	}

	variant, err := h.productService.GetVariant(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			httpx.NotFoundResponse(w, r)
		default:
			h.logger.Error("error while fetching variant", "error", err, "id", id)
			httpx.ServerErrorResponse(w, r, err)
		}

		return
	}

	err = httpx.WriteJSON(w, http.StatusOK, httpx.Envelope{"variant": variant}, nil)
	if err != nil {
		h.logger.Error("error writing a json", "error", err)
		httpx.ServerErrorResponse(w, r, err)
	}
}

func (h *Handler) createVariantHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("craete variant request received")
	id, err := httpx.ReadIDParam(r)
	if err != nil {
		httpx.NotFoundResponse(w, r)
		return
	}

	input := &service.CreateVariantInput{}

	err = httpx.ReadJSON(w, r, input)
	if err != nil {
		httpx.BadRequestResponse(w, r, err)
		return
	}

	variant, err := h.productService.CreateVariant(r.Context(), id, input)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateSlug):
			httpx.BadRequestResponse(w, r, errors.New("variant with this slug already exists"))
		default:
			h.logger.Error("error while creating variant", "error", err)
			httpx.ServerErrorResponse(w, r, err)
		}
		return
	}

	err = httpx.WriteJSON(w, http.StatusCreated, httpx.Envelope{"variant": variant}, nil)
	if err != nil {
		h.logger.Error("error writing a json", "error", err)
		httpx.ServerErrorResponse(w, r, err)
	}
}

func (h *Handler) updateVariantHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("update variant request received")
	id, err := httpx.ReadIDParam(r)
	if err != nil {
		httpx.NotFoundResponse(w, r)
		return
	}

	input := &service.UpdateVariantInput{}

	err = httpx.ReadJSON(w, r, input)
	if err != nil {
		httpx.ServerErrorResponse(w, r, err)
		return
	}

	variant, err := h.productService.UpdateVariant(r.Context(), id, input)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			httpx.NotFoundResponse(w, r)
		default:
			h.logger.Error("error while updating variant", "error", err)
			httpx.ServerErrorResponse(w, r, err)
		}
		return
	}

	err = httpx.WriteJSON(w, http.StatusOK, httpx.Envelope{"variant": variant}, nil)
	if err != nil {
		h.logger.Error("error writing a json", "error", err)
		httpx.ServerErrorResponse(w, r, err)
	}
}

func (h *Handler) deleteVariantHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("delete variant request received")
	id, err := httpx.ReadIDParam(r)
	if err != nil {
		httpx.NotFoundResponse(w, r)
		return
	}

	err = h.productService.DeleteVariant(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			httpx.NotFoundResponse(w, r)
		default:
			h.logger.Error("error while deleting a variant", "error", err)
			httpx.ServerErrorResponse(w, r, err)
		}
		return
	}

	err = httpx.WriteJSON(w, http.StatusNoContent, nil, nil)
	if err != nil {
		h.logger.Error("error writing a json", "error", err)
		httpx.ServerErrorResponse(w, r, err)
	}
}
