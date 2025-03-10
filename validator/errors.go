package validator

type ErrorMessages struct {
	Required     *string `json:"required,omitempty"`
	TypeMismatch *string `json:"type_mismatch,omitempty"`
	MinLength    *string `json:"min_length,omitempty"`
	MaxLength    *string `json:"max_length,omitempty"`
	Min          *string `json:"min,omitempty"`
	Max          *string `json:"max,omitempty"`
	Regex        *string `json:"regex,omitempty"`
	CustomError  *string `json:"custom_error,omitempty"`
}
