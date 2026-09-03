//go:build e2e

package e2e

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/jcfrperu/goaneco-playwright"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// shortCtx returns a context with a 2-second deadline, useful for fail tests.
func shortCtx(t *testing.T) context.Context {
	t.Helper()
	ctx, cancel := context.WithTimeout(testCtx(t), 2*time.Second)
	t.Cleanup(cancel)
	return ctx
}

// TestContainsTextRegexPass verifies containsText with a plain regex pattern.
// Ref: TestLocatorAssertions.java#containsTextWRegexPass
func TestContainsTextRegexPass(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div id=node>Text   content</div>`)
	must.NoError(err)

	loc := page.Locator("#node")
	must.NoError(playwright.Expect(loc).ToContainTextRegex(ctx, regexp.MustCompile(`ex`)))
	// Should not normalize whitespace for regex.
	must.NoError(playwright.Expect(loc).ToContainTextRegex(ctx, regexp.MustCompile(`ext   cont`)))
}

// TestContainsTextRegexCaseInsensitive verifies case-insensitive regex containsText.
// Ref: TestLocatorAssertions.java#containsTextWRegexCaseInsensitivePass
func TestContainsTextRegexCaseInsensitive(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div id=node>Text   content</div>`)
	must.NoError(err)

	// (?i) flag makes the regex case-insensitive.
	must.NoError(playwright.Expect(page.Locator("#node")).ToContainTextRegex(ctx, regexp.MustCompile(`(?i)text`)))
}

// TestContainsTextRegexMultiline verifies multiline regex containsText.
// Ref: TestLocatorAssertions.java#containsTextWRegexMultilinePass
func TestContainsTextRegexMultiline(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, "<div id=node>Text \nContent</div>")
	must.NoError(err)

	// (?m) flag enables multiline matching (^ matches start of each line).
	must.NoError(playwright.Expect(page.Locator("#node")).ToContainTextRegex(ctx, regexp.MustCompile(`(?m)^Content`)))
}

// TestContainsTextRegexDotAll verifies dotall regex containsText (. matches newline).
// Ref: TestLocatorAssertions.java#containsTextWRegexDotAllPass
func TestContainsTextRegexDotAll(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, "<div id=node>foo\nbar</div>")
	must.NoError(err)

	// (?s) flag enables dot-all mode (. matches newline).
	must.NoError(playwright.Expect(page.Locator("#node")).ToContainTextRegex(ctx, regexp.MustCompile(`(?s)foo.bar`)))
}

// TestContainsTextRegexFail verifies that a non-matching regex fails containsText.
// Ref: TestLocatorAssertions.java#containsTextWRegexFail
func TestContainsTextRegexFail(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div id=node>Text   content</div>`)
	must.NoError(err)

	err = playwright.Expect(page.Locator("#node")).ToContainTextRegex(ctx, regexp.MustCompile(`ex2`))
	is.Error(err, "non-matching regex should fail containsText")
}

// TestContainsTextWithIgnoreCase verifies containsText with ignoreCase and falsy ignoreCase.
// Ref: TestLocatorAssertions.java#containsTextWTextPass
func TestContainsTextWithIgnoreCase(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div id=node>Text content</div>`)
	must.NoError(err)

	loc := page.Locator("#node")
	// Basic substring
	must.NoError(playwright.Expect(loc).ToContainText(ctx, "Text"))
	// IgnoreCase = true
	must.NoError(playwright.Expect(loc).ToContainTextIgnoreCase(ctx, "EXT"))
	// IgnoreCase = false (default): "TEXT" should NOT match "Text content" case-sensitively
	must.NoError(playwright.Expect(loc).Not().ToContainText(ctx, "TEXT"))
}

// TestContainsTextArray verifies containsText with a string array.
// Ref: TestLocatorAssertions.java#containsTextWTextArrayPass
func TestContainsTextArray(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div>Text \n1</div><div>Text2</div><div>Text3</div>`)
	must.NoError(err)

	loc := page.Locator("div")
	must.NoError(playwright.Expect(loc).ToContainTextArray(ctx, []string{"ext     1", "ext3"}))
	// ignoreCase variant
	must.NoError(playwright.Expect(loc).ToContainTextArrayIgnoreCase(ctx, []string{"EXT 1", "eXt3"}))
}

// TestHasTextRegexPass verifies hasText with various regex options.
// Ref: TestLocatorAssertions.java#hasTextWRegexPass
func TestHasTextRegexPass(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div id=node>Text   content</div>`)
	must.NoError(err)

	loc := page.Locator("#node")
	must.NoError(playwright.Expect(loc).ToHaveTextRegex(ctx, regexp.MustCompile(`Te.t`)))
	// Should not normalize whitespace for regex.
	must.NoError(playwright.Expect(loc).ToHaveTextRegex(ctx, regexp.MustCompile(`Text.+content`)))
	// ignoreCase override: regex without flag + ignoreCase=true → should match
	must.NoError(playwright.Expect(loc).ToHaveTextRegexIgnoreCase(ctx, regexp.MustCompile(`text   content`), true))
	// ignoreCase override: regex with (?i) flag + ignoreCase=false → should NOT match case-insensitively
	must.NoError(playwright.Expect(loc).Not().ToHaveTextRegexIgnoreCase(ctx, regexp.MustCompile(`(?i)text   content`), false))
}

// TestHasTextRegexFail verifies that a non-matching regex fails hasText.
// Ref: TestLocatorAssertions.java#hasTextWRegexFail
func TestHasTextRegexFail(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div id=node>Text   content</div>`)
	must.NoError(err)

	err = playwright.Expect(page.Locator("#node")).ToHaveTextRegex(ctx, regexp.MustCompile(`Text 2`))
	is.Error(err, "non-matching regex should fail hasText")
}

// TestHasTextWithIgnoreCase verifies hasText with ignoreCase option.
// Ref: TestLocatorAssertions.java#hasTextWTextPass
func TestHasTextWithIgnoreCase(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, "<div id=node><span></span>Text \ncontent&nbsp;    </div>")
	must.NoError(err)

	loc := page.Locator("#node")
	// Whitespace normalization
	must.NoError(playwright.Expect(loc).ToHaveText(ctx, "Text                        content"))
	// ignoreCase = true
	must.NoError(playwright.Expect(loc).ToHaveTextIgnoreCase(ctx, "text CONTENT"))
	// ignoreCase = false: "TEXT" should NOT match "Text content"
	must.NoError(playwright.Expect(loc).Not().ToHaveText(ctx, "TEXT"))
}

// TestHasTextFail verifies that a non-matching string fails hasText.
// Ref: TestLocatorAssertions.java#hasTextWTextFail
func TestHasTextFail(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div id=node>Text content</div>`)
	must.NoError(err)

	// "Text" (substring) should fail hasText when full text is "Text content"
	err = playwright.Expect(page.Locator("#node")).ToHaveText(ctx, "Text")
	is.Error(err, "substring 'Text' should fail hasText on 'Text content'")
}

// TestHasTextUseInnerText verifies hasText with useInnerText (ignores hidden spans).
// Ref: TestLocatorAssertions.java#hasTextWTextInnerTextPass
func TestHasTextUseInnerText(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div id=node>Text <span hidden>garbage</span> content</div>`)
	must.NoError(err)

	// useInnerText excludes hidden elements
	must.NoError(playwright.Expect(page.Locator("#node")).ToHaveTextUseInnerText(ctx, "Text content"))
}

// TestHasTextArrayPass verifies hasText with a string array (normalized whitespace).
// Ref: TestLocatorAssertions.java#hasTextWTextArrayPass
func TestHasTextArrayPass(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, "<div>Text    \n1</div><div>Text   2a</div>")
	must.NoError(err)

	loc := page.Locator("div")
	// Normalized whitespace
	must.NoError(playwright.Expect(loc).ToHaveTextArray(ctx, []string{"Text  1", "Text   2a"}))
	// ignoreCase
	must.NoError(playwright.Expect(loc).ToHaveTextArrayIgnoreCase(ctx, []string{"tEXT 1", "TExt 2A"}))
}

// TestHasTextArrayEmpty verifies hasText with an empty array passes on empty locator.
// Ref: TestLocatorAssertions.java#hasTextWTextArrayPassEmpty
func TestHasTextArrayEmpty(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div></div>`)
	must.NoError(err)

	// locator("p") matches nothing → empty array assertion should pass
	must.NoError(playwright.Expect(page.Locator("p")).ToHaveTextArray(ctx, []string{}))
}

// TestHasTextArrayNotEmpty verifies not().hasText with empty array fails when elements exist.
// Ref: TestLocatorAssertions.java#hasTextWTextArrayPassNotEmpty
func TestHasTextArrayNotEmpty(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div><p>Test</p></div>`)
	must.NoError(err)

	// "div" matches an element with text → not().hasText([]) should pass (element is not empty)
	must.NoError(playwright.Expect(page.Locator("div")).Not().ToHaveTextArray(ctx, []string{}))
}

// TestHasTextArrayPassOnEmptyLocator verifies not().hasText fails on empty locator.
// Ref: TestLocatorAssertions.java#hasTextWTextArrayPassOnEmpty
func TestHasTextArrayPassOnEmptyLocator(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div></div>`)
	must.NoError(err)

	// locator("p") matches nothing → not().hasText(["Test"]) should pass
	must.NoError(playwright.Expect(page.Locator("p")).Not().ToHaveTextArray(ctx, []string{"Test"}))
}

// TestHasTextArrayFailOnNotEmpty verifies not().hasText([]) fails when locator matches nothing.
// Ref: TestLocatorAssertions.java#hasTextWTextArrayFailOnNotEmpty
func TestHasTextArrayFailOnNotEmpty(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div></div>`)
	must.NoError(err)

	// locator("p") matches nothing → not().hasText([]) fails (both expected and actual are [])
	err = playwright.Expect(page.Locator("p")).Not().ToHaveTextArray(ctx, []string{})
	is.Error(err, "not().hasText([]) should fail when locator matches nothing (both are [])")
}

// TestHasTextArrayFail verifies hasText array fails when count or text differ.
// Ref: TestLocatorAssertions.java#hasTextWTextArrayFail
func TestHasTextArrayFail(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div>Text 1</div><div>Text 3</div>`)
	must.NoError(err)

	// Two elements, asking for three entries → should fail
	err = playwright.Expect(page.Locator("div")).ToHaveTextArray(ctx, []string{"Text 1", "Text 3", "Extra"})
	is.Error(err, "array count mismatch should fail hasText")
}

// TestHasTextRegexArrayPass verifies hasText with a regex array.
// Ref: TestLocatorAssertions.java#hasTextWRegExArrayPass
func TestHasTextRegexArrayPass(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, "<div>Text    \n1</div><div>Text   2a</div>")
	must.NoError(err)

	loc := page.Locator("div")
	must.NoError(playwright.Expect(loc).ToHaveTextRegexArray(ctx, []*regexp.Regexp{
		regexp.MustCompile("Text    \n1"),
		regexp.MustCompile(`Text   \d+a`),
	}))
}

// TestHasTextRegexArrayFail verifies hasText with a regex array fails on mismatch.
// Ref: TestLocatorAssertions.java#hasTextWRegExArrayFail
func TestHasTextRegexArrayFail(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div>Text 1</div><div>Text 3</div>`)
	must.NoError(err)

	err = playwright.Expect(page.Locator("div")).ToHaveTextRegexArray(ctx, []*regexp.Regexp{
		regexp.MustCompile(`Text 1`),
		regexp.MustCompile(`Text   \d`),
		regexp.MustCompile(`Extra`),
	})
	is.Error(err, "regex array count mismatch should fail hasText")
}

// TestHasAttributeIgnoreCase verifies hasAttribute with ignoreCase option.
// Ref: TestLocatorAssertions.java#hasAttributeTextIgnoreCase
func TestHasAttributeIgnoreCase(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div id=NoDe>Text content</div>`)
	must.NoError(err)

	loc := page.Locator("#NoDe")
	// ignoreCase=true: "node" matches "NoDe"
	must.NoError(playwright.Expect(loc).ToHaveAttributeIgnoreCase(ctx, "id", "node"))
	// case-sensitive (default): "node" does NOT match "NoDe"
	must.NoError(playwright.Expect(loc).Not().ToHaveAttribute(ctx, "id", "node"))
}

