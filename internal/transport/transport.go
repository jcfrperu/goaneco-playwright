package transport

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"sync"
)

// MaxMessageSize limits incoming message allocation to 256 MiB, matching Playwright's own limit.
const MaxMessageSize uint32 = 256 * 1024 * 1024

// PipeTransport implements the Playwright IPC framing protocol over a pair of OS pipes.
// Each message is prefixed with a 4-byte little-endian length, followed by the JSON payload.
// Incoming messages are delivered to onMessage; read errors are forwarded to onError.
type PipeTransport struct {
	streamReader io.ReadCloser
	streamWriter io.WriteCloser
	mu           sync.Mutex
	once         sync.Once
	wg           sync.WaitGroup
	onMessage    chan<- []byte
	onError      func(error)
}

// NewPipeTransport creates a PipeTransport backed by reader (stdin of the driver process) and
// writer (stdout of the driver process). Incoming decoded messages are sent to onMessage;
// any read error is forwarded to onError. Call Start to begin the read loop.
func NewPipeTransport(reader io.ReadCloser, writer io.WriteCloser, onMessage chan<- []byte, onError func(error)) *PipeTransport {
	return &PipeTransport{
		streamReader: reader,
		streamWriter: writer,
		onMessage:    onMessage,
		onError:      onError,
	}
}

// Start launches the read loop in a background goroutine. The loop reads length-prefixed
// messages from the pipe until an error occurs or ctx is canceled, then calls Close.
// Call Wait to block until the goroutine has exited.
func (t *PipeTransport) Start(ctx context.Context) {
	t.wg.Add(1)
	go func() {
		defer t.wg.Done()
		defer t.Close()
		for {
			var length uint32
			if err := binary.Read(t.streamReader, binary.LittleEndian, &length); err != nil {
				t.notifyError(err)
				return
			}

			if length > MaxMessageSize {
				t.notifyError(fmt.Errorf("transport: message length %d exceeds MaxMessageSize %d", length, MaxMessageSize))
				return
			}

			buf := make([]byte, length)
			if _, err := io.ReadFull(t.streamReader, buf); err != nil {
				t.notifyError(err)
				return
			}

			if t.onMessage != nil {
				select {
				case t.onMessage <- buf:
				case <-ctx.Done():
					return
				}
			}
		}
	}()
}

// Send writes a length-prefixed message to the pipe. The 4-byte little-endian length header
// is written atomically with the payload under a mutex to prevent interleaving from concurrent callers.
func (t *PipeTransport) Send(message []byte) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if err := binary.Write(t.streamWriter, binary.LittleEndian, uint32(len(message))); err != nil {
		return err
	}
	_, err := t.streamWriter.Write(message)
	return err
}

// Close closes both the reader and writer pipes exactly once (guarded by sync.Once).
// It is safe to call from multiple goroutines and is idempotent.
func (t *PipeTransport) Close() {
	t.once.Do(func() {
		// Close() has no return value; errors from the underlying streams cannot be propagated.
		_ = t.streamReader.Close()
		_ = t.streamWriter.Close()
	})
}

// Wait blocks until the read loop goroutine has finished.
func (t *PipeTransport) Wait() {
	t.wg.Wait()
}

func (t *PipeTransport) notifyError(err error) {
	if t.onError != nil {
		t.onError(err)
	}
}
