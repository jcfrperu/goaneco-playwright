package playwright

import "fmt"

// FrameLocator represents an iframe element and allows locating elements inside it.
// All locator methods are scoped to the iframe's document.
type FrameLocator struct {
	frame          ChannelOwner
	frameSelector  string // selector of the iframe itself
	selector       string // frameSelector + " >> internal:control=enter-frame"
	browserContext *BrowserContext
}

func newFrameLocator(frame ChannelOwner, frameSelector string, bc ...*BrowserContext) *FrameLocator {
	var b *BrowserContext
	if len(bc) > 0 {
		b = bc[0]
	}
	return &FrameLocator{
		frame:          frame,
		frameSelector:  frameSelector,
		selector:       fmt.Sprintf("%s >> internal:control=enter-frame", frameSelector),
		browserContext: b,
	}
}

func (f *FrameLocator) locatorWithSelector(innerSelector string) *Locator {
	return &Locator{
		frame:    f.frame,
		selector: fmt.Sprintf("%s >> %s", f.selector, innerSelector),
	}
}

// Owner returns a Locator that points to the frame element itself.
func (f *FrameLocator) Owner() *Locator {
	return &Locator{
		frame:    f.frame,
		selector: f.frameSelector,
	}
}

// First returns a FrameLocator pointing to the first matching frame.
func (f *FrameLocator) First() *FrameLocator {
	return f.Nth(0)
}

// Last returns a FrameLocator pointing to the last matching frame.
func (f *FrameLocator) Last() *FrameLocator {
	return f.Nth(-1)
}

// Nth returns a FrameLocator pointing to the nth matching frame (0-indexed; -1 for last).
func (f *FrameLocator) Nth(index int) *FrameLocator {
	var sel string
	if index == -1 {
		sel = fmt.Sprintf("%s >> nth=-1", f.frameSelector)
	} else {
		sel = fmt.Sprintf("%s >> nth=%d", f.frameSelector, index)
	}
	return &FrameLocator{
		frame:          f.frame,
		frameSelector:  sel,
		selector:       fmt.Sprintf("%s >> internal:control=enter-frame", sel),
		browserContext: f.browserContext,
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

// GetByPlaceholder returns a Locator scoped to this frame matching input elements by placeholder text.
func (f *FrameLocator) GetByPlaceholder(text string, opts ...*GetByPlaceholderOptions) *Locator {
	exact := false
	if len(opts) > 0 && opts[0] != nil && opts[0].Exact != nil {
		exact = *opts[0].Exact
	}
	return f.locatorWithSelector(locatorFromAttr(f.frame, "placeholder", text, exact).selector)
}

// GetByAltText returns a Locator scoped to this frame matching elements by alt text.
func (f *FrameLocator) GetByAltText(text string, opts ...*GetByAltTextOptions) *Locator {
	exact := false
	if len(opts) > 0 && opts[0] != nil && opts[0].Exact != nil {
		exact = *opts[0].Exact
	}
	return f.locatorWithSelector(locatorFromAttr(f.frame, "alt", text, exact).selector)
}

// GetByTitle returns a Locator scoped to this frame matching elements by title attribute.
func (f *FrameLocator) GetByTitle(text string, opts ...*GetByTitleOptions) *Locator {
	exact := false
	if len(opts) > 0 && opts[0] != nil && opts[0].Exact != nil {
		exact = *opts[0].Exact
	}
	return f.locatorWithSelector(locatorFromAttr(f.frame, "title", text, exact).selector)
}

// GetByTestId returns a Locator scoped to this frame matching elements by the configured testid attribute.
// The attribute name defaults to "data-testid" and can be changed with BrowserContext.SetTestIdAttribute.
func (f *FrameLocator) GetByTestId(testId any) *Locator {
	attr := "data-testid"
	if f.browserContext != nil {
		attr = f.browserContext.TestIdAttributeName()
	}
	return f.locatorWithSelector(getByTestIdSelector(attr, testId))
}

// FrameLocator returns a nested FrameLocator for locating elements within an iframe inside this frame.
func (f *FrameLocator) FrameLocator(selector string) *FrameLocator {
	return newFrameLocator(f.frame, fmt.Sprintf("%s >> %s", f.selector, selector), f.browserContext)
}
