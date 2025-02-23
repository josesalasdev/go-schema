package validator

type CustomValidator func(value interface{}) *string

type Rule struct {
	Type          string          `json:"type"`
	Required      bool            `json:"required"`
	Default       interface{}     `json:"default,omitempty"` // Nuevo campo
	MinLength     int             `json:"minlength,omitempty"`
	MaxLength     int             `json:"maxlength,omitempty"`
	Min           int             `json:"min,omitempty"`
	Max           int             `json:"max,omitempty"`
	Regex         string          `json:"regex,omitempty"`
	AllowedValues []interface{}   `json:"allowed_values,omitempty"`
	List          *Rule           `json:"list,omitempty"`
	Schema        *Schema         `json:"schema,omitempty"`
	CheckWith     CustomValidator `json:"-"`
	Messages      *ErrorMessages  `json:"messages,omitempty"`
}

type Schema map[string]Rule
