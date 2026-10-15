//go:build e2e

// E2E tests for: TestPageWaitForResponse.java.
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

// Ref: TestPageWaitForResponse.java#shouldWork  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageWaitForResponse.java
func TestPageWaitForResponseShouldWork(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)
	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	resultCh := make(chan *playwright.NetworkResponse, 1)
	errCh := make(chan error, 1)
	go func() {
		resp, err := page.WaitForResponse(ctx, srv.Prefix()+"/empty.html")
		if err != nil {
			errCh <- err
			return
		}
		resultCh <- resp
	}()

	_, err := page.Evaluate(ctx, `() => fetch("/empty.html")`)
	must.NoError(err)

	select {
	case resp := <-resultCh:
		is.Equal(srv.Prefix()+"/empty.html", resp.URL())
	case err := <-errCh:
		must.NoError(err)
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for response")
	}
}

// Ref: TestPageWaitForResponse.java#shouldRespectTimeout  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageWaitForResponse.java
func TestPageWaitForResponseShouldRespectTimeout(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)
	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	ctx, cancel := context.WithTimeout(ctx, 200*time.Millisecond)
	defer cancel()
	_, err := page.WaitForResponse(ctx, srv.Prefix()+"/missing.html", 750*time.Millisecond)
	is.Error(err)
}

// Ref: TestPageWaitForResponse.java#shouldWorkWithPredicate  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageWaitForResponse.java
func TestPageWaitForResponseShouldWorkWithPredicate(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)
	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	resultCh := make(chan *playwright.NetworkResponse, 1)
	errCh := make(chan error, 1)
	go func() {
		resp, err := page.WaitForResponse(ctx, func(r *playwright.NetworkResponse) bool {
			return strings.Contains(r.URL(), "/empty.html")
		})
		if err != nil {
			errCh <- err
			return
		}
		resultCh <- resp
	}()

	_, err := page.Evaluate(ctx, `() => fetch("/empty.html")`)
	must.NoError(err)

	select {
	case resp := <-resultCh:
		is.Contains(resp.URL(), "/empty.html")
	case err := <-errCh:
		must.NoError(err)
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for response predicate")
	}
}

// Ref: TestPageWaitForResponse.java#shouldRespectDefaultTimeout  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageWaitForResponse.java
func TestPageWaitForResponseShouldRespectDefaultTimeout(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)
	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	must.NoError(page.SetDefaultTimeout(ctx, 200))

	_, err := page.WaitForResponse(ctx, srv.Prefix()+"/missing.html")
	is.Error(err)
}

// Ref: TestPageWaitForResponse.java#shouldWorkWithNoTimeout  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageWaitForResponse.java
func TestPageWaitForResponseShouldWorkWithNoTimeout(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)
	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	resultCh := make(chan *playwright.NetworkResponse, 1)
	errCh := make(chan error, 1)
	go func() {
		resp, err := page.WaitForResponse(ctx, 0)
		if err != nil {
			errCh <- err
			return
		}
		resultCh <- resp
	}()

	_, err := page.Evaluate(ctx, `() => fetch("/empty.html")`)
	must.NoError(err)

	select {
	case resp := <-resultCh:
		is.Contains(resp.URL(), "/empty.html")
	case err := <-errCh:
		must.NoError(err)
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for no-timeout response")
	}
}
