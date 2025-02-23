package validator

import (
	"regexp"
)

func Validate(data map[string]interface{}, schema Schema) ValidationResult {
	var validationErrors []ValidationError

	for field, rule := range schema {
		value, exists := data[field]

		if rule.Required && !exists {
			validationErrors = append(validationErrors, ValidationError{Field: field, Message: "Campo requerido"})
			continue
		}

		// Validar tipo de dato
		if exists {
			switch rule.Type {
			case "string":
				str, ok := value.(string)
				if !ok {
					validationErrors = append(validationErrors, ValidationError{Field: field, Message: "Debe ser un string"})
					continue
				}
				if rule.MinLength != nil && len(str) < *rule.MinLength {
					validationErrors = append(validationErrors, ValidationError{Field: field, Message: "String demasiado corto"})
				}
				if rule.MaxLength != nil && len(str) > *rule.MaxLength {
					validationErrors = append(validationErrors, ValidationError{Field: field, Message: "String demasiado largo"})
				}
				if rule.Regex != nil {
					re := regexp.MustCompile(*rule.Regex)
					if !re.MatchString(str) {
						validationErrors = append(validationErrors, ValidationError{Field: field, Message: "Formato inválido"})
					}
				}

			case "int":
				num, ok := value.(int)
				if !ok {
					validationErrors = append(validationErrors, ValidationError{Field: field, Message: "Debe ser un número entero"})
					continue
				}
				if rule.Min != nil && num < *rule.Min {
					validationErrors = append(validationErrors, ValidationError{Field: field, Message: "Número demasiado pequeño"})
				}
				if rule.Max != nil && num > *rule.Max {
					validationErrors = append(validationErrors, ValidationError{Field: field, Message: "Número demasiado grande"})
				}
			}
		}
	}

	return ValidationResult{
		IsValid: len(validationErrors) == 0,
		Errors:  validationErrors,
	}
}
