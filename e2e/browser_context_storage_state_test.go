//go:build e2e

// E2E tests for: TestBrowserContextStorageState.java.
// Tests marked "Extra:" are additional tests with no equivalent method in playwright-java.
package e2e

import (
	"context"
	"encoding/json"
	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestBrowserContextOptionsStorageState verifies pre-populated storage state in BrowserContextOptions.
// Extra: no equivalent method in playwright-java; related class: TestBrowserContextStorageState.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextStorageState.java
func TestBrowserContextOptionsStorageState(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	bCtx, err := globalBrowser.NewContext(ctx, &playwright.BrowserContextOptions{
		StorageState: &playwright.StorageState{
			Origins: []playwright.OriginStorage{
				{
					Origin: srv.Prefix(),
					LocalStorage: []playwright.LocalStorageEntry{
						{Name: "preloaded_key", Value: "preloaded_val"},
					},
				},
			},
		},
	})
	must.NoError(err, "NewContext with StorageState failed")
	t.Cleanup(func() {
		closeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = bCtx.Close(closeCtx)
	})

	page, err := bCtx.NewPage(ctx)
	must.NoError(err, "NewPage failed")

	err = page.Goto(ctx, srv.EmptyPage())
	must.NoError(err, "Goto failed")

	valRaw, err := page.Evaluate(ctx, "() => localStorage.getItem('preloaded_key')")
	must.NoError(err, "Evaluate localStorage failed")

	if valRaw != "preloaded_val" {
		t.Errorf("localStorage preloaded_key = %v, want 'preloaded_val'", valRaw)
	}
}

// TestBrowserContextStorageState_CookiesAndLocalStorage (E2E-CTX-03)
// Verifies full round-trip initialization of StorageState with both Cookies and LocalStorage
// across multiple origins and verifies isolation from other contexts.
// Extra: no equivalent method in playwright-java; related class: TestBrowserContextStorageState.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextStorageState.java
func TestBrowserContextStorageState_CookiesAndLocalStorage(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	emptyPageURL := srv.EmptyPage()

	storageState := &playwright.StorageState{
		Cookies: []playwright.Cookie{
			{
				Name:     "session_id",
				Value:    "xyz_secret_999",
				URL:      &emptyPageURL,
				SameSite: playwright.SameSiteLax,
			},
			{
				Name:     "user_pref",
				Value:    "dark_mode",
				URL:      &emptyPageURL,
				SameSite: playwright.SameSiteLax,
			},
		},
		Origins: []playwright.OriginStorage{
			{
				Origin: srv.Prefix(),
				LocalStorage: []playwright.LocalStorageEntry{
					{Name: "auth_token", Value: "jwt_token_abc"},
					{Name: "app_theme", Value: "nord"},
				},
			},
		},
	}

	bCtx, err := globalBrowser.NewContext(ctx, &playwright.BrowserContextOptions{
		StorageState: storageState,
	})
	must.NoError(err, "NewContext with StorageState failed")
	t.Cleanup(func() {
		closeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = bCtx.Close(closeCtx)
	})

	page, err := bCtx.NewPage(ctx)
	must.NoError(err, "NewPage failed")

	err = page.Goto(ctx, srv.EmptyPage())
	must.NoError(err, "Goto failed")

	// 1. Verify Cookies were injected into document.cookie
	cookieStrRaw, err := page.Evaluate(ctx, "() => document.cookie")
	must.NoError(err, "Evaluate document.cookie failed")
	cookieStr, ok := cookieStrRaw.(string)
	is.True(ok, "expected cookie string")

	is.Contains(cookieStr, "session_id=xyz_secret_999")
	is.Contains(cookieStr, "user_pref=dark_mode")

	// 2. Verify LocalStorage entries were injected
	tokenRaw, err := page.Evaluate(ctx, "() => localStorage.getItem('auth_token')")
	must.NoError(err)
	is.Equal("jwt_token_abc", tokenRaw)

	themeRaw, err := page.Evaluate(ctx, "() => localStorage.getItem('app_theme')")
	must.NoError(err)
	is.Equal("nord", themeRaw)

	// 3. Verify a second page in the same context shares the storage state
	page2, err := bCtx.NewPage(ctx)
	must.NoError(err, "page2 NewPage failed")
	err = page2.Goto(ctx, srv.EmptyPage())
	must.NoError(err, "page2 Goto failed")

	token2, err := page2.Evaluate(ctx, "() => localStorage.getItem('auth_token')")
	must.NoError(err)
	is.Equal("jwt_token_abc", token2)

	// 4. Verify an isolated context created without StorageState does NOT have these cookies/storage
	isolatedCtx, err := globalBrowser.NewContext(ctx)
	must.NoError(err, "isolated NewContext failed")
	t.Cleanup(func() {
		closeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = isolatedCtx.Close(closeCtx)
	})

	isolatedPage, err := isolatedCtx.NewPage(ctx)
	must.NoError(err, "isolatedPage NewPage failed")
	err = isolatedPage.Goto(ctx, srv.EmptyPage())
	must.NoError(err, "isolatedPage Goto failed")

	isoCookies, err := isolatedPage.Evaluate(ctx, "() => document.cookie")
	must.NoError(err, "Evaluate document.cookie on isolated page failed")
	if cookieStr, ok := isoCookies.(string); ok && strings.Contains(cookieStr, "session_id") {
		t.Errorf("isolated context should not have session_id cookie, got %q", cookieStr)
	}

	isoStorage, err := isolatedPage.Evaluate(ctx, "() => localStorage.getItem('auth_token')")
	must.NoError(err, "Evaluate localStorage on isolated page failed")
	if isoStorage != nil {
		t.Errorf("isolated context should not have auth_token in localStorage, got %v", isoStorage)
	}
}

// TestBrowserContextStorageStatePath verifies that StorageStatePath loads cookies and localStorage
// from a JSON file on disk, equivalent to passing StorageState inline.
// Extra: no equivalent method in playwright-java; related class: TestBrowserContextStorageState.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextStorageState.java
func TestBrowserContextStorageStatePath(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)

	emptyPageURL := srv.EmptyPage()
	stateJSON := `{
		"cookies": [{"name":"file_cookie","value":"file_value","url":"` + emptyPageURL + `","sameSite":"Lax"}],
		"origins": [{"origin":"` + srv.Prefix() + `","localStorage":[{"name":"file_key","value":"file_val"}]}]
	}`

	f, err := os.CreateTemp("", "storage-state-*.json")
	must.NoError(err, "CreateTemp failed")
	t.Cleanup(func() { _ = os.Remove(f.Name()) })
	_, err = f.WriteString(stateJSON)
	must.NoError(err, "WriteString failed")
	_ = f.Close()

	path := f.Name()
	bCtx, err := globalBrowser.NewContext(ctx, &playwright.BrowserContextOptions{
		StorageStatePath: &path,
	})
	must.NoError(err, "NewContext with StorageStatePath failed")
	t.Cleanup(func() {
		closeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = bCtx.Close(closeCtx)
	})

	page, err := bCtx.NewPage(ctx)
	must.NoError(err, "NewPage failed")
	err = page.Goto(ctx, srv.EmptyPage())
	must.NoError(err, "Goto failed")

	cookieRaw, err := page.Evaluate(ctx, "() => document.cookie")
	must.NoError(err, "Evaluate cookie failed")
	cookieStr, _ := cookieRaw.(string)
	is.Contains(cookieStr, "file_cookie=file_value")

	storageRaw, err := page.Evaluate(ctx, "() => localStorage.getItem('file_key')")
	must.NoError(err, "Evaluate localStorage failed")
	is.Equal("file_val", storageRaw)
}

