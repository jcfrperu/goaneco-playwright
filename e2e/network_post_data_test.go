//go:build e2e

// E2E tests for: TestNetworkPostData.java.
package e2e

import (
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

// TestNetworkPostDataShouldReturnCorrectPostDataBufferForUtf8Body verifies PostData() returns the correct UTF-8 encoded body.
// Ref: TestNetworkPostData.java#shouldReturnCorrectPostDataBufferForUtf8Body  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestNetworkPostData.java
func TestNetworkPostDataShouldReturnCorrectPostDataBufferForUtf8Body(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	value := "baẞ"
	req := captureRequestTo(t, page, "title.html", func() {
		_, err := page.Evaluate(ctx, `({url, value}) => {
			const request = new Request(url, {
				method: 'POST',
				body: JSON.stringify(value),
			});
			request.headers.set('content-type', 'application/json;charset=UTF-8');
			return fetch(request);
		}`, map[string]any{"url": srv.Prefix() + "/title.html", "value": value})
		if err != nil {
			t.Logf("Evaluate: %v", err)
		}
	})

	postData := req.PostData()
	is.NotEmpty(postData, "PostData() should not be empty")
	is.Contains(postData, "baẞ", "PostData should contain the unicode value")
}

// TestNetworkPostDataShouldReturnPostDataWOContentType verifies PostData() works when content-type is empty.
// Ref: TestNetworkPostData.java#shouldReturnPostDataWOContentType  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestNetworkPostData.java
func TestNetworkPostDataShouldReturnPostDataWOContentType(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	req := captureRequestTo(t, page, "title.html", func() {
		_, err := page.Evaluate(ctx, `({url}) => {
			const request = new Request(url, {
				method: 'POST',
				body: JSON.stringify({ value: 42 }),
			});
			request.headers.set('content-type', '');
			return fetch(request);
		}`, map[string]any{"url": srv.Prefix() + "/title.html"})
		if err != nil {
			t.Logf("Evaluate: %v", err)
		}
	})

	postData := req.PostData()
	is.Contains(postData, "42", "PostData should contain the JSON-encoded value 42")
}

// TestNetworkPostDataShouldReturnPostDataForPUTRequests verifies PostData() works with PUT requests.
// Ref: TestNetworkPostData.java#shouldReturnPostDataForPUTRequests  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestNetworkPostData.java
func TestNetworkPostDataShouldReturnPostDataForPUTRequests(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	req := captureRequestTo(t, page, "title.html", func() {
		_, err := page.Evaluate(ctx, `({url}) => {
			const request = new Request(url, {
				method: 'PUT',
				body: JSON.stringify({ value: 42 }),
			});
			return fetch(request);
		}`, map[string]any{"url": srv.Prefix() + "/title.html"})
		if err != nil {
			t.Logf("Evaluate: %v", err)
		}
	})

	postData := req.PostData()
	is.Contains(postData, "42", "PostData should contain value 42 for PUT request")
}
