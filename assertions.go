package playwright

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/jcfrperu/goaneco-playwright/protocol"
)

// flagPrefixRegexp matches pure inline flag prefixes like (?i), (?im), (?ims) in Go regexp strings.
var flagPrefixRegexp = regexp.MustCompile(`^\(\?([imsU-]+)\)`)

// LocatorAssertions provides web-first assertion methods for a Locator.
// Assertions use the Frame.expect IPC, which retries until the assertion passes
// or the timeout is exceeded.
type LocatorAssertions struct {
	locator *Locator
	isNot   bool
}

// Expect creates a new LocatorAssertions for the given Locator.
func Expect(locator *Locator) *LocatorAssertions {
	return &LocatorAssertions{locator: locator}
}

// Not returns a negated version of the assertion.
func (a *LocatorAssertions) Not() *LocatorAssertions {
	return &LocatorAssertions{locator: a.locator, isNot: true}
}

// frameExpectRequestWire is the wire format for Frame.expect IPC.
// It uses omitempty on fields that must be absent (not null) when not provided,
// particularly the Pseudo field which the server validates as "before"|"after"|absent.
// The server requires the Timeout field (must be a float, not undefined).
type frameExpectRequestWire struct {
	ExpectedNumber *float64                     `json:"expectedNumber,omitempty"`
	ExpectedText   []protocol.ExpectedTextValue `json:"expectedText"`
	ExpectedValue  json.RawMessage              `json:"expectedValue,omitempty"` // raw JSON for flexible serialization
	Expression     string                       `json:"expression"`
	ExpressionArg  any                          `json:"expressionArg,omitempty"`
	IsNot          bool                         `json:"isNot"`
	Pseudo         *string                      `json:"pseudo,omitempty"` // omitempty: absent when not "before"/"after"
	Selector       *string                      `json:"selector,omitempty"`
	Timeout        float64                      `json:"timeout"` // required by server; 0 = use server default (no timeout)
	UseInnerText   *bool                        `json:"useInnerText,omitempty"`
}

// defaultExpectTimeout is the default timeout for assertions in milliseconds.
// Matches Playwright's default expect timeout of 5 seconds.
const defaultExpectTimeout = 5000.0

// sendExpect sends a Frame.expect IPC request and returns any error.
func (a *LocatorAssertions) sendExpect(ctx context.Context, req frameExpectRequestWire) error {
	req.IsNot = a.isNot
	req.Selector = &a.locator.selector
	if req.Timeout == 0 {
		req.Timeout = defaultExpectTimeout
	}

	_, err := a.locator.frame.SendMessageRequest(ctx, "expect", req)
	if err != nil {
		notStr := ""
		if a.isNot {
			notStr = "not."
		}
		return fmt.Errorf("expect(%s%s) failed: %w", notStr, req.Expression, err)
	}
	return nil
}

// checkedExpectedValue builds the serialized argument for to.be.checked assertion.
// Playwright's server expects: { value: { o: [{k: "checked", v: {b: <bool>}}] }, handles: [] }
func checkedExpectedValue(checked bool) json.RawMessage {
	msg, err := json.Marshal(map[string]any{
		"value": map[string]any{
			"o": []any{
				map[string]any{
					"k": "checked",
					"v": map[string]any{"b": checked},
				},
			},
		},
		"handles": []any{},
	})
	if err != nil {
		return nil
	}
	return msg
}

// ToBeVisible asserts that the element is visible.
func (a *LocatorAssertions) ToBeVisible(ctx context.Context) error {
	return a.sendExpect(ctx, frameExpectRequestWire{
		Expression:   "to.be.visible",
		ExpectedText: []protocol.ExpectedTextValue{},
	})
}

// ToBeEnabled asserts that the element is enabled.
func (a *LocatorAssertions) ToBeEnabled(ctx context.Context) error {
	return a.sendExpect(ctx, frameExpectRequestWire{
		Expression:   "to.be.enabled",
		ExpectedText: []protocol.ExpectedTextValue{},
	})
}

// ToBeDisabled asserts that the element is disabled.
func (a *LocatorAssertions) ToBeDisabled(ctx context.Context) error {
	return a.sendExpect(ctx, frameExpectRequestWire{
		Expression:   "to.be.disabled",
		ExpectedText: []protocol.ExpectedTextValue{},
	})
}

// ToBeEditable asserts that the element is editable.
func (a *LocatorAssertions) ToBeEditable(ctx context.Context) error {
	return a.sendExpect(ctx, frameExpectRequestWire{
		Expression:   "to.be.editable",
		ExpectedText: []protocol.ExpectedTextValue{},
	})
}

// ToBeChecked asserts that the element is checked.
// The expectedValue must include { checked: true } in serialized form so the
// Playwright server can extract the boolean via `const { checked } = options.expectedValue`.
func (a *LocatorAssertions) ToBeChecked(ctx context.Context) error {
	return a.sendExpect(ctx, frameExpectRequestWire{
		Expression:    "to.be.checked",
		ExpectedText:  []protocol.ExpectedTextValue{},
		ExpectedValue: checkedExpectedValue(true),
	})
}