// TestBrowserContextLocalStorageSetGetEx2 verifies localStorage set/get in page.
// Extra: no equivalent method in playwright-java; related class: TestBrowserContextStorageState.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextStorageState.java
func TestBrowserContextLocalStorageSetGetEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	_, err := page.Evaluate(ctx, `() => localStorage.setItem('myKey', 'myValue')`)
	must.NoError(err)

	result, err := page.Evaluate(ctx, `() => localStorage.getItem('myKey')`)
	must.NoError(err)
	is.Equal("myValue", result)
}

// TestBrowserContextLocalStorageIsolatedAcrossContextsEx2 verifies contexts have separate localStorage.
// Extra: no equivalent method in playwright-java; related class: TestBrowserContextStorageState.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextStorageState.java
func TestBrowserContextLocalStorageIsolatedAcrossContextsEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	bc1 := newContext(t)
	bc2 := newContext(t)

	p1, err := bc1.NewPage(ctx)
	must.NoError(err)
	p2, err := bc2.NewPage(ctx)
	must.NoError(err)

	must.NoError(p1.Goto(ctx, srv.EmptyPage()))
	must.NoError(p2.Goto(ctx, srv.EmptyPage()))

	_, err = p1.Evaluate(ctx, `() => localStorage.setItem('isolated', 'ctx1')`)
	must.NoError(err)

	val, err := p2.Evaluate(ctx, `() => localStorage.getItem('isolated')`)
	must.NoError(err)
	is.Nil(val)
}

