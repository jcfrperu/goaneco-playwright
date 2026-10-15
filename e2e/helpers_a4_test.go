//go:build e2e

package e2e

import (
	"bytes"
	"compress/gzip"
	"context"
	"image"
	"image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/require"
)

// a4Ptr returns a pointer to v.
func a4Ptr[T any](v T) *T { return &v }

// a4NewContext creates a context with the given options and registers cleanup.
func a4NewContext(t *testing.T, opts ...*playwright.BrowserContextOptions) *playwright.BrowserContext {
	t.Helper()
	bCtx, err := globalBrowser.NewContext(testCtx(t), opts...)
	require.NoError(t, err, "NewContext failed")
	t.Cleanup(func() {
		c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = bCtx.Close(c)
	})
	return bCtx
}

// a4NewPage creates a page in a fresh context with the given options and registers cleanup.
func a4NewPage(t *testing.T, opts ...*playwright.BrowserContextOptions) *playwright.Page {
	t.Helper()
	bCtx := a4NewContext(t, opts...)
	page, err := bCtx.NewPage(testCtx(t))
	require.NoError(t, err, "NewPage failed")
	return page
}

// a4Navigate navigates page to url and returns the status of the first response whose URL equals url.
func a4Navigate(t *testing.T, page *playwright.Page, url string) (int, []byte) {
	t.Helper()
	ctx := testCtx(t)
	type res struct {
		status int
		body   []byte
	}
	ch := make(chan res, 4)
	off := page.OnResponse(func(r *playwright.NetworkResponse) {
		if r.URL() == url {
			b, _ := r.Body(ctx)
			select {
			case ch <- res{r.Status(), b}:
			default:
			}
		}
	})
	defer off()
	_ = page.Goto(ctx, url)
	select {
	case r := <-ch:
		return r.status, r.body
	case <-time.After(10 * time.Second):
		t.Fatalf("no response observed for %s", url)
	}
	return 0, nil
}

// a4NewAPICtx creates a standalone APIRequestContext and registers cleanup.
func a4NewAPICtx(t *testing.T, opts ...*playwright.APIRequestContextOptions) *playwright.APIRequestContext {
	t.Helper()
	apiCtx, err := globalPW.NewAPIRequestContext(testCtx(t), opts...)
	require.NoError(t, err, "NewAPIRequestContext")
	t.Cleanup(func() {
		c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = apiCtx.Dispose(c)
	})
	return apiCtx
}

// a4Gzip serves body gzip-compressed at path.
func a4Gzip(srv *testserver.Server, path, contentType string, body []byte) {
	srv.SetRoute(path, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", contentType)
		w.Header().Set("Content-Encoding", "gzip")
		zw := gzip.NewWriter(w)
		_, _ = zw.Write(body)
		_ = zw.Close()
	})
}

// a4Log records network events of a page.
type a4Log struct {
	mu       sync.Mutex
	reqs     []*playwright.NetworkRequest
	resps    []*playwright.NetworkResponse
	finished []*playwright.NetworkRequest
	failed   []*playwright.NetworkRequest
}

// a4Watch starts recording network events of the page.
func a4Watch(t *testing.T, page *playwright.Page) *a4Log {
	l := &a4Log{}
	offs := []func(){
		page.OnRequest(func(r *playwright.NetworkRequest) { l.mu.Lock(); l.reqs = append(l.reqs, r); l.mu.Unlock() }),
		page.OnResponse(func(r *playwright.NetworkResponse) { l.mu.Lock(); l.resps = append(l.resps, r); l.mu.Unlock() }),
		page.OnRequestFinished(func(r *playwright.NetworkRequest) { l.mu.Lock(); l.finished = append(l.finished, r); l.mu.Unlock() }),
		page.OnRequestFailed(func(r *playwright.NetworkRequest) { l.mu.Lock(); l.failed = append(l.failed, r); l.mu.Unlock() }),
	}
	t.Cleanup(func() {
		for _, off := range offs {
			off()
		}
	})
	return l
}

// Requests returns a snapshot of requests seen so far.
func (l *a4Log) Requests() []*playwright.NetworkRequest {
	l.mu.Lock()
	defer l.mu.Unlock()
	return append([]*playwright.NetworkRequest(nil), l.reqs...)
}