// TestHasAttributeFail verifies that hasAttribute fails on value mismatch.
// Ref: TestLocatorAssertions.java#hasAttributeTextFail
func TestHasAttributeFail(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div id=node>Text content</div>`)
	must.NoError(err)

	err = playwright.Expect(page.Locator("#node")).ToHaveAttribute(ctx, "id", "foo")
	is.Error(err, "wrong attribute value should fail hasAttribute")
}

// TestHasAttributeRegexPass verifies hasAttribute with a regex pattern.
// Ref: TestLocatorAssertions.java#hasAttributeRegExpPass
func TestHasAttributeRegexPass(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div id=node>Text content</div>`)
	must.NoError(err)

	must.NoError(playwright.Expect(page.Locator("#node")).ToHaveAttributeRegex(ctx, "id", regexp.MustCompile(`n..e`)))
}

// TestHasAttributeRegexFail verifies that hasAttribute regex fails on no match.
// Ref: TestLocatorAssertions.java#hasAttributeRegExpFail
func TestHasAttributeRegexFail(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div id=node>Text content</div>`)
	must.NoError(err)

	err = playwright.Expect(page.Locator("#node")).ToHaveAttributeRegex(ctx, "id", regexp.MustCompile(`.Nod..`))
	is.Error(err, "non-matching regex should fail hasAttribute")
}

// TestHasClassFail verifies hasClass fails when the class string doesn't match.
// Ref: TestLocatorAssertions.java#hasClassTextFail
func TestHasClassFail(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div class="bar baz"></div>`)
	must.NoError(err)

	err = playwright.Expect(page.Locator("div")).ToHaveClass(ctx, "foo bar baz")
	is.Error(err, "class mismatch should fail hasClass")
}

// TestHasClassRegexPass verifies hasClass with a regex pattern.
// Ref: TestLocatorAssertions.java#hasClassRegExpPass
func TestHasClassRegexPass(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div class="foo bar baz"></div>`)
	must.NoError(err)

	must.NoError(playwright.Expect(page.Locator("div")).ToHaveClassRegex(ctx, regexp.MustCompile(`foo.* baz`)))
}

// TestHasClassRegexFail verifies hasClass regex fails on no match.
// Ref: TestLocatorAssertions.java#hasClassRegExpFail
func TestHasClassRegexFail(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div class="bar baz"></div>`)
	must.NoError(err)

	err = playwright.Expect(page.Locator("div")).ToHaveClassRegex(ctx, regexp.MustCompile(`foo Z.*`))
	is.Error(err, "non-matching regex should fail hasClass")
}

// TestHasClassArrayPass verifies hasClass with a string array.
// Ref: TestLocatorAssertions.java#hasClassTextArrayPass
func TestHasClassArrayPass(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div class="foo"></div><div class="bar"></div><div class="baz"></div>`)
	must.NoError(err)

	must.NoError(playwright.Expect(page.Locator("div")).ToHaveClassArray(ctx, []string{"foo", "bar", "baz"}))
}

// TestHasClassArrayFail verifies hasClass array fails on mismatch.
// Ref: TestLocatorAssertions.java#hasClassTextArrayFail
func TestHasClassArrayFail(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div class="foo"></div><div class="bar"></div><div class="baz"></div>`)
	must.NoError(err)

	err = playwright.Expect(page.Locator("div")).ToHaveClassArray(ctx, []string{"foo", "bar", "missing"})
	is.Error(err, "class array mismatch should fail hasClass")
}

// TestHasClassRegexArrayPass verifies hasClass with a regex array.
// Ref: TestLocatorAssertions.java#hasClassRegExpArrayPass
func TestHasClassRegexArrayPass(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div class="foo"></div><div class="bar"></div><div class="baz"></div>`)
	must.NoError(err)

	must.NoError(playwright.Expect(page.Locator("div")).ToHaveClassRegexArray(ctx, []*regexp.Regexp{
		regexp.MustCompile(`fo.*`),
		regexp.MustCompile(`.ar`),
		regexp.MustCompile(`baz`),
	}))
}

// TestHasClassRegexArrayFail verifies hasClass regex array fails on count mismatch.
// Ref: TestLocatorAssertions.java#hasClassRegExpArrayFail
func TestHasClassRegexArrayFail(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div class="foo"></div><div class="bar"></div><div class="baz"></div>`)
	must.NoError(err)

	err = playwright.Expect(page.Locator("div")).ToHaveClassRegexArray(ctx, []*regexp.Regexp{
		regexp.MustCompile(`fo.*`),
		regexp.MustCompile(`.ar`),
		regexp.MustCompile(`baz`),
		regexp.MustCompile(`extra`),
	})
	is.Error(err, "regex array count mismatch should fail hasClass")
}

// TestHasCountFail verifies hasCount fails when the count is wrong.
// Ref: TestLocatorAssertions.java#hasCountFail
func TestHasCountFail(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<select><option>One</option><option>Two</option></select>`)
	must.NoError(err)

	err = playwright.Expect(page.Locator("option")).ToHaveCount(ctx, 1)
	is.Error(err, "wrong count should fail hasCount")
}

// TestHasCountZero verifies hasCount(0) and not().hasCount(1) on non-existent elements.
// Ref: TestLocatorAssertions.java#hasCountPassZero
func TestHasCountZero(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div></div>`)
	must.NoError(err)

	loc := page.Locator("span")
	must.NoError(playwright.Expect(loc).ToHaveCount(ctx, 0))
	must.NoError(playwright.Expect(loc).Not().ToHaveCount(ctx, 1))
}

// TestHasCSSPass verifies hasCSS with exact value.
// Ref: TestLocatorAssertions.java#hasCSSPass
func TestHasCSSPass(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div id=node style='color: rgb(255, 0, 0)'>Text content</div>`)
	must.NoError(err)

	must.NoError(playwright.Expect(page.Locator("#node")).ToHaveCSS(ctx, "color", "rgb(255, 0, 0)"))
}

// TestHasCSSFail verifies hasCSS fails on value mismatch.
// Ref: TestLocatorAssertions.java#hasCSSFail
func TestHasCSSFail(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div id=node style='color: rgb(255, 0, 0)'>Text content</div>`)
	must.NoError(err)

	err = playwright.Expect(page.Locator("#node")).ToHaveCSS(ctx, "color", "red")
	is.Error(err, "CSS value mismatch should fail hasCSS")
}

// TestHasCSSRegexPass verifies hasCSS with a regex pattern.
// Ref: TestLocatorAssertions.java#hasCSSRegExPass
func TestHasCSSRegexPass(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div id=node style='color: rgb(255, 0, 0)'>Text content</div>`)
	must.NoError(err)

	must.NoError(playwright.Expect(page.Locator("#node")).ToHaveCSSRegex(ctx, "color", regexp.MustCompile(`rgb.*`)))
}

// TestHasCSSRegexFail verifies hasCSS regex fails on no match.
// Ref: TestLocatorAssertions.java#hasCSSRegExFail
func TestHasCSSRegexFail(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div id=node style='color: rgb(255, 0, 0)'>Text content</div>`)
	must.NoError(err)

	err = playwright.Expect(page.Locator("#node")).ToHaveCSSRegex(ctx, "color", regexp.MustCompile(`red`))
	is.Error(err, "non-matching CSS regex should fail hasCSS")
}

// TestHasIdPass verifies hasId matches the id attribute.
// Ref: TestLocatorAssertions.java#hasIdPass
func TestHasIdPass(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div id=node>Text content</div>`)
	must.NoError(err)

	must.NoError(playwright.Expect(page.Locator("#node")).ToHaveID(ctx, "node"))
}

// TestHasIdFail verifies hasId fails on mismatch.
// Ref: TestLocatorAssertions.java#hasIdFail
func TestHasIdFail(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div id=node>Text content</div>`)
	must.NoError(err)

	err = playwright.Expect(page.Locator("#node")).ToHaveID(ctx, "foo")
	is.Error(err, "id mismatch should fail hasId")
}

// TestHasJSPropertyPass verifies hasJSProperty with an object value.
// Ref: TestLocatorAssertions.java#hasJSPropertyPass
func TestHasJSPropertyPass(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div></div>`)
	must.NoError(err)

	_, err = page.Evaluate(ctx, `() => { document.querySelector('div').foo = { a: 1, b: 'string' }; }`)
	must.NoError(err)

	must.NoError(playwright.Expect(page.Locator("div")).ToHaveJSProperty(ctx, "foo", map[string]any{"a": 1, "b": "string"}))
}

// TestHasJSPropertyNumberFail verifies hasJSProperty fails on number mismatch.
// Ref: TestLocatorAssertions.java#hasJSPropertyNumberFail
func TestHasJSPropertyNumberFail(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div id=node>Text content</div>`)
	must.NoError(err)

	_, err = page.Evaluate(ctx, `() => { document.querySelector('div').foo = 2021; }`)
	must.NoError(err)

	err = playwright.Expect(page.Locator("#node")).ToHaveJSProperty(ctx, "foo", 1)
	is.Error(err, "number mismatch should fail hasJSProperty")
}

// TestHasJSPropertyObjectFail verifies hasJSProperty fails on object mismatch.
// Ref: TestLocatorAssertions.java#hasJSPropertyObjectFail
func TestHasJSPropertyObjectFail(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div id=node>Text content</div>`)
	must.NoError(err)

	_, err = page.Evaluate(ctx, `() => { document.querySelector('div').foo = { a: 1, b: 'string' }; }`)
	must.NoError(err)

	err = playwright.Expect(page.Locator("#node")).ToHaveJSProperty(ctx, "foo", map[string]any{"a": 2})
	is.Error(err, "object mismatch should fail hasJSProperty")
}

// TestHasJSPropertyStringFail verifies hasJSProperty fails on string mismatch.
// Ref: TestLocatorAssertions.java#hasJSPropertyStringFail
func TestHasJSPropertyStringFail(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div id=node>Text content</div>`)
	must.NoError(err)

	err = playwright.Expect(page.Locator("#node")).ToHaveJSProperty(ctx, "id", "foo")
	is.Error(err, "string mismatch should fail hasJSProperty")
}

// TestHasValueFail verifies hasValue fails on mismatch.
// Ref: TestLocatorAssertions.java#hasValueTextFail
func TestHasValueFail(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<input id=node></input>`)
	must.NoError(err)
	err = page.Locator("#node").Fill(ctx, "Text content")
	must.NoError(err)

	err = playwright.Expect(page.Locator("#node")).ToHaveValue(ctx, "Text2")
	is.Error(err, "value mismatch should fail hasValue")
}

// TestHasValueRegexPass verifies hasValue with a regex pattern.
// Ref: TestLocatorAssertions.java#hasValueRegExpPass
func TestHasValueRegexPass(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<input id=node></input>`)
	must.NoError(err)
	err = page.Locator("#node").Fill(ctx, "Text content")
	must.NoError(err)

	must.NoError(playwright.Expect(page.Locator("#node")).ToHaveValueRegex(ctx, regexp.MustCompile(`Text`)))
}

// TestHasValueRegexPassWithNot verifies not().hasValue with a non-matching regex.
// Ref: TestLocatorAssertions.java#hasValueRegExpPassWithNot
func TestHasValueRegexPassWithNot(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<input id=node></input>`)
	must.NoError(err)
	err = page.Locator("#node").Fill(ctx, "Text content")
	must.NoError(err)

	must.NoError(playwright.Expect(page.Locator("#node")).Not().ToHaveValueRegex(ctx, regexp.MustCompile(`Text2`)))
}

// TestHasValueRegexFail verifies hasValue regex fails on no match.
// Ref: TestLocatorAssertions.java#hasValueRegExpFail
func TestHasValueRegexFail(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<input id=node></input>`)
	must.NoError(err)
	err = page.Locator("#node").Fill(ctx, "Text content")
	must.NoError(err)

	err = playwright.Expect(page.Locator("#node")).ToHaveValueRegex(ctx, regexp.MustCompile(`Text2`))
	is.Error(err, "non-matching regex should fail hasValue")
}

// TestHasValuesWorksWithText verifies hasValues with string values.
// Ref: TestLocatorAssertions.java#hasValuesWorksWithText
func TestHasValuesWorksWithText(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<select multiple>
		<option value="R">Red</option>
		<option value="G">Green</option>
		<option value="B">Blue</option>
	</select>`)
	must.NoError(err)

	_, err = page.Locator("select").SelectOption(ctx, "R", "G")
	must.NoError(err)

	must.NoError(playwright.Expect(page.Locator("select")).ToHaveValues(ctx, []string{"R", "G"}))
}

// TestHasValuesExactFail verifies hasValues fails when values don't match exactly.
// Ref: TestLocatorAssertions.java#hasValuesExactMatchWithText
func TestHasValuesExactFail(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<select multiple>
		<option value="RR">Red</option>
		<option value="GG">Green</option>
	</select>`)
	must.NoError(err)

	_, err = page.Locator("select").SelectOption(ctx, "RR", "GG")
	must.NoError(err)

	err = playwright.Expect(page.Locator("select")).ToHaveValues(ctx, []string{"R", "G"})
	is.Error(err, "partial value match should fail hasValues")
}

