package httpx

import (
	"net/http"
)

func ErrorReponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	env := Envelope{"error": message}

	err := WriteJSON(w, status, env, nil)
	if err != nil {
		w.WriteHeader(500)
	}
}

func ServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	message := "the server encountered a problem and could not process your request"
	ErrorReponse(w, r, http.StatusInternalServerError, message)
}

func BadRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	ErrorReponse(w, r, http.StatusBadRequest, err.Error())
}

func FailedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	ErrorReponse(w, r, http.StatusUnprocessableEntity, errors)
}

func NotFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	ErrorReponse(w, r, http.StatusNotFound, message)
}
