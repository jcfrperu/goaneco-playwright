package playwright

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var textualMimeTypeRe = regexp.MustCompile(`^(text\/.*?|application\/(json|(x-)?javascript|xml.*?|ecmascript|graphql|x-www-form-urlencoded)|image\/svg(\+xml)?|application\/.*?(\+json|\+xml))(;\s*charset=.*)?$`)

// APIResponseAssertions provides assertion methods for an APIResponse.
type APIResponseAssertions struct {
	response *APIResponse
	isNot    bool
}

// ExpectAPIResponse creates a new APIResponseAssertions for the given response.
func ExpectAPIResponse(response *APIResponse) *APIResponseAssertions {
	return &APIResponseAssertions{response: response}
}

// Not returns a negated version of the assertion.
func (a *APIResponseAssertions) Not() *APIResponseAssertions {
	return &APIResponseAssertions{response: a.response, isNot: true}
}

// ToBeOK asserts that the response status is in the 200-299 range.
func (a *APIResponseAssertions) ToBeOK(ctx context.Context) error {
	if a.isNot != a.response.OK() {
		return nil
	}
	message := fmt.Sprintf("Response status expected to be within [200..299] range, was %d", a.response.Status())
	if a.isNot {
		message = strings.ReplaceAll(message, "expected to", "expected not to")
	}
	logLine := fmt.Sprintf("→ %s %s\n← %d %s",
		a.response.method, a.response.URL(),
		a.response.Status(), a.response.StatusText())
	message += "\nCall log:\n" + logLine

	contentType := a.response.headers["content-type"]
	if isTextualMimeType(contentType) {
		text, err := a.response.Text(ctx)
		if err == nil {
			message += fmt.Sprintf("\nResponse text:\n%s", subString(text, 0, 1000))
		}
	}
	return errors.New(message)
}

func isTextualMimeType(mimeType string) bool {
	return textualMimeTypeRe.MatchString(mimeType)
}

func subString(s string, start, length int) string {
	if start < 0 {
		start = 0
	}
	if length < 0 {
		length = 0
	}
	rs := []rune(s)
	end := start + length
	if end > len(rs) {
		end = len(rs)
	}
	return string(rs[start:end])
}
