//go:build e2e

// E2E tests for: TestPageRequestIntercept.java.
package e2e

import (
	"testing"

	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/require"
)

// TestPageRequestInterceptShouldFulfillPopupMainRequestUsingAlias is not ported yet.
// Ref: TestPageRequestIntercept.java#shouldFulfillPopupMainRequestUsingAlias  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageRequestIntercept.java
func TestPageRequestInterceptShouldFulfillPopupMainRequestUsingAlias(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Context().Route(ctx, "**/*", func(route *playwright.Route) {
		resp, err := route.Fetch(ctx)
		if err != nil {
			_ = route.Abort(ctx)
			return
		}
		hello := "hello"
		_ = route.Fulfill(ctx, &playwright.RouteFulfillOptions{
			Response: resp,
			Body:     &hello,
		})
	}))

	must.NoError(page.SetContent(ctx, "<a target=_blank href='"+srv.EmptyPage()+"'>click me</a>"))
	popup, err := page.WaitForPopup(ctx, func() error {
		return page.GetByText("click me").Click(ctx)
	})
	must.NoError(err)
	must.NoError(playwright.Expect(popup.Locator("body")).ToHaveText(ctx, "hello"))
}
