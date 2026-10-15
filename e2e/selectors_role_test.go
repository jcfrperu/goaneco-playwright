//go:build e2e

// E2E tests for: TestSelectorsRole.java.
// Tests marked "Extra:" are additional tests with no equivalent method in playwright-java.
package e2e

import (
	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

// TestRoleSelectorFilterHidden verifies includeHidden option.
// Extra: no equivalent method in playwright-java; related class: TestSelectorsRole.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsRole.java
func TestRoleSelectorFilterHidden(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<button>Visible</button>
	<button style="display:none">Hidden</button>`)
	must.NoError(err)

	// By default, hidden elements are not matched
	visible, err := page.GetByRole(playwright.AriaRoleButton).EvaluateAll(ctx, "els => els.length")
	must.NoError(err)
	is.Equal(float64(1), visible)

	// With includeHidden=true, hidden elements are matched
	all, err := page.GetByRole(playwright.AriaRoleButton, &playwright.GetByRoleOptions{
		IncludeHidden: roleBoolPtr(true),
	}).EvaluateAll(ctx, "els => els.length")
	must.NoError(err)
	is.Equal(float64(2), all)
}

// TestSelectorsRoleShouldInheritDisabledFromTheAncestor verifies that buttons inherit aria-disabled from
// ancestor elements, but non-control roles like h1 do not.
// Ref: TestSelectorsRole.java#shouldInheritDisabledFromTheAncestor  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsRole.java
func TestSelectorsRoleShouldInheritDisabledFromTheAncestor(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<span aria-disabled="true"><button>Click me!</button></span>`)
	must.NoError(err)

	disabled, err := page.Locator("button").IsDisabled(ctx)
	must.NoError(err)
	is.True(disabled)

	err = page.SetContent(ctx, `<span aria-disabled="true"><h1>Heading</h1></span>`)
	must.NoError(err)

	enabled, err := page.Locator("h1").IsEnabled(ctx)
	must.NoError(err)
	is.True(enabled)
}

// TestSelectorsRoleShouldSupportDisabledFieldset verifies disabled fieldset behaviour: only the first legend
// element exempts its descendants from the disabled state.
// Ref: TestSelectorsRole.java#shouldSupportDisabledFieldset  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsRole.java
func TestSelectorsRoleShouldSupportDisabledFieldset(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `
<fieldset disabled>
  <input></input>
  <button data-testid="inside-fieldset-element">x</button>
  <legend>
    <button data-testid="inside-legend-element">legend</button>
  </legend>
</fieldset>

<fieldset disabled>
  <legend>
    <div>
      <button data-testid="nested-inside-legend-element">x</button>
    </div>
  </legend>
</fieldset>

<fieldset disabled>
  <div></div>
  <legend>
    <button data-testid="first-legend-element">x</button>
  </legend>
  <legend>
    <button data-testid="second-legend-element">x</button>
  </legend>
</fieldset>

<fieldset disabled>
  <fieldset>
    <button data-testid="deep-button">x</button>
  </fieldset>
</fieldset>
`)
	must.NoError(err)

	insideLegend, err := page.GetByTestId("inside-legend-element").IsEnabled(ctx)
	must.NoError(err)
	is.True(insideLegend)

	nestedInsideLegend, err := page.GetByTestId("nested-inside-legend-element").IsEnabled(ctx)
	must.NoError(err)
	is.True(nestedInsideLegend)

	firstLegend, err := page.GetByTestId("first-legend-element").IsEnabled(ctx)
	must.NoError(err)
	is.True(firstLegend)

	secondLegend, err := page.GetByTestId("second-legend-element").IsDisabled(ctx)
	must.NoError(err)
	is.True(secondLegend)

	deepButton, err := page.GetByTestId("deep-button").IsDisabled(ctx)
	must.NoError(err)
	is.True(deepButton)
}

