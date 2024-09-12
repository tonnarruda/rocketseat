package api

import (
	"testing"
)

func TestGencode(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected bool
	}{
		{
			name:     "Valid code",
			code:     gencode(),
			expected: true,
		},
		{
			name:     "Invalid code with special characters",
			code:     "ABCD@#@", // Um código inválido fictício
			expected: false,
		},
		{
			name:     "Invalid code with incorrect length",
			code:     "ABCDE", // Código com comprimento incorreto
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid := len(tt.code) == 8 // Verifica comprimento

			if isValid {
				for _, char := range tt.code {
					if !isValidCharacter(char) {
						isValid = false
						t.Errorf("Invalid character in code: %c", char)
						break
					}
				}
			}

			if isValid != tt.expected {
				t.Errorf("Expected validity %v, but got %v", tt.expected, isValid)
			}
		})
	}
}

// isValidCharacter verifica se um caractere é válido
func isValidCharacter(char rune) bool {
	for _, validChar := range characters {
		if char == validChar {
			return true
		}
	}
	return false
}
