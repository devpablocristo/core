package httperr

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/devpablocristo/core/saas/go/domainerr"
)

const (
	CodeUnauthorized = "UNAUTHORIZED"
	CodeNotFound     = "NOT_FOUND"
	CodeValidation   = "VALIDATION_ERROR"
	CodeRateLimited  = "RATE_LIMITED"
	CodeInternal     = "INTERNAL"
)

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	RequestID string   `json:"request_id,omitempty"`
	Error     APIError `json:"error"`
}

type HTTPError struct {
	Status  int
	Code    string
	Message string
}

func (e HTTPError) Error() string { return e.Code + ": " + e.Message }

func New(status int, code, message string) HTTPError {
	if status == 0 {
		status = http.StatusInternalServerError
	}
	return HTTPError{Status: status, Code: code, Message: message}
}

var statusByKind = map[domainerr.Kind]int{
	domainerr.KindUnauthorized: http.StatusUnauthorized,
	domainerr.KindForbidden:    http.StatusForbidden,
	domainerr.KindNotFound:     http.StatusNotFound,
	domainerr.KindValidation:   http.StatusBadRequest,
	domainerr.KindConflict:     http.StatusConflict,
	domainerr.KindInternal:     http.StatusInternalServerError,
}

func Normalize(err error) (int, APIError) {
	var de domainerr.Error
	if errors.As(err, &de) {
		status := statusByKind[de.Kind()]
		if status == 0 {
			status = http.StatusInternalServerError
		}
		return status, APIError{Code: string(de.Kind()), Message: de.Message()}
	}

	var he HTTPError
	if errors.As(err, &he) {
		status := he.Status
		if status == 0 {
			status = http.StatusInternalServerError
		}
		return status, APIError{Code: he.Code, Message: he.Message}
	}

	return http.StatusInternalServerError, APIError{Code: CodeInternal, Message: "internal error"}
}

func WriteJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		slog.Error("httperr: failed to encode json response", "error", err)
	}
}

func Write(w http.ResponseWriter, status int, code, message string) {
	if status == 0 {
		status = http.StatusInternalServerError
	}
	WriteJSON(w, status, ErrorResponse{
		Error: APIError{
			Code:    code,
			Message: message,
		},
	})
}

func WriteFrom(w http.ResponseWriter, err error) {
	status, apiErr := Normalize(err)
	Write(w, status, apiErr.Code, apiErr.Message)
}

func BadRequest(w http.ResponseWriter, message string) {
	Write(w, http.StatusBadRequest, CodeValidation, message)
}

func Unauthorized(w http.ResponseWriter, message string) {
	Write(w, http.StatusUnauthorized, CodeUnauthorized, message)
}
