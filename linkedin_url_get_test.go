package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Org struct {
	Name  string
	LIURL string
}

func (o Org) GetLinkedinURL(key string) string {
	return o.LIURL
}

func TestRemoveWithEmptyOrNotEqualLinkedinURL_NonEmptyLinkedInURL(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		orgs []Org
		want []Org
	}{
		{
			name: "Given an empty array, we return an empty array",
			orgs: []Org{},
			want: []Org{},
		},
		{
			name: "Given an array with orgs having different LIURL with different formats, we return an empty array",
			orgs: []Org{
				{Name: "Org1", LIURL: "https://www.linkedin.com/company/allbirds/"}, // https + www + trailing /
				{Name: "Org2", LIURL: "http://www.linkedin.com/company/allbirds/"},  // http + www + trailing /
				{Name: "Org3", LIURL: "https://linkedin.com/company/allbirds/"},     // https + trailing /
				{Name: "Org4", LIURL: "http://linkedin.com/company/allbirds/"},      // http + trailing /
				{Name: "Org5", LIURL: "https://www.linkedin.com/company/allbirds"},  // https + www
				{Name: "Org6", LIURL: "http://www.linkedin.com/company/allbirds"},   // http + www
				{Name: "Org7", LIURL: "https://linkedin.com/company/allbirds"},      // https
				{Name: "Org8", LIURL: "http://linkedin.com/company/allbirds"},       // http
			},
			want: []Org{},
		},
		{
			name: "Given an array with orgs having empty LIURL, we return the same array",
			orgs: []Org{
				{Name: "Org1", LIURL: ""},
				{Name: "Org2", LIURL: ""},
			},
			want: []Org{
				{Name: "Org1", LIURL: ""},
				{Name: "Org2", LIURL: ""},
			},
		},
		{
			name: "Given an array with orgs having the same LIURL with different formats, we return the same array",
			orgs: []Org{
				{Name: "Org1", LIURL: "https://www.linkedin.com/company/surfe/"}, // https + www + trailing /
				{Name: "Org2", LIURL: "http://www.linkedin.com/company/surfe/"},  // http + www + trailing /
				{Name: "Org3", LIURL: "https://linkedin.com/company/surfe/"},     // https + trailing /
				{Name: "Org4", LIURL: "http://linkedin.com/company/surfe/"},      // http + trailing /
				{Name: "Org5", LIURL: "https://www.linkedin.com/company/surfe"},  // https + www
				{Name: "Org6", LIURL: "http://www.linkedin.com/company/surfe"},   // http + www
				{Name: "Org7", LIURL: "https://linkedin.com/company/surfe"},      // https
				{Name: "Org8", LIURL: "http://linkedin.com/company/surfe"},       // http
			},
			want: []Org{
				{Name: "Org1", LIURL: "https://www.linkedin.com/company/surfe/"}, // https + www + trailing /
				{Name: "Org2", LIURL: "http://www.linkedin.com/company/surfe/"},  // http + www + trailing /
				{Name: "Org3", LIURL: "https://linkedin.com/company/surfe/"},     // https + trailing /
				{Name: "Org4", LIURL: "http://linkedin.com/company/surfe/"},      // http + trailing /
				{Name: "Org5", LIURL: "https://www.linkedin.com/company/surfe"},  // https + www
				{Name: "Org6", LIURL: "http://www.linkedin.com/company/surfe"},   // http + www
				{Name: "Org7", LIURL: "https://linkedin.com/company/surfe"},      // https
				{Name: "Org8", LIURL: "http://linkedin.com/company/surfe"},       // http
			},
		},
		{
			name: "Given an array with one org having the same LIURL, and anoter org having a different LIURL, we return an array with only the org with the matching URL",
			orgs: []Org{
				{Name: "Org1", LIURL: "https://www.linkedin.com/company/allbirds"},
				{Name: "Org2", LIURL: "https://www.linkedin.com/company/surfe/"},
			},
			want: []Org{
				{Name: "Org2", LIURL: "https://www.linkedin.com/company/surfe/"},
			},
		},
		{
			name: "Given an array with one org having a different LIURL and one org having empty LIURL, we return an array with only the org with empty URL",
			orgs: []Org{
				{Name: "Org1", LIURL: "https://www.linkedin.com/company/allbirds"},
				{Name: "Org2", LIURL: ""},
			},
			want: []Org{
				{Name: "Org2", LIURL: ""},
			},
		},
		{
			name: "Given an array with one org having the same LIURL, one org having a different LIURL and one org having empty LIURL, we return an array with the orgs with empty and matching URLs",
			orgs: []Org{
				{Name: "Org1", LIURL: "https://www.linkedin.com/company/surfe"},
				{Name: "Org2", LIURL: "https://www.linkedin.com/company/allbirds/"},
				{Name: "Org3", LIURL: ""},
			},
			want: []Org{
				{Name: "Org1", LIURL: "https://www.linkedin.com/company/surfe"},
				{Name: "Org3", LIURL: ""},
			},
		},
	}
	for _, tt := range tests {
		t.Run("Given a linkedInURL without trailing slash - "+tt.name, func(t *testing.T) {
			t.Parallel()
			got := RemoveWithEmptyOrNotEqualLinkedinURL[Org](tt.orgs, "key", "https://www.linkedin.com/company/surfe")
			assert.Equal(t, tt.want, got)
		})

		t.Run("Given a linkedInURL with trailing slash - "+tt.name, func(t *testing.T) {
			t.Parallel()
			got := RemoveWithEmptyOrNotEqualLinkedinURL[Org](tt.orgs, "key", "https://www.linkedin.com/company/surfe/")
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestRemoveWithEmptyOrNotEqualLinkedinURL_EmptyLinkedInURL(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		orgs []Org
		want []Org
	}{
		{
			name: "Given an empty array, we return an empty array",
			orgs: []Org{},
			want: []Org{},
		},
		{
			name: "Given an array with any contents, we return the same array",
			orgs: []Org{
				{Name: "Org1", LIURL: "https://www.linkedin.com/company/surfe"},
				{Name: "Org2", LIURL: "https://www.linkedin.com/company/allbirds/"},
				{Name: "Org3", LIURL: ""},
			},
			want: []Org{
				{Name: "Org1", LIURL: "https://www.linkedin.com/company/surfe"},
				{Name: "Org2", LIURL: "https://www.linkedin.com/company/allbirds/"},
				{Name: "Org3", LIURL: ""},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := RemoveWithEmptyOrNotEqualLinkedinURL[Org](tt.orgs, "key", "")
			assert.Equal(t, tt.want, got)
		})
	}
}