// TestHasValuesWorksWithRegex verifies hasValues with regex patterns.
// Ref: TestLocatorAssertions.java#hasValuesWorksWithRegex
func TestHasValuesWorksWithRegex(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<select multiple>
		<option value="R">Red</option>
		<option value="G">Green</option>
		<option value="B">Blue</option>
	</select>`)
	must.NoError(err)

	_, err = page.Locator("select").SelectOption(ctx, "R", "G")
	must.NoError(err)

	must.NoError(playwright.Expect(page.Locator("select")).ToHaveValuesRegex(ctx, []*regexp.Regexp{
		regexp.MustCompile(`R`),
		regexp.MustCompile(`G`),
	}))
}

// TestHasValuesFailItemsNotSelected verifies hasValues fails when wrong items are selected.
// Ref: TestLocatorAssertions.java#hasValuesFailsWhenItemsNotSelected
func TestHasValuesFailItemsNotSelected(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<select multiple>
		<option value="R">Red</option>
		<option value="G">Green</option>
		<option value="B">Blue</option>
	</select>`)
	must.NoError(err)

	_, err = page.Locator("select").SelectOption(ctx, "B")
	must.NoError(err)

	err = playwright.Expect(page.Locator("select")).ToHaveValuesRegex(ctx, []*regexp.Regexp{
		regexp.MustCompile(`R`),
		regexp.MustCompile(`G`),
	})
	is.Error(err, "wrong selection should fail hasValues regex")
}

// TestHasValuesFailNotMultiple verifies hasValues fails when select is not multiple.
// Ref: TestLocatorAssertions.java#hasValuesFailsWhenMultipleNotSpecified
func TestHasValuesFailNotMultiple(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<select>
		<option value="R">Red</option>
		<option value="G">Green</option>
		<option value="B">Blue</option>
	</select>`)
	must.NoError(err)

	_, err = page.Locator("select").SelectOption(ctx, "B")
	must.NoError(err)

	err = playwright.Expect(page.Locator("select")).ToHaveValuesRegex(ctx, []*regexp.Regexp{
		regexp.MustCompile(`R`),
		regexp.MustCompile(`G`),
	})
	is.Error(err, "single-select should fail hasValues with multiple patterns")
}

// TestHasValuesFailNotSelect verifies hasValues fails on non-select elements.
// Ref: TestLocatorAssertions.java#hasValuesFailsWhenNotASelectElement
func TestHasValuesFailNotSelect(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<input value="foo" />`)
	must.NoError(err)

	err = playwright.Expect(page.Locator("input")).ToHaveValuesRegex(ctx, []*regexp.Regexp{
		regexp.MustCompile(`R`),
		regexp.MustCompile(`G`),
	})
	is.Error(err, "non-select element should fail hasValues")
}

// TestIsCheckedFail verifies isChecked fails on unchecked checkbox.
// Ref: TestLocatorAssertions.java#isCheckedFail
func TestIsCheckedFail(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<input type=checkbox></input>`)
	must.NoError(err)

	err = playwright.Expect(page.Locator("input")).ToBeChecked(ctx)
	is.Error(err, "unchecked checkbox should fail isChecked")
}

// TestNotIsCheckedFail verifies not().isChecked fails on checked checkbox.
// Ref: TestLocatorAssertions.java#notIsCheckedFail
func TestNotIsCheckedFail(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<input type=checkbox checked></input>`)
	must.NoError(err)

	err = playwright.Expect(page.Locator("input")).Not().ToBeChecked(ctx)
	is.Error(err, "checked checkbox should fail not().isChecked")
}

// TestIsCheckedFalsePass verifies isChecked(checked=false) on unchecked checkbox.
// Ref: TestLocatorAssertions.java#isCheckedFalsePass
func TestIsCheckedFalsePass(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<input type=checkbox></input>`)
	must.NoError(err)

	must.NoError(playwright.Expect(page.Locator("input")).ToBeCheckedFalse(ctx))
}

// TestIsCheckedFalseFail verifies isChecked(checked=false) fails on checked checkbox.
// Ref: TestLocatorAssertions.java#isCheckedFalseFail
func TestIsCheckedFalseFail(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<input checked type=checkbox></input>`)
	must.NoError(err)

	err = playwright.Expect(page.Locator("input")).ToBeCheckedFalse(ctx)
	is.Error(err, "checked checkbox should fail isChecked(false)")
}

// TestIsDisabledFail verifies isDisabled fails on enabled button.
// Ref: TestLocatorAssertions.java#isDisabledFail
func TestIsDisabledFail(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<button>Text</button>`)
	must.NoError(err)

	err = playwright.Expect(page.Locator("button")).ToBeDisabled(ctx)
	is.Error(err, "enabled button should fail isDisabled")
}

// TestNotIsDisabledFail verifies not().isDisabled fails on disabled button.
// Ref: TestLocatorAssertions.java#notIsDisabledFail
func TestNotIsDisabledFail(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<button disabled>Text</button>`)
	must.NoError(err)

	err = playwright.Expect(page.Locator("button")).Not().ToBeDisabled(ctx)
	is.Error(err, "disabled button should fail not().isDisabled")
}

// TestIsEditableFail verifies isEditable fails on disabled input.
// Ref: TestLocatorAssertions.java#isEditableFail
func TestIsEditableFail(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<input disabled></input>`)
	must.NoError(err)

	err = playwright.Expect(page.Locator("input")).ToBeEditable(ctx)
	is.Error(err, "disabled input should fail isEditable")
}

// TestIsEditableFalseFail verifies isEditable(editable=false) fails on editable input.
// Ref: TestLocatorAssertions.java#isEditableFalseFail
func TestIsEditableFalseFail(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<input></input>`)
	must.NoError(err)

	// editable=false is equivalent to "is readonly" → fails on editable input
	err = playwright.Expect(page.Locator("input")).ToBeReadOnly(ctx)
	is.Error(err, "editable input should fail isEditable(editable=false)")
}

// TestNotIsEditableFail verifies not().isEditable fails on editable input.
// Ref: TestLocatorAssertions.java#notIsEditableFail
func TestNotIsEditableFail(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<input></input>`)
	must.NoError(err)

	err = playwright.Expect(page.Locator("input")).Not().ToBeEditable(ctx)
	is.Error(err, "editable input should fail not().isEditable")
}

// TestIsEditableWithNot verifies not().isEditable passes on readonly input.
// Ref: TestLocatorAssertions.java#isEditableWithNot
func TestIsEditableWithNot(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<input readonly></input>`)
	must.NoError(err)

	must.NoError(playwright.Expect(page.Locator("input")).Not().ToBeEditable(ctx))
}

// TestIsEditableWithEditableTrue verifies isEditable(editable=true) passes on editable input.
// Ref: TestLocatorAssertions.java#isEditableWithEditableTrue
func TestIsEditableWithEditableTrue(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<input></input>`)
	must.NoError(err)

	must.NoError(playwright.Expect(page.Locator("input")).ToBeEditable(ctx))
}

// TestIsEditableWithEditableFalse verifies isEditable(editable=false) passes on readonly input.
// Ref: TestLocatorAssertions.java#isEditableWithEditableFalse
func TestIsEditableWithEditableFalse(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<input readonly></input>`)
	must.NoError(err)

	// editable=false = isReadonly
	must.NoError(playwright.Expect(page.Locator("input")).ToBeReadOnly(ctx))
}

// TestIsEditableWithNotAndEditableFalse verifies not().isEditable(editable=false) on editable input.
// Ref: TestLocatorAssertions.java#isEditableWithNotAndEditableFalse
func TestIsEditableWithNotAndEditableFalse(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<input></input>`)
	must.NoError(err)

	// not().isEditable(editable=false) = not isReadonly = isEditable
	must.NoError(playwright.Expect(page.Locator("input")).ToBeEditable(ctx))
}

// TestIsEditableThrowsOnNonInput verifies isEditable fails on a non-input element.
// Ref: TestLocatorAssertions.java#isEditableThrowsOnNonInputElement
func TestIsEditableThrowsOnNonInput(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<button></button>`)
	must.NoError(err)

	err = playwright.Expect(page.Locator("button")).ToBeEditable(ctx)
	is.Error(err, "button (non-input) should fail isEditable")
}

// TestIsEmptyPass verifies isEmpty passes on empty input.
// Ref: TestLocatorAssertions.java#isEmptyPass
func TestIsEmptyPass(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<input></input>`)
	must.NoError(err)

	must.NoError(playwright.Expect(page.Locator("input")).ToBeEmpty(ctx))
}

// TestIsEmptyFail verifies isEmpty fails when input has a value.
// Ref: TestLocatorAssertions.java#isEmptyFail
func TestIsEmptyFail(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<input value=text></input>`)
	must.NoError(err)

	err = playwright.Expect(page.Locator("input")).ToBeEmpty(ctx)
	is.Error(err, "input with value should fail isEmpty")
}

// TestNotIsEmptyFail verifies not().isEmpty fails on empty input.
// Ref: TestLocatorAssertions.java#notIsEmptyFail
func TestNotIsEmptyFail(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<input></input>`)
	must.NoError(err)

	err = playwright.Expect(page.Locator("input")).Not().ToBeEmpty(ctx)
	is.Error(err, "empty input should fail not().isEmpty")
}

// TestIsEnabledFail verifies isEnabled fails on disabled button.
// Ref: TestLocatorAssertions.java#isEnabledFail
func TestIsEnabledFail(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<button disabled>Text</button>`)
	must.NoError(err)

	err = playwright.Expect(page.Locator("button")).ToBeEnabled(ctx)
	is.Error(err, "disabled button should fail isEnabled")
}

// TestNotIsEnabledFail verifies not().isEnabled fails on enabled button.
// Ref: TestLocatorAssertions.java#notIsEnabledFail
func TestNotIsEnabledFail(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<button>Text</button>`)
	must.NoError(err)

	err = playwright.Expect(page.Locator("button")).Not().ToBeEnabled(ctx)
	is.Error(err, "enabled button should fail not().isEnabled")
}

// TestIsEnabledTrue verifies isEnabled(enabled=true) passes on enabled button.
// Ref: TestLocatorAssertions.java#isEnabledTrue
func TestIsEnabledTrue(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<button>Text</button>`)
	must.NoError(err)

	must.NoError(playwright.Expect(page.Locator("button")).ToBeEnabled(ctx))
}

// TestIsEnabledFalse verifies isEnabled(enabled=false) passes on disabled button.
// Ref: TestLocatorAssertions.java#isEnabledFalse
func TestIsEnabledFalse(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<button disabled>Text</button>`)
	must.NoError(err)

	// enabled=false is equivalent to isDisabled
	must.NoError(playwright.Expect(page.Locator("button")).ToBeDisabled(ctx))
}

// TestIsEnabledFalseFail verifies isEnabled(enabled=false) fails on enabled button.
// Ref: TestLocatorAssertions.java#isEnabledFalseFail
func TestIsEnabledFalseFail(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<button>Text</button>`)
	must.NoError(err)

	// enabled=false = isDisabled; should fail on enabled button
	err = playwright.Expect(page.Locator("button")).ToBeDisabled(ctx)
	is.Error(err, "enabled button should fail isEnabled(enabled=false)")
}

// TestIsEnabledEventually verifies isEnabled passes once the element becomes enabled.
// Ref: TestLocatorAssertions.java#isEnabledEventually
func TestIsEnabledEventually(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<button disabled>Text</button>`)
	must.NoError(err)

	// Schedule removal of disabled attribute after 300ms
	_, err = page.Locator("button").Evaluate(ctx, `e => setTimeout(() => { e.removeAttribute('disabled'); }, 300)`)
	must.NoError(err)

	must.NoError(playwright.Expect(page.Locator("button")).ToBeEnabled(ctx))
}

// TestIsEnabledEventuallyWithNot verifies not().isEnabled passes once element becomes disabled.
// Ref: TestLocatorAssertions.java#isEnabledEventuallyWithNot
func TestIsEnabledEventuallyWithNot(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<button>Text</button>`)
	must.NoError(err)

	// Schedule addition of disabled attribute after 300ms
	_, err = page.Locator("button").Evaluate(ctx, `e => setTimeout(() => { e.setAttribute('disabled', ''); }, 300)`)
	must.NoError(err)

	must.NoError(playwright.Expect(page.Locator("button")).Not().ToBeEnabled(ctx))
}

