# goaneco-playwright

基于 [Microsoft Playwright](https://playwright.dev) 的 Go 测试自动化库。
使用同一工具在 Go 中编写 UI 和 API 测试，支持 **Chromium**、**Firefox** 和 **WebKit**。

本项目为开源项目，基于 [MIT 许可证](LICENSE) 发布。该项目主要为个人使用而创建，但欢迎提出想法、
建议和改进意见——能使项目更有价值的贡献都会被认真考虑。您也可以自由地 fork 本项目并按自己的
方向发展。

---

## 我可以用它测试什么？

| 测试类型 | 说明 | 包含示例 |
|---|---|---|
| **UI / E2E** | 自动化真实浏览器：点击、表单、导航、截图 | [SauceDemo](examples/ui-automation/) — 10 个场景 |
| **REST API** | 发送 HTTP 请求（GET、POST、PUT、DELETE）并验证 JSON 响应——无需浏览器 | [Petstore API](examples/api-automation/) — 30 个场景 |

---

## 环境配置

### 1. 安装前提条件

| 工具 | 最低版本 | 安装方式 |
|---|---|---|
| [Go](https://go.dev/dl/) | 1.21 | 从 go.dev 下载安装 |
| [Node.js](https://nodejs.org/) | 18 | 从 nodejs.org 下载安装 |
| Playwright | — | `npm install -g playwright` |

### 2. 安装浏览器

```bash
npx playwright install chromium
```

在 Linux 上，还需安装操作系统级依赖：

```bash
npx playwright install --with-deps chromium
```

### 3. 设置环境变量

所有测试都需要知道 Playwright CLI 的位置：

```bash
# Linux / macOS
export PLAYWRIGHT_CLI_PATH=$(npm root -g)/playwright/node_modules/playwright-core/cli.js

# Windows PowerShell
$env:PLAYWRIGHT_CLI_PATH = "$env:APPDATA\npm\node_modules\playwright\node_modules\playwright-core\cli.js"
```

将此行添加到 shell 配置文件（`~/.bashrc`、`~/.zshrc` 等）中，避免每次会话重复设置。

---

## 测试代码示例

测试以**场景**为单位组织——每个测试用例对应一个文件。每个测试都是一个 Go 函数，
清晰地描述其验证内容。

### UI 测试（使用 Page Object Model）

目标：[https://www.saucedemo.com](https://www.saucedemo.com) — 用于 UI 自动化练习的演示电商网站。

```go
func TestLoginAndInventory(t *testing.T) {
    must := require.New(t)
    is   := assert.New(t)
    page := newPage(t)           // 打开新的浏览器标签页
    loginStandardUser(t, page)   // 使用标准凭据登录的辅助函数

    ctx := testCtx(t)
    inv := pom.NewInventoryPage(page)

    count, err := inv.ProductCount(ctx)
    must.NoError(err)
    is.Equal(6, count)   // 商品列表必须显示恰好 6 个产品
}
```

**Page Object Model（POM）** 将每个页面的操作封装到可复用的结构体中
（`pom.InventoryPage`、`pom.CartPage` 等），使测试代码简洁易读。

### API 测试（完整生命周期）

目标：[https://petstore.swagger.io/v2](https://petstore.swagger.io/v2) — 用于 API 自动化练习的公共 REST API 沙箱。

```go
func TestScenario12PetFullLifecycle(t *testing.T) {
    must := require.New(t)
    is   := assert.New(t)
    ctx    := testCtx(t)
    apiCtx := newAPICtx(t)   // 指向 API 基础 URL 的 HTTP 客户端

    // 创建宠物
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

    // 验证可以检索到
    getResp, err := apiCtx.Get(ctx, fmt.Sprintf("/pet/%d", created.ID))
    must.NoError(err)
    defer getResp.Dispose(ctx)
    is.Equal(200, getResp.Status())

    // 将状态更新为 "sold"
    created.Status = "sold"
    putResp, err := apiCtx.Put(ctx, "/pet", &playwright.APIFetchOptions{Headers: jsonHeaders(), Data: mustJSON(t, created)})
    must.NoError(err)
    defer putResp.Dispose(ctx)
    is.Equal(200, putResp.Status())

    // 删除宠物
    delResp, err := apiCtx.Delete(ctx, fmt.Sprintf("/pet/%d", created.ID))
    must.NoError(err)
    defer delResp.Dispose(ctx)
    is.Equal(200, delResp.Status())

    // 确认已不存在（404）
    confirmResp, err := apiCtx.Get(ctx, fmt.Sprintf("/pet/%d", created.ID))
    must.NoError(err)
    defer confirmResp.Dispose(ctx)
    is.Equal(404, confirmResp.Status())
}
```

---

## 包含的示例

### UI 自动化 — SauceDemo（`examples/ui-automation/`）

使用 **Page Object Model** 模式对 [SauceDemo](https://www.saucedemo.com) 进行浏览器自动化。
10 个场景覆盖了在线商店的完整流程。

```bash
# 运行所有场景
PLAYWRIGHT_CLI_PATH="..." go test -tags e2e -v -timeout 300s ./examples/ui-automation/...

# 运行单个场景
PLAYWRIGHT_CLI_PATH="..." go test -tags e2e -v -timeout 120s -run TestFullCheckoutFlow ./examples/ui-automation/...
```

| # | 场景 | 验证内容 |
|---|---|---|
| 01 | 登录和商品列表 | 身份验证冒烟测试；确认 6 个产品可见 |
| 02 | 添加商品到购物车 | 添加 3 个产品；验证购物车徽章数量和商品列表 |
| 03 | 按价格从低到高排序 | 最便宜的商品（Onesie，$7.99）排在第一位 |
| 04 | 按名称从 Z 到 A 排序 | 产品按字母逆序排列 |
| 05 | 完整结账流程 | 添加 → 购物车 → 填写信息 → 概览 → 确认订单 |
| 06 | 从购物车删除商品 | 添加 3 个，删除 1 个，确认剩余 2 个 |
| 07 | 产品详情页 | 验证产品名称、价格和描述 |
| 08 | 通过汉堡菜单退出登录 | 打开侧边导航，点击登出，确认显示登录页面 |
| 09 | 添加最贵的商品 | 按价格从高到低排序；添加 Fleece Jacket（$49.99） |
| 10 | 页脚社交媒体链接 | 验证 Twitter、Facebook 和 LinkedIn 链接可见 |

页面被抽象为 `examples/ui-automation/pom/` 中的 **Page Objects**：

| Page Object | 可用操作 |
|---|---|
| `LoginPage` | `Navigate`, `Login`, `IsVisible`, `ErrorMessage` |
| `InventoryPage` | `ProductCount`, `AddToCart`, `CartBadgeCount`, `Sort`, `Logout` |
| `CartPage` | `ItemCount`, `GetItemNames`, `GetItemPrices`, `RemoveItem`, `Checkout` |
| `CheckoutInfoPage` | `FillInfo`, `Continue`, `ErrorMessage` |
| `CheckoutOverviewPage` | `GetItemCount`, `GetSubtotal`, `GetTotal`, `Finish` |
| `CheckoutCompletePage` | `GetHeader`, `BackHome` |
| `ProductDetailPage` | `GetName`, `GetPrice`, `GetDescription`, `AddToCart`, `BackToInventory` |

完整场景说明：[`examples/ui-automation/ui-automation.md`](examples/ui-automation/ui-automation.md)

---

### API 自动化 — Petstore（`examples/api-automation/`）

针对公共 [Petstore API](https://petstore.swagger.io/v2) 的 30 个测试场景。
无需启动浏览器——仅发送 HTTP 请求并验证响应。

```bash
# 运行所有场景
PLAYWRIGHT_CLI_PATH="..." go test -tags e2e -v -timeout 300s ./examples/api-automation/...

# 运行单个场景
PLAYWRIGHT_CLI_PATH="..." go test -tags e2e -v -timeout 60s -run TestScenario01AddPet ./examples/api-automation/...
```

#### 宠物 — `/pet` 接口（01–12）

| # | 场景 | 接口 | 验证 |
|---|---|---|---|
| 01 | 添加宠物 | `POST /pet` | `id > 0`，名称和状态正确 |
| 02 | 按 ID 获取 | `GET /pet/{id}` | 返回字段与创建的宠物一致 |
| 03 | 更新宠物 | `PUT /pet` | 状态更新为 `"sold"` |
| 04 | 删除宠物 | `DELETE /pet/{id}` | 后续 GET 返回 404 |
| 05–07 | 按状态查询 | `GET /pet/findByStatus` | available / pending / sold |
| 08 | 按多个状态查询 | `GET /pet/findByStatus?status=available&status=pending` | 仅返回这两种状态 |
| 09 | 用表单数据更新 | `POST /pet/{id}`（form-encoded） | 状态 200 |
| 10 | 按标签查询 | `GET /pet/findByTags?tags=goaneco` | 有效的 JSON 数组 |
| 11 | 不存在的宠物 | `GET /pet/999999999` | 状态 404 |
| 12 | 完整生命周期 | `POST → GET → PUT → DELETE → GET` | 最终 GET 返回 404 |

#### 商店 — `/store` 接口（13–20）

| # | 场景 | 接口 | 验证 |
|---|---|---|---|
| 13 | 获取库存 | `GET /store/inventory` | 非空的状态-数量映射 |
| 14 | 下订单 | `POST /store/order` | `id > 0`，`status == "placed"` |
| 15 | 按 ID 获取订单 | `GET /store/order/{id}` | 字段与订单一致 |
| 16 | 删除订单 | `DELETE /store/order/{id}` | 后续 GET 返回 404 |
| 17 | 下已完成的订单 | `POST /store/order`（`complete: true`） | `complete == true` |
| 18 | 不存在的订单 | `GET /store/order/999999` | 状态 404 |
| 19 | 带数量的订单 | `POST /store/order`（`quantity: 5`） | `quantity == 5` |
| 20 | 完整订单生命周期 | `POST → GET → DELETE → GET` | 最终 GET 返回 404 |

#### 用户 — `/user` 接口（21–30）

| # | 场景 | 接口 | 验证 |
|---|---|---|---|
| 21 | 创建用户 | `POST /user` + `GET /user/{username}` | 字段一致 |
| 22 | 按用户名获取 | `GET /user/{username}` | `username` 匹配 |
| 23 | 更新用户 | `PUT /user/{username}` + `GET` | 返回更新后的邮箱 |
| 24 | 删除用户 | `DELETE /user/{username}` | 后续 GET 返回 404 |
| 25 | 有效凭据登录 | `GET /user/login` | 状态 200，body 中含 token |
| 26 | 登出 | `GET /user/logout` | 状态 200 |
| 27 | 用数组批量创建用户 | `POST /user/createWithArray` | 每个用户均可通过 GET 获取 |
| 28 | 用列表批量创建用户 | `POST /user/createWithList` | 每个用户均可通过 GET 获取 |
| 29 | 无效凭据登录 | `GET /user/login`（错误凭据） | 响应不是 2xx |
| 30 | 完整用户生命周期 | `POST → login → GET → PUT → DELETE → GET` | 最终 GET 返回 404 |

完整场景说明：[`examples/api-automation/api-automation.md`](examples/api-automation/api-automation.md)

---

## 在多个浏览器上运行

UI 测试可在 Chromium（默认）、Firefox 或 WebKit 上运行：

```bash
# Firefox
PLAYWRIGHT_CLI_PATH="..." PLAYWRIGHT_BROWSER=firefox go test -tags e2e -v -timeout 300s ./examples/ui-automation/...

# WebKit（Safari 引擎）
PLAYWRIGHT_CLI_PATH="..." PLAYWRIGHT_BROWSER=webkit go test -tags e2e -v -timeout 300s ./examples/ui-automation/...
```

---

## 测试覆盖

本库通过分布在三个层次的 **3,318 个测试**进行验证：

| 层次 | 测试数 | 覆盖内容 |
|---|---|---|
| **库 E2E 测试套件** (`e2e/`) | 3,278 | Playwright 核心能力：浏览器、页面、定位器、框架、网络、API 请求、断言、追踪、截图等——124 个测试文件，覆盖完整 API 范围 |
| **UI 自动化示例** (`examples/ui-automation/`) | 10 | 使用 Page Object Model 模式对 SauceDemo 进行端到端浏览器流程测试 |
| **API 自动化示例** (`examples/api-automation/`) | 30 | 针对公共 Petstore API 的 REST API 场景测试——无需浏览器 |

库 E2E 测试套件以 [playwright-java](https://github.com/microsoft/playwright-java) 为参考构建，将其测试用例移植并适配到 Go，以确保等效的行为覆盖。

---

## 许可证

MIT 许可证——详见 [LICENSE](LICENSE)。

[![Go Reference](https://pkg.go.dev/badge/github.com/jcfrperu/goaneco-playwright.svg)](https://pkg.go.dev/github.com/jcfrperu/goaneco-playwright)
[![Go Report Card](https://goreportcard.com/badge/github.com/jcfrperu/goaneco-playwright)](https://goreportcard.com/report/github.com/jcfrperu/goaneco-playwright)