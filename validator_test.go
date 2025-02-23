package validator

import "testing"

func TestValidate(t *testing.T) {
	schema := Schema{
		"nombre": {Type: "string", MinLength: intPtr(2), Required: true},
		"edad":   {Type: "int", Min: intPtr(18), Max: intPtr(99), Required: true},
	}

	data := map[string]interface{}{
		"nombre": "A",
		"edad":   17,
	}

	result := Validate(data, schema)

	if result.IsValid {
		t.Errorf("Se esperaba error, pero la validación pasó")
	}

	if len(result.Errors) != 2 {
		t.Errorf("Se esperaban 2 errores, pero se encontraron %d", len(result.Errors))
	}
}
