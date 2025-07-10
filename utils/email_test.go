package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsPersonalEmail(t *testing.T) {
	t.Parallel()

	tests := []struct {
		email    string
		expected bool
	}{
		{"test@gmail.com", true},
		{"user@surfe.com", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.email, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expected, IsPersonalEmail(tt.email))
		})
	}
}

func TestIsPersonalEmailDomain(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		domain   string
		expected bool
	}{
		{"Gmail email", "gmail.com", true},
		{"Company email (@surfe.com)", "surfe.com", false},
		{"Empty email", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.expected, IsPersonalEmailDomain(tt.domain))
		})
	}
}
