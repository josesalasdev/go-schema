package validator

import (
	"regexp"
)

func Validate(data map[string]interface{}, schema Schema) ValidationResult {
	var validationErrors []ValidationError

	for field, rule := range schema {
		value, exists := data[field]

		if rule.Required && !exists {
			msg := "Field is required"
			if rule.Messages != nil && rule.Messages.Required != nil {
				msg = *rule.Messages.Required
			}
			validationErrors = append(validationErrors, ValidationError{Field: field, Message: msg})
			continue
		}

		if exists {
			switch rule.Type {
			case "string":
				str, ok := value.(string)
				if !ok {
					msg := "Must be a string"
					if rule.Messages != nil && rule.Messages.TypeMismatch != nil {
						msg = *rule.Messages.TypeMismatch
					}
					validationErrors = append(validationErrors, ValidationError{Field: field, Message: msg})
					continue
				}
				if rule.MinLength != nil && len(str) < *rule.MinLength {
					msg := "String too short"
					if rule.Messages != nil && rule.Messages.MinLength != nil {
						msg = *rule.Messages.MinLength
					}
					validationErrors = append(validationErrors, ValidationError{Field: field, Message: msg})
				}
				if rule.MaxLength != nil && len(str) > *rule.MaxLength {
					msg := "String too long"
					if rule.Messages != nil && rule.Messages.MaxLength != nil {
						msg = *rule.Messages.MaxLength
					}
					validationErrors = append(validationErrors, ValidationError{Field: field, Message: msg})
				}
				if rule.Regex != nil {
					re := regexp.MustCompile(*rule.Regex)
					if !re.MatchString(str) {
						msg := "Invalid format"
						if rule.Messages != nil && rule.Messages.Regex != nil {
							msg = *rule.Messages.Regex
						}
						validationErrors = append(validationErrors, ValidationError{Field: field, Message: msg})
					}
				}

			case "int":
				num, ok := value.(int)
				if !ok {
					msg := "Must be an integer"
					if rule.Messages != nil && rule.Messages.TypeMismatch != nil {
						msg = *rule.Messages.TypeMismatch
					}
					validationErrors = append(validationErrors, ValidationError{Field: field, Message: msg})
					continue
				}
				if rule.Min != nil && num < *rule.Min {
					msg := "Number too small"
					if rule.Messages != nil && rule.Messages.Min != nil {
						msg = *rule.Messages.Min
					}
					validationErrors = append(validationErrors, ValidationError{Field: field, Message: msg})
				}
				if rule.Max != nil && num > *rule.Max {
					msg := "Number too large"
					if rule.Messages != nil && rule.Messages.Max != nil {
						msg = *rule.Messages.Max
					}
					validationErrors = append(validationErrors, ValidationError{Field: field, Message: msg})
				}
			}

			// Apply custom validation function
			if rule.CheckWith != nil {
				if errorMsg := rule.CheckWith(value); errorMsg != nil {
					msg := *errorMsg
					if rule.Messages != nil && rule.Messages.CustomError != nil {
						msg = *rule.Messages.CustomError
					}
					validationErrors = append(validationErrors, ValidationError{Field: field, Message: msg})
				}
			}
		}
	}

	return ValidationResult{
		IsValid: len(validationErrors) == 0,
		Errors:  validationErrors,
	}
}
