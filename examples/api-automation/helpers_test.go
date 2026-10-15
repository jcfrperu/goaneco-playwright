//go:build e2e

package apiautomation

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/examples/api-automation/models"
	"github.com/stretchr/testify/require"
)

// mustJSON serialises v to JSON, failing the test on error.
func mustJSON(t *testing.T, v any) []byte {
	t.Helper()
	b, err := json.Marshal(v)
	require.NoError(t, err, "json.Marshal failed")
	return b
}

// mustUnmarshalPet deserialises body into a Pet, failing the test on error.
func mustUnmarshalPet(t *testing.T, body []byte) *models.Pet {
	t.Helper()
	var pet models.Pet
	require.NoError(t, json.Unmarshal(body, &pet), "unmarshal Pet failed")
	return &pet
}

// mustUnmarshalPets deserialises body into a []Pet, failing the test on error.
func mustUnmarshalPets(t *testing.T, body []byte) []models.Pet {
	t.Helper()
	var pets []models.Pet
	require.NoError(t, json.Unmarshal(body, &pets), "unmarshal []Pet failed")
	return pets
}

// mustUnmarshalOrder deserialises body into an Order, failing the test on error.
func mustUnmarshalOrder(t *testing.T, body []byte) *models.Order {
	t.Helper()
	var order models.Order
	require.NoError(t, json.Unmarshal(body, &order), "unmarshal Order failed")
	return &order
}

// mustUnmarshalUser deserialises body into a User, failing the test on error.
func mustUnmarshalUser(t *testing.T, body []byte) *models.User {
	t.Helper()
	var user models.User
	require.NoError(t, json.Unmarshal(body, &user), "unmarshal User failed")
	return &user
}

// jsonHeaders returns the standard JSON content-type header map.
func jsonHeaders() map[string]string {
	return map[string]string{"Content-Type": "application/json"}
}

// createTestPet POSTs a new pet and registers a cleanup to delete it.
func createTestPet(t *testing.T) *models.Pet {
	t.Helper()
	apiCtx := newAPICtx(t)
	ctx := testCtx(t)

	pet := &models.Pet{
		Name:      "GoAneco-" + t.Name(),
		PhotoURLs: []string{"https://example.com/pet.jpg"},
		Status:    "available",
		Tags:      []models.Tag{{Name: "goaneco"}},
	}

	resp, err := apiCtx.Post(ctx, "/pet", &playwright.APIFetchOptions{
		Headers: jsonHeaders(),
		Data:    mustJSON(t, pet),
	})
	require.NoError(t, err, "createTestPet: POST /pet failed")
	defer resp.Dispose(ctx) //nolint:errcheck
	require.Equal(t, 200, resp.Status(), "createTestPet: unexpected status")

	body, err := resp.Body(ctx)
	require.NoError(t, err, "createTestPet: read body failed")
	created := mustUnmarshalPet(t, body)

	t.Cleanup(func() {
		cleanCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		cleanAPICtx := newCleanupAPICtx()
		if cleanAPICtx != nil {
			deletePetIfExists(t, cleanAPICtx, cleanCtx, created.ID)
			_ = cleanAPICtx.Dispose(cleanCtx)
		}
	})
	return created
}

// createTestOrder POSTs a new order and registers a cleanup to delete it.
func createTestOrder(t *testing.T, petID int64) *models.Order {
	t.Helper()
	apiCtx := newAPICtx(t)
	ctx := testCtx(t)

	order := &models.Order{
		PetID:    petID,
		Quantity: 1,
		Status:   "placed",
	}

	resp, err := apiCtx.Post(ctx, "/store/order", &playwright.APIFetchOptions{
		Headers: jsonHeaders(),
		Data:    mustJSON(t, order),
	})
	require.NoError(t, err, "createTestOrder: POST /store/order failed")
	defer resp.Dispose(ctx) //nolint:errcheck
	require.Equal(t, 200, resp.Status(), "createTestOrder: unexpected status")

	body, err := resp.Body(ctx)
	require.NoError(t, err, "createTestOrder: read body failed")
	created := mustUnmarshalOrder(t, body)

	t.Cleanup(func() {
		cleanCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		cleanAPICtx := newCleanupAPICtx()
		if cleanAPICtx != nil {
			deleteOrderIfExists(t, cleanAPICtx, cleanCtx, created.ID)
			_ = cleanAPICtx.Dispose(cleanCtx)
		}
	})
	return created
}

// createTestUser POSTs a new user and registers a cleanup to delete it.
// suffix is appended to the username to ensure uniqueness across parallel tests.
func createTestUser(t *testing.T, suffix string) *models.User {
	t.Helper()
	apiCtx := newAPICtx(t)
	ctx := testCtx(t)

	username := fmt.Sprintf("goaneco-%s", suffix)
	user := &models.User{
		Username:  username,
		FirstName: "Go",
		LastName:  "Aneco",
		Email:     username + "@example.com",
		Password:  "pass123",
		Phone:     "555-0100",
	}

	resp, err := apiCtx.Post(ctx, "/user", &playwright.APIFetchOptions{
		Headers: jsonHeaders(),
		Data:    mustJSON(t, user),
	})
	require.NoError(t, err, "createTestUser: POST /user failed")
	defer resp.Dispose(ctx) //nolint:errcheck
	require.Equal(t, 200, resp.Status(), "createTestUser: unexpected status")

	t.Cleanup(func() {
		cleanCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		cleanAPICtx := newCleanupAPICtx()
		if cleanAPICtx != nil {
			deleteUserIfExists(t, cleanAPICtx, cleanCtx, username)
			_ = cleanAPICtx.Dispose(cleanCtx)
		}
	})
	user.Username = username
	return user
}

// deletePetIfExists deletes a pet by ID, ignoring errors.
func deletePetIfExists(t *testing.T, apiCtx *playwright.APIRequestContext, ctx context.Context, id int64) {
	t.Helper()
	resp, err := apiCtx.Delete(ctx, fmt.Sprintf("/pet/%d", id))
	if err != nil {
		return
	}
	_ = resp.Dispose(ctx)
}

// deleteOrderIfExists deletes an order by ID, ignoring errors.
func deleteOrderIfExists(t *testing.T, apiCtx *playwright.APIRequestContext, ctx context.Context, id int64) {
	t.Helper()
	resp, err := apiCtx.Delete(ctx, fmt.Sprintf("/store/order/%d", id))
	if err != nil {
		return
	}
	_ = resp.Dispose(ctx)
}

// deleteUserIfExists deletes a user by username, ignoring errors.
func deleteUserIfExists(t *testing.T, apiCtx *playwright.APIRequestContext, ctx context.Context, username string) {
	t.Helper()
	resp, err := apiCtx.Delete(ctx, "/user/"+username)
	if err != nil {
		return
	}
	_ = resp.Dispose(ctx)
}

// newCleanupAPICtx creates a standalone APIRequestContext for cleanup purposes.
// Returns nil if Playwright is not available (e.g., already stopped).
func newCleanupAPICtx() *playwright.APIRequestContext {
	if globalPW == nil {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	base := baseURL
	apiCtx, err := globalPW.NewAPIRequestContext(ctx, &playwright.APIRequestContextOptions{
		BaseURL: &base,
	})
	if err != nil {
		return nil
	}
	return apiCtx
}
