package validation

import (
	"net/http"
	"strings"

	"github.com/devpablocristo/core/backend/httpjson"
)

// FieldError representa una violación de validación tipada.
type FieldError struct {
	Field   string `json:"field"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Errors agrupa errores de validación y satisface `error`.
type Errors []FieldError

func (e Errors) Error() string {
	if len(e) == 0 {
		return "validation failed"
	}
	return e[0].Field + ": " + e[0].Message
}

func (e Errors) HasAny() bool {
	return len(e) > 0
}

// Append agrega un error si no es nil.
func (e Errors) Append(item *FieldError) Errors {
	if item == nil {
		return e
	}
	return append(e, *item)
}

// WriteHTTP serializa el envelope de validación estándar.
func (e Errors) WriteHTTP(w http.ResponseWriter) {
	httpjson.WriteJSON(w, http.StatusBadRequest, map[string]any{
		"error": map[string]any{
			"code":    "VALIDATION",
			"message": "validation failed",
			"fields":  []FieldError(e),
		},
	})
}

func RequiredString(field, value string) *FieldError {
	if strings.TrimSpace(value) != "" {
		return nil
	}
	return &FieldError{
		Field:   field,
		Code:    "required",
		Message: "field is required",
	}
}

func MaxLen(field, value string, max int) *FieldError {
	if max <= 0 || len(value) <= max {
		return nil
	}
	return &FieldError{
		Field:   field,
		Code:    "max_length",
		Message: "field exceeds maximum length",
	}
}

func PositiveInt(field string, value int) *FieldError {
	if value > 0 {
		return nil
	}
	return &FieldError{
		Field:   field,
		Code:    "positive",
		Message: "field must be positive",
	}
}

func OneOf(field, value string, accepted ...string) *FieldError {
	current := strings.TrimSpace(strings.ToLower(value))
	for _, item := range accepted {
		if current == strings.TrimSpace(strings.ToLower(item)) {
			return nil
		}
	}
	return &FieldError{
		Field:   field,
		Code:    "one_of",
		Message: "field has invalid value",
	}
}
