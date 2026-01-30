package http

import "net/http"

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /v1/users", h.createUserHandler)
	mux.HandleFunc("GET /v1/users/{id}", h.getUserByIDHandler)
	mux.HandleFunc("PATCH /v1/users/{id}", h.updateUserHandler)
	mux.HandleFunc("DELETE /v1/users/{id}", h.deleteUserHandler)

	return mux
}