// ToHaveText asserts that the element's text matches the expected string.
// Uses exact matching with whitespace normalization.
func (a *LocatorAssertions) ToHaveText(ctx context.Context, expected string) error {
	matchSubstring := false
	normalizeWhiteSpace := true
	return a.sendExpect(ctx, frameExpectRequestWire{
		Expression: "to.have.text",
		ExpectedText: []protocol.ExpectedTextValue{
			{
				String:              &expected,
				MatchSubstring:      &matchSubstring,
				NormalizeWhiteSpace: &normalizeWhiteSpace,
			},
		},
	})
}

// ToHaveValue asserts that an input element has the expected value.
func (a *LocatorAssertions) ToHaveValue(ctx context.Context, expected string) error {
	return a.sendExpect(ctx, frameExpectRequestWire{
		Expression: "to.have.value",
		ExpectedText: []protocol.ExpectedTextValue{
			{
				String: &expected,
			},
		},
	})
}

// ToHaveAttribute asserts that the element has an attribute with the given name and value.
func (a *LocatorAssertions) ToHaveAttribute(ctx context.Context, name, value string) error {
	return a.sendExpect(ctx, frameExpectRequestWire{
		Expression:    "to.have.attribute.value",
		ExpressionArg: name,
		ExpectedText: []protocol.ExpectedTextValue{
			{
				String: &value,
			},
		},
	})
}

// ToHaveCount asserts that the number of elements matching the locator equals count.
func (a *LocatorAssertions) ToHaveCount(ctx context.Context, count int) error {
	expectedNumber := float64(count)
	return a.sendExpect(ctx, frameExpectRequestWire{
		Expression:     "to.have.count",
		ExpectedNumber: &expectedNumber,
		ExpectedText:   []protocol.ExpectedTextValue{},
	})
}

// ToBeAttached asserts that the element is attached to the DOM.
func (a *LocatorAssertions) ToBeAttached(ctx context.Context) error {
	return a.sendExpect(ctx, frameExpectRequestWire{
		Expression:   "to.be.attached",
		ExpectedText: []protocol.ExpectedTextValue{},
	})
}

// ToBeInViewport asserts that the element intersects the visible viewport.
func (a *LocatorAssertions) ToBeInViewport(ctx context.Context) error {
	return a.sendExpect(ctx, frameExpectRequestWire{
		Expression:   "to.be.in.viewport",
		ExpectedText: []protocol.ExpectedTextValue{},
	})
}

// ToHaveAccessibleName asserts that the element has the given ARIA accessible name.
func (a *LocatorAssertions) ToHaveAccessibleName(ctx context.Context, expected string) error {
	return a.sendExpect(ctx, frameExpectRequestWire{
		Expression:   "to.have.accessible.name",
		ExpectedText: []protocol.ExpectedTextValue{{String: &expected}},
	})
}

// ToHaveAccessibleDescription asserts that the element has the given ARIA accessible description.
func (a *LocatorAssertions) ToHaveAccessibleDescription(ctx context.Context, expected string) error {
	return a.sendExpect(ctx, frameExpectRequestWire{
		Expression:   "to.have.accessible.description",
		ExpectedText: []protocol.ExpectedTextValue{{String: &expected}},
	})
}

// ToHaveRole asserts that the element has the given ARIA role.
func (a *LocatorAssertions) ToHaveRole(ctx context.Context, role string) error {
	return a.sendExpect(ctx, frameExpectRequestWire{
		Expression:   "to.have.role",
		ExpectedText: []protocol.ExpectedTextValue{{String: &role}},
	})
}

// ToBeHidden asserts that the element is not visible.
func (a *LocatorAssertions) ToBeHidden(ctx context.Context) error {
	return a.sendExpect(ctx, frameExpectRequestWire{
		Expression:   "to.be.hidden",
		ExpectedText: []protocol.ExpectedTextValue{},
	})
}

// ToBeFocused asserts that the element is focused.
func (a *LocatorAssertions) ToBeFocused(ctx context.Context) error {
	return a.sendExpect(ctx, frameExpectRequestWire{
		Expression:   "to.be.focused",
		ExpectedText: []protocol.ExpectedTextValue{},
	})
}

// ToBeEmpty asserts that the element (input/textarea) is empty or has no child nodes.
func (a *LocatorAssertions) ToBeEmpty(ctx context.Context) error {
	return a.sendExpect(ctx, frameExpectRequestWire{
		Expression:   "to.be.empty",
		ExpectedText: []protocol.ExpectedTextValue{},
	})
}

