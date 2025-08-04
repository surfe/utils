package utils

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLinkedinURLCleaner(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		rawUrl string
		want   string
	}{
		{
			name:   "url with escaped path and extra params",
			rawUrl: "https://www.linkedin.com/in/cl%C3%A9mence-decaup-682517102/?miniProfileUrn=urn%3Ali%3Afs_miniProfile%3AACoAABoPoAIBxhYFfNI8BkSMy68iINpliothTLE",
			want:   "https://www.linkedin.com/in/cl%C3%A9mence-decaup-682517102",
		},
		{
			name:   "url with unescaped path and extra params",
			rawUrl: "https://www.linkedin.com/in/clémence-decaup-682517102/?miniProfileUrn=urn%3Ali%3Afs_miniProfile%3AACoAABoPoAIBxhYFfNI8BkSMy68iINpliothTLE",
			want:   "https://www.linkedin.com/in/clémence-decaup-682517102",
		},
		{
			name:   "url with space in profile path",
			rawUrl: "https://www.linkedin.com/in/clémence-decaup 682517102/?miniProfileUrn=urn%3Ali%3Afs_miniProfile%3AACoAABoPoAIBxhYFfNI8BkSMy68iINpliothTLE",
			want:   "https://www.linkedin.com/in/clémence-decaup",
		},
		{
			name:   "wrong hostname",
			rawUrl: "https://leadjet.io/in/clémence-decaup 682517102/?miniProfileUrn=urn%3Ali%3Afs_miniProfile%3AACoAABoPoAIBxhYFfNI8BkSMy68iINpliothTLE",
			want:   "",
		},
		{
			name:   "empty url",
			rawUrl: "",
			want:   "",
		},
		{
			name:   "url without scheme",
			rawUrl: "www.linkedin.com/in/clémence-decaup 682517102/?miniProfileUrn=urn%3Ali%3Afs_miniProfile%3AACoAABoPoAIBxhYFfNI8BkSMy68iINpliothTLE",
			want:   "",
		},
		{
			name:   "url with more path blocks",
			rawUrl: "https://www.linkedin.com/in/clémence-decaup/some-dummy-path",
			want:   "https://www.linkedin.com/in/clémence-decaup",
		},
		{
			name:   "company url",
			rawUrl: "https://www.linkedin.com/company/leadjet?testing",
			want:   "https://www.linkedin.com/company/leadjet",
		},
		{
			name:   "company url with special char '&'",
			rawUrl: "https://www.linkedin.com/company/leadjet&friends?testing=test/test",
			want:   "https://www.linkedin.com/company/leadjet&friends",
		},
		{
			name:   "company url with special char '/'",
			rawUrl: "https://www.linkedin.com/company/leadjet&friends/test",
			want:   "https://www.linkedin.com/company/leadjet&friends",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := LinkedinURLCleaner(tt.rawUrl)
			if got != tt.want {
				t.Errorf("CleanLinkedinURL() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLinkedinURLCleanerErr(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		input        string
		escapeHandle bool
		want         string
		wantErr      error
	}{
		{
			name:         "Valid URL with trailing slash should return input without the trailing slash",
			input:        "https://www.linkedin.com/in/jude-don-4aaa01234/",
			escapeHandle: true,
			want:         "https://www.linkedin.com/in/jude-don-4aaa01234",
			wantErr:      nil,
		},
		{
			name:         "Invalid URL should return input as is with error",
			input:        "https://google.com",
			escapeHandle: true,
			want:         "https://google.com",
			wantErr:      errors.New("not a LinkedIn URL"),
		},
		{
			name:         "Valid URL with escaped handle should return escaped",
			input:        "https://www.linkedin.com/company/sek-s%C3%BCt-end%C3%BCstri%CC%87si%CC%87-kurumu-anoni%CC%87m-%C5%9Fi%CC%87rketi%CC%87",
			escapeHandle: true,
			want:         "https://www.linkedin.com/company/sek-s%C3%BCt-end%C3%BCstri%CC%87si%CC%87-kurumu-anoni%CC%87m-%C5%9Fi%CC%87rketi%CC%87",
			wantErr:      nil,
		},
		{
			name:         "Valid URL with unescaped handle with `escapeHandle` option false, should return unescaped",
			input:        "https://www.linkedin.com/company/sek-süt-endüstri̇si̇-kurumu-anoni̇m-şi̇rketi̇/",
			escapeHandle: false,
			want:         "https://www.linkedin.com/company/sek-süt-endüstri̇si̇-kurumu-anoni̇m-şi̇rketi̇",
			wantErr:      nil,
		},
		{
			name:         "Valid URL with unescaped handle with `escapeHandle` option true, should return escaped",
			input:        "https://www.linkedin.com/company/sek-süt-endüstri̇si̇-kurumu-anoni̇m-şi̇rketi̇/",
			escapeHandle: true,
			want:         "https://www.linkedin.com/company/sek-s%C3%BCt-end%C3%BCstri%CC%87si%CC%87-kurumu-anoni%CC%87m-%C5%9Fi%CC%87rketi%CC%87",
			wantErr:      nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result, err := LinkedinURLCleanerErr(tt.input, tt.escapeHandle)
			require.Equal(t, tt.want, result)
			require.Equal(t, tt.wantErr, err)
		})
	}
}

func TestTruncateString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "should not truncate if not longer",
			input: "Fermentum Tortor Condi",
			want:  "Fermentum Tortor Condi",
		},
		{
			name:  "truncates simple longer text",
			input: "Fermentum Tortor Condimentum",
			want:  "Fermentum Tortor Cond" + DefaultOmission,
		},
		{
			name:  "truncates text with unicode chars",
			input: "Göz göre göre şarkı söyleyemek vardı",
			want:  "Göz göre göre şarkı s" + DefaultOmission,
		},
		{
			name:  "truncates text with emojis",
			input: "Venenatis V😛estib👻ulum Lorem Ligula",
			want:  "Venenatis Vestibulum " + DefaultOmission,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := TruncateString(tt.input, 22)
			if got != tt.want {
				t.Errorf("TruncateString() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSanitizeForSOQL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		str  string
		want string
	}{
		{
			name: "Single-quote",
			str:  "L'occitane grou'p",
			want: `L\'occitane grou\'p`,
		},
		{
			name: "Double-quote",
			str:  `L"occitane gro"up`,
			want: `L\"occitane gro\"up`,
		},
		{
			name: "Underscore",
			str:  `L_occitane grou_p`,
			want: `L\_occitane grou\_p`,
		},
		{
			name: "Percent sign",
			str:  `L%occitane grou%p`,
			want: `L\%occitane grou\%p`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			res := SanitizeForSOQL(tt.str)
			require.Equal(t, tt.want, res)
		})
	}
}

func TestDomainFromURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "Should return both domain and tld",
			input: "https://leadjet.io/blog/some-article",
			want:  "leadjet.io",
		},
		{
			name:  "Should return both domain and tld",
			input: "https://www.leadjet.design?param=testing",
			want:  "leadjet.design",
		},
		{
			name:  "Should return both domain and tld",
			input: "http://leadjet.design?param=testing",
			want:  "leadjet.design",
		},
		{
			name:  "Should detect URLs without protocol",
			input: "gelsenwasser.de",
			want:  "gelsenwasser.de",
		},
		{
			name:  "Should return empty when url is not valid",
			input: "xxxx",
			want:  "",
		},
		{
			name:  "Should return empty when url is not valid",
			input: "https:///dddd.com",
			want:  "",
		},
		{
			name:  "Should return redirected url domain when url is an urlshortened url for bit.ly",
			input: "https://bit.ly/3qs9ftN",
			want:  "hawkeyeinnovations.com",
		},
		{
			name:  "Should return redirected url domain when url is an urlshortened url for bit.ly",
			input: "https://bit.ly/333332323233",
			want:  "surfe.com",
		},
		{
			name:  "Should return redirected url domain when url is an urlshortened url for tinyurl.com",
			input: "https://tinyurl.com/rrrrr2323232",
			want:  "surfe.com",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, tt.want, DomainFromURL(tt.input))
		})
	}
}

func TestDomainFromURLNoFiltering(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		input  string
		want   string
		hasErr bool
	}{
		{
			name:  "Should return both domain and tld",
			input: "https://leadjet.io/blog/some-article",
			want:  "leadjet.io",
		},
		{
			name:  "Should return both domain and tld",
			input: "https://www.leadjet.design?param=testing",
			want:  "leadjet.design",
		},
		{
			name:  "Should return both domain and tld",
			input: "http://leadjet.design?param=testing",
			want:  "leadjet.design",
		},
		{
			name:  "Should detect URLs without protocol",
			input: "gelsenwasser.de",
			want:  "gelsenwasser.de",
		},
		{
			name:   "Should return empty when url is not valid",
			input:  "xxxx",
			want:   "",
			hasErr: true,
		},
		{
			name:   "Should return empty when url is not valid",
			input:  "https:///dddd.com",
			want:   "",
			hasErr: true,
		},
		{
			name:  "Should handle multi-TLD domains correctly",
			input: "www.some.gov.eg/resource",
			want:  "some.gov.eg",
		},
		{
			name:  "Should handle URLs with a trailing slash correctly",
			input: "www.some.net/resource/",
			want:  "some.net",
		},
		{
			name:  "Should ignore fragments",
			input: "https://www.surfe.com/resource?query=value#fragment",
			want:  "surfe.com",
		},
		{
			name:  "Should ignore subdomains",
			input: "https://app.some.co.uk",
			want:  "some.co.uk",
		},
		{
			name:  "Should return domain without filtering",
			input: "https://bit.ly/3qs9ftN",
			want:  "bit.ly",
		},
		{
			name:  "Should handle a multi-part public suffix as a valid domain",
			input: "uk.com",
			want:  "uk.com",
		},
		{
			name:  "Should handle a different multi-part public suffix with a protocol",
			input: "https://co.uk",
			want:  "co.uk",
		},
		{
			name:  "Should correctly parse a full domain that uses a multi-part suffix",
			input: "my-business.co.uk",
			want:  "my-business.co.uk",
		},
		{
			name:  "Should correctly parse a full domain that uses uk.com",
			input: "https://another-example.uk.com/path",
			want:  "another-example.uk.com",
		},
		{
			name:   "Should fail for a single-label top-level domain",
			input:  "com",
			want:   "",
			hasErr: true,
		},
		{
			name:   "Should fail for 'localhost' as it has no TLD",
			input:  "localhost",
			want:   "",
			hasErr: true,
		},
		{
			name:   "Should fail for 'localhost' with a port and protocol",
			input:  "http://localhost:3000",
			want:   "",
			hasErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			domain, err := DomainFromURLNoFiltering(tt.input)
			if !tt.hasErr {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}

			require.Equal(t, tt.want, domain)
		})
	}
}

func TestFormatDomainURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "Should return domain without https and /",
			input: "https://leadjet.io/",
			want:  "leadjet.io",
		},
		{
			name:  "Should return domain without https",
			input: "https://leadjet.io",
			want:  "leadjet.io",
		},
		{
			name:  "Should return domain without http",
			input: "http://surfe.com",
			want:  "surfe.com",
		},
		{
			name:  "Should return domain without http but www",
			input: "https://www.surfe.com",
			want:  "www.surfe.com",
		},
		{
			name:  "Should detect URLs without protocol",
			input: "gelsenwasser.de",
			want:  "gelsenwasser.de",
		},
		{
			name:  "Should return itself when url is not valid",
			input: "xxxx",
			want:  "xxxx",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, tt.want, FormatDomainURL(tt.input))
		})
	}
}

func TestRemoveQueryParams(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "Base URL should return itself",
			input: "https://www.surfe.com",
			want:  "https://www.surfe.com",
		},
		{
			name:  "Base URL should return itself without trailing /",
			input: "https://www.surfe.com/",
			want:  "https://www.surfe.com",
		},
		{
			name:  "Should return domain without query params",
			input: "https://www.surfe.com?utm_source=linkedin&utm_medium=companypage&utm_campaign=maincta",
			want:  "https://www.surfe.com",
		},
		{
			name:  "Should return domain without query params and /",
			input: "https://www.surfe.com/?utm_source=linkedin&utm_medium=companypage&utm_campaign=maincta",
			want:  "https://www.surfe.com",
		},
		{
			name:  "Should work with http",
			input: "http://www.surfe.com/?utm_source=linkedin&utm_medium=companypage&utm_campaign=maincta",
			want:  "http://www.surfe.com",
		},
		{
			name:  "Should return itself when url is not valid",
			input: "xxxx",
			want:  "xxxx",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, tt.want, RemoveQueryParams(tt.input))
		})
	}
}

