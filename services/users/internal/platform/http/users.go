package http

import (
	"errors"
	"net/http"

	"github.com/amane15/ecom_microservice/pkg/httpx"
	"github.com/amane15/ecom_microservice/services/users/internal/data"
	"github.com/amane15/ecom_microservice/services/users/internal/service"
)

func (h *Handler) createUserHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("create user request received")

	input := &service.CreateUserInput{}
	err := httpx.ReadJSON(w, r, &input)
	if err != nil {
		h.logger.Error("error while reading body", "error", err)
		httpx.BadRequestResponse(w, r, err)
		return
	}

	h.logger.Info("creating a user with", "email", input.Email)
	user, err := h.userService.CreateUser(r.Context(), input)
	if err != nil {
		var validationErrors *data.ValidationError
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			httpx.BadRequestResponse(w, r, errors.New("user already exist"))
		case errors.As(err, &validationErrors):
			httpx.FailedValidationResponse(w, r, validationErrors.Fields)
		default:
			httpx.ServerErrorResponse(w, r, err)
		}
		return
	}
	h.logger.Info("user created", "id", user.ID)

	err = httpx.WriteJSON(w, http.StatusCreated, httpx.Envelope{"user": user}, nil)
	if err != nil {
		h.logger.Error("error while writing a json", "error", err)
		httpx.ServerErrorResponse(w, r, err)
	}
}

func (h *Handler) getUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("get user request received")
	id, err := httpx.ReadIDParam(r)
	if err != nil {
		httpx.NotFoundResponse(w, r)
		return
	}

	h.logger.Info("fetching user", "id", id)
	user, err := h.userService.GetUser(r.Context(), id)
	if err != nil {
		var validationErrors *data.ValidationError
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			httpx.NotFoundResponse(w, r)
		case errors.As(err, &validationErrors):
			httpx.FailedValidationResponse(w, r, validationErrors.Fields)
		default:
			h.logger.Error("error while fetching user", "error", err, "id", id)
			httpx.ServerErrorResponse(w, r, err)
		}
		return
	}
	h.logger.Info("fetched user", "id", id)

	err = httpx.WriteJSON(w, http.StatusOK, httpx.Envelope{"user": user}, nil)
	if err != nil {
		h.logger.Error("error while writing a json", "error", err, "id", id)
		httpx.ServerErrorResponse(w, r, err)
	}
}

func (h *Handler) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("update user request received")
	id, err := httpx.ReadIDParam(r)
	if err != nil {
		httpx.NotFoundResponse(w, r)
		return
	}

	input := &service.UpdateUserInput{}

	err = httpx.ReadJSON(w, r, input)
	if err != nil {
		h.logger.Error("error while reading body", "error", err)
		httpx.BadRequestResponse(w, r, err)
		return
	}

	h.logger.Info("updating user", "id", id)
	user, err := h.userService.UpdateUser(r.Context(), id, input)
	if err != nil {
		var validationErrors *data.ValidationError
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			httpx.NotFoundResponse(w, r)
		case errors.As(err, &validationErrors):
			httpx.FailedValidationResponse(w, r, validationErrors.Fields)
		}
		httpx.ServerErrorResponse(w, r, err)
		return
	}
	h.logger.Info("updated user", "id", id)

	err = httpx.WriteJSON(w, http.StatusOK, httpx.Envelope{"user": user}, nil)
	if err != nil {
		h.logger.Error("error while writing a json", "error", err)
		httpx.ServerErrorResponse(w, r, err)
	}
}

func (h *Handler) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("delete user request received")
	id, err := httpx.ReadIDParam(r)
	if err != nil {
		httpx.NotFoundResponse(w, r)
		return
	}

	h.logger.Info("deleting user", "id", id)
	err = h.userService.DeleteUser(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			httpx.NotFoundResponse(w, r)
		default:
			h.logger.Error("error while deleting useer", "error", err, "id", id)
			httpx.ServerErrorResponse(w, r, err)
		}
	}

	err = httpx.WriteJSON(w, http.StatusNoContent, httpx.Envelope{"message": "user has been deleted"}, nil)
	if err != nil {
		h.logger.Error("error while writing a json", "error", err)
		httpx.ServerErrorResponse(w, r, err)
	}
}