// TestSelectorsRoleErrors verifies that malformed role selector syntax produces descriptive error messages.
// Ref: TestSelectorsRole.java#errors  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsRole.java
func TestSelectorsRoleErrors(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	_, err := page.QuerySelector(ctx, "role=[bar]")
	is.Error(err)
	is.Contains(err.Error(), "Role must not be empty")

	_, err = page.QuerySelector(ctx, "role=foo[sElected]")
	is.Error(err)
	is.Contains(err.Error(), "Unknown attribute \"sElected\"")

	_, err = page.QuerySelector(ctx, "role=foo[bar . qux=true]")
	is.Error(err)
	is.Contains(err.Error(), "Unknown attribute")

	_, err = page.QuerySelector(ctx, "role=heading[level='bar']")
	is.Error(err)
	is.Contains(err.Error(), "\"level\" attribute must be compared to a number")

	_, err = page.QuerySelector(ctx, "role=checkbox[checked='bar']")
	is.Error(err)
	is.Contains(err.Error(), "\"checked\" must be one of true, false")

	_, err = page.QuerySelector(ctx, "role=checkbox[checked~=true]")
	is.Error(err)
	is.Contains(err.Error(), "cannot use ~=")

	_, err = page.QuerySelector(ctx, "role=button[level=3]")
	is.Error(err)
	is.Contains(err.Error(), "\"level\" attribute is only supported for roles")

	_, err = page.QuerySelector(ctx, "role=button[name]")
	is.Error(err)
	is.Contains(err.Error(), "\"name\" attribute must have a value")
}

// TestSelectorsRoleShouldDetectRoles verifies that role selectors match expected elements.
// Ref: TestSelectorsRole.java#shouldDetectRoles  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsRole.java
func TestSelectorsRoleShouldDetectRoles(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<button>Hello</button>
	<select multiple size="2"></select>
	<select></select>
	<h3>Heading</h3>
	<details><summary>Hello</summary></details>
	<div role="dialog">I am a dialog</div>`)
	must.NoError(err)

	buttons, err := page.Locator("role=button").EvaluateAll(ctx, "els => els.map(e => e.tagName)")
	must.NoError(err)
	is.Equal([]any{"BUTTON"}, buttons)

	headings, err := page.GetByRole(playwright.AriaRoleHeading).EvaluateAll(ctx, "els => els.map(e => e.tagName)")
	must.NoError(err)
	is.Equal([]any{"H3"}, headings)

	dialogs, err := page.GetByRole(playwright.AriaRoleDialog).EvaluateAll(ctx, "els => els.map(e => e.getAttribute('role'))")
	must.NoError(err)
	is.Equal([]any{"dialog"}, dialogs)

	menuitems, err := page.GetByRole(playwright.AriaRoleMenuItem).EvaluateAll(ctx, "els => els.length")
	must.NoError(err)
	is.Equal(float64(0), menuitems)
}

// TestSelectorsRoleShouldSupportChecked verifies [checked] filter on role=checkbox.
// Ref: TestSelectorsRole.java#shouldSupportChecked  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsRole.java
func TestSelectorsRoleShouldSupportChecked(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<input type="checkbox">
	<input type="checkbox" checked>
	<div role="checkbox" aria-checked="true">Hi</div>
	<div role="checkbox" aria-checked="false">Hello</div>`)
	must.NoError(err)

	checked, err := page.GetByRole(playwright.AriaRoleCheckbox, &playwright.GetByRoleOptions{
		Checked: roleBoolPtr(true),
	}).EvaluateAll(ctx, "els => els.length")
	must.NoError(err)
	is.Equal(float64(2), checked)

	unchecked, err := page.GetByRole(playwright.AriaRoleCheckbox, &playwright.GetByRoleOptions{
		Checked: roleBoolPtr(false),
	}).EvaluateAll(ctx, "els => els.length")
	must.NoError(err)
	is.Equal(float64(2), unchecked)
}

