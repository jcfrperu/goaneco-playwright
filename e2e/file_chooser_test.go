//go:build e2e

// E2E tests for FileChooser and Locator.SetInputFiles.
package e2e

import (
	"context"
	"math"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const fileInputHTML = `<input id="file" type="file"><span id="out"></span>
<script>
  document.getElementById('file').addEventListener('change', function() {
    document.getElementById('out').textContent = this.files[0].name;
  });
</script>`

// TestLocatorSetInputFiles verifies that Locator.SetInputFiles sets a file on a file input
// and that the browser reflects the chosen file name.
func TestLocatorSetInputFiles(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	srv.ServeWithBody("/upload", "text/html", fileInputHTML)
	err := page.Goto(ctx, srv.Prefix()+"/upload")
	must.NoError(err, "Goto failed")

	// Create a temporary file to upload
	tmp, err := os.CreateTemp("", "playwright-*.txt")
	must.NoError(err, "failed to create temp file")
	defer os.Remove(tmp.Name())
	_, err = tmp.WriteString("hello playwright")
	must.NoError(err, "failed to write temp file")
	tmp.Close()

	err = page.Locator("#file").SetInputFiles(ctx, []string{tmp.Name()})
	must.NoError(err, "SetInputFiles failed")

	// The browser should have picked up the file; check via JS
	rawName, err := page.Evaluate(ctx, "() => document.querySelector('#file').files[0].name")
	must.NoError(err, "Evaluate failed")
	name, _ := rawName.(string)
	if !strings.HasSuffix(name, ".txt") {
		t.Errorf("expected a .txt filename, got %q", name)
	}
}

// TestPageOnFileChooser verifies that OnFileChooser fires when a file input is activated
// and that FileChooser.SetFiles correctly sets the selected file.
func TestPageOnFileChooser(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	srv.ServeWithBody("/upload-chooser", "text/html", fileInputHTML)
	err := page.Goto(ctx, srv.Prefix()+"/upload-chooser")
	must.NoError(err, "Goto failed")

	tmp, err := os.CreateTemp("", "playwright-chooser-*.txt")
	must.NoError(err, "failed to create temp file")
	defer os.Remove(tmp.Name())
	tmp.WriteString("chooser content")
	tmp.Close()

	chooserCh := make(chan *playwright.FileChooser, 1)
	cancel := page.OnFileChooser(func(fc *playwright.FileChooser) {
		select {
		case chooserCh <- fc:
		default:
		}
	})
	defer cancel()

	// Clicking the file input opens the native file chooser; Playwright intercepts it.
	err = page.Locator("#file").Click(ctx)
	must.NoError(err, "Click failed")

	var fc *playwright.FileChooser
	select {
	case fc = <-chooserCh:
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for fileChooser event")
	}

	err = fc.SetFiles(ctx, []string{tmp.Name()})
	must.NoError(err, "FileChooser.SetFiles failed")

	rawName, err := page.Evaluate(ctx, "() => document.querySelector('#file').files[0].name")
	must.NoError(err, "Evaluate failed")
	name, _ := rawName.(string)
	if !strings.HasSuffix(name, ".txt") {
		t.Errorf("expected a .txt filename, got %q", name)
	}
}

// TestFileChooserPage verifies that FileChooser.Page() returns the page that opened the chooser.
func TestFileChooserPage(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	srv.ServeWithBody("/fc-page", "text/html", `<input id="f" type="file">`)
	err := page.Goto(ctx, srv.Prefix()+"/fc-page")
	must.NoError(err, "Goto failed")

	chooserCh := make(chan *playwright.FileChooser, 1)
	cancel := page.OnFileChooser(func(fc *playwright.FileChooser) {
		select {
		case chooserCh <- fc:
		default:
		}
	})
	defer cancel()

	err = page.Locator("#f").Click(ctx)
	must.NoError(err, "Click failed")

	select {
	case fc := <-chooserCh:
		if fc.Page() != page {
			t.Error("FileChooser.Page() does not match the page that opened the chooser")
		}
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for fileChooser event")
	}
}

// TestFileChooserIsMultiple verifies that IsMultiple reflects the multiple attribute.
func TestFileChooserIsMultiple(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	srv.ServeWithBody("/upload-multi", "text/html", `<input id="file" type="file" multiple>`)
	err := page.Goto(ctx, srv.Prefix()+"/upload-multi")
	must.NoError(err, "Goto failed")

	chooserCh := make(chan *playwright.FileChooser, 1)
	cancel := page.OnFileChooser(func(fc *playwright.FileChooser) {
		select {
		case chooserCh <- fc:
		default:
		}
	})
	defer cancel()

	err = page.Locator("#file").Click(ctx)
	must.NoError(err, "Click failed")

	select {
	case fc := <-chooserCh:
		if !fc.IsMultiple() {
			t.Error("expected IsMultiple() = true for input with multiple attribute")
		}
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for fileChooser event")
	}
}

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

// TestSetInputFilesWork verifies SetInputFiles sets a file on a plain file input.
// Ref: TestPageSetInputFiles.java#shouldWork
func TestSetInputFilesWork(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)
	filePath := fileToUploadPath(t)

	err := page.SetContent(ctx, `<input type=file>`)
	must.NoError(err)

	err = page.Locator("input").SetInputFiles(ctx, []string{filePath})
	must.NoError(err)

	count, err := page.Evaluate(ctx, `() => document.querySelector('input').files.length`)
	must.NoError(err)
	is.Equal(float64(1), count)

	name, err := page.Evaluate(ctx, `() => document.querySelector('input').files[0].name`)
	must.NoError(err)
	is.NotEmpty(name)
}