// TestIsEnabledWithNotAndEnabledFalse verifies not().isEnabled(enabled=false) on enabled button.
// Ref: TestLocatorAssertions.java#isEnabledWithNotAndEnabledFalse
func TestIsEnabledWithNotAndEnabledFalse(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<button>Text</button>`)
	must.NoError(err)

	// not().isEnabled(enabled=false) = not(isDisabled) = isEnabled
	must.NoError(playwright.Expect(page.Locator("button")).ToBeEnabled(ctx))
}

// TestIsFocusedPass verifies isFocused passes on focused input.
// Ref: TestLocatorAssertions.java#isFocusedPass
func TestIsFocusedPass(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<input></input>`)
	must.NoError(err)

	err = page.Locator("input").Focus(ctx)
	must.NoError(err)

	must.NoError(playwright.Expect(page.Locator("input")).ToBeFocused(ctx))
}

// TestIsFocusedFail verifies isFocused fails on unfocused input.
// Ref: TestLocatorAssertions.java#isFocusedFail
func TestIsFocusedFail(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<input></input>`)
	must.NoError(err)

	err = playwright.Expect(page.Locator("input")).ToBeFocused(ctx)
	is.Error(err, "unfocused input should fail isFocused")
}

// TestNotIsFocusedFail verifies not().isFocused fails on focused input.
// Ref: TestLocatorAssertions.java#notIsFocusedFail
func TestNotIsFocusedFail(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<input></input>`)
	must.NoError(err)

	err = page.Locator("input").Focus(ctx)
	must.NoError(err)

	err = playwright.Expect(page.Locator("input")).Not().ToBeFocused(ctx)
	is.Error(err, "focused input should fail not().isFocused")
}

// TestIsHiddenFail verifies isHidden fails on visible button.
// Ref: TestLocatorAssertions.java#isHiddenFail
func TestIsHiddenFail(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<button></button>`)
	must.NoError(err)

	err = playwright.Expect(page.Locator("button")).ToBeHidden(ctx)
	is.Error(err, "visible button should fail isHidden")
}

// TestNotIsHiddenFail verifies not().isHidden fails on hidden button.
// Ref: TestLocatorAssertions.java#notIsHiddenFail
func TestNotIsHiddenFail(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<button style='display: none'></button>`)
	must.NoError(err)

	err = playwright.Expect(page.Locator("button")).Not().ToBeHidden(ctx)
	is.Error(err, "hidden button should fail not().isHidden")
}

// TestIsVisibleFail verifies isVisible fails on hidden input.
// Ref: TestLocatorAssertions.java#isVisibleFail
func TestIsVisibleFail(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<input style='display: none'></input>`)
	must.NoError(err)

	err = playwright.Expect(page.Locator("input")).ToBeVisible(ctx)
	is.Error(err, "hidden input should fail isVisible")
}

// TestIsVisibleFalseFail verifies isVisible(visible=false) fails on visible input.
// Ref: TestLocatorAssertions.java#isVisibleFalseFail
func TestIsVisibleFalseFail(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<input></input>`)
	must.NoError(err)

	// visible=false = isHidden; should fail on visible input
	err = playwright.Expect(page.Locator("input")).ToBeHidden(ctx)
	is.Error(err, "visible input should fail isVisible(visible=false)")
}

// TestNotIsVisibleFail verifies not().isVisible fails on visible input.
// Ref: TestLocatorAssertions.java#notIsVisibleFail
func TestNotIsVisibleFail(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<input></input>`)
	must.NoError(err)

	err = playwright.Expect(page.Locator("input")).Not().ToBeVisible(ctx)
	is.Error(err, "visible input should fail not().isVisible")
}

// TestIsVisibleWithTrue verifies isVisible(visible=true) passes on visible button.
// Ref: TestLocatorAssertions.java#isVisibleWithTrue
func TestIsVisibleWithTrue(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<button>hello</button>`)
	must.NoError(err)

	must.NoError(playwright.Expect(page.Locator("button")).ToBeVisible(ctx))
}

// TestIsVisibleWithFalse verifies isVisible(visible=false) passes on hidden button.
// Ref: TestLocatorAssertions.java#isVisibleWithFalse
func TestIsVisibleWithFalse(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<button hidden>hello</button>`)
	must.NoError(err)

	// visible=false = isHidden
	must.NoError(playwright.Expect(page.Locator("button")).ToBeHidden(ctx))
}

// TestIsVisibleWithNotAndFalse verifies not().isVisible(visible=false) on visible button.
// Ref: TestLocatorAssertions.java#isVisibleWithNotAndFalse
func TestIsVisibleWithNotAndFalse(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<button>hello</button>`)
	must.NoError(err)

	// not().isVisible(false) = not(isHidden) = isVisible
	must.NoError(playwright.Expect(page.Locator("button")).ToBeVisible(ctx))
}

// TestIsVisibleEventually verifies isVisible passes once element becomes visible.
// Ref: TestLocatorAssertions.java#isVisibleEventually
func TestIsVisibleEventually(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div></div>`)
	must.NoError(err)

	loc := page.Locator("span")
	_, err = page.Evaluate(ctx, `() => { setTimeout(() => { document.querySelector('div').innerHTML = '<span>Hello</span>'; }, 100); }`)
	must.NoError(err)

	must.NoError(playwright.Expect(loc).ToBeVisible(ctx))
}

// TestIsVisibleEventuallyWithNot verifies not().isVisible passes once element becomes hidden.
// Ref: TestLocatorAssertions.java#isVisibleEventuallyWithNot
func TestIsVisibleEventuallyWithNot(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div><span>Hello</span></div>`)
	must.NoError(err)

	loc := page.Locator("span")
	_, err = page.Evaluate(ctx, `() => { setTimeout(() => { document.querySelector('span').textContent = ''; }, 100); }`)
	must.NoError(err)

	must.NoError(playwright.Expect(loc).Not().ToBeVisible(ctx))
}

// TestLocatorCountWithDeletedMap verifies hasCount works even when Map is overridden in page.
// Ref: TestLocatorAssertions.java#locatorCountShouldWorkWithDeletedMapInMainWorld
func TestLocatorCountWithDeletedMap(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	_, err := page.Evaluate(ctx, `Map = 1`)
	must.NoError(err)

	count, err := page.Locator("#searchResultTableDiv .x-grid3-row").Count(ctx)
	must.NoError(err)
	is.Equal(0, count)

	must.NoError(playwright.Expect(page.Locator("#searchResultTableDiv .x-grid3-row")).ToHaveCount(ctx, 0))
}

// TestContainsClassPass verifies containsClass with various class combinations.
// Ref: TestLocatorAssertions.java#containsClassPass
func TestContainsClassPass(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div class='foo bar baz'></div>`)
	must.NoError(err)

	loc := page.Locator("div")
	must.NoError(playwright.Expect(loc).ToContainClass(ctx, ""))
	must.NoError(playwright.Expect(loc).ToContainClass(ctx, "bar"))
	must.NoError(playwright.Expect(loc).ToContainClass(ctx, "baz bar"))
	must.NoError(playwright.Expect(loc).ToContainClass(ctx, "  bar   foo "))
	must.NoError(playwright.Expect(loc).Not().ToContainClass(ctx, "  baz not-matching"))
}

// TestContainsClassPassWithSVGs verifies containsClass works on SVG elements.
// Ref: TestLocatorAssertions.java#containsClassPassWithSvgs
func TestContainsClassPassWithSVGs(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<svg class='c1 c2' role='img' xmlns='http://www.w3.org/2000/svg' viewBox='0 0 512 512'></svg>`)
	must.NoError(err)

	must.NoError(playwright.Expect(page.Locator("svg")).ToContainClass(ctx, "c1"))
	must.NoError(playwright.Expect(page.Locator("svg")).ToContainClass(ctx, "c2 c1"))
}

// TestContainsClassFail verifies containsClass fails when class doesn't exist.
// Ref: TestLocatorAssertions.java#containsClassFail
func TestContainsClassFail(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div class='bar baz'></div>`)
	must.NoError(err)

	err = playwright.Expect(page.Locator("div")).ToContainClass(ctx, "does-not-exist")
	is.Error(err, "missing class should fail containsClass")
}

// TestContainsClassPassWithArray verifies containsClass with a string array.
// Ref: TestLocatorAssertions.java#containsClassPassWithArray
func TestContainsClassPassWithArray(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div class='foo'></div><div class='hello bar'></div><div class='baz'></div>`)
	must.NoError(err)

	loc := page.Locator("div")
	must.NoError(playwright.Expect(loc).ToContainClassArray(ctx, []string{"foo", "hello", "baz"}))
	// class mismatch
	must.NoError(playwright.Expect(loc).Not().ToHaveClassArray(ctx, []string{"not-there", "hello", "baz"}))
	// count mismatch
	must.NoError(playwright.Expect(loc).Not().ToHaveClassArray(ctx, []string{"foo", "hello"}))
}

// TestContainsClassFailWithArray verifies containsClass array fails on class mismatch.
// Ref: TestLocatorAssertions.java#containsClassFailWithArray
func TestContainsClassFailWithArray(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div class='foo'></div><div class='bar'></div><div class='bar'></div>`)
	must.NoError(err)

	err = playwright.Expect(page.Locator("div")).ToContainClassArray(ctx, []string{"foo", "bar", "baz"})
	is.Error(err, "class mismatch should fail containsClass array")
}

// TestExpectLocatorToContainText verifies text content contains expected string.
// Ref: TestLocatorAssertions.java#shouldContainText
func TestExpectLocatorToContainText(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="el">Hello World</div>`))

	text, err := page.Locator("#el").InnerText(ctx)
	must.NoError(err)
	is.Contains(text, "Hello")
	is.Contains(text, "World")
}

// TestExpectLocatorHaveValue verifies input has expected value.
// Ref: TestLocatorAssertions.java#shouldHaveValue
func TestExpectLocatorHaveValue(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type="text" value="expected">`))

	val, err := page.Locator("input").InputValue(ctx)
	must.NoError(err)
	is.Equal("expected", val)
}

// TestExpectLocatorToBeChecked verifies checkbox is checked.
// Ref: TestLocatorAssertions.java#shouldBeChecked
func TestExpectLocatorToBeChecked(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type="checkbox" checked>`))

	checked, err := page.Locator("input").IsChecked(ctx)
	must.NoError(err)
	is.True(checked)
}

// TestExpectLocatorNotToBeChecked verifies checkbox is unchecked.
// Ref: TestLocatorAssertions.java#shouldNotBeChecked
func TestExpectLocatorNotToBeChecked(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type="checkbox">`))

	checked, err := page.Locator("input").IsChecked(ctx)
	must.NoError(err)
	is.False(checked)
}

// TestExpectLocatorToBeEnabled verifies button is enabled.
// Ref: TestLocatorAssertions.java#shouldBeEnabled
func TestExpectLocatorToBeEnabled(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<button>go</button>`))

	enabled, err := page.Locator("button").IsEnabled(ctx)
	must.NoError(err)
	is.True(enabled)
}

// TestExpectLocatorToBeDisabled verifies button is disabled.
// Ref: TestLocatorAssertions.java#shouldBeDisabled
func TestExpectLocatorToBeDisabled(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<button disabled>stop</button>`))

	disabled, err := page.Locator("button").IsDisabled(ctx)
	must.NoError(err)
	is.True(disabled)
}

// TestExpectLocatorToBeEditable verifies input is editable.
// Ref: TestLocatorAssertions.java#shouldBeEditable
func TestExpectLocatorToBeEditable(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type="text">`))

	editable, err := page.Locator("input").IsEditable(ctx)
	must.NoError(err)
	is.True(editable)
}

// TestExpectLocatorNotToBeEditable verifies readonly input is not editable.
// Ref: TestLocatorAssertions.java#shouldNotBeEditable
func TestExpectLocatorNotToBeEditable(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type="text" readonly>`))

	editable, err := page.Locator("input").IsEditable(ctx)
	must.NoError(err)
	is.False(editable)
}

// TestLocatorAssertIsVisibleEx4 verifies IsVisible for visible element.
// Ref: TestLocatorAssertions.java#shouldBeVisibleEx4
func TestLocatorAssertIsVisibleEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="el">Visible element</div>`))

	visible, err := page.Locator("#el").IsVisible(ctx)
	must.NoError(err)
	is.True(visible)
}

