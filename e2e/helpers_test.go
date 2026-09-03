//go:build e2e

package e2e

import (
	"bytes"
	"sync"
)

// threadSafeWriter is a concurrency-safe wrapper around bytes.Buffer for testing log outputs.
type threadSafeWriter struct {
	buf *bytes.Buffer
	mu  *sync.Mutex
}

func (w *threadSafeWriter) Write(p []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.buf.Write(p)
}
