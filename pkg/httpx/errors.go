package httpx

import (
	"log/slog"
	"net/http"
)

type ErrorResponse struct {
	logger *slog.Logger
}

func NewErrorResponseSender(logger *slog.Logger) *ErrorResponse {
	return &ErrorResponse{logger: logger}
}

func (el *ErrorResponse) logError(r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)

	el.logger.Error(err.Error(), "method", method, "uri", uri)
}

func (el *ErrorResponse) sendErrorReponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	env := Envelope{"error": message}

	err := WriteJSON(w, status, env, nil)
	if err != nil {
		el.logError(r, err)
		w.WriteHeader(500)
	}
}

func (el *ErrorResponse) ServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	message := "the server encountered a problem and could not process your request"
	el.sendErrorReponse(w, r, http.StatusInternalServerError, message)
}

func (el *ErrorResponse) BadRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	el.sendErrorReponse(w, r, http.StatusBadRequest, err.Error())
}

func (el *ErrorResponse) FailedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	el.sendErrorReponse(w, r, http.StatusUnprocessableEntity, errors)
}

func (el *ErrorResponse) NotFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	el.sendErrorReponse(w, r, http.StatusNotFound, message)
}
