package router

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

// HTTPError represents an HTTP error with a status code and message
type HTTPError struct {
	Code    int    `json:"-"`
	Message string `json:"error"`
	Details any    `json:"details,omitempty"`
}

func (e *HTTPError) Error() string {
	return e.Message
}

// NewHTTPError creates a new HTTPError
func NewHTTPError(code int, message string) *HTTPError {
	return &HTTPError{
		Code:    code,
		Message: message,
	}
}

// WithDetails adds details to an HTTPError
func (e *HTTPError) WithDetails(details any) *HTTPError {
	e.Details = details
	return e
}

// Common HTTP errors
var (
	ErrBadRequest          = NewHTTPError(http.StatusBadRequest, "bad request")
	ErrUnauthorized        = NewHTTPError(http.StatusUnauthorized, "unauthorized")
	ErrForbidden           = NewHTTPError(http.StatusForbidden, "forbidden")
	ErrNotFound            = NewHTTPError(http.StatusNotFound, "not found")
	ErrMethodNotAllowed    = NewHTTPError(http.StatusMethodNotAllowed, "method not allowed")
	ErrConflict            = NewHTTPError(http.StatusConflict, "conflict")
	ErrUnprocessableEntity = NewHTTPError(http.StatusUnprocessableEntity, "unprocessable entity")
	ErrTooManyRequests     = NewHTTPError(http.StatusTooManyRequests, "too many requests")
	ErrInternalServer      = NewHTTPError(http.StatusInternalServerError, "internal server error")
	ErrServiceUnavailable  = NewHTTPError(http.StatusServiceUnavailable, "service unavailable")
)

// HandlerFunc is a custom handler type that can return errors
type HandlerFunc func(w http.ResponseWriter, r *http.Request) error

// Handle wraps a HandlerFunc to handle errors automatically
func Handle(h HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			handleError(w, r, err)
		}
	}
}

// handleError handles errors returned by handlers
func handleError(w http.ResponseWriter, r *http.Request, err error) {
	requestID := r.Header.Get("X-Request-ID")

	// Check if it's an HTTPError
	if httpErr, ok := err.(*HTTPError); ok {
		// Log the error with context
		if httpErr.Code >= 500 {
			slog.Error("internal server error",
				"error", httpErr.Message,
				"details", httpErr.Details,
				"path", r.URL.Path,
				"method", r.Method,
				"request_id", requestID,
			)
		} else {
			slog.Warn("client error",
				"error", httpErr.Message,
				"code", httpErr.Code,
				"details", httpErr.Details,
				"path", r.URL.Path,
				"method", r.Method,
				"request_id", requestID,
			)
		}

		// Write error response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(httpErr.Code)
		json.NewEncoder(w).Encode(httpErr)
		return
	}

	// Unknown error - log and return 500
	slog.Error("unexpected error",
		"error", err.Error(),
		"path", r.URL.Path,
		"method", r.Method,
		"request_id", requestID,
	)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(map[string]string{
		"error": "internal server error",
	})
}
