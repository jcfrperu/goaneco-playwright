package driver

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"time"

	"github.com/jcfrperu/goaneco-playwright/internal/connection"
	"github.com/jcfrperu/goaneco-playwright/internal/transport"
)

// Driver encapsulates the Node.js playwright-cli process.
type Driver struct {
	cmd       *exec.Cmd
	transport *transport.PipeTransport
	conn      *connection.Connection
	cancel    context.CancelFunc
	done      chan struct{}
	waitErr   error
}

// StartDriver launches the playwright driver using Node and wires its IPC pipes to a Connection.
func StartDriver(ctx context.Context, cliPath string) (*Driver, error) {
	driverCtx, cancel := context.WithCancel(ctx)

	nodePath := os.Getenv("PLAYWRIGHT_NODEJS_PATH")
	if nodePath == "" {
		nodePath = "node"
	}
	cmd := exec.Command(nodePath, cliPath, "run-driver")

	stderr, err := cmd.StderrPipe()
	if err != nil {
		cancel()
		return nil, fmt.Errorf("could not get stderr pipe: %w", err)
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		_ = stderr.Close() // primary error already captured; best-effort pipe cleanup
		cancel()
		return nil, fmt.Errorf("could not get stdin pipe: %w", err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		_ = stderr.Close() // primary error already captured; best-effort pipe cleanup
		_ = stdin.Close()
		cancel()
		return nil, fmt.Errorf("could not get stdout pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		_ = stderr.Close() // primary error already captured; best-effort pipe cleanup
		_ = stdin.Close()
		_ = stdout.Close()
		cancel()
		return nil, fmt.Errorf("could not start playwright driver: %w", err)
	}

	stderrDone := make(chan struct{})
	go func() {
		defer close(stderrDone)
		// No driverCtx needed: the OS closes the stderr pipe when the process dies,
		// which causes scanner.Scan to return false. The watchdog waits on stderrDone
		// before closing onMessage, so no messages are dropped.
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			slog.Default().Info("playwright-driver", "msg", scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			slog.Default().Error("playwright-driver: stderr scanner error", "error", err)
		}
	}()

	// dispatchBufSize buffers IPC messages between the transport reader and the dispatcher.
	// Playwright can emit bursts of __create__ events during page load; 256 provides headroom.
	const dispatchBufSize = 256
	onMessage := make(chan []byte, dispatchBufSize)
	conn := connection.NewConnection()

	onError := func(err error) {
		if driverCtx.Err() != nil || errors.Is(err, io.EOF) || errors.Is(err, os.ErrClosed) {
			return // normal shutdown
		}
		slog.Default().Error("transport error", "error", err)
	}
	tr := transport.NewPipeTransport(stdout, stdin, onMessage, onError)
	conn.SetTransportSend(tr.Send)
	tr.Start(driverCtx)

	done := make(chan struct{})

	d := &Driver{
		cmd:       cmd,
		transport: tr,
		conn:      conn,
		cancel:    cancel,
		done:      done,
	}

	// Watchdog: closes transport and reaps the process; signals done when finished.
	go func() {
		defer close(done)

		d.waitErr = cmd.Wait()
		if d.waitErr != nil && driverCtx.Err() == nil {
			slog.Default().Error("playwright driver process exited with error", "error", d.waitErr)
		}

		cancel()   // Cancel context so reader and dispatcher goroutines initiate exit
		tr.Close() // Close pipes
		tr.Wait()  // Ensure transport reader goroutine has completely exited before closing channel
		<-stderrDone
		close(onMessage)
	}()

	// Dispatcher: routes incoming IPC messages from the transport into the connection.
	// Exits only when onMessage is closed by the watchdog, ensuring no buffered
	// responses are dropped when Stop() cancels the context.
	// After draining, conn.Close unblocks any goroutines waiting in SendRequest.
	go func() {
		for msg := range onMessage {
			conn.Dispatch(msg)
		}
		conn.Close(d.waitErr)
	}()

	return d, nil
}

// Conn returns the underlying Connection used for IPC with the playwright driver.
func (d *Driver) Conn() *connection.Connection {
	return d.conn
}

// Stop cancels the driver context and waits up to 5 seconds for the process to exit
// gracefully before sending a kill signal.
func (d *Driver) Stop() error {
	d.cancel()
	d.transport.Close() // unblocks binary.Read immediately

	timer := time.NewTimer(5 * time.Second)
	defer timer.Stop()
	select {
	case <-d.done:
		return d.waitErr
	case <-timer.C:
	}

	if d.cmd != nil && d.cmd.Process != nil {
		// Ignore Kill syscall error — process exit status from cmd.Wait is authoritative.
		_ = d.cmd.Process.Kill()
		// Wait for watchdog so d.waitErr is populated by cmd.Wait().
		<-d.done
	}
	return d.waitErr
}
