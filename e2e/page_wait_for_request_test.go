//go:build e2e

// E2E tests for: TestPageWaitForRequest.java.
package e2e

import (
	"context"
	"strings"
	"testing"
	"time"

	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Ref: TestPageWaitForRequest.java#shouldWork  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageWaitForRequest.java
func TestPageWaitForRequestShouldWork(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)
	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	resultCh := make(chan *playwright.NetworkRequest, 1)
	errCh := make(chan error, 1)
	go func() {
		req, err := page.WaitForRequest(ctx, srv.Prefix()+"/empty.html")
		if err != nil {
			errCh <- err
			return
		}
		resultCh <- req
	}()

	_, err := page.Evaluate(ctx, `() => fetch("/empty.html")`)
	must.NoError(err)

	select {
	case req := <-resultCh:
		is.Equal(srv.Prefix()+"/empty.html", req.URL())
	case err := <-errCh:
		must.NoError(err)
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for request")
	}
}

// Ref: TestPageWaitForRequest.java#shouldWorkWithPredicate  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageWaitForRequest.java
func TestPageWaitForRequestShouldWorkWithPredicate(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)
	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	resultCh := make(chan *playwright.NetworkRequest, 1)
	errCh := make(chan error, 1)
	go func() {
		req, err := page.WaitForRequest(ctx, func(r *playwright.NetworkRequest) bool {
			return strings.Contains(r.URL(), "/empty.html")
		})
		if err != nil {
			errCh <- err
			return
		}
		resultCh <- req
	}()

	_, err := page.Evaluate(ctx, `() => fetch("/empty.html")`)
	must.NoError(err)

	select {
	case req := <-resultCh:
		is.Contains(req.URL(), "/empty.html")
	case err := <-errCh:
		must.NoError(err)
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for predicate match")
	}
}

// Ref: TestPageWaitForRequest.java#shouldRespectTimeout  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageWaitForRequest.java
func TestPageWaitForRequestShouldRespectTimeout(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)
	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	ctx, cancel := context.WithTimeout(ctx, 200*time.Millisecond)
	defer cancel()
	_, err := page.WaitForRequest(ctx, srv.Prefix()+"/missing.html", 750*time.Millisecond)
	is.Error(err)
}

// Ref: TestPageWaitForRequest.java#shouldRespectDefaultTimeout  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageWaitForRequest.java
func TestPageWaitForRequestShouldRespectDefaultTimeout(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)
	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	must.NoError(page.SetDefaultTimeout(ctx, 200))

	_, err := page.WaitForRequest(ctx, srv.Prefix()+"/missing.html")
	is.Error(err)
}

// Ref: TestPageWaitForRequest.java#shouldWorkWithNoTimeout  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageWaitForRequest.java
func TestPageWaitForRequestShouldWorkWithNoTimeout(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)
	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	resultCh := make(chan *playwright.NetworkRequest, 1)
	errCh := make(chan error, 1)
	go func() {
		req, err := page.WaitForRequest(ctx, 0)
		if err != nil {
			errCh <- err
			return
		}
		resultCh <- req
	}()

	_, err := page.Evaluate(ctx, `() => fetch("/empty.html")`)
	must.NoError(err)

	select {
	case req := <-resultCh:
		is.Contains(req.URL(), "/empty.html")
	case err := <-errCh:
		must.NoError(err)
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for no-timeout request")
	}
}

// Ref: TestPageWaitForRequest.java#shouldWorkWithUrlMatch  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageWaitForRequest.java
func TestPageWaitForRequestShouldWorkWithUrlMatch(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)
	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	resultCh := make(chan *playwright.NetworkRequest, 1)
	errCh := make(chan error, 1)
	go func() {
		req, err := page.WaitForRequest(ctx, srv.Prefix()+"*")
		if err != nil {
			errCh <- err
			return
		}
		resultCh <- req
	}()

	_, err := page.Evaluate(ctx, `() => fetch("/empty.html")`)
	must.NoError(err)

	select {
	case req := <-resultCh:
		is.True(strings.HasPrefix(req.URL(), srv.Prefix()))
	case err := <-errCh:
		must.NoError(err)
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for URL match")
	}
}