// ToContainText asserts that the element contains the expected text (substring match).
// Uses "to.contain.text.array" which allows the element count to differ from expected count.
func (a *LocatorAssertions) ToContainText(ctx context.Context, expected string) error {
	matchSubstring := true
	normalizeWhiteSpace := true
	return a.sendExpect(ctx, frameExpectRequestWire{
		Expression: "to.contain.text.array",
		ExpectedText: []protocol.ExpectedTextValue{
			{
				String:              &expected,
				MatchSubstring:      &matchSubstring,
				NormalizeWhiteSpace: &normalizeWhiteSpace,
			},
		},
	})
}

// ToHaveClass asserts that the element has the expected CSS class.
func (a *LocatorAssertions) ToHaveClass(ctx context.Context, expected string) error {
	return a.sendExpect(ctx, frameExpectRequestWire{
		Expression: "to.have.class",
		ExpectedText: []protocol.ExpectedTextValue{
			{String: &expected},
		},
	})
}

// ToHaveID asserts that the element has the expected id attribute.
func (a *LocatorAssertions) ToHaveID(ctx context.Context, expected string) error {
	return a.sendExpect(ctx, frameExpectRequestWire{
		Expression: "to.have.id",
		ExpectedText: []protocol.ExpectedTextValue{
			{String: &expected},
		},
	})
}

// ToHaveCSS asserts that the element has a computed CSS property with the expected value.
func (a *LocatorAssertions) ToHaveCSS(ctx context.Context, name, value string) error {
	return a.sendExpect(ctx, frameExpectRequestWire{
		Expression:    "to.have.css",
		ExpressionArg: name,
		ExpectedText: []protocol.ExpectedTextValue{
			{String: &value},
		},
	})
}

// ToHaveJSProperty asserts that the element has a JavaScript property with the expected value.
// The server expects a SerializedArgument envelope {"value": <SerializedValue>, "handles": []},
// then deserializes it and compares with the element property via deepEquals.
func (a *LocatorAssertions) ToHaveJSProperty(ctx context.Context, name string, value any) error {
	envelope, err := json.Marshal(map[string]any{
		"value":   serializeValue(value),
		"handles": []any{},
	})
	if err != nil {
		return fmt.Errorf("ToHaveJSProperty: marshal failed: %w", err)
	}
	return a.sendExpect(ctx, frameExpectRequestWire{
		Expression:    "to.have.property",
		ExpressionArg: name,
		ExpectedValue: envelope,
		ExpectedText:  []protocol.ExpectedTextValue{},
	})
}

// ToHaveValues asserts that a multi-select element has all the specified selected values.
func (a *LocatorAssertions) ToHaveValues(ctx context.Context, values []string) error {
	expectedText := make([]protocol.ExpectedTextValue, len(values))
	for i, v := range values {
		v := v
		expectedText[i] = protocol.ExpectedTextValue{String: &v}
	}
	return a.sendExpect(ctx, frameExpectRequestWire{
		Expression:   "to.have.values",
		ExpectedText: expectedText,
	})
}

// extractRegexpInfo splits a Go *regexp.Regexp into a plain source pattern and JS-compatible
// flags string. Go uses inline flags like (?i) and negations like (?-i); we track the active
// set of i, m, s flags and strip all prefix directives from the source.
func extractRegexpInfo(re *regexp.Regexp) (source, flags string) {
	if re == nil {
		return "", ""
	}
	s := re.String()
	activeFlags := make(map[rune]bool)
	// Only strip pure flag-directive prefixes like (?i), (?-i), (?im) — NOT (?:...) groups.
	for {
		loc := flagPrefixRegexp.FindStringSubmatchIndex(s)
		if loc == nil {
			break
		}
		inner := s[loc[2]:loc[3]]
		negated := false
		for _, c := range inner {
			if c == '-' {
				negated = true
				continue
			}
			switch c {
			case 'i', 'm', 's':
				if negated {
					delete(activeFlags, c)
				} else {
					activeFlags[c] = true
				}
			}
		}
		s = s[loc[1]:]
	}
	var b strings.Builder
	for _, c := range []rune{'i', 'm', 's'} {
		if activeFlags[c] {
			b.WriteRune(c)
		}
	}
	return s, b.String()
}

// ToContainTextRegex asserts that the element contains text matching the given regular expression.
// Ref: TestLocatorAssertions.java#containsTextWRegexPass
func (a *LocatorAssertions) ToContainTextRegex(ctx context.Context, re *regexp.Regexp) error {
	src, flags := extractRegexpInfo(re)
	matchSubstring := true
	return a.sendExpect(ctx, frameExpectRequestWire{
		Expression: "to.contain.text.array",
		ExpectedText: []protocol.ExpectedTextValue{
			{RegexSource: &src, RegexFlags: &flags, MatchSubstring: &matchSubstring},
		},
	})
}

// ToHaveTextRegex asserts that the element's text matches the given regular expression.
// Ref: TestLocatorAssertions.java#hasTextWRegexPass
func (a *LocatorAssertions) ToHaveTextRegex(ctx context.Context, re *regexp.Regexp) error {
	src, flags := extractRegexpInfo(re)
	return a.sendExpect(ctx, frameExpectRequestWire{
		Expression: "to.have.text",
		ExpectedText: []protocol.ExpectedTextValue{
			{RegexSource: &src, RegexFlags: &flags},
		},
	})
}

