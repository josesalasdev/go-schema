package validator

import "fmt"

type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("Error en '%s': %s", e.Field, e.Message)
}

type ValidationResult struct {
	IsValid bool
	Errors  []ValidationError
}