func TestConvertMDToHTML_ConvertHTMLToMD_Compatibility(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		text string
		want string
	}{
		{
			name: "Preserve 1 & 2 line breaks",
			text: "L1\nL2\n\nL3\n\n\nL4",
			want: "L1\nL2\n\nL3\n\nL4",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			html := ConvertMDToHTML(tt.text)
			md := ConvertHTMLToMD(html, "")
			require.Equal(t, tt.want, md)
		})
	}
}

func TestURNExtractor(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "Input follows entityURN format urn:li:fs_normalized_company:13205888",
			input: "urn:li:fs_normalized_company:13205888",
			want:  "13205888",
		},
		{
			name:  "Input follows objectURN format urn:li:member:22719531",
			input: "urn:li:member:22719531",
			want:  "22719531",
		},
		{
			name:  "Input is empty string",
			input: "",
			want:  "",
		},
		{
			name:  "Input is a non-empty string without any colon (:) character",
			input: "non-empty string without colon",
			want:  "non-empty string without colon",
		},
		{
			name:  "Should extract the SalesNav ID",
			input: "urn:li:fs_salesProfile:(ACwAAAJlc6wBYdHGFmVJDHu,NAME_SEARCH,ij9X)",
			want:  "ACwAAAJlc6wBYdHGFmVJDHu",
		},
		{
			name:  "Should extract the SalesNav ID even there is space",
			input: "urn:li:fs_salesProfile:( ACwAAAJlc6wBYdHGFmVJDHu , NAME_SEARCH , ij9X )",
			want:  "ACwAAAJlc6wBYdHGFmVJDHu",
		},
		{
			name:  "Should extract the SalesNav ID with no comma",
			input: "urn:li:fs_salesProfile:(ACwAAAJlc6wBYdHGFmVJDHu)",
			want:  "ACwAAAJlc6wBYdHGFmVJDHu",
		},
		{
			name:  "Should not return value when no content in parenthesis",
			input: "urn:li:fs_salesProfile:()",
			want:  "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := URNExtractor(tt.input)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestToValidEmail(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		email string
		want  bool
	}{
		{
			name:  "Valid email should return input email back and true",
			email: "demo-pd@leadjet.com",
			want:  true,
		},
		{
			name:  "Invalid email should return input email back and true",
			email: "demo-pd@leadjet",
			want:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result, got := ToValidEmail(tt.email)
			require.Equal(t, tt.want, got)
			require.Equal(t, tt.email, result)
		})
	}
}

func TestDomainFromEmail(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		email  string
		domain string
	}{
		{
			name:   "Valid email should return domain",
			email:  "demo-pd@leadjet.com",
			domain: "leadjet.com",
		},
		{
			name:   "Invalid email should return input email back and true",
			email:  "demo-pdd",
			domain: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := DomainFromEmail(tt.email)
			require.Equal(t, tt.domain, result)
		})
	}
}

func TestSameDomainsl(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		linkedinDomain string
		crmDomain      string
		result         bool
	}{
		{
			name:           "Same domains should return true",
			linkedinDomain: "leadjet.com",
			crmDomain:      "leadjet.com",
			result:         true,
		},
		{
			name:           "Different domains should return false",
			linkedinDomain: "leadjet.io",
			crmDomain:      "surfe.com",
			result:         false,
		},
		{
			name:           "Empty crm domain should return false",
			linkedinDomain: "surfe.com",
			crmDomain:      "",
			result:         false,
		},
		{
			name:           "Empty linkedin domain should return false",
			linkedinDomain: "",
			crmDomain:      "surfe.com",
			result:         false,
		},
		{
			name:           "Empty domains should return false",
			linkedinDomain: "",
			crmDomain:      "",
			result:         false,
		},
		{
			name:           "FQDN Crm domain and bare linkedin domain should return true",
			linkedinDomain: "surfe.com",
			crmDomain:      "http://www.surfe.com",
			result:         true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := SameDomains(tt.linkedinDomain, tt.crmDomain)
			require.Equal(t, tt.result, result)
		})
	}
}

func TestRemoveAccents(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		want    string
		wantErr error
	}{
		{
			name:    "Valid input with no accent should return expectred",
			input:   "abc",
			want:    "abc",
			wantErr: nil,
		},
		{
			name:    "Valid input with accent should return expected",
			input:   "abc üğişçöıÜĞİŞÇÖI def",
			want:    "abc ugiscoıUGISCOI def",
			wantErr: nil,
		},
		{
			name:    "Valid input with accent should return expected",
			input:   "ḇoḇ čÁr DĀāvid Éṽḙ Ḟolḏ Ğod Ḧelp J̌ÚlĬe Ķit-kat ḼÒve ṂoṀ Ńuń ÓṔeŔa quiť ŝlow ẃẀwÍ",
			want:    "bob cAr DAavid Eve Fold God Help JUlIe Kit-kat LOve MoM Nun OPeRa quit slow wWwI",
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result, err := RemoveAccents(tt.input)
			require.Equal(t, tt.want, result)
			require.Equal(t, tt.wantErr, err)
		})
	}
}

func TestMatchLinkedinURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		url1 string
		url2 string
		want bool
	}{
		{
			name: "Same urls should return true",
			url1: "https://www.linkedin.com/company/sek",
			url2: "https://www.linkedin.com/company/sek/",
			want: true,
		},
		{
			name: "Same encoded urls should return true",
			url1: "https://www.linkedin.com/company/sek-s%C3%BCt-end%C3%BCstri%CC%87si%CC%87-kurumu-anoni%CC%87m-%C5%9Fi%CC%87rketi%CC%87",
			url2: "https://www.linkedin.com/company/sek-s%C3%BCt-end%C3%BCstri%CC%87si%CC%87-kurumu-anoni%CC%87m-%C5%9Fi%CC%87rketi%CC%87/",
			want: true,
		},
		{
			name: "Different urls should return false",
			url1: "https://www.linkedin.com/company/fff",
			url2: "https://www.linkedin.com/company/vvvv",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			res := MatchLinkedinURL(tt.url1, tt.url2)
			require.Equal(t, tt.want, res)
		})
	}
}