// TestSetInputFilesFromMemory verifies SetInputFilesPayload with in-memory FilePayload.
// Ref: TestPageSetInputFiles.java#shouldSetFromMemory
func TestSetInputFilesFromMemory(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	err := page.SetContent(ctx, `<input type=file>`)
	must.NoError(err)

	err = page.Locator("input").SetInputFilesPayload(ctx, []playwright.FilePayload{
		{Name: "test.txt", MimeType: "text/plain", Buffer: []byte("this is a test")},
	})
	must.NoError(err)

	count, err := page.Evaluate(ctx, `() => document.querySelector('input').files.length`)
	must.NoError(err)
	is.Equal(float64(1), count)

	name, err := page.Evaluate(ctx, `() => document.querySelector('input').files[0].name`)
	must.NoError(err)
	is.Equal("test.txt", name)
}

// TestSetInputFilesEmitInputAndChangeEvents verifies that input and change events fire after setInputFiles.
// Ref: TestPageSetInputFiles.java#shouldEmitInputAndChangeEvents
func TestSetInputFilesEmitInputAndChangeEvents(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)
	filePath := fileToUploadPath(t)

	err := page.SetContent(ctx, `<input id=input type=file></input>
<script>
  window._events = [];
  input.addEventListener('input', e => window._events.push(e.type));
  input.addEventListener('change', e => window._events.push(e.type));
</script>`)
	must.NoError(err)

	err = page.Locator("input").SetInputFiles(ctx, []string{filePath})
	must.NoError(err)

	raw, err := page.Evaluate(ctx, "() => window._events")
	must.NoError(err)
	got, ok := raw.([]any)
	is.True(ok, "expected array of events, got %T", raw)
	is.Len(got, 2)
	is.Equal("input", got[0])
	is.Equal("change", got[1])
}

// TestFileChooserWorkWhenNotAttachedToDOM verifies FileChooser fires for dynamically created input.
// Ref: TestPageSetInputFiles.java#shouldWorkWhenFileInputIsNotAttachedToDOM
func TestFileChooserWorkWhenNotAttachedToDOM(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx, cancel := testCtxWithTimeout(t, 10*time.Second)
	defer cancel()
	page := newPage(t)

	fc, err := page.WaitForFileChooser(ctx, func() {
		_, _ = page.Evaluate(ctx, `() => {
			const el = document.createElement('input');
			el.type = 'file';
			el.click();
		}`)
	})
	must.NoError(err)
	must.NotNil(fc)
}

