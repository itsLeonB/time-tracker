package apperror

import "net/http"

type AppError struct {
	Err            error
	HttpStatusCode int
	Type           string
	Message        string
}

func (ae *AppError) Error() string {
	return ae.Err.Error()
}

func InternalServerError(err error) *AppError {
	return &AppError{
		Err:            err,
		HttpStatusCode: http.StatusInternalServerError,
		Type:           "InternalServerError",
		Message:        "Unexpected error occurred",
	}
}

func ConflictError(err error) *AppError {
	return &AppError{
		Err:            err,
		HttpStatusCode: http.StatusConflict,
		Type:           "ConflictError",
		Message:        "Resource already exists",
	}
}

func NotFoundError(err error) *AppError {
	return &AppError{
		Err:            err,
		HttpStatusCode: http.StatusNotFound,
		Type:           "NotFoundError",
		Message:        "Resource not found",
	}
}

func BadRequestError(err error) *AppError {
	return &AppError{
		Err:            err,
		HttpStatusCode: http.StatusBadRequest,
		Type:           "BadRequestError",
		Message:        "Bad request",
	}
}
