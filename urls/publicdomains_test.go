package urls

import "testing"

func TestIsPublicDomain(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		domain string
		want   bool
	}{
		{
			name:   "Detects social media domains",
			domain: "instagram.com",
			want:   true,
		},
		{
			name:   "Detects non-public domains",
			domain: "sonereker.com",
			want:   false,
		},
		{
			name:   "Detects website builders",
			domain: "wix.com",
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := IsPublicDomain(tt.domain); got != tt.want {
				t.Errorf("IsPublicDomain() = %v, want %v", got, tt.want)
			}
		})
	}
}
