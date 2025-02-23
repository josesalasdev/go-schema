package validator

import (
	"testing"
)

// Custom validation function (must be even)
func validateEvenNumber(value interface{}) *string {
	num, ok := value.(int)
	if !ok || num%2 != 0 {
		msg := "Must be an even number"
		return &msg
	}
	return nil
}

func TestListValidation(t *testing.T) {
	schema := Schema{
		"tags": {
			Type: "list",
			Items: &Rule{
				Type: "string",
			},
		},
	}

	data := map[string]interface{}{
		"tags": []interface{}{"golang", "schema"},
	}

	result := Validate(data, schema)
	if !result.IsValid {
		t.Errorf("Expected valid data, got errors: %v", result.Errors)
	}

	dataInvalid := map[string]interface{}{
		"tags": []interface{}{123, "schema"},
	}

	resultInvalid := Validate(dataInvalid, schema)
	if !resultInvalid.IsValid {
		t.Errorf("Expected validation to fail, but passed")
	}
}

func TestMapValidation(t *testing.T) {
	schema := Schema{
		"user": {
			Type: "map",
			Schema: &Schema{
				"name": {Type: "string", Required: true},
				"age":  {Type: "int", Min: intPtr(18)},
			},
		},
	}

	validData := map[string]interface{}{
		"user": map[string]interface{}{
			"name": "John",
			"age":  25,
		},
	}

	result := Validate(validData, schema)
	if !result.IsValid {
		t.Errorf("Expected valid data, got errors: %v", result.Errors)
	}

	invalidData := map[string]interface{}{
		"user": map[string]interface{}{
			"age": 16, // Too young
		},
	}

	resultInvalid := Validate(invalidData, schema)
	if !resultInvalid.IsValid {
		t.Errorf("Expected validation to fail, but passed")
	}
}

func TestCustomValidation(t *testing.T) {
	schema := Schema{
		"age": {
			Type:      "int",
			CheckWith: validateEvenNumber,
		},
	}

	validData := map[string]interface{}{
		"age": 24,
	}

	result := Validate(validData, schema)
	if !result.IsValid {
		t.Errorf("Expected valid data, got errors: %v", result.Errors)
	}

	invalidData := map[string]interface{}{
		"age": 25, // Odd number
	}

	resultInvalid := Validate(invalidData, schema)
	if resultInvalid.IsValid {
		t.Errorf("Expected validation to fail, but passed")
	}
}

func TestCustomMessages(t *testing.T) {
	schema := Schema{
		"name": {
			Type:     "string",
			Required: true,
			Messages: &ErrorMessages{
				Required: strPtr("The name field is mandatory"),
			},
		},
	}

	data := map[string]interface{}{}

	result := Validate(data, schema)
	if result.IsValid {
		t.Errorf("Expected validation to fail, but passed")
	}

	expectedMsg := "The name field is mandatory"
	if result.Errors[0].Message != expectedMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedMsg, result.Errors[0].Message)
	}
}