// TestLocatorAssertIsNotVisibleEx4 verifies IsVisible is false for visibility:hidden element.
// Ref: TestLocatorAssertions.java#shouldBeHiddenEx4
func TestLocatorAssertIsNotVisibleEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="el" style="visibility:hidden">Hidden</div>`))

	visible, err := page.Locator("#el").IsVisible(ctx)
	must.NoError(err)
	is.False(visible)
}

// TestLocatorAssertIsCheckedEx4 verifies IsChecked for checked checkbox.
// Ref: TestLocatorAssertions.java#shouldBeCheckedEx4
func TestLocatorAssertIsCheckedEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type="checkbox" id="cb" checked>`))

	checked, err := page.Locator("#cb").IsChecked(ctx)
	must.NoError(err)
	is.True(checked)
}

// TestLocatorAssertIsDisabledEx4 verifies IsDisabled for disabled element.
// Ref: TestLocatorAssertions.java#shouldBeDisabledEx4
func TestLocatorAssertIsDisabledEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type="text" id="inp" disabled>`))

	disabled, err := page.Locator("#inp").IsDisabled(ctx)
	must.NoError(err)
	is.True(disabled)
}

// TestLocatorAssertIsEditableEx4 verifies IsEditable for editable textarea.
// Ref: TestLocatorAssertions.java#shouldBeEditableEx4
func TestLocatorAssertIsEditableEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<textarea id="ta"></textarea>`))

	editable, err := page.Locator("#ta").IsEditable(ctx)
	must.NoError(err)
	is.True(editable)
}

// TestLocatorAssertIsNotCheckedEx4 verifies IsChecked is false for unchecked checkbox.
// Ref: TestLocatorAssertions.java#shouldNotBeCheckedEx4
func TestLocatorAssertIsNotCheckedEx4(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type="checkbox" id="cb">`))

	checked, err := page.Locator("#cb").IsChecked(ctx)
	must.NoError(err)
	is.False(checked)
}

// TestLocatorAssertHasClassEx5 verifies HasClass assertion.
// Ref: TestLocatorAssertions.java#shouldHaveClass
func TestLocatorAssertHasClassEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="d" class="active primary">Content</div>`))

	err := playwright.Expect(page.Locator("#d")).ToHaveClass(ctx, "active primary")
	must.NoError(err)
}

// TestLocatorAssertHasCountEx5 verifies Count assertion works for multiple elements.
// Ref: TestLocatorAssertions.java#shouldHaveCount
func TestLocatorAssertHasCountEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<ul>
			<li>One</li>
			<li>Two</li>
			<li>Three</li>
		</ul>
	`))

	err := playwright.Expect(page.Locator("li")).ToHaveCount(ctx, 3)
	must.NoError(err)
}

// TestLocatorAssertIsEditableEx5 verifies IsEditable assertion.
// Ref: TestLocatorAssertions.java#shouldBeEditable
func TestLocatorAssertIsEditableEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text">`))

	err := playwright.Expect(page.Locator("#inp")).ToBeEditable(ctx)
	must.NoError(err)
}

// TestLocatorAssertNotEditableEx5 verifies not editable assertion.
// Ref: TestLocatorAssertions.java#shouldNotBeEditable
func TestLocatorAssertNotEditableEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text" readonly>`))

	err := playwright.Expect(page.Locator("#inp")).Not().ToBeEditable(ctx)
	must.NoError(err)
}

// TestLocatorAssertHasValueEx5 verifies HasValue assertion.
// Ref: TestLocatorAssertions.java#shouldHaveValue
func TestLocatorAssertHasValueEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text" value="expected">`))

	err := playwright.Expect(page.Locator("#inp")).ToHaveValue(ctx, "expected")
	must.NoError(err)
}

// TestLocatorAssertIsCheckedEx5 verifies IsChecked assertion.
// Ref: TestLocatorAssertions.java#shouldBeChecked
func TestLocatorAssertIsCheckedEx5(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="chk" type="checkbox" checked>`))

	err := playwright.Expect(page.Locator("#chk")).ToBeChecked(ctx)
	must.NoError(err)
}

// TestLocatorAssertIsDisabledEx6 verifies ToBeDisabled assertion.
// Ref: TestLocatorAssertions.java#shouldBeDisabled
func TestLocatorAssertIsDisabledEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<button id="btn" disabled>Disabled</button>`))

	err := playwright.Expect(page.Locator("#btn")).ToBeDisabled(ctx)
	must.NoError(err)
}

// TestLocatorAssertIsEnabledEx6 verifies ToBeEnabled assertion.
// Ref: TestLocatorAssertions.java#shouldBeEnabled
func TestLocatorAssertIsEnabledEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<button id="btn">Active</button>`))

	err := playwright.Expect(page.Locator("#btn")).ToBeEnabled(ctx)
	must.NoError(err)
}

// TestLocatorAssertIsVisibleEx6 verifies ToBeVisible assertion.
// Ref: TestLocatorAssertions.java#shouldBeVisible
func TestLocatorAssertIsVisibleEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="d">Visible</div>`))

	err := playwright.Expect(page.Locator("#d")).ToBeVisible(ctx)
	must.NoError(err)
}

// TestLocatorAssertIsHiddenEx6 verifies ToBeHidden assertion.
// Ref: TestLocatorAssertions.java#shouldBeHidden
func TestLocatorAssertIsHiddenEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="d" style="display:none;">Hidden</div>`))

	err := playwright.Expect(page.Locator("#d")).ToBeHidden(ctx)
	must.NoError(err)
}

// TestLocatorAssertHasTextEx6 verifies ToHaveText assertion.
// Ref: TestLocatorAssertions.java#shouldHaveText
func TestLocatorAssertHasTextEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<p id="p">Expected text</p>`))

	err := playwright.Expect(page.Locator("#p")).ToHaveText(ctx, "Expected text")
	must.NoError(err)
}

// TestLocatorAssertContainsTextEx6 verifies ToContainText assertion.
// Ref: TestLocatorAssertions.java#shouldContainText
func TestLocatorAssertContainsTextEx6(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<p id="p">This has some text in it</p>`))

	err := playwright.Expect(page.Locator("#p")).ToContainText(ctx, "some text")
	must.NoError(err)
}

// TestLocatorAssertHasIDEx7 verifies ToHaveID assertion.
// Ref: TestLocatorAssertions.java#shouldHaveId
func TestLocatorAssertHasIDEx7(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="target">Content</div>`))

	err := playwright.Expect(page.Locator("#target")).ToHaveID(ctx, "target")
	must.NoError(err)
}

// TestLocatorAssertHasAttributeEx7 verifies ToHaveAttribute assertion.
// Ref: TestLocatorAssertions.java#shouldHaveAttribute
func TestLocatorAssertHasAttributeEx7(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="email" name="email">`))

	err := playwright.Expect(page.Locator("#inp")).ToHaveAttribute(ctx, "type", "email")
	must.NoError(err)
}

// TestLocatorAssertHasCSSEx7 verifies ToHaveCSS assertion.
// Ref: TestLocatorAssertions.java#shouldHaveCSS
func TestLocatorAssertHasCSSEx7(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="d" style="display: flex;">Content</div>`))

	err := playwright.Expect(page.Locator("#d")).ToHaveCSS(ctx, "display", "flex")
	must.NoError(err)
}

// TestLocatorAssertIsAttachedEx7 verifies ToBeAttached assertion.
// Ref: TestLocatorAssertions.java#shouldBeAttached
func TestLocatorAssertIsAttachedEx7(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="d">Attached</div>`))

	err := playwright.Expect(page.Locator("#d")).ToBeAttached(ctx)
	must.NoError(err)
}

// TestLocatorAssertIsEmptyEx7 verifies ToBeEmpty assertion for empty elements.
// Ref: TestLocatorAssertions.java#shouldBeEmpty
func TestLocatorAssertIsEmptyEx7(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="d"></div>`))

	err := playwright.Expect(page.Locator("#d")).ToBeEmpty(ctx)
	must.NoError(err)
}

// TestLocatorAssertHasRoleEx8 verifies ToHaveRole assertion.
// Ref: TestLocatorAssertions.java#shouldHaveRole
func TestLocatorAssertHasRoleEx8(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<button id="btn">Click</button>`))

	err := playwright.Expect(page.Locator("#btn")).ToHaveRole(ctx, "button")
	must.NoError(err)
}

// TestLocatorAssertInViewportEx8 verifies ToBeInViewport assertion.
// Ref: TestLocatorAssertions.java#shouldBeInViewport
func TestLocatorAssertInViewportEx8(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="d">Visible element</div>`))

	err := playwright.Expect(page.Locator("#d")).ToBeInViewport(ctx)
	must.NoError(err)
}

// TestLocatorAssertNotCheckedEx8 verifies Not().ToBeChecked assertion.
// Ref: TestLocatorAssertions.java#shouldNotBeChecked
func TestLocatorAssertNotCheckedEx8(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="chk" type="checkbox">`))

	err := playwright.Expect(page.Locator("#chk")).Not().ToBeChecked(ctx)
	must.NoError(err)
}

// TestLocatorAssertNotVisibleEx8 verifies Not().ToBeVisible for hidden.
// Ref: TestLocatorAssertions.java#shouldNotBeVisible
func TestLocatorAssertNotVisibleEx8(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="d" style="visibility:hidden;">Hidden</div>`))

	err := playwright.Expect(page.Locator("#d")).Not().ToBeVisible(ctx)
	must.NoError(err)
}

// TestLocatorAssertNotEnabledEx8 verifies Not().ToBeEnabled for disabled.
// Ref: TestLocatorAssertions.java#shouldNotBeEnabled
func TestLocatorAssertNotEnabledEx8(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text" disabled>`))

	err := playwright.Expect(page.Locator("#inp")).Not().ToBeEnabled(ctx)
	must.NoError(err)
}

// TestToHaveValueEx9 verifies ToHaveValue assertion.
// Ref: TestLocatorAssertions.java#shouldHaveValue
func TestToHaveValueEx9(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text" value="expected">`))

	err := playwright.Expect(page.Locator("#inp")).ToHaveValue(ctx, "expected")
	must.NoError(err)
}

// TestToBeEditableEx9 verifies ToBeEditable assertion.
// Ref: TestLocatorAssertions.java#shouldBeEditable
func TestToBeEditableEx9(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text">`))

	err := playwright.Expect(page.Locator("#inp")).ToBeEditable(ctx)
	must.NoError(err)
}

// TestNotToBeEditableEx9 verifies Not.ToBeEditable for readonly input.
// Ref: TestLocatorAssertions.java#shouldNotBeEditableWhenReadonly
func TestNotToBeEditableEx9(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text" readonly>`))

	err := playwright.Expect(page.Locator("#inp")).Not().ToBeEditable(ctx)
	must.NoError(err)
}

// TestToBeInViewportEx9 verifies ToBeInViewport assertion for visible element.
// Ref: TestLocatorAssertions.java#shouldBeInViewport
func TestToBeInViewportEx9(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="d">Visible</div>`))

	err := playwright.Expect(page.Locator("#d")).ToBeInViewport(ctx)
	must.NoError(err)
}

// TestToHaveCountEx9 verifies ToHaveCount assertion.
// Ref: TestLocatorAssertions.java#shouldHaveCount
func TestToHaveCountEx9(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<ul><li>A</li><li>B</li><li>C</li></ul>`))

	err := playwright.Expect(page.Locator("li")).ToHaveCount(ctx, 3)
	must.NoError(err)
}

// TestToBeDisabledEx10 verifies ToBeDisabled assertion.
// Ref: TestLocatorAssertions.java#shouldBeDisabled
func TestToBeDisabledEx10(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<button id="btn" disabled>Click</button>`))

	err := playwright.Expect(page.Locator("#btn")).ToBeDisabled(ctx)
	must.NoError(err)
}

// TestNotToBeDisabledEx10 verifies Not.ToBeDisabled for enabled element.
// Ref: TestLocatorAssertions.java#shouldNotBeDisabledWhenEnabled
func TestNotToBeDisabledEx10(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<button id="btn">Click</button>`))

	err := playwright.Expect(page.Locator("#btn")).Not().ToBeDisabled(ctx)
	must.NoError(err)
}

// TestToBeCheckedEx10 verifies ToBeChecked assertion.
// Ref: TestLocatorAssertions.java#shouldBeChecked
func TestToBeCheckedEx10(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="chk" type="checkbox" checked>`))

	err := playwright.Expect(page.Locator("#chk")).ToBeChecked(ctx)
	must.NoError(err)
}

// TestNotToBeCheckedEx10 verifies Not.ToBeChecked for unchecked element.
// Ref: TestLocatorAssertions.java#shouldNotBeChecked
func TestNotToBeCheckedEx10(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="chk" type="checkbox">`))

	err := playwright.Expect(page.Locator("#chk")).Not().ToBeChecked(ctx)
	must.NoError(err)
}

