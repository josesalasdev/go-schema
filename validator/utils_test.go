package validator

import (
	"math"
	"testing"
)

func TestExtractIntValue(t *testing.T) {
	tests := []struct {
		name      string
		input     interface{}
		wantValue int64
		wantOk    bool
	}{
		{
			name:      "Extract from int",
			input:     42,
			wantValue: 42,
			wantOk:    true,
		},
		{
			name:      "Extract from int8",
			input:     int8(127),
			wantValue: 127,
			wantOk:    true,
		},
		{
			name:      "Extract from int16",
			input:     int16(32000),
			wantValue: 32000,
			wantOk:    true,
		},
		{
			name:      "Extract from int32",
			input:     int32(2147483647),
			wantValue: 2147483647,
			wantOk:    true,
		},
		{
			name:      "Extract from int64",
			input:     int64(9223372036854775807),
			wantValue: 9223372036854775807,
			wantOk:    true,
		},
		{
			name:      "Extract from uint",
			input:     uint(42),
			wantValue: 42,
			wantOk:    true,
		},
		{
			name:      "Extract from uint8",
			input:     uint8(255),
			wantValue: 255,
			wantOk:    true,
		},
		{
			name:      "Extract from uint16",
			input:     uint16(65535),
			wantValue: 65535,
			wantOk:    true,
		},
		{
			name:      "Extract from uint32",
			input:     uint32(4294967295),
			wantValue: 4294967295,
			wantOk:    true,
		},
		{
			name:      "Extract from float32 (whole number)",
			input:     float32(100.0),
			wantValue: 100,
			wantOk:    true,
		},
		{
			name:      "Extract from float64 (whole number)",
			input:     float64(100.0),
			wantValue: 100,
			wantOk:    true,
		},
		{
			name:      "Fail from float32 (fractional)",
			input:     float32(100.5),
			wantValue: 0,
			wantOk:    false,
		},
		{
			name:      "Fail from float64 (fractional)",
			input:     float64(100.5),
			wantValue: 0,
			wantOk:    false,
		},
		{
			name:      "Fail from string",
			input:     "42",
			wantValue: 0,
			wantOk:    false,
		},
		{
			name:      "Fail from bool",
			input:     true,
			wantValue: 0,
			wantOk:    false,
		},
		{
			name:      "Fail from nil",
			input:     nil,
			wantValue: 0,
			wantOk:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotValue, gotOk := extractIntValue(tt.input)
			if gotValue != tt.wantValue || gotOk != tt.wantOk {
				t.Errorf("extractIntValue(%v) = (%v, %v), want (%v, %v)",
					tt.input, gotValue, gotOk, tt.wantValue, tt.wantOk)
			}
		})
	}
}

func TestExtractFloatValue(t *testing.T) {
	tests := []struct {
		name      string
		input     interface{}
		wantValue float64
		wantOk    bool
	}{
		{
			name:      "Extract from float64",
			input:     float64(3.14159),
			wantValue: 3.14159,
			wantOk:    true,
		},
		{
			name:      "Extract from float32",
			input:     float32(3.14),
			wantValue: float64(float32(3.14)), // Account for precision loss
			wantOk:    true,
		},
		{
			name:      "Extract from int",
			input:     42,
			wantValue: 42.0,
			wantOk:    true,
		},
		{
			name:      "Extract from int8",
			input:     int8(127),
			wantValue: 127.0,
			wantOk:    true,
		},
		{
			name:      "Extract from int16",
			input:     int16(32000),
			wantValue: 32000.0,
			wantOk:    true,
		},
		{
			name:      "Extract from int32",
			input:     int32(2147483647),
			wantValue: 2147483647.0,
			wantOk:    true,
		},
		{
			name:      "Extract from int64",
			input:     int64(9223372036854775807),
			wantValue: 9223372036854775807.0,
			wantOk:    true,
		},
		{
			name:      "Extract from uint",
			input:     uint(42),
			wantValue: 42.0,
			wantOk:    true,
		},
		{
			name:      "Extract from uint8",
			input:     uint8(255),
			wantValue: 255.0,
			wantOk:    true,
		},
		{
			name:      "Extract from uint16",
			input:     uint16(65535),
			wantValue: 65535.0,
			wantOk:    true,
		},
		{
			name:      "Extract from uint32",
			input:     uint32(4294967295),
			wantValue: 4294967295.0,
			wantOk:    true,
		},
		{
			name:      "Fail from string",
			input:     "3.14",
			wantValue: 0,
			wantOk:    false,
		},
		{
			name:      "Fail from bool",
			input:     true,
			wantValue: 0,
			wantOk:    false,
		},
		{
			name:      "Fail from nil",
			input:     nil,
			wantValue: 0,
			wantOk:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotValue, gotOk := extractFloatValue(tt.input)
			if math.Abs(gotValue-tt.wantValue) > 0.0001 || gotOk != tt.wantOk {
				t.Errorf("extractFloatValue(%v) = (%v, %v), want (%v, %v)",
					tt.input, gotValue, gotOk, tt.wantValue, tt.wantOk)
			}
		})
	}
}

func TestIsValidJSONKey(t *testing.T) {
	tests := []struct {
		name  string
		key   string
		valid bool
	}{
		{
			name:  "Valid simple key",
			key:   "name",
			valid: true,
		},
		{
			name:  "Valid key with numbers",
			key:   "user123",
			valid: true,
		},
		{
			name:  "Valid key with special chars",
			key:   "user_name-123",
			valid: true,
		},
		{
			name:  "Invalid key with space",
			key:   "user name",
			valid: false,
		},
		{
			name:  "Invalid key with tab",
			key:   "user\tname",
			valid: false,
		},
		{
			name:  "Invalid key with newline",
			key:   "user\nname",
			valid: false,
		},
		{
			name:  "Invalid key with control character",
			key:   "user\u0000name",
			valid: false,
		},
		{
			name:  "Empty key",
			key:   "",
			valid: true, // Empty keys are technically valid in JSON
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isValidJSONKey(tt.key)
			if got != tt.valid {
				t.Errorf("isValidJSONKey(%q) = %v, want %v", tt.key, got, tt.valid)
			}
		})
	}
}