// Finished returns a snapshot of finished requests seen so far.
func (l *a4Log) Finished() []*playwright.NetworkRequest {
	l.mu.Lock()
	defer l.mu.Unlock()
	return append([]*playwright.NetworkRequest(nil), l.finished...)
}

// Failed returns a snapshot of failed requests seen so far.
func (l *a4Log) Failed() []*playwright.NetworkRequest {
	l.mu.Lock()
	defer l.mu.Unlock()
	return append([]*playwright.NetworkRequest(nil), l.failed...)
}

// WaitResponse waits for a response whose URL contains substr (and optional status when > 0).
func (l *a4Log) WaitResponse(t *testing.T, substr string, status int) *playwright.NetworkResponse {
	t.Helper()
	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		l.mu.Lock()
		for _, r := range l.resps {
			if strings.Contains(r.URL(), substr) && (status == 0 || r.Status() == status) {
				l.mu.Unlock()
				return r
			}
		}
		l.mu.Unlock()
		time.Sleep(20 * time.Millisecond)
	}
	t.Fatalf("response containing %q not observed", substr)
	return nil
}

// WaitRequest waits for a request whose URL contains substr.
func (l *a4Log) WaitRequest(t *testing.T, substr string) *playwright.NetworkRequest {
	t.Helper()
	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		l.mu.Lock()
		for _, r := range l.reqs {
			if strings.Contains(r.URL(), substr) {
				l.mu.Unlock()
				return r
			}
		}
		l.mu.Unlock()
		time.Sleep(20 * time.Millisecond)
	}
	t.Fatalf("request containing %q not observed", substr)
	return nil
}

// a4Rec records the first request received on a server path.
type a4Rec struct {
	got    chan struct{}
	once   sync.Once
	Method string
	Body   []byte
	Header http.Header
}

// a4RecordRoute registers a route at path that records the request and replies 200 (with CORS enabled).
func a4RecordRoute(srv *testserver.Server, path string) *a4Rec {
	r := &a4Rec{got: make(chan struct{})}
	srv.SetRoute(path, func(w http.ResponseWriter, req *http.Request) {
		b, _ := io.ReadAll(req.Body)
		r.once.Do(func() {
			r.Method = req.Method
			r.Body = b
			r.Header = req.Header.Clone()
			close(r.got)
		})
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "text/html")
		_, _ = w.Write([]byte("ok"))
	})
	return r
}

// Wait blocks until the request was recorded.
func (r *a4Rec) Wait(t *testing.T) {
	t.Helper()
	select {
	case <-r.got:
	case <-time.After(10 * time.Second):
		t.Fatal("server did not receive the request")
	}
}

// a4Bytes256 returns the bytes 0..255.
func a4Bytes256() []byte {
	b := make([]byte, 256)
	for i := range b {
		b[i] = byte(i)
	}
	return b
}

// a4PNGSize returns the dimensions of a PNG image.
func a4PNGSize(t *testing.T, data []byte) (int, int) {
	t.Helper()
	cfg, err := png.DecodeConfig(bytes.NewReader(data))
	require.NoError(t, err)
	return cfg.Width, cfg.Height
}

var _ = image.Rect

// a4ServeGrid serves /grid.html, a simple grid of coloured divs with an image.
func a4ServeGrid(srv *testserver.Server) {
	var sb strings.Builder
	sb.WriteString("<style>body{margin:0}div{display:inline-block;width:50px;height:50px;border:1px solid #333;box-sizing:border-box;background:#9cf}</style>")
	for i := 0; i < 100; i++ {
		sb.WriteString("<div><img width=20 height=20 src=\"data:image/gif;base64,R0lGODlhAQABAIAAAP8AAAAAACH5BAAAAAAALAAAAAABAAEAAAICRAEAOw==\"></div>")
	}
	srv.ServeWithBody("/grid.html", "text/html", sb.String())
}