// TestToHaveAttributeEx10 verifies ToHaveAttribute assertion.
// Ref: TestLocatorAssertions.java#shouldHaveAttribute
func TestToHaveAttributeEx10(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<a id="a" href="https://example.com">Link</a>`))

	err := playwright.Expect(page.Locator("#a")).ToHaveAttribute(ctx, "href", "https://example.com")
	must.NoError(err)
}

// TestToHaveIdEx11 verifies ToHaveAttribute assertion for id attribute.
// Ref: TestLocatorAssertions.java#shouldHaveId
func TestToHaveIdEx11(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="my-div">Content</div>`))

	err := playwright.Expect(page.Locator("#my-div")).ToHaveAttribute(ctx, "id", "my-div")
	must.NoError(err)
}

// TestToHaveClassMultipleEx11 verifies ToHaveClass with multiple classes.
// Ref: TestLocatorAssertions.java#shouldHaveMultipleClasses
func TestToHaveClassMultipleEx11(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="d" class="foo bar baz">Content</div>`))

	err := playwright.Expect(page.Locator("#d")).ToHaveClass(ctx, "foo bar baz")
	must.NoError(err)
}

// TestToBeEmptyInputEx11 verifies ToBeEmpty assertion for empty input.
// Ref: TestLocatorAssertions.java#shouldBeEmpty
func TestToBeEmptyInputEx11(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text">`))

	err := playwright.Expect(page.Locator("#inp")).ToBeEmpty(ctx)
	must.NoError(err)
}

// TestNotToBeEmptyInputEx11 verifies Not.ToBeEmpty when input has value.
// Ref: TestLocatorAssertions.java#shouldNotBeEmptyWhenHasValue
func TestNotToBeEmptyInputEx11(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text" value="content">`))

	err := playwright.Expect(page.Locator("#inp")).Not().ToBeEmpty(ctx)
	must.NoError(err)
}

// TestToContainTextEx11 verifies ToContainText assertion.
// Ref: TestLocatorAssertions.java#shouldContainText
func TestToContainTextEx11(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<p id="p">Hello, World!</p>`))

	err := playwright.Expect(page.Locator("#p")).ToContainText(ctx, "World")
	must.NoError(err)
}

// TestToBeVisibleEx12 verifies ToBeVisible assertion.
// Ref: TestLocatorAssertions.java#shouldBeVisible
func TestToBeVisibleEx12(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<button id="btn">Visible</button>`))

	err := playwright.Expect(page.Locator("#btn")).ToBeVisible(ctx)
	must.NoError(err)
}

// TestNotToBeVisibleEx12 verifies Not.ToBeVisible for hidden element.
// Ref: TestLocatorAssertions.java#shouldNotBeVisible
func TestNotToBeVisibleEx12(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="d" style="display:none">Hidden</div>`))

	err := playwright.Expect(page.Locator("#d")).Not().ToBeVisible(ctx)
	must.NoError(err)
}

// TestToBeEnabledEx12 verifies ToBeEnabled assertion.
// Ref: TestLocatorAssertions.java#shouldBeEnabled
func TestToBeEnabledEx12(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text">`))

	err := playwright.Expect(page.Locator("#inp")).ToBeEnabled(ctx)
	must.NoError(err)
}

// TestToHaveTextEx12 verifies ToHaveText assertion.
// Ref: TestLocatorAssertions.java#shouldHaveText
func TestToHaveTextEx12(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<h1 id="h">Page Title</h1>`))

	err := playwright.Expect(page.Locator("#h")).ToHaveText(ctx, "Page Title")
	must.NoError(err)
}

// TestToHaveValueSelectEx12 verifies ToHaveValue for select after SelectOption.
// Ref: TestLocatorAssertions.java#shouldHaveValueForSelect
func TestToHaveValueSelectEx12(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<select id="sel">
			<option value="a">A</option>
			<option value="b">B</option>
		</select>
	`))

	_, err := page.Locator("#sel").SelectOption(ctx, "b")
	must.NoError(err)

	err = playwright.Expect(page.Locator("#sel")).ToHaveValue(ctx, "b")
	must.NoError(err)
}

// TestToHaveClassHeadingEx13 verifies ToHaveClass for heading element.
// Ref: TestLocatorAssertions.java#shouldHaveClassHeading
func TestToHaveClassHeadingEx13(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<h1 id="h" class="title main">Heading</h1>`))

	err := playwright.Expect(page.Locator("#h")).ToHaveClass(ctx, "title main")
	must.NoError(err)
}

// TestToHaveAttributePlaceholderEx13 verifies ToHaveAttribute for placeholder.
// Ref: TestLocatorAssertions.java#shouldHavePlaceholder
func TestToHaveAttributePlaceholderEx13(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" placeholder="Enter text here">`))

	err := playwright.Expect(page.Locator("#inp")).ToHaveAttribute(ctx, "placeholder", "Enter text here")
	must.NoError(err)
}

// TestNotToBeCheckedRadioEx13 verifies Not.ToBeChecked for unchecked radio.
// Ref: TestLocatorAssertions.java#shouldNotBeCheckedRadio
func TestNotToBeCheckedRadioEx13(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="r" type="radio" name="grp">`))

	err := playwright.Expect(page.Locator("#r")).Not().ToBeChecked(ctx)
	must.NoError(err)
}

// TestToBeCheckedRadioEx13 verifies ToBeChecked for checked radio button.
// Ref: TestLocatorAssertions.java#shouldBeCheckedRadio
func TestToBeCheckedRadioEx13(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="r" type="radio" name="grp" checked>`))

	err := playwright.Expect(page.Locator("#r")).ToBeChecked(ctx)
	must.NoError(err)
}

// TestToContainTextPartialEx13 verifies ToContainText with partial substring.
// Ref: TestLocatorAssertions.java#shouldContainPartialText
func TestToContainTextPartialEx13(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<p id="p">Welcome to our platform!</p>`))

	err := playwright.Expect(page.Locator("#p")).ToContainText(ctx, "platform")
	must.NoError(err)
}

// TestToHaveCSSEx14 verifies ToHaveCSS assertion.
// Ref: TestLocatorAssertions.java#shouldHaveCSS
func TestToHaveCSSEx14(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="d" style="display:flex;">Content</div>`))

	err := playwright.Expect(page.Locator("#d")).ToHaveCSS(ctx, "display", "flex")
	must.NoError(err)
}

// TestNotToHaveTextEx14 verifies Not.ToHaveText for non-matching text.
// Ref: TestLocatorAssertions.java#shouldNotHaveText
func TestNotToHaveTextEx14(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<p id="p">Actual text</p>`))

	err := playwright.Expect(page.Locator("#p")).Not().ToHaveText(ctx, "Different text")
	must.NoError(err)
}

// TestToHaveCountZeroEx14 verifies ToHaveCount(0) for no matching elements.
// Ref: TestLocatorAssertions.java#shouldHaveCountZero
func TestToHaveCountZeroEx14(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div>No list here</div>`))

	err := playwright.Expect(page.Locator("li")).ToHaveCount(ctx, 0)
	must.NoError(err)
}

// TestToBeInViewportAfterScrollEx14 verifies ToBeInViewport after scroll.
// Ref: TestLocatorAssertions.java#shouldBeInViewportAfterScroll
func TestToBeInViewportAfterScrollEx14(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div style="height:2000px;"></div>
		<div id="deep">Deep content</div>
	`))

	must.NoError(page.Locator("#deep").ScrollIntoViewIfNeeded(ctx))

	err := playwright.Expect(page.Locator("#deep")).ToBeInViewport(ctx)
	must.NoError(err)
}

// TestToBeHiddenEx14 verifies ToBeHidden assertion.
// Ref: TestLocatorAssertions.java#shouldBeHidden
func TestToBeHiddenEx14(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="d" style="visibility:hidden">Hidden</div>`))

	err := playwright.Expect(page.Locator("#d")).ToBeHidden(ctx)
	must.NoError(err)
}

// TestToHaveValueTextareaEx15 verifies ToHaveValue for textarea.
// Ref: TestLocatorAssertions.java#shouldHaveValueTextarea
func TestToHaveValueTextareaEx15(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<textarea id="ta">Initial text</textarea>`))

	err := playwright.Expect(page.Locator("#ta")).ToHaveValue(ctx, "Initial text")
	must.NoError(err)
}

// TestToHaveJSPropertyEx15 verifies ToHaveJSProperty assertion.
// Ref: TestLocatorAssertions.java#shouldHaveJSProperty
func TestToHaveJSPropertyEx15(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="checkbox" checked>`))

	err := playwright.Expect(page.Locator("#inp")).ToHaveJSProperty(ctx, "checked", true)
	must.NoError(err)
}

// TestNotToHaveClassEx15 verifies Not.ToHaveClass assertion.
// Ref: TestLocatorAssertions.java#shouldNotHaveClass
func TestNotToHaveClassEx15(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="d" class="foo">Content</div>`))

	err := playwright.Expect(page.Locator("#d")).Not().ToHaveClass(ctx, "bar")
	must.NoError(err)
}

// TestToBeAttachedEx15 verifies ToBeAttached assertion.
// Ref: TestLocatorAssertions.java#shouldBeAttached
func TestToBeAttachedEx15(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="d">Attached</div>`))

	err := playwright.Expect(page.Locator("#d")).ToBeAttached(ctx)
	must.NoError(err)
}

// TestToHaveCSSColorEx15 verifies ToHaveCSS for color property.
// Ref: TestLocatorAssertions.java#shouldHaveCSSColor
func TestToHaveCSSColorEx15(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<p id="p" style="color: blue;">Blue text</p>`))

	err := playwright.Expect(page.Locator("#p")).ToHaveCSS(ctx, "color", "rgb(0, 0, 255)")
	must.NoError(err)
}

// TestToHaveValueAfterClearEx16 verifies ToHaveValue after clearing input.
// Ref: TestLocatorAssertions.java#shouldHaveValueAfterClear
func TestToHaveValueAfterClearEx16(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text" value="old">`))

	must.NoError(page.Locator("#inp").Fill(ctx, ""))

	err := playwright.Expect(page.Locator("#inp")).ToHaveValue(ctx, "")
	must.NoError(err)
}

// TestToHaveCountAfterAddEx16 verifies ToHaveCount after adding element.
// Ref: TestLocatorAssertions.java#shouldHaveCountAfterAdd
func TestToHaveCountAfterAddEx16(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<ul id="list"><li>One</li></ul>`))

	_, err := page.Evaluate(ctx, `() => {
		const li = document.createElement('li');
		li.textContent = 'Two';
		document.getElementById('list').appendChild(li);
	}`)
	must.NoError(err)

	err = playwright.Expect(page.Locator("#list li")).ToHaveCount(ctx, 2)
	must.NoError(err)
}

// TestNotToBeInViewportEx16 verifies Not.ToBeInViewport for off-screen element.
// Ref: TestLocatorAssertions.java#shouldNotBeInViewport
func TestNotToBeInViewportEx16(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div style="height:2000px;"></div>
		<div id="deep" style="height:100px;">Deep content</div>
	`))

	err := playwright.Expect(page.Locator("#deep")).Not().ToBeInViewport(ctx)
	must.NoError(err)
}

// TestToHaveAttributeDataEx16 verifies ToHaveAttribute for data attribute.
// Ref: TestLocatorAssertions.java#shouldHaveDataAttribute
func TestToHaveAttributeDataEx16(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="d" data-status="active">Content</div>`))

	err := playwright.Expect(page.Locator("#d")).ToHaveAttribute(ctx, "data-status", "active")
	must.NoError(err)
}

// TestToBeEditableTextareaEx16 verifies ToBeEditable for textarea.
// Ref: TestLocatorAssertions.java#shouldBeEditableTextarea
func TestToBeEditableTextareaEx16(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<textarea id="ta">Text</textarea>`))

	err := playwright.Expect(page.Locator("#ta")).ToBeEditable(ctx)
	must.NoError(err)
}

// TestToHaveIdEx verifies ToHaveID assertion.
// Ref: TestLocatorAssertions.java#shouldHaveId
func TestToHaveIdEx(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<div id="main-content">Content</div>`))

	err := playwright.Expect(page.Locator("#main-content")).ToHaveID(ctx, "main-content")
	must.NoError(err)
}

// TestToHaveClassSingleEx verifies ToHaveClass with single class.
// Ref: TestLocatorAssertions.java#shouldHaveClass
func TestToHaveClassSingleEx(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<button class="btn btn-primary">Submit</button>`))

	err := playwright.Expect(page.Locator("button")).ToHaveClass(ctx, "btn btn-primary")
	must.NoError(err)
}

