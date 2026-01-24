package http

import "net/http"

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /v1/carts", h.createCartHandler)
	mux.HandleFunc("GET /v1/carts/{id}", h.getCartHandler)
	mux.HandleFunc("DELETE /v1/carts/{id}", h.deleteCartHandler)

	mux.HandleFunc("POST /v1/carts/{id}/items", h.createCartItemHandler)
	mux.HandleFunc("DELETE /v1/items/{id}", h.deleteCartItemHandler)

	return mux
}
