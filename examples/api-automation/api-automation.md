# API Automation Examples

End-to-end API test automation for the public [Petstore API](https://petstore.swagger.io/v2) using the
[goaneco-playwright](https://github.com/jcfrperu/goaneco-playwright) Go client.  
All examples use Playwright's `APIRequestContext` for HTTP requests and are tagged with `//go:build e2e`.

---

## Prerequisites

| Requirement | Details |
|---|---|
| Go | 1.21+ |
| Node.js | 18+ |
| Playwright | `npm install -g playwright` or `npm install -g playwright-core` |
| Internet access | Tests reach `https://petstore.swagger.io/v2` |

Set the path to the Playwright CLI before running:

```bash
# Linux / macOS
export PLAYWRIGHT_CLI_PATH=$(npm root -g)/playwright/node_modules/playwright-core/cli.js

# Windows PowerShell
$env:PLAYWRIGHT_CLI_PATH = "$env:APPDATA\npm\node_modules\playwright\node_modules\playwright-core\cli.js"
```

---

## Running the examples

```bash
# Run all 30 scenarios
go test -tags e2e -v -timeout 300s ./examples/api-automation/...

# Run a single scenario by test name
go test -tags e2e -v -timeout 60s -run TestScenario01AddPet ./examples/api-automation/...
```

No browser is launched. `TestMain` starts a single `playwright.Playwright` instance, creates an
`APIRequestContext` pointed at `https://petstore.swagger.io/v2`, and tears it down when the suite finishes.

---

## Project structure

```
examples/api-automation/
├── doc.go                                      # package apiautomation
├── main_test.go                                # TestMain: Playwright lifecycle, newAPICtx, testCtx
├── helpers_test.go                             # mustJSON, mustUnmarshal*, createTest*, deleteIfExists*
│
├── scenario_01_add_pet_test.go
├── scenario_02_get_pet_by_id_test.go
├── scenario_03_update_pet_test.go
├── scenario_04_delete_pet_test.go
├── scenario_05_find_by_status_available_test.go  # table-driven: covers statuses 05/06/07
├── scenario_08_find_by_status_multiple_test.go
├── scenario_09_update_pet_with_form_test.go
├── scenario_10_find_by_tags_test.go
├── scenario_11_get_pet_not_found_test.go
├── scenario_12_pet_full_lifecycle_test.go
├── scenario_13_get_inventory_test.go
├── scenario_14_place_order_test.go
├── scenario_15_get_order_by_id_test.go
├── scenario_16_delete_order_test.go
├── scenario_17_place_order_complete_test.go
├── scenario_18_get_order_not_found_test.go
├── scenario_19_place_order_with_quantity_test.go
├── scenario_20_order_full_lifecycle_test.go
├── scenario_21_create_user_test.go
├── scenario_22_get_user_by_username_test.go
├── scenario_23_update_user_test.go
├── scenario_24_delete_user_test.go
├── scenario_25_login_valid_test.go
├── scenario_26_logout_test.go
├── scenario_27_create_users_with_array_test.go
├── scenario_28_create_users_with_list_test.go
├── scenario_29_login_invalid_test.go
├── scenario_30_user_full_lifecycle_test.go
│
├── features/
│   ├── 01_add_pet.feature
│   ├── 02_get_pet_by_id.feature
│   ├── ...  (30 feature files, one per scenario)
│   └── 30_user_full_lifecycle.feature
│
├── models/
│   ├── doc.go
│   ├── pet.go      # Pet, Category, Tag
│   ├── order.go    # Order
│   └── user.go     # User, ApiResponse
│
└── api-automation.md
```

---

## Shared infrastructure

### `main_test.go`

`TestMain` starts a single `*playwright.Playwright` for the whole package (no browser launched).
Two helpers are available in every test:

| Helper | Returns | Description |
|---|---|---|
| `newAPICtx(t)` | `*playwright.APIRequestContext` | New context with `BaseURL = "https://petstore.swagger.io/v2"`; auto-disposed at `t.Cleanup` |
| `testCtx(t)` | `context.Context` | 30-second timeout; cancelled at `t.Cleanup` |

### `helpers_test.go`

| Helper | Description |
|---|---|
| `mustJSON(t, v)` | Marshals `v` to JSON; fails the test on error |
| `mustUnmarshalPet(t, body)` | Deserialises response bytes into `*models.Pet` |
| `mustUnmarshalPets(t, body)` | Deserialises into `[]models.Pet` |
| `mustUnmarshalOrder(t, body)` | Deserialises into `*models.Order` |
| `mustUnmarshalUser(t, body)` | Deserialises into `*models.User` |
| `createTestPet(t)` | Creates a pet via POST /pet; registers DELETE cleanup |
| `createTestOrder(t, petID)` | Creates an order via POST /store/order; registers DELETE cleanup |
| `createTestUser(t, suffix)` | Creates a user via POST /user with a unique username; registers DELETE cleanup |
| `deletePetIfExists(t, apiCtx, id)` | Best-effort DELETE /pet/{id}; errors ignored |
| `deleteOrderIfExists(t, apiCtx, id)` | Best-effort DELETE /store/order/{id} |
| `deleteUserIfExists(t, apiCtx, username)` | Best-effort DELETE /user/{username} |

### `models/`

Data models mirroring the Petstore API schema:

```go
// Pet
type Pet struct {
    ID        int64     `json:"id,omitempty"`
    Category  *Category `json:"category,omitempty"`
    Name      string    `json:"name"`
    PhotoURLs []string  `json:"photoUrls"`
    Tags      []Tag     `json:"tags,omitempty"`
    Status    string    `json:"status,omitempty"` // available | pending | sold
}

// Order
type Order struct {
    ID       int64  `json:"id,omitempty"`
    PetID    int64  `json:"petId"`
    Quantity int    `json:"quantity"`
    ShipDate string `json:"shipDate,omitempty"`
    Status   string `json:"status,omitempty"` // placed | approved | delivered
    Complete bool   `json:"complete,omitempty"`
}

// User
type User struct {
    ID         int64  `json:"id,omitempty"`
    Username   string `json:"username"`
    FirstName  string `json:"firstName,omitempty"`
    LastName   string `json:"lastName,omitempty"`
    Email      string `json:"email,omitempty"`
    Password   string `json:"password,omitempty"`
    Phone      string `json:"phone,omitempty"`
    UserStatus int    `json:"userStatus,omitempty"`
}
```

---

## Scenarios

### Pet endpoints (scenarios 01–12)

---

### 01 — Add a new pet
**Test:** `TestScenario01AddPet`  
**File:** `scenario_01_add_pet_test.go`  
**Feature:** `features/01_add_pet.feature`

Posts a new pet named `"Buddy"` with status `"available"` to `POST /pet`.  
Verifies: status 200, `id > 0`, `name == "Buddy"`, `status == "available"`.

---

### 02 — Get a pet by ID
**Test:** `TestScenario02GetPetByID`  
**File:** `scenario_02_get_pet_by_id_test.go`  
**Feature:** `features/02_get_pet_by_id.feature`

Creates a pet via `createTestPet`, then fetches it with `GET /pet/{id}`.  
Verifies: status 200, returned `id`, `name`, and `status` match the created pet.

---

### 03 — Update a pet
**Test:** `TestScenario03UpdatePet`  
**File:** `scenario_03_update_pet_test.go`  
**Feature:** `features/03_update_pet.feature`

Creates a pet, then sends `PUT /pet` with `status = "sold"`.  
Verifies: status 200, `id` unchanged, `status == "sold"`.

---

### 04 — Delete a pet
**Test:** `TestScenario04DeletePet`  
**File:** `scenario_04_delete_pet_test.go`  
**Feature:** `features/04_delete_pet.feature`

Creates a pet, deletes it with `DELETE /pet/{id}`, then attempts `GET /pet/{id}`.  
Verifies: DELETE returns 200; subsequent GET returns 404.

---

### 05, 06, 07 — Find pets by status
**Test:** `TestScenario05FindByStatus` (table-driven with subtests `available`, `pending`, `sold`)  
**File:** `scenario_05_find_by_status_available_test.go`  
**Features:** `features/05_find_by_status_available.feature`, `features/06_find_by_status_pending.feature`, `features/07_find_by_status_sold.feature`

Queries `GET /pet/findByStatus?status={status}` for each of the three valid statuses in parallel subtests.  
Verifies: status 200, non-empty array, every pet in the response carries the queried status.

---

### 08 — Find pets by multiple statuses
**Test:** `TestScenario08FindByStatusMultiple`  
**File:** `scenario_08_find_by_status_multiple_test.go`  
**Feature:** `features/08_find_by_status_multiple.feature`

Queries `GET /pet/findByStatus?status=available&status=pending` in a single request.  
Verifies: status 200, non-empty array, every pet's status is either `"available"` or `"pending"`.

---

### 09 — Update pet with form data
**Test:** `TestScenario09UpdatePetWithForm`  
**File:** `scenario_09_update_pet_with_form_test.go`  
**Feature:** `features/09_update_pet_with_form.feature`

Creates a pet, then sends `POST /pet/{id}` with `Content-Type: application/x-www-form-urlencoded`
and form fields `name=UpdatedName`, `status=pending`.  
Verifies: status 200.

---

### 10 — Find pets by tag
**Test:** `TestScenario10FindByTags`  
**File:** `scenario_10_find_by_tags_test.go`  
**Feature:** `features/10_find_by_tags.feature`

Queries `GET /pet/findByTags?tags=goaneco` (the tag applied by `createTestPet`).  
Verifies: status 200, response body is a parseable JSON array.

---

### 11 — Get non-existent pet
**Test:** `TestScenario11GetPetNotFound`  
**File:** `scenario_11_get_pet_not_found_test.go`  
**Feature:** `features/11_get_pet_not_found.feature`

Requests `GET /pet/999999999` — an ID that cannot exist.  
Verifies: status 404.

---

### 12 — Pet full lifecycle
**Test:** `TestScenario12PetFullLifecycle`  
**File:** `scenario_12_pet_full_lifecycle_test.go`  
**Feature:** `features/12_pet_full_lifecycle.feature`

Exercises the complete pet lifecycle in sequence:  
`POST /pet` → `GET /pet/{id}` → `PUT /pet` (status "sold") → `DELETE /pet/{id}` → `GET /pet/{id}`.  
Verifies: each step succeeds; the final GET returns 404.

---

### Store endpoints (scenarios 13–20)

---

### 13 — Get store inventory
**Test:** `TestScenario13GetInventory`  
**File:** `scenario_13_get_inventory_test.go`  
**Feature:** `features/13_get_inventory.feature`

Requests `GET /store/inventory`.  
Verifies: status 200, body deserialises to a non-empty `map[string]int`.

---

### 14 — Place a store order
**Test:** `TestScenario14PlaceOrder`  
**File:** `scenario_14_place_order_test.go`  
**Feature:** `features/14_place_order.feature`

Posts `{petId: 1, quantity: 1, status: "placed"}` to `POST /store/order`.  
Verifies: status 200, `id > 0`, `status == "placed"`.

---

### 15 — Get an order by ID
**Test:** `TestScenario15GetOrderByID`  
**File:** `scenario_15_get_order_by_id_test.go`  
**Feature:** `features/15_get_order_by_id.feature`

Creates a pet and an order, then fetches the order with `GET /store/order/{id}`.  
Verifies: status 200, returned `id`, `petId`, and `quantity` match the placed order.

---

### 16 — Delete an order
**Test:** `TestScenario16DeleteOrder`  
**File:** `scenario_16_delete_order_test.go`  
**Feature:** `features/16_delete_order.feature`

Creates an order, deletes it with `DELETE /store/order/{id}`, then requests `GET /store/order/{id}`.  
Verifies: DELETE returns 200; subsequent GET returns 404.

---

### 17 — Place a completed order
**Test:** `TestScenario17PlaceOrderComplete`  
**File:** `scenario_17_place_order_complete_test.go`  
**Feature:** `features/17_place_order_complete.feature`

Posts an order with `complete: true` to `POST /store/order`.  
Verifies: status 200, `id > 0`, `complete == true`.

---

### 18 — Get non-existent order
**Test:** `TestScenario18GetOrderNotFound`  
**File:** `scenario_18_get_order_not_found_test.go`  
**Feature:** `features/18_get_order_not_found.feature`

Requests `GET /store/order/999999` — an ID that cannot exist.  
Verifies: status 404.

---

### 19 — Place order with quantity
**Test:** `TestScenario19PlaceOrderWithQuantity`  
**File:** `scenario_19_place_order_with_quantity_test.go`  
**Feature:** `features/19_place_order_with_quantity.feature`

Posts `{petId: 1, quantity: 5, status: "placed"}` to `POST /store/order`.  
Verifies: status 200, `quantity == 5`.

---

### 20 — Order full lifecycle
**Test:** `TestScenario20OrderFullLifecycle`  
**File:** `scenario_20_order_full_lifecycle_test.go`  
**Feature:** `features/20_order_full_lifecycle.feature`

Exercises the complete order lifecycle:  
`POST /store/order` → `GET /store/order/{id}` → `DELETE /store/order/{id}` → `GET /store/order/{id}`.  
Verifies: each step succeeds; the final GET returns 404.

---

### User endpoints (scenarios 21–30)

---

### 21 — Create a user
**Test:** `TestScenario21CreateUser`  
**File:** `scenario_21_create_user_test.go`  
**Feature:** `features/21_create_user.feature`

Posts a new user to `POST /user`, then retrieves it with `GET /user/{username}`.  
Verifies: POST returns 200; GET returns 200 with matching `username` and `email`.

---

### 22 — Get user by username
**Test:** `TestScenario22GetUserByUsername`  
**File:** `scenario_22_get_user_by_username_test.go`  
**Feature:** `features/22_get_user_by_username.feature`

Creates a user via `createTestUser`, then fetches it with `GET /user/{username}`.  
Verifies: status 200, `username` matches.

---

### 23 — Update a user
**Test:** `TestScenario23UpdateUser`  
**File:** `scenario_23_update_user_test.go`  
**Feature:** `features/23_update_user.feature`

Creates a user, sends `PUT /user/{username}` with an updated email, then re-fetches with `GET /user/{username}`.  
Verifies: PUT returns 200; GET returns the updated email.

---

### 24 — Delete a user
**Test:** `TestScenario24DeleteUser`  
**File:** `scenario_24_delete_user_test.go`  
**Feature:** `features/24_delete_user.feature`

Creates a user, deletes with `DELETE /user/{username}`, then requests `GET /user/{username}`.  
Verifies: DELETE returns 200; subsequent GET returns 404.

---

### 25 — Login with valid credentials
**Test:** `TestScenario25LoginValid`  
**File:** `scenario_25_login_valid_test.go`  
**Feature:** `features/25_login_valid.feature`

Sends `GET /user/login?username=test&password=abc123`.  
Verifies: status 200, response body is non-empty (session token present).

---

### 26 — Logout
**Test:** `TestScenario26Logout`  
**File:** `scenario_26_logout_test.go`  
**Feature:** `features/26_logout.feature`

Sends `GET /user/logout`.  
Verifies: status 200.

---

### 27 — Create users with array
**Test:** `TestScenario27CreateUsersWithArray`  
**File:** `scenario_27_create_users_with_array_test.go`  
**Feature:** `features/27_create_users_with_array.feature`

Posts a JSON array of 2 users to `POST /user/createWithArray`, then fetches each with `GET /user/{username}`.  
Verifies: POST returns 200; each subsequent GET returns 200.

---

### 28 — Create users with list
**Test:** `TestScenario28CreateUsersWithList`  
**File:** `scenario_28_create_users_with_list_test.go`  
**Feature:** `features/28_create_users_with_list.feature`

Posts a JSON array of 2 users to `POST /user/createWithList`, then fetches each with `GET /user/{username}`.  
Verifies: POST returns 200; each subsequent GET returns 200.

---

### 29 — Login with invalid credentials
**Test:** `TestScenario29LoginInvalid`  
**File:** `scenario_29_login_invalid_test.go`  
**Feature:** `features/29_login_invalid.feature`

Sends `GET /user/login` with a deliberately invalid username and password.  
Verifies: response status is not 2xx (`resp.OK() == false`).

---

### 30 — User full lifecycle
**Test:** `TestScenario30UserFullLifecycle`  
**File:** `scenario_30_user_full_lifecycle_test.go`  
**Feature:** `features/30_user_full_lifecycle.feature`

Exercises the complete user lifecycle:  
`POST /user` → `GET /user/login` → `GET /user/{username}` → `PUT /user/{username}` → `DELETE /user/{username}` → `GET /user/{username}`.  
Verifies: each step succeeds; the final GET returns 404.

---

## Test patterns

### Each test follows this structure

```go
func TestScenarioNNDescription(t *testing.T) {
    ctx    := testCtx(t)      // 30-second context, auto-cancelled
    apiCtx := newAPICtx(t)    // APIRequestContext → https://petstore.swagger.io/v2

    resp, err := apiCtx.Get(ctx, "/endpoint")
    require.NoError(t, err)
    defer resp.Dispose(ctx)

    require.Equal(t, 200, resp.Status())
    body, err := resp.Body(ctx)
    require.NoError(t, err)
    pet := mustUnmarshalPet(t, body)
    assert.Greater(t, pet.ID, int64(0))
}
```

### Cleanup is automatic

`createTestPet`, `createTestOrder`, and `createTestUser` register a `t.Cleanup` that deletes the
created resource after the test ends — whether the test passed or failed.

### Tests run in parallel

Every scenario calls `t.Parallel()` at the start, allowing the suite to run concurrently and
finish significantly faster than a sequential run against a remote API.
