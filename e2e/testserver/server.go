//go:build e2e

// Package testserver provides an embedded HTTP test server for Playwright-Go E2E tests.
// It is functionally equivalent to the Server.java class in the playwright-java repository.
//
// Typical usage in tests:
//
//	srv := testserver.New(t)
//	page.Goto(ctx, srv.EmptyPage())
package testserver

import (
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
)

// Server is an embedded test HTTP server.
// Provides URLs and routing utilities equivalent to the Server.java class in the reference project.
type Server struct {
	t      testing.TB
	srv    *httptest.Server
	mu     sync.RWMutex
	routes map[string]http.HandlerFunc
}

// New creates and starts an HTTP test server.
// The server registers for automatic shutdown upon test cleanup.
func New(t testing.TB) *Server {
	t.Helper()
	s := &Server{
		t:      t,
		routes: make(map[string]http.HandlerFunc),
	}
	s.srv = httptest.NewServer(http.HandlerFunc(s.handle))
	t.Cleanup(s.Close)
	return s
}

// Prefix returns the base URL prefix of the server, e.g. "http://127.0.0.1:PORT".
func (s *Server) Prefix() string {
	return s.srv.URL
}

// EmptyPage returns the URL for an empty HTML page.
func (s *Server) EmptyPage() string {
	return s.srv.URL + "/empty.html"
}

// KeyboardPage returns the URL for the keyboard test page.
func (s *Server) KeyboardPage() string {
	return s.srv.URL + "/input/keyboard.html"
}

// DomPage returns the URL for the dom.html test page.
func (s *Server) DomPage() string {
	return s.srv.URL + "/dom.html"
}

// ButtonPage returns the URL for the input/button.html test page.
func (s *Server) ButtonPage() string {
	return s.srv.URL + "/input/button.html"
}

// Port returns the numeric TCP port on which the server is listening.
func (s *Server) Port() int {
	addr := s.srv.Listener.Addr().(*net.TCPAddr)
	return addr.Port
}

// SetRoute registers an HTTP handler function for the specified path.
func (s *Server) SetRoute(path string, handler http.HandlerFunc) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.routes[path] = handler
}

// RemoveRoute deletes the registered handler for the specified path.
func (s *Server) RemoveRoute(path string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.routes, path)
}

// Reset clears all dynamically registered route handlers.
// Equivalent to server.reset() in Java.
func (s *Server) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.routes = make(map[string]http.HandlerFunc)
}

// Close stops the embedded server. It is idempotent.
func (s *Server) Close() {
	s.srv.Close()
}

// handle is the main HandlerFunc for the test server.
// Serves dynamically registered routes and provides built-in static mock pages.
func (s *Server) handle(w http.ResponseWriter, r *http.Request) {
	// 1. Check dynamically registered routes
	s.mu.RLock()
	handler, ok := s.routes[r.URL.Path]
	s.mu.RUnlock()

	if ok {
		handler(w, r)
		return
	}

	// 2. Built-in static pages
	switch r.URL.Path {
	case "/empty.html":
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		io.WriteString(w, "<!DOCTYPE html><html><head></head><body></body></html>\n")

	case "/title.html":
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		io.WriteString(w, "<!DOCTYPE html><html><head><title>Woof-Woof</title></head><body></body></html>\n")

	case "/beforeunload.html":
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		io.WriteString(w, `<!DOCTYPE html><html><body><script>
window.addEventListener('beforeunload', e => { e.returnValue = 'Leave?'; });
</script></body></html>`)

	case "/input/textarea.html":
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		io.WriteString(w, `<!DOCTYPE html><html><body><textarea></textarea></body></html>`)

	case "/input/keyboard.html":
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		io.WriteString(w, `<!DOCTYPE html><html><head><title>Keyboard test</title></head><body><textarea></textarea><script>
window.result = "";
let textarea = document.querySelector('textarea');
textarea.focus();
textarea.addEventListener('keydown', event => { log('Keydown:', event.key, event.code, event.which, modifiers(event)); });
textarea.addEventListener('keypress', event => { log('Keypress:', event.key, event.code, event.which, event.charCode, modifiers(event)); });
textarea.addEventListener('keyup', event => { log('Keyup:', event.key, event.code, event.which, modifiers(event)); });
function modifiers(event) {
  let m = [];
  if (event.altKey) m.push('Alt');
  if (event.ctrlKey) m.push('Control');
  if (event.shiftKey) m.push('Shift');
  return '[' + m.join(' ') + ']';
}
function log(...args) { result += args.join(' ') + '\n'; }
function getResult() { let temp = result.trim(); result = ""; return temp; }
</script></body></html>`)

	case "/dom.html":
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		io.WriteString(w, `<div id="outer" name="value"><div id="inner">Text,`+"\n"+`more text</div></div><input id="check" type=checkbox checked foo="bar&quot;"><input id="input"></input><textarea id="textarea"></textarea><select id="select"><option></option><option value="foo"></option></select>`)

	case "/input/button.html":
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		io.WriteString(w, `<!DOCTYPE html><html><head><title>Button test</title></head><body><button>Click target</button><script>
window.result = 'Was not clicked';
window.offsetX = undefined;
window.offsetY = undefined;
window.pageX = undefined;
window.pageY = undefined;
window.shiftKey = undefined;
window.bubbles = undefined;
document.querySelector('button').addEventListener('click', e => {
  result = 'Clicked';
  offsetX = e.offsetX;
  offsetY = e.offsetY;
  pageX = e.pageX;
  pageY = e.pageY;
  shiftKey = e.shiftKey;
  bubbles = e.bubbles;
  cancelable = e.cancelable;
  composed = e.composed;
}, false);
</script></body></html>`)

	case "/drag-n-drop.html":
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		io.WriteString(w, `<!DOCTYPE html><html><body>
<div id="source" draggable="true">source</div>
<div id="target">target</div>
<script>
document.querySelector('#source').addEventListener('dragstart', e => e.dataTransfer.setData('text', 'dragged'));
document.querySelector('#target').addEventListener('dragover', e => e.preventDefault());
document.querySelector('#target').addEventListener('drop', e => { e.preventDefault(); e.target.appendChild(document.querySelector('#source')); });
</script></body></html>`)

	default:
		http.NotFound(w, r)
		s.t.Logf("testserver: 404 for path %s", r.URL.Path)
	}
}