func TestExtractPlanInfo(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		planID       string
		planName     string
		planInterval string
	}{
		{
			name:         "Given starter_eur_monthly plan id, should return starter and monthly",
			planID:       "starter_eur_monthly",
			planName:     "starter",
			planInterval: "monthly",
		},
		{
			name:         "Given professional_usd_yearly plan id, should return professional and yearly",
			planID:       "professional_usd_yearly",
			planName:     "professional",
			planInterval: "yearly",
		},
		{
			name:         "Given surfe_yearly_basic_usd plan id, should return basic and yearly",
			planID:       "surfe_yearly_basic_usd",
			planName:     "basic",
			planInterval: "yearly",
		},
		{
			name:         "Given surfe_monthly_business_eur plan id, should return business_ and monthly",
			planID:       "surfe_monthly_business_eur",
			planName:     "business",
			planInterval: "monthly",
		},
		{
			name:         "Given entreprise_eur_monthly plan id, should return enterprise and monthly",
			planID:       "entreprise_eur_monthly",
			planName:     "enterprise",
			planInterval: "monthly",
		},
		{
			name:         "Given enrich_monthly_eur plan id, should return enrich and monthly",
			planID:       "enrich_monthly_eur",
			planName:     "enrich",
			planInterval: "monthly",
		},
		{
			name:         "Given notavalidplanid plan id, should return empty and empty",
			planID:       "notavalidplanid",
			planName:     "",
			planInterval: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			resPlanName, resPlanInterval := ExtractPlanInfo(tt.planID)
			require.Equal(t, tt.planName, resPlanName)
			require.Equal(t, tt.planInterval, resPlanInterval)
		})
	}
}

func TestURLHostnameExtractor(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		str  string
		want string
	}{
		{
			name: "https with www",
			str:  "https://www.surfe.com",
			want: "surfe.com",
		},
		{
			name: "http with www",
			str:  "http://www.surfe.com",
			want: "surfe.com",
		},
		{
			name: "only www",
			str:  "www.surfe.com",
			want: "surfe.com",
		},
		{
			name: "only https",
			str:  "https://surfe.com",
			want: "surfe.com",
		},
		{
			name: "only http",
			str:  "http://surfe.com",
			want: "surfe.com",
		},
		{
			name: "with http and path",
			str:  "http://surfe.com/test",
			want: "surfe.com",
		},
		{
			name: "http ending with /",
			str:  "http://surfe.com/",
			want: "surfe.com",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			res := URLHostnameExtractor(tt.str)
			require.Equal(t, tt.want, res)
		})
	}
}

func Test_SimplifyName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		given string
		want  string
	}{
		{
			name:  "Firstname without abbreviation",
			given: "Marie",
			want:  "Marie",
		},
		{
			name:  "Lastname without abbreviation",
			given: "DRAPEAU",
			want:  "DRAPEAU",
		},
		{
			name:  "Firstname with abbreviation",
			given: "Dr John",
			want:  "John",
		},
		{
			name:  "Firstname with uppercase abbreviation",
			given: "PR. DR. John",
			want:  "John",
		},
		{
			name:  "Firstname with abbreviation",
			given: "Dr.John",
			want:  "John",
		},
		{
			name:  "Firstname with abbreviation",
			given: "PMP, John",
			want:  "John",
		},
		{
			name:  "Firstname with comma separated abbreviation",
			given: "PMP,John",
			want:  "John",
		},
		{
			name:  "Lastname with comma+space separated abbreviation",
			given: "John, CXAP",
			want:  "John",
		},
		{
			name:  "Lastname with comma separated abbreviation",
			given: "John,mba",
			want:  "John",
		},
		{
			name:  "Lastname with space separated abbreviation",
			given: "John MBA",
			want:  "John",
		},
		{
			name:  "Lastname with multiple comma|dot|space separated abbreviations",
			given: "Doe II. III, IV FRSA, F.R.S.A",
			want:  "Doe",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			n := SimplifyName(tt.given)
			require.Equal(t, tt.want, n)
		})
	}
}

func Test_SimplifyCompanyName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		given string
		want  string
	}{
		{
			name:  "no edit",
			given: "Surfe (ex-leadjet)",
			want:  "Surfe (ex-leadjet)",
		},
		{
			name:  "no edit",
			given: "Inclu",
			want:  "Inclu",
		},
		{
			name:  ", INC.",
			given: "Surfe (ex-leadjet), INC.",
			want:  "Surfe (ex-leadjet)",
		},
		{
			name:  "INC.",
			given: "Company INC.",
			want:  "Company",
		},
		{
			name:  "inc,",
			given: "inc, Company",
			want:  "Company",
		},
		{
			name:  " INC ",
			given: "Limited INC Company",
			want:  "Company",
		},
		{
			name:  ", INC,",
			given: "Company, INC, Incorporated",
			want:  "Company, Incorporated",
		},
		{
			name:  ". INC.",
			given: "Company. INC.",
			want:  "Company",
		},
		{
			name:  "S.A.S.",
			given: "S.A.S. Company. INC.",
			want:  "Company",
		},
		{
			name:  "S.A.",
			given: "S.A. Company. INC.",
			want:  "Company",
		},
		{
			name:  "A/S",
			given: "A/S Company. INC.",
			want:  "Company",
		},
		{
			name:  "Pty Ltd",
			given: "Pty Ltd Company. INC.",
			want:  "Company",
		},
		{
			name:  "SE & Co. KGaA",
			given: "SE & Co. KGaA Company. INC.",
			want:  "Company",
		},
		{
			name:  "GmbH & Co. KGaA",
			given: "GmbH & Co. KGaA Company. INC.",
			want:  "Company",
		},
		{
			name:  "With multiple abbreviations",
			given: "STIHL Kettenwerk GmbH & Co KG, Waiblingen (DE), Zweigniederlassung Wil SG",
			want:  "STIHL Kettenwerk & Waiblingen (DE), Wil SG",
		},
		{
			name:  "Limited",
			given: "Cardinal Financial Company, Limited Partnership",
			want:  "Cardinal Financial Company, Partnership",
		},
		{
			// I was not able to omit this, we need to live with this exception
			name:  "Dot exception",
			given: "INC.orporated",
			want:  "orporated",
		},
		{
			name:  "Should not remove SE from name if there is no space/dot before",
			given: "Pathrise",
			want:  "Pathrise",
		},
		{
			name:  "Should not remove OY from name if there is no space/dot before",
			given: "Le Roy Logistique",
			want:  "Le Roy Logistique",
		},
		{
			name:  "Should remove SE from name if there is a space/dot before",
			given: "SE Pathrise INC.",
			want:  "Pathrise",
		},
		{
			name:  "If the simplified company name is shorter than 3 characters, revert to the original",
			given: "SAS",
			want:  "SAS",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			n := SimplifyCompanyName(tt.given)
			require.Equal(t, tt.want, n)
		})
	}
}

