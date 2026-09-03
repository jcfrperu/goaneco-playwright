//go:build e2e

// Ref: TestBrowserContextStorageState.java
package e2e

import (
	"context"
	"testing"
	"time"

	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestStorageStateLoneSurrogates verifies StorageState serialises lone surrogate characters
// as the Unicode replacement character U+FFFD.
// Ref: TestBrowserContextStorageState.java#shouldSerialiseStorageStateWithLoneSurrogates
func TestStorageStateLoneSurrogates(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	bc := newContext(t)

	pg, err := bc.NewPage(ctx)
	must.NoError(err)

	must.NoError(pg.Goto(ctx, srv.EmptyPage()))

	// Store a lone surrogate (code point 55934 = 0xDA7E) in localStorage.
	_, err = pg.Evaluate(ctx, `() => window.localStorage.setItem('foo', String.fromCharCode(55934))`)
	must.NoError(err)

	state, err := bc.StorageState(ctx)
	must.NoError(err)
	must.NotNil(state)

	is.Len(state.Origins, 1)
	is.Len(state.Origins[0].LocalStorage, 1)
	is.Equal("foo", state.Origins[0].LocalStorage[0].Name)

	// Lone surrogates cannot be represented in UTF-8; the browser replaces them
	// with the Unicode replacement character U+FFFD (decimal 65533).
	expected := string(rune(65533))
	is.Equal(expected, state.Origins[0].LocalStorage[0].Value)
}

// TestStorageStateCaptureMultipleOrigins verifies StorageState captures localStorage from multiple origins.
// Ref: TestBrowserContextStorageState.java#shouldCaptureLocalStorage
func TestStorageStateCaptureMultipleOrigins(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)

	bc := newContext(t)
	pg, err := bc.NewPage(ctx)
	must.NoError(err)

	// Route all navigation to return a simple HTML page so we can visit any origin.
	must.NoError(pg.Route(ctx, "**/*", func(r *playwright.Route) {
		body := "<html></html>"
		_ = r.Fulfill(ctx, &playwright.RouteFulfillOptions{Body: &body})
	}))

	must.NoError(pg.Goto(ctx, "https://www.example.com"))
	_, err = pg.Evaluate(ctx, `() => localStorage['name1'] = 'value1'`)
	must.NoError(err)

	must.NoError(pg.Goto(ctx, "https://www.domain.com"))
	_, err = pg.Evaluate(ctx, `() => localStorage['name2'] = 'value2'`)
	must.NoError(err)

	state, err := bc.StorageState(ctx)
	must.NoError(err)
	must.NotNil(state)

	// Find the two origins in storage state (order may vary).
	originMap := make(map[string][]playwright.LocalStorageEntry)
	for _, o := range state.Origins {
		originMap[o.Origin] = o.LocalStorage
	}

	is.Contains(originMap, "https://www.example.com")
	is.Contains(originMap, "https://www.domain.com")

	findEntry := func(entries []playwright.LocalStorageEntry, name string) string {
		for _, e := range entries {
			if e.Name == name {
				return e.Value
			}
		}
		return ""
	}

	is.Equal("value1", findEntry(originMap["https://www.example.com"], "name1"))
	is.Equal("value2", findEntry(originMap["https://www.domain.com"], "name2"))
}

// TestStorageStateSetLocalStorage verifies StorageState can pre-populate localStorage.
// Ref: TestBrowserContextStorageState.java#shouldSetLocalStorage
func TestStorageStateSetLocalStorage(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)

	storage := &playwright.StorageState{
		Origins: []playwright.OriginStorage{
			{
				Origin: "https://www.example.com",
				LocalStorage: []playwright.LocalStorageEntry{
					{Name: "name1", Value: "value1"},
				},
			},
		},
	}

	bCtx, err := globalBrowser.NewContext(ctx, &playwright.BrowserContextOptions{StorageState: storage})
	must.NoError(err)
	t.Cleanup(func() {
		c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = bCtx.Close(c)
	})

	pg, err := bCtx.NewPage(ctx)
	must.NoError(err)

	must.NoError(pg.Route(ctx, "**/*", func(r *playwright.Route) {
		body := "<html></html>"
		_ = r.Fulfill(ctx, &playwright.RouteFulfillOptions{Body: &body})
	}))

	must.NoError(pg.Goto(ctx, "https://www.example.com"))

	val, err := pg.Evaluate(ctx, "window.localStorage")
	must.NoError(err)

	m, ok := val.(map[string]any)
	is.True(ok, "localStorage should be a map, got %T", val)
	is.Equal("value1", m["name1"])
}
