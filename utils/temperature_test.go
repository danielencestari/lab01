package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCelsiusToFahrenheit(t *testing.T) {
	tests := []struct {
		name     string
		celsius  float64
		expected float64
	}{
		{
			name:     "Zero Celsius should be 32 Fahrenheit",
			celsius:  0.0,
			expected: 32.0,
		},
		{
			name:     "100 Celsius should be 212 Fahrenheit",
			celsius:  100.0,
			expected: 212.0,
		},
		{
			name:     "25 Celsius should be 77 Fahrenheit",
			celsius:  25.0,
			expected: 77.0,
		},
		{
			name:     "Negative temperature -10 Celsius should be 14 Fahrenheit",
			celsius:  -10.0,
			expected: 14.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CelsiusToFahrenheit(tt.celsius)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCelsiusToKelvin(t *testing.T) {
	tests := []struct {
		name     string
		celsius  float64
		expected float64
	}{
		{
			name:     "Zero Celsius should be 273 Kelvin",
			celsius:  0.0,
			expected: 273.0,
		},
		{
			name:     "25 Celsius should be 298 Kelvin",
			celsius:  25.0,
			expected: 298.0,
		},
		{
			name:     "100 Celsius should be 373 Kelvin",
			celsius:  100.0,
			expected: 373.0,
		},
		{
			name:     "Negative temperature -10 Celsius should be 263 Kelvin",
			celsius:  -10.0,
			expected: 263.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CelsiusToKelvin(tt.celsius)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestConvertTemperatures(t *testing.T) {
	tests := []struct {
		name          string
		celsius       float64
		expectedC     float64
		expectedF     float64
		expectedK     float64
	}{
		{
			name:      "Convert 25 Celsius to all scales",
			celsius:   25.0,
			expectedC: 25.0,
			expectedF: 77.0,
			expectedK: 298.0,
		},
		{
			name:      "Convert 0 Celsius to all scales",
			celsius:   0.0,
			expectedC: 0.0,
			expectedF: 32.0,
			expectedK: 273.0,
		},
		{
			name:      "Convert -10 Celsius to all scales",
			celsius:   -10.0,
			expectedC: -10.0,
			expectedF: 14.0,
			expectedK: 263.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, f, k := ConvertTemperatures(tt.celsius)
			assert.Equal(t, tt.expectedC, c)
			assert.Equal(t, tt.expectedF, f)
			assert.Equal(t, tt.expectedK, k)
		})
	}
} 