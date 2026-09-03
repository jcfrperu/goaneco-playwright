# goaneco-playwright

Biblioteca Go para automação de testes com [Microsoft Playwright](https://playwright.dev).
Escreva testes de UI e API em Go usando a mesma ferramenta, executando em **Chromium**,
**Firefox** e **WebKit**.

Este projeto é open source, licenciado sob a [Licença MIT](LICENSE). Foi criado principalmente
para uso pessoal, mas ideias, sugestões e melhorias são bem-vindas — contribuições que tornem
o projeto mais útil serão sempre consideradas. Você também é livre para fazer um fork e
desenvolver na direção que preferir.

---

## O que posso testar com isso?

| Tipo de teste | Descrição | Exemplo incluído |
|---|---|---|
| **UI / E2E** | Automatiza um navegador real: cliques, formulários, navegação, capturas de tela | [SauceDemo](examples/ui-automation/) — 10 cenários |
| **REST API** | Envia requisições HTTP (GET, POST, PUT, DELETE) e valida respostas JSON — sem navegador | [Petstore API](examples/api-automation/) — 30 cenários |

---

## Configuração

### 1. Instalar os pré-requisitos

| Ferramenta | Versão mínima | Como instalar |
|---|---|---|
| [Go](https://go.dev/dl/) | 1.21 | Baixar e instalar em go.dev |
| [Node.js](https://nodejs.org/) | 18 | Baixar e instalar em nodejs.org |
| Playwright | — | `npm install -g playwright` |

### 2. Instalar os navegadores

```bash
npx playwright install chromium
```

No Linux, instale também as dependências do sistema operacional:

```bash
npx playwright install --with-deps chromium
```

### 3. Configurar a variável de ambiente

Todos os testes precisam saber onde está o CLI do Playwright:

```bash
# Linux / macOS
export PLAYWRIGHT_CLI_PATH=$(npm root -g)/playwright/node_modules/playwright-core/cli.js

# Windows PowerShell
$env:PLAYWRIGHT_CLI_PATH = "$env:APPDATA\npm\node_modules\playwright\node_modules\playwright-core\cli.js"
```

Adicione essa linha ao seu perfil de shell (`~/.bashrc`, `~/.zshrc`, etc.) para não precisar
repeti-la a cada sessão.

---

## Como os testes se parecem

Os testes são organizados como **cenários** — um arquivo por caso de teste. Cada teste é uma
função Go que descreve claramente o que está verificando.

### Teste de UI (com Page Object Model)

Alvo: [https://www.saucedemo.com](https://www.saucedemo.com) — um site de e-commerce de demonstração para prática de automação de UI.

```go
func TestLoginAndInventory(t *testing.T) {
    must := require.New(t)
    is   := assert.New(t)
    page := newPage(t)           // abre uma nova aba no navegador
    loginStandardUser(t, page)   // helper que faz login com credenciais padrão

    ctx := testCtx(t)
    inv := pom.NewInventoryPage(page)

    count, err := inv.ProductCount(ctx)
    must.NoError(err)
    is.Equal(6, count)   // o inventário deve exibir exatamente 6 produtos
}
```

O **Page Object Model (POM)** encapsula as ações de cada tela em uma struct reutilizável
(`pom.InventoryPage`, `pom.CartPage`, etc.), mantendo os testes limpos e fáceis de ler.

### Teste de API (ciclo de vida completo)

Alvo: [https://petstore.swagger.io/v2](https://petstore.swagger.io/v2) — uma REST API pública usada como sandbox para prática de automação de API.

```go
func TestScenario12PetFullLifecycle(t *testing.T) {
    must := require.New(t)
    is   := assert.New(t)
    ctx    := testCtx(t)
    apiCtx := newAPICtx(t)   // cliente HTTP apontando para a URL base da API

    // Criar um pet
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

    // Verificar que pode ser recuperado
    getResp, err := apiCtx.Get(ctx, fmt.Sprintf("/pet/%d", created.ID))
    must.NoError(err)
    defer getResp.Dispose(ctx)
    is.Equal(200, getResp.Status())

    // Atualizar o status para "sold"
    created.Status = "sold"
    putResp, err := apiCtx.Put(ctx, "/pet", &playwright.APIFetchOptions{Headers: jsonHeaders(), Data: mustJSON(t, created)})
    must.NoError(err)
    defer putResp.Dispose(ctx)
    is.Equal(200, putResp.Status())

    // Deletar o pet
    delResp, err := apiCtx.Delete(ctx, fmt.Sprintf("/pet/%d", created.ID))
    must.NoError(err)
    defer delResp.Dispose(ctx)
    is.Equal(200, delResp.Status())

    // Confirmar que não existe mais (404)
    confirmResp, err := apiCtx.Get(ctx, fmt.Sprintf("/pet/%d", created.ID))
    must.NoError(err)
    defer confirmResp.Dispose(ctx)
    is.Equal(404, confirmResp.Status())
}
```

---

## Exemplos incluídos

### Automação de UI — SauceDemo (`examples/ui-automation/`)

Automação de navegador para [SauceDemo](https://www.saucedemo.com) usando o padrão
**Page Object Model**. Os 10 cenários cobrem o fluxo completo de uma loja online.

```bash
# Executar todos os cenários
PLAYWRIGHT_CLI_PATH="..." go test -tags e2e -v -timeout 300s ./examples/ui-automation/...

# Executar um cenário específico
PLAYWRIGHT_CLI_PATH="..." go test -tags e2e -v -timeout 120s -run TestFullCheckoutFlow ./examples/ui-automation/...
```

| # | Cenário | O que verifica |
|---|---|---|
| 01 | Login e inventário | Smoke test de autenticação; confirma 6 produtos visíveis |
| 02 | Adicionar itens ao carrinho | Adiciona 3 produtos; verifica badge e lista do carrinho |
| 03 | Ordenar por preço (menor ao maior) | O item mais barato (Onesie, $7.99) aparece primeiro |
| 04 | Ordenar por nome (Z a A) | Produtos em ordem alfabética reversa |
| 05 | Fluxo completo de checkout | Adicionar → carrinho → preencher dados → resumo → confirmação |
| 06 | Remover item do carrinho | Adiciona 3, remove 1, confirma que restam 2 |
| 07 | Página de detalhe do produto | Verifica nome, preço e descrição |
| 08 | Logout pelo menu hambúrguer | Abre nav lateral, clica em Logout, confirma tela de login |
| 09 | Adicionar o item mais caro | Ordena do maior ao menor; adiciona Fleece Jacket ($49.99) |
| 10 | Links de redes sociais no rodapé | Verifica links visíveis do Twitter, Facebook e LinkedIn |

As telas são abstraídas como **Page Objects** em `examples/ui-automation/pom/`:

| Page Object | Ações disponíveis |
|---|---|
| `LoginPage` | `Navigate`, `Login`, `IsVisible`, `ErrorMessage` |
| `InventoryPage` | `ProductCount`, `AddToCart`, `CartBadgeCount`, `Sort`, `Logout` |
| `CartPage` | `ItemCount`, `GetItemNames`, `GetItemPrices`, `RemoveItem`, `Checkout` |
| `CheckoutInfoPage` | `FillInfo`, `Continue`, `ErrorMessage` |
| `CheckoutOverviewPage` | `GetItemCount`, `GetSubtotal`, `GetTotal`, `Finish` |
| `CheckoutCompletePage` | `GetHeader`, `BackHome` |
| `ProductDetailPage` | `GetName`, `GetPrice`, `GetDescription`, `AddToCart`, `BackToInventory` |

Descrições detalhadas dos cenários: [`examples/ui-automation/ui-automation.md`](examples/ui-automation/ui-automation.md)

---

### Automação de API — Petstore (`examples/api-automation/`)

30 cenários de teste contra a [Petstore API](https://petstore.swagger.io/v2) pública.
Nenhum navegador é iniciado — apenas requisições HTTP e validações de resposta.

```bash
# Executar todos os cenários
PLAYWRIGHT_CLI_PATH="..." go test -tags e2e -v -timeout 300s ./examples/api-automation/...

# Executar um cenário específico
PLAYWRIGHT_CLI_PATH="..." go test -tags e2e -v -timeout 60s -run TestScenario01AddPet ./examples/api-automation/...
```

#### Pets — endpoints `/pet` (01–12)

| # | Cenário | Endpoint | Verifica |
|---|---|---|---|
| 01 | Adicionar pet | `POST /pet` | `id > 0`, nome e status corretos |
| 02 | Obter por ID | `GET /pet/{id}` | Campos correspondem ao pet criado |
| 03 | Atualizar pet | `PUT /pet` | Status atualizado para `"sold"` |
| 04 | Deletar pet | `DELETE /pet/{id}` | GET subsequente retorna 404 |
| 05–07 | Buscar por status | `GET /pet/findByStatus` | available / pending / sold |
| 08 | Buscar por múltiplos status | `GET /pet/findByStatus?status=available&status=pending` | Apenas esses dois status |
| 09 | Atualizar com form data | `POST /pet/{id}` (form-encoded) | Status 200 |
| 10 | Buscar por tag | `GET /pet/findByTags?tags=goaneco` | Array JSON válido |
| 11 | Pet inexistente | `GET /pet/999999999` | Status 404 |
| 12 | Ciclo de vida completo | `POST → GET → PUT → DELETE → GET` | GET final retorna 404 |

#### Loja — endpoints `/store` (13–20)

| # | Cenário | Endpoint | Verifica |
|---|---|---|---|
| 13 | Obter inventário | `GET /store/inventory` | Mapa status-contagem não vazio |
| 14 | Fazer pedido | `POST /store/order` | `id > 0`, `status == "placed"` |
| 15 | Obter pedido por ID | `GET /store/order/{id}` | Campos correspondem ao pedido |
| 16 | Deletar pedido | `DELETE /store/order/{id}` | GET subsequente retorna 404 |
| 17 | Pedido completo | `POST /store/order` (`complete: true`) | `complete == true` |
| 18 | Pedido inexistente | `GET /store/order/999999` | Status 404 |
| 19 | Pedido com quantidade | `POST /store/order` (`quantity: 5`) | `quantity == 5` |
| 20 | Ciclo de vida completo do pedido | `POST → GET → DELETE → GET` | GET final retorna 404 |

#### Usuários — endpoints `/user` (21–30)

| # | Cenário | Endpoint | Verifica |
|---|---|---|---|
| 21 | Criar usuário | `POST /user` + `GET /user/{username}` | Campos correspondem |
| 22 | Obter por username | `GET /user/{username}` | `username` corresponde |
| 23 | Atualizar usuário | `PUT /user/{username}` + `GET` | E-mail atualizado retornado |
| 24 | Deletar usuário | `DELETE /user/{username}` | GET subsequente retorna 404 |
| 25 | Login com credenciais válidas | `GET /user/login` | Status 200, token no body |
| 26 | Logout | `GET /user/logout` | Status 200 |
| 27 | Criar usuários com array | `POST /user/createWithArray` | Cada usuário recuperável via GET |
| 28 | Criar usuários com lista | `POST /user/createWithList` | Cada usuário recuperável via GET |
| 29 | Login com credenciais inválidas | `GET /user/login` (credenciais erradas) | Resposta não é 2xx |
| 30 | Ciclo de vida completo do usuário | `POST → login → GET → PUT → DELETE → GET` | GET final retorna 404 |

Descrições detalhadas dos cenários: [`examples/api-automation/api-automation.md`](examples/api-automation/api-automation.md)

---

## Executar em múltiplos navegadores

Os testes de UI podem rodar no Chromium (padrão), Firefox ou WebKit:

```bash
# Firefox
PLAYWRIGHT_CLI_PATH="..." PLAYWRIGHT_BROWSER=firefox go test -tags e2e -v -timeout 300s ./examples/ui-automation/...

# WebKit (motor do Safari)
PLAYWRIGHT_CLI_PATH="..." PLAYWRIGHT_BROWSER=webkit go test -tags e2e -v -timeout 300s ./examples/ui-automation/...
```

---

## Cobertura de testes

Esta biblioteca é validada por **3.318 testes** distribuídos em três camadas:

| Camada | Testes | O que cobre |
|---|---|---|
| **Suite E2E da biblioteca** (`e2e/`) | 3.278 | Capacidades core do Playwright: browser, page, locator, frames, rede, requisições API, assertions, tracing, screenshots e mais — 124 arquivos de teste cobrindo toda a superfície da API |
| **Exemplos de automação UI** (`examples/ui-automation/`) | 10 | Fluxos completos de navegador no SauceDemo usando o padrão Page Object Model |
| **Exemplos de automação API** (`examples/api-automation/`) | 30 | Cenários REST API contra a Petstore API pública — sem necessidade de navegador |

A suite E2E da biblioteca foi construída usando [playwright-java](https://github.com/microsoft/playwright-java) como referência, portando e adaptando seus casos de teste para Go para garantir uma cobertura de comportamento equivalente.

---

## Licença

Licença MIT — veja [LICENSE](LICENSE) para detalhes.

[![Go Reference](https://pkg.go.dev/badge/github.com/jcfrperu/goaneco-playwright.svg)](https://pkg.go.dev/github.com/jcfrperu/goaneco-playwright)
[![Go Report Card](https://goreportcard.com/badge/github.com/jcfrperu/goaneco-playwright)](https://goreportcard.com/report/github.com/jcfrperu/goaneco-playwright)