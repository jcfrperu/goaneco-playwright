//go:build e2e

// E2E tests for: TestPageEventRequest.java.
package e2e

import (
	"fmt"
	"net/http"
	"testing"

	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Ref: TestPageEventRequest.java#shouldReturnLastRequests  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageEventRequest.java
func TestPageEventRequestShouldReturnLastRequests(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	// Set up routes for 200 requests
	for i := 0; i < 200; i++ {
		path := fmt.Sprintf("/fetch-%d", i)
		srv.SetRoute(path, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
			_, _ = w.Write([]byte("url:" + srv.Prefix() + r.URL.RequestURI()))
		})
	}

	must.NoError(page.Goto(ctx, srv.Prefix()+"/title.html"))

	// #0 is the navigation request, so start with #1.
	for i := 0; i < 99; i++ {
		_, err := page.Evaluate(ctx, "url => fetch(url)", srv.Prefix()+fmt.Sprintf("/fetch-%d", i))
		must.NoError(err)
	}

	first99WithNav, err := page.Requests(ctx)
	must.NoError(err)
	must.GreaterOrEqual(len(first99WithNav), 100)
	// Remove navigation request (first element)
	first99Requests := first99WithNav[1:]

	for i := 99; i < 199; i++ {
		_, err := page.Evaluate(ctx, "url => fetch(url)", srv.Prefix()+fmt.Sprintf("/fetch-%d", i))
		must.NoError(err)
	}

	last100Requests, err := page.Requests(ctx)
	must.NoError(err)

	allRequests := make([]*playwright.NetworkRequest, 0, len(first99Requests)+len(last100Requests))
	allRequests = append(allRequests, first99Requests...)
	allRequests = append(allRequests, last100Requests...)

	// All 199 requests are fully functional.
	for index, request := range allRequests {
		response, err := request.Response(ctx)
		must.NoError(err)
		must.NotNil(response)
		text, err := response.Text(ctx)
		must.NoError(err)
		expectedURL := srv.Prefix() + fmt.Sprintf("/fetch-%d", index)
		is.Equal("url:"+expectedURL, text)
		is.Equal(expectedURL, request.URL())
	}
}