// TestFileChooserAcceptSingleFile verifies FileChooser.SetFiles sets a single file correctly.
// Ref: TestPageSetInputFiles.java#shouldAcceptSingleFile
func TestFileChooserAcceptSingleFile(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx, cancel := testCtxWithTimeout(t, 10*time.Second)
	defer cancel()
	page := newPage(t)
	filePath := fileToUploadPath(t)

	err := page.SetContent(ctx, `<input type=file>`)
	must.NoError(err)

	fc, err := page.WaitForFileChooser(ctx, func() {
		_ = page.Locator("input").Click(ctx)
	})
	must.NoError(err)
	must.NotNil(fc)
	is.Equal(page, fc.Page())

	err = fc.SetFiles(ctx, []string{filePath})
	must.NoError(err)

	count, err := page.Evaluate(ctx, `() => document.querySelector('input').files.length`)
	must.NoError(err)
	is.Equal(float64(1), count)
}

// TestFileChooserAcceptSingleFilePayload verifies FileChooser.SetFilesPayload with in-memory data.
// Ref: TestPageSetInputFiles.java#shouldAcceptSingleFilePayload
func TestFileChooserAcceptSingleFilePayload(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx, cancel := testCtxWithTimeout(t, 10*time.Second)
	defer cancel()
	page := newPage(t)

	err := page.SetContent(ctx, `<input type=file>`)
	must.NoError(err)

	fc, err := page.WaitForFileChooser(ctx, func() {
		_ = page.Locator("input").Click(ctx)
	})
	must.NoError(err)
	must.NotNil(fc)

	err = fc.SetFilesPayload(ctx, []playwright.FilePayload{
		{Name: "test.txt", MimeType: "text/plain", Buffer: []byte("Hello!")},
	})
	must.NoError(err)

	count, err := page.Evaluate(ctx, `() => document.querySelector('input').files.length`)
	must.NoError(err)
	is.Equal(float64(1), count)

	name, err := page.Evaluate(ctx, `() => document.querySelector('input').files[0].name`)
	must.NoError(err)
	is.Equal("test.txt", name)
}

// TestFileChooserIsNotMultiple verifies that a single-file input reports isMultiple=false.
// Ref: TestPageSetInputFiles.java#shouldWorkForSingleFilePick
func TestFileChooserIsNotMultiple(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx, cancel := testCtxWithTimeout(t, 10*time.Second)
	defer cancel()
	page := newPage(t)

	err := page.SetContent(ctx, `<input type=file>`)
	must.NoError(err)

	fc, err := page.WaitForFileChooser(ctx, func() {
		_ = page.Locator("input").Click(ctx)
	})
	must.NoError(err)
	is.False(fc.IsMultiple(), "single-file input should report IsMultiple=false")
}

// TestFileChooserIsMultipleWithMultiple verifies that a multiple-file input reports isMultiple=true.
// Ref: TestPageSetInputFiles.java#shouldWorkForMultiple
func TestFileChooserIsMultipleWithMultiple(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx, cancel := testCtxWithTimeout(t, 10*time.Second)
	defer cancel()
	page := newPage(t)

	err := page.SetContent(ctx, `<input multiple type=file>`)
	must.NoError(err)

	fc, err := page.WaitForFileChooser(ctx, func() {
		_ = page.Locator("input").Click(ctx)
	})
	must.NoError(err)
	is.True(fc.IsMultiple(), "multiple-file input should report IsMultiple=true")
}

// TestFileChooserIsMultipleWithWebkitdirectory verifies that a webkitdirectory input reports isMultiple=true.
// Ref: TestPageSetInputFiles.java#shouldWorkForWebkitdirectory
func TestFileChooserIsMultipleWithWebkitdirectory(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx, cancel := testCtxWithTimeout(t, 10*time.Second)
	defer cancel()
	page := newPage(t)

	err := page.SetContent(ctx, `<input multiple webkitdirectory type=file>`)
	must.NoError(err)

	fc, err := page.WaitForFileChooser(ctx, func() {
		_ = page.Locator("input").Click(ctx)
	})
	must.NoError(err)
	is.True(fc.IsMultiple(), "webkitdirectory input should report IsMultiple=true")
}

