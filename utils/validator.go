package utils

import (
	"regexp"
	"strings"
)

// IsValidCEP valida se o CEP tem formato correto (8 dígitos com ou sem hífen)
func IsValidCEP(cep string) bool {
	// Remove espaços em branco
	cep = strings.TrimSpace(cep)
	
	// Verifica formato: 12345678 ou 12345-678
	cepRegex := regexp.MustCompile(`^\d{5}-?\d{3}$`)
	return cepRegex.MatchString(cep)
}

// NormalizeCEP remove hífen e formata o CEP para 8 dígitos
func NormalizeCEP(cep string) string {
	// Remove hífen e espaços
	normalized := strings.ReplaceAll(cep, "-", "")
	normalized = strings.TrimSpace(normalized)
	return normalized
}

// FormatCEP adiciona hífen no CEP (formato: 12345-678)
func FormatCEP(cep string) string {
	normalized := NormalizeCEP(cep)
	if len(normalized) == 8 {
		return normalized[:5] + "-" + normalized[5:]
	}
	return cep
} 