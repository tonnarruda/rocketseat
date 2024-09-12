package api

import (
	"net/http"
	"net/http/httptest"
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
func TestSendJSON(t *testing.T) {
	tests := []struct {
		name     string
		resp     Response
		status   int
		expected string
	}{
		{
			name:     "Valid response with data",
			resp:     Response{Data: "testdata"},
			status:   http.StatusOK,
			expected: `{"data":"testdata"}`,
		},
		{
			name:     "Valid response with error",
			resp:     Response{Error: "testerror"},
			status:   http.StatusBadRequest,
			expected: `{"error":"testerror"}`,
		},
		{
			name:     "Empty response",
			resp:     Response{},
			status:   http.StatusNoContent,
			expected: `{}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			SendJSON(w, tt.resp, tt.status)

			if w.Code != tt.status {
				t.Errorf("Expected status %v, but got %v", tt.status, w.Code)
			}

			if w.Body.String() != tt.expected {
				t.Errorf("Expected body %v, but got %v", tt.expected, w.Body.String())
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