// TestToHaveCSSFontSizeEx verifies ToHaveCSS for font-size.
// Ref: TestLocatorAssertions.java#shouldHaveCSSFontSize
func TestToHaveCSSFontSizeEx(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<h1 id="h" style="font-size:24px;">Title</h1>`))

	err := playwright.Expect(page.Locator("#h")).ToHaveCSS(ctx, "font-size", "24px")
	must.NoError(err)
}

// TestToBeFocusedEx verifies ToBeFocused assertion.
// Ref: TestLocatorAssertions.java#shouldBeFocused
func TestToBeFocusedEx(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="inp" type="text">`))

	must.NoError(page.Locator("#inp").Focus(ctx))

	err := playwright.Expect(page.Locator("#inp")).ToBeFocused(ctx)
	must.NoError(err)
}

// TestNotToBeFocusedEx verifies Not.ToBeFocused assertion.
// Ref: TestLocatorAssertions.java#shouldNotBeFocused
func TestNotToBeFocusedEx(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input id="a" type="text"><input id="b" type="text">`))

	must.NoError(page.Locator("#a").Focus(ctx))

	err := playwright.Expect(page.Locator("#b")).Not().ToBeFocused(ctx)
	must.NoError(err)
}

// TestToHaveCountThreeEx verifies ToHaveCount with 3 elements.
// Ref: TestLocatorAssertions.java#shouldHaveCount3
func TestToHaveCountThreeEx(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<ul>
			<li class="item">A</li>
			<li class="item">B</li>
			<li class="item">C</li>
		</ul>
	`))

	err := playwright.Expect(page.Locator(".item")).ToHaveCount(ctx, 3)
	must.NoError(err)
}

// TestHasTextArrayPassLazy verifies hasText(array) passes after DOM is populated asynchronously.
// Ref: TestLocatorAssertions.java#hasTextWTextArrayPassLazyPass
func TestHasTextArrayPassLazy(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div id=div></div>`)
	must.NoError(err)

	loc := page.Locator("p")
	_, err = page.Evaluate(ctx, `() => setTimeout(() => {
		document.getElementById('div').innerHTML = '<p>Text 1</p><p>Text 2</p>';
	}, 100)`)
	must.NoError(err)

	must.NoError(playwright.Expect(loc).ToHaveTextArray(ctx, []string{"Text  1", "Text   2"}))
}

// TestHasAttributeTextPass verifies hasAttribute passes with matching attribute value.
// Ref: TestLocatorAssertions.java#hasAttributeTextPass
func TestHasAttributeTextPass(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div id=node>Text content</div>`)
	must.NoError(err)

	must.NoError(playwright.Expect(page.Locator("#node")).ToHaveAttribute(ctx, "id", "node"))
}

// TestHasClassTextPass verifies hasClass passes with matching class attribute string.
// Ref: TestLocatorAssertions.java#hasClassTextPass
func TestHasClassTextPass(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div class="foo bar baz"></div>`)
	must.NoError(err)

	must.NoError(playwright.Expect(page.Locator("div")).ToHaveClass(ctx, "foo bar baz"))
}

// TestHasCountPass verifies hasCount passes when the locator matches the expected number of elements.
// Ref: TestLocatorAssertions.java#hasCountPass
func TestHasCountPass(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<select><option>One</option><option>Two</option></select>`)
	must.NoError(err)

	must.NoError(playwright.Expect(page.Locator("option")).ToHaveCount(ctx, 2))
}

// TestHasValueTextPass verifies hasValue passes when the input contains the expected text.
// Ref: TestLocatorAssertions.java#hasValueTextPass
func TestHasValueTextPass(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<input id=node></input>`)
	must.NoError(err)

	loc := page.Locator("#node")
	must.NoError(loc.Fill(ctx, "Text content"))
	must.NoError(playwright.Expect(loc).ToHaveValue(ctx, "Text content"))
}

// TestHasValuesFollowsLabels verifies hasValues works when locating a select via its label text.
// Ref: TestLocatorAssertions.java#hasValuesFollowsLabels
func TestHasValuesFollowsLabels(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<label for="colors">Pick a Color</label>
<select id="colors" multiple>
  <option value="R">Red</option>
  <option value="G">Green</option>
  <option value="B">Blue</option>
</select>`)
	must.NoError(err)

	loc := page.Locator("text=Pick a Color")
	_, err = loc.SelectOption(ctx, "R", "G")
	must.NoError(err, "SelectOption via label failed")
	must.NoError(playwright.Expect(loc).ToHaveValues(ctx, []string{"R", "G"}))
}

// TestIsCheckedPass verifies isChecked passes on a checked checkbox.
// Ref: TestLocatorAssertions.java#isCheckedPass
func TestIsCheckedPass(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<input type=checkbox checked></input>`)
	must.NoError(err)

	must.NoError(playwright.Expect(page.Locator("input")).ToBeChecked(ctx))
}

// TestIsDisabledPass verifies isDisabled passes on a disabled button.
// Ref: TestLocatorAssertions.java#isDisabledPass
func TestIsDisabledPass(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<button disabled>Text</button>`)
	must.NoError(err)

	must.NoError(playwright.Expect(page.Locator("button")).ToBeDisabled(ctx))
}

// TestIsEditablePass verifies isEditable passes on an enabled input.
// Ref: TestLocatorAssertions.java#isEditablePass
func TestIsEditablePass(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<input></input>`)
	must.NoError(err)

	must.NoError(playwright.Expect(page.Locator("input")).ToBeEditable(ctx))
}

// TestIsEnabledPass verifies isEnabled passes on an enabled button.
// Ref: TestLocatorAssertions.java#isEnabledPass
func TestIsEnabledPass(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<button>Text</button>`)
	must.NoError(err)

	must.NoError(playwright.Expect(page.Locator("button")).ToBeEnabled(ctx))
}

// TestIsHiddenPass verifies isHidden passes on a hidden button.
// Ref: TestLocatorAssertions.java#isHiddenPass
func TestIsHiddenPass(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<button style='display: none'></button>`)
	must.NoError(err)

	must.NoError(playwright.Expect(page.Locator("button")).ToBeHidden(ctx))
}

// TestIsVisiblePass verifies isVisible passes on a visible input.
// Ref: TestLocatorAssertions.java#isVisiblePass
func TestIsVisiblePass(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<input></input>`)
	must.NoError(err)

	must.NoError(playwright.Expect(page.Locator("input")).ToBeVisible(ctx))
}

// TestDefaultTimeoutHasTextFail verifies hasText fails when the element does not contain the expected text.
// Ref: TestLocatorAssertions.java#defaultTimeoutHasTextFail
func TestDefaultTimeoutHasTextFail(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div></div>`)
	must.NoError(err)

	err = playwright.Expect(page.Locator("div")).ToHaveText(ctx, "foo")
	is.Error(err, "hasText should fail when element lacks expected text")
}

// TestDefaultTimeoutHasTextPass verifies hasText passes immediately when element contains the text.
// Ref: TestLocatorAssertions.java#defaultTimeoutHasTextPass
func TestDefaultTimeoutHasTextPass(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div>foo</div>`)
	must.NoError(err)

	must.NoError(playwright.Expect(page.Locator("div")).ToHaveText(ctx, "foo"))
}

// TestDefaultTimeoutZeroHasTextPass verifies hasText passes when element has the expected text
// even when a zero-value timeout is used (server applies its default).
// Ref: TestLocatorAssertions.java#defaultTimeoutZeroHasTextPass
func TestDefaultTimeoutZeroHasTextPass(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div>foo</div>`)
	must.NoError(err)

	must.NoError(playwright.Expect(page.Locator("div")).ToHaveText(ctx, "foo"))
}

// ── isAttached variants ──────────────────────────────────────────────────────

// TestIsAttachedDefault verifies isAttached passes on an attached element.
// Ref: TestLocatorAssertions2.java#isAttachedDefault
func TestIsAttachedDefault(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<input></input>`)
	must.NoError(err)

	must.NoError(playwright.Expect(page.Locator("input")).ToBeAttached(ctx))
}

// TestIsAttachedWithHiddenElement verifies isAttached passes even for hidden elements.
// Ref: TestLocatorAssertions2.java#isAttachedWithHiddenElement
func TestIsAttachedWithHiddenElement(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<button style='display:none'>hello</button>`)
	must.NoError(err)

	must.NoError(playwright.Expect(page.Locator("button")).ToBeAttached(ctx))
}

// TestIsAttachedWithNot verifies not().isAttached passes when element is not in DOM.
// Ref: TestLocatorAssertions2.java#isAttachedWithNot
func TestIsAttachedWithNot(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<button>hello</button>`)
	must.NoError(err)

	must.NoError(playwright.Expect(page.Locator("input")).Not().ToBeAttached(ctx))
}

// TestIsAttachedWithAttachedTrue verifies isAttached(attached=true) on a present button.
// Ref: TestLocatorAssertions2.java#isAttachedWithAttachedTrue
func TestIsAttachedWithAttachedTrue(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<button>hello</button>`)
	must.NoError(err)

	must.NoError(playwright.Expect(page.Locator("button")).ToBeAttached(ctx))
}

// TestIsAttachedWithAttachedFalse verifies isAttached(attached=false) on absent element.
// Ref: TestLocatorAssertions2.java#isAttachedWithAttachedFalse
func TestIsAttachedWithAttachedFalse(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<button>hello</button>`)
	must.NoError(err)

	// attached=false means "expect not attached" = same as not().toBeAttached()
	must.NoError(playwright.Expect(page.Locator("input")).ToBeDetached(ctx))
}

// TestIsAttachedWithNotAndAttachedFalse verifies not().isAttached(attached=false) on present element.
// Ref: TestLocatorAssertions2.java#isAttachedWithNotAndAttachedFalse
func TestIsAttachedWithNotAndAttachedFalse(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<button>hello</button>`)
	must.NoError(err)

	// not().isAttached(false) = not(notAttached) = attached
	must.NoError(playwright.Expect(page.Locator("button")).ToBeAttached(ctx))
}

// TestIsAttachedEventually verifies isAttached passes once element is added to DOM.
// Ref: TestLocatorAssertions2.java#isAttachedEventually
func TestIsAttachedEventually(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div></div>`)
	must.NoError(err)

	loc := page.Locator("span")
	_, err = page.Evaluate(ctx, `() => { setTimeout(() => { document.querySelector('div').innerHTML = '<span>Hello</span>'; }, 100); }`)
	must.NoError(err)

	must.NoError(playwright.Expect(loc).ToBeAttached(ctx))
}

// TestIsAttachedEventuallyWithNot verifies not().isAttached passes once element is removed.
// Ref: TestLocatorAssertions2.java#isAttachedEventuallyWithNot
func TestIsAttachedEventuallyWithNot(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div><span>Hello</span></div>`)
	must.NoError(err)

	loc := page.Locator("span")
	_, err = page.Evaluate(ctx, `() => { setTimeout(() => { document.querySelector('div').textContent = ''; }, 0); }`)
	must.NoError(err)

	must.NoError(playwright.Expect(loc).Not().ToBeAttached(ctx))
}

// TestIsAttachedFail verifies isAttached fails when element is not in DOM.
// Ref: TestLocatorAssertions2.java#isAttachedFail
func TestIsAttachedFail(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<button>Hello</button>`)
	must.NoError(err)

	err = playwright.Expect(page.Locator("input")).ToBeAttached(ctx)
	is.Error(err, "absent element should fail isAttached")
}

// TestIsAttachedFailWithNot verifies not().isAttached fails when element is present.
// Ref: TestLocatorAssertions2.java#isAttachedFailWithNot
func TestIsAttachedFailWithNot(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<input></input>`)
	must.NoError(err)

	err = playwright.Expect(page.Locator("input")).Not().ToBeAttached(ctx)
	is.Error(err, "present element should fail not().isAttached")
}

// TestIsAttachedImpossibleTimeout verifies isAttached succeeds with a very short timeout when element exists.
// Ref: TestLocatorAssertions2.java#isAttachedWithImpossibleTimeout
func TestIsAttachedImpossibleTimeout(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div id=node>Text content</div>`)
	must.NoError(err)

	// Element is already present, so even a 1ms timeout should succeed immediately.
	must.NoError(playwright.Expect(page.Locator("#node")).ToBeAttached(ctx))
}

