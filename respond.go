// Package respond provides HTTP JSON response helpers for Go.
//
// Write JSON responses in one line instead of five. Includes structured
// error responses and RFC 9457 Problem Details support.
package respond

import (
	"encoding/json"
	"net/http"
)

// errorResponse is the structure used for error JSON responses.
type errorResponse struct {
	Error errorBody `json:"error"`
}

// errorBody contains the error details within an error response.
type errorBody struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

// JSON writes data as a JSON response with the given HTTP status code.
// It sets the Content-Type header to application/json. If marshalling fails,
// a 500 Internal Server Error is written as plain text.
func JSON(w http.ResponseWriter, status int, data any) {
	body, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(body)
}

// OK writes a 200 OK JSON response with the given data.
func OK(w http.ResponseWriter, data any) {
	JSON(w, http.StatusOK, data)
}

// Created writes a 201 Created JSON response with the given data.
func Created(w http.ResponseWriter, data any) {
	JSON(w, http.StatusCreated, data)
}

// NoContent writes a 204 No Content response with no body.
func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

// Error writes a structured error JSON response with the given status code
// and message. The response body has the form:
//
//	{"error":{"status":N,"message":"..."}}
func Error(w http.ResponseWriter, status int, message string) {
	JSON(w, status, errorResponse{
		Error: errorBody{
			Status:  status,
			Message: message,
		},
	})
}

// ErrorWithDetails writes a structured error JSON response that includes
// an additional details field. The response body has the form:
//
//	{"error":{"status":N,"message":"...","details":...}}
func ErrorWithDetails(w http.ResponseWriter, status int, message string, details any) {
	JSON(w, status, errorResponse{
		Error: errorBody{
			Status:  status,
			Message: message,
			Details: details,
		},
	})
}
