package apierror

import (
	"encoding/json"
	"errors"
	"net/http"
	"todo_api/middleware"
)

type APIError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Cause   error  `json:"-"`
}

func (e *APIError) Error() string {
	if e.Cause != nil {
		return e.Message + ": " + e.Cause.Error()
	}
	return e.Message
}

func New(status int, message string, cause error) *APIError {
	return &APIError{
		Status:  status,
		Message: message,
		Cause:   cause,
	}
}

func BadRequest(message string, cause error) *APIError {
	return New(http.StatusBadRequest, message, cause)
}

func NotFound(message string, cause error) *APIError {
	return New(http.StatusNotFound, message, cause)
}

func Internal(message string, cause error) *APIError {
	return New(http.StatusInternalServerError, message, cause)
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func Write(w http.ResponseWriter, r *http.Request, err error) {
	status := http.StatusInternalServerError
	message := "Internal server error"

	middleware.AddError(r, err)

	var apiErr *APIError
	if errors.As(err, &apiErr) {
		status = apiErr.Status
		message = apiErr.Message
	} else if err != nil {
		message = err.Error()
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}