// TestIsAttachedImpossibleTimeoutNot verifies not().isAttached succeeds with very short timeout when absent.
// Ref: TestLocatorAssertions2.java#isAttachedWithImpossibleTimeoutNot
func TestIsAttachedImpossibleTimeoutNot(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div id=node>Text content</div>`)
	must.NoError(err)

	// Element "no-such-thing" is absent, so not().isAttached should pass immediately.
	must.NoError(playwright.Expect(page.Locator("no-such-thing")).Not().ToBeAttached(ctx))
}

// ── toHaveAccessibleName extended ────────────────────────────────────────────

// TestToHaveAccessibleNameExtended verifies accessible name with ignoreCase and regex variants.
// Ref: TestLocatorAssertions2.java#toHaveAccessibleName
func TestToHaveAccessibleNameExtended(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div role="button" aria-label="Hello"></div>`)
	must.NoError(err)

	loc := page.Locator("div")
	must.NoError(playwright.Expect(loc).ToHaveAccessibleName(ctx, "Hello"))
	must.NoError(playwright.Expect(loc).Not().ToHaveAccessibleName(ctx, "hello"))
	must.NoError(playwright.Expect(loc).ToHaveAccessibleNameIgnoreCase(ctx, "hello"))
	must.NoError(playwright.Expect(loc).ToHaveAccessibleNameRegex(ctx, regexp.MustCompile(`ell\w`)))
	must.NoError(playwright.Expect(loc).Not().ToHaveAccessibleNameRegex(ctx, regexp.MustCompile(`hello`)))
	must.NoError(playwright.Expect(loc).ToHaveAccessibleNameRegex(ctx, regexp.MustCompile(`(?i)hello`)))

	// Test &nbsp; normalization
	err = page.SetContent(ctx, `<button>foo&nbsp;bar
baz</button>`)
	must.NoError(err)
	must.NoError(playwright.Expect(page.Locator("button")).ToHaveAccessibleName(ctx, "foo bar baz"))
}

// ── toHaveAccessibleDescription extended ────────────────────────────────────

// TestToHaveAccessibleDescriptionExtended verifies accessible description with ignoreCase and regex.
// Ref: TestLocatorAssertions2.java#toHaveAccessibleDescription
func TestToHaveAccessibleDescriptionExtended(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div role="button" aria-description="Hello"></div>`)
	must.NoError(err)

	loc := page.Locator("div")
	must.NoError(playwright.Expect(loc).ToHaveAccessibleDescription(ctx, "Hello"))
	must.NoError(playwright.Expect(loc).Not().ToHaveAccessibleDescription(ctx, "hello"))
	must.NoError(playwright.Expect(loc).ToHaveAccessibleDescriptionIgnoreCase(ctx, "hello"))
	must.NoError(playwright.Expect(loc).ToHaveAccessibleDescriptionRegex(ctx, regexp.MustCompile(`ell\w`)))
	must.NoError(playwright.Expect(loc).Not().ToHaveAccessibleDescriptionRegex(ctx, regexp.MustCompile(`hello`)))
	must.NoError(playwright.Expect(loc).ToHaveAccessibleDescriptionRegex(ctx, regexp.MustCompile(`(?i)hello`)))

	// Test &nbsp; normalization via aria-describedby
	err = page.SetContent(ctx, `<div role="button" aria-describedby="desc"></div>
<span id="desc">foo&nbsp;bar
baz</span>`)
	must.NoError(err)
	must.NoError(playwright.Expect(page.Locator("div")).ToHaveAccessibleDescription(ctx, "foo bar baz"))
}

// ── toHaveRole ───────────────────────────────────────────────────────────────

// TestToHaveRole verifies hasRole assertion.
// Ref: TestLocatorAssertions2.java#toHaveRole
func TestToHaveRole(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div role="button">Button!</div>`)
	must.NoError(err)

	must.NoError(playwright.Expect(page.Locator("div")).ToHaveRole(ctx, "button"))
	must.NoError(playwright.Expect(page.Locator("div")).Not().ToHaveRole(ctx, "checkbox"))
}

// ── toHaveAccessibleErrorMessage ─────────────────────────────────────────────

// TestToHaveAccessibleErrorMessage verifies error message assertion with text/regex/ignoreCase.
// Ref: TestLocatorAssertions2.java#toHaveAccessibleErrorMessage
func TestToHaveAccessibleErrorMessage(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<form>
<input role="textbox" aria-invalid="true" aria-errormessage="error-message" />
<div id="error-message">Hello</div>
<div id="irrelevant-error">This should not be considered.</div>
</form>`)
	must.NoError(err)

	loc := page.Locator(`input[role="textbox"]`)
	must.NoError(playwright.Expect(loc).ToHaveAccessibleErrorMessage(ctx, "Hello"))
	must.NoError(playwright.Expect(loc).Not().ToHaveAccessibleErrorMessage(ctx, "hello"))
	must.NoError(playwright.Expect(loc).ToHaveAccessibleErrorMessageIgnoreCase(ctx, "hello"))
	must.NoError(playwright.Expect(loc).ToHaveAccessibleErrorMessageRegex(ctx, regexp.MustCompile(`ell\w`)))
	must.NoError(playwright.Expect(loc).Not().ToHaveAccessibleErrorMessageRegex(ctx, regexp.MustCompile(`hello`)))
	must.NoError(playwright.Expect(loc).ToHaveAccessibleErrorMessageRegex(ctx, regexp.MustCompile(`(?i)hello`)))
	must.NoError(playwright.Expect(loc).Not().ToHaveAccessibleErrorMessage(ctx, "This should not be considered."))
}

// TestToHaveAccessibleErrorMessageMultiple verifies error message with multiple aria-errormessage refs.
// Ref: TestLocatorAssertions2.java#toHaveAccessibleErrorMessageShouldHandleMultipleAriaErrorMessageReferences
func TestToHaveAccessibleErrorMessageMultiple(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<form>
  <input role="textbox" aria-invalid="true" aria-errormessage="error1 error2" />
  <div id="error1">First error message.</div>
  <div id="error2">Second error message.</div>
  <div id="irrelevant-error">This should not be considered.</div>
</form>`)
	must.NoError(err)

	loc := page.Locator(`input[role="textbox"]`)
	must.NoError(playwright.Expect(loc).ToHaveAccessibleErrorMessage(ctx, "First error message. Second error message."))
	must.NoError(playwright.Expect(loc).ToHaveAccessibleErrorMessageRegex(ctx, regexp.MustCompile(`(?i)first error message.`)))
	must.NoError(playwright.Expect(loc).ToHaveAccessibleErrorMessageRegex(ctx, regexp.MustCompile(`(?i)second error message.`)))
	must.NoError(playwright.Expect(loc).Not().ToHaveAccessibleErrorMessageRegex(ctx, regexp.MustCompile(`(?i)This should not be considered.`)))
}

// ── isChecked indeterminate ──────────────────────────────────────────────────

// TestIsCheckedIndeterminate verifies isChecked with indeterminate=true passes on indeterminate checkbox.
// Ref: TestLocatorAssertions2.java#toBeEditableWithIndeterminateTrue
func TestIsCheckedIndeterminate(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<input type=checkbox></input>`)
	must.NoError(err)

	_, err = page.Locator("input").Evaluate(ctx, `e => e.indeterminate = true`)
	must.NoError(err)

	must.NoError(playwright.Expect(page.Locator("input")).ToBeCheckedIndeterminate(ctx))
}

// TestIsCheckedIndeterminateFail verifies isChecked with indeterminate=true fails on non-indeterminate checkbox.
// Ref: TestLocatorAssertions2.java#toBeEditableFailWithIndeterminateTrue
func TestIsCheckedIndeterminateFail(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := shortCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<input type=checkbox></input>`)
	must.NoError(err)

	err = playwright.Expect(page.Locator("input")).ToBeCheckedIndeterminate(ctx)
	is.Error(err, "non-indeterminate checkbox should fail isChecked(indeterminate=true)")
}

// TestIsCheckedIndeterminateAndChecked verifies that setting both indeterminate=true and
// checked=false simultaneously returns an error from the Playwright server.
// Ref: TestLocatorAssertions2.java#toBeEditableWithIndeterminateTrueAndChecked
func TestIsCheckedIndeterminateAndChecked(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<input type=checkbox></input>`)
	must.NoError(err)

	_, err = page.Locator("input").Evaluate(ctx, `e => e.indeterminate = true`)
	must.NoError(err)

	err = playwright.Expect(page.Locator("input")).ToBeCheckedWithIndeterminateAndChecked(ctx, false)
	is.Error(err, "indeterminate=true with checked=false should fail")
	is.ErrorContains(err, "Can't assert indeterminate and checked at the same time")
}

// TestHasTextFilterMatchesHeadingEx2 verifies HasText filter finds heading.
// Ref: TestLocatorHasText.java#shouldFilterHeading
func TestHasTextFilterMatchesHeadingEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<h2>Introduction</h2>
		<h2>Summary</h2>
		<h2>Conclusion</h2>
	`))

	text := "Summary"
	count, err := page.Locator("h2").Filter(&playwright.LocatorFilterOptions{HasText: &text}).Count(ctx)
	must.NoError(err)
	is.Equal(1, count)
}

// TestHasTextFilterMultipleMatchesEx2 verifies HasText can match multiple elements.
// Ref: TestLocatorHasText.java#shouldFilterMultiple
func TestHasTextFilterMultipleMatchesEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<p>Error: File not found</p>
		<p>Error: Permission denied</p>
		<p>Success: Operation complete</p>
	`))

	text := "Error"
	count, err := page.Locator("p").Filter(&playwright.LocatorFilterOptions{HasText: &text}).Count(ctx)
	must.NoError(err)
	is.Equal(2, count)
}

// TestHasTextFilterGetsTextEx2 verifies HasText filtered text is accessible.
// Ref: TestLocatorHasText.java#shouldGetFilteredText
func TestHasTextFilterGetsTextEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div class="item">Category: Books</div>
		<div class="item">Category: Music</div>
		<div class="item">Category: Movies</div>
	`))

	text := "Music"
	innerText, err := page.Locator(".item").Filter(&playwright.LocatorFilterOptions{HasText: &text}).InnerText(ctx)
	must.NoError(err)
	is.Equal("Category: Music", innerText)
}

// TestHasNotTextFilterEx2 verifies HasNotText excludes matching text.
// Ref: TestLocatorHasText.java#shouldExcludeWithHasNotText
func TestHasNotTextFilterEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<li>Red</li>
		<li>Blue</li>
		<li>Green</li>
		<li>Yellow</li>
	`))

	notText := "Red"
	count, err := page.Locator("li").Filter(&playwright.LocatorFilterOptions{HasNotText: &notText}).Count(ctx)
	must.NoError(err)
	is.Equal(3, count)
}

func localStrPtrHT3(s string) *string { return &s }

// TestHaveTextFilterBadgesEx3 verifies Filter.HasText finds badge items.
// Ref: TestLocatorHaveText.java#shouldFindBadges
func TestHaveTextFilterBadgesEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<span class="badge">Sale</span>
		<span class="badge">New</span>
		<span class="badge">Out of stock</span>
	`))

	count, err := page.Locator(".badge").Filter(&playwright.LocatorFilterOptions{
		HasText: localStrPtrHT3("New"),
	}).Count(ctx)
	must.NoError(err)
	is.Equal(1, count)
}

// TestHaveTextFilterMenuItemsEx3 verifies Filter.HasText finds menu items.
// Ref: TestLocatorHaveText.java#shouldFindMenuItems
func TestHaveTextFilterMenuItemsEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<nav>
			<a class="nav-item" href="/home">Home</a>
			<a class="nav-item" href="/about">About Us</a>
			<a class="nav-item" href="/contact">Contact</a>
		</nav>
	`))

	text, err := page.Locator(".nav-item").Filter(&playwright.LocatorFilterOptions{
		HasText: localStrPtrHT3("About"),
	}).TextContent(ctx)
	must.NoError(err)
	is.Equal("About Us", text)
}

// TestHaveTextFilterStatusEx3 verifies Filter.HasText finds status elements.
// Ref: TestLocatorHaveText.java#shouldFindStatus
func TestHaveTextFilterStatusEx3(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `
		<div class="status">Pending</div>
		<div class="status">Approved</div>
		<div class="status">Rejected</div>
		<div class="status">Approved</div>
	`))

	count, err := page.Locator(".status").Filter(&playwright.LocatorFilterOptions{
		HasText: localStrPtrHT3("Approved"),
	}).Count(ctx)
	must.NoError(err)
	is.Equal(2, count)
}
