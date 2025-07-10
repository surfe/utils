package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsPersonalPhone(t *testing.T) {
	t.Parallel()

	tests := []struct {
		phone    string
		expected bool
	}{
		{"+33612345678", true},
		{"0612345678", true},
		{"1234567890", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.phone, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, tt.expected, IsPersonalPhone(tt.phone))
		})
	}
}