// TestSelectorsRoleShouldSupportDisabled verifies [disabled] filter.
// Ref: TestSelectorsRole.java#shouldSupportDisabled  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsRole.java
func TestSelectorsRoleShouldSupportDisabled(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<button>Enabled</button>
	<button disabled>Disabled</button>
	<div role="button" aria-disabled="true">AriaDisabled</div>`)
	must.NoError(err)

	disabled, err := page.GetByRole(playwright.AriaRoleButton, &playwright.GetByRoleOptions{
		Disabled: roleBoolPtr(true),
	}).EvaluateAll(ctx, "els => els.length")
	must.NoError(err)
	is.Equal(float64(2), disabled)

	enabled, err := page.GetByRole(playwright.AriaRoleButton, &playwright.GetByRoleOptions{
		Disabled: roleBoolPtr(false),
	}).EvaluateAll(ctx, "els => els.length")
	must.NoError(err)
	is.Equal(float64(1), enabled)
}

// TestSelectorsRoleShouldSupportName verifies [name] filter on role selectors.
// Ref: TestSelectorsRole.java#shouldSupportName  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsRole.java
func TestSelectorsRoleShouldSupportName(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<button>Hello</button><button>World</button>`)
	must.NoError(err)

	hello, err := page.GetByRole(playwright.AriaRoleButton, &playwright.GetByRoleOptions{
		Name: roleStrPtr("Hello"),
	}).EvaluateAll(ctx, "els => els.map(e => e.textContent)")
	must.NoError(err)
	is.Equal([]any{"Hello"}, hello)

	world, err := page.GetByRole(playwright.AriaRoleButton, &playwright.GetByRoleOptions{
		Name: roleStrPtr("World"),
	}).EvaluateAll(ctx, "els => els.map(e => e.textContent)")
	must.NoError(err)
	is.Equal([]any{"World"}, world)
}

// TestSelectorsRoleShouldSupportLevel verifies heading level filter.
// Ref: TestSelectorsRole.java#shouldSupportLevel  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsRole.java
func TestSelectorsRoleShouldSupportLevel(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<h1>H1</h1><h2>H2</h2><h3>H3</h3>`)
	must.NoError(err)

	lvl1 := 1
	h1, err := page.GetByRole(playwright.AriaRoleHeading, &playwright.GetByRoleOptions{
		Level: &lvl1,
	}).EvaluateAll(ctx, "els => els.map(e => e.textContent)")
	must.NoError(err)
	is.Equal([]any{"H1"}, h1)

	lvl2 := 2
	h2, err := page.GetByRole(playwright.AriaRoleHeading, &playwright.GetByRoleOptions{
		Level: &lvl2,
	}).EvaluateAll(ctx, "els => els.map(e => e.textContent)")
	must.NoError(err)
	is.Equal([]any{"H2"}, h2)
}

// TestSelectorsRoleShouldSupportExpanded verifies [expanded] filter.
// Ref: TestSelectorsRole.java#shouldSupportExpanded  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsRole.java
func TestSelectorsRoleShouldSupportExpanded(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<div role="treeitem">Hi</div>
	<div role="treeitem" aria-expanded="true">Hello</div>
	<div role="treeitem" aria-expanded="false">Bye</div>`)
	must.NoError(err)

	expanded, err := page.GetByRole(playwright.AriaRoleTreeItem, &playwright.GetByRoleOptions{
		Expanded: roleBoolPtr(true),
	}).EvaluateAll(ctx, "els => els.map(e => e.textContent)")
	must.NoError(err)
	is.Equal([]any{"Hello"}, expanded)

	collapsed, err := page.GetByRole(playwright.AriaRoleTreeItem, &playwright.GetByRoleOptions{
		Expanded: roleBoolPtr(false),
	}).EvaluateAll(ctx, "els => els.map(e => e.textContent)")
	must.NoError(err)
	is.Equal([]any{"Bye"}, collapsed)
}

// TestSelectorsRoleShouldSupportSelected verifies [selected] filter on role=option.
// Ref: TestSelectorsRole.java#shouldSupportSelected  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsRole.java
func TestSelectorsRoleShouldSupportSelected(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<select>
		<option>Hi</option>
		<option selected>Hello</option>
	</select>
	<div>
		<div role="option" aria-selected="true">Hi</div>
		<div role="option" aria-selected="false">Hello</div>
	</div>`)
	must.NoError(err)

	count, err := page.GetByRole(playwright.AriaRoleOption, &playwright.GetByRoleOptions{
		Selected: roleBoolPtr(true),
	}).EvaluateAll(ctx, "els => els.length")
	must.NoError(err)
	is.Equal(float64(2), count)
}

// TestSelectorsRoleShouldSupportPressed verifies [pressed] filter on role=button.
// Ref: TestSelectorsRole.java#shouldSupportPressed  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsRole.java
func TestSelectorsRoleShouldSupportPressed(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<button>Hi</button>
	<button aria-pressed="true">Hello</button>
	<button aria-pressed="false">Bye</button>`)
	must.NoError(err)

	pressed, err := page.GetByRole(playwright.AriaRoleButton, &playwright.GetByRoleOptions{
		Pressed: roleBoolPtr(true),
	}).EvaluateAll(ctx, "els => els.map(e => e.textContent)")
	must.NoError(err)
	is.Equal([]any{"Hello"}, pressed)

	notPressed, err := page.GetByRole(playwright.AriaRoleButton, &playwright.GetByRoleOptions{
		Pressed: roleBoolPtr(false),
	}).EvaluateAll(ctx, "els => els.length")
	must.NoError(err)
	is.Equal(float64(2), notPressed)
}

// Ref: TestSelectorsRole.java#shouldFilterHiddenUnlessExplicitlyAskedFor  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestSelectorsRole.java
func TestSelectorsRoleShouldFilterHiddenUnlessExplicitlyAskedFor(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, "<button>Hi</button>\n"+
		"    <button hidden>Hello</button>\n"+
		"    <button aria-hidden=\"true\">Yay</button>\n"+
		"    <button aria-hidden=\"false\">Nay</button>\n"+
		"    <button style=\"visibility:hidden\">Bye</button>\n"+
		"    <div style=\"visibility:hidden\">\n"+
		"      <button>Oh</button>\n"+
		"    </div>\n"+
		"    <div style=\"visibility:hidden\">\n"+
		"      <button style=\"visibility:visible\">Still here</button>\n"+
		"    </div>\n"+
		"    <button style=\"display:none\">Never</button>\n"+
		"    <div id=host1></div>\n"+
		"    <div id=host2 style=\"display:none\"></div>\n"+
		"    <script>\n"+
		"      function addButton(host, text) {\n"+
		"        const root = host.attachShadow({ mode: 'open' });\n"+
		"        const button = document.createElement('button');\n"+
		"        button.textContent = text;\n"+
		"        root.appendChild(button);\n"+
		"      }\n"+
		"      addButton(document.getElementById('host1'), 'Shadow1');\n"+
		"      addButton(document.getElementById('host2'), 'Shadow2');\n"+
		"    </script>"))

	visible := []any{
		"<button>Hi</button>",
		"<button aria-hidden=\"false\">Nay</button>",
		"<button style=\"visibility:visible\">Still here</button>",
		"<button>Shadow1</button>",
	}
	all := []any{
		"<button>Hi</button>",
		"<button hidden=\"\">Hello</button>",
		"<button aria-hidden=\"true\">Yay</button>",
		"<button aria-hidden=\"false\">Nay</button>",
		"<button style=\"visibility:hidden\">Bye</button>",
		"<button>Oh</button>",
		"<button style=\"visibility:visible\">Still here</button>",
		"<button style=\"display:none\">Never</button>",
		"<button>Shadow1</button>",
		"<button>Shadow2</button>",
	}
	check := func(selector string, expected []any) {
		res, err := page.Locator(selector).EvaluateAll(ctx, "els => els.map(e => e.outerHTML)")
		must.NoError(err, selector)
		assert.Equal(t, expected, res, selector)
	}
	check("role=button", visible)
	check("role=button[include-hidden]", all)
	check("role=button[include-hidden=true]", all)
	check("role=button[include-hidden=false]", visible)
}
