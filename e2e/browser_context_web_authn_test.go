//go:build e2e

// E2E tests for: TestBrowserContextWebAuthn.java.
package e2e

import (
	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

// TestBrowserContextWebAuthnShouldNotInterceptNavigatorCredentialsWithoutInstall verifies credentials are not intercepted without install.
// Ref: TestBrowserContextWebAuthn.java#shouldNotInterceptNavigatorCredentialsWithoutInstall  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextWebAuthn.java
func TestBrowserContextWebAuthnShouldNotInterceptNavigatorCredentialsWithoutInstall(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	bCtx := newContextWithCleanup(t)

	creds := bCtx.Credentials()
	if creds == nil {
		t.Skip("Credentials not available in this browser context (server did not include credentials GUID)")
		return
	}

	// Without Install(), the virtual authenticator is not active.
	// get() should return empty list since no credentials are registered.
	list, err := creds.Get(ctx)
	must.NoError(err, "Credentials.Get() failed before Install()")
	must.Empty(list, "expected no credentials before Install()")
}

// TestBrowserContextWebAuthnShouldSeedKnownCredentialAndAuthenticate verifies seeding a credential and authenticating.
// Ref: TestBrowserContextWebAuthn.java#shouldSeedKnownCredentialAndAuthenticate  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextWebAuthn.java
func TestBrowserContextWebAuthnShouldSeedKnownCredentialAndAuthenticate(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	bCtx := newContextWithCleanup(t)

	creds := bCtx.Credentials()
	if creds == nil {
		t.Skip("Credentials not available in this browser context (server did not include credentials GUID)")
		return
	}

	must.NoError(creds.Install(ctx), "Credentials.Install() failed")

	// Create a credential for a relying party.
	cred, err := creds.Create(ctx, "example.com")
	must.NoError(err, "Credentials.Create() failed")
	must.NotNil(cred, "Credentials.Create() returned nil")
	must.NotEmpty(cred.ID, "credential ID should not be empty")
	must.Equal("example.com", cred.RpID, "credential RpID should match")

	// Verify it appears in Get().
	list, err := creds.Get(ctx)
	must.NoError(err, "Credentials.Get() failed")
	must.Len(list, 1, "expected exactly one credential")
	must.Equal(cred.ID, list[0].ID, "listed credential ID should match created ID")
}

// TestBrowserContextWebAuthnShouldCapturePageCreatedCredentialAndReuseItInAnotherContext verifies capturing a credential and reusing it in another context.
// Ref: TestBrowserContextWebAuthn.java#shouldCapturePageCreatedCredentialAndReuseItInAnotherContext  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextWebAuthn.java
func TestBrowserContextWebAuthnShouldCapturePageCreatedCredentialAndReuseItInAnotherContext(t *testing.T) {
	t.Parallel()
	must := require.New(t)

	ctx := testCtx(t)
	bCtx1 := newContextWithCleanup(t)
	bCtx2 := newContextWithCleanup(t)

	creds1 := bCtx1.Credentials()
	creds2 := bCtx2.Credentials()
	if creds1 == nil || creds2 == nil {
		t.Skip("Credentials not available in this browser context")
		return
	}

	must.NoError(creds1.Install(ctx))
	must.NoError(creds2.Install(ctx))

	// Create a credential in context 1.
	cred, err := creds1.Create(ctx, "example.com")
	must.NoError(err, "Credentials.Create() in context 1 failed")
	must.NotNil(cred)

	// Delete it from context 1 and re-create in context 2 to simulate "reuse".
	must.NoError(creds1.Delete(ctx, cred.ID), "Credentials.Delete() failed")

	_, err = creds2.Create(ctx, "example.com")
	must.NoError(err, "Credentials.Create() in context 2 failed")

	list, err := creds2.Get(ctx)
	must.NoError(err)
	must.NotEmpty(list, "context 2 should have at least one credential")
}

// Ref: TestBrowserContextWebAuthn.java#shouldSeedAndReportSignCount  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextWebAuthn.java
func TestBrowserContextWebAuthnShouldSeedAndReportSignCount(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	srv := testserver.New(t)
	bCtx := newContextWithCleanup(t)

	creds := bCtx.Credentials()
	if creds == nil {
		t.Skip("Credentials not available in this browser context")
		return
	}

	fresh, err := creds.Create(ctx, "fresh.example.com")
	must.NoError(err)
	is.Equal(0, fresh.SignCount)

	signCount41 := 41
	seeded, err := creds.Create(ctx, "localhost", &playwright.CredentialsCreateOptions{SignCount: &signCount41})
	must.NoError(err)
	is.Equal(41, seeded.SignCount)

	must.NoError(creds.Install(ctx))
	page, err := bCtx.NewPage(ctx)
	must.NoError(err)
	must.NoError(page.Goto(ctx, srv.EmptyPage()))

	assertAndGetSignCount := func(pg *playwright.Page, rpId string) int {
		raw, err := pg.Evaluate(ctx, `async ({ rpId }) => {
			const challenge = crypto.getRandomValues(new Uint8Array(32));
			const cred = await navigator.credentials.get({
				publicKey: { challenge, rpId, userVerification: 'preferred' },
			});
			return new DataView(cred.response.authenticatorData).getUint32(33);
		}`, map[string]any{"rpId": rpId})
		must.NoError(err)
		num, ok := raw.(float64)
		must.True(ok)
		return int(num)
	}

	is.Equal(42, assertAndGetSignCount(page, "localhost"))
	is.Equal(43, assertAndGetSignCount(page, "localhost"))

	credentials, err := creds.Get(ctx, &playwright.CredentialsGetOptions{ID: &seeded.ID})
	must.NoError(err)
	must.Len(credentials, 1)
	captured := credentials[0]
	is.Equal(43, captured.SignCount)

	freshList, err := creds.Get(ctx, &playwright.CredentialsGetOptions{ID: &fresh.ID})
	must.NoError(err)
	must.Len(freshList, 1)
	is.Equal(0, freshList[0].SignCount)

	// A captured credential continues from the same counter in another context.
	bCtx2 := newContextWithCleanup(t)
	creds2 := bCtx2.Credentials()
	if creds2 != nil {
		_, err = creds2.Create(ctx, captured.RpID, &playwright.CredentialsCreateOptions{
			ID:         &captured.ID,
			UserHandle: &captured.UserHandle,
			PrivateKey: &captured.PrivateKey,
			PublicKey:  &captured.PublicKey,
			SignCount:  &captured.SignCount,
		})
		must.NoError(err)
		must.NoError(creds2.Install(ctx))
		page2, err := bCtx2.NewPage(ctx)
		must.NoError(err)
		must.NoError(page2.Goto(ctx, srv.EmptyPage()))
		is.Equal(44, assertAndGetSignCount(page2, "localhost"))
	}
}

// Ref: TestBrowserContextWebAuthn.java#shouldRejectInvalidSignCount  Folder: playwright-reference-projects/playwright-java/playwright/src/test/java/com/microsoft/playwright/TestBrowserContextWebAuthn.java
func TestBrowserContextWebAuthnShouldRejectInvalidSignCount(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)
	ctx := testCtx(t)
	bCtx := newContextWithCleanup(t)

	creds := bCtx.Credentials()
	if creds == nil {
		t.Skip("Credentials not available in this browser context")
		return
	}

	neg1 := -1
	_, err := creds.Create(ctx, "example.com", &playwright.CredentialsCreateOptions{SignCount: &neg1})
	must.Error(err)
	is.Contains(err.Error(), "signCount must be between 0 and 4294967295, got -1")

	list, err := creds.Get(ctx)
	must.NoError(err)
	is.Empty(list)
}
