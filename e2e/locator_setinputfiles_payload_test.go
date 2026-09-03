//go:build e2e

// E2E tests for Locator.SetInputFilesPayload with FilePayload.
// Migration of: TestPageSetInputFiles.java (payload/buffer cases)
package e2e

import (
	"testing"

	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSetInputFilesPayloadSingleFile verifies SetInputFilesPayload with a single file payload.
// Ref: TestPageSetInputFiles.java#shouldWorkWithFilePayload
func TestSetInputFilesPayloadSingleFile(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type="file">`))
	must.NoError(page.Locator("input").SetInputFilesPayload(ctx, []playwright.FilePayload{
		{Name: "test.txt", MimeType: "text/plain", Buffer: []byte("hello world")},
	}))

	count, err := page.Evaluate(ctx, `() => document.querySelector('input').files.length`)
	must.NoError(err)
	is.Equal(float64(1), count)
}

// TestSetInputFilesPayloadFileName verifies the filename is preserved.
// Ref: TestPageSetInputFiles.java#shouldPreserveFileName
func TestSetInputFilesPayloadFileName(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type="file">`))
	must.NoError(page.Locator("input").SetInputFilesPayload(ctx, []playwright.FilePayload{
		{Name: "myfile.csv", MimeType: "text/csv", Buffer: []byte("a,b,c")},
	}))

	name, err := page.Evaluate(ctx, `() => document.querySelector('input').files[0].name`)
	must.NoError(err)
	is.Equal("myfile.csv", name)
}

// TestSetInputFilesPayloadMimeType verifies the MIME type is set correctly.
// Ref: TestPageSetInputFiles.java#shouldPreserveMimeType
func TestSetInputFilesPayloadMimeType(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type="file">`))
	must.NoError(page.Locator("input").SetInputFilesPayload(ctx, []playwright.FilePayload{
		{Name: "image.png", MimeType: "image/png", Buffer: []byte("\x89PNG")},
	}))

	mimeType, err := page.Evaluate(ctx, `() => document.querySelector('input').files[0].type`)
	must.NoError(err)
	is.Equal("image/png", mimeType)
}

// TestSetInputFilesPayloadFileSize verifies the file size matches buffer length.
// Ref: TestPageSetInputFiles.java#shouldHaveCorrectSize
func TestSetInputFilesPayloadFileSize(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	content := []byte("content of the file")
	must.NoError(page.SetContent(ctx, `<input type="file">`))
	must.NoError(page.Locator("input").SetInputFilesPayload(ctx, []playwright.FilePayload{
		{Name: "data.txt", MimeType: "text/plain", Buffer: content},
	}))

	size, err := page.Evaluate(ctx, `() => document.querySelector('input').files[0].size`)
	must.NoError(err)
	is.Equal(float64(len(content)), size)
}

// TestSetInputFilesPayloadMultipleFiles verifies multiple file payloads can be set.
// Ref: TestPageSetInputFiles.java#shouldWorkWithMultipleFilePayloads
func TestSetInputFilesPayloadMultipleFiles(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	must.NoError(page.SetContent(ctx, `<input type="file" multiple>`))
	must.NoError(page.Locator("input").SetInputFilesPayload(ctx, []playwright.FilePayload{
		{Name: "a.txt", MimeType: "text/plain", Buffer: []byte("aaa")},
		{Name: "b.txt", MimeType: "text/plain", Buffer: []byte("bbb")},
	}))

	count, err := page.Evaluate(ctx, `() => document.querySelector('input').files.length`)
	must.NoError(err)
	is.Equal(float64(2), count)
}
