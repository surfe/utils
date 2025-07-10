package utils

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/JohannesKaufmann/html-to-markdown/plugin"
	"github.com/PuerkitoBio/goquery"
	"github.com/gosimple/unidecode"
	"github.com/jpillora/go-tld"
	"github.com/kennygrant/sanitize"
	"github.com/pariz/gountries"
	"github.com/surfe/logger/v2"
	"github.com/surfe/utils/utils/urls"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

const (
	clearbitLogoURL = "https://logo.clearbit.com/"
)

var (
	ErrEmptyURL    = errors.New("empty URL")
	ErrNoRedirects = errors.New("no redirects found")
)

func MaskEmail(str string) string {
	parts := strings.Split(str, "@")
	if len(parts) != 2 || len(parts[0]) < 3 || len(parts[1]) < 3 {
		return "Invalid email"
	}

	username := parts[0]
	domain := parts[1]
	return fmt.Sprintf("%s%s%s@%s",
		username[:1],
		strings.Repeat("*", len(username)-2),
		username[len(username)-1:],
		domain,
	)
}

func MaskStringWithAsterisks(str string) string {
	if len(str) <= 2 {
		return "**"
	}
	return str[0:1] + strings.Repeat("*", len(str)-2) + str[len(str)-1:]
}

func ToValidEmail(email string) (string, bool) {
	email = strings.TrimSpace(strings.ToLower(email))
	return email, reEmail.MatchString(email)
}

// Title capitalizes the first letter of a string
func Title(s string) string {
	if len(s) == 0 {
		return s
	}

	parts := strings.Split(s, " ")
	runes := make([]rune, 0, len(s))
	for _, p := range parts {
		runesP := []rune(p)
		if len(runesP) == 0 {
			continue
		}

		runesP[0] = unicode.ToUpper(runesP[0])
		runes = append(runes, runesP...)
		runes = append(runes, ' ')
	}
	return strings.TrimSpace(string(runes))
}

func DomainFromURL(s string) string {
	if strings.TrimSpace(s) == "" {
		return ""
	}

	domain, err := DomainFromURLNoFiltering(s)
	if err != nil {
		if !errors.Is(err, ErrEmptyURL) {
			logger.Log(context.Background()).Err(err).Infof("DomainFromURLNoFiltering failed for %s", s)
		}
		return ""
	}

	if urls.IsURLShortenerDomain(domain) {
		if redirectedDomain, err := getRedirectedDomain(s); err == nil {
			domain = redirectedDomain
		} else {
			logger.Log(context.Background()).Err(err).Errorf("GetRedirectedDomain: %s", s)
		}
	}

	if urls.IsPublicDomain(domain) {
		return ""
	}

	if knownDomain, exists := knownDomains[domain]; exists {
		return knownDomain
	}

	return domain
}

func SubdomainWithDomainFromURL(s string) (string, error) {
	s = strings.ToLower(s)
	if !reWebSchema.MatchString(s) {
		s = "//" + s // URL needs to prefixed with `//` to be parseable
	}
	u, err := tld.Parse(s)
	if err != nil {
		return "", fmt.Errorf("parse domain from URL: %w", err)
	}

	if u.Domain == "" && u.TLD == "" {
		return "", fmt.Errorf("empty domain and TLD")
	}

	if (u.Subdomain == "") || (strings.ToLower(u.Subdomain) == "www") {
		return fmt.Sprintf("%s.%s", u.Domain, u.TLD), nil
	}

	return fmt.Sprintf("%s.%s.%s", u.Subdomain, u.Domain, u.TLD), nil
}

func DomainFromURLNoFiltering(s string) (string, error) {
	if s == "" {
		return "", ErrEmptyURL
	}

	s = strings.ToLower(s)
	if !reWebSchema.MatchString(s) {
		s = "//" + s // URL needs to be prefixed with `//` to be parseable
	}

	u, err := tld.Parse(s)
	if err != nil {
		if strings.Contains(err.Error(), "cannot derive eTLD+1") {
			tempURL, urlErr := url.Parse(s)

			if urlErr == nil && tempURL.Host != "" && strings.Contains(tempURL.Host, ".") {
				return tempURL.Host, nil
			}
		}

		return "", fmt.Errorf("parse domain from URL: %w", err)
	}

	if u.Domain == "" && u.TLD == "" {
		return "", fmt.Errorf("empty domain and TLD")
	}

	return u.Domain + "." + u.TLD, nil
}

func DomainFromEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return ""
	}
	return parts[1]
}

func FormatDomainURL(domainURL string) string {
	parsedURL, err := url.Parse(domainURL)
	if err != nil || parsedURL == nil || parsedURL.Hostname() == "" {
		return domainURL
	}

	return parsedURL.Hostname()
}

// DomainNameWithoutTLD extracts the name of the domain, e.g. input: `https://www.surfe.com/some-path`, output: `surfe`
func DomainNameWithoutTLD(rawURL string) string {
	if !strings.HasPrefix(rawURL, "http") {
		rawURL = "http://" + rawURL
	}
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}

	host := parsedURL.Hostname()
	parts := strings.Split(host, ".")

	if len(parts) > 1 {
		return parts[0]
	}
	return host
}

func RemoveQueryParams(s string) string {
	u, err := url.Parse(s)
	if err != nil {
		return s
	}

	u.RawQuery = ""
	u.Path = strings.TrimSuffix(u.Path, "/")
	return u.String()
}

func SameDomains(linkedinDomain, crmDomain string) bool {
	if linkedinDomain == "" || crmDomain == "" {
		return false
	}

	if v := DomainFromURL(crmDomain); v != "" {
		return strings.EqualFold(v, linkedinDomain)
	}

	return strings.EqualFold(crmDomain, linkedinDomain)
}

func RemoveAccents(s string) (string, error) {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	output, _, err := transform.String(t, s)
	return output, err
}

func SanitizeName(s string) string {
	s = reNameSanitize.ReplaceAllString(s, "")
	s = allEmojisRe.ReplaceAllString(s, "")  // Remove all emojis
	s = strings.Join(strings.Fields(s), " ") // Replace consecutive white space characters with one
	s = strings.TrimSpace(s)                 // Remove leading and trailing white spaces
	s = strings.ToValidUTF8(s, "")

	return s
}

// SimplifyName removes leading and trailing abbreviations from given name after sanitizing
func SimplifyName(s string) string {
	s = SanitizeName(s)

	// Firstname; e.g. Dr. John
	if reAbbrvPrefix.MatchString(s) {
		s = reAbbrvPrefix.ReplaceAllString(s, "")
	}

	// Lastname; e.g. Poirot, MBA, II, III, IV FRSA, F.R.S.A
	if reAbbrvSuffix.MatchString(s) {
		s = reAbbrvSuffix.ReplaceAllString(s, "")
	}

	s = strings.TrimSpace(s)
	s = strings.TrimFunc(s, func(r rune) bool { // Trim non-letter leading and trailing chars
		return !unicode.IsLetter(r)
	})

	return s
}

// SimplifyCompanyName removes abbreviations from given company name
func SimplifyCompanyName(originalName string) string {
	s := originalName
	if reRemoveAbbreviationsWithDots.MatchString(s) {
		s = reRemoveAbbreviationsWithDots.ReplaceAllString(s, "")
	}

	if reRemoveAbbreviationsWithoutDots.MatchString(s) {
		s = reRemoveAbbreviationsWithoutDots.ReplaceAllString(s, "")
	}

	s = reRemove.ReplaceAllString(s, "")

	if len(s) < 3 {
		return originalName
	}
	return s
}

// SanitizeForSOQL escapes chars which are valid in SOQL queries
func SanitizeForSOQL(s string) string {
	s = strings.ReplaceAll(s, "'", `\'`) // Single-quote
	s = strings.ReplaceAll(s, `"`, `\"`) // Double-quote
	s = strings.ReplaceAll(s, "_", `\_`) // Underscore
	s = strings.ReplaceAll(s, "%", `\%`) // Percent sign
	return s
}

