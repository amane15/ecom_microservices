package http

import (
	"errors"
	"net/http"

	"github.com/amane15/ecom_microservice/pkg/httpx"
	"github.com/amane15/ecom_microservice/services/products/internal/data"
	"github.com/amane15/ecom_microservice/services/products/internal/service"
	"github.com/jackc/pgx/v5"
)

func (h *Handler) createProductHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("create product request received")
	input := &service.CreateProductInput{}

	err := httpx.ReadJSON(w, r, &input)
	if err != nil {
		h.logger.Error("error reading json", "error", err)
		httpx.BadRequestResponse(w, r, err)
		return
	}
	h.logger.Debug("successfully read json")

	product, err := h.productService.CreateProduct(r.Context(), input)
	if err != nil {
		var validationErrors *data.ValidationErrors
		switch {
		case errors.Is(err, data.ErrDuplicateSlug):
			httpx.BadRequestResponse(w, r, errors.New("product already exists"))
		case errors.As(err, &validationErrors):
			httpx.FailedValidationResponse(w, r, validationErrors.Fields)
		default:
			h.logger.Error("error while creating product", "error", err)
			httpx.ServerErrorResponse(w, r, err)
		}
		return
	}

	err = httpx.WriteJSON(w, http.StatusCreated, httpx.Envelope{"product": product}, nil)
	if err != nil {
		h.logger.Error("errror while sending response", "error", err)
		httpx.ServerErrorResponse(w, r, err)
	}
	h.logger.Debug("response sent successfully")
}

func (h *Handler) getProductHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("get product request received")
	id, err := httpx.ReadIDParam(r)
	if err != nil {
		httpx.NotFoundResponse(w, r)
		return
	}

	product, err := h.productService.GetProduct(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			httpx.NotFoundResponse(w, r)
		default:
			h.logger.Error("error while fetching product", "error", err)
			httpx.ServerErrorResponse(w, r, err)
		}
		return
	}

	err = httpx.WriteJSON(w, http.StatusOK, httpx.Envelope{"product": product}, nil)
	if err != nil {
		h.logger.Error("error while sending response", "error", err)
		httpx.ServerErrorResponse(w, r, err)
	}
	h.logger.Debug("successfully sent response")
}

func (h *Handler) updateProductHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("update product request received")
	id, err := httpx.ReadIDParam(r)
	if err != nil {
		httpx.NotFoundResponse(w, r)
		return
	}

	input := &service.UpdateProductInput{}

	err = httpx.ReadJSON(w, r, &input)
	if err != nil {
		h.logger.Error("error while reading json", "error", err)
		httpx.BadRequestResponse(w, r, err)
		return
	}
	h.logger.Debug("read json successfully")

	h.logger.Info("updating product", "id", id, "body", input)
	product, err := h.productService.UpdateProduct(r.Context(), id, input)
	if err != nil {
		var validationErrors *data.ValidationErrors
		switch {
		case errors.Is(err, data.ErrNoFieldsToUpdate):
			httpx.BadRequestResponse(w, r, errors.New("no fields were provided"))
		case errors.As(err, &validationErrors):
			httpx.FailedValidationResponse(w, r, validationErrors.Fields)
		case errors.Is(err, data.ErrRecordNotFound):
			httpx.NotFoundResponse(w, r)
		default:
			h.logger.Error("error while updating product", "error", err)
			httpx.ServerErrorResponse(w, r, err)
		}

		return
	}
	h.logger.Info("updated product", "id", id, "body", input)

	err = httpx.WriteJSON(w, http.StatusOK, httpx.Envelope{"product": product}, nil)
	if err != nil {
		h.logger.Error("error writing response", "error", err)
		httpx.ServerErrorResponse(w, r, err)
	}
	h.logger.Debug("successfully sent response")
}

