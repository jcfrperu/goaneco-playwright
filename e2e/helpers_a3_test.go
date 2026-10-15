//go:build e2e

// Shared helpers for the ported playwright-java tests (agent 3 batch).
package e2e

import (
	"context"
	"math"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	playwright "github.com/jcfrperu/goaneco-playwright"
)

// a3IsWebKit reports whether the page runs on WebKit.
func a3IsWebKit(page *playwright.Page) bool {
	return page.Context().Browser().BrowserType().Name() == "webkit"
}

// a3IsFirefox reports whether the page runs on Firefox.
func a3IsFirefox(page *playwright.Page) bool {
	return page.Context().Browser().BrowserType().Name() == "firefox"
}

// a3InputValue returns the value of the first <input> of the page.
func a3InputValue(ctx context.Context, page *playwright.Page) (any, error) {
	return page.EvalOnSelector(ctx, "input", "input => input.value")
}

// a3GiveItAChanceToResolve mirrors the Java helper with the same name.
func a3GiveItAChanceToResolve(ctx context.Context, page *playwright.Page) error {
	for i := 0; i < 5; i++ {
		if _, err := page.Evaluate(ctx, "() => new Promise(f => requestAnimationFrame(() => requestAnimationFrame(f)))"); err != nil {
			return err
		}
	}
	return nil
}

// a3Unshift removes the common leading whitespace and blank lines (TestPageAriaSnapshot.unshift).
func a3Unshift(snapshot string) string {
	lines := strings.Split(snapshot, "\n")
	prefix := 100
	re := regexp.MustCompile(`^(\s*).*`)
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		if m := re.FindStringSubmatch(line); m != nil && len(m[1]) < prefix {
			prefix = len(m[1])
		}
	}
	var out []string
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		out = append(out, line[prefix:])
	}
	return strings.Join(out, "\n")
}

// a3CheckAndMatchSnapshot mirrors TestPageAriaSnapshot.checkAndMatchSnapshot.
func a3CheckAndMatchSnapshot(ctx context.Context, t *testing.T, locator *playwright.Locator, snapshot string) {
	t.Helper()
	got, err := locator.AriaSnapshot(ctx)
	require.NoError(t, err)
	assert.Equal(t, a3Unshift(snapshot), got)
	require.NoError(t, playwright.Expect(locator).ToMatchAriaSnapshot(ctx, snapshot))
}

// a3AriaBody sets the content on a fresh page and returns the body locator.
func a3AriaBody(ctx context.Context, t *testing.T, html string) (*playwright.Page, *playwright.Locator) {
	t.Helper()
	page := newPage(t)
	require.NoError(t, page.SetContent(ctx, html))
	return page, page.Locator("body")
}

// a3NormalizeQuotes turns single quotes into double quotes (newer drivers format selectors in error messages with single quotes).
func a3NormalizeQuotes(s string) string {
	return strings.ReplaceAll(s, "'", "\"")
}

// a3Num converts a numeric value returned by Evaluate into float64.
func a3Num(v any) float64 {
	switch n := v.(type) {
	case float64:
		return n
	case float32:
		return float64(n)
	case int:
		return float64(n)
	case int64:
		return float64(n)
	case int32:
		return float64(n)
	}
	return math.NaN()
}

// a3Round mirrors Math.round of the Java tests.
func a3Round(v float64) int64 {
	return int64(math.Round(v))
}

// a3GridHTML is a stand-in for the grid.html fixture: 10 boxes of 50x50 per row in a 500px viewport.
const a3GridHTML = `<style>
  body { margin: 0; }
  .box { float: left; width: 50px; height: 50px; box-sizing: border-box; }
</style>
<body>
<script>
  for (let i = 0; i < 100; i++) {
    const div = document.createElement('div');
    div.className = 'box';
    div.textContent = 'box ' + i;
    document.body.appendChild(div);
  }
</script>`

// a3SelectHTML is a stand-in for the input/select.html fixture.
const a3SelectHTML = `<select>
  <option value="blue">Blue</option>
  <option value="green">Green</option>
  <option value="red">Red</option>
</select>
<script>
  window.result = { onInput: null, onChange: null };
  const select = document.querySelector('select');
  select.addEventListener('input', () => {
    window.result.onInput = Array.from(select.selectedOptions).map(o => o.value);
  }, false);
  select.addEventListener('change', () => {
    window.result.onChange = Array.from(select.selectedOptions).map(o => o.value);
  }, false);
</script>`

// a3ScrollableHTML is a stand-in for the input/scrollable.html fixture.
const a3ScrollableHTML = `<body>
<script>
  for (let i = 0; i < 100; i++) {
    const button = document.createElement('button');
    button.id = 'button-' + i;
    button.textContent = 'Button ' + i;
    document.body.appendChild(button);
    document.body.appendChild(document.createElement('br'));
  }
</script>`

// a3InputResultHTML is a stand-in for the input/textarea.html fixture with an <input> element.
const a3InputResultHTML = `<input type="text">
<script>
  window.result = '';
  const input = document.querySelector('input');
  input.addEventListener('input', () => window.result = input.value, false);
</script>`

// a3AtomicSelector is a selector engine that modifies the found elements right after querying them.
const a3AtomicSelector = `{
  query(root, selector) {
    const result = root.querySelector(selector);
    if (result)
      Promise.resolve().then(() => result.textContent = 'modified');
    return result;
  },
  queryAll(root, selector) {
    const result = Array.from(root.querySelectorAll(selector));
    for (const e of result)
      Promise.resolve().then(() => e.textContent = 'modified');
    return result;
  }
}`