func LooseString(s string) (text string) {
	text = strings.ReplaceAll(s, "&", "and")
	text = strings.ReplaceAll(text, "@", "at")
	text = unidecode.Unidecode(text)
	text = strings.ToLower(text)
	text = removeNonAlphanumericCharacters(text)
	return text
}

func SanitizedLooseString(s string) string {
	return LooseString(sanitize.HTML(s))
}

func TrimSuffix(s, suffix string) (string, bool) {
	if strings.HasSuffix(s, suffix) {
		return s[:len(s)-len(suffix)], true
	}
	return s, false
}

// URLHostnameExtractor extracts hostname from given URL string, and removes `www.` prefix if exists
func URLHostnameExtractor(s string) string {
	if s == "" {
		return ""
	}
	if !reWebSchema.MatchString(s) {
		s = "//" + s // URL needs to prefixed with `//` to be parseable
	}
	u, err := url.Parse(s)
	if err != nil {
		return ""
	}
	return reWWW.ReplaceAllString(u.Hostname(), "")
}

func URLProfileExtract(s string) string {
	if s == "" {
		return ""
	}

	// Return as is if already a handle
	if reHandle.MatchString(s) {
		return s
	}

	m := reLinkedinType.FindStringSubmatch(s)
	if len(m) < 3 {
		return ""
	}
	res, _ := url.QueryUnescape(m[2])

	return res
}

// ExtractLinkedInSlug extracts the LinkedIn slug from a given string.
func ExtractLinkedInSlug(s string) string {
	if s == "" {
		return ""
	}

	// Clean the URL by removing fragments and trailing slashes.
	s = strings.SplitN(s, "#", 2)[0]
	s = strings.TrimSuffix(s, "/")

	// Return if it's already a handle.
	if reHandle.MatchString(s) {
		return s
	}

	// Remove trailing text added by the user.
	s = strings.Fields(s)[0]

	// Extract the LinkedIn slug using regex.
	matches := reLinkedinType.FindStringSubmatch(s)
	if len(matches) < 3 {
		return ""
	}

	// Decode any URL-encoded characters in the result.
	slug, _ := url.QueryUnescape(matches[2])

	return slug
}

// ContactProfileURL creates link with provided ID. Handle, MemberID or SalesNavID will work well with this link.
func ContactProfileURL(id string) string {
	return fmt.Sprintf("https://linkedin.com/in/%s", id)
}

// OrganizationProfileURL creates link with provided ID. Handle, MemberID or SalesNavID will work well with this link.
// When ID is provided link with /company/ will also work for schools and redirect to correct organization.
func OrganizationProfileURL(id string) string {
	if id == "" {
		return ""
	}
	return fmt.Sprintf("https://linkedin.com/company/%s", id)
}

// URNExtractor extracts ID from URN e.g. 13205888 from urn:li:fs_normalized_company:13205888
// or ACwAAAJlc6wBYdHGFmVJDHu from urn:li:fs_salesProfile:(ACwAAAJlc6wBYdHGFmVJDHu,NAME_SEARCH,ij9X)
var URNExtractor = func(s string) string {
	a := strings.Split(s, ":")
	if len(a) == 0 {
		return ""
	}

	lastSegment := a[len(a)-1]

	re := regexp.MustCompile(`\((.*?)\)`) // urn:li:fs_salesProfile:(ACwAAAJlc6wBYdHGFmVJDHu,NAME_SEARCH,ij9X)
	matches := re.FindStringSubmatch(lastSegment)
	if len(matches) < 2 {
		return lastSegment
	}

	splitContent := strings.Split(matches[1], ",")
	return strings.TrimSpace(splitContent[0])
}

func SalesProfileURLFromURN(s string) string {
	urn, err := EntityURN(s)
	if err != nil {
		return ""
	}
	return salesPeopleProfilePageRoot + "/" + urn.ProfileID
}

// SalesCompanyURLFromURN converts URN like `urn:li:fs_salesCompany:34307789` to URL `https://www.linkedin.com/sales/company/34307789`
func SalesCompanyURLFromURN(s string) string {
	id := ExtractSalesProfileIDFromURN(s)
	if id == "" {
		return ""
	}
	return salesCompanyProfilePageRoot + "/" + id
}

