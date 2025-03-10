package validator

import (
	"testing"
)

// Helper function for string pointers
func strPtr(s string) *string {
	return &s
}

func TestNumberValidation(t *testing.T) {
	schema := Schema{
		"user": {
			Type: "int",
			Max:  1,
		},
	}

	validData := map[string]interface{}{
		"user": 3,
	}

	result := Validate(validData, schema)
	if result.IsValid {
		t.Errorf("Expected valid data, got errors: %v", result.Errors)
	}
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
	if resultInvalid.IsValid {
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
	if resultInvalid.IsValid {
		t.Errorf("Expected validation to fail, but passed")
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

// TestBasicValidation tests the basic validation functionality
func TestBasicValidation(t *testing.T) {
	t.Run("RequiredField", func(t *testing.T) {
		schema := Schema{
			"name": {Type: "string", Required: true},
		}

		validData := map[string]interface{}{
			"name": "John",
		}

		invalidData := map[string]interface{}{
			"age": 25, // Missing required field
		}

		// Test valid data
		result := Validate(validData, schema)
		if !result.IsValid {
			t.Errorf("Expected valid data, got errors: %v", result.Errors)
		}

		// Test invalid data
		result = Validate(invalidData, schema)
		if result.IsValid {
			t.Errorf("Expected validation to fail for missing required field")
		}
	})

	t.Run("TypeValidation", func(t *testing.T) {
		schema := Schema{
			"name":   {Type: "string"},
			"age":    {Type: "int"},
			"active": {Type: "bool"},
		}

		validData := map[string]interface{}{
			"name":   "John",
			"age":    30,
			"active": true,
		}

		invalidData := map[string]interface{}{
			"name":   123, // Wrong type
			"age":    "thirty",
			"active": "yes",
		}

		// Test valid data
		result := Validate(validData, schema)
		if !result.IsValid {
			t.Errorf("Expected valid data, got errors: %v", result.Errors)
		}

		// Test invalid data
		result = Validate(invalidData, schema)
		if result.IsValid {
			t.Errorf("Expected validation to fail for type mismatch")
		}
	})
}

// TestNumericValidation tests numeric validation rules
func TestNumericValidation(t *testing.T) {
	tests := []struct {
		name       string
		schema     Schema
		data       map[string]interface{}
		shouldPass bool
	}{
		{
			name: "Valid Int Within Range",
			schema: Schema{
				"age": {Type: "int", Min: 18, Max: 100},
			},
			data: map[string]interface{}{
				"age": 30,
			},
			shouldPass: true,
		},
		{
			name: "Int Below Minimum",
			schema: Schema{
				"age": {Type: "int", Min: 18, Max: 100},
			},
			data: map[string]interface{}{
				"age": 15,
			},
			shouldPass: false,
		},
		{
			name: "Int Above Maximum",
			schema: Schema{
				"age": {Type: "int", Min: 18, Max: 100},
			},
			data: map[string]interface{}{
				"age": 120,
			},
			shouldPass: false,
		},
		{
			name: "Valid Float Within Range",
			schema: Schema{
				"price": {Type: "float", Min: 0.1, Max: 999.99},
			},
			data: map[string]interface{}{
				"price": 149.99,
			},
			shouldPass: true,
		},
		{
			name: "Float Below Minimum",
			schema: Schema{
				"price": {Type: "float", Min: 0.1, Max: 999.99},
			},
			data: map[string]interface{}{
				"price": 0.05,
			},
			shouldPass: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Validate(tt.data, tt.schema)
			if result.IsValid != tt.shouldPass {
				if tt.shouldPass {
					t.Errorf("Expected validation to pass, but got errors: %v", result.Errors)
				} else {
					t.Errorf("Expected validation to fail, but it passed")
				}
			}
		})
	}
}

// TestStringValidation tests string validation rules
func TestStringValidation(t *testing.T) {
	tests := []struct {
		name       string
		schema     Schema
		data       map[string]interface{}
		shouldPass bool
	}{
		{
			name: "Valid String Length",
			schema: Schema{
				"username": {Type: "string", MinLength: 3, MaxLength: 20},
			},
			data: map[string]interface{}{
				"username": "johndoe",
			},
			shouldPass: true,
		},
		{
			name: "String Too Short",
			schema: Schema{
				"username": {Type: "string", MinLength: 3, MaxLength: 20},
			},
			data: map[string]interface{}{
				"username": "jo",
			},
			shouldPass: false,
		},
		{
			name: "String Too Long",
			schema: Schema{
				"username": {Type: "string", MinLength: 3, MaxLength: 20},
			},
			data: map[string]interface{}{
				"username": "thisusernameiswaytoolong",
			},
			shouldPass: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Validate(tt.data, tt.schema)
			if result.IsValid != tt.shouldPass {
				if tt.shouldPass {
					t.Errorf("Expected validation to pass, but got errors: %v", result.Errors)
				} else {
					t.Errorf("Expected validation to fail, but it passed")
				}
			}
		})
	}
}

// TestNestedStructures tests validation of nested data structures
func TestNestedStructures(t *testing.T) {
	t.Run("ListValidation", func(t *testing.T) {
		schema := Schema{
			"tags": {
				Type: "list",
				List: &Rule{
					Type: "string",
				},
			},
		}

		validData := map[string]interface{}{
			"tags": []interface{}{"golang", "schema"},
		}

		invalidData := map[string]interface{}{
			"tags": []interface{}{123, "schema"},
		}

		// Test valid data
		result := Validate(validData, schema)
		if !result.IsValid {
			t.Errorf("Expected valid data, got errors: %v", result.Errors)
		}

		// Test invalid data
		result = Validate(invalidData, schema)
		if result.IsValid {
			t.Errorf("Expected validation to fail for invalid list item type")
		}
	})

	t.Run("MapValidation", func(t *testing.T) {
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

		invalidData := map[string]interface{}{
			"user": map[string]interface{}{
				"age": 16, // Missing required name & age too young
			},
		}

		// Test valid data
		result := Validate(validData, schema)
		if !result.IsValid {
			t.Errorf("Expected valid data, got errors: %v", result.Errors)
		}

		// Test invalid data
		result = Validate(invalidData, schema)
		if result.IsValid {
			t.Errorf("Expected validation to fail for invalid nested map")
		}
	})

	t.Run("DeepNesting", func(t *testing.T) {
		schema := Schema{
			"data": {
				Type: "map",
				Schema: &Schema{
					"users": {
						Type: "list",
						List: &Rule{
							Type: "map",
							Schema: &Schema{
								"id": {Type: "int", Required: true},
								"details": {
									Type: "map",
									Schema: &Schema{
										"name":  {Type: "string", Required: true},
										"email": {Type: "string"},
									},
								},
							},
						},
					},
				},
			},
		}

		validData := map[string]interface{}{
			"data": map[string]interface{}{
				"users": []interface{}{
					map[string]interface{}{
						"id": 1,
						"details": map[string]interface{}{
							"name":  "John",
							"email": "john@example.com",
						},
					},
				},
			},
		}

		// Test valid data
		result := Validate(validData, schema)
		if !result.IsValid {
			t.Errorf("Expected valid data, got errors: %v", result.Errors)
		}
	})
}

// TestCustomMessages tests custom error message functionality
func TestCustomMessages(t *testing.T) {
	schema := Schema{
		"name": {
			Type:     "string",
			Required: true,
			Messages: &Messages{
				Required: strPtr("The name field is mandatory"),
			},
		},
		"age": {
			Type: "int",
			Min:  18,
			Messages: &Messages{
				Range: strPtr("You must be at least 18 years old"),
			},
		},
	}

	data := map[string]interface{}{}

	result := Validate(data, schema)
	if result.IsValid {
		t.Errorf("Expected validation to fail, but passed")
	}

	expectedMsg := "The name field is mandatory"
	var foundExpectedMsg bool
	for _, err := range result.Errors {
		if err.Field == "name" && err.Message == expectedMsg {
			foundExpectedMsg = true
			break
		}
	}

	if !foundExpectedMsg {
		t.Errorf("Expected error message '%s', but didn't find it in: %v", expectedMsg, result.Errors)
	}
}

// TestSchemaValidation tests schema self-validation functionality
func TestSchemaValidation(t *testing.T) {
	tests := []struct {
		name        string
		schema      Schema
		expectError bool
	}{
		{
			name: "Valid Schema",
			schema: Schema{
				"name": {Type: "string", Default: "John"},
				"age":  {Type: "int", Min: 18, Max: 99},
				"tags": {Type: "list", List: &Rule{Type: "string"}},
				"meta": {Type: "map", Schema: &Schema{"version": {Type: "string"}}},
			},
			expectError: false,
		},
		{
			name: "Invalid Type",
			schema: Schema{
				"age": {Type: "number"}, // "number" isn't supported
			},
			expectError: true,
		},
		{
			name: "Invalid Field Name",
			schema: Schema{
				"invalid key": {Type: "string"}, // Contains space
			},
			expectError: true,
		},
		{
			name: "Type Mismatch in Default Value",
			schema: Schema{
				"active": {Type: "bool", Default: "yes"}, // "yes" isn't a bool
			},
			expectError: true,
		},
		{
			name: "Min/Max on Non-numeric Type",
			schema: Schema{
				"username": {Type: "string", Min: 3, Max: 10}, // Min/Max only for numbers
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateSchema(tt.schema)
			if (err == nil) == tt.expectError {
				if tt.expectError {
					t.Errorf("Expected error for invalid schema, but got nil")
				} else {
					t.Errorf("Expected valid schema, got error: %v", err)
				}
			}
		})
	}
}