func TestURLProfileExtract(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		handle     string
		wantHandle string
	}{
		{
			name:       "Common case",
			handle:     "abc-xyz",
			wantHandle: "abc-xyz",
		},
		{
			name:       "Common case with trailing slash",
			handle:     "abc-xyz/",
			wantHandle: "abc-xyz",
		},
		{
			name:       "Trailing text added by customer",
			handle:     "abc-xyz (Surfe)",
			wantHandle: "abc-xyz",
		},
		{
			name:       "Trailing text added by customer with slash",
			handle:     "abc-xyz/ (Surfe)",
			wantHandle: "abc-xyz",
		},
	}
	for _, tt := range tests {
		t.Run("Person LI URL - "+tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.wantHandle, URLProfileExtract("https://www.linkedin.com/in/"+tt.handle))
		})
		t.Run("Company LI URL - "+tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.wantHandle, URLProfileExtract("https://www.linkedin.com/company/"+tt.handle))
		})
		t.Run("School LI URL - "+tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.wantHandle, URLProfileExtract("https://www.linkedin.com/school/"+tt.handle))
		})
	}
}

func TestExtractHostAndPath(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		fullURL string
		wantURL string
	}{
		{
			name:    "Provided a full HTTPS URL, it returns only the host+path, without the w3 subdomain and without trailing /",
			fullURL: "https://www.linkedin.com/company/surfe/",
			wantURL: "linkedin.com/company/surfe",
		},
		{
			name:    "Provided a HTTPS URL without the w3 subdomain, it returns only the host+path without trailing /",
			fullURL: "https://linkedin.com/company/surfe/",
			wantURL: "linkedin.com/company/surfe",
		},
		{
			name:    "Provided a full HTTP URL, it returns only the host+path, without the w3 subdomain and without trailing /",
			fullURL: "https://www.linkedin.com/company/surfe/",
			wantURL: "linkedin.com/company/surfe",
		},
		{
			name:    "Provided a HTTP URL without the w3 subdomain, it returns only the host+path without trailing /",
			fullURL: "https://linkedin.com/company/surfe/",
			wantURL: "linkedin.com/company/surfe",
		},
		{
			name:    "Provided a URL without the scheme but including the w3 subdomain, it returns the same URL without trailing /",
			fullURL: "www.linkedin.com/company/surfe/",
			wantURL: "www.linkedin.com/company/surfe",
		},
		{
			name:    "Provided a URL without the scheme, it returns the same URL without trailing /",
			fullURL: "linkedin.com/company/surfe/",
			wantURL: "linkedin.com/company/surfe",
		},
		{
			name:    "Provided an invalid URL, it returns the same string without trailing /",
			fullURL: ":/linkedin.com/company/surfe/",
			wantURL: ":/linkedin.com/company/surfe",
		},
		{
			name:    "Provided a random string, it returns the same string without trailing /",
			fullURL: "this is not an URL",
			wantURL: "this is not an URL",
		},
		{
			name:    "Provided an empty string, it returns an empty string without trailing /",
			fullURL: "",
			wantURL: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.wantURL, ExtractHostAndPath(tt.fullURL))
		})
	}
}

func TestGenerateURLCombinations(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		url  string
		want []string
	}{
		{
			name: "Provided a full LI URL, it returns all the possible combinations",
			url:  "https://www.linkedin.com/company/surfe/",
			want: []string{
				"https://www.linkedin.com/company/surfe/", // https + www + trailing /
				"http://www.linkedin.com/company/surfe/",  // http + www + trailing /
				"https://linkedin.com/company/surfe/",     // https + trailing /
				"http://linkedin.com/company/surfe/",      // http + trailing /
				"https://www.linkedin.com/company/surfe",  // https + www
				"http://www.linkedin.com/company/surfe",   // http + www
				"https://linkedin.com/company/surfe",      // https
				"http://linkedin.com/company/surfe",       // http
			},
		},
		{
			name: "Provided a any string, it returns all the possible combinations",
			url:  "any string",
			want: []string{
				"https://www.any string/", // https + www + trailing /
				"http://www.any string/",  // http + www + trailing /
				"https://any string/",     // https + trailing /
				"http://any string/",      // http + trailing /
				"https://www.any string",  // https + www
				"http://www.any string",   // http + www
				"https://any string",      // https
				"http://any string",       // http
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.ElementsMatch(t, tt.want, GenerateURLCombinations(tt.url))
		})
	}
}

func TestExtractSalesNavIDFromURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		url  string
		want string
	}{
		{
			name: "Test with SalesNav URL",
			url:  "https://www.linkedin.com/sales/people/ACwAAAB0jjIBhkHdtLPXAimZXL_SAtpk7UFZbik,OUT_OF_NETWORK,2N7m",
			want: "ACwAAAB0jjIBhkHdtLPXAimZXL_SAtpk7UFZbik",
		},
		{
			name: "Test simple SalesNav URL without extra params",
			url:  "https://www.linkedin.com/sales/people/ACwAAAB0jjIBhkHdtLPXAimZXL_SAtpk7UFZbik",
			want: "ACwAAAB0jjIBhkHdtLPXAimZXL_SAtpk7UFZbik",
		},
		{
			name: "Invalid URL returns empty string",
			url:  "https://www.linkedin.com/sales/people/",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := ExtractSalesNavIDFromURL(tt.url); got != tt.want {
				t.Errorf("ExtractSalesNavIDFromURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMatchLIURLByIDOrHandle(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		name       string
		s          string
		idOrHandle string
		expected   bool
	}{
		{
			"https://www.",
			"https://www.linkedin.com/in/johndoe",
			"johndoe",
			true,
		},
		{
			"https://www. with query params",
			"https://www.linkedin.com/in/johndoe?locale=en_US",
			"johndoe",
			true,
		},
		{
			"https:// without www",
			"https://linkedin.com/in/johndoe",
			"johndoe",
			true,
		},
		{
			"https://www. with dash in handle",
			"https://www.linkedin.com/company/xyz-corp",
			"xyz-corp",
			true,
		},
		{
			"Unmatched handle",
			"https://www.linkedin.com/in/johndoe",
			"janedoe",
			false,
		},
		{
			"Not a LI URL",
			"https://www.example.com/johndoe",
			"johndoe",
			false,
		},
		{
			"idOrHandle not provided",
			"https://www.linkedin.com/in/johndoe",
			"",
			false,
		},
		{
			"Slash at the end",
			"https://www.linkedin.com/in/johndoe/",
			"johndoe",
			true,
		},
		{
			"Similar word",
			"https://www.linkedin.com/in/johndoe2",
			"johndoe",
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if output := MatchLIURLByIDOrHandle(tt.s, tt.idOrHandle); output != tt.expected {
				t.Errorf("Test failed for %s %s, expected output %t, but got %t", tt.s, tt.idOrHandle, tt.expected, output)
			}
		})
	}
}

func TestLooseString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		str  string
		want string
	}{
		{
			name: "Given long name with special characters we return simplified string",
			str:  " Sirąj Uddiń--Chowdhury\n & Rubel (CV Doctor)?   _-",
			want: "sirajuddinchowdhuryandrubelcvdoctor",
		},
		{
			name: "Given html tags, they are not cleaned up",
			str:  "<br>",
			want: "br",
		},
		{
			name: "Given accent characters, they are getting unicodified",
			str:  "ß",
			want: "ss",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := LooseString(tt.str)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestSanitizedLooseString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		str  string
		want string
	}{
		{
			name: "Given long name with special characters we return simplified string",
			str:  " Sirąj Uddiń--Chowdhury\n & Rußel (CV Doctor)?   _-",
			want: "sirajuddinchowdhuryandrusselcvdoctor",
		},
		{
			name: "Given html tags, they are cleaned up",
			str:  " Sirąj Uddiń <br> Cho",
			want: "sirajuddincho",
		},
		{
			name: "Given accent characters, they are getting unicodified",
			str:  "ß",
			want: "ss",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := SanitizedLooseString(tt.str)
			require.Equal(t, tt.want, got)
		})
	}
}

