package http

import (
	"errors"
	"net/http"

	"github.com/amane15/ecom_microservice/pkg/httpx"
	"github.com/amane15/ecom_microservice/services/products/internal/data"
	"github.com/amane15/ecom_microservice/services/products/internal/service"
)

func (h *Handler) getVariantHandler(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.ReadIDParam(r)
	if err != nil {
		h.httpErrRes.NotFoundResponse(w, r)
		return
	}

	variant, err := h.productService.GetVariant(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			h.httpErrRes.NotFoundResponse(w, r)
		default:
			h.httpErrRes.ServerErrorResponse(w, r, err)
		}

		return
	}

	err = httpx.WriteJSON(w, http.StatusOK, httpx.Envelope{"variant": variant}, nil)
	if err != nil {
		h.httpErrRes.ServerErrorResponse(w, r, err)
	}
}

func (h *Handler) createVariantHandler(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.ReadIDParam(r)
	if err != nil {
		h.httpErrRes.NotFoundResponse(w, r)
		return
	}

	input := &service.CreateVariantInput{}

	err = httpx.ReadJSON(w, r, input)
	if err != nil {
		h.httpErrRes.BadRequestResponse(w, r, err)
		return
	}

	variant, err := h.productService.CreateVariant(r.Context(), id, input)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateSlug):
			h.httpErrRes.BadRequestResponse(w, r, errors.New("variant with this slug already exists"))
		default:
			h.httpErrRes.ServerErrorResponse(w, r, err)
		}
		return
	}

	err = httpx.WriteJSON(w, http.StatusCreated, httpx.Envelope{"variant": variant}, nil)
	if err != nil {
		h.httpErrRes.ServerErrorResponse(w, r, err)
	}
}

func (h *Handler) updateVariantHandler(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.ReadIDParam(r)
	if err != nil {
		h.httpErrRes.NotFoundResponse(w, r)
		return
	}

	input := &service.UpdateVariantInput{}

	err = httpx.ReadJSON(w, r, input)
	if err != nil {
		h.httpErrRes.ServerErrorResponse(w, r, err)
		return
	}

	variant, err := h.productService.UpdateVariant(r.Context(), id, input)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			h.httpErrRes.NotFoundResponse(w, r)
		default:
			h.httpErrRes.ServerErrorResponse(w, r, err)
		}
		return
	}

	err = httpx.WriteJSON(w, http.StatusOK, httpx.Envelope{"variant": variant}, nil)
	if err != nil {
		h.httpErrRes.ServerErrorResponse(w, r, err)
	}
}

func (h *Handler) deleteVariantHandler(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.ReadIDParam(r)
	if err != nil {
		h.httpErrRes.NotFoundResponse(w, r)
		return
	}

	err = h.productService.DeleteVariant(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			h.httpErrRes.NotFoundResponse(w, r)
		default:
			h.httpErrRes.ServerErrorResponse(w, r, err)
		}
		return
	}

	err = httpx.WriteJSON(w, http.StatusNoContent, nil, nil)
	if err != nil {
		h.httpErrRes.ServerErrorResponse(w, r, err)
	}
}
