package transport

import (
	"bytes"
	"context"
	"encoding/binary"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type nopWriteCloser struct {
	io.Writer
}

func (nopWriteCloser) Close() error {
	return nil
}

func TestPipeTransport_Send(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	reader, writer := io.Pipe()

	dummyReader := io.NopCloser(bytes.NewReader(nil))
	tr := NewPipeTransport(dummyReader, writer, nil, nil)

	msg := []byte(`{"message": "hello playwright"}`)

	errCh := make(chan error, 1)
	go func() {
		defer writer.Close()
		if err := tr.Send(msg); err != nil {
			errCh <- err
			return
		}
		errCh <- nil
	}()

	var length uint32
	err := binary.Read(reader, binary.LittleEndian, &length)
	must.NoError(err)
	is.Equal(uint32(len(msg)), length)

	buf := make([]byte, length)
	_, err = io.ReadFull(reader, buf)
	must.NoError(err)
	is.Equal(msg, buf)

	select {
	case err := <-errCh:
		must.NoError(err)
	case <-time.After(2 * time.Second):
		t.Fatal("timeout in Send")
	}
}

func TestPipeTransport_Receive(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	reader, writer := io.Pipe()

	onMessage := make(chan []byte, 1)
	dummyWriter := nopWriteCloser{bytes.NewBuffer(nil)}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	t.Cleanup(cancel)

	tr := NewPipeTransport(reader, dummyWriter, onMessage, nil)
	tr.Start(ctx)

	msg := []byte(`{"event": "test event"}`)

	errCh := make(chan error, 1)
	go func() {
		if err := binary.Write(writer, binary.LittleEndian, uint32(len(msg))); err != nil {
			errCh <- err
			return
		}
		if _, err := writer.Write(msg); err != nil {
			errCh <- err
			return
		}
		errCh <- nil
		_ = writer.Close() // signals EOF so the transport goroutine exits; success already reported
	}()

	select {
	case err := <-errCh:
		must.NoError(err)
	case <-ctx.Done():
		t.Fatal("timeout writing payload")
	}

	select {
	case received := <-onMessage:
		is.Equal(msg, received)
	case <-ctx.Done():
		t.Fatal("Timeout waiting for message to be decoded by PipeTransport")
	}
}

func TestPipeTransport_OversizedMessage(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	reader, writer := io.Pipe()
	dummyWriter := nopWriteCloser{bytes.NewBuffer(nil)}

	errCh := make(chan error, 1)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	t.Cleanup(cancel)

	tr := NewPipeTransport(reader, dummyWriter, nil, func(err error) { errCh <- err })
	tr.Start(ctx)

	go func() {
		_ = binary.Write(writer, binary.LittleEndian, MaxMessageSize+1)
		_ = writer.Close() // transport errors on the oversized header before reaching this close
	}()

	select {
	case err := <-errCh:
		is.Error(err)
	case <-ctx.Done():
		t.Fatal("Timeout waiting for oversized message error")
	}
}

func TestPipeTransport_TruncatedPayload(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	reader, writer := io.Pipe()
	dummyWriter := nopWriteCloser{bytes.NewBuffer(nil)}

	errCh := make(chan error, 1)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	t.Cleanup(cancel)

	tr := NewPipeTransport(reader, dummyWriter, nil, func(err error) { errCh <- err })
	tr.Start(ctx)

	go func() {
		_ = binary.Write(writer, binary.LittleEndian, uint32(100))
		_, _ = writer.Write([]byte("short"))
		_ = writer.Close() // transport errors on the incomplete read before reaching this close
	}()

	select {
	case err := <-errCh:
		is.Error(err)
	case <-ctx.Done():
		t.Fatal("Timeout waiting for truncated payload error")
	}
}

func TestPipeTransport_EOF(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	reader, writer := io.Pipe()
	dummyWriter := nopWriteCloser{bytes.NewBuffer(nil)}

	errCh := make(chan error, 1)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	t.Cleanup(cancel)

	tr := NewPipeTransport(reader, dummyWriter, nil, func(err error) { errCh <- err })
	tr.Start(ctx)

	_ = writer.Close() // triggers the EOF that the transport will report via errCh

	select {
	case err := <-errCh:
		is.Error(err)
	case <-ctx.Done():
		t.Fatal("Timeout waiting for EOF error")
	}
}

func TestPipeTransport_Close(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	reader, writer := io.Pipe()
	dummyWriter := nopWriteCloser{bytes.NewBuffer(nil)}

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	tr := NewPipeTransport(reader, dummyWriter, nil, nil)
	tr.Start(ctx)

	must.NotPanics(func() {
		tr.Close()
		tr.Close() // second call must be idempotent
	})
	_ = writer.Close()
}

func TestPipeTransport_NilOnMessage(t *testing.T) {
	t.Parallel()
	reader, writer := io.Pipe()
	dummyWriter := nopWriteCloser{bytes.NewBuffer(nil)}

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	errCh := make(chan error, 1)
	tr := NewPipeTransport(reader, dummyWriter, nil, func(err error) { errCh <- err })
	tr.Start(ctx)

	msg := []byte(`{"event": "ignored"}`)
	go func() {
		_ = binary.Write(writer, binary.LittleEndian, uint32(len(msg)))
		_, _ = writer.Write(msg)
		_ = writer.Close() // signals EOF so the transport goroutine exits; error irrelevant to test
	}()

	select {
	case <-errCh: // EOF expected
	case <-ctx.Done():
		t.Fatal("Timeout waiting for EOF")
	}
}
