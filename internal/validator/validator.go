package validator

import (
	"slices"
	"strings"
	"unicode/utf8"
)

// Validator struct to hold a map of validation errors
type Validator struct {
	FieldErrors map[string]string
}

// Valid returns true if the map is empty
func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0
}

func (v *Validator) AddFieldError(key, message string) {
	// init map if not already initialized
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}
	if _, exist := v.FieldErrors[key]; !exist {
		v.FieldErrors[key] = message
	}
}

func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

func PermittedValue[T comparable](value T, permittedValues ...T) bool {
	return slices.Contains(permittedValues, value)
}