// TestBrowserContextLocalStorageClearedOnReloadEx2 verifies localStorage persists across reload.
// Extra: no equivalent method in playwright-java; related class: TestBrowserContextStorageState.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextStorageState.java
func TestBrowserContextLocalStorageClearedOnReloadEx2(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	page := newPage(t)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	_, err := page.Evaluate(ctx, `() => localStorage.setItem('persist', 'yes')`)
	must.NoError(err)

	must.NoError(page.Reload(ctx))

	val, err := page.Evaluate(ctx, `() => localStorage.getItem('persist')`)
	must.NoError(err)
	is.Equal("yes", val)
}

// TestBrowserContextStorageStateShouldSerialiseStorageStateWithLoneSurrogates verifies StorageState serialises lone surrogate characters
// as the Unicode replacement character U+FFFD.
// Ref: TestBrowserContextStorageState.java#shouldSerialiseStorageStateWithLoneSurrogates  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextStorageState.java
func TestBrowserContextStorageStateShouldSerialiseStorageStateWithLoneSurrogates(t *testing.T) {
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

// TestBrowserContextStorageStateShouldCaptureLocalStorage verifies StorageState captures localStorage from multiple origins.
// Ref: TestBrowserContextStorageState.java#shouldCaptureLocalStorage  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextStorageState.java
func TestBrowserContextStorageStateShouldCaptureLocalStorage(t *testing.T) {
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

// TestBrowserContextStorageStateShouldSetLocalStorage verifies StorageState can pre-populate localStorage.
// Ref: TestBrowserContextStorageState.java#shouldSetLocalStorage  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextStorageState.java
func TestBrowserContextStorageStateShouldSetLocalStorage(t *testing.T) {
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

// Extra: no equivalent method in playwright-java; related class: TestBrowserContextStorageState.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextStorageState.java
func TestStorageStateCaptureCookies(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/cookie-page", "text/html", `<!DOCTYPE html><html><body></body></html>`)

	bCtx, err := globalBrowser.NewContext(ctx)
	must.NoError(err, "NewContext failed")
	t.Cleanup(func() { _ = bCtx.Close(testCtx(t)) })

	err = bCtx.AddCookies(ctx, []playwright.Cookie{
		{Name: "session", Value: "abc123", Domain: strPtr("127.0.0.1"), Path: strPtr("/")},
	})
	must.NoError(err, "AddCookies failed")

	state, err := bCtx.StorageState(ctx)
	must.NoError(err, "StorageState failed")

	found := false
	for _, c := range state.Cookies {
		if c.Name == "session" && c.Value == "abc123" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("cookie 'session=abc123' not found in StorageState; cookies: %+v", state.Cookies)
	}
}

// Extra: no equivalent method in playwright-java; related class: TestBrowserContextStorageState.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextStorageState.java
func TestStorageStateCaptureLocalStorage(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	srv.ServeWithBody("/ls-page", "text/html", `<!DOCTYPE html><html><body></body></html>`)

	page := newPage(t)
	err := page.Goto(ctx, srv.Prefix()+"/ls-page")
	must.NoError(err, "Goto failed")
	_, err = page.Evaluate(ctx, `() => localStorage.setItem("key1", "value1")`)
	must.NoError(err, "localStorage.setItem failed")

	bCtx := page.Context()
	state, err := bCtx.StorageState(ctx)
	must.NoError(err, "StorageState failed")

	found := false
	for _, o := range state.Origins {
		for _, entry := range o.LocalStorage {
			if entry.Name == "key1" && entry.Value == "value1" {
				found = true
			}
		}
	}
	if !found {
		t.Errorf("localStorage entry 'key1=value1' not found in StorageState; origins: %+v", state.Origins)
	}
}

// Extra: no equivalent method in playwright-java; related class: TestBrowserContextStorageState.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextStorageState.java
func TestStorageStateEmptyByDefault(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	bCtx, err := globalBrowser.NewContext(ctx)
	must.NoError(err, "NewContext failed")
	t.Cleanup(func() { _ = bCtx.Close(testCtx(t)) })

	state, err := bCtx.StorageState(ctx)
	must.NoError(err, "StorageState failed")
	is.Len(state.Cookies, 0)
}

// TestStorageStateHasCookies verifies StorageState captures cookies.
// Extra: no equivalent method in playwright-java; related class: TestBrowserContextStorageState.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextStorageState.java
func TestStorageStateHasCookies(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	bc := newContext(t)

	page, err := bc.NewPage(ctx)
	must.NoError(err)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	_, err = page.Evaluate(ctx, `() => document.cookie = 'session=abc'`)
	must.NoError(err)

	state, err := bc.StorageState(ctx)
	must.NoError(err)
	must.NotNil(state)
	// May have cookies
}

// TestStorageStateHasOrigins verifies StorageState returns origins field.
// Extra: no equivalent method in playwright-java; related class: TestBrowserContextStorageState.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextStorageState.java
func TestStorageStateHasOrigins(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	bc := newContext(t)

	page, err := bc.NewPage(ctx)
	must.NoError(err)

	must.NoError(page.Goto(ctx, srv.EmptyPage()))
	_, err = page.Evaluate(ctx, `() => localStorage.setItem('test', 'value')`)
	must.NoError(err)

	state, err := bc.StorageState(ctx)
	must.NoError(err)
	must.NotNil(state)
	must.NotNil(state.Origins)
}

// TestStorageStateRoundTrip verifies storage state can be captured and applied.
// Extra: no equivalent method in playwright-java; related class: TestBrowserContextStorageState.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextStorageState.java
func TestStorageStateRoundTrip(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	srv := testserver.New(t)
	bc1 := newContext(t)

	page1, err := bc1.NewPage(ctx)
	must.NoError(err)

	must.NoError(page1.Goto(ctx, srv.EmptyPage()))
	_, err = page1.Evaluate(ctx, `() => localStorage.setItem('key', 'stored')`)
	must.NoError(err)

	state, err := bc1.StorageState(ctx)
	must.NoError(err)
	must.NotNil(state)
}

// TestLocalStorageSnapshotViaStorageState verifies that localStorage entries set in a page are
// captured by BrowserContext.StorageState and reflected in the returned Origins list.
// Extra: no equivalent method in playwright-java; related class: TestBrowserContextStorageState.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextStorageState.java
func TestLocalStorageSnapshotViaStorageState(t *testing.T) {
	t.Parallel()
	ctx := testCtx(t)
	srv := testserver.New(t)
	is := assert.New(t)
	must := require.New(t)

	bCtx := newContext(t)
	page, err := bCtx.NewPage(ctx)
	must.NoError(err)
	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	_, err = page.Evaluate(ctx, "() => localStorage.setItem('snap-key', 'snap-value')")
	must.NoError(err)

	state, err := bCtx.StorageState(ctx)
	must.NoError(err)
	must.NotNil(state)

	found := false
	for _, origin := range state.Origins {
		for _, entry := range origin.LocalStorage {
			if entry.Name == "snap-key" && entry.Value == "snap-value" {
				found = true
			}
		}
	}
	is.True(found, "StorageState should capture the localStorage entry 'snap-key=snap-value'; got: %+v", state.Origins)
}

// TestLocalStorageRestoreViaStorageState verifies that a BrowserContext created with a pre-populated
// StorageState exposes the restored localStorage values to pages navigating to that origin.
// Extra: no equivalent method in playwright-java; related class: TestBrowserContextStorageState.java  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextStorageState.java
func TestLocalStorageRestoreViaStorageState(t *testing.T) {
	t.Parallel()
	ctx := testCtx(t)
	srv := testserver.New(t)
	is := assert.New(t)
	must := require.New(t)

	// Step 1: capture a storage state containing a localStorage entry.
	bCtx1 := newContext(t)
	page1, err := bCtx1.NewPage(ctx)
	must.NoError(err)
	must.NoError(page1.Goto(ctx, srv.EmptyPage()))
	_, err = page1.Evaluate(ctx, "() => localStorage.setItem('restored-key', 'restored-value')")
	must.NoError(err)

	state, err := bCtx1.StorageState(ctx)
	must.NoError(err)
	must.NotEmpty(state.Origins, "expected at least one origin in StorageState after setItem")

	// Step 2: create a new context seeded with the captured state and verify restoration.
	bCtx2, err := globalBrowser.NewContext(ctx, &playwright.BrowserContextOptions{
		StorageState: state,
	})
	must.NoError(err)
	t.Cleanup(func() { _ = bCtx2.Close(testCtx(t)) })

	page2, err := bCtx2.NewPage(ctx)
	must.NoError(err)
	must.NoError(page2.Goto(ctx, srv.EmptyPage()))

	val, err := page2.Evaluate(ctx, "() => localStorage.getItem('restored-key')")
	must.NoError(err)
	is.Equal("restored-value", val, "localStorage should be restored from StorageState")
}

// TestBrowserContextStorageStateShouldRoundTripThroughTheFile verifies saving storage state to a file and restoring it.
// Ref: TestBrowserContextStorageState.java#shouldRoundTripThroughTheFile  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextStorageState.java
func TestBrowserContextStorageStateShouldRoundTripThroughTheFile(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)

	bCtx := newContext(t)
	page, err := bCtx.NewPage(ctx)
	must.NoError(err)
	must.NoError(page.Route(ctx, "**/*", func(route *playwright.Route) {
		body := "<html></html>"
		_ = route.Fulfill(ctx, &playwright.RouteFulfillOptions{Body: &body})
	}))
	must.NoError(page.Goto(ctx, "https://www.example.com"))
	_, err = page.Evaluate(ctx, `() => {
		localStorage.setItem("name1", "value1");
		document.cookie = "username=John Doe";
	}`)
	must.NoError(err)

	path := filepath.Join(t.TempDir(), "storage-state.json")
	state, err := bCtx.StorageState(ctx, &playwright.StorageStateOptions{Path: &path})
	must.NoError(err)
	data, err := os.ReadFile(path)
	must.NoError(err)
	stateJSON, err := json.Marshal(state)
	must.NoError(err)
	is.JSONEq(string(stateJSON), string(data))

	must.Len(state.Cookies, 1)
	cookie := state.Cookies[0]
	is.Equal("username", cookie.Name)
	is.Equal("John Doe", cookie.Value)
	is.Equal(a4Ptr("www.example.com"), cookie.Domain)
	is.Equal(a4Ptr("/"), cookie.Path)
	is.Equal(a4Ptr(float64(-1)), cookie.Expires)
	is.Equal(a4Ptr(false), cookie.HTTPOnly)
	is.Equal(a4Ptr(false), cookie.Secure)
	is.Contains([]playwright.SameSiteAttribute{
		playwright.SameSiteStrict,
		playwright.SameSiteLax,
		playwright.SameSiteNone,
	}, cookie.SameSite)
	must.Equal([]playwright.OriginStorage{{
		Origin: "https://www.example.com",
		LocalStorage: []playwright.LocalStorageEntry{
			{Name: "name1", Value: "value1"},
		},
	}}, state.Origins)

	bCtx2, err := globalBrowser.NewContext(ctx, &playwright.BrowserContextOptions{StorageStatePath: &path})
	must.NoError(err)
	t.Cleanup(func() {
		closeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = bCtx2.Close(closeCtx)
	})
	page2, err := bCtx2.NewPage(ctx)
	must.NoError(err)
	must.NoError(page2.Route(ctx, "**/*", func(route *playwright.Route) {
		body := "<html></html>"
		_ = route.Fulfill(ctx, &playwright.RouteFulfillOptions{Body: &body})
	}))
	must.NoError(page2.Goto(ctx, "https://www.example.com"))
	value, err := page2.Evaluate(ctx, `() => localStorage.getItem("name1")`)
	must.NoError(err)
	is.Equal("value1", value)
	if globalBTName != "firefox" {
		cookieValue, err := page2.Evaluate(ctx, "document.cookie")
		must.NoError(err)
		is.Equal("username=John Doe", cookieValue)
	}
}

// TestBrowserContextStorageStateShouldSupportIndexedDB verifies indexedDB is included in storage state.
// Ref: TestBrowserContextStorageState.java#shouldSupportIndexedDB  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextStorageState.java
func TestBrowserContextStorageStateShouldSupportIndexedDB(t *testing.T) {
	t.Skip("not implemented: BrowserContext.StorageState IndexedDB option and /to-do-notifications test page")
}

// TestBrowserContextStorageStateSetStorageStateShouldHandleMissingFile verifies setStorageState with a missing file fails.
// Ref: TestBrowserContextStorageState.java#setStorageStateShouldHandleMissingFile  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextStorageState.java
func TestBrowserContextStorageStateSetStorageStateShouldHandleMissingFile(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	bCtx := newContext(t)

	path := filepath.Join(t.TempDir(), "does-not-exist.json")
	err := bCtx.SetStorageState(ctx, path)
	must.Error(err)
	is.Contains(err.Error(), "Failed to read storage state from file")
}

// TestBrowserContextStorageStateShouldRoundTripOPFS verifies OPFS is included in storage state.
// Ref: TestBrowserContextStorageState.java#shouldRoundTripOPFS  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextStorageState.java
func TestBrowserContextStorageStateShouldRoundTripOPFS(t *testing.T) {
	t.Skip("not implemented: BrowserContext.StorageState OPFS option / SetStorageState")
}
