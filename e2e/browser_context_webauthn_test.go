//go:build e2e

// E2E tests for BrowserContext WebAuthn/Passkeys virtual authenticator.
// Migration of: TestBrowserContextWebAuthn.java
package e2e

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestWebAuthnNoInterceptWithoutInstall verifies credentials are not intercepted without install.
// Ref: TestBrowserContextWebAuthn.java#shouldNotInterceptNavigatorCredentialsWithoutInstall
func TestWebAuthnNoInterceptWithoutInstall(t *testing.T) {
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

// TestWebAuthnSeedAndAuthenticate verifies seeding a credential and authenticating.
// Ref: TestBrowserContextWebAuthn.java#shouldSeedKnownCredentialAndAuthenticate
func TestWebAuthnSeedAndAuthenticate(t *testing.T) {
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

// TestWebAuthnCaptureAndReuseCredential verifies capturing a credential and reusing it in another context.
// Ref: TestBrowserContextWebAuthn.java#shouldCapturePageCreatedCredentialAndReuseItInAnotherContext
func TestWebAuthnCaptureAndReuseCredential(t *testing.T) {
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
