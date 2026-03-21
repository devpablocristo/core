package domainerr

import "fmt"

// Kind categoriza errores de dominio sin acoplarlos a transporte.
type Kind string

const (
	KindUnauthorized Kind = "UNAUTHORIZED"
	KindForbidden    Kind = "FORBIDDEN"
	KindNotFound     Kind = "NOT_FOUND"
	KindValidation   Kind = "VALIDATION_ERROR"
	KindConflict     Kind = "CONFLICT"
	KindInternal     Kind = "INTERNAL"
)

// Error representa un error de dominio categorizado.
type Error struct {
	kind    Kind
	message string
}

func (e Error) Error() string { return string(e.kind) + ": " + e.message }

func (e Error) Kind() Kind { return e.kind }

func (e Error) Message() string { return e.message }

func New(kind Kind, message string) Error {
	return Error{kind: kind, message: message}
}

func Newf(kind Kind, format string, args ...any) Error {
	return Error{kind: kind, message: fmt.Sprintf(format, args...)}
}

func Unauthorized(message string) Error { return New(KindUnauthorized, message) }
func Forbidden(message string) Error    { return New(KindForbidden, message) }
func NotFound(message string) Error     { return New(KindNotFound, message) }
func Validation(message string) Error   { return New(KindValidation, message) }
func Conflict(message string) Error     { return New(KindConflict, message) }
func Internal(message string) Error     { return New(KindInternal, message) }

func (e Error) Is(target error) bool {
	typed, ok := target.(Error)
	if !ok {
		return false
	}
	return e.kind == typed.kind
}
