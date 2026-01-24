package http

import (
	"errors"
	"net/http"
	"unicode/utf8"

	"github.com/amane15/ecom_microservice/pkg/httpx"
	"github.com/amane15/ecom_microservice/pkg/validator"
	"github.com/amane15/ecom_microservice/services/products/internal/data"
)

func (h *Handler) createProductHandler(w http.ResponseWriter, r *http.Request) {
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
		h.httpErrRes.BadRequestResponse(w, r, err)
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
		h.httpErrRes.FailedValidationResponse(w, r, v.Errors)
		return
	}

	err = h.app.Models.Products.Create(product)
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

	product, err := h.app.Models.Products.Get(id)
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
	productInput := data.UpdateProductInput{}

	err = httpx.ReadJSON(w, r, &productInput)
	if err != nil {
		h.httpErrRes.BadRequestResponse(w, r, err)
		return
	}

	if productInput.Name == nil && productInput.Description == nil &&
		productInput.ShortDescription == nil && productInput.MetaTitle == nil &&
		productInput.MetaDescription == nil {
		h.httpErrRes.BadRequestResponse(w, r, errors.New("at least one field must be provided"))
		return
	}

	v := validator.New()
	if data.ValidateUpdateProduct(v, &productInput); !v.Valid() {
		h.httpErrRes.FailedValidationResponse(w, r, v.Errors)
		return
	}

	product, err := h.app.Models.Products.Update(id, &productInput)
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
		Status *data.ProductStatus `json:"status"`
	}

	err = httpx.ReadJSON(w, r, &input)
	if err != nil {
		h.httpErrRes.BadRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	if input.Status == nil {
		v.AddError("status", "status must be provided")
		h.httpErrRes.FailedValidationResponse(w, r, v.Errors)
		return
	}
	updateInput := &data.UpdateProductInput{}
	status := *input.Status

	switch status {
	case data.ProductStatusDraft, data.ProductStatusArchived, data.ProductStatusActive:
		if status == data.ProductStatusActive {
			count, err := h.app.Models.Variants.GetVariantCountForProduct(id)
			if err != nil {
				switch {
				case errors.Is(err, data.ErrRecordNotFound):
					count = 0
				default:
					h.httpErrRes.ServerErrorResponse(w, r, err)
					return
				}
			}

			if status == data.ProductStatusActive && count == 0 {
				h.httpErrRes.BadRequestResponse(w, r, errors.New("to make a product active you need have at least 1 variant"))
				return
			}
		}
		updateInput.Status = &status
	default:
		v.AddError("status", "must be a valid status")
		h.httpErrRes.FailedValidationResponse(w, r, v.Errors)
		return

	}

	product, err := h.app.Models.Products.Update(id, updateInput)
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

	v := validator.New()

	if input.DefaultVariantID == nil {
		v.AddError("default_variant_id", "must be provided")
		h.httpErrRes.FailedValidationResponse(w, r, v.Errors)
		return
	}

	variantID := *input.DefaultVariantID
	exists, err := h.app.Models.Variants.IsVariantExists(variantID)
	if err != nil {
		h.httpErrRes.ServerErrorResponse(w, r, err)
		return
	}
	if !exists {
		h.httpErrRes.BadRequestResponse(w, r, errors.New("variant does not exists"))
		return
	}

	updateInput := &data.UpdateProductInput{}
	updateInput.DefaultVariantID = input.DefaultVariantID

	product, err := h.app.Models.Products.Update(id, updateInput)
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

	products, err := h.app.Models.Products.GetAll(limit, offset)
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

	variants, err := h.app.Models.Variants.GetVariantsByProduct(id)
	if err != nil {
		h.httpErrRes.ServerErrorResponse(w, r, err)
		return
	}

	err = httpx.WriteJSON(w, http.StatusOK, httpx.Envelope{"variants": variants}, nil)
	if err != nil {
		h.httpErrRes.ServerErrorResponse(w, r, err)
	}
}