const salesPeopleProfilePageRoot = "https://www.linkedin.com/sales/people"
const salesCompanyProfilePageRoot = "https://www.linkedin.com/sales/company"

// ExtractSalesProfileIDFromURN extracts e.g. 34307789 from urn:li:fs_salesCompany:34307789
func ExtractSalesProfileIDFromURN(s string) string {
	urnArr := strings.Split(s, ":")
	if len(urnArr) != 4 {
		return ""
	}
	return urnArr[3]
}

// ExtractSalesNavIDFromURL extracts e.g. ACwAAAB0jjIBhkHdtLPXAimZXL_SAtpk7UFZbik from https://www.linkedin.com/sales/people/ACwAAAB0jjIBhkHdtLPXAimZXL_SAtpk7UFZbik,OUT_OF_NETWORK,2N7m
func ExtractSalesNavIDFromURL(s string) string {
	if s == "" {
		return ""
	}

	parts := strings.Split(s, ",")
	if len(parts) == 0 {
		return ""
	}
	urlParts := strings.Split(parts[0], "/")
	if len(urlParts) == 0 {
		return ""
	}
	return urlParts[len(urlParts)-1]
}

func Hash(o interface{}) string {
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%v", o)))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func HashByte(o interface{}) []byte {
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%v", o)))
	return h.Sum(nil)
}

