# goaneco-playwright

Librería Go para automatización de pruebas con [Microsoft Playwright](https://playwright.dev).
Permite escribir tests de UI y de API en Go con la misma herramienta, ejecutándolos en
**Chromium**, **Firefox** y **WebKit**.

Este proyecto es open source, publicado bajo la [licencia MIT](LICENSE). Fue creado principalmente
para uso personal, pero se aceptan ideas, sugerencias y mejoras — las contribuciones que hagan el
proyecto más útil siempre serán consideradas. También eres libre de hacer un fork y llevarlo en
la dirección que prefieras.

---

## ¿Qué puedo probar con esto?

| Tipo de test | Descripción | Ejemplo incluido |
|---|---|---|
| **UI / E2E** | Automatiza un navegador real: clics, formularios, navegación, capturas de pantalla | [SauceDemo](examples/ui-automation/) — 10 escenarios |
| **API REST** | Envía peticiones HTTP (GET, POST, PUT, DELETE) y valida respuestas JSON — sin navegador | [Petstore API](examples/api-automation/) — 30 escenarios |

---

## Configuración

### 1. Instalar los prerrequisitos

| Herramienta | Versión mínima | Cómo instalar |
|---|---|---|
| [Go](https://go.dev/dl/) | 1.21 | Descargar e instalar desde go.dev |
| [Node.js](https://nodejs.org/) | 18 | Descargar e instalar desde nodejs.org |
| Playwright | — | `npm install -g playwright` |

### 2. Instalar los navegadores

```bash
npx playwright install chromium
```

En Linux, instalar también las dependencias del sistema operativo:

```bash
npx playwright install --with-deps chromium
```

### 3. Configurar la variable de entorno

Todos los tests necesitan saber dónde está el CLI de Playwright:

```bash
# Linux / macOS
export PLAYWRIGHT_CLI_PATH=$(npm root -g)/playwright/node_modules/playwright-core/cli.js

# Windows PowerShell
$env:PLAYWRIGHT_CLI_PATH = "$env:APPDATA\npm\node_modules\playwright\node_modules\playwright-core\cli.js"
```

Agrega esta línea a tu perfil de shell (`~/.bashrc`, `~/.zshrc`, etc.) para no repetirla en cada sesión.

---

## Cómo se ven los tests

Los tests se organizan en **escenarios** — un archivo por caso de prueba. Cada test es una función
Go que describe claramente qué está verificando.

### Test de UI (con Page Object Model)

Objetivo: [https://www.saucedemo.com](https://www.saucedemo.com) — sitio de e-commerce de demo para práctica de automatización UI.

```go
func TestLoginAndInventory(t *testing.T) {
    must := require.New(t)
    is   := assert.New(t)
    page := newPage(t)           // abre una nueva pestaña del navegador
    loginStandardUser(t, page)   // helper que hace login con credenciales estándar

    ctx := testCtx(t)
    inv := pom.NewInventoryPage(page)

    count, err := inv.ProductCount(ctx)
    must.NoError(err)
    is.Equal(6, count)   // el inventario debe mostrar exactamente 6 productos
}
```

El **Page Object Model (POM)** encapsula las acciones de cada pantalla en un struct reutilizable
(`pom.InventoryPage`, `pom.CartPage`, etc.), manteniendo los tests limpios y fáciles de leer.

### Test de API (ciclo completo)

Objetivo: [https://petstore.swagger.io/v2](https://petstore.swagger.io/v2) — API REST pública usada como sandbox para práctica de automatización de APIs.

```go
func TestScenario12PetFullLifecycle(t *testing.T) {
    must := require.New(t)
    is   := assert.New(t)
    ctx    := testCtx(t)
    apiCtx := newAPICtx(t)   // cliente HTTP apuntando a la URL base de la API

    // Crear una mascota
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

    // Verificar que se puede recuperar
    getResp, err := apiCtx.Get(ctx, fmt.Sprintf("/pet/%d", created.ID))
    must.NoError(err)
    defer getResp.Dispose(ctx)
    is.Equal(200, getResp.Status())

    // Actualizar el estado a "sold"
    created.Status = "sold"
    putResp, err := apiCtx.Put(ctx, "/pet", &playwright.APIFetchOptions{Headers: jsonHeaders(), Data: mustJSON(t, created)})
    must.NoError(err)
    defer putResp.Dispose(ctx)
    is.Equal(200, putResp.Status())

    // Eliminar la mascota
    delResp, err := apiCtx.Delete(ctx, fmt.Sprintf("/pet/%d", created.ID))
    must.NoError(err)
    defer delResp.Dispose(ctx)
    is.Equal(200, delResp.Status())

    // Confirmar que ya no existe (404)
    confirmResp, err := apiCtx.Get(ctx, fmt.Sprintf("/pet/%d", created.ID))
    must.NoError(err)
    defer confirmResp.Dispose(ctx)
    is.Equal(404, confirmResp.Status())
}
```

---

## Ejemplos incluidos

### Automatización UI — SauceDemo (`examples/ui-automation/`)

Automatización de navegador para [SauceDemo](https://www.saucedemo.com) usando el patrón
**Page Object Model**. Los 10 escenarios cubren el flujo completo de una tienda online.

```bash
# Ejecutar todos los escenarios
PLAYWRIGHT_CLI_PATH="..." go test -tags e2e -v -timeout 300s ./examples/ui-automation/...

# Ejecutar un escenario específico
PLAYWRIGHT_CLI_PATH="..." go test -tags e2e -v -timeout 120s -run TestFullCheckoutFlow ./examples/ui-automation/...
```

| # | Escenario | Qué verifica |
|---|---|---|
| 01 | Login e inventario | Smoke test de autenticación; confirma 6 productos visibles |
| 02 | Agregar ítems al carrito | Agrega 3 productos; verifica badge y lista del carrito |
| 03 | Ordenar por precio (menor a mayor) | El ítem más barato (Onesie, $7.99) aparece primero |
| 04 | Ordenar por nombre (Z a A) | Productos en orden alfabético inverso |
| 05 | Flujo completo de checkout | Agregar → carrito → datos → resumen → confirmación |
| 06 | Eliminar ítem del carrito | Agrega 3, elimina 1, confirma que quedan 2 |
| 07 | Página de detalle del producto | Verifica nombre, precio y descripción |
| 08 | Logout desde el menú | Abre menú lateral, hace logout, confirma pantalla de login |
| 09 | Agregar el ítem más caro | Ordena de mayor a menor; agrega Fleece Jacket ($49.99) |
| 10 | Links de redes sociales en footer | Verifica links visibles de Twitter, Facebook y LinkedIn |

Las pantallas están abstraídas como **Page Objects** en `examples/ui-automation/pom/`:

| Page Object | Acciones disponibles |
|---|---|
| `LoginPage` | `Navigate`, `Login`, `IsVisible`, `ErrorMessage` |
| `InventoryPage` | `ProductCount`, `AddToCart`, `CartBadgeCount`, `Sort`, `Logout` |
| `CartPage` | `ItemCount`, `GetItemNames`, `GetItemPrices`, `RemoveItem`, `Checkout` |
| `CheckoutInfoPage` | `FillInfo`, `Continue`, `ErrorMessage` |
| `CheckoutOverviewPage` | `GetItemCount`, `GetSubtotal`, `GetTotal`, `Finish` |
| `CheckoutCompletePage` | `GetHeader`, `BackHome` |
| `ProductDetailPage` | `GetName`, `GetPrice`, `GetDescription`, `AddToCart`, `BackToInventory` |

Descripción detallada de escenarios: [`examples/ui-automation/ui-automation.md`](examples/ui-automation/ui-automation.md)

---

### Automatización API — Petstore (`examples/api-automation/`)

30 escenarios de prueba contra la [Petstore API](https://petstore.swagger.io/v2) pública.
No se lanza ningún navegador — solo peticiones HTTP y validaciones de respuesta.

```bash
# Ejecutar todos los escenarios
PLAYWRIGHT_CLI_PATH="..." go test -tags e2e -v -timeout 300s ./examples/api-automation/...

# Ejecutar un escenario específico
PLAYWRIGHT_CLI_PATH="..." go test -tags e2e -v -timeout 60s -run TestScenario01AddPet ./examples/api-automation/...
```

#### Mascotas — endpoints `/pet` (01–12)

| # | Escenario | Endpoint | Verifica |
|---|---|---|---|
| 01 | Agregar mascota | `POST /pet` | `id > 0`, nombre y estado correctos |
| 02 | Obtener por ID | `GET /pet/{id}` | Campos coinciden con la mascota creada |
| 03 | Actualizar mascota | `PUT /pet` | Estado actualizado a `"sold"` |
| 04 | Eliminar mascota | `DELETE /pet/{id}` | GET posterior retorna 404 |
| 05–07 | Buscar por estado | `GET /pet/findByStatus` | available / pending / sold |
| 08 | Buscar por múltiples estados | `GET /pet/findByStatus?status=available&status=pending` | Solo esos dos estados |
| 09 | Actualizar con form data | `POST /pet/{id}` (form-encoded) | Status 200 |
| 10 | Buscar por etiqueta | `GET /pet/findByTags?tags=goaneco` | Array JSON válido |
| 11 | Mascota inexistente | `GET /pet/999999999` | Status 404 |
| 12 | Ciclo completo | `POST → GET → PUT → DELETE → GET` | GET final retorna 404 |

#### Tienda — endpoints `/store` (13–20)

| # | Escenario | Endpoint | Verifica |
|---|---|---|---|
| 13 | Obtener inventario | `GET /store/inventory` | Mapa de estados no vacío |
| 14 | Crear orden | `POST /store/order` | `id > 0`, `status == "placed"` |
| 15 | Obtener orden por ID | `GET /store/order/{id}` | Campos coinciden con la orden |
| 16 | Eliminar orden | `DELETE /store/order/{id}` | GET posterior retorna 404 |
| 17 | Orden completada | `POST /store/order` (`complete: true`) | `complete == true` |
| 18 | Orden inexistente | `GET /store/order/999999` | Status 404 |
| 19 | Orden con cantidad | `POST /store/order` (`quantity: 5`) | `quantity == 5` |
| 20 | Ciclo completo de orden | `POST → GET → DELETE → GET` | GET final retorna 404 |

#### Usuarios — endpoints `/user` (21–30)

| # | Escenario | Endpoint | Verifica |
|---|---|---|---|
| 21 | Crear usuario | `POST /user` + `GET /user/{username}` | Campos coinciden |
| 22 | Obtener por username | `GET /user/{username}` | `username` coincide |
| 23 | Actualizar usuario | `PUT /user/{username}` + `GET` | Email actualizado |
| 24 | Eliminar usuario | `DELETE /user/{username}` | GET posterior retorna 404 |
| 25 | Login válido | `GET /user/login` | Status 200, token en body |
| 26 | Logout | `GET /user/logout` | Status 200 |
| 27 | Crear usuarios con array | `POST /user/createWithArray` | Cada usuario recuperable |
| 28 | Crear usuarios con lista | `POST /user/createWithList` | Cada usuario recuperable |
| 29 | Login inválido | `GET /user/login` (credenciales incorrectas) | Respuesta no es 2xx |
| 30 | Ciclo completo de usuario | `POST → login → GET → PUT → DELETE → GET` | GET final retorna 404 |

Descripción detallada de escenarios: [`examples/api-automation/api-automation.md`](examples/api-automation/api-automation.md)

---

## Ejecutar en múltiples navegadores

Los tests de UI pueden correr en Chromium (por defecto), Firefox o WebKit:

```bash
# Firefox
PLAYWRIGHT_CLI_PATH="..." PLAYWRIGHT_BROWSER=firefox go test -tags e2e -v -timeout 300s ./examples/ui-automation/...

# WebKit (motor de Safari)
PLAYWRIGHT_CLI_PATH="..." PLAYWRIGHT_BROWSER=webkit go test -tags e2e -v -timeout 300s ./examples/ui-automation/...
```

---

## Cobertura de tests

Esta librería está validada por **3.318 tests** distribuidos en tres capas:

| Capa | Tests | Qué cubre |
|---|---|---|
| **Suite E2E de la librería** (`e2e/`) | 3.278 | Capacidades core de Playwright: browser, page, locator, frames, red, solicitudes API, assertions, tracing, capturas de pantalla y más — 124 archivos de test que cubren toda la superficie de la API |
| **Ejemplos UI automation** (`examples/ui-automation/`) | 10 | Flujos completos de navegador en SauceDemo usando el patrón Page Object Model |
| **Ejemplos API automation** (`examples/api-automation/`) | 30 | Escenarios REST API contra la Petstore API pública — sin necesidad de navegador |

La suite E2E de la librería fue construida usando [playwright-java](https://github.com/microsoft/playwright-java) como referencia, portando y adaptando sus casos de prueba a Go para garantizar una cobertura de comportamiento equivalente.

---

## Licencia

Licencia MIT — ver [LICENSE](LICENSE) para más detalles.

[![Go Reference](https://pkg.go.dev/badge/github.com/jcfrperu/goaneco-playwright.svg)](https://pkg.go.dev/github.com/jcfrperu/goaneco-playwright)
[![Go Report Card](https://goreportcard.com/badge/github.com/jcfrperu/goaneco-playwright)](https://goreportcard.com/report/github.com/jcfrperu/goaneco-playwright)