func BenchmarkLooseString(b *testing.B) {
	for range b.N {
		LooseString(" Sirąj Uddiń--Chowdhury\n & Rußel (CV Doctor)?   _-")
	}
}

func BenchmarkSanitizedLooseString(b *testing.B) {
	for range b.N {
		SanitizedLooseString(" Sirąj<br> Uddiń--Chowdhury\n & Rußel (CV Doctor)?   _-")
	}
}

func TestMaskStringWithAsterisks(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		str  string
		want string
	}{
		{
			name: "Given String length less than or equal to 2 then return **",
			str:  "a",
			want: "**",
		},
		{
			name: "Given String length greater than 2 then return first and last character with *",
			str:  "abc",
			want: "a*c",
		},
		{
			name: "Given String length greater than 2 with special characters then return first and last character with *",
			str:  "a#c",
			want: "a*c",
		},
		{
			name: "Given Empty string then return **",
			str:  "",
			want: "**",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			actual := MaskStringWithAsterisks(tt.str)
			require.Equal(t, tt.want, actual)
		})
	}
}

func TestFirstAndLastNameFromFullName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		input             string
		expectedFirstName string
		expectedLastName  string
	}{
		{
			name:              "returns empty if input is empty",
			input:             "",
			expectedFirstName: "",
			expectedLastName:  "",
		},
		{
			name:              "returns empty if input is only spaces",
			input:             "   ",
			expectedFirstName: "",
			expectedLastName:  "",
		},
		{
			name:              "returns first name only if name has no last name",
			input:             "Nour",
			expectedFirstName: "Nour",
			expectedLastName:  "",
		},
		{
			name:              "returns first name only if name has no last name and has spaces before/after",
			input:             "     Evren    ",
			expectedFirstName: "Evren",
			expectedLastName:  "",
		},
		{
			name:              "returns first and last name if name has at least 2 words with a space between them",
			input:             "Aleksander Piotrowski",
			expectedFirstName: "Aleksander",
			expectedLastName:  "Piotrowski",
		},
		{
			name:              "returns first and last name if name has at least 2 words with a space between them, ignores additional spaces",
			input:             " Soner   Eker ",
			expectedFirstName: "Soner",
			expectedLastName:  "Eker",
		},
		{
			name:              "returns first word as the first name and remainding words as the last name if name has at more than 2 words with a space between them, ignores additional spaces",
			input:             " Laura   Osset   Blanch ",
			expectedFirstName: "Laura",
			expectedLastName:  "Osset Blanch",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			actualFirstName, actualLastName := FirstAndLastNameFromFullName(tt.input)
			require.Equal(t, tt.expectedFirstName, actualFirstName)
			require.Equal(t, tt.expectedLastName, actualLastName)
		})
	}
}

func TestGetRedirectedDomain(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		url            string
		expectedDomain string
		expectedError  string
	}{
		{
			name:           "Given Valid redirection then return domain",
			url:            "https://bit.ly/333332323233",
			expectedDomain: "surfe.com",
			expectedError:  "",
		},
		{
			name:           "Given Invalid URL then return empty domain and error",
			url:            "http://invalid-url",
			expectedDomain: "",
			expectedError:  "failed to get redirection URL for http://invalid-url",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			domain, err := getRedirectedDomain(tt.url)
			if tt.expectedError == "" {
				require.NoError(t, err)
				require.Equal(t, tt.expectedDomain, domain)
			} else {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.expectedError)
			}
		})
	}
}

func TestExtractLinkedInSlug(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		url  string
		slug string
	}{
		{
			name: "Empty string",
			url:  "",
			slug: "",
		},
		{
			name: "Already a handle",
			url:  "johndoe",
			slug: "johndoe",
		},
		{
			name: "Valid LinkedIn profile URL",
			url:  "https://www.linkedin.com/in/johndoe",
			slug: "johndoe",
		},
		{
			name: "Valid LinkedIn profile URL with trailing slash",
			url:  "https://www.linkedin.com/in/johndoe/",
			slug: "johndoe",
		},
		{
			name: "Valid LinkedIn company URL",
			url:  "https://www.linkedin.com/company/xyz-corp",
			slug: "xyz-corp",
		},
		{
			name: "Valid LinkedIn school URL",
			url:  "https://www.linkedin.com/school/abc-university",
			slug: "abc-university",
		},
		{
			name: "Invalid LinkedIn URL",
			url:  "https://www.linkedin.com/invalid/xyz",
			slug: "",
		},
		{
			name: "URL with special characters",
			url:  "https://www.linkedin.com/in/john%20doe",
			slug: "john doe",
		},
		{
			name: "URL with query parameters",
			url:  "https://www.linkedin.com/in/johndoe?param=value",
			slug: "johndoe",
		},
		{
			name: "URL with fragment",
			url:  "https://www.linkedin.com/in/johndoe#section",
			slug: "johndoe",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := ExtractLinkedInSlug(tt.url)
			require.Equal(t, tt.slug, got)
		})
	}
}

