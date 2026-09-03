# goaneco-playwright

Go-Bibliothek für Testautomatisierung mit [Microsoft Playwright](https://playwright.dev).
Schreibe UI- und API-Tests in Go mit demselben Werkzeug, ausgeführt auf **Chromium**, **Firefox**
und **WebKit**.

Dieses Projekt ist Open Source und steht unter der [MIT-Lizenz](LICENSE). Es wurde hauptsächlich
für den persönlichen Gebrauch entwickelt, aber Ideen, Vorschläge und Verbesserungen sind willkommen
— Beiträge, die das Projekt nützlicher machen, werden immer berücksichtigt. Es steht dir auch frei,
das Projekt zu forken und in eine eigene Richtung weiterzuentwickeln.

---

## Was kann ich damit testen?

| Testtyp | Beschreibung | Enthaltenes Beispiel |
|---|---|---|
| **UI / E2E** | Automatisiert einen echten Browser: Klicks, Formulare, Navigation, Screenshots | [SauceDemo](examples/ui-automation/) — 10 Szenarien |
| **REST API** | Sendet HTTP-Anfragen (GET, POST, PUT, DELETE) und validiert JSON-Antworten — kein Browser nötig | [Petstore API](examples/api-automation/) — 30 Szenarien |

---

## Einrichtung

### 1. Voraussetzungen installieren

| Werkzeug | Mindestversion | Installation |
|---|---|---|
| [Go](https://go.dev/dl/) | 1.21 | Herunterladen und installieren von go.dev |
| [Node.js](https://nodejs.org/) | 18 | Herunterladen und installieren von nodejs.org |
| Playwright | — | `npm install -g playwright` |

### 2. Browser installieren

```bash
npx playwright install chromium
```

Unter Linux auch die Abhängigkeiten des Betriebssystems installieren:

```bash
npx playwright install --with-deps chromium
```

### 3. Umgebungsvariable setzen

Alle Tests müssen wissen, wo sich der Playwright CLI befindet:

```bash
# Linux / macOS
export PLAYWRIGHT_CLI_PATH=$(npm root -g)/playwright/node_modules/playwright-core/cli.js

# Windows PowerShell
$env:PLAYWRIGHT_CLI_PATH = "$env:APPDATA\npm\node_modules\playwright\node_modules\playwright-core\cli.js"
```

Diese Zeile zum Shell-Profil hinzufügen (`~/.bashrc`, `~/.zshrc` usw.), um sie nicht bei jeder Sitzung wiederholen zu müssen.

---

## Wie Tests aussehen

Tests sind als **Szenarien** organisiert — eine Datei pro Testfall. Jeder Test ist eine Go-Funktion,
die klar beschreibt, was sie überprüft.

### UI-Test (mit Page Object Model)

Ziel: [https://www.saucedemo.com](https://www.saucedemo.com) — eine Demo-E-Commerce-Seite für UI-Automatisierungsübungen.

```go
func TestLoginAndInventory(t *testing.T) {
    must := require.New(t)
    is   := assert.New(t)
    page := newPage(t)           // öffnet einen neuen Browser-Tab
    loginStandardUser(t, page)   // Hilfsfunktion, die sich mit Standarddaten anmeldet

    ctx := testCtx(t)
    inv := pom.NewInventoryPage(page)

    count, err := inv.ProductCount(ctx)
    must.NoError(err)
    is.Equal(6, count)   // das Inventar muss genau 6 Produkte anzeigen
}
```

Das **Page Object Model (POM)** kapselt die Aktionen jedes Bildschirms in einer wiederverwendbaren
Struct (`pom.InventoryPage`, `pom.CartPage` usw.) und hält Tests sauber und lesbar.

### API-Test (vollständiger Lebenszyklus)

Ziel: [https://petstore.swagger.io/v2](https://petstore.swagger.io/v2) — eine öffentliche REST API als Sandbox für API-Automatisierungsübungen.

```go
func TestScenario12PetFullLifecycle(t *testing.T) {
    must := require.New(t)
    is   := assert.New(t)
    ctx    := testCtx(t)
    apiCtx := newAPICtx(t)   // HTTP-Client, der auf die API-Basis-URL zeigt

    // Ein Tier erstellen
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

    // Prüfen, ob es abgerufen werden kann
    getResp, err := apiCtx.Get(ctx, fmt.Sprintf("/pet/%d", created.ID))
    must.NoError(err)
    defer getResp.Dispose(ctx)
    is.Equal(200, getResp.Status())

    // Status auf "sold" aktualisieren
    created.Status = "sold"
    putResp, err := apiCtx.Put(ctx, "/pet", &playwright.APIFetchOptions{Headers: jsonHeaders(), Data: mustJSON(t, created)})
    must.NoError(err)
    defer putResp.Dispose(ctx)
    is.Equal(200, putResp.Status())

    // Das Tier löschen
    delResp, err := apiCtx.Delete(ctx, fmt.Sprintf("/pet/%d", created.ID))
    must.NoError(err)
    defer delResp.Dispose(ctx)
    is.Equal(200, delResp.Status())

    // Bestätigen, dass es nicht mehr existiert (404)
    confirmResp, err := apiCtx.Get(ctx, fmt.Sprintf("/pet/%d", created.ID))
    must.NoError(err)
    defer confirmResp.Dispose(ctx)
    is.Equal(404, confirmResp.Status())
}
```

---

## Enthaltene Beispiele

### UI-Automatisierung — SauceDemo (`examples/ui-automation/`)

Browser-Automatisierung für [SauceDemo](https://www.saucedemo.com) mit dem **Page Object Model**.
Die 10 Szenarien decken den vollständigen Ablauf eines Online-Shops ab.

```bash
# Alle Szenarien ausführen
PLAYWRIGHT_CLI_PATH="..." go test -tags e2e -v -timeout 300s ./examples/ui-automation/...

# Ein einzelnes Szenario ausführen
PLAYWRIGHT_CLI_PATH="..." go test -tags e2e -v -timeout 120s -run TestFullCheckoutFlow ./examples/ui-automation/...
```

| # | Szenario | Was wird geprüft |
|---|---|---|
| 01 | Login und Inventar | Authentifizierungs-Smoke-Test; bestätigt 6 sichtbare Produkte |
| 02 | Artikel in den Warenkorb legen | Fügt 3 Produkte hinzu; prüft Warenkorb-Badge und Artikelliste |
| 03 | Nach Preis sortieren (aufsteigend) | Günstigstes Produkt (Onesie, $7.99) erscheint zuerst |
| 04 | Nach Name sortieren (Z bis A) | Produkte in umgekehrter alphabetischer Reihenfolge |
| 05 | Vollständiger Checkout-Ablauf | Hinzufügen → Warenkorb → Daten → Übersicht → Bestätigung |
| 06 | Artikel aus Warenkorb entfernen | Fügt 3 hinzu, entfernt 1, bestätigt dass 2 verbleiben |
| 07 | Produktdetailseite | Prüft Name, Preis und Beschreibung eines Produkts |
| 08 | Abmeldung über Burger-Menü | Öffnet Seitennavigation, klickt Logout, bestätigt Login-Seite |
| 09 | Teuerstes Produkt hinzufügen | Sortiert absteigend; fügt Fleece Jacket ($49.99) hinzu |
| 10 | Social-Media-Links im Footer | Prüft sichtbare Twitter-, Facebook- und LinkedIn-Links |

Bildschirme sind als **Page Objects** in `examples/ui-automation/pom/` abstrahiert:

| Page Object | Verfügbare Aktionen |
|---|---|
| `LoginPage` | `Navigate`, `Login`, `IsVisible`, `ErrorMessage` |
| `InventoryPage` | `ProductCount`, `AddToCart`, `CartBadgeCount`, `Sort`, `Logout` |
| `CartPage` | `ItemCount`, `GetItemNames`, `GetItemPrices`, `RemoveItem`, `Checkout` |
| `CheckoutInfoPage` | `FillInfo`, `Continue`, `ErrorMessage` |
| `CheckoutOverviewPage` | `GetItemCount`, `GetSubtotal`, `GetTotal`, `Finish` |
| `CheckoutCompletePage` | `GetHeader`, `BackHome` |
| `ProductDetailPage` | `GetName`, `GetPrice`, `GetDescription`, `AddToCart`, `BackToInventory` |

Ausführliche Szenariobeschreibungen: [`examples/ui-automation/ui-automation.md`](examples/ui-automation/ui-automation.md)

---

### API-Automatisierung — Petstore (`examples/api-automation/`)

30 Testszenarien gegen die öffentliche [Petstore API](https://petstore.swagger.io/v2).
Es wird kein Browser gestartet — nur HTTP-Anfragen und Antwortvalidierungen.

```bash
# Alle Szenarien ausführen
PLAYWRIGHT_CLI_PATH="..." go test -tags e2e -v -timeout 300s ./examples/api-automation/...

# Ein einzelnes Szenario ausführen
PLAYWRIGHT_CLI_PATH="..." go test -tags e2e -v -timeout 60s -run TestScenario01AddPet ./examples/api-automation/...
```

#### Tiere — `/pet`-Endpunkte (01–12)

| # | Szenario | Endpunkt | Prüft |
|---|---|---|---|
| 01 | Tier hinzufügen | `POST /pet` | `id > 0`, Name und Status korrekt |
| 02 | Nach ID abrufen | `GET /pet/{id}` | Felder stimmen mit erstelltem Tier überein |
| 03 | Tier aktualisieren | `PUT /pet` | Status auf `"sold"` aktualisiert |
| 04 | Tier löschen | `DELETE /pet/{id}` | Nachfolgendes GET gibt 404 zurück |
| 05–07 | Nach Status suchen | `GET /pet/findByStatus` | available / pending / sold |
| 08 | Nach mehreren Status suchen | `GET /pet/findByStatus?status=available&status=pending` | Nur diese zwei Status |
| 09 | Mit Formulardaten aktualisieren | `POST /pet/{id}` (form-encoded) | Status 200 |
| 10 | Nach Tag suchen | `GET /pet/findByTags?tags=goaneco` | Gültiges JSON-Array |
| 11 | Nicht vorhandenes Tier | `GET /pet/999999999` | Status 404 |
| 12 | Vollständiger Lebenszyklus | `POST → GET → PUT → DELETE → GET` | Letztes GET gibt 404 |

#### Shop — `/store`-Endpunkte (13–20)

| # | Szenario | Endpunkt | Prüft |
|---|---|---|---|
| 13 | Inventar abrufen | `GET /store/inventory` | Nicht leere Status-Anzahl-Map |
| 14 | Bestellung aufgeben | `POST /store/order` | `id > 0`, `status == "placed"` |
| 15 | Bestellung nach ID abrufen | `GET /store/order/{id}` | Felder stimmen mit Bestellung überein |
| 16 | Bestellung löschen | `DELETE /store/order/{id}` | Nachfolgendes GET gibt 404 |
| 17 | Abgeschlossene Bestellung | `POST /store/order` (`complete: true`) | `complete == true` |
| 18 | Nicht vorhandene Bestellung | `GET /store/order/999999` | Status 404 |
| 19 | Bestellung mit Menge | `POST /store/order` (`quantity: 5`) | `quantity == 5` |
| 20 | Vollständiger Bestelllebenszyklus | `POST → GET → DELETE → GET` | Letztes GET gibt 404 |

#### Benutzer — `/user`-Endpunkte (21–30)

| # | Szenario | Endpunkt | Prüft |
|---|---|---|---|
| 21 | Benutzer erstellen | `POST /user` + `GET /user/{username}` | Felder stimmen überein |
| 22 | Nach Benutzername abrufen | `GET /user/{username}` | `username` stimmt überein |
| 23 | Benutzer aktualisieren | `PUT /user/{username}` + `GET` | Aktualisierte E-Mail zurückgegeben |
| 24 | Benutzer löschen | `DELETE /user/{username}` | Nachfolgendes GET gibt 404 |
| 25 | Login mit gültigen Daten | `GET /user/login` | Status 200, Token im Body |
| 26 | Logout | `GET /user/logout` | Status 200 |
| 27 | Benutzer mit Array erstellen | `POST /user/createWithArray` | Jeder Benutzer per GET abrufbar |
| 28 | Benutzer mit Liste erstellen | `POST /user/createWithList` | Jeder Benutzer per GET abrufbar |
| 29 | Login mit ungültigen Daten | `GET /user/login` (falsche Daten) | Antwort ist nicht 2xx |
| 30 | Vollständiger Benutzerlebenszyklus | `POST → login → GET → PUT → DELETE → GET` | Letztes GET gibt 404 |

Ausführliche Szenariobeschreibungen: [`examples/api-automation/api-automation.md`](examples/api-automation/api-automation.md)

---

## Auf mehreren Browsern ausführen

UI-Tests können auf Chromium (Standard), Firefox oder WebKit ausgeführt werden:

```bash
# Firefox
PLAYWRIGHT_CLI_PATH="..." PLAYWRIGHT_BROWSER=firefox go test -tags e2e -v -timeout 300s ./examples/ui-automation/...

# WebKit (Safari-Engine)
PLAYWRIGHT_CLI_PATH="..." PLAYWRIGHT_BROWSER=webkit go test -tags e2e -v -timeout 300s ./examples/ui-automation/...
```

---

## Testabdeckung

Diese Bibliothek wird durch **3.318 Tests** auf drei Ebenen validiert:

| Ebene | Tests | Was abgedeckt wird |
|---|---|---|
| **Bibliotheks-E2E-Suite** (`e2e/`) | 3.278 | Kern-Playwright-Funktionen: Browser, Page, Locator, Frames, Netzwerk, API-Anfragen, Assertions, Tracing, Screenshots und mehr — 124 Testdateien, die die gesamte API-Oberfläche abdecken |
| **UI-Automatisierungsbeispiele** (`examples/ui-automation/`) | 10 | Vollständige Browser-Flows auf SauceDemo mit dem Page-Object-Model-Muster |
| **API-Automatisierungsbeispiele** (`examples/api-automation/`) | 30 | REST-API-Szenarien gegen die öffentliche Petstore-API — kein Browser erforderlich |

Die Bibliotheks-E2E-Suite wurde unter Verwendung von [playwright-java](https://github.com/microsoft/playwright-java) als Referenz entwickelt, wobei dessen Testfälle nach Go portiert und angepasst wurden, um eine gleichwertige Verhaltensabdeckung sicherzustellen.

---

## Lizenz

MIT-Lizenz — siehe [LICENSE](LICENSE) für Details.

[![Go Reference](https://pkg.go.dev/badge/github.com/jcfrperu/goaneco-playwright.svg)](https://pkg.go.dev/github.com/jcfrperu/goaneco-playwright)
[![Go Report Card](https://goreportcard.com/badge/github.com/jcfrperu/goaneco-playwright)](https://goreportcard.com/report/github.com/jcfrperu/goaneco-playwright)