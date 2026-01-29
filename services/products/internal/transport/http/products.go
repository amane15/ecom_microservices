package http

import (
	"errors"
	"net/http"

	"github.com/amane15/ecom_microservice/pkg/httpx"
	"github.com/amane15/ecom_microservice/services/products/internal/data"
	data1 "github.com/amane15/ecom_microservice/services/products/internal/data"
	"github.com/amane15/ecom_microservice/services/products/internal/service"
)

func (h *Handler) createProductHandler(w http.ResponseWriter, r *http.Request) {
	input := &service.CreateProductInput{}

	err := httpx.ReadJSON(w, r, &input)
	if err != nil {
		h.httpErrRes.BadRequestResponse(w, r, err)
		return
	}

	product, err := h.productService.CreateProduct(r.Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateSlug):
			h.httpErrRes.BadRequestResponse(w, r, errors.New("product already exists"))
		default:
			h.httpErrRes.ServerErrorResponse(w, r, err)
		}
		return
	}

	err = httpx.WriteJSON(w, http.StatusCreated, httpx.Envelope{"product": product}, nil)
	if err != nil {
		h.httpErrRes.ServerErrorResponse(w, r, err)
	}
}

func (h *Handler) getProductHandler(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.ReadIDParam(r)
	if err != nil {
		h.httpErrRes.NotFoundResponse(w, r)
		return
	}

	product, err := h.productService.GetProduct(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			h.httpErrRes.NotFoundResponse(w, r)
		default:
			h.httpErrRes.ServerErrorResponse(w, r, err)
		}
		return
	}

	err = httpx.WriteJSON(w, http.StatusOK, httpx.Envelope{"product": product}, nil)
	if err != nil {
		h.httpErrRes.ServerErrorResponse(w, r, err)
	}
}

func (h *Handler) updateProductHandler(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.ReadIDParam(r)
	if err != nil {
		h.httpErrRes.NotFoundResponse(w, r)
		return
	}

	input := &service.UpdateProductInput{}

	err = httpx.ReadJSON(w, r, &input)
	if err != nil {
		h.httpErrRes.BadRequestResponse(w, r, err)
		return
	}

	product, err := h.productService.UpdateProduct(r.Context(), id, input)
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

	err = httpx.WriteJSON(w, http.StatusOK, httpx.Envelope{"product": product}, nil)
	if err != nil {
		h.httpErrRes.ServerErrorResponse(w, r, err)
	}
}

func (h *Handler) changeProductStatusHandler(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.ReadIDParam(r)
	if err != nil {
		h.httpErrRes.NotFoundResponse(w, r)
		return
	}

	var input struct {
		Status *data1.ProductStatus `json:"status"`
	}

	err = httpx.ReadJSON(w, r, &input)
	if err != nil {
		h.httpErrRes.BadRequestResponse(w, r, err)
		return
	}

	product, err := h.productService.UpdateProduct(r.Context(),
		id,
		&service.UpdateProductInput{Status: input.Status})
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			h.httpErrRes.NotFoundResponse(w, r)
		default:
			h.httpErrRes.ServerErrorResponse(w, r, err)
		}
		return

	}

	err = httpx.WriteJSON(w, http.StatusOK, httpx.Envelope{"product": product}, nil)
	if err != nil {
		h.httpErrRes.ServerErrorResponse(w, r, err)
	}
}

func (h *Handler) setDefaultVariantHandler(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.ReadIDParam(r)
	if err != nil {
		h.httpErrRes.NotFoundResponse(w, r)
		return
	}

	var input struct {
		DefaultVariantID *int64 `json:"default_variant_id"`
	}

	err = httpx.ReadJSON(w, r, &input)
	if err != nil {
		h.httpErrRes.BadRequestResponse(w, r, err)
		return
	}

	product, err := h.productService.UpdateProduct(r.Context(),
		id,
		&service.UpdateProductInput{DefaultVariantID: input.DefaultVariantID})
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			h.httpErrRes.BadRequestResponse(w, r, errors.New("product does not exists"))
		default:
			h.httpErrRes.ServerErrorResponse(w, r, err)
		}
		return
	}

	err = httpx.WriteJSON(w, http.StatusOK, httpx.Envelope{"product": product}, nil)
	if err != nil {
		h.httpErrRes.ServerErrorResponse(w, r, err)
	}
}

func (h *Handler) listProductsHandler(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	limit := httpx.ReadInt(queryParams, "limit", 10)
	offset := httpx.ReadInt(queryParams, "offset", 0)

	products, err := h.productService.ListProducts(r.Context(), int32(limit), int32(offset))
	if err != nil {
		h.httpErrRes.ServerErrorResponse(w, r, err)
		return
	}

	err = httpx.WriteJSON(w, http.StatusOK, httpx.Envelope{"products": products}, nil)
	if err != nil {
		h.httpErrRes.ServerErrorResponse(w, r, err)
	}
}

func (h *Handler) listProductVariantsHandler(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.ReadIDParam(r)
	if err != nil {
		h.httpErrRes.NotFoundResponse(w, r)
		return
	}

	variants, err := h.productService.ListProductVariants(r.Context(), id)
	if err != nil {
		h.httpErrRes.ServerErrorResponse(w, r, err)
		return
	}

	err = httpx.WriteJSON(w, http.StatusOK, httpx.Envelope{"variants": variants}, nil)
	if err != nil {
		h.httpErrRes.ServerErrorResponse(w, r, err)
	}
}

func (h *Handler) deleteProductHandler(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.ReadIDParam(r)
	if err != nil {
		h.httpErrRes.NotFoundResponse(w, r)
		return
	}

	err = h.productService.DeleteProduct(r.Context(), id)
	if err != nil {
		h.httpErrRes.ServerErrorResponse(w, r, err)
		return
	}

	err = httpx.WriteJSON(w, http.StatusNoContent, nil, nil)
	if err != nil {
		h.httpErrRes.ServerErrorResponse(w, r, err)
	}
}