func TestGetLocationWithoutCountry(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		location string
		want     string
	}{
		{
			name:     "returns an empty string if location is empty",
			location: "",
			want:     "",
		},
		{
			name:     "returns city, region if location is city, region",
			location: "Poznan, Poland",
			want:     "Poznan",
		},
		{
			name:     "returns city, region if location is city, region, country",
			location: "Paris, Île-de-France, France",
			want:     "Paris, Île-de-France",
		},
		{
			name:     "returns city if location is city, country",
			location: "Cairo, Egypt",
			want:     "Cairo",
		},
		{
			name:     "handles special characters correctly",
			location: "Hyvinkää, Uusimaa, Finland",
			want:     "Hyvinkää, Uusimaa",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := GetLocationWithoutCountry(context.Background(), tt.location)

			require.Equal(t, tt.want, got)
		})
	}
}

func TestTitle(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "Empty string",
			input: "",
			want:  "",
		},
		{
			name:  "Empty space",
			input: " ",
			want:  "",
		},
		{
			name:  "Single word",
			input: "hello",
			want:  "Hello",
		},
		{
			name:  "Multiple words",
			input: "hello world",
			want:  "Hello World",
		},
		{
			name:  "Already capitalized",
			input: "Hello World",
			want:  "Hello World",
		},
		{
			name:  "Mixed case",
			input: "hElLo WoRLd",
			want:  "HElLo WoRLd",
		},
		{
			name:  "Single letter words",
			input: "a b c",
			want:  "A B C",
		},
		{
			name:  "Words with punctuation",
			input: "hello, world!",
			want:  "Hello, World!",
		},
		{
			name:  "Words with numbers",
			input: "hello 123 world",
			want:  "Hello 123 World",
		},
		{
			name:  "Words with special characters",
			input: "hello @world",
			want:  "Hello @world",
		},
		{
			name:  "Words with unicode characters",
			input: "héllo wörld",
			want:  "Héllo Wörld",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := Title(tt.input)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestEntityURN(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		want    URN
		wantErr string
	}{
		{
			name:  "valid sales profile URN",
			input: "urn:li:fs_salesProfile:(ACwAAAKWZe8BZ8gXVKS6ePAs8I4GWmjW4Tjm-7w,NAME_SEARCH,xqt8)",
			want: URN{
				ProfileID: "ACwAAAKWZe8BZ8gXVKS6ePAs8I4GWmjW4Tjm-7w",
				AuthType:  "NAME_SEARCH",
				AuthToken: "xqt8",
			},
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: "empty URN",
		},
		{
			name:    "invalid format - no colons",
			input:   "invalid-urn",
			wantErr: "invalid URN format",
		},
		{
			name:    "invalid format - no parentheses",
			input:   "urn:li:fs_salesProfile:ACwAAAKWZe8BZ8gXVKS6ePAs,NAME_SEARCH,xqt8",
			wantErr: "invalid URN format",
		},
		{
			name:    "invalid format - missing components",
			input:   "urn:li:fs_salesProfile:(ACwAAAKWZe8BZ8gXVKS6ePAs)",
			wantErr: "invalid URN format",
		},
		{
			name:    "invalid format - too many components",
			input:   "urn:li:fs_salesProfile:(part1,part2,part3,part4)",
			wantErr: "invalid URN format",
		},
		{
			name:    "malformed URN - empty parentheses",
			input:   "urn:li:fs_salesProfile:()",
			wantErr: "invalid URN format",
		},
		{
			name:    "malformed URN - missing closing parenthesis",
			input:   "urn:li:fs_salesProfile:(ACwAAAKWZe8BZ8gXVKS6ePAs,NAME_SEARCH,xqt8",
			wantErr: "invalid URN format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := EntityURN(tt.input)

			if tt.wantErr != "" {
				require.Error(t, err)
				require.Equal(t, tt.wantErr, err.Error())

				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func Test_EntityURNPanic(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("The code panicked with %v", r)
		}
	}()

	// These should not panic
	malformedInputs := []string{
		"urn:li:fs_salesProfile:(",
		"urn:li:fs_salesProfile:)",
		"urn:li:fs_salesProfile:(,)",
		"urn:li:fs_salesProfile:(,,)",
		"urn:",
		"random string",
		"urn:li:fs_salesProfile:(malformed",
	}

	for _, input := range malformedInputs {
		_, _ = EntityURN(input) // Should not panic
	}
}

func TestSanitizeForDB(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "valid UTF-8 string",
			input:    "Hello, world!",
			expected: "Hello, world!",
		},
		{
			name:     "string with invalid byte sequence",
			input:    "Invalid \x86 byte sequence",
			expected: "Invalid  byte sequence",
		}, {
			name:     "string with invalid byte sequence",
			input:    "Picus Security is the pioneer of Breach and Attack Simulation (\u0000BAS)\u0000.",
			expected: "Picus Security is the pioneer of Breach and Attack Simulation (BAS).",
		},
		{
			name:     "string with multiple invalid sequences",
			input:    "Multiple \x86 invalid \x92 sequences",
			expected: "Multiple  invalid  sequences",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "string with non-ASCII but valid UTF-8",
			input:    "こんにちは世界",
			expected: "こんにちは世界",
		},
		{
			name:     "string with emojis",
			input:    "Hello 👋 world 🌍",
			expected: "Hello 👋 world 🌍",
		},
		{
			name:     "string with BOM",
			input:    "\uFEFFHello world",
			expected: "\ufeffHello world",
		},
		{
			name:     "string with only invalid sequences",
			input:    "\x86\x92\x80",
			expected: "",
		},
		{
			name: "string with surrogate pairs",
			input: `Experience the world.

𝗦𝗼𝗹𝘂𝘁𝗶𝗼𝗻𝘀 𝗱𝗲𝘀𝗶𝗴𝗻𝗲𝗱 𝘁𝗼 𝗮𝗰𝗵𝗶𝗲𝘃𝗲 𝘁𝗵𝗲 𝗲𝘅𝗰𝗲𝗽𝘁𝗶𝗼𝗻𝗮𝗹:

Experience elevation
Commerce acceleration
Enterprise transformation
Marketing Creativity & Performance
Data Evolution

We are a network.`,
			expected: "Experience the world.\n\n     :\n\nExperience elevation\nCommerce acceleration\nEnterprise transformation\nMarketing Creativity & Performance\nData Evolution\n\nWe are a network.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := SanitizeForDB(tt.input)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestSalesCompanyURLFromURN(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		urn      string
		expected string
	}{
		{
			name:     "Valid URN",
			urn:      "urn:li:fs_salesCompany:34307789",
			expected: "https://www.linkedin.com/sales/company/34307789",
		},
		{
			name:     "Valid URN with different ID",
			urn:      "urn:li:fs_salesCompany:12345",
			expected: "https://www.linkedin.com/sales/company/12345",
		},
		{
			name:     "Invalid URN - too few parts",
			urn:      "urn:li:fs_salesCompany",
			expected: "",
		},
		{
			name:     "Invalid URN - too many parts",
			urn:      "urn:li:fs:salesCompany:34307789:extra",
			expected: "",
		},
		{
			name:     "Empty URN",
			urn:      "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := SalesCompanyURLFromURN(tt.urn)
			require.Equal(t, tt.expected, result, "salesCompanyURLFromURN(%s) result should match expected value", tt.urn)
		})
	}
}
func Test_getCountry_FranceVariants(t *testing.T) {
	t.Parallel()

	type testCase struct {
		location string
	}

	cases := []testCase{
		// Country names
		{"france"},
		{"France"},
		{"FRANCE"},
		{"fRaNcE"},
		// Alpha2 codes
		{"FR"},
		{"fr"},
		{"fR"},
		{"Fr"},
		// Alpha3 codes
		{"FRA"},
		{"fra"},
		{"fRa"},
		{"Fra"},
	}

	for _, tc := range cases {
		t.Run(tc.location, func(t *testing.T) {
			t.Parallel()

			ct, err := getCountry(context.Background(), tc.location)
			require.NoError(t, err, "getCountry should not error for %q", tc.location)
			require.Equal(t, "France", ct.Name.Common, "ct.Name.Common should be 'France' for %q", tc.location)
		})
	}
}