func StructToStringJSON(s interface{}) (string, error) {
	if s == nil {
		return "", nil
	}
	out, err := json.Marshal(s)
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func StringToStructJSON(s string, i interface{}) error {
	if s == "" {
		return nil
	}
	return json.Unmarshal([]byte(s), i)
}

func StructToStringGob(s interface{}) (string, error) {
	if s == nil {
		return "", nil
	}
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(s)
	return buf.String(), err
}

func StringToStructGob(s string, i interface{}) error {
	if s == "" {
		return nil
	}
	return gob.NewDecoder(bytes.NewBuffer([]byte(s))).Decode(i)
}

func FirstNameFromFullName(name string) string {
	n := strings.Split(name, " ")
	if len(n) == 0 {
		return ""
	}
	return n[0]
}

func FirstAndLastNameFromFullName(name string) (string, string) {
	name = strings.TrimSpace(name)
	n := strings.Split(name, " ")
	if len(n) == 0 {
		return "", ""
	}

	var lastNameParts []string
	for i := range n {
		n[i] = strings.TrimSpace(n[i])
		if i > 0 && n[i] != "" {
			lastNameParts = append(lastNameParts, n[i])
		}
	}

	return n[0], strings.Join(lastNameParts, " ")
}

// ConsiderNonEmpty returns first param if it's not empty, second otherwise
func ConsiderNonEmpty(str string, def string) string {
	if str != "" {
		return str
	}

	return def
}

// LinkedinURLCleaner returns clean url with only scheme, hostname, and the path
var LinkedinURLCleaner = func(rawUrl string) string {
	cleanUrl, err := LinkedinURLCleanerErr(rawUrl, false)
	if err != nil {
		return ""
	}

	return cleanUrl
}

// LinkedinURLCleanerErr cleans (and escapes handle fragment if requested) of the LinkedIn URL to get a consistent id.
// If provided rawUrl is not a LinkedIn URL then returns the rawUrl and an error.
func LinkedinURLCleanerErr(rawUrl string, escapeHandle bool) (string, error) {
	cleanUrl := reLinkedinURL.FindString(rawUrl)
	if cleanUrl == "" {
		return rawUrl, errors.New("not a LinkedIn URL")
	}

	sArr := strings.Split(cleanUrl, "/")
	urlArr := sArr[:len(sArr)-1]
	handle := sArr[len(sArr)-1]

	if escapeHandle {
		// Unescape first to make sure we won't double escape
		handle, _ = url.PathUnescape(handle)
		handle = url.PathEscape(handle)
	}

	urlArr = append(urlArr, handle)
	return strings.Join(urlArr, "/"), nil
}

// ExtractHostAndPath takes a string containing a URL, and returns another string with the same URL without
// the scheme (http/https), the www. subdomain and any trailing /
// If the provided string cannot be parsed as URL, it gets returned without any trailing /
func ExtractHostAndPath(fullURL string) string {
	urlWithoutTrailingSlash, _ := strings.CutSuffix(fullURL, "/")
	parsedURL, err := url.Parse(urlWithoutTrailingSlash)
	if err != nil {
		return urlWithoutTrailingSlash
	}

	host, _ := strings.CutPrefix(parsedURL.Host, "www.")
	return host + parsedURL.Path
}

// GenerateURLCombinations takes a string containing a URL, and returns an array with all different formats
// that same URL could take while remaining valid. The function does not check if the provided string is an actual URL
func GenerateURLCombinations(fullURL string) []string {
	if fullURL == "" {
		return []string{}
	}

	hostAndPath := ExtractHostAndPath(fullURL)
	combinations := []string{
		"https://" + "www." + hostAndPath + "/", // https + www + trailing /
		"http://" + "www." + hostAndPath + "/",  // http + www + trailing /
		"https://" + hostAndPath + "/",          // https + trailing /
		"http://" + hostAndPath + "/",           // http + trailing /
		"https://" + "www." + hostAndPath,       // https + www
		"http://" + "www." + hostAndPath,        // http + www
		"https://" + hostAndPath,                // https
		"http://" + hostAndPath,                 // http
	}
	return combinations
}

// TwitterURLBuilder builds Twitter profile URL from username
var TwitterURLBuilder = func(username string) string {
	return "https://twitter.com/" + username
}

// MatchLinkedinURL checks if provided URLs are matching
func MatchLinkedinURL(url1, url2 string) bool {
	// Delete trailing slash
	url1 = strings.TrimRight(url1, "/")
	url2 = strings.TrimRight(url2, "/")

	// Check if valid linkedin URLs
	return url1 == url2 && LinkedinURLCleaner(url1) != ""
}

// MatchLIURLByIDOrHandle checks if the given string matches a LinkedIn URL that contains the provided ID or handle.
func MatchLIURLByIDOrHandle(s string, idOrHandle string) bool {
	if idOrHandle == "" {
		return false
	}

	urlRegex := `^(?:https?:\/\/)?(?:www\.)?linkedin\.com.*?\b` + idOrHandle + `\b`
	reURL := regexp.MustCompile(urlRegex)
	return reURL.MatchString(s)
}

// TruncateString truncates given string to given length after removing emojis
func TruncateString(str string, length int) string {
	if str == "" {
		return ""
	}
	str = allEmojisRe.ReplaceAllString(str, "") // Length of emojis are not consistent between browsers/langs, so better remove all first

	r := []rune(str)
	if length > 0 && utf8.RuneCountInString(str) > length {
		return string(r[:length-utf8.RuneCountInString(DefaultOmission)]) + DefaultOmission
	}

	return string(r)
}

func ExtractNumbersFromString(str string) string {
	nums := reNum.FindAllString(str, -1)
	return strings.Join(nums, "")
}

// ConvertMDToHTML converts Markdown syntax to HTML
func ConvertMDToHTML(txt string) string {
	if len(txt) == 0 {
		return txt
	}

	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,         // GitHub flavored markdown
			extension.Typographer, // Replace punctuations with typographic entities
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
	)
	var buf bytes.Buffer
	if err := md.Convert([]byte(txt), &buf); err != nil {
		logger.Log(context.Background()).Err(err).Error("Failed to convert MD to HTML")
		return txt
	}

	return buf.String()
}

// ConvertHTMLToMD converts HTML to Markdown syntax
func ConvertHTMLToMD(html, domain string) string {
	if len(html) == 0 {
		return html
	}

	u, err := url.Parse(domain)
	if err != nil {
		domain = ""
	} else {
		domain = u.Host
	}

	converter := md.NewConverter(domain, true, nil)
	converter.Use(plugin.GitHubFlavored())
	converter.AddRules(
		md.Rule{
			Filter: []string{"div"},
			Replacement: func(content string, selec *goquery.Selection, opt *md.Options) *string {
				content = strings.TrimSpace(content)
				return md.String("\n" + content + "\n")
			},
		},
		md.Rule{
			Filter: []string{"br"},
			Replacement: func(content string, selec *goquery.Selection, opt *md.Options) *string {
				return md.String("\n")
			},
		},
	)

	markdown, err := converter.ConvertString(html)
	if err != nil {
		logger.Log(context.Background()).Err(err).Error("Failed to convert HTML to MD")
		return html
	}

	return markdown
}

