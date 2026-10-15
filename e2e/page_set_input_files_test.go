//go:build e2e

// E2E tests for: TestPageSetInputFiles.java.
// Tests marked "Extra:" are additional tests with no equivalent method in playwright-java.
package e2e

import (
	"context"
	"fmt"
	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"
)

// TestLocatorSetInputFiles verifies that Locator.SetInputFiles sets a file on a file input
// and that the browser reflects the chosen file name.
// Extra: no equivalent method in playwright-java; related class: TestPageSetInputFiles.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSetInputFiles.java
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
// Extra: no equivalent method in playwright-java; related class: TestPageSetInputFiles.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSetInputFiles.java
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
// Extra: no equivalent method in playwright-java; related class: TestPageSetInputFiles.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSetInputFiles.java
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
// Extra: no equivalent method in playwright-java; related class: TestPageSetInputFiles.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSetInputFiles.java
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

// TestSetInputFilesWithMultipleFiles verifies SetInputFiles with multiple files works.
// Extra: no equivalent method in playwright-java; related class: TestPageSetInputFiles.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSetInputFiles.java
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
// Extra: no equivalent method in playwright-java; related class: TestPageSetInputFiles.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSetInputFiles.java
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
// Extra: no equivalent method in playwright-java; related class: TestPageSetInputFiles.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSetInputFiles.java
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

// TestSetInputFilesPayloadSingleFile verifies SetInputFilesPayload with a single file payload.
// Extra: no equivalent method in playwright-java; related class: TestPageSetInputFiles.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSetInputFiles.java
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
// Extra: no equivalent method in playwright-java; related class: TestPageSetInputFiles.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSetInputFiles.java
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
// Extra: no equivalent method in playwright-java; related class: TestPageSetInputFiles.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSetInputFiles.java
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
// Extra: no equivalent method in playwright-java; related class: TestPageSetInputFiles.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSetInputFiles.java
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
// Extra: no equivalent method in playwright-java; related class: TestPageSetInputFiles.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSetInputFiles.java
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

