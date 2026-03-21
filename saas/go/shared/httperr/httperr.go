package httperr

import (
	"net/http"

	base "github.com/devpablocristo/core/saas/go/httperr"
)

type APIError = base.APIError
type ErrorResponse = base.ErrorResponse
type HTTPError = base.HTTPError

const (
	CodeUnauthorized = base.CodeUnauthorized
	CodeNotFound     = base.CodeNotFound
	CodeValidation   = base.CodeValidation
	CodeRateLimited  = base.CodeRateLimited
	CodeInternal     = base.CodeInternal
)

var (
	New          = base.New
	Normalize    = base.Normalize
	WriteJSON    = base.WriteJSON
	Write        = base.Write
	WriteFrom    = base.WriteFrom
	BadRequest   = base.BadRequest
	Unauthorized = base.Unauthorized
)

func Forbidden(w http.ResponseWriter, message string) {
	base.Write(w, http.StatusForbidden, CodeUnauthorized, message)
}

func NotFound(w http.ResponseWriter, message string) {
	base.Write(w, http.StatusNotFound, CodeNotFound, message)
}