// ExtractPlanInfo returns plan name and plan interval from plan id
func ExtractPlanInfo(planID string) (planName string, planInterval string) {
	// Use the regular expressions to extract the plan name and plan interval from the plan ID string
	planInterval = rePlanInterval.FindString(planID)
	planName = rePlanName.FindString(planID)
	if planName == "entreprise" {
		planName = "enterprise"
	}

	return
}

func AddPrefixIfMissing(str, prefix string) string {
	if !strings.HasPrefix(str, prefix) {
		return prefix + str
	}

	return str
}

func ContainsIgnoreCase(dataMap map[string]string, target string) bool {
	for key := range dataMap {
		if strings.EqualFold(key, target) {
			return true
		}
	}
	return false
}

func GetNewValueIfEmpty(currentValue *string, newValue *string) *string {
	if currentValue == nil || *currentValue == "" {
		return newValue
	}
	return currentValue
}

func GetErrorString(err error) string {
	if err != nil {
		return err.Error()
	}
	return ""
}

func GetRedirectedDomainFromDomain(domain string) (string, error) {
	domain, err := DomainFromURLNoFiltering(domain)
	if err != nil {
		return "", fmt.Errorf("domain form url no filtering: %w", err)
	}

	if !strings.HasPrefix(domain, "http") {
		domain = "https://" + domain
	}

	redirectedDomain, err := getRedirectedDomain(domain)
	if err != nil {
		return "", fmt.Errorf("get redirected domain: %w", err)
	}

	return redirectedDomain, nil
}

func getRedirectedDomain(url string) (string, error) {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Head(url)
	if err != nil {
		return "", fmt.Errorf("failed to get redirection URL for %s: %w", url, err)
	}

	defer resp.Body.Close()

	if resp.StatusCode < 300 || resp.StatusCode >= 400 {
		return "", fmt.Errorf("no redirection found for %s with status code: %d, %w", url, resp.StatusCode, ErrNoRedirects)
	}

	redirectURL := resp.Header.Get("Location")
	domain, err := DomainFromURLNoFiltering(redirectURL)
	if err != nil {
		return "", fmt.Errorf("failed to get domain from URL %s: %w", redirectURL, err)
	}

	return domain, nil
}

const DefaultOmission = "…"

var knownDomains = map[string]string{
	"goo.gle": "google",
}

// Abbreviations which we don't want in contact first or last names
var abbrvs = "PH.D.|PHD|MD|CPA|CMA|PROF|PR|MBA|PHR|MA|BFA|PMP|MSM|TMP|RN|CFRE|PLS|MSW|CEC|HCS|CFP|AAMS|CLU|ChFC|M.P.A.|MLEC|MAQP|MSHR|SHRM-SCP|MG|MS|CSP|CAS|MAS|LDN|LPN|DC|JR|SR|CIR|A.C.C.|M.Ed.|M.A.I.|AI-GRS|JD|PE|CCP|CAA|LUTCF|FSS|MHR|FACS|MHA|PT|DPT|CDAL|CVM|LPC|CIC|SIOR|CPM|GC|CHHC|AADP|MPA|PE|BASI|CFRE|CMPE|FACHE|CAPS|CEPA|MSOM|IPMA-SCP|CME|ITIL|PMA|DR|II|III|IV|FRSA|F.R.S.A|LL.M.|CFA|MFE|CXAP"

