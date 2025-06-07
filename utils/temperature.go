package utils

// CelsiusToFahrenheit converte temperatura de Celsius para Fahrenheit
// Fórmula: F = C * 1.8 + 32
func CelsiusToFahrenheit(celsius float64) float64 {
	return celsius*1.8 + 32
}

// CelsiusToKelvin converte temperatura de Celsius para Kelvin  
// Fórmula: K = C + 273
func CelsiusToKelvin(celsius float64) float64 {
	return celsius + 273
}

// ConvertTemperatures converte uma temperatura em Celsius para todas as escalas
func ConvertTemperatures(celsius float64) (float64, float64, float64) {
	fahrenheit := CelsiusToFahrenheit(celsius)
	kelvin := CelsiusToKelvin(celsius)
	return celsius, fahrenheit, kelvin
} 