// ToHaveTextArray asserts that each element in the matched set has a text matching the
// corresponding entry in the expected slice (ordered, substring match).
// Ref: TestLocatorAssertions.java#hasTextWTextArrayPass
func (a *LocatorAssertions) ToHaveTextArray(ctx context.Context, expected []string) error {
	matchSubstring := true
	normalizeWhiteSpace := true
	expectedText := make([]protocol.ExpectedTextValue, len(expected))
	for i, s := range expected {
		s := s
		expectedText[i] = protocol.ExpectedTextValue{
			String:              &s,
			MatchSubstring:      &matchSubstring,
			NormalizeWhiteSpace: &normalizeWhiteSpace,
		}
	}
	return a.sendExpect(ctx, frameExpectRequestWire{
		Expression:   "to.have.text.array",
		ExpectedText: expectedText,
	})
}

// ToHaveAttributeRegex asserts that the element has an attribute whose value matches the regex.
// Ref: TestLocatorAssertions.java#hasAttributeRegExpPass
func (a *LocatorAssertions) ToHaveAttributeRegex(ctx context.Context, name string, re *regexp.Regexp) error {
	src, flags := extractRegexpInfo(re)
	return a.sendExpect(ctx, frameExpectRequestWire{
		Expression:    "to.have.attribute.value",
		ExpressionArg: name,
		ExpectedText: []protocol.ExpectedTextValue{
			{RegexSource: &src, RegexFlags: &flags},
		},
	})
}

// ToHaveAttributeIgnoreCase asserts that the element has an attribute with the given value,
// ignoring case differences.
// Ref: TestLocatorAssertions.java#hasAttributeTextIgnoreCase
func (a *LocatorAssertions) ToHaveAttributeIgnoreCase(ctx context.Context, name, value string) error {
	ignoreCase := true
	return a.sendExpect(ctx, frameExpectRequestWire{
		Expression:    "to.have.attribute.value",
		ExpressionArg: name,
		ExpectedText: []protocol.ExpectedTextValue{
			{String: &value, IgnoreCase: &ignoreCase},
		},
	})
}

// ToHaveClassRegex asserts that the element's class attribute matches the given regular expression.
// Ref: TestLocatorAssertions.java#hasClassRegExpPass
func (a *LocatorAssertions) ToHaveClassRegex(ctx context.Context, re *regexp.Regexp) error {
	src, flags := extractRegexpInfo(re)
	return a.sendExpect(ctx, frameExpectRequestWire{
		Expression: "to.have.class",
		ExpectedText: []protocol.ExpectedTextValue{
			{RegexSource: &src, RegexFlags: &flags},
		},
	})
}

// ToHaveCSSRegex asserts that the element's computed CSS property value matches the regex.
// Ref: TestLocatorAssertions.java#hasCSSRegExPass
func (a *LocatorAssertions) ToHaveCSSRegex(ctx context.Context, name string, re *regexp.Regexp) error {
	src, flags := extractRegexpInfo(re)
	return a.sendExpect(ctx, frameExpectRequestWire{
		Expression:    "to.have.css",
		ExpressionArg: name,
		ExpectedText: []protocol.ExpectedTextValue{
			{RegexSource: &src, RegexFlags: &flags},
		},
	})
}

// ToHaveValueRegex asserts that an input element's value matches the given regular expression.
// Ref: TestLocatorAssertions.java#hasValueRegExpPass
func (a *LocatorAssertions) ToHaveValueRegex(ctx context.Context, re *regexp.Regexp) error {
	src, flags := extractRegexpInfo(re)
	return a.sendExpect(ctx, frameExpectRequestWire{
		Expression: "to.have.value",
		ExpectedText: []protocol.ExpectedTextValue{
			{RegexSource: &src, RegexFlags: &flags},
		},
	})
}

// ToBeCheckedFalse asserts that the element (checkbox/radio) is NOT checked.
// This is equivalent to expect(locator).not().toBeChecked() but uses the checked:false option.
// Ref: TestLocatorAssertions.java#isCheckedFalsePass
func (a *LocatorAssertions) ToBeCheckedFalse(ctx context.Context) error {
	return a.sendExpect(ctx, frameExpectRequestWire{
		Expression:    "to.be.checked",
		ExpectedText:  []protocol.ExpectedTextValue{},
		ExpectedValue: checkedExpectedValue(false),
	})
}

// ToContainClass asserts that the element's class list contains the given class name.
// Unlike ToHaveClass which matches the full class string, ToContainClass checks for
// inclusion of a specific class within the class attribute.
// Ref: TestLocatorAssertions.java#containsClassPass
func (a *LocatorAssertions) ToContainClass(ctx context.Context, class string) error {
	return a.sendExpect(ctx, frameExpectRequestWire{
		Expression: "to.contain.class",
		ExpectedText: []protocol.ExpectedTextValue{
			{String: &class},
		},
	})
}

// ToContainTextIgnoreCase asserts that the element contains the expected text, ignoring case.
func (a *LocatorAssertions) ToContainTextIgnoreCase(ctx context.Context, expected string) error {
	matchSubstring := true
	normalizeWhiteSpace := true
	ignoreCase := true
	return a.sendExpect(ctx, frameExpectRequestWire{
		Expression: "to.contain.text.array",
		ExpectedText: []protocol.ExpectedTextValue{
			{
				String:              &expected,
				MatchSubstring:      &matchSubstring,
				NormalizeWhiteSpace: &normalizeWhiteSpace,
				IgnoreCase:          &ignoreCase,
			},
		},
	})
}

// ToContainTextArray asserts that elements contain each expected text (subset, ordered match).
func (a *LocatorAssertions) ToContainTextArray(ctx context.Context, expected []string) error {
	matchSubstring := true
	normalizeWhiteSpace := true
	expectedText := make([]protocol.ExpectedTextValue, len(expected))
	for i, s := range expected {
		s := s
		expectedText[i] = protocol.ExpectedTextValue{
			String:              &s,
			MatchSubstring:      &matchSubstring,
			NormalizeWhiteSpace: &normalizeWhiteSpace,
		}
	}
	return a.sendExpect(ctx, frameExpectRequestWire{
		Expression:   "to.contain.text.array",
		ExpectedText: expectedText,
	})
}

// ToContainTextArrayIgnoreCase asserts elements contain each expected text, case-insensitively.
func (a *LocatorAssertions) ToContainTextArrayIgnoreCase(ctx context.Context, expected []string) error {
	matchSubstring := true
	normalizeWhiteSpace := true
	ignoreCase := true
	expectedText := make([]protocol.ExpectedTextValue, len(expected))
	for i, s := range expected {
		s := s
		expectedText[i] = protocol.ExpectedTextValue{
			String:              &s,
			MatchSubstring:      &matchSubstring,
			NormalizeWhiteSpace: &normalizeWhiteSpace,
			IgnoreCase:          &ignoreCase,
		}
	}
	return a.sendExpect(ctx, frameExpectRequestWire{
		Expression:   "to.contain.text.array",
		ExpectedText: expectedText,
	})
}

// ToHaveTextIgnoreCase asserts that the element's text matches the expected string, ignoring case.
func (a *LocatorAssertions) ToHaveTextIgnoreCase(ctx context.Context, expected string) error {
	matchSubstring := true
	normalizeWhiteSpace := true
	ignoreCase := true
	return a.sendExpect(ctx, frameExpectRequestWire{
		Expression: "to.have.text",
		ExpectedText: []protocol.ExpectedTextValue{
			{
				String:              &expected,
				MatchSubstring:      &matchSubstring,
				NormalizeWhiteSpace: &normalizeWhiteSpace,
				IgnoreCase:          &ignoreCase,
			},
		},
	})
}

// ToHaveTextUseInnerText asserts the element's innerText matches the expected string.
func (a *LocatorAssertions) ToHaveTextUseInnerText(ctx context.Context, expected string) error {
	useInnerText := true
	return a.sendExpect(ctx, frameExpectRequestWire{
		Expression: "to.have.text",
		ExpectedText: []protocol.ExpectedTextValue{
			{String: &expected},
		},
		UseInnerText: &useInnerText,
	})
}

// ToHaveTextArrayIgnoreCase asserts each element has text matching the array entry, case-insensitively.
func (a *LocatorAssertions) ToHaveTextArrayIgnoreCase(ctx context.Context, expected []string) error {
	ignoreCase := true
	normalizeWhiteSpace := true
	expectedText := make([]protocol.ExpectedTextValue, len(expected))
	for i, s := range expected {
		s := s
		expectedText[i] = protocol.ExpectedTextValue{
			String:              &s,
			IgnoreCase:          &ignoreCase,
			NormalizeWhiteSpace: &normalizeWhiteSpace,
		}
	}
	return a.sendExpect(ctx, frameExpectRequestWire{
		Expression:   "to.have.text.array",
		ExpectedText: expectedText,
	})
}

// ToHaveTextRegexIgnoreCase asserts element text matches the regex source with an explicit ignoreCase override.
func (a *LocatorAssertions) ToHaveTextRegexIgnoreCase(ctx context.Context, re *regexp.Regexp, ignoreCase bool) error {
	src, _ := extractRegexpInfo(re)
	ic := ignoreCase
	return a.sendExpect(ctx, frameExpectRequestWire{
		Expression: "to.have.text",
		ExpectedText: []protocol.ExpectedTextValue{
			{RegexSource: &src, IgnoreCase: &ic},
		},
	})
}

// ToHaveTextRegexArray asserts each element's text matches the corresponding regex in the slice.
func (a *LocatorAssertions) ToHaveTextRegexArray(ctx context.Context, res []*regexp.Regexp) error {
	expectedText := make([]protocol.ExpectedTextValue, len(res))
	for i, re := range res {
		src, flags := extractRegexpInfo(re)
		s, f := src, flags
		expectedText[i] = protocol.ExpectedTextValue{
			RegexSource: &s,
			RegexFlags:  &f,
		}
	}
	return a.sendExpect(ctx, frameExpectRequestWire{
		Expression:   "to.have.text.array",
		ExpectedText: expectedText,
	})
}

// ToHaveClassArray asserts that each element has the corresponding class string (full match).
func (a *LocatorAssertions) ToHaveClassArray(ctx context.Context, classes []string) error {
	expectedText := make([]protocol.ExpectedTextValue, len(classes))
	for i, s := range classes {
		s := s
		expectedText[i] = protocol.ExpectedTextValue{String: &s}
	}
	return a.sendExpect(ctx, frameExpectRequestWire{
		Expression:   "to.have.class.array",
		ExpectedText: expectedText,
	})
}

// ToHaveClassRegexArray asserts each element's class attribute matches the corresponding regex.
func (a *LocatorAssertions) ToHaveClassRegexArray(ctx context.Context, res []*regexp.Regexp) error {
	expectedText := make([]protocol.ExpectedTextValue, len(res))
	for i, re := range res {
		src, flags := extractRegexpInfo(re)
		s, f := src, flags
		expectedText[i] = protocol.ExpectedTextValue{
			RegexSource: &s,
			RegexFlags:  &f,
		}
	}
	return a.sendExpect(ctx, frameExpectRequestWire{
		Expression:   "to.have.class.array",
		ExpectedText: expectedText,
	})
}

// ToHaveValuesRegex asserts that a multi-select element's selected values match the given regexes.
func (a *LocatorAssertions) ToHaveValuesRegex(ctx context.Context, res []*regexp.Regexp) error {
	expectedText := make([]protocol.ExpectedTextValue, len(res))
	for i, re := range res {
		src, flags := extractRegexpInfo(re)
		s, f := src, flags
		expectedText[i] = protocol.ExpectedTextValue{
			RegexSource: &s,
			RegexFlags:  &f,
		}
	}
	return a.sendExpect(ctx, frameExpectRequestWire{
		Expression:   "to.have.values",
		ExpectedText: expectedText,
	})
}

// ToBeDetached asserts that the element is NOT attached to the DOM.
func (a *LocatorAssertions) ToBeDetached(ctx context.Context) error {
	return a.sendExpect(ctx, frameExpectRequestWire{
		Expression:   "to.be.detached",
		ExpectedText: []protocol.ExpectedTextValue{},
	})
}

// ToBeReadOnly asserts that the element is not editable (is readonly or disabled).
func (a *LocatorAssertions) ToBeReadOnly(ctx context.Context) error {
	return a.sendExpect(ctx, frameExpectRequestWire{
		Expression:   "to.be.readonly",
		ExpectedText: []protocol.ExpectedTextValue{},
	})
}

// ToHaveAccessibleNameRegex asserts the element has an accessible name matching the regex.
func (a *LocatorAssertions) ToHaveAccessibleNameRegex(ctx context.Context, re *regexp.Regexp) error {
	src, flags := extractRegexpInfo(re)
	return a.sendExpect(ctx, frameExpectRequestWire{
		Expression:   "to.have.accessible.name",
		ExpectedText: []protocol.ExpectedTextValue{{RegexSource: &src, RegexFlags: &flags}},
	})
}

// ToHaveAccessibleNameIgnoreCase asserts the element has the given accessible name, ignoring case.
func (a *LocatorAssertions) ToHaveAccessibleNameIgnoreCase(ctx context.Context, text string) error {
	ignoreCase := true
	return a.sendExpect(ctx, frameExpectRequestWire{
		Expression:   "to.have.accessible.name",
		ExpectedText: []protocol.ExpectedTextValue{{String: &text, IgnoreCase: &ignoreCase}},
	})
}

// ToHaveAccessibleDescriptionRegex asserts the element has an accessible description matching the regex.
func (a *LocatorAssertions) ToHaveAccessibleDescriptionRegex(ctx context.Context, re *regexp.Regexp) error {
	src, flags := extractRegexpInfo(re)
	return a.sendExpect(ctx, frameExpectRequestWire{
		Expression:   "to.have.accessible.description",
		ExpectedText: []protocol.ExpectedTextValue{{RegexSource: &src, RegexFlags: &flags}},
	})
}

// ToHaveAccessibleDescriptionIgnoreCase asserts the element has the given accessible description, ignoring case.
func (a *LocatorAssertions) ToHaveAccessibleDescriptionIgnoreCase(ctx context.Context, text string) error {
	ignoreCase := true
	return a.sendExpect(ctx, frameExpectRequestWire{
		Expression:   "to.have.accessible.description",
		ExpectedText: []protocol.ExpectedTextValue{{String: &text, IgnoreCase: &ignoreCase}},
	})
}

