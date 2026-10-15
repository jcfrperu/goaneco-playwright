package connection

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConnection_Dispatch(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	conn := NewConnection()
	id := 42

	conn.SetTransportSend(func(payload []byte) error {
		go func() {
			msg := []byte(fmt.Sprintf(`{"id": %d, "method": "success", "params": {}}`, id))
			conn.Dispatch(msg)
		}()
		return nil
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	t.Cleanup(cancel)
	res, err := conn.SendRequest(ctx, id, []byte(`{}`))
	must.NoError(err)

	must.NotNil(res.ID)
	is.Equal(id, *res.ID)
	is.Equal("success", res.Method)
}

func TestConnection_Dispatch_NilID(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	conn := NewConnection()
	msg := []byte(`{"guid": "abc", "method": "event", "params": {}}`)
	must.NotPanics(func() { conn.Dispatch(msg) })
}

func TestConnection_Dispatch_InvalidJSON(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	conn := NewConnection()
	must.NotPanics(func() { conn.Dispatch([]byte("not valid json")) })
}

func TestConnection_Dispatch_UnknownID(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	conn := NewConnection()
	msg := []byte(`{"id": 999, "method": "orphan", "params": {}}`)
	must.NotPanics(func() { conn.Dispatch(msg) })
}

func TestConnection_SendRequest_Timeout(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	conn := NewConnection()

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately so SendRequest returns with context error

	conn.SetTransportSend(func([]byte) error { return nil })
	_, err := conn.SendRequest(ctx, 1, []byte(`{}`))
	is.ErrorIs(err, context.Canceled)
}

func TestConnection_SendRequest_Concurrent(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	conn := NewConnection()
	const n = 20

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	t.Cleanup(cancel)
	var wg sync.WaitGroup
	results := make([]Message, n)
	conn.SetTransportSend(func(payload []byte) error {
		var req struct {
			ID int `json:"id"`
		}
		if err := json.Unmarshal(payload, &req); err != nil {
			t.Errorf("failed to unmarshal transport payload: %v", err)
			return nil
		}
		go conn.Dispatch([]byte(fmt.Sprintf(`{"id":%d,"method":"ok","params":{}}`, req.ID)))
		return nil
	})

	for i := range n {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			msg, err := conn.SendRequest(ctx, id, []byte(fmt.Sprintf(`{"id":%d}`, id)))
			if err != nil {
				t.Errorf("expected no error, got: %v", err)
				return
			}
			results[id] = msg
		}(i)
	}

	wg.Wait()

	for i := range n {
		must.NotNil(results[i].ID, "result %d should have an ID", i)
		is.Equal(i, *results[i].ID)
		is.Equal("ok", results[i].Method)
	}
}

func TestConnection_OnEvent_OffEvent(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	conn := NewConnection()

	var eventCount int
	var mu sync.Mutex

	handler := func(payload json.RawMessage) {
		mu.Lock()
		defer mu.Unlock()
		eventCount++
	}

	id := conn.OnEvent("test-guid", "test-method", handler)

	// Dispatch event
	msg := []byte(`{"guid":"test-guid", "method":"test-method", "params":{}}`)
	conn.Dispatch(msg)
	conn.WaitDispatch()

	mu.Lock()
	is.Equal(1, eventCount)
	mu.Unlock()

	// OffEvent
	conn.OffEvent("test-guid", "test-method", id)

	// Verify key is removed from map when no handlers remain
	conn.mu.Lock()
	_, exists := conn.eventListeners[eventKey{GUID: "test-guid", Method: "test-method"}]
	conn.mu.Unlock()
	is.False(exists, "OffEvent should remove key from map when no handlers remain")

	// Dispatch again, should not trigger handler
	conn.Dispatch(msg)
	conn.WaitDispatch()

	mu.Lock()
	is.Equal(1, eventCount)
	mu.Unlock()
}

// TestConnection_OffEvent_DeletesMapKey verifies that when the last handler for a key is removed with OffEvent,
// dispatching the same event does not call any handler.
func TestConnection_OffEvent_DeletesMapKey(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	conn := NewConnection()

	var callCountA, callCountB int
	var mu sync.Mutex

	handlerA := func(_ json.RawMessage) { mu.Lock(); callCountA++; mu.Unlock() }
	handlerB := func(_ json.RawMessage) { mu.Lock(); callCountB++; mu.Unlock() }

	idA := conn.OnEvent("guid-x", "evt", handlerA)
	idB := conn.OnEvent("guid-x", "evt", handlerB)

	msg := []byte(`{"guid":"guid-x","method":"evt","params":{}}`)

	// Both handlers should receive the event
	conn.Dispatch(msg)
	conn.WaitDispatch()
	mu.Lock()
	is.Equal(1, callCountA)
	is.Equal(1, callCountB)
	mu.Unlock()

	// Remove handlerA: slice still has handlerB, key should remain
	conn.OffEvent("guid-x", "evt", idA)
	conn.Dispatch(msg)
	conn.WaitDispatch()
	mu.Lock()
	is.Equal(1, callCountA, "handlerA should not be called after OffEvent")
	is.Equal(2, callCountB, "handlerB should continue receiving events")
	mu.Unlock()

	// Remove handlerB: empty slice -> key must be deleted from map
	conn.OffEvent("guid-x", "evt", idB)
	conn.Dispatch(msg)
	conn.WaitDispatch()
	mu.Lock()
	is.Equal(1, callCountA, "handlerA should not be called")
	is.Equal(2, callCountB, "handlerB should not be called after OffEvent")
	mu.Unlock()
}

func TestConnection_CreateAndDispose(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	conn := NewConnection()

	// Dispatch __create__
	createMsg := []byte(`{"method":"__create__", "params":{"guid":"obj1", "type":"TestType", "initializer":{}}}`)
	conn.Dispatch(createMsg)

	obj, ok := conn.GetObject("obj1")
	is.True(ok)
	is.Equal("obj1", obj.GUID())

	// Register event listener on obj1
	called := false
	conn.OnEvent("obj1", "someEvent", func(_ json.RawMessage) {
		called = true
	})

	// Dispatch __dispose__
	disposeMsg := []byte(`{"method":"__dispose__", "guid":"obj1"}`)
	conn.Dispatch(disposeMsg)

	_, ok = conn.GetObject("obj1")
	is.False(ok)

	// Verify event listener map key for obj1 was removed and dispatching event does not trigger callback
	conn.mu.Lock()
	_, exists := conn.eventListeners[eventKey{GUID: "obj1", Method: "someEvent"}]
	conn.mu.Unlock()
	is.False(exists, "event listener key should be deleted on __dispose__")

	conn.Dispatch([]byte(`{"guid":"obj1", "method":"someEvent", "params":{}}`))
	conn.WaitDispatch()
	is.False(called, "event listener should not be called after __dispose__")
}

func TestConnection_WaitPlaywright(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	conn := NewConnection()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	t.Cleanup(cancel)

	go func() {
		// Simulate late __create__ for Playwright object
		createMsg := []byte(`{"method":"__create__", "params":{"guid":"pw-guid", "type":"Playwright", "initializer":{}}}`)
		conn.Dispatch(createMsg)
	}()

	obj, err := conn.WaitPlaywright(ctx)
	must.NoError(err)
	must.NotNil(obj)
	is.Equal("pw-guid", obj.GUID())
}

func TestConnection_SendRequest_NilTransport(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	conn := NewConnection()
	// transportSend not initialized

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	t.Cleanup(cancel)

	_, err := conn.SendRequest(ctx, 1, []byte(`{}`))
	must.EqualError(err, "transport send function not initialized")
}

// TestConnection_FailedSendDoesNotOrphanCallback (UNIT-IPC-01)
// Verifies that if transportSend returns an error, the callback channel is removed
// immediately from the callbacks map and does not leak in memory.
func TestConnection_FailedSendDoesNotOrphanCallback(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	conn := NewConnection()
	expectedErr := fmt.Errorf("transport write broken pipe")
	conn.SetTransportSend(func([]byte) error {
		return expectedErr
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	t.Cleanup(cancel)

	_, err := conn.SendRequest(ctx, 101, []byte(`{}`))
	is.ErrorIs(err, expectedErr)

	// Verify callback id was removed from callbacks map
	conn.mu.Lock()
	_, exists := conn.callbacks[101]
	callbackCount := len(conn.callbacks)
	conn.mu.Unlock()

	is.False(exists, "failed SendRequest must not retain its callback ID in callbacks map")
	is.Equal(0, callbackCount, "callbacks map must be empty after failed send")
}

// TestConnection_EventHandlerReentrancy (UNIT-EVT-01)
// Verifies that registering new handlers or removing handlers from within an event callback
// does not cause a deadlock on Connection mutex.
func TestConnection_EventHandlerReentrancy(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	conn := NewConnection()
	done := make(chan struct{})

	var nestedID ListenerID
	var id ListenerID

	id = conn.OnEvent("reentrant-guid", "eventA", func(p json.RawMessage) {
		// Reentrant call: register a new handler for another event
		nestedID = conn.OnEvent("reentrant-guid", "eventB", func(p2 json.RawMessage) {})
		// Reentrant call: unregister itself
		conn.OffEvent("reentrant-guid", "eventA", id)
		close(done)
	})

	// Dispatch event
	msg := []byte(`{"guid":"reentrant-guid","method":"eventA","params":{}}`)
	conn.Dispatch(msg)

	select {
	case <-done:
		// Succeeded without deadlock
	case <-time.After(2 * time.Second):
		t.Fatal("event handler that re-entered Connection deadlocked")
	}

	must.NotZero(nestedID)

	// Verify eventA listener was removed
	conn.mu.Lock()
	_, existsA := conn.eventListeners[eventKey{GUID: "reentrant-guid", Method: "eventA"}]
	handlersB := conn.eventListeners[eventKey{GUID: "reentrant-guid", Method: "eventB"}]
	conn.mu.Unlock()

	is.False(existsA, "eventA handler should be removed")
	is.Len(handlersB, 1, "eventB handler should be registered")
}

// TestConnection_OnceHandlerStrictlyOnce (UNIT-EVT-01)
// Verifies that a one-shot handler runs strictly once even when multiple goroutines
// dispatch events concurrently.
func TestConnection_OnceHandlerStrictlyOnce(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	conn := NewConnection()

	var count int
	var countMu sync.Mutex

	handler := func(_ json.RawMessage) {
		countMu.Lock()
		count++
		countMu.Unlock()
	}

	conn.OnEventOnce("target-guid", "onceEvent", handler)

	const concurrentDispatches = 50
	var wg sync.WaitGroup
	wg.Add(concurrentDispatches)

	for range concurrentDispatches {
		go func() {
			defer wg.Done()
			conn.Dispatch([]byte(`{"guid":"target-guid","method":"onceEvent","params":{}}`))
		}()
	}

	wg.Wait()

	countMu.Lock()
	finalCount := count
	countMu.Unlock()

	// Handler registered with OnEventOnce must have executed exactly once
	is.Equal(1, finalCount, "one-shot handler should execute strictly once")

	// Verify listener is cleaned up
	conn.mu.Lock()
	_, exists := conn.eventListeners[eventKey{GUID: "target-guid", Method: "onceEvent"}]
	conn.mu.Unlock()
	is.False(exists, "OnEventOnce must remove listener after execution")
}

// TestConnection_SendRequest_SimultaneousTimeoutAndDispatch (UNIT-WAIT-01 / UNIT-WAIT-02)
// Verifies that when context timeout/cancellation coincides with message dispatching,
// neither SendRequest nor Dispatch blocks or panics.
func TestConnection_SendRequest_SimultaneousTimeoutAndDispatch(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	const iterations = 50

	for i := range iterations {
		conn := NewConnection()
		reqID := 500 + i

		conn.SetTransportSend(func([]byte) error {
			go func() {
				// Rapid dispatch attempting to race with context cancel
				time.Sleep(time.Duration(i%5) * time.Millisecond)
				respMsg := []byte(fmt.Sprintf(`{"id":%d,"result":{}}`, reqID))
				conn.Dispatch(respMsg)
			}()
			return nil
		})

		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(i%5)*time.Millisecond)
		_, _ = conn.SendRequest(ctx, reqID, []byte(`{}`))
		cancel()

		// Verify callback cleanup
		conn.mu.Lock()
		_, exists := conn.callbacks[reqID]
		conn.mu.Unlock()
		is.False(exists, "callbacks map must not contain request after SendRequest finishes")
	}
}

// TestConnection_SendRequest_NonBlockingCallbackDrop (UNIT-WAIT-01 / UNIT-WAIT-02)
// Verifies that if a callback channel is full or receiver is gone, Dispatch drops the message
// gracefully without blocking the dispatch thread.
func TestConnection_SendRequest_NonBlockingCallbackDrop(t *testing.T) {
	t.Parallel()
	conn := NewConnection()

	// Manually register a callback channel with capacity 1 and fill it
	ch := make(chan Message, 1)
	ch <- Message{Method: "existing"}

	conn.mu.Lock()
	conn.callbacks[9999] = ch
	conn.mu.Unlock()

	done := make(chan struct{})
	go func() {
		// Dispatching another message for id 9999 must not block even though ch is full
		msg := []byte(`{"id":9999,"result":{}}`)
		conn.Dispatch(msg)
		close(done)
	}()

	select {
	case <-done:
		// Succeeded without blocking
	case <-time.After(1 * time.Second):
		t.Fatal("Dispatch blocked on full callback channel")
	}
}