// WaitForRequest returns a channel that receives an incoming request and blocks the handler
// until inspected, equivalent to server.waitForRequest() in Java.
func WaitForRequest(s *Server, path string) <-chan *http.Request {
	ch := make(chan *http.Request, 1)
	s.SetRoute(path, func(w http.ResponseWriter, r *http.Request) {
		ch <- r
		w.WriteHeader(http.StatusOK)
	})
	return ch
}

// ServeWithBody serves a path with a specific content-type and body string.
func (s *Server) ServeWithBody(path, contentType, body string) {
	s.SetRoute(path, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", contentType)
		io.WriteString(w, body)
	})
}

// ServeWithRedirect serves a path that redirects to another URL.
func (s *Server) ServeWithRedirect(from, to string) {
	s.SetRoute(from, func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, to, http.StatusFound)
	})
}

// BlockRequest blocks an incoming request indefinitely until the client cancels.
func (s *Server) BlockRequest(path string) {
	s.SetRoute(path, func(w http.ResponseWriter, r *http.Request) {
		<-r.Context().Done()
	})
}

// RequestReceived returns a channel that is closed when the first HTTP request
// arrives at the given path. The request remains blocked until the client cancels.
func (s *Server) RequestReceived(path string) <-chan struct{} {
	ch := make(chan struct{})
	var once sync.Once
	s.SetRoute(path, func(w http.ResponseWriter, r *http.Request) {
		once.Do(func() { close(ch) })
		<-r.Context().Done()
	})
	return ch
}

// ServeWithBasicAuth serves a path protected by HTTP Basic Auth.
// Requests without valid credentials receive a 401 WWW-Authenticate challenge.
// Requests with valid credentials receive the given contentType and body.
func (s *Server) ServeWithBasicAuth(path, username, password, contentType, body string) {
	s.SetRoute(path, func(w http.ResponseWriter, r *http.Request) {
		u, p, ok := r.BasicAuth()
		if !ok || u != username || p != password {
			w.Header().Set("WWW-Authenticate", `Basic realm="test"`)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", contentType)
		io.WriteString(w, body)
	})
}

// HasHTTPSSupport returns false as this test server is HTTP-only.
func (s *Server) HasHTTPSSupport() bool { return false }

// CrossProcessPrefix returns the base URL replacing 127.0.0.1 with localhost (or vice versa)
// to simulate cross-origin requests (equivalent to Java's CROSS_PROCESS_PREFIX).
func (s *Server) CrossProcessPrefix() string {
	if strings.Contains(s.Prefix(), "127.0.0.1") {
		return strings.Replace(s.Prefix(), "127.0.0.1", "localhost", 1)
	}
	return strings.Replace(s.Prefix(), "localhost", "127.0.0.1", 1)
}