// TestPageSetInputFilesShouldWork verifies SetInputFiles sets a file on a plain file input.
// Ref: TestPageSetInputFiles.java#shouldWork  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSetInputFiles.java
func TestPageSetInputFilesShouldWork(t *testing.T) {
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

// TestPageSetInputFilesShouldSetFromMemory verifies SetInputFilesPayload with in-memory FilePayload.
// Ref: TestPageSetInputFiles.java#shouldSetFromMemory  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSetInputFiles.java
func TestPageSetInputFilesShouldSetFromMemory(t *testing.T) {
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

// TestPageSetInputFilesShouldEmitInputAndChangeEvents verifies that input and change events fire after setInputFiles.
// Ref: TestPageSetInputFiles.java#shouldEmitInputAndChangeEvents  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSetInputFiles.java
func TestPageSetInputFilesShouldEmitInputAndChangeEvents(t *testing.T) {
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

// TestPageSetInputFilesShouldWorkWhenFileInputIsNotAttachedToDOM verifies FileChooser fires for dynamically created input.
// Ref: TestPageSetInputFiles.java#shouldWorkWhenFileInputIsNotAttachedToDOM  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSetInputFiles.java
func TestPageSetInputFilesShouldWorkWhenFileInputIsNotAttachedToDOM(t *testing.T) {
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

// TestPageSetInputFilesShouldAcceptSingleFile verifies FileChooser.SetFiles sets a single file correctly.
// Ref: TestPageSetInputFiles.java#shouldAcceptSingleFile  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSetInputFiles.java
func TestPageSetInputFilesShouldAcceptSingleFile(t *testing.T) {
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

// TestPageSetInputFilesShouldAcceptSingleFilePayload verifies FileChooser.SetFilesPayload with in-memory data.
// Ref: TestPageSetInputFiles.java#shouldAcceptSingleFilePayload  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSetInputFiles.java
func TestPageSetInputFilesShouldAcceptSingleFilePayload(t *testing.T) {
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

// TestPageSetInputFilesShouldWorkForSingleFilePick verifies that a single-file input reports isMultiple=false.
// Ref: TestPageSetInputFiles.java#shouldWorkForSingleFilePick  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSetInputFiles.java
func TestPageSetInputFilesShouldWorkForSingleFilePick(t *testing.T) {
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

// TestPageSetInputFilesShouldWorkForMultiple verifies that a multiple-file input reports isMultiple=true.
// Ref: TestPageSetInputFiles.java#shouldWorkForMultiple  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSetInputFiles.java
func TestPageSetInputFilesShouldWorkForMultiple(t *testing.T) {
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

// TestPageSetInputFilesShouldWorkForWebkitdirectory verifies that a webkitdirectory input reports isMultiple=true.
// Ref: TestPageSetInputFiles.java#shouldWorkForWebkitdirectory  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSetInputFiles.java
func TestPageSetInputFilesShouldWorkForWebkitdirectory(t *testing.T) {
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

// TestPageSetInputFilesShouldPreserveLastModifiedTimestamp verifies that the browser preserves the file's last modified timestamp.
// Ref: TestPageSetInputFiles.java#shouldPreserveLastModifiedTimestamp  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSetInputFiles.java
func TestPageSetInputFilesShouldPreserveLastModifiedTimestamp(t *testing.T) {
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

// TestPageSetInputFilesShouldUploadTheFile verifies uploading a file and reading it back.
// Ref: TestPageSetInputFiles.java#shouldUploadTheFile  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSetInputFiles.java
func TestPageSetInputFilesShouldUploadTheFile(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)
	a4ServeUpload(srv)
	must.NoError(page.Goto(ctx, srv.Prefix()+"/input/fileupload.html"))
	input := page.Locator("input").First()
	path := fileToUploadPath(t)
	must.NoError(input.SetInputFiles(ctx, []string{path}))
	name, err := input.Evaluate(ctx, "e => e.files[0].name")
	must.NoError(err)
	is.Equal(filepath.Base(path), name)
	content, err := input.Evaluate(ctx, a4ReadFileJS)
	must.NoError(err)
	is.Equal("contents of the file", content)
}

// a4LargeUpload is the shared body of the 200MB upload tests.
func a4LargeUpload(t *testing.T, relative bool) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)
	a4ServeUpload(srv)
	must.NoError(page.Goto(ctx, srv.Prefix()+"/input/fileupload.html"))
	uploadFile := filepath.Join(t.TempDir(), "200MB.zip")
	a4WriteBig(t, uploadFile, 200*1024*1024)
	input := page.Locator("input[type='file']")
	// Locator.EvaluateHandle is not exposed; the events array is kept on window instead of a JSHandle.
	_, err := input.Evaluate(ctx, `e => {
    const events = [];
    window.__a4events = events;
    e.addEventListener('input', () => events.push('input'));
    e.addEventListener('change', () => events.push('change'));
  }`)
	must.NoError(err)
	setPath := uploadFile
	if relative {
		wd, err := os.Getwd()
		must.NoError(err)
		setPath, err = filepath.Rel(wd, uploadFile)
		if err != nil {
			t.Skip("cannot build a relative path: " + err.Error())
		}
		is.False(filepath.IsAbs(setPath))
	}
	must.NoError(input.SetInputFiles(ctx, []string{setPath}))
	name, err := input.Evaluate(ctx, "e => e.files[0].name")
	must.NoError(err)
	is.Equal("200MB.zip", name)
	evs, err := page.Evaluate(ctx, "() => window.__a4events")
	must.NoError(err)
	is.Equal([]any{"input", "change"}, evs)
	up := a4UploadRoute(srv)
	clickCtx, cancel := context.WithTimeout(ctx, 90*time.Second)
	defer cancel()
	must.NoError(page.Locator("input[type=submit]").Click(clickCtx))
	files := up.Wait(t, 90*time.Second)
	is.Equal([]a4UploadedFile{{Name: "200MB.zip", Size: 200 * 1024 * 1024}}, files)
}

// TestPageSetInputFilesShouldUploadLargeFile verifies uploading a 200MB file.
// Ref: TestPageSetInputFiles.java#shouldUploadLargeFile  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSetInputFiles.java
func TestPageSetInputFilesShouldUploadLargeFile(t *testing.T) {
	a4LargeUpload(t, false)
}

// TestPageSetInputFilesShouldUploadLargeFileWithRelativePath verifies uploading a 200MB file by relative path.
// Ref: TestPageSetInputFiles.java#shouldUploadLargeFileWithRelativePath  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSetInputFiles.java
func TestPageSetInputFilesShouldUploadLargeFileWithRelativePath(t *testing.T) {
	a4LargeUpload(t, true)
}

// TestPageSetInputFilesShouldUploadMultipleLargeFiles verifies choosing ten 49MB files via a file chooser.
// Ref: TestPageSetInputFiles.java#shouldUploadMultipleLargeFiles  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSetInputFiles.java
func TestPageSetInputFilesShouldUploadMultipleLargeFiles(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)
	a4ServeUploadMulti(srv)
	must.NoError(page.Goto(ctx, srv.Prefix()+"/input/fileupload-multi.html"))
	dir := t.TempDir()
	first := filepath.Join(dir, "50MB_0.zip")
	a4WriteBig(t, first, 49*1024*1024)
	files := []string{first}
	data, err := os.ReadFile(first)
	must.NoError(err)
	for i := 1; i < 10; i++ {
		p := filepath.Join(dir, fmt.Sprintf("50MB_%d.zip", i))
		must.NoError(os.WriteFile(p, data, 0o644))
		files = append(files, p)
	}
	input := page.Locator("input[type='file']")
	fc, err := page.WaitForFileChooser(ctx, func() { _ = input.Click(ctx) })
	must.NoError(err)
	must.NoError(fc.SetFiles(ctx, files))
	n, err := input.Evaluate(ctx, "e => e.files.length")
	must.NoError(err)
	is.True(fc.IsMultiple())
	is.EqualValues(10, n)
}

// TestPageSetInputFilesShouldEmitEventOnce verifies a file chooser is emitted on click.
// Ref: TestPageSetInputFiles.java#shouldEmitEventOnce  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSetInputFiles.java
func TestPageSetInputFilesShouldEmitEventOnce(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx := testCtx(t)
	page := newPage(t)
	must.NoError(page.SetContent(ctx, "<input type=file>"))
	fc, err := page.WaitForFileChooser(ctx, func() { _ = page.Locator("input").Click(ctx) })
	must.NoError(err)
	must.NotNil(fc)
}

// TestPageSetInputFilesShouldEmitEventAddListenerRemoveListener verifies a listener can remove itself.
// Ref: TestPageSetInputFiles.java#shouldEmitEventAddListenerRemoveListener  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSetInputFiles.java
func TestPageSetInputFilesShouldEmitEventAddListenerRemoveListener(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx := testCtx(t)
	page := newPage(t)
	must.NoError(page.SetContent(ctx, "<input type=file>"))
	ch := make(chan *playwright.FileChooser, 4)
	var off func()
	var mu sync.Mutex
	off = page.OnFileChooser(func(fc *playwright.FileChooser) {
		select {
		case ch <- fc:
		default:
		}
		mu.Lock()
		o := off
		mu.Unlock()
		if o != nil {
			o()
		}
	})
	must.NoError(page.Locator("input").Click(ctx))
	select {
	case fc := <-ch:
		must.NotNil(fc)
	case <-time.After(10 * time.Second):
		t.Fatal("file chooser event not received")
	}
}

// TestPageSetInputFilesShouldWorkWhenFileInputIsAttachedToDOM verifies a chooser opens for an attached input.
// Ref: TestPageSetInputFiles.java#shouldWorkWhenFileInputIsAttachedToDOM  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSetInputFiles.java
func TestPageSetInputFilesShouldWorkWhenFileInputIsAttachedToDOM(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx := testCtx(t)
	page := newPage(t)
	must.NoError(page.SetContent(ctx, "<input type=file>"))
	fc, err := page.WaitForFileChooser(ctx, func() { _ = page.Locator("input").Click(ctx) })
	must.NoError(err)
	must.NotNil(fc)
}

// TestPageSetInputFilesShouldWorkWithCSP verifies setting files on a page with a strict CSP.
// Ref: TestPageSetInputFiles.java#shouldWorkWithCSP  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSetInputFiles.java
func TestPageSetInputFilesShouldWorkWithCSP(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)
	srv.SetRoute("/empty.html", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", "default-src 'none'")
		w.Header().Set("Content-Type", "text/html")
	})
	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	must.NoError(page.SetContent(ctx, "<input type=file>"))
	path := fileToUploadPath(t)
	must.NoError(page.Locator("input").SetInputFiles(ctx, []string{path}))
	n, err := page.EvalOnSelector(ctx, "input", "input => input.files.length")
	must.NoError(err)
	is.EqualValues(1, n)
	name, err := page.EvalOnSelector(ctx, "input", "input => input.files[0].name")
	must.NoError(err)
	is.Equal(filepath.Base(path), name)
}

// TestPageSetInputFilesShouldRespectTimeout verifies WaitForFileChooser honours a 1ms deadline.
// Ref: TestPageSetInputFiles.java#shouldRespectTimeout  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSetInputFiles.java
func TestPageSetInputFilesShouldRespectTimeout(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx := testCtx(t)
	page := newPage(t)
	tctx, cancel := context.WithTimeout(ctx, time.Millisecond)
	defer cancel()
	_, err := page.WaitForFileChooser(tctx, func() {})
	must.Error(err)
	must.ErrorIs(err, context.DeadlineExceeded)
}

// TestPageSetInputFilesShouldRespectDefaultTimeoutWhenThereIsNoCustomTimeout is not ported.
// Ref: TestPageSetInputFiles.java#shouldRespectDefaultTimeoutWhenThereIsNoCustomTimeout  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSetInputFiles.java
func TestPageSetInputFilesShouldRespectDefaultTimeoutWhenThereIsNoCustomTimeout(t *testing.T) {
	t.Skip("not implemented: Page.WaitForFileChooser ignores Page.SetDefaultTimeout (only the ctx deadline applies)")
}

// TestPageSetInputFilesShouldPrioritizeExactTimeoutOverDefaultTimeout verifies the explicit deadline wins over default timeout 0.
// Ref: TestPageSetInputFiles.java#shouldPrioritizeExactTimeoutOverDefaultTimeout  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSetInputFiles.java
func TestPageSetInputFilesShouldPrioritizeExactTimeoutOverDefaultTimeout(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx := testCtx(t)
	page := newPage(t)
	must.NoError(page.SetDefaultTimeout(ctx, 0))
	tctx, cancel := context.WithTimeout(ctx, time.Millisecond)
	defer cancel()
	_, err := page.WaitForFileChooser(tctx, func() {})
	must.Error(err)
	must.ErrorIs(err, context.DeadlineExceeded)
}

// TestPageSetInputFilesShouldWorkWithNoTimeout verifies a delayed chooser is received.
// Ref: TestPageSetInputFiles.java#shouldWorkWithNoTimeout  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSetInputFiles.java
func TestPageSetInputFilesShouldWorkWithNoTimeout(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx := testCtx(t)
	page := newPage(t)
	fc, err := page.WaitForFileChooser(ctx, func() {
		_, _ = page.Evaluate(ctx, `() => setTimeout(() => {
  const el = document.createElement('input');
  el.type = 'file';
  el.click();
}, 50)`)
	})
	must.NoError(err)
	must.NotNil(fc)
}

// TestPageSetInputFilesShouldReturnTheSameFileChooserWhenThereAreManyWatchdogsSimultaneously verifies nested waiters get one chooser.
// Ref: TestPageSetInputFiles.java#shouldReturnTheSameFileChooserWhenThereAreManyWatchdogsSimultaneously  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSetInputFiles.java
func TestPageSetInputFilesShouldReturnTheSameFileChooserWhenThereAreManyWatchdogsSimultaneously(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	page := newPage(t)
	must.NoError(page.SetContent(ctx, "<input type=file>"))
	var inner *playwright.FileChooser
	var innerErr error
	outer, err := page.WaitForFileChooser(ctx, func() {
		inner, innerErr = page.WaitForFileChooser(ctx, func() {
			_, _ = page.EvalOnSelector(ctx, "input", "input => input.click()")
		})
	})
	must.NoError(err)
	must.NoError(innerErr)
	// The Go client builds one FileChooser per listener, so equality is checked by value (same element/page) rather than by identity.
	is.Equal(outer, inner)
}

// TestPageSetInputFilesShouldDetectMimeType is not ported.
// Ref: TestPageSetInputFiles.java#shouldDetectMimeType  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSetInputFiles.java
func TestPageSetInputFilesShouldDetectMimeType(t *testing.T) {
	t.Skip("not implemented: the Java test body is commented out (TODO: parse form fields on server)")
}

// TestPageSetInputFilesShouldBeAbleToReadSelectedFile verifies a file chosen via OnFileChooser can be read.
// Ref: TestPageSetInputFiles.java#shouldBeAbleToReadSelectedFile  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSetInputFiles.java
func TestPageSetInputFilesShouldBeAbleToReadSelectedFile(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	page := newPage(t)
	must.NoError(page.SetContent(ctx, "<input type=file>"))
	path := fileToUploadPath(t)
	off := page.OnFileChooser(func(fc *playwright.FileChooser) { _ = fc.SetFiles(ctx, []string{path}) })
	defer off()
	content, err := page.EvalOnSelector(ctx, "input", `async picker => {
  picker.click();
  await new Promise(x => picker.oninput = x);
  const reader = new FileReader();
  const promise = new Promise(fulfill => reader.onload = fulfill);
  reader.readAsText(picker.files[0]);
  return promise.then(() => reader.result);
}`)
	must.NoError(err)
	is.Equal("contents of the file", content)
}

// TestPageSetInputFilesShouldBeAbleToResetSelectedFilesWithEmptyFileList verifies an empty list resets the selection.
// Ref: TestPageSetInputFiles.java#shouldBeAbleToResetSelectedFilesWithEmptyFileList  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSetInputFiles.java
func TestPageSetInputFilesShouldBeAbleToResetSelectedFilesWithEmptyFileList(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	page := newPage(t)
	must.NoError(page.SetContent(ctx, "<input type=file>"))
	path := fileToUploadPath(t)
	const pick = `async picker => {
  picker.click();
  await new Promise(x => picker.oninput = x);
  return picker.files.length;
}`
	off := page.OnFileChooser(func(fc *playwright.FileChooser) { _ = fc.SetFiles(ctx, []string{path}) })
	n1, err := page.EvalOnSelector(ctx, "input", pick)
	must.NoError(err)
	is.EqualValues(1, n1)
	off()
	off2 := page.OnFileChooser(func(fc *playwright.FileChooser) { _ = fc.SetFiles(ctx, []string{}) })
	defer off2()
	n2, err := page.EvalOnSelector(ctx, "input", pick)
	must.NoError(err)
	is.EqualValues(0, n2)
}

// TestPageSetInputFilesShouldNotAcceptMultipleFilesForSingleFileInput verifies a non-multiple input rejects several files.
// Ref: TestPageSetInputFiles.java#shouldNotAcceptMultipleFilesForSingleFileInput  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSetInputFiles.java
func TestPageSetInputFilesShouldNotAcceptMultipleFilesForSingleFileInput(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	page := newPage(t)
	must.NoError(page.SetContent(ctx, "<input type=file>"))
	fc, err := page.WaitForFileChooser(ctx, func() { _ = page.Locator("input").Click(ctx) })
	must.NoError(err)
	err = fc.SetFiles(ctx, []string{fileToUploadPath(t), fileToUploadPath(t)})
	must.Error(err)
	is.Contains(err.Error(), "Non-multiple file input can only accept single file")
}

// a4FolderSetup opens the folder upload page and returns its input and a fresh directory path.
func a4FolderSetup(t *testing.T) (context.Context, *playwright.Locator, string) {
	t.Helper()
	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)
	a4ServeUpload(srv)
	require.NoError(t, page.Goto(ctx, srv.Prefix()+"/input/folderupload.html"))
	return ctx, page.Locator("input[name=\"file1\"]"), filepath.Join(t.TempDir(), "file-upload-test")
}

// TestPageSetInputFilesShouldUploadAFolder verifies uploading a directory.
// Ref: TestPageSetInputFiles.java#shouldUploadAFolder  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSetInputFiles.java
func TestPageSetInputFilesShouldUploadAFolder(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx, input, dir := a4FolderSetup(t)
	a4WriteText(t, filepath.Join(dir, "file1.txt"), "file1 content")
	a4WriteText(t, filepath.Join(dir, "file2"), "file2 content")
	a4WriteText(t, filepath.Join(dir, "sub-dir", "really.txt"), "sub-dir file content")
	must.NoError(input.SetInputFiles(ctx, []string{dir}))
	v, err := input.Evaluate(ctx, "e => [...e.files].map(f => f.webkitRelativePath)")
	must.NoError(err)
	var rels []string
	for _, x := range v.([]any) {
		rels = append(rels, x.(string))
	}
	sorted := append([]string(nil), rels...)
	sort.Strings(sorted)
	is.Equal([]string{"file-upload-test/file1.txt", "file-upload-test/file2", "file-upload-test/sub-dir/really.txt"}, sorted)
	for i, rel := range rels {
		content, err := input.Evaluate(ctx, `(e, i) => {
      const reader = new FileReader();
      const promise = new Promise(fulfill => reader.onload = fulfill);
      reader.readAsText(e.files[i]);
      return promise.then(() => reader.result);
    }`, i)
		must.NoError(err)
		expected, err := os.ReadFile(filepath.Join(filepath.Dir(dir), filepath.FromSlash(rel)))
		must.NoError(err)
		is.Equal(string(expected), content)
	}
}

// TestPageSetInputFilesShouldUploadAFolderAndThrowForMultipleDirectories verifies several directories are rejected.
// Ref: TestPageSetInputFiles.java#shouldUploadAFolderAndThrowForMultipleDirectories  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSetInputFiles.java
func TestPageSetInputFilesShouldUploadAFolderAndThrowForMultipleDirectories(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, input, dir := a4FolderSetup(t)
	a4WriteText(t, filepath.Join(dir, "folder1", "file1.txt"), "file1 content")
	a4WriteText(t, filepath.Join(dir, "folder2", "file2.txt"), "file2 content")
	err := input.SetInputFiles(ctx, []string{filepath.Join(dir, "folder1"), filepath.Join(dir, "folder2")})
	must.Error(err)
	must.Contains(err.Error(), "Multiple directories are not supported")
}

// TestPageSetInputFilesShouldThrowIfADirectoryAndFilesArePassed verifies mixing a directory and files is rejected.
// Ref: TestPageSetInputFiles.java#shouldThrowIfADirectoryAndFilesArePassed  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSetInputFiles.java
func TestPageSetInputFilesShouldThrowIfADirectoryAndFilesArePassed(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx, input, dir := a4FolderSetup(t)
	a4WriteText(t, filepath.Join(dir, "folder1", "file1.txt"), "file1 content")
	err := input.SetInputFiles(ctx, []string{filepath.Join(dir, "folder1"), filepath.Join(dir, "folder1", "file1.txt")})
	must.Error(err)
	must.Contains(err.Error(), "File paths must be all files or a single directory")
}

// TestPageSetInputFilesShouldThrowWhenUploadingAFolderInANormalFileUploadInput verifies a normal input rejects directories.
// Ref: TestPageSetInputFiles.java#shouldThrowWhenUploadingAFolderInANormalFileUploadInput  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestPageSetInputFiles.java
func TestPageSetInputFilesShouldThrowWhenUploadingAFolderInANormalFileUploadInput(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)
	a4ServeUpload(srv)
	must.NoError(page.Goto(ctx, srv.Prefix()+"/input/fileupload.html"))
	input := page.Locator("input[name=\"file1\"]")
	dir := filepath.Join(t.TempDir(), "file-upload-test")
	a4WriteText(t, filepath.Join(dir, "file1.txt"), "file1 content")
	err := input.SetInputFiles(ctx, []string{dir})
	must.Error(err)
	must.Contains(err.Error(), "File input does not support directories, pass individual files instead")
}
