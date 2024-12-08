package model

import (
	"fmt"
)

type JSONResponse struct {
	Success bool  `json:"success"`
	Data    any   `json:"data,omitempty"`
	Error   error `json:"error,omitempty"`
}

func NewSuccessJSON(data any) *JSONResponse {
	return &JSONResponse{
		Success: true,
		Data:    data,
	}
}

func NewErrorJSON(err error) *JSONResponse {
	return &JSONResponse{
		Success: false,
		Error:   err,
	}
}

type ErrorResponse struct {
	Code    int    `json:"-"`
	Type    string `json:"type"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

func (er *ErrorResponse) Error() string {
	return fmt.Sprintf("%s: %s", er.Type, er.Details)
}
