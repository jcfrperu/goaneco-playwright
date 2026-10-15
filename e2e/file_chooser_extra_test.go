//go:build e2e

// Additional E2E tests without an equivalent in playwright-java.
package e2e

import (
	"context"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

const fileInputHTML = `<input id="file" type="file"><span id="out"></span>
<script>
  document.getElementById('file').addEventListener('change', function() {
    document.getElementById('out').textContent = this.files[0].name;
  });
</script>`

// fileToUploadPath creates a temp file with "contents of the file" and returns its path.
func fileToUploadPath(t *testing.T) string {
	t.Helper()
	tmp, err := os.CreateTemp("", "playwright-upload-*.txt")
	require.NoError(t, err)
	t.Cleanup(func() { os.Remove(tmp.Name()) })
	_, err = tmp.WriteString("contents of the file")
	require.NoError(t, err)
	require.NoError(t, tmp.Close())
	return tmp.Name()
}

// testCtxWithTimeout creates a test context with a specific timeout.
func testCtxWithTimeout(t *testing.T, d time.Duration) (context.Context, context.CancelFunc) {
	t.Helper()
	return context.WithTimeout(testCtx(t), d)
}
