package validator

type Rule struct {
	Type          string        `json:"type"`                     // Tipo de dato: string, int, bool, list...
	Required      bool          `json:"required"`                 // ¿Es obligatorio?
	MinLength     *int          `json:"minlength,omitempty"`      // Longitud mínima (para strings)
	MaxLength     *int          `json:"maxlength,omitempty"`      // Longitud máxima
	Min           *int          `json:"min,omitempty"`            // Mínimo (para números)
	Max           *int          `json:"max,omitempty"`            // Máximo (para números)
	Regex         *string       `json:"regex,omitempty"`          // Validación con regex
	AllowedValues []interface{} `json:"allowed_values,omitempty"` // Lista de valores permitidos
}

type Schema map[string]Rule
