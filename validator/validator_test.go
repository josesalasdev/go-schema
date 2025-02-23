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
			List: &Rule{
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
				"age":  {Type: "int", Min: 18},
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
	if !resultInvalid.IsValid {
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

func TestValidateSchema(t *testing.T) {
	// ✅ Caso 1: Schema válido
	validSchema := Schema{
		"name": {Type: "string", Default: "John"},
		"age":  {Type: "int", Min: 18, Max: 99},
		"tags": {Type: "list", List: &Rule{Type: "string"}},
		"meta": {Type: "map", Schema: &Schema{"version": {Type: "string"}}},
	}

	if err := ValidateSchema(validSchema); err != nil {
		t.Errorf("Expected valid schema, got error: %v", err)
	}

	// ❌ Caso 2: Tipo inválido
	invalidTypeSchema := Schema{
		"age": {Type: "number"}, // "number" no está soportado
	}

	if err := ValidateSchema(invalidTypeSchema); err == nil {
		t.Errorf("Expected error for invalid type, but got nil")
	}

	// ❌ Caso 3: Nombre de campo inválido
	invalidKeySchema := Schema{
		"invalid key": {Type: "string"}, // Contiene espacio
	}

	if err := ValidateSchema(invalidKeySchema); err == nil {
		t.Errorf("Expected error for invalid field name, but got nil")
	}

	// ❌ Caso 4: Default value con tipo incorrecto
	invalidDefaultSchema := Schema{
		"active": {Type: "bool", Default: "yes"}, // "yes" no es bool
	}

	if err := ValidateSchema(invalidDefaultSchema); err == nil {
		t.Errorf("Expected error for mismatched default value, but got nil")
	}

	// ❌ Caso 5: Min/Max en tipo no numérico
	invalidMinMaxSchema := Schema{
		"username": {Type: "string", Min: 3, Max: 10}, // Min/Max solo en números
	}

	if err := ValidateSchema(invalidMinMaxSchema); err == nil {
		t.Errorf("Expected error for min/max in non-numeric field, but got nil")
	}

	// ❌ Caso 6: Lista con esquema inválido
	invalidListSchema := Schema{
		"list": {Type: "list", List: &Rule{Type: "invalid_type"}},
	}

	if err := ValidateSchema(invalidListSchema); err == nil {
		t.Errorf("Expected error for invalid list schema, but got nil")
	}

	// ❌ Caso 7: Mapa con esquema inválido
	invalidMapSchema := Schema{
		"data": {Type: "map", Schema: &Schema{"key": {Type: "unknown"}}},
	}

	if err := ValidateSchema(invalidMapSchema); err == nil {
		t.Errorf("Expected error for invalid map schema, but got nil")
	}
}
