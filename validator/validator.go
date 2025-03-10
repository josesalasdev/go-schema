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
		// Permitir int, int64 y detectar float64 con valores enteros
		if t.Kind() == reflect.Int || t.Kind() == reflect.Int64 {
			return true
		}
		if t.Kind() == reflect.Float64 {
			return value == float64(int(value.(float64))) // Verifica si es un entero real
		}
		return false
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

	// Validate provided data against schema
	for field, value := range data {
		rule, exists := schema[field]
		if !exists {
			continue // Skip fields not in schema
		}

		// Type validation
		if !matchesType(value, rule.Type) {
			msg := fmt.Sprintf("Invalid type: expected %s, got %T", rule.Type, value)
			if rule.Messages != nil && rule.Messages.TypeMismatch != nil {
				msg = *rule.Messages.TypeMismatch
			}
			validationErrors = append(validationErrors, ValidationError{Field: field, Message: msg})
			continue
		}

		// Type-specific validations
		switch rule.Type {
		case "int", "float":
			if valid, errMsg := validateNumeric(value, rule); !valid {
				if rule.Messages != nil && rule.Messages.Range != nil {
					errMsg = *rule.Messages.Range
				}
				validationErrors = append(validationErrors, ValidationError{Field: field, Message: errMsg})
			}
		case "string":
			if strVal, ok := value.(string); ok {
				if valid, errMsg := validateString(strVal, rule); !valid {
					if rule.Messages != nil && rule.Messages.Length != nil {
						errMsg = *rule.Messages.Length
					}
					validationErrors = append(validationErrors, ValidationError{Field: field, Message: errMsg})
				}
			}
		case "list":
			if listVal, ok := value.([]interface{}); ok && rule.List != nil {
				for i, item := range listVal {
					itemData := map[string]interface{}{"items": item}
					result := Validate(itemData, Schema{"items": *rule.List})
					if !result.IsValid {
						for _, err := range result.Errors {
							validationErrors = append(validationErrors, ValidationError{
								Field:   fmt.Sprintf("%s[%d]", field, i),
								Message: err.Message,
							})
						}
					}
				}
			}
		case "map":
			if mapVal, ok := value.(map[string]interface{}); ok && rule.Schema != nil {
				result := Validate(mapVal, *rule.Schema)
				if !result.IsValid {
					for _, err := range result.Errors {
						validationErrors = append(validationErrors, ValidationError{
							Field:   fmt.Sprintf("%s.%s", field, err.Field),
							Message: err.Message,
						})
					}
				}
			}
		}
	}

	// Check for required fields
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

// Add numeric validation (Min/Max) - call this from Validate
func validateNumeric(value interface{}, rule Rule) (bool, string) {
	if rule.Type == "int" {
		intVal, ok := extractIntValue(value)
		if !ok {
			return false, fmt.Sprintf("Failed to convert %v to integer", value)
		}
		if rule.Min != 0 && intVal < int64(rule.Min) {
			return false, fmt.Sprintf("Value %d is less than minimum %d", intVal, int64(rule.Min))
		}
		if rule.Max != 0 && intVal > int64(rule.Max) {
			return false, fmt.Sprintf("Value %d is greater than maximum %d", intVal, int64(rule.Max))
		}
	} else if rule.Type == "float" {
		floatVal, ok := extractFloatValue(value)
		if !ok {
			return false, fmt.Sprintf("Failed to convert %v to float", value)
		}
		if rule.Min != 0 && floatVal < float64(rule.Min) {
			return false, fmt.Sprintf("Value %f is less than minimum %f", floatVal, rule.Min)
		}
		if rule.Max != 0 && floatVal > float64(rule.Max) {
			return false, fmt.Sprintf("Value %f is greater than maximum %f", floatVal, rule.Max)
		}
	}
	return true, ""
}

// Add string validation (Length/Regex) - call this from Validate
func validateString(value string, rule Rule) (bool, string) {
	if rule.MinLength != 0 && len(value) < rule.MinLength {
		return false, fmt.Sprintf("String length %d is less than minimum %d", len(value), rule.MinLength)
	}
	if rule.MaxLength != 0 && len(value) > rule.MaxLength {
		return false, fmt.Sprintf("String length %d is greater than maximum %d", len(value), rule.MaxLength)
	}
	if rule.Regex != nil && !rule.Regex.MatchString(value) {
		return false, "String does not match pattern"
	}
	return true, ""
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
		// 1. Validar nombre del campo
		if !isValidJSONKey(field) {
			return fmt.Errorf("invalid field name: '%s'", field)
		}

		// 2. Validar tipo
		if _, ok := validTypes[rule.Type]; !ok {
			return fmt.Errorf("invalid type '%s' for field '%s'", rule.Type, field)
		}

		// 3. Validar valores por defecto
		if rule.Default != nil && !matchesType(rule.Default, rule.Type) {
			return fmt.Errorf("default value for '%s' does not match type '%s'", field, rule.Type)
		}

		// 4. Validar Min y Max solo en números
		if (rule.Min != 0 || rule.Max != 0) && rule.Type != "int" && rule.Type != "float" {
			return fmt.Errorf("min/max can only be used for numeric fields, but found in '%s'", field)
		}

		// 5. Validar listas y mapas anidados
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

	return nil
}
