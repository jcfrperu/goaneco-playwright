//go:build e2e

// E2E tests for APIResponseAssertions.
// Migration of: TestAPIResponseAssertions.java
package e2e

import (
	"io"
	"net/http"
	"testing"

	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAPIResponseAssertionsPassWithResponse verifies that ToBeOK passes for a 200 response.
// Ref: TestAPIResponseAssertions.java#passWithResponse
func TestAPIResponseAssertionsPassWithResponse(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	apiCtx := newAPICtx(t)
	resp, err := apiCtx.Get(ctx, srv.EmptyPage())
	must.NoError(err)
	defer resp.Dispose(ctx) //nolint:errcheck

	must.NoError(playwright.ExpectAPIResponse(resp).ToBeOK(ctx))
}

// TestAPIResponseAssertionsPassWithNot verifies that Not().ToBeOK passes for a non-200 response.
// Ref: TestAPIResponseAssertions.java#passWithNot
func TestAPIResponseAssertionsPassWithNot(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	apiCtx := newAPICtx(t)
	resp, err := apiCtx.Get(ctx, srv.Prefix()+"/unknown")
	must.NoError(err)
	defer resp.Dispose(ctx) //nolint:errcheck

	must.NoError(playwright.ExpectAPIResponse(resp).Not().ToBeOK(ctx))
}

// TestAPIResponseAssertionsFail verifies that ToBeOK fails for a 404 response and the error
// message includes the request method/URL and response status/statusText.
// Ref: TestAPIResponseAssertions.java#fail
func TestAPIResponseAssertionsFail(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	apiCtx := newAPICtx(t)
	resp, err := apiCtx.Get(ctx, srv.Prefix()+"/unknown")
	must.NoError(err)
	defer resp.Dispose(ctx) //nolint:errcheck

	assertErr := playwright.ExpectAPIResponse(resp).ToBeOK(ctx)
	is.Error(assertErr)
	is.Contains(assertErr.Error(), "→ GET "+srv.Prefix()+"/unknown")
	is.Contains(assertErr.Error(), "← 404 Not Found")
}

// TestAPIResponseAssertionsShouldPrintResponseTextIfIsOkFails verifies that the error message
// includes the response body when the content-type is textual (text/plain; charset=utf-8).
// Ref: TestAPIResponseAssertions.java#shouldPrintResponseTextIfIdOkFails
func TestAPIResponseAssertionsShouldPrintResponseTextIfIsOkFails(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	apiCtx := newAPICtx(t)
	resp, err := apiCtx.Get(ctx, srv.Prefix()+"/unknown")
	must.NoError(err)
	defer resp.Dispose(ctx) //nolint:errcheck

	assertErr := playwright.ExpectAPIResponse(resp).ToBeOK(ctx)
	is.Error(assertErr)
	// Go's http.NotFound returns "404 page not found\n" with Content-Type text/plain; charset=utf-8.
	is.Contains(assertErr.Error(), "404 page not found")
}

// TestAPIResponseAssertionsShouldOnlyPrintResponseWithTextContentTypeIfIsOkFails verifies that
// the error message includes the body only for textual content types (text/*, image/svg+xml, etc.)
// and not for binary or missing content-types.
// Ref: TestAPIResponseAssertions.java#shouldOnlyPrintResponseWithTextContentTypeIfIsOkFails
func TestAPIResponseAssertionsShouldOnlyPrintResponseWithTextContentTypeIfIsOkFails(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	apiCtx := newAPICtx(t)

	// text/plain → body must be included in the error message
	srv.SetRoute("/text-content-type", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(404)
		io.WriteString(w, "Text error") //nolint:errcheck
	})
	resp1, err := apiCtx.Get(ctx, srv.Prefix()+"/text-content-type")
	must.NoError(err)
	err1 := playwright.ExpectAPIResponse(resp1).ToBeOK(ctx)
	resp1.Dispose(ctx) //nolint:errcheck
	is.Error(err1)
	is.Contains(err1.Error(), "Text error")

	// image/svg+xml is a textual MIME type → body must be included
	srv.SetRoute("/svg-xml-content-type", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/svg+xml")
		w.WriteHeader(404)
		io.WriteString(w, "Json error") //nolint:errcheck
	})
	resp2, err := apiCtx.Get(ctx, srv.Prefix()+"/svg-xml-content-type")
	must.NoError(err)
	err2 := playwright.ExpectAPIResponse(resp2).ToBeOK(ctx)
	resp2.Dispose(ctx) //nolint:errcheck
	is.Error(err2)
	is.Contains(err2.Error(), "Json error")

	// no content-type → body must NOT be included
	srv.SetRoute("/no-content-type", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		io.WriteString(w, "No content type error") //nolint:errcheck
	})
	resp3, err := apiCtx.Get(ctx, srv.Prefix()+"/no-content-type")
	must.NoError(err)
	err3 := playwright.ExpectAPIResponse(resp3).ToBeOK(ctx)
	resp3.Dispose(ctx) //nolint:errcheck
	is.Error(err3)
	is.NotContains(err3.Error(), "No content type error")

	// image/bmp is a binary type → body must NOT be included
	srv.SetRoute("/image-content-type", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/bmp")
		w.WriteHeader(404)
		io.WriteString(w, "Image type error") //nolint:errcheck
	})
	resp4, err := apiCtx.Get(ctx, srv.Prefix()+"/image-content-type")
	must.NoError(err)
	err4 := playwright.ExpectAPIResponse(resp4).ToBeOK(ctx)
	resp4.Dispose(ctx) //nolint:errcheck
	is.Error(err4)
	is.NotContains(err4.Error(), "Image type error")
}
