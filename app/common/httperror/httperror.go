package httperror

import "net/http"

type HttpError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func NewBadRequestError(message string) HttpError {
	return HttpError{
		Status:  http.StatusBadRequest,
		Message: message,
	}
}
func NewInternalServerError(message string) HttpError {
	return HttpError{
		Status:  http.StatusInternalServerError,
		Message: message,
	}
}
