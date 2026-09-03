package playwright

import (
	"context"
	"encoding/json"
	"fmt"
)

// VirtualCredential represents a virtual WebAuthn credential.
type VirtualCredential struct {
	ID         string `json:"id"`
	RpID       string `json:"rpId"`
	UserHandle string `json:"userHandle"`
	PrivateKey string `json:"privateKey"`
	PublicKey  string `json:"publicKey"`
}

// CredentialsCreateOptions configures Credentials.Create.
type CredentialsCreateOptions struct {
	// ID is the credential identifier (base64url-encoded). If empty, the server generates one.
	ID *string
	// UserHandle is the user handle (base64url-encoded).
	UserHandle *string
	// PrivateKey is the private key in PKCS#8 PEM format.
	PrivateKey *string
	// PublicKey is the public key in SPKI PEM format.
	PublicKey *string
}

// CredentialsGetOptions filters the credentials returned by Credentials.Get.
type CredentialsGetOptions struct {
	// RpID filters credentials to those registered for this relying party ID.
	RpID *string
	// ID filters to a single credential by its ID.
	ID *string
}

// Credentials provides a virtual WebAuthn authenticator scoped to a BrowserContext.
// Obtain via BrowserContext.Credentials().
type Credentials struct {
	owner ChannelOwner
}

// Install activates the virtual authenticator for this browser context.
// After calling Install, WebAuthn operations use the virtual authenticator
// instead of the platform authenticator.
func (cr *Credentials) Install(ctx context.Context) error {
	_, err := cr.owner.SendMessageRequest(ctx, "credentialsInstall", struct{}{})
	if err != nil {
		return fmt.Errorf("credentials.install failed: %w", err)
	}
	return nil
}

// Create registers a new virtual credential for the given relying party ID.
func (cr *Credentials) Create(ctx context.Context, rpID string, opts ...*CredentialsCreateOptions) (*VirtualCredential, error) {
	params := map[string]any{
		"rpId": rpID,
	}
	if len(opts) > 0 && opts[0] != nil {
		o := opts[0]
		if o.ID != nil {
			params["id"] = *o.ID
		}
		if o.UserHandle != nil {
			params["userHandle"] = *o.UserHandle
		}
		if o.PrivateKey != nil {
			params["privateKey"] = *o.PrivateKey
		}
		if o.PublicKey != nil {
			params["publicKey"] = *o.PublicKey
		}
	}
	result, err := cr.owner.SendMessageRequest(ctx, "credentialsCreate", params)
	if err != nil {
		return nil, fmt.Errorf("credentials.create failed: %w", err)
	}
	var resp struct {
		Credential VirtualCredential `json:"credential"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse credentialsCreate response: %w", err)
	}
	return &resp.Credential, nil
}

// Get returns virtual credentials matching the given filter options.
func (cr *Credentials) Get(ctx context.Context, opts ...*CredentialsGetOptions) ([]*VirtualCredential, error) {
	params := map[string]any{}
	if len(opts) > 0 && opts[0] != nil {
		o := opts[0]
		if o.RpID != nil {
			params["rpId"] = *o.RpID
		}
		if o.ID != nil {
			params["id"] = *o.ID
		}
	}
	result, err := cr.owner.SendMessageRequest(ctx, "credentialsGet", params)
	if err != nil {
		return nil, fmt.Errorf("credentials.get failed: %w", err)
	}
	var resp struct {
		Credentials []*VirtualCredential `json:"credentials"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse credentialsGet response: %w", err)
	}
	return resp.Credentials, nil
}

// Delete removes the virtual credential with the given ID.
func (cr *Credentials) Delete(ctx context.Context, id string) error {
	_, err := cr.owner.SendMessageRequest(ctx, "credentialsDelete", map[string]any{
		"id": id,
	})
	if err != nil {
		return fmt.Errorf("credentials.delete failed: %w", err)
	}
	return nil
}

// Credentials returns the Credentials object for this BrowserContext.
// The Credentials GUID comes from the BrowserContext initializer.
func (c *BrowserContext) Credentials() *Credentials {
	var init struct {
		Credentials struct {
			Guid string `json:"guid"`
		} `json:"credentials"`
	}
	if raw := c.initializer; len(raw) > 0 {
		_ = json.Unmarshal(raw, &init)
	}
	if init.Credentials.Guid == "" {
		return nil
	}
	return &Credentials{owner: c.owner.child(init.Credentials.Guid)}
}