// a4ServeSelect serves /input/select.html, a colour select with input/change tracking.
func a4ServeSelect(srv *testserver.Server) {
	colors := []string{"black", "blue", "brown", "cyan", "gray", "green", "indigo", "magenta", "orange", "pink", "red", "violet", "white", "yellow"}
	var sb strings.Builder
	sb.WriteString("<!DOCTYPE html><html><head><title>Selection Test</title></head><body><select>")
	for _, c := range colors {
		sb.WriteString("<option id=\"" + c + "Option\" value=\"" + c + "\">" + strings.ToUpper(c[:1]) + c[1:] + "</option>")
	}
	sb.WriteString(`</select><script>
window.result = {onInput: null, onChange: null};
window.makeMultiple = () => document.querySelector('select').setAttribute('multiple', true);
const select = document.querySelector('select');
select.addEventListener('input', () => { result.onInput = Array.from(select.querySelectorAll('option:checked')).map(o => o.value); });
select.addEventListener('change', () => { result.onChange = Array.from(select.querySelectorAll('option:checked')).map(o => o.value); });
</script></body></html>`)
	srv.ServeWithBody("/input/select.html", "text/html", sb.String())
}

const a4ReadFileJS = `e => {
  const reader = new FileReader();
  const promise = new Promise(fulfill => reader.onload = fulfill);
  reader.readAsText(e.files[0]);
  return promise.then(() => reader.result);
}`

// a4WriteBig writes size bytes of 'A' to path.
func a4WriteBig(t *testing.T, path string, size int) {
	t.Helper()
	f, err := os.Create(path)
	require.NoError(t, err)
	defer f.Close()
	chunk := bytes.Repeat([]byte("A"), 1<<20)
	for written := 0; written < size; written += len(chunk) {
		n := len(chunk)
		if size-written < n {
			n = size - written
		}
		_, err := f.Write(chunk[:n])
		require.NoError(t, err)
	}
}

// a4WriteText writes content to path, creating parent directories.
func a4WriteText(t *testing.T, path, content string) {
	t.Helper()
	require.NoError(t, os.MkdirAll(filepath.Dir(path), 0o755))
	require.NoError(t, os.WriteFile(path, []byte(content), 0o644))
}

// a4ServeUpload serves the file/folder upload pages used by the set-input-files tests.
func a4ServeUpload(srv *testserver.Server) {
	form := func(attrs string) string {
		return `<form action="/upload" method="post" enctype="multipart/form-data"><input type="file" name="file1" ` + attrs + `><input type="submit"></form>`
	}
	srv.ServeWithBody("/input/fileupload.html", "text/html", form(""))
	srv.ServeWithBody("/input/folderupload.html", "text/html", form("webkitdirectory"))
	a4ServeUploadMulti(srv)
}

// a4ServeUploadMulti serves the multi-file upload page.
func a4ServeUploadMulti(srv *testserver.Server) {
	srv.ServeWithBody("/input/fileupload-multi.html", "text/html", `<input type="file" multiple>`)
}

// a4UploadedFile describes one uploaded multipart file.
type a4UploadedFile struct {
	Name string
	Size int64
}

// a4Upload collects the files posted to /upload.
type a4Upload struct{ ch chan []a4UploadedFile }

// a4UploadRoute registers /upload and records the multipart files.
func a4UploadRoute(srv *testserver.Server) *a4Upload {
	u := &a4Upload{ch: make(chan []a4UploadedFile, 1)}
	srv.SetRoute("/upload", func(w http.ResponseWriter, r *http.Request) {
		mr, err := r.MultipartReader()
		var files []a4UploadedFile
		if err == nil {
			for {
				part, err := mr.NextPart()
				if err != nil {
					break
				}
				n, _ := io.Copy(io.Discard, part)
				files = append(files, a4UploadedFile{Name: part.FileName(), Size: n})
			}
		}
		select {
		case u.ch <- files:
		default:
		}
		w.WriteHeader(200)
	})
	return u
}

// Wait returns the uploaded files.
func (u *a4Upload) Wait(t *testing.T, d time.Duration) []a4UploadedFile {
	t.Helper()
	select {
	case f := <-u.ch:
		return f
	case <-time.After(d):
		t.Fatal("upload not received")
		return nil
	}
}
