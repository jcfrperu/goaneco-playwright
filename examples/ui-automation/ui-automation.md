# UI Automation Examples

End-to-end automation examples for [SauceDemo](https://www.saucedemo.com) using the
[goaneco-playwright](https://github.com/jcfrperu/goaneco-playwright) Go client.  
All examples follow the **Page Object Model (POM)** pattern and are tagged with `//go:build e2e`.

---

## Prerequisites

| Requirement | Details |
|---|---|
| Go | 1.21+ |
| Node.js | 18+ |
| Playwright | `npm install -g playwright` or `npm install -g playwright-core` |
| Internet access | Tests reach `https://www.saucedemo.com` |

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
# Run all 10 scenarios
go test -tags e2e -v -timeout 300s ./examples/ui-automation/...

# Run a single scenario by test name
go test -tags e2e -v -timeout 120s -run TestFullCheckoutFlow ./examples/ui-automation/...
```

The test suite launches a **headless Chromium** browser once via `TestMain`, shares it across
all tests, and closes it when the suite finishes.

---

## Project structure

```
examples/ui-automation/
├── doc.go                                  # package uiautomation — package declaration
├── main_test.go                            # TestMain: browser lifecycle, newPage, testCtx
├── helpers_test.go                         # Shared loginStandardUser helper
│
├── scenario_01_login_inventory_test.go
├── scenario_02_add_items_to_cart_test.go
├── scenario_03_sort_price_low_high_test.go
├── scenario_04_sort_name_z_to_a_test.go
├── scenario_05_full_checkout_flow_test.go
├── scenario_06_remove_item_from_cart_test.go
├── scenario_07_product_detail_page_test.go
├── scenario_08_logout_test.go
├── scenario_09_add_most_expensive_item_test.go
├── scenario_10_footer_social_links_test.go
│
├── pom/                                    # package pom — Page Object Model structs
│   ├── doc.go
│   ├── login_page.go
│   ├── inventory_page.go
│   ├── cart_page.go
│   ├── checkout_info_page.go
│   ├── checkout_overview_page.go
│   ├── checkout_complete_page.go
│   └── product_detail_page.go
│
├── ui-automation.md
└── features/
    ├── 01_login_inventory.feature
    ├── 02_add_items_to_cart.feature
    ├── 03_sort_price_low_high.feature
    ├── 04_sort_name_z_to_a.feature
    ├── 05_full_checkout_flow.feature
    ├── 06_remove_item_from_cart.feature
    ├── 07_product_detail_page.feature
    ├── 08_logout.feature
    ├── 09_add_most_expensive_item.feature
    └── 10_footer_social_links.feature
```

The `pom/` subpackage is a non-test package (`package pom`), importable at:

```
github.com/jcfrperu/goaneco-playwright/examples/ui-automation/pom
```

This separates reusable page abstractions from test orchestration, following the same pattern
used by the `e2e/testserver` package in this repository.

---

## Scenarios

### 01 — Login and verify inventory
**File:** `scenario_01_login_inventory_test.go`  
**Feature:** `features/01_login_inventory.feature`

Logs in with `standard_user` and verifies that the inventory page lists exactly **6 products**
and that the sort dropdown is visible. Serves as a basic smoke test for authentication.

---

### 02 — Add items to cart
**File:** `scenario_02_add_items_to_cart_test.go`  
**Feature:** `features/02_add_items_to_cart.feature`

Adds three products (Backpack, Bike Light, Bolt T-Shirt) to the cart from the inventory page.
Verifies the **cart badge count** (3), navigates to the cart, and confirms all item names and
prices are listed correctly.

---

### 03 — Sort by price low to high
**File:** `scenario_03_sort_price_low_high_test.go`  
**Feature:** `features/03_sort_price_low_high.feature`

Selects the "Price (low to high)" sort option and asserts that the first product shown is
**Sauce Labs Onesie at $7.99** — the cheapest item in the catalog.

---

### 04 — Sort by name Z to A
**File:** `scenario_04_sort_name_z_to_a_test.go`  
**Feature:** `features/04_sort_name_z_to_a.feature`

Selects "Name (Z to A)" and verifies that the list starts with **Test.allTheThings() T-Shirt (Red)**
and ends with **Sauce Labs Backpack**, confirming reverse-alphabetical ordering.

---

### 05 — Full checkout flow
**File:** `scenario_05_full_checkout_flow_test.go`  
**Feature:** `features/05_full_checkout_flow.feature`

Exercises the complete purchase path: add one item → proceed to cart → fill checkout info
(first name, last name, postal code) → review order overview → finish.  
Confirms the order confirmation page displays **"Thank you for your order!"**.

---

### 06 — Remove item from cart
**File:** `scenario_06_remove_item_from_cart_test.go`  
**Feature:** `features/06_remove_item_from_cart.feature`

Adds 3 items to the cart, navigates to the cart, removes **Sauce Labs Bike Light**, and
confirms that only 2 items remain and the removed item is no longer listed.

---

### 07 — Product detail page
**File:** `scenario_07_product_detail_page_test.go`  
**Feature:** `features/07_product_detail_page.feature`

Clicks on **Sauce Labs Backpack** from the inventory to open its detail page. Verifies the
product name, price ($29.99), and non-empty description. Then navigates back to the inventory
and confirms all 6 products are still displayed.

---

### 08 — Logout via burger menu
**File:** `scenario_08_logout_test.go`  
**Feature:** `features/08_logout.feature`

Opens the hamburger side-navigation menu and clicks **Logout**. Verifies that the browser is
redirected to the login page and that the login form is visible again.

---

### 09 — Add the most expensive item
**File:** `scenario_09_add_most_expensive_item_test.go`  
**Feature:** `features/09_add_most_expensive_item.feature`

Sorts by "Price (high to low)" and adds the first result — **Sauce Labs Fleece Jacket at $49.99** —
to the cart. Confirms the item name and price appear correctly in the cart page.

---

### 10 — Footer social media links
**File:** `scenario_10_footer_social_links_test.go`  
**Feature:** `features/10_footer_social_links.feature`

Verifies that the inventory page footer contains **visible, non-empty links** to Twitter,
Facebook, and LinkedIn — a structural check to detect footer regressions.

---

## Page Object Model overview

Page objects live in `pom/` (import `github.com/jcfrperu/goaneco-playwright/examples/ui-automation/pom`).  
Each struct holds a `*playwright.Page` reference. Constructors follow the pattern `pom.NewXxxPage(page)`.

| Page object | Struct | Key methods |
|---|---|---|
| Login | `pom.LoginPage` | `Navigate`, `Login`, `IsVisible`, `ErrorMessage` |
| Inventory | `pom.InventoryPage` | `ProductCount`, `GetProductNames`, `AddToCart`, `CartBadgeCount`, `GoToCart`, `Sort`, `Logout`, `ClickProduct` |
| Cart | `pom.CartPage` | `ItemCount`, `GetItemNames`, `GetItemPrices`, `RemoveItem`, `Checkout`, `ContinueShopping` |
| Checkout info | `pom.CheckoutInfoPage` | `FillInfo`, `Continue`, `ErrorMessage` |
| Checkout overview | `pom.CheckoutOverviewPage` | `GetItemCount`, `GetItemNames`, `GetSubtotal`, `GetTotal`, `Finish`, `Cancel` |
| Checkout complete | `pom.CheckoutCompletePage` | `GetHeader`, `GetSubHeader`, `BackHome` |
| Product detail | `pom.ProductDetailPage` | `GetName`, `GetPrice`, `GetDescription`, `AddToCart`, `BackToInventory` |

All methods accept a `context.Context` as the first argument and return an `error`.
Use `testCtx(t)` (90-second timeout) for contexts in tests.
