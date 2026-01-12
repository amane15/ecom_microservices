package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	// mux.HandleFunc("GET /v1/healthcheck", app.healthCheckHandler)
	mux.HandleFunc("POST /v1/users", app.createUserHandler)
	mux.HandleFunc("GET /v1/users/{id}", app.getUserByIDHandler)
	mux.HandleFunc("PATCH /v1/users/{id}", app.updateUserHandler)

	return mux
}
