//go:build e2e

// E2E tests for: TestPageDrop.java.
package e2e

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupDropzone(t *testing.T, ctx context.Context, page *playwright.Page) {
	t.Helper()
	err := page.SetContent(ctx, `<style>#dropzone { width: 300px; height: 200px; border: 2px dashed #888; }</style>
<div id="dropzone"></div>
<script>
  window.__dropInfo = null;
  const zone = document.getElementById('dropzone');
  zone.addEventListener('dragenter', e => e.preventDefault());
  zone.addEventListener('dragover', e => e.preventDefault());
  zone.addEventListener('drop', async e => {
    e.preventDefault();
    const files = [];
    for (const file of e.dataTransfer.files)
      files.push({ name: file.name, type: file.type, size: file.size, text: await file.text() });
    const data = {};
    for (const t of e.dataTransfer.types) {
      if (t !== 'Files')
        data[t] = e.dataTransfer.getData(t);
    }
    window.__dropInfo = { files, data };
  });
</script>`)
	require.NoError(t, err)
}

func waitForDropInfo(t *testing.T, ctx context.Context, page *playwright.Page) map[string]any {
	t.Helper()
	_, err := page.WaitForFunction(ctx, "() => window.__dropInfo !== null", nil)
	require.NoError(t, err)
	info, err := page.Evaluate(ctx, "() => window.__dropInfo")
	require.NoError(t, err)
	res, ok := info.(map[string]any)
	require.True(t, ok)
	return res
}

// Ref: TestPageDrop.java#shouldDropFilePayload  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageDrop.java
func TestPageDropShouldDropFilePayload(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	setupDropzone(t, ctx, page)
	must.NoError(page.Locator("#dropzone").Drop(ctx, playwright.DropPayload{
		FilePayloads: []playwright.FilePayload{
			{Name: "note.txt", MimeType: "text/plain", Buffer: []byte("hello")},
		},
	}))
	info := waitForDropInfo(t, ctx, page)
	files, ok := info["files"].([]any)
	must.True(ok)
	must.Len(files, 1)
	f0, ok := files[0].(map[string]any)
	must.True(ok)
	is.Equal("note.txt", f0["name"])
	is.Equal("text/plain", f0["type"])
	is.Equal("hello", f0["text"])
}

// Ref: TestPageDrop.java#shouldDropMultipleFilePayloads  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageDrop.java
func TestPageDropShouldDropMultipleFilePayloads(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	setupDropzone(t, ctx, page)
	must.NoError(page.Locator("#dropzone").Drop(ctx, playwright.DropPayload{
		FilePayloads: []playwright.FilePayload{
			{Name: "a.txt", MimeType: "text/plain", Buffer: []byte("AAA")},
			{Name: "b.txt", MimeType: "text/plain", Buffer: []byte("BB")},
		},
	}))
	info := waitForDropInfo(t, ctx, page)
	files, ok := info["files"].([]any)
	must.True(ok)
	must.Len(files, 2)
	f0 := files[0].(map[string]any)
	is.Equal("a.txt", f0["name"])
	is.Equal("AAA", f0["text"])
	f1 := files[1].(map[string]any)
	is.Equal("b.txt", f1["name"])
	is.Equal("BB", f1["text"])
}

// Ref: TestPageDrop.java#shouldDropClipboardLikeData  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageDrop.java
func TestPageDropShouldDropClipboardLikeData(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	setupDropzone(t, ctx, page)
	must.NoError(page.Locator("#dropzone").Drop(ctx, playwright.DropPayload{
		Data: map[string]string{
			"text/plain":    "hello world",
			"text/uri-list": "https://example.com",
		},
	}))
	info := waitForDropInfo(t, ctx, page)
	files, ok := info["files"].([]any)
	must.True(ok)
	is.Empty(files)
	droppedData, ok := info["data"].(map[string]any)
	must.True(ok)
	is.Equal("hello world", droppedData["text/plain"])
	is.Equal("https://example.com", droppedData["text/uri-list"])
}

// Ref: TestPageDrop.java#shouldDropFileByLocalPath  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageDrop.java
func TestPageDropShouldDropFileByLocalPath(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	setupDropzone(t, ctx, page)
	dir := t.TempDir()
	filePath := filepath.Join(dir, "hello.txt")
	must.NoError(os.WriteFile(filePath, []byte("path-content"), 0o600))

	must.NoError(page.Locator("#dropzone").Drop(ctx, playwright.DropPayload{
		Files: []string{filePath},
	}))
	info := waitForDropInfo(t, ctx, page)
	files, ok := info["files"].([]any)
	must.True(ok)
	must.Len(files, 1)
	f0 := files[0].(map[string]any)
	is.Equal("hello.txt", f0["name"])
	is.Equal("path-content", f0["text"])
}