// Abbreviations which we don't want in companyName
var companyAbbrvs = []string{
	"A/S", "AG", "AB", "AE", "ApS", "AS", "BV", "Co", "Corp", "CV", "EEIG", "GmbH", "Inc", "K/S", "Ltd", "Oy", "PLC",
	"Pty Ltd", "SE", "SP", "SRL", "KGaA", "LLP", "SARL", "SàRL", "SCE", "SCOP", "SCEA", "SNC", "SOCIMI",
	"SpA", "Zrt", "Coop", "GbR", "HUF", "IBC", "I/S", "JSC", "KDA", "KG", "Kommanditgesellschaft", "LDA", "LLLP", "NUF",
	"OJSC", "OOD", "OÜ", "LLC", "AG", "SA", "limited", "Zweigniederlassung", "Genossenschaft", "Verein", "Stiftung", "SAS",
}

var companyAbbrvsWithDots = []string{
	"A.G.", "K.K.", "L.L.C.", "L.P.", "N.V.", "S.A.S.", "S.A.", "S.C.S.", "S.C.", "S.P.A.", "S.R.L.", "S.à r.l.",
	"Sp. z o.o.", "V.O.F.", "B.V.B.A.", "C.A.", "C.C.", "C.V.A.", "C.F.", "C.L.", "C.S.", "E.A.D.", "F.Z.E.",
	"K.F.T.", "K.D.", "M.B.", "N.P.", "O.E.", "P.C.", "P.L.L.C.", "P.S.C.", "P.C.C.", "P.C.G.", "Q.S.C.", "R.L.",
	"S.A.U.", "S.A.R.F.", "SE & Co. KGaA", "SE & Co. KG", "AG & Co. KGaA", "AG & Co. KG", "GmbH & Co. KGaA", "GmbH & Co. KG",
}

