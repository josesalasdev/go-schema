package validator

type CustomValidator func(value interface{}) *string

type Rule struct {
	Type          string          `json:"type"`
	Required      bool            `json:"required"`
	MinLength     *int            `json:"minlength,omitempty"`
	MaxLength     *int            `json:"maxlength,omitempty"`
	Min           *int            `json:"min,omitempty"`
	Max           *int            `json:"max,omitempty"`
	Regex         *string         `json:"regex,omitempty"`
	AllowedValues []interface{}   `json:"allowed_values,omitempty"`
	Items         *Rule           `json:"items,omitempty"`    // For lists
	Schema        *Schema         `json:"schema,omitempty"`   // For maps
	CheckWith     CustomValidator `json:"-"`                  // Custom validation function
	Messages      *ErrorMessages  `json:"messages,omitempty"` // Custom error messages
}

type Schema map[string]Rule