// TestCoalesce serves as the main entry point for all Coalesce tests.
// It uses t.Run to group tests by the data type being tested.
func TestCoalesce(t *testing.T) {
	t.Parallel()

	t.Run("Strings", testCoalesceForStrings)
	t.Run("Integers", testCoalesceForIntegers)
	t.Run("Pointers", testCoalesceForPointers)
	t.Run("Structs", testCoalesceForStructs)
}

// testCoalesceForStrings tests the Coalesce function with string types.
func testCoalesceForStrings(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		values   []string
		expected string
	}{
		{
			name:     "No arguments should return empty string",
			values:   []string{},
			expected: "",
		},
		{
			name:     "Single non-empty value",
			values:   []string{"hello"},
			expected: "hello",
		},
		{
			name:     "Single empty value",
			values:   []string{""},
			expected: "",
		},
		{
			name:     "All empty values should return empty string",
			values:   []string{"", "", ""},
			expected: "",
		},
		{
			name:     "First value is non-empty",
			values:   []string{"first", "second", "third"},
			expected: "first",
		},
		{
			name:     "Middle value is non-empty",
			values:   []string{"", "second", "third"},
			expected: "second",
		},
		{
			name:     "Last value is non-empty",
			values:   []string{"", "", "third"},
			expected: "third",
		},
		{
			name:     "All values are non-empty",
			values:   []string{"eins", "zwei", "drei"},
			expected: "eins",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			actual := Coalesce(tc.values...)
			require.Equal(t, tc.expected, actual)
		})
	}
}

// testCoalesceForIntegers tests the Coalesce function with int types.
func testCoalesceForIntegers(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		values   []int
		expected int
	}{
		{
			name:     "No arguments should return zero",
			values:   []int{},
			expected: 0,
		},
		{
			name:     "Single non-zero value",
			values:   []int{42},
			expected: 42,
		},
		{
			name:     "Single zero value",
			values:   []int{0},
			expected: 0,
		},
		{
			name:     "All zero values should return zero",
			values:   []int{0, 0, 0},
			expected: 0,
		},
		{
			name:     "First value is non-zero",
			values:   []int{100, 200, 300},
			expected: 100,
		},
		{
			name:     "Middle value is non-zero",
			values:   []int{0, -50, 100},
			expected: -50,
		},
		{
			name:     "Last value is non-zero",
			values:   []int{0, 0, 99},
			expected: 99,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			actual := Coalesce(tc.values...)
			require.Equal(t, tc.expected, actual)
		})
	}
}

// testCoalesceForPointers tests the Coalesce function with pointer types (*int).
func testCoalesceForPointers(t *testing.T) {
	t.Parallel()

	val1, val2 := 10, 20

	testCases := []struct {
		name     string
		values   []*int
		expected *int
	}{
		{
			name:     "No arguments should return nil",
			values:   []*int{},
			expected: nil,
		},
		{
			name:     "Single non-nil value",
			values:   []*int{&val1},
			expected: &val1,
		},
		{
			name:     "Single nil value",
			values:   []*int{nil},
			expected: nil,
		},
		{
			name:     "All nil values should return nil",
			values:   []*int{nil, nil, nil},
			expected: nil,
		},
		{
			name:     "First value is non-nil",
			values:   []*int{&val1, &val2, nil},
			expected: &val1,
		},
		{
			name:     "Middle value is non-nil",
			values:   []*int{nil, &val2, &val1},
			expected: &val2,
		},
		{
			name:     "Last value is non-nil",
			values:   []*int{nil, nil, &val1},
			expected: &val1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			actual := Coalesce(tc.values...)
			require.Equal(t, tc.expected, actual)
		})
	}
}

// A simple comparable struct for testing.
type user struct {
	ID   int
	Name string
}

// testCoalesceForStructs tests the Coalesce function with struct types.
func testCoalesceForStructs(t *testing.T) {
	t.Parallel()

	user1 := user{ID: 1, Name: "Alice"}
	user2 := user{ID: 2, Name: "Bob"}
	emptyUser := user{} // This is the zero value for the struct

	testCases := []struct {
		name     string
		values   []user
		expected user
	}{
		{
			name:     "No arguments should return empty struct",
			values:   []user{},
			expected: emptyUser,
		},
		{
			name:     "Single non-empty value",
			values:   []user{user1},
			expected: user1,
		},
		{
			name:     "Single empty value",
			values:   []user{emptyUser},
			expected: emptyUser,
		},
		{
			name:     "All empty values should return empty struct",
			values:   []user{emptyUser, {}, user{}},
			expected: emptyUser,
		},
		{
			name:     "First value is non-empty",
			values:   []user{user1, user2, emptyUser},
			expected: user1,
		},
		{
			name:     "Middle value is non-empty",
			values:   []user{emptyUser, user2, user1},
			expected: user2,
		},
		{
			name:     "Last value is non-empty",
			values:   []user{emptyUser, user{}, user1},
			expected: user1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			actual := Coalesce(tc.values...)
			require.True(t, reflect.DeepEqual(actual, tc.expected))
		})
	}
}
