//go:build e2e

// Additional E2E tests without an equivalent in playwright-java.
package e2e

import (
	"fmt"
	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

// writeHAR writes a single-entry HAR file and returns its path.
func writeHAR(t *testing.T, url, mimeType, body string, status int) string {
	t.Helper()
	content := fmt.Sprintf(`{"log":{"version":"1.2","creator":{"name":"test","version":"1"},"entries":[{`+
		`"startedDateTime":"2024-01-01T00:00:00.000Z",`+
		`"time":0,`+
		`"request":{"method":"GET","url":%q,"headers":[],"queryString":[],"cookies":[],"headersSize":-1,"bodySize":-1},`+
		`"response":{"status":%d,"statusText":"OK",`+
		`"headers":[{"name":"content-type","value":%q}],`+
		`"cookies":[],"content":{"size":%d,"mimeType":%q,"text":%q},`+
		`"redirectURL":"","headersSize":-1,"bodySize":-1},`+
		`"cache":{},"timings":{"send":0,"wait":0,"receive":0}`+
		`}]}}`,
		url, status, mimeType, len(body), mimeType, body)

	dir, err := os.MkdirTemp("", "har-*")
	require.NoError(t, err, "MkdirTemp")
	t.Cleanup(func() { _ = os.RemoveAll(dir) })

	path := filepath.Join(dir, "test.har")
	err = os.WriteFile(path, []byte(content), 0600)
	require.NoError(t, err, "WriteFile")
	return path
}

// Ensure imports are used.
var _ *playwright.RouteFromHAROptions
