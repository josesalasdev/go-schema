package validator

import (
	"regexp"
)

// Schema defines validation rules for a set of fields
type Schema map[string]Rule

// Rule defines validation rules for a single field
type Rule struct {
	Type         string         `json:"type"`
	Required     bool           `json:"required,omitempty"`
	Default      interface{}    `json:"default,omitempty"`
	Min          float64        `json:"min,omitempty"`
	Max          float64        `json:"max,omitempty"`
	MinLength    int            `json:"min_length,omitempty"`
	MaxLength    int            `json:"max_length,omitempty"`
	Regex        *regexp.Regexp `json:"-"`
	RegexPattern string         `json:"regex,omitempty"`
	List         *Rule          `json:"list,omitempty"`
	Schema       *Schema        `json:"schema,omitempty"`
	Messages     *Messages      `json:"messages,omitempty"`
}

// Messages provides customized error messages
type Messages struct {
	Required     *string `json:"required,omitempty"`
	TypeMismatch *string `json:"type_mismatch,omitempty"`
	Range        *string `json:"range,omitempty"`
	Length       *string `json:"length,omitempty"`
	Pattern      *string `json:"pattern,omitempty"`
}

// ValidationResult represents the result of validation
type ValidationResult struct {
	IsValid bool
	Errors  []ValidationError
}

// ValidationError represents a single validation error
type ValidationError struct {
	Field   string
	Message string
}