var (
	reNum          = regexp.MustCompile("[0-9]+")
	reWWW          = regexp.MustCompile(`(?:www\.)?`)
	reWebSchema    = regexp.MustCompile(`^http(s)?:\/\/`)
	reLinkedinURL  = regexp.MustCompile(`http(s)?://([\w]+\.)?linkedin\.com/(pub|in|profile|company|school)/[^/?\s]+`)
	rePlanName     = regexp.MustCompile("basic|starter|professional|enrich|business|entreprise|enterprise|pro|essential")
	rePlanInterval = regexp.MustCompile("monthly|yearly")
	reEmail        = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,12}$`)
	reNameSanitize = regexp.MustCompile(`\([^\)]*\)`)
	reAbbrvPrefix  = regexp.MustCompile(`^(?i)((` + abbrvs + `)+[ |,|.]+)+`)
	reAbbrvSuffix  = regexp.MustCompile(`(?i)([ |,|.]+(` + abbrvs + `))+$`)
	reLinkedinType = regexp.MustCompile(`linkedin\.com\/(pub|in|profile|company|school)\/([^\/\ ?]+)`)
	reHandle       = regexp.MustCompile(`^[\p{L}0-9-]+(?:-[\p{L}0-9]+)*$`)

	// Remove dots, commas, and spaces from the beginning and end of the string
	reRemove = regexp.MustCompile(`^[.,\s]+|[.,\s]+$`)

	// Remove abbreviations with dots
	reRemoveAbbreviationsWithDots = regexp.MustCompile(fmt.Sprintf(`(?i)\b%s\b`,
		strings.ReplaceAll(strings.Join(companyAbbrvsWithDots, "|"), ".", "\\.")))

	// Remove abbreviations without dots
	reRemoveAbbreviationsWithoutDots = regexp.MustCompile(`(?i)([ ]|^)(` + strings.Join(companyAbbrvs, "|") + `)($|\.|\,|\b)`)
)

func removeNonAlphanumericCharacters(s string) string {
	var result strings.Builder
	for i := 0; i < len(s); i++ {
		b := s[i]
		if ('a' <= b && b <= 'z') || ('A' <= b && b <= 'Z') || ('0' <= b && b <= '9') {
			result.WriteByte(b)
		}
	}
	return result.String()
}

func GetCountryAlpha2(ctx context.Context, country, location string) string {
	if country != "" {
		return country
	}

	if location == "" {
		return ""
	}

	ct, err := getCountry(ctx, location)
	if err != nil {
		return ""
	}

	return ct.Alpha2
}

func GetCountryCommonName(ctx context.Context, country, location string) string {
	if country != "" {
		return country
	}

	if location == "" {
		return ""
	}

	ct, err := getCountry(ctx, location)
	if err != nil {
		return ""
	}

	return ct.Name.Common
}

func getCountry(ctx context.Context, location string) (gountries.Country, error) {
	locs := strings.Split(location, ",")
	loc := strings.TrimSpace(locs[len(locs)-1])
	loc = strings.ToLower(loc)

	query := gountries.New()
	ct, err := query.FindCountryByName(loc)
	if err != nil && !strings.Contains(err.Error(), "Could not find country with name") {
		logger.Log(ctx).Err(err).Error("Find country by name")
		return gountries.Country{}, err
	}

	if err != nil {
		ct, err = query.FindCountryByAlpha(loc)
		if err != nil {
			logger.Log(ctx).Err(err).Error("Find country by alpha")
			return gountries.Country{}, err
		}
	}

	return ct, nil
}

func GetLocationWithoutCountry(ctx context.Context, location string) string {
	if location == "" {
		return ""
	}

	if GetCountryAlpha2(ctx, "", location) == "" {
		return location
	}

	locs := strings.Split(location, ",")

	return strings.Join(locs[0:len(locs)-1], ",")
}

func GetCountryNameByAlpha2(ctx context.Context, alpha2 string) (string, error) {
	if len(alpha2) != 2 {
		return "", errors.New("code passed is not alpha2")
	}

	query := gountries.New()
	ct, err := query.FindCountryByAlpha(alpha2)
	if err != nil {
		return "", fmt.Errorf("gountries: find country by alpha: %w", err)
	}

	return ct.Name.Common, nil
}

func GetClearbitURL(url string) string {
	domain, err := DomainFromURLNoFiltering(url)
	if err != nil || domain == "" {
		return ""
	}

	return fmt.Sprintf("%s%s", clearbitLogoURL, domain)
}

func IsLinkedInURL(rawUrl string) bool {
	if rawUrl == "" {
		return false
	}
	return reLinkedinURL.MatchString(rawUrl)
}

// EntityURN converts urn string formatted like `urn:li:fs_salesProfile:(ACwAAAKWZe8BZ8gXVKS6ePAs8I4GWmjW4Tjm-7w,NAME_SEARCH,xqt8)` to URN struct
func EntityURN(s string) (URN, error) {
	if s == "" {
		return URN{}, errors.New("empty URN")
	}

	a := strings.Split(s, ":")
	if len(a) == 0 {
		return URN{}, errors.New("invalid URN format")
	}

	lastSegment := a[len(a)-1]

	re := regexp.MustCompile(`\((.*?)\)`) // urn:li:fs_salesProfile:(ACwAAAJlc6wBYdHGFmVJDHu,NAME_SEARCH,ij9X)
	matches := re.FindStringSubmatch(lastSegment)
	if len(matches) == 0 {
		return URN{}, errors.New("invalid URN format")
	}

	split := strings.Split(matches[1], ",")
	if len(split) != 3 {
		return URN{}, errors.New("invalid URN format")
	}

	return URN{
		ProfileID: split[0],
		AuthType:  split[1],
		AuthToken: split[2],
	}, nil
}

type URN struct {
	ProfileID string
	AuthType  string
	AuthToken string
}

func SanitizeForDB(s string) string {
	var builder strings.Builder
	builder.Grow(len(s))

	for _, r := range s {
		switch {
		case r == '\u0000':
			// Remove null bytes completely
			continue
		case unicode.IsControl(r) && r != '\n' && r != '\t' && r != '\r':
			// Replace other control characters (except newlines and tabs) with spaces
			builder.WriteRune(' ')
		case r == '\uFFFD':
			// Replace the Unicode replacement character
			continue
		case r >= 0xD800 && r <= 0xDFFF:
			// Remove surrogate pairs
			continue
		case r >= 0x1D400 && r <= 0x1D7FF:
			// Remove Mathematical Script/Bold characters
			continue
		default:
			builder.WriteRune(r)
		}
	}

	return builder.String()
}
