# goaneco-playwright

Go library for test automation with [Microsoft Playwright](https://playwright.dev).
Write UI and API tests in Go using the same tool, running against **Chromium**, **Firefox**,
and **WebKit**.

This project is open source, released under the [MIT License](LICENSE). It was built primarily
for personal use, but ideas, suggestions, and improvements are welcome — contributions that make
the project more useful are always considered. You are also free to fork it and take it in
whatever direction you like.

---

## What can I test with this?

| Test type | Description | Included example |
|---|---|---|
| **UI / E2E** | Automates a real browser: clicks, forms, navigation, screenshots | [SauceDemo](examples/ui-automation/) — 10 scenarios |
| **API REST** | Sends HTTP requests (GET, POST, PUT, DELETE) and validates JSON responses — no browser needed | [Petstore API](examples/api-automation/) — 30 scenarios |

---

## Setup

### 1. Install the prerequisites

| Tool | Min version | How to install |
|---|---|---|
| [Go](https://go.dev/dl/) | 1.21 | Download and install from go.dev |
| [Node.js](https://nodejs.org/) | 18 | Download and install from nodejs.org |
| Playwright | — | `npm install -g playwright` |

### 2. Install the browsers

```bash
npx playwright install chromium
```

On Linux, also install the OS-level dependencies:

```bash
npx playwright install --with-deps chromium
```

### 3. Set the environment variable

All tests need to know where the Playwright CLI is located:

```bash
# Linux / macOS
export PLAYWRIGHT_CLI_PATH=$(npm root -g)/playwright/node_modules/playwright-core/cli.js

# Windows PowerShell
$env:PLAYWRIGHT_CLI_PATH = "$env:APPDATA\npm\node_modules\playwright\node_modules\playwright-core\cli.js"
```

Add this line to your shell profile (`~/.bashrc`, `~/.zshrc`, etc.) so you don't have to repeat it every session.

---

## What tests look like

Tests are organized as **scenarios** — one file per test case. Each test is a Go function that
clearly describes what it is verifying.

### UI test (with Page Object Model)

Target: [https://www.saucedemo.com](https://www.saucedemo.com) — a demo e-commerce site used for UI automation practice.

```go
func TestLoginAndInventory(t *testing.T) {
    must := require.New(t)
    is   := assert.New(t)
    page := newPage(t)           // opens a new browser tab
    loginStandardUser(t, page)   // helper that logs in with standard credentials

    ctx := testCtx(t)
    inv := pom.NewInventoryPage(page)

    count, err := inv.ProductCount(ctx)
    must.NoError(err)
    is.Equal(6, count)   // the inventory must show exactly 6 products
}
```

The **Page Object Model (POM)** wraps the actions of each screen into a reusable struct
(`pom.InventoryPage`, `pom.CartPage`, etc.), keeping tests clean and easy to read.

### API test (full lifecycle)

Target: [https://petstore.swagger.io/v2](https://petstore.swagger.io/v2) — a public REST API used as a sandbox for API automation practice.

```go
func TestScenario12PetFullLifecycle(t *testing.T) {
    must := require.New(t)
    is   := assert.New(t)
    ctx    := testCtx(t)
    apiCtx := newAPICtx(t)   // HTTP client pointed at the API base URL

    // Create a pet
    pet := &models.Pet{Name: "LifecyclePet", PhotoURLs: []string{"https://example.com/img.jpg"}, Status: "available"}
    postResp, err := apiCtx.Post(ctx, "/pet", &playwright.APIFetchOptions{
        Headers: jsonHeaders(),
        Data:    mustJSON(t, pet),
    })
    must.NoError(err)
    defer postResp.Dispose(ctx)
    must.Equal(200, postResp.Status())
    postBody, err := postResp.Body(ctx)
    must.NoError(err)
    created := mustUnmarshalPet(t, postBody)
    must.Greater(created.ID, int64(0))

    // Verify it can be retrieved
    getResp, err := apiCtx.Get(ctx, fmt.Sprintf("/pet/%d", created.ID))
    must.NoError(err)
    defer getResp.Dispose(ctx)
    is.Equal(200, getResp.Status())

    // Update the status to "sold"
    created.Status = "sold"
    putResp, err := apiCtx.Put(ctx, "/pet", &playwright.APIFetchOptions{Headers: jsonHeaders(), Data: mustJSON(t, created)})
    must.NoError(err)
    defer putResp.Dispose(ctx)
    is.Equal(200, putResp.Status())

    // Delete the pet
    delResp, err := apiCtx.Delete(ctx, fmt.Sprintf("/pet/%d", created.ID))
    must.NoError(err)
    defer delResp.Dispose(ctx)
    is.Equal(200, delResp.Status())

    // Confirm it no longer exists (404)
    confirmResp, err := apiCtx.Get(ctx, fmt.Sprintf("/pet/%d", created.ID))
    must.NoError(err)
    defer confirmResp.Dispose(ctx)
    is.Equal(404, confirmResp.Status())
}
```

---

## Included examples

### UI automation — SauceDemo (`examples/ui-automation/`)

Browser automation for [SauceDemo](https://www.saucedemo.com) using the **Page Object Model**
pattern. The 10 scenarios cover the complete flow of an online store.

```bash
# Run all scenarios
PLAYWRIGHT_CLI_PATH="..." go test -tags e2e -v -timeout 300s ./examples/ui-automation/...

# Run a single scenario
PLAYWRIGHT_CLI_PATH="..." go test -tags e2e -v -timeout 120s -run TestFullCheckoutFlow ./examples/ui-automation/...
```

| # | Scenario | What it verifies |
|---|---|---|
| 01 | Login and inventory | Authentication smoke test; confirms 6 products are visible |
| 02 | Add items to cart | Adds 3 products; verifies cart badge count and item list |
| 03 | Sort by price low to high | The cheapest item (Onesie, $7.99) appears first |
| 04 | Sort by name Z to A | Products are listed in reverse alphabetical order |
| 05 | Full checkout flow | Add → cart → fill info → overview → order confirmation |
| 06 | Remove item from cart | Adds 3 items, removes 1, confirms 2 remain |
| 07 | Product detail page | Verifies name, price, and description of a product |
| 08 | Logout via burger menu | Opens side nav, clicks Logout, confirms login screen |
| 09 | Add most expensive item | Sorts high to low; adds Fleece Jacket ($49.99) to cart |
| 10 | Footer social links | Verifies visible Twitter, Facebook, and LinkedIn links |

Screens are abstracted as **Page Objects** in `examples/ui-automation/pom/`:

| Page Object | Available actions |
|---|---|
| `LoginPage` | `Navigate`, `Login`, `IsVisible`, `ErrorMessage` |
| `InventoryPage` | `ProductCount`, `AddToCart`, `CartBadgeCount`, `Sort`, `Logout` |
| `CartPage` | `ItemCount`, `GetItemNames`, `GetItemPrices`, `RemoveItem`, `Checkout` |
| `CheckoutInfoPage` | `FillInfo`, `Continue`, `ErrorMessage` |
| `CheckoutOverviewPage` | `GetItemCount`, `GetSubtotal`, `GetTotal`, `Finish` |
| `CheckoutCompletePage` | `GetHeader`, `BackHome` |
| `ProductDetailPage` | `GetName`, `GetPrice`, `GetDescription`, `AddToCart`, `BackToInventory` |

Full scenario descriptions: [`examples/ui-automation/ui-automation.md`](examples/ui-automation/ui-automation.md)

---

### API automation — Petstore (`examples/api-automation/`)

30 test scenarios against the public [Petstore API](https://petstore.swagger.io/v2).
No browser is launched — only HTTP requests and response assertions.

```bash
# Run all scenarios
PLAYWRIGHT_CLI_PATH="..." go test -tags e2e -v -timeout 300s ./examples/api-automation/...

# Run a single scenario
PLAYWRIGHT_CLI_PATH="..." go test -tags e2e -v -timeout 60s -run TestScenario01AddPet ./examples/api-automation/...
```

#### Pet — `/pet` endpoints (01–12)

| # | Scenario | Endpoint | Verifies |
|---|---|---|---|
| 01 | Add a pet | `POST /pet` | `id > 0`, name and status match |
| 02 | Get by ID | `GET /pet/{id}` | Returned fields match created pet |
| 03 | Update a pet | `PUT /pet` | Status updated to `"sold"` |
| 04 | Delete a pet | `DELETE /pet/{id}` | Subsequent GET returns 404 |
| 05–07 | Find by status | `GET /pet/findByStatus` | Table-driven: available / pending / sold |
| 08 | Find by multiple statuses | `GET /pet/findByStatus?status=available&status=pending` | Only those two statuses in response |
| 09 | Update with form data | `POST /pet/{id}` (form-encoded) | Status 200 |
| 10 | Find by tag | `GET /pet/findByTags?tags=goaneco` | Valid JSON array |
| 11 | Non-existent pet | `GET /pet/999999999` | Status 404 |
| 12 | Full lifecycle | `POST → GET → PUT → DELETE → GET` | Final GET returns 404 |

#### Store — `/store` endpoints (13–20)

| # | Scenario | Endpoint | Verifies |
|---|---|---|---|
| 13 | Get inventory | `GET /store/inventory` | Non-empty status-count map |
| 14 | Place an order | `POST /store/order` | `id > 0`, `status == "placed"` |
| 15 | Get order by ID | `GET /store/order/{id}` | Fields match placed order |
| 16 | Delete an order | `DELETE /store/order/{id}` | Subsequent GET returns 404 |
| 17 | Place completed order | `POST /store/order` (`complete: true`) | `complete == true` in response |
| 18 | Non-existent order | `GET /store/order/999999` | Status 404 |
| 19 | Order with quantity | `POST /store/order` (`quantity: 5`) | `quantity == 5` in response |
| 20 | Full order lifecycle | `POST → GET → DELETE → GET` | Final GET returns 404 |

#### Users — `/user` endpoints (21–30)

| # | Scenario | Endpoint | Verifies |
|---|---|---|---|
| 21 | Create a user | `POST /user` + `GET /user/{username}` | Fields match |
| 22 | Get by username | `GET /user/{username}` | `username` matches |
| 23 | Update a user | `PUT /user/{username}` + `GET` | Updated email returned |
| 24 | Delete a user | `DELETE /user/{username}` | Subsequent GET returns 404 |
| 25 | Login with valid credentials | `GET /user/login` | Status 200, token present in body |
| 26 | Logout | `GET /user/logout` | Status 200 |
| 27 | Create users with array | `POST /user/createWithArray` | Each user retrievable via GET |
| 28 | Create users with list | `POST /user/createWithList` | Each user retrievable via GET |
| 29 | Login with invalid credentials | `GET /user/login` (bad credentials) | Response is not 2xx |
| 30 | Full user lifecycle | `POST → login → GET → PUT → DELETE → GET` | Final GET returns 404 |

Full scenario descriptions: [`examples/api-automation/api-automation.md`](examples/api-automation/api-automation.md)

---

## Running against multiple browsers

UI tests can run on Chromium (default), Firefox, or WebKit:

```bash
# Firefox
PLAYWRIGHT_CLI_PATH="..." PLAYWRIGHT_BROWSER=firefox go test -tags e2e -v -timeout 300s ./examples/ui-automation/...

# WebKit (Safari engine)
PLAYWRIGHT_CLI_PATH="..." PLAYWRIGHT_BROWSER=webkit go test -tags e2e -v -timeout 300s ./examples/ui-automation/...
```

---

## Test coverage

This library is validated by **3,318 tests** spread across three layers:

| Layer | Tests | What it covers |
|---|---|---|
| **Library E2E suite** (`e2e/`) | 3,278 | Core Playwright capabilities: browser, page, locator, frames, network, API requests, assertions, tracing, screenshots, and more — 124 test files covering the full API surface |
| **UI automation examples** (`examples/ui-automation/`) | 10 | End-to-end browser flows on SauceDemo using the Page Object Model pattern |
| **API automation examples** (`examples/api-automation/`) | 30 | REST API scenarios against the public Petstore API — no browser required |

The library E2E suite was built using [playwright-java](https://github.com/microsoft/playwright-java) as a reference, porting and adapting its test cases to Go to ensure equivalent behavior coverage.

---

## License

MIT License — see [LICENSE](LICENSE) for details.

[![Go Reference](https://pkg.go.dev/badge/github.com/jcfrperu/goaneco-playwright.svg)](https://pkg.go.dev/github.com/jcfrperu/goaneco-playwright)
[![Go Report Card](https://goreportcard.com/badge/github.com/jcfrperu/goaneco-playwright)](https://goreportcard.com/report/github.com/jcfrperu/goaneco-playwright)