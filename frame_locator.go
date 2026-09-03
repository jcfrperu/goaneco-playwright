package playwright

import "fmt"

// FrameLocator represents an iframe element and allows locating elements inside it.
// All locator methods are scoped to the iframe's document.
type FrameLocator struct {
	frame          ChannelOwner
	selector       string // includes the iframe selector + ">> internal:control=enter-frame"
	browserContext *BrowserContext
}

func (f *FrameLocator) locatorWithSelector(innerSelector string) *Locator {
	return &Locator{
		frame:    f.frame,
		selector: fmt.Sprintf("%s >> %s", f.selector, innerSelector),
	}
}

// Locator returns a Locator scoped to this frame, matching elements by CSS/internal selector.
func (f *FrameLocator) Locator(selector string) *Locator {
	return f.locatorWithSelector(selector)
}

// GetByRole returns a Locator scoped to this frame matching elements by ARIA role.
func (f *FrameLocator) GetByRole(role AriaRole, opts ...*GetByRoleOptions) *Locator {
	var options *GetByRoleOptions
	if len(opts) > 0 {
		options = opts[0]
	}
	return f.locatorWithSelector(locatorFromRole(f.frame, role, options).selector)
}

// GetByText returns a Locator scoped to this frame matching elements by visible text.
func (f *FrameLocator) GetByText(text string, opts ...*GetByTextOptions) *Locator {
	exact := false
	if len(opts) > 0 && opts[0] != nil && opts[0].Exact != nil {
		exact = *opts[0].Exact
	}
	return f.locatorWithSelector(locatorFromText(f.frame, text, exact).selector)
}

// GetByLabel returns a Locator scoped to this frame matching form elements by associated label text.
func (f *FrameLocator) GetByLabel(text string, opts ...*GetByLabelOptions) *Locator {
	exact := false
	if len(opts) > 0 && opts[0] != nil && opts[0].Exact != nil {
		exact = *opts[0].Exact
	}
	return f.locatorWithSelector(locatorFromLabel(f.frame, text, exact).selector)
}

// GetByTestId returns a Locator scoped to this frame matching elements by the configured testid attribute.
// The attribute name defaults to "data-testid" and can be changed with BrowserContext.SetTestIdAttribute.
func (f *FrameLocator) GetByTestId(testId string) *Locator {
	attr := "data-testid"
	if f.browserContext != nil {
		attr = f.browserContext.TestIdAttributeName()
	}
	escaped := escapeForAttributeSelector(testId)
	return f.locatorWithSelector(fmt.Sprintf(`internal:testid=[%s=%s]`, attr, escaped))
}

// FrameLocator returns a nested FrameLocator for locating elements within an iframe inside this frame.
func (f *FrameLocator) FrameLocator(selector string) *FrameLocator {
	return &FrameLocator{
		frame:          f.frame,
		selector:       fmt.Sprintf("%s >> %s >> internal:control=enter-frame", f.selector, selector),
		browserContext: f.browserContext,
	}
}
