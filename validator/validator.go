package validator

import (
	"fmt"
	"reflect"
)

func matchesType(value interface{}, expectedType string) bool {
	t := reflect.TypeOf(value)

	switch expectedType {
	case "string":
		return t.Kind() == reflect.String
	case "int":
		return t.Kind() == reflect.Int || t.Kind() == reflect.Int64
	case "float":
		return t.Kind() == reflect.Float32 || t.Kind() == reflect.Float64
	case "bool":
		return t.Kind() == reflect.Bool
	case "list":
		return t.Kind() == reflect.Slice || t.Kind() == reflect.Array
	case "map":
		return t.Kind() == reflect.Map
	default:
		return false
	}
}

func Validate(data map[string]interface{}, schema Schema) ValidationResult {
	var validationErrors []ValidationError

	// Iteramos sobre los datos, validando solo los campos presentes
	for field, value := range data {
		rule, exists := schema[field]
		if !exists {
			continue // Si el campo no está en el esquema, lo ignoramos
		}

		// Validar tipo
		if !matchesType(value, rule.Type) {
			msg := fmt.Sprintf("Invalid type: expected %s, got %T", rule.Type, value)
			if rule.Messages != nil && rule.Messages.TypeMismatch != nil {
				msg = *rule.Messages.TypeMismatch
			}
			validationErrors = append(validationErrors, ValidationError{Field: field, Message: msg})
			continue
		}
	}

	// Ahora verificamos los campos requeridos en el schema que no están en data
	for field, rule := range schema {
		if rule.Required {
			if _, exists := data[field]; !exists {
				msg := "Field is required"
				if rule.Messages != nil && rule.Messages.Required != nil {
					msg = *rule.Messages.Required
				}
				validationErrors = append(validationErrors, ValidationError{Field: field, Message: msg})
			}
		}
	}

	return ValidationResult{
		IsValid: len(validationErrors) == 0,
		Errors:  validationErrors,
	}
}

func isValidJSONKey(key string) bool {
	// Un JSON key válido no debe contener caracteres de control ni espacios en blanco
	for _, r := range key {
		if r <= 0x1F || r == ' ' { // Caracteres de control y espacio
			return false
		}
	}
	return true
}

func ValidateSchema(schema Schema) error {
	validTypes := map[string]bool{
		"string": true, "int": true, "float": true,
		"bool": true, "list": true, "map": true,
	}

	for field, rule := range schema {
		// 1. Validar que el nombre del campo sea un string válido en JSON
		if !isValidJSONKey(field) {
			return fmt.Errorf("invalid field name: '%s'", field)
		}

		// 2. Validar que el tipo esté en la lista de tipos permitidos
		if _, ok := validTypes[rule.Type]; !ok {
			return fmt.Errorf("invalid type '%s' for field '%s'", rule.Type, field)
		}

		// 3. Validar que los valores por defecto sean del tipo correcto
		if rule.Default != nil && !matchesType(rule.Default, rule.Type) {
			return fmt.Errorf("default value for '%s' does not match type '%s'", field, rule.Type)
		}

		// 4. Validar que Min y Max solo existan en tipos numéricos
		if (rule.Min != 0 || rule.Max != 0) && rule.Type != "int" && rule.Type != "float" {
			return fmt.Errorf("min/max can only be used for numeric fields, but found in '%s'", field)
		}

		// 5. Validar esquemas anidados en listas y mapas
		if rule.Type == "list" && rule.List != nil {
			if err := ValidateSchema(Schema{"items": *rule.List}); err != nil {
				return fmt.Errorf("invalid list schema in '%s': %v", field, err)
			}
		}

		if rule.Type == "map" && rule.Schema != nil {
			if err := ValidateSchema(*rule.Schema); err != nil {
				return fmt.Errorf("invalid map schema in '%s': %v", field, err)
			}
		}
	}

	return nil // Si todo está correcto
}