// ToHaveAccessibleErrorMessage asserts the element has the given ARIA error message.
func (a *LocatorAssertions) ToHaveAccessibleErrorMessage(ctx context.Context, text string) error {
	return a.sendExpect(ctx, frameExpectRequestWire{
		Expression:   "to.have.accessible.error.message",
		ExpectedText: []protocol.ExpectedTextValue{{String: &text}},
	})
}

// ToHaveAccessibleErrorMessageRegex asserts the element's ARIA error message matches the regex.
func (a *LocatorAssertions) ToHaveAccessibleErrorMessageRegex(ctx context.Context, re *regexp.Regexp) error {
	src, flags := extractRegexpInfo(re)
	return a.sendExpect(ctx, frameExpectRequestWire{
		Expression:   "to.have.accessible.error.message",
		ExpectedText: []protocol.ExpectedTextValue{{RegexSource: &src, RegexFlags: &flags}},
	})
}

// ToHaveAccessibleErrorMessageIgnoreCase asserts the element's ARIA error message matches the text, ignoring case.
func (a *LocatorAssertions) ToHaveAccessibleErrorMessageIgnoreCase(ctx context.Context, text string) error {
	ignoreCase := true
	return a.sendExpect(ctx, frameExpectRequestWire{
		Expression:   "to.have.accessible.error.message",
		ExpectedText: []protocol.ExpectedTextValue{{String: &text, IgnoreCase: &ignoreCase}},
	})
}

// indeterminateCheckedValue builds the expectedValue for a to.be.checked assertion with indeterminate=true.
func indeterminateCheckedValue() json.RawMessage {
	msg, err := json.Marshal(map[string]any{
		"value": map[string]any{
			"o": []any{
				map[string]any{
					"k": "indeterminate",
					"v": map[string]any{"b": true},
				},
			},
		},
		"handles": []any{},
	})
	if err != nil {
		return nil
	}
	return msg
}

// ToBeCheckedIndeterminate asserts that the element (checkbox) is in indeterminate state.
func (a *LocatorAssertions) ToBeCheckedIndeterminate(ctx context.Context) error {
	return a.sendExpect(ctx, frameExpectRequestWire{
		Expression:    "to.be.checked",
		ExpectedText:  []protocol.ExpectedTextValue{},
		ExpectedValue: indeterminateCheckedValue(),
	})
}

// ToBeCheckedWithIndeterminateAndChecked sends both indeterminate=true and checked options together.
// This combination is invalid per Playwright's protocol; the server returns an error:
// "Can't assert indeterminate and checked at the same time".
// Used only in tests that verify proper error handling.
func (a *LocatorAssertions) ToBeCheckedWithIndeterminateAndChecked(ctx context.Context, checked bool) error {
	msg, _ := json.Marshal(map[string]any{
		"value": map[string]any{
			"o": []any{
				map[string]any{"k": "indeterminate", "v": map[string]any{"b": true}},
				map[string]any{"k": "checked", "v": map[string]any{"b": checked}},
			},
		},
		"handles": []any{},
	})
	return a.sendExpect(ctx, frameExpectRequestWire{
		Expression:    "to.be.checked",
		ExpectedText:  []protocol.ExpectedTextValue{},
		ExpectedValue: msg,
	})
}

// ToContainClassArray asserts that each element in the matched set contains the corresponding class.
func (a *LocatorAssertions) ToContainClassArray(ctx context.Context, classes []string) error {
	expectedText := make([]protocol.ExpectedTextValue, len(classes))
	for i, s := range classes {
		s := s
		expectedText[i] = protocol.ExpectedTextValue{String: &s}
	}
	return a.sendExpect(ctx, frameExpectRequestWire{
		Expression:   "to.contain.class.array",
		ExpectedText: expectedText,
	})
}

// ToMatchAriaSnapshot asserts that the element's ARIA accessibility tree matches the
// expected YAML snapshot string. Delegates to the Playwright server for structured YAML
// subset matching (requires Playwright v1.49+).
func (a *LocatorAssertions) ToMatchAriaSnapshot(ctx context.Context, expected string) error {
	return a.sendExpect(ctx, frameExpectRequestWire{
		Expression:   "to.match.aria.snapshot",
		ExpectedText: []protocol.ExpectedTextValue{{String: &expected}},
	})
}

// PageAssertions provides web-first assertion methods for a Page.
// All assertions retry automatically until the condition is met or the context deadline is reached.
type PageAssertions struct {
	page  *Page
	isNot bool
}

// ExpectPage creates a new PageAssertions for the given Page.
func ExpectPage(page *Page) *PageAssertions {
	return &PageAssertions{page: page}
}

// Not returns a negated version of the assertion.
func (a *PageAssertions) Not() *PageAssertions {
	return &PageAssertions{page: a.page, isNot: true}
}

