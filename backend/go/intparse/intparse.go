// Package intparse ofrece parseo mínimo de enteros para query params y similares,
// sin acoplar a HTTP (pagination y httpjson reutilizan esto).
package intparse

import (
	"errors"
	"strconv"
)

// ParsePositiveInt parsea un entero positivo o devuelve defaultValue si raw está vacío.
func ParsePositiveInt(raw string, defaultValue int) (int, error) {
	if raw == "" {
		return defaultValue, nil
	}
	value, err := strconv.Atoi(raw)
	if err != nil || value <= 0 {
		return 0, errors.New("invalid positive integer")
	}
	return value, nil
}