// TestSetInputFilesPreserveTimestamp verifies that the browser preserves the file's last modified timestamp.
// Ref: TestPageSetInputFiles.java#shouldPreserveLastModifiedTimestamp
func TestSetInputFilesPreserveTimestamp(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)
	filePath := fileToUploadPath(t)

	// Get the file's modification time
	info, err := os.Stat(filePath)
	must.NoError(err)
	expectedMs := float64(info.ModTime().UnixMilli())

	err = page.SetContent(ctx, `<input type=file multiple=true/>`)
	must.NoError(err)

	err = page.Locator("input").SetInputFiles(ctx, []string{filePath})
	must.NoError(err)

	timestamps, err := page.Evaluate(ctx, `() => [...document.querySelector('input').files].map(f => f.lastModified)`)
	must.NoError(err)

	tsSlice, ok := timestamps.([]any)
	is.True(ok, "expected array of timestamps")
	is.Len(tsSlice, 1)

	actualMs, _ := tsSlice[0].(float64)
	// Browser may round timestamps; allow ±1 second tolerance
	is.True(math.Abs(actualMs-expectedMs) < 1000,
		"timestamp mismatch: expected %v, got %v", expectedMs, actualMs)
}

// testCtxWithTimeout creates a test context with a specific timeout.
func testCtxWithTimeout(t *testing.T, d time.Duration) (context.Context, context.CancelFunc) {
	t.Helper()
	return context.WithTimeout(testCtx(t), d)
}

// TestSetInputFilesWithMultipleFiles verifies SetInputFiles with multiple files works.
// Ref: TestPageSetInputFiles.java#shouldWorkWithMultipleFiles
func TestSetInputFilesWithMultipleFiles(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	dir := t.TempDir()
	f1 := filepath.Join(dir, "file1.txt")
	f2 := filepath.Join(dir, "file2.txt")
	must.NoError(os.WriteFile(f1, []byte("content1"), 0o600))
	must.NoError(os.WriteFile(f2, []byte("content2"), 0o600))

	must.NoError(page.SetContent(ctx, `<input type="file" multiple>`))
	must.NoError(page.Locator("input").SetInputFiles(ctx, []string{f1, f2}))

	count, err := page.Evaluate(ctx, `() => document.querySelector('input').files.length`)
	must.NoError(err)
	is.Equal(float64(2), count)
}

// TestSetInputFilesClearsOnEmptySlice verifies SetInputFiles with empty slice clears selection.
// Ref: TestPageSetInputFiles.java#shouldClearFiles
func TestSetInputFilesClearsOnEmptySlice(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	dir := t.TempDir()
	f1 := filepath.Join(dir, "file.txt")
	must.NoError(os.WriteFile(f1, []byte("content"), 0o600))

	must.NoError(page.SetContent(ctx, `<input type="file">`))
	must.NoError(page.Locator("input").SetInputFiles(ctx, []string{f1}))

	count1, err := page.Evaluate(ctx, `() => document.querySelector('input').files.length`)
	must.NoError(err)
	is.Equal(float64(1), count1)

	// Clear files by setting empty slice
	must.NoError(page.Locator("input").SetInputFiles(ctx, []string{}))

	count2, err := page.Evaluate(ctx, `() => document.querySelector('input').files.length`)
	must.NoError(err)
	is.Equal(float64(0), count2)
}

// TestSetInputFilesFiresChangeEvent verifies SetInputFiles fires a change event.
// Ref: TestPageSetInputFiles.java#shouldFireChangeEvent
func TestSetInputFilesFiresChangeEvent(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	page := newPage(t)

	dir := t.TempDir()
	f1 := filepath.Join(dir, "changed.txt")
	must.NoError(os.WriteFile(f1, []byte("file"), 0o600))

	must.NoError(page.SetContent(ctx, `
		<input type="file" onchange="document.getElementById('result').textContent = this.files.length">
		<div id="result"></div>
	`))

	must.NoError(page.Locator("input").SetInputFiles(ctx, []string{f1}))

	text, err := page.Locator("#result").InnerText(ctx)
	must.NoError(err)
	is.Equal("1", text)
}