func (h *Handler) changeProductStatusHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("change product status request received")
	id, err := httpx.ReadIDParam(r)
	if err != nil {
		httpx.NotFoundResponse(w, r)
		return
	}

	var input struct {
		Status *data.ProductStatus `json:"status"`
	}

	err = httpx.ReadJSON(w, r, &input)
	if err != nil {
		h.logger.Error("error while reading body", "error", err)
		httpx.BadRequestResponse(w, r, err)
		return
	}

	h.logger.Info("changing product status", "id", id, "status=", input.Status)
	product, err := h.productService.UpdateProduct(r.Context(),
		id,
		&service.UpdateProductInput{Status: input.Status})
	if err != nil {
		var validationErrors *data.ValidationErrors
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			httpx.NotFoundResponse(w, r)
		case errors.As(err, &validationErrors):
			httpx.FailedValidationResponse(w, r, validationErrors.Fields)
		default:
			h.logger.Info("error while changing product status", "error", err, "id", id)
			httpx.ServerErrorResponse(w, r, err)
		}
		return

	}
	h.logger.Info("changed product status", "id", id, "status=", product.Status)

	err = httpx.WriteJSON(w, http.StatusOK, httpx.Envelope{"product": product}, nil)
	if err != nil {
		h.logger.Error("error sending response", "error", err)
		httpx.ServerErrorResponse(w, r, err)
	}
}

func (h *Handler) setDefaultVariantHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("set default variant request received")
	id, err := httpx.ReadIDParam(r)
	if err != nil {
		httpx.NotFoundResponse(w, r)
		return
	}

	var input struct {
		DefaultVariantID *int64 `json:"default_variant_id"`
	}

	err = httpx.ReadJSON(w, r, &input)
	if err != nil {
		h.logger.Error("error reading body", "error", err)
		httpx.BadRequestResponse(w, r, err)
		return
	}

	product, err := h.productService.UpdateProduct(r.Context(),
		id,
		&service.UpdateProductInput{DefaultVariantID: input.DefaultVariantID})
	if err != nil {
		var validationErrors *data.ValidationErrors
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			httpx.BadRequestResponse(w, r, errors.New("product does not exists"))
		case errors.As(err, &validationErrors):
			httpx.FailedValidationResponse(w, r, validationErrors.Fields)
		default:
			h.logger.Error("error setting default variant", "error", err, "product_id", id, "variant_id", input.DefaultVariantID)
			httpx.ServerErrorResponse(w, r, err)
		}
		return
	}

	err = httpx.WriteJSON(w, http.StatusOK, httpx.Envelope{"product": product}, nil)
	if err != nil {
		h.logger.Error("error sending json response", "error", err)
		httpx.ServerErrorResponse(w, r, err)
	}
}

func (h *Handler) listProductsHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("list products request received")
	queryParams := r.URL.Query()
	limit := httpx.ReadInt(queryParams, "limit", 10)
	offset := httpx.ReadInt(queryParams, "offset", 0)

	products, err := h.productService.ListProducts(r.Context(), int32(limit), int32(offset))
	if err != nil {
		h.logger.Error("error fetching products", "error", err)
		httpx.ServerErrorResponse(w, r, err)
		return
	}

	err = httpx.WriteJSON(w, http.StatusOK, httpx.Envelope{"products": products}, nil)
	if err != nil {
		h.logger.Error("error sending writing json", "error", err)
		httpx.ServerErrorResponse(w, r, err)
	}
}

func (h *Handler) listProductVariantsHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("list product variants request received")
	id, err := httpx.ReadIDParam(r)
	if err != nil {
		httpx.NotFoundResponse(w, r)
		return
	}

	variants, err := h.productService.ListProductVariants(r.Context(), id)
	if err != nil {
		h.logger.Error("error while fetching product variants", "error", err, "id", id)
		httpx.ServerErrorResponse(w, r, err)
		return
	}

	err = httpx.WriteJSON(w, http.StatusOK, httpx.Envelope{"variants": variants}, nil)
	if err != nil {
		h.logger.Error("error while writing json", "error", err)
		httpx.ServerErrorResponse(w, r, err)
	}
}

func (h *Handler) deleteProductHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("delete product request received")
	id, err := httpx.ReadIDParam(r)
	if err != nil {
		httpx.NotFoundResponse(w, r)
		return
	}

	err = h.productService.DeleteProduct(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			httpx.NotFoundResponse(w, r)
		default:
			h.logger.Error("error deleting product", "error", err)
			httpx.ServerErrorResponse(w, r, err)
		}
		return
	}

	err = httpx.WriteJSON(w, http.StatusNoContent, nil, nil)
	if err != nil {
		h.logger.Error("error writing json", "error", err)
		httpx.ServerErrorResponse(w, r, err)
	}
}
