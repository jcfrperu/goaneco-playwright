//go:build e2e

// Ref: TestSelectorsRole.java

package e2e

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestRoleSelectorInheritDisabledFromAncestor verifies that buttons inherit aria-disabled from
// ancestor elements, but non-control roles like h1 do not.
// Ref: TestSelectorsRole.java#shouldInheritDisabledFromTheAncestor
func TestRoleSelectorInheritDisabledFromAncestor(t *testing.T) {
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

// TestRoleSelectorDisabledFieldset verifies disabled fieldset behaviour: only the first legend
// element exempts its descendants from the disabled state.
// Ref: TestSelectorsRole.java#shouldSupportDisabledFieldset
func TestRoleSelectorDisabledFieldset(t *testing.T) {
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

// TestRoleSelectorErrors verifies that malformed role selector syntax produces descriptive error messages.
// Ref: TestSelectorsRole.java#errors
func TestRoleSelectorErrors(t *testing.T) {
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
