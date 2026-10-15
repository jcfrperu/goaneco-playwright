package playwright

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResolveUnsetTimeouts(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	for _, g := range []string{"context@t1", "page@t1", "frame@t1"} {
		timeoutScopes.Delete(g)
	}
	page := ChannelOwner{guid: "page@t1"}
	frame := ChannelOwner{guid: "frame@t1"}
	ctxScope := ensureTimeoutScope("context@t1", "")
	pageScope := ensureTimeoutScope(page.guid, "context@t1")
	ensureTimeoutScope(frame.guid, page.guid)

	// Nothing configured: Playwright's 30s fallback.
	is.Equal(30000.0, frame.effectiveTimeout("click"))

	// Context-level default applies through page and frame.
	ctxScope.setTimeout(1000)
	is.Equal(1000.0, frame.effectiveTimeout("click"))
	is.Equal(1000.0, frame.effectiveTimeout("goto"))

	// Page-level overrides context; navigation timeout only affects navigations.
	pageScope.setTimeout(2000)
	pageScope.setNavTimeout(3000)
	is.Equal(2000.0, frame.effectiveTimeout("click"))
	is.Equal(3000.0, frame.effectiveTimeout("goto"))

	// The sentinel is replaced, explicit values are untouched.
	out := string(frame.resolveUnsetTimeouts("click", []byte(`{"selector":"a","timeout":-1,"strict":true}`)))
	is.Equal(`{"selector":"a","timeout":2000,"strict":true}`, out)
	out = string(frame.resolveUnsetTimeouts("goto", []byte(`{"url":"x","timeout":-1}`)))
	is.Equal(`{"url":"x","timeout":3000}`, out)
	out = string(frame.resolveUnsetTimeouts("click", []byte(`{"timeout":5000}`)))
	is.Equal(`{"timeout":5000}`, out)
}
