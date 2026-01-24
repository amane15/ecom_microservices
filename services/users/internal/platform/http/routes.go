package http

import "net/http"

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()

	// mux.HandleFunc("GET /v1/healthcheck", app.healthCheckHandler)
	mux.HandleFunc("POST /v1/users", h.createUserHandler)
	mux.HandleFunc("GET /v1/users/{id}", h.getUserByIDHandler)
	mux.HandleFunc("PATCH /v1/users/{id}", h.updateUserHandler)

	return mux
}
