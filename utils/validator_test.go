package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidCEP(t *testing.T) {
	tests := []struct {
		name     string
		cep      string
		expected bool
	}{
		{
			name:     "Valid CEP with hyphen",
			cep:      "01310-100",
			expected: true,
		},
		{
			name:     "Valid CEP without hyphen",
			cep:      "01310100",
			expected: true,
		},
		{
			name:     "Invalid CEP - too short",
			cep:      "0131010",
			expected: false,
		},
		{
			name:     "Invalid CEP - too long",
			cep:      "013101000",
			expected: false,
		},
		{
			name:     "Invalid CEP - contains letters",
			cep:      "01310a00",
			expected: false,
		},
		{
			name:     "Invalid CEP - empty string",
			cep:      "",
			expected: false,
		},
		{
			name:     "Invalid CEP - only spaces",
			cep:      "   ",
			expected: false,
		},
		{
			name:     "Valid CEP with spaces",
			cep:      " 01310-100 ",
			expected: true,
		},
		{
			name:     "Invalid CEP - wrong hyphen position",
			cep:      "013-10100",
			expected: false,
		},
		{
			name:     "Invalid CEP - multiple hyphens",
			cep:      "01310-1-00",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidCEP(tt.cep)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNormalizeCEP(t *testing.T) {
	tests := []struct {
		name     string
		cep      string
		expected string
	}{
		{
			name:     "CEP with hyphen should remove hyphen",
			cep:      "01310-100",
			expected: "01310100",
		},
		{
			name:     "CEP without hyphen should remain the same",
			cep:      "01310100",
			expected: "01310100",
		},
		{
			name:     "CEP with spaces should remove spaces",
			cep:      " 01310-100 ",
			expected: "01310100",
		},
		{
			name:     "CEP with multiple hyphens should remove all",
			cep:      "01-31-0-100",
			expected: "01310100",
		},
		{
			name:     "Empty CEP should return empty",
			cep:      "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NormalizeCEP(tt.cep)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFormatCEP(t *testing.T) {
	tests := []struct {
		name     string
		cep      string
		expected string
	}{
		{
			name:     "8-digit CEP should add hyphen",
			cep:      "01310100",
			expected: "01310-100",
		},
		{
			name:     "CEP with hyphen should maintain format",
			cep:      "01310-100",
			expected: "01310-100",
		},
		{
			name:     "CEP with spaces should format properly",
			cep:      " 01310100 ",
			expected: "01310-100",
		},
		{
			name:     "Invalid length CEP should return original",
			cep:      "123456",
			expected: "123456",
		},
		{
			name:     "Empty CEP should return empty",
			cep:      "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatCEP(tt.cep)
			assert.Equal(t, tt.expected, result)
		})
	}
} 