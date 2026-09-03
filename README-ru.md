# goaneco-playwright

Библиотека Go для автоматизации тестирования с [Microsoft Playwright](https://playwright.dev).
Пишите UI и API тесты на Go с единым инструментом, запускаемым в **Chromium**, **Firefox**
и **WebKit**.

Этот проект является открытым исходным кодом и распространяется под [лицензией MIT](LICENSE).
Он был создан в первую очередь для личного использования, но идеи, предложения и улучшения
приветствуются — вклад, делающий проект более полезным, всегда будет рассмотрен. Вы также
свободны сделать fork и развивать проект в любом направлении.

---

## Что можно тестировать?

| Тип теста | Описание | Включённый пример |
|---|---|---|
| **UI / E2E** | Автоматизирует реальный браузер: клики, формы, навигация, скриншоты | [SauceDemo](examples/ui-automation/) — 10 сценариев |
| **REST API** | Отправляет HTTP-запросы (GET, POST, PUT, DELETE) и проверяет JSON-ответы — без браузера | [Petstore API](examples/api-automation/) — 30 сценариев |

---

## Настройка

### 1. Установка предварительных требований

| Инструмент | Минимальная версия | Как установить |
|---|---|---|
| [Go](https://go.dev/dl/) | 1.21 | Скачать и установить с go.dev |
| [Node.js](https://nodejs.org/) | 18 | Скачать и установить с nodejs.org |
| Playwright | — | `npm install -g playwright` |

### 2. Установка браузеров

```bash
npx playwright install chromium
```

На Linux также установите системные зависимости:

```bash
npx playwright install --with-deps chromium
```

### 3. Установка переменной среды

Все тесты должны знать, где находится Playwright CLI:

```bash
# Linux / macOS
export PLAYWRIGHT_CLI_PATH=$(npm root -g)/playwright/node_modules/playwright-core/cli.js

# Windows PowerShell
$env:PLAYWRIGHT_CLI_PATH = "$env:APPDATA\npm\node_modules\playwright\node_modules\playwright-core\cli.js"
```

Добавьте эту строку в профиль shell (`~/.bashrc`, `~/.zshrc` и т.д.), чтобы не повторять её
в каждой сессии.

---

## Как выглядят тесты

Тесты организованы в виде **сценариев** — один файл на тестовый случай. Каждый тест — это
функция Go, которая чётко описывает, что именно проверяется.

### UI тест (с Page Object Model)

Цель: [https://www.saucedemo.com](https://www.saucedemo.com) — демо-сайт интернет-магазина для практики UI-автоматизации.

```go
func TestLoginAndInventory(t *testing.T) {
    must := require.New(t)
    is   := assert.New(t)
    page := newPage(t)           // открывает новую вкладку браузера
    loginStandardUser(t, page)   // вспомогательная функция для входа со стандартными данными

    ctx := testCtx(t)
    inv := pom.NewInventoryPage(page)

    count, err := inv.ProductCount(ctx)
    must.NoError(err)
    is.Equal(6, count)   // каталог должен показывать ровно 6 товаров
}
```

**Page Object Model (POM)** инкапсулирует действия каждого экрана в переиспользуемую структуру
(`pom.InventoryPage`, `pom.CartPage` и т.д.), делая тесты чистыми и лёгкими для чтения.

### API тест (полный жизненный цикл)

Цель: [https://petstore.swagger.io/v2](https://petstore.swagger.io/v2) — публичный REST API, используемый как песочница для практики API-автоматизации.

```go
func TestScenario12PetFullLifecycle(t *testing.T) {
    must := require.New(t)
    is   := assert.New(t)
    ctx    := testCtx(t)
    apiCtx := newAPICtx(t)   // HTTP-клиент, указывающий на базовый URL API

    // Создать питомца
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

    // Проверить, что питомец можно получить
    getResp, err := apiCtx.Get(ctx, fmt.Sprintf("/pet/%d", created.ID))
    must.NoError(err)
    defer getResp.Dispose(ctx)
    is.Equal(200, getResp.Status())

    // Обновить статус на "sold"
    created.Status = "sold"
    putResp, err := apiCtx.Put(ctx, "/pet", &playwright.APIFetchOptions{Headers: jsonHeaders(), Data: mustJSON(t, created)})
    must.NoError(err)
    defer putResp.Dispose(ctx)
    is.Equal(200, putResp.Status())

    // Удалить питомца
    delResp, err := apiCtx.Delete(ctx, fmt.Sprintf("/pet/%d", created.ID))
    must.NoError(err)
    defer delResp.Dispose(ctx)
    is.Equal(200, delResp.Status())

    // Убедиться, что питомец больше не существует (404)
    confirmResp, err := apiCtx.Get(ctx, fmt.Sprintf("/pet/%d", created.ID))
    must.NoError(err)
    defer confirmResp.Dispose(ctx)
    is.Equal(404, confirmResp.Status())
}
```

---

## Включённые примеры

### UI-автоматизация — SauceDemo (`examples/ui-automation/`)

Автоматизация браузера для [SauceDemo](https://www.saucedemo.com) с использованием паттерна
**Page Object Model**. 10 сценариев охватывают полный поток интернет-магазина.

```bash
# Запустить все сценарии
PLAYWRIGHT_CLI_PATH="..." go test -tags e2e -v -timeout 300s ./examples/ui-automation/...

# Запустить один сценарий
PLAYWRIGHT_CLI_PATH="..." go test -tags e2e -v -timeout 120s -run TestFullCheckoutFlow ./examples/ui-automation/...
```

| # | Сценарий | Что проверяется |
|---|---|---|
| 01 | Вход и каталог | Smoke-тест аутентификации; подтверждает 6 видимых товаров |
| 02 | Добавление товаров в корзину | Добавляет 3 товара; проверяет счётчик значка и список |
| 03 | Сортировка по цене (по возрастанию) | Самый дешёвый товар (Onesie, $7.99) отображается первым |
| 04 | Сортировка по названию (Z–A) | Товары в обратном алфавитном порядке |
| 05 | Полный процесс оформления заказа | Добавить → корзина → данные → обзор → подтверждение |
| 06 | Удаление товара из корзины | Добавляет 3, удаляет 1, подтверждает наличие 2 |
| 07 | Страница деталей товара | Проверяет название, цену и описание |
| 08 | Выход через меню-гамбургер | Открывает боковую панель, нажимает Logout, подтверждает страницу входа |
| 09 | Добавление самого дорогого товара | Сортирует по убыванию; добавляет Fleece Jacket ($49.99) |
| 10 | Ссылки соцсетей в подвале | Проверяет видимые ссылки Twitter, Facebook и LinkedIn |

Экраны абстрагированы как **Page Objects** в `examples/ui-automation/pom/`:

| Page Object | Доступные действия |
|---|---|
| `LoginPage` | `Navigate`, `Login`, `IsVisible`, `ErrorMessage` |
| `InventoryPage` | `ProductCount`, `AddToCart`, `CartBadgeCount`, `Sort`, `Logout` |
| `CartPage` | `ItemCount`, `GetItemNames`, `GetItemPrices`, `RemoveItem`, `Checkout` |
| `CheckoutInfoPage` | `FillInfo`, `Continue`, `ErrorMessage` |
| `CheckoutOverviewPage` | `GetItemCount`, `GetSubtotal`, `GetTotal`, `Finish` |
| `CheckoutCompletePage` | `GetHeader`, `BackHome` |
| `ProductDetailPage` | `GetName`, `GetPrice`, `GetDescription`, `AddToCart`, `BackToInventory` |

Подробное описание сценариев: [`examples/ui-automation/ui-automation.md`](examples/ui-automation/ui-automation.md)

---

### API-автоматизация — Petstore (`examples/api-automation/`)

30 тестовых сценариев против публичного [Petstore API](https://petstore.swagger.io/v2).
Браузер не запускается — только HTTP-запросы и проверка ответов.

```bash
# Запустить все сценарии
PLAYWRIGHT_CLI_PATH="..." go test -tags e2e -v -timeout 300s ./examples/api-automation/...

# Запустить один сценарий
PLAYWRIGHT_CLI_PATH="..." go test -tags e2e -v -timeout 60s -run TestScenario01AddPet ./examples/api-automation/...
```

#### Питомцы — эндпоинты `/pet` (01–12)

| # | Сценарий | Эндпоинт | Проверяет |
|---|---|---|---|
| 01 | Добавить питомца | `POST /pet` | `id > 0`, имя и статус корректны |
| 02 | Получить по ID | `GET /pet/{id}` | Поля совпадают с созданным |
| 03 | Обновить питомца | `PUT /pet` | Статус обновлён до `"sold"` |
| 04 | Удалить питомца | `DELETE /pet/{id}` | Следующий GET возвращает 404 |
| 05–07 | Поиск по статусу | `GET /pet/findByStatus` | available / pending / sold |
| 08 | Поиск по нескольким статусам | `GET /pet/findByStatus?status=available&status=pending` | Только эти два статуса |
| 09 | Обновление через form data | `POST /pet/{id}` (form-encoded) | Статус 200 |
| 10 | Поиск по тегу | `GET /pet/findByTags?tags=goaneco` | Корректный JSON-массив |
| 11 | Несуществующий питомец | `GET /pet/999999999` | Статус 404 |
| 12 | Полный жизненный цикл | `POST → GET → PUT → DELETE → GET` | Последний GET возвращает 404 |

#### Магазин — эндпоинты `/store` (13–20)

| # | Сценарий | Эндпоинт | Проверяет |
|---|---|---|---|
| 13 | Получить инвентарь | `GET /store/inventory` | Непустая карта статус-количество |
| 14 | Сделать заказ | `POST /store/order` | `id > 0`, `status == "placed"` |
| 15 | Получить заказ по ID | `GET /store/order/{id}` | Поля совпадают с заказом |
| 16 | Удалить заказ | `DELETE /store/order/{id}` | Следующий GET возвращает 404 |
| 17 | Выполненный заказ | `POST /store/order` (`complete: true`) | `complete == true` |
| 18 | Несуществующий заказ | `GET /store/order/999999` | Статус 404 |
| 19 | Заказ с количеством | `POST /store/order` (`quantity: 5`) | `quantity == 5` |
| 20 | Полный жизненный цикл заказа | `POST → GET → DELETE → GET` | Последний GET возвращает 404 |

#### Пользователи — эндпоинты `/user` (21–30)

| # | Сценарий | Эндпоинт | Проверяет |
|---|---|---|---|
| 21 | Создать пользователя | `POST /user` + `GET /user/{username}` | Поля совпадают |
| 22 | Получить по имени | `GET /user/{username}` | `username` совпадает |
| 23 | Обновить пользователя | `PUT /user/{username}` + `GET` | Обновлённый email возвращён |
| 24 | Удалить пользователя | `DELETE /user/{username}` | Следующий GET возвращает 404 |
| 25 | Вход с корректными данными | `GET /user/login` | Статус 200, token в body |
| 26 | Выход | `GET /user/logout` | Статус 200 |
| 27 | Создание пользователей массивом | `POST /user/createWithArray` | Каждый пользователь доступен через GET |
| 28 | Создание пользователей списком | `POST /user/createWithList` | Каждый пользователь доступен через GET |
| 29 | Вход с неверными данными | `GET /user/login` (неверные данные) | Ответ не является 2xx |
| 30 | Полный жизненный цикл пользователя | `POST → login → GET → PUT → DELETE → GET` | Последний GET возвращает 404 |

Подробное описание сценариев: [`examples/api-automation/api-automation.md`](examples/api-automation/api-automation.md)

---

## Запуск в нескольких браузерах

UI-тесты можно запускать в Chromium (по умолчанию), Firefox или WebKit:

```bash
# Firefox
PLAYWRIGHT_CLI_PATH="..." PLAYWRIGHT_BROWSER=firefox go test -tags e2e -v -timeout 300s ./examples/ui-automation/...

# WebKit (движок Safari)
PLAYWRIGHT_CLI_PATH="..." PLAYWRIGHT_BROWSER=webkit go test -tags e2e -v -timeout 300s ./examples/ui-automation/...
```

---

## Покрытие тестами

Библиотека проверяется **3 318 тестами**, распределёнными по трём уровням:

| Уровень | Тесты | Что охватывает |
|---|---|---|
| **E2E-набор библиотеки** (`e2e/`) | 3 278 | Основные возможности Playwright: браузер, страница, локатор, фреймы, сеть, API-запросы, утверждения, трассировка, скриншоты и многое другое — 124 тестовых файла, покрывающих всю поверхность API |
| **Примеры UI-автоматизации** (`examples/ui-automation/`) | 10 | Полные браузерные сценарии на SauceDemo с использованием паттерна Page Object Model |
| **Примеры API-автоматизации** (`examples/api-automation/`) | 30 | REST API-сценарии против публичного Petstore API — без запуска браузера |

E2E-набор библиотеки был создан с использованием [playwright-java](https://github.com/microsoft/playwright-java) в качестве справочника — его тестовые случаи были портированы и адаптированы для Go, чтобы обеспечить эквивалентное покрытие поведения.

---

## Лицензия

Лицензия MIT — подробности в [LICENSE](LICENSE).

[![Go Reference](https://pkg.go.dev/badge/github.com/jcfrperu/goaneco-playwright.svg)](https://pkg.go.dev/github.com/jcfrperu/goaneco-playwright)
[![Go Report Card](https://goreportcard.com/badge/github.com/jcfrperu/goaneco-playwright)](https://goreportcard.com/report/github.com/jcfrperu/goaneco-playwright)