// pollPageAssertion retries check every 100 ms until it returns nil or the context is done.
// If the context has no deadline, a default assertion timeout is applied.
func pollPageAssertion(ctx context.Context, check func() error) error {
	if _, hasDeadline := ctx.Deadline(); !hasDeadline {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Duration(defaultActionTimeoutMs*float64(time.Millisecond)))
		defer cancel()
	}
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	var lastErr error
	for {
		lastErr = check()
		if lastErr == nil {
			return nil
		}
		select {
		case <-ctx.Done():
			return lastErr
		case <-ticker.C:
		}
	}
}

// HasURL asserts that the page URL matches the given string (exact match).
func (a *PageAssertions) HasURL(ctx context.Context, expected string) error {
	return pollPageAssertion(ctx, func() error {
		url, err := a.page.Evaluate(ctx, "() => window.location.href")
		if err != nil {
			return fmt.Errorf("pageAssert.hasURL: evaluate failed: %w", err)
		}
		actual, _ := url.(string)
		if a.isNot {
			if actual == expected {
				return fmt.Errorf("expected URL NOT to be %q, but it was", expected)
			}
		} else {
			if actual != expected {
				return fmt.Errorf("expected URL %q, got %q", expected, actual)
			}
		}
		return nil
	})
}

// HasURLContains asserts that the page URL contains the expected substring.
func (a *PageAssertions) HasURLContains(ctx context.Context, expected string) error {
	return pollPageAssertion(ctx, func() error {
		url, err := a.page.Evaluate(ctx, "() => window.location.href")
		if err != nil {
			return fmt.Errorf("pageAssert.hasURLContains: evaluate failed: %w", err)
		}
		actual, _ := url.(string)
		if a.isNot {
			if strings.Contains(actual, expected) {
				return fmt.Errorf("expected URL NOT to contain %q, but got %q", expected, actual)
			}
		} else {
			if !strings.Contains(actual, expected) {
				return fmt.Errorf("expected URL to contain %q, got %q", expected, actual)
			}
		}
		return nil
	})
}

// HasURLRegex asserts that the page URL matches the given regular expression.
func (a *PageAssertions) HasURLRegex(ctx context.Context, pattern string) error {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return fmt.Errorf("pageAssert.hasURLRegex: invalid pattern %q: %w", pattern, err)
	}
	return pollPageAssertion(ctx, func() error {
		url, err := a.page.Evaluate(ctx, "() => window.location.href")
		if err != nil {
			return fmt.Errorf("pageAssert.hasURLRegex: evaluate failed: %w", err)
		}
		actual, _ := url.(string)
		if a.isNot {
			if re.MatchString(actual) {
				return fmt.Errorf("expected URL NOT to match %q, but got %q", pattern, actual)
			}
		} else {
			if !re.MatchString(actual) {
				return fmt.Errorf("expected URL to match %q, got %q", pattern, actual)
			}
		}
		return nil
	})
}

// HasTitle asserts that the page title matches the given string (exact match).
func (a *PageAssertions) HasTitle(ctx context.Context, expected string) error {
	return pollPageAssertion(ctx, func() error {
		actual, err := a.page.Title(ctx)
		if err != nil {
			return fmt.Errorf("pageAssert.hasTitle: Title() failed: %w", err)
		}
		if a.isNot {
			if actual == expected {
				return fmt.Errorf("expected title NOT to be %q, but it was", expected)
			}
		} else {
			if actual != expected {
				return fmt.Errorf("expected title %q, got %q", expected, actual)
			}
		}
		return nil
	})
}

// HasTitleContains asserts that the page title contains the expected substring.
func (a *PageAssertions) HasTitleContains(ctx context.Context, expected string) error {
	return pollPageAssertion(ctx, func() error {
		actual, err := a.page.Title(ctx)
		if err != nil {
			return fmt.Errorf("pageAssert.hasTitleContains: Title() failed: %w", err)
		}
		if a.isNot {
			if strings.Contains(actual, expected) {
				return fmt.Errorf("expected title NOT to contain %q, but got %q", expected, actual)
			}
		} else {
			if !strings.Contains(actual, expected) {
				return fmt.Errorf("expected title to contain %q, got %q", expected, actual)
			}
		}
		return nil
	})
}

// HasTitleRegex asserts that the page title matches the given regular expression.
func (a *PageAssertions) HasTitleRegex(ctx context.Context, pattern string) error {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return fmt.Errorf("pageAssert.hasTitleRegex: invalid pattern %q: %w", pattern, err)
	}
	return pollPageAssertion(ctx, func() error {
		actual, err := a.page.Title(ctx)
		if err != nil {
			return fmt.Errorf("pageAssert.hasTitleRegex: Title() failed: %w", err)
		}
		if a.isNot {
			if re.MatchString(actual) {
				return fmt.Errorf("expected title NOT to match %q, but got %q", pattern, actual)
			}
		} else {
			if !re.MatchString(actual) {
				return fmt.Errorf("expected title to match %q, got %q", pattern, actual)
			}
		}
		return nil
	})
}
