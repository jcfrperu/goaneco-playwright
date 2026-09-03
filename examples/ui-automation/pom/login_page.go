//go:build e2e

package pom

import (
	"context"
	"fmt"

	playwright "github.com/jcfrperu/goaneco-playwright"
)

const (
	LoginURL         = "https://www.saucedemo.com/"
	StandardUsername = "standard_user"
	StandardPassword = "secret_sauce"
)

// LoginPage encapsulates all interactions with the SauceDemo login page.
type LoginPage struct {
	page *playwright.Page
}

// NewLoginPage returns a LoginPage bound to the given page.
func NewLoginPage(page *playwright.Page) *LoginPage {
	return &LoginPage{page: page}
}

// Navigate navigates to the SauceDemo login URL.
func (lp *LoginPage) Navigate(ctx context.Context) error {
	if err := lp.page.Goto(ctx, LoginURL); err != nil {
		return fmt.Errorf("login page navigate: %w", err)
	}
	return nil
}

// Login fills the credentials form and submits it.
func (lp *LoginPage) Login(ctx context.Context, username, password string) error {
	if err := lp.page.Locator("#user-name").Fill(ctx, username); err != nil {
		return fmt.Errorf("login fill username: %w", err)
	}
	if err := lp.page.Locator("#password").Fill(ctx, password); err != nil {
		return fmt.Errorf("login fill password: %w", err)
	}
	if err := lp.page.Locator("#login-button").Click(ctx); err != nil {
		return fmt.Errorf("login click submit: %w", err)
	}
	return nil
}

// IsVisible reports whether the login form is currently displayed.
func (lp *LoginPage) IsVisible(ctx context.Context) (bool, error) {
	return lp.page.Locator("#login-button").IsVisible(ctx)
}

// ErrorMessage returns the text of the login error banner, or empty string if none.
func (lp *LoginPage) ErrorMessage(ctx context.Context) (string, error) {
	loc := lp.page.Locator("[data-test='error']")
	visible, err := loc.IsVisible(ctx)
	if err != nil || !visible {
		return "", err
	}
	return loc.InnerText(ctx)
}
