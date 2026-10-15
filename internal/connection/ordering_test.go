package connection

import (
	"encoding/json"
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Events must reach a handler in the order they were dispatched.
func TestConnection_EventsDeliveredInOrder(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	conn := NewConnection()
	const n = 500
	var mu sync.Mutex
	got := make([]int, 0, n)
	conn.OnEvent("g", "console", func(p json.RawMessage) {
		var v struct {
			N int `json:"n"`
		}
		_ = json.Unmarshal(p, &v)
		mu.Lock()
		got = append(got, v.N)
		mu.Unlock()
	})

	for i := 0; i < n; i++ {
		conn.Dispatch([]byte(fmt.Sprintf(`{"guid":"g","method":"console","params":{"n":%d}}`, i)))
	}
	conn.WaitDispatch()

	mu.Lock()
	defer mu.Unlock()
	is.Len(got, n)
	for i, v := range got {
		if v != i {
			t.Fatalf("event %d delivered at position %d", v, i)
		}
	}
}

// A panicking handler must not crash the process nor stop later events from being delivered.
func TestConnection_HandlerPanicIsRecovered(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	conn := NewConnection()
	var mu sync.Mutex
	calls := 0
	conn.OnEvent("g", "ev", func(json.RawMessage) {
		mu.Lock()
		calls++
		mu.Unlock()
		panic("boom")
	})

	must.NotPanics(func() {
		conn.Dispatch([]byte(`{"guid":"g","method":"ev","params":{}}`))
		conn.Dispatch([]byte(`{"guid":"g","method":"ev","params":{}}`))
		conn.WaitDispatch()
	})
	mu.Lock()
	defer mu.Unlock()
	must.Equal(2, calls)
}
