package validator

import (
	"regexp"
	"testing"
)

func TestValidationError(t *testing.T) {
	tests := []struct {
		name     string
		err      ValidationError
		expected string
	}{
		{
			name: "Simple error",
			err: ValidationError{
				Field:   "username",
				Message: "is required",
			},
			expected: "username: is required",
		},
		{
			name: "Nested field error",
			err: ValidationError{
				Field:   "user.address.zipcode",
				Message: "must be numeric",
			},
			expected: "user.address.zipcode: must be numeric",
		},
		{
			name: "Array field error",
			err: ValidationError{
				Field:   "items[2]",
				Message: "invalid value",
			},
			expected: "items[2]: invalid value",
		},
		{
			name: "Empty field",
			err: ValidationError{
				Field:   "",
				Message: "general error",
			},
			expected: ": general error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.err.Error()
			if result != tt.expected {
				t.Errorf("Error() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestSchemaCreation(t *testing.T) {
	// Test creating a schema with various rule types
	schema := Schema{
		"name": {
			Type:      "string",
			Required:  true,
			MinLength: 2,
			MaxLength: 50,
		},
		"age": {
			Type: "int",
			Min:  18,
			Max:  120,
		},
		"email": {
			Type:         "string",
			Required:     true,
			RegexPattern: `^[^@]+@[^@]+\.[^@]+$`,
			Regex:        regexp.MustCompile(`^[^@]+@[^@]+\.[^@]+$`),
		},
		"tags": {
			Type: "list",
			List: &Rule{
				Type: "string",
			},
		},
		"address": {
			Type: "map",
			Schema: &Schema{
				"street": {Type: "string", Required: true},
				"city":   {Type: "string", Required: true},
				"zip":    {Type: "string", Required: true},
			},
		},
		"active": {
			Type:    "bool",
			Default: true,
		},
	}

	// Test that schema creation works as expected
	if len(schema) != 6 {
		t.Errorf("Expected schema with 6 fields, got %d", len(schema))
	}

	// Test individual field rules
	if !schema["name"].Required {
		t.Errorf("Expected 'name' to be required")
	}

	if schema["age"].Min != 18 {
		t.Errorf("Expected 'age' minimum to be 18, got %f", schema["age"].Min)
	}

	if schema["email"].RegexPattern != `^[^@]+@[^@]+\.[^@]+$` {
		t.Errorf("Expected email regex pattern didn't match")
	}

	if schema["tags"].List == nil || schema["tags"].List.Type != "string" {
		t.Errorf("Expected 'tags' to be a list of strings")
	}

	if schema["address"].Schema == nil || len(*schema["address"].Schema) != 3 {
		t.Errorf("Expected 'address' to be a map with 3 fields")
	}

	if schema["active"].Default != true {
		t.Errorf("Expected 'active' default to be true")
	}
}

func TestValidationResultCreation(t *testing.T) {
	// Valid result
	validResult := ValidationResult{
		IsValid: true,
		Errors:  []ValidationError{},
	}

	if !validResult.IsValid {
		t.Errorf("Expected result to be valid")
	}

	if len(validResult.Errors) != 0 {
		t.Errorf("Expected valid result to have no errors")
	}

	// Invalid result
	invalidResult := ValidationResult{
		IsValid: false,
		Errors: []ValidationError{
			{Field: "name", Message: "is required"},
			{Field: "age", Message: "must be at least 18"},
		},
	}

	if invalidResult.IsValid {
		t.Errorf("Expected result to be invalid")
	}

	if len(invalidResult.Errors) != 2 {
		t.Errorf("Expected invalid result to have 2 errors, got %d", len(invalidResult.Errors))
	}

	// Check error messages
	if invalidResult.Errors[0].Field != "name" ||
		invalidResult.Errors[0].Message != "is required" {
		t.Errorf("First error doesn't match expected values")
	}

	if invalidResult.Errors[1].Field != "age" ||
		invalidResult.Errors[1].Message != "must be at least 18" {
		t.Errorf("Second error doesn't match expected values")
	}
}

func TestMessagesCreation(t *testing.T) {
	requiredMsg := "This field is absolutely required"
	typeMsg := "Wrong type provided"
	rangeMsg := "Value out of acceptable range"

	messages := Messages{
		Required:     &requiredMsg,
		TypeMismatch: &typeMsg,
		Range:        &rangeMsg,
	}

	if *messages.Required != requiredMsg {
		t.Errorf("Required message doesn't match")
	}

	if *messages.TypeMismatch != typeMsg {
		t.Errorf("Type mismatch message doesn't match")
	}

	if *messages.Range != rangeMsg {
		t.Errorf("Range message doesn't match")
	}

	if messages.Length != nil {
		t.Errorf("Length message should be nil")
	}

	if messages.Pattern != nil {
		t.Errorf("Pattern message should be nil")
	}
}

func TestRuleWithCustomMessages(t *testing.T) {
	requiredMsg := "Name is required"
	lengthMsg := "Name must be between 2 and 50 characters"

	rule := Rule{
		Type:      "string",
		Required:  true,
		MinLength: 2,
		MaxLength: 50,
		Messages: &Messages{
			Required: &requiredMsg,
			Length:   &lengthMsg,
		},
	}

	if rule.Required != true {
		t.Errorf("Rule should be required")
	}

	if rule.Type != "string" {
		t.Errorf("Rule type should be string")
	}

	if rule.Messages == nil {
		t.Errorf("Messages should not be nil")
	}

	if *rule.Messages.Required != requiredMsg {
		t.Errorf("Required message doesn't match")
	}

	if *rule.Messages.Length != lengthMsg {
		t.Errorf("Length message doesn't match")
	}
}
