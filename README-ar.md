# goaneco-playwright

مكتبة Go لأتمتة الاختبارات باستخدام [Microsoft Playwright](https://playwright.dev).
اكتب اختبارات UI وAPI بلغة Go بأداة واحدة، تعمل على **Chromium** و**Firefox** و**WebKit**.

هذا المشروع مفتوح المصدر، مُرخَّص بموجب [رخصة MIT](LICENSE). تم إنشاؤه بشكل رئيسي للاستخدام
الشخصي، لكن الأفكار والاقتراحات والتحسينات مرحَّب بها — أي مساهمة تجعل المشروع أكثر فائدة
ستُؤخَذ دائمًا بعين الاعتبار. كما أنك حر في نسخ المشروع (fork) وتطويره في الاتجاه الذي تريده.

---

## ماذا يمكنني اختباره بهذه المكتبة؟

| نوع الاختبار | الوصف | مثال مُضمَّن |
|---|---|---|
| **UI / E2E** | أتمتة متصفح حقيقي: نقرات، نماذج، تنقل، لقطات شاشة | [SauceDemo](examples/ui-automation/) — 10 سيناريوهات |
| **REST API** | إرسال طلبات HTTP (GET, POST, PUT, DELETE) والتحقق من استجابات JSON — بدون متصفح | [Petstore API](examples/api-automation/) — 30 سيناريو |

---

## الإعداد

### 1. تثبيت المتطلبات الأساسية

| الأداة | الإصدار الأدنى | طريقة التثبيت |
|---|---|---|
| [Go](https://go.dev/dl/) | 1.21 | تنزيل وتثبيت من go.dev |
| [Node.js](https://nodejs.org/) | 18 | تنزيل وتثبيت من nodejs.org |
| Playwright | — | `npm install -g playwright` |

### 2. تثبيت المتصفحات

```bash
npx playwright install chromium
```

على Linux، ثبِّت أيضًا التبعيات على مستوى نظام التشغيل:

```bash
npx playwright install --with-deps chromium
```

### 3. تعيين متغير البيئة

تحتاج جميع الاختبارات إلى معرفة مكان Playwright CLI:

```bash
# Linux / macOS
export PLAYWRIGHT_CLI_PATH=$(npm root -g)/playwright/node_modules/playwright-core/cli.js

# Windows PowerShell
$env:PLAYWRIGHT_CLI_PATH = "$env:APPDATA\npm\node_modules\playwright\node_modules\playwright-core\cli.js"
```

أضف هذا السطر إلى ملف إعداد الـ shell (`~/.bashrc`، `~/.zshrc`، إلخ) لتجنب تكراره في كل جلسة.

---

## كيف تبدو الاختبارات؟

تُنظَّم الاختبارات كـ**سيناريوهات** — ملف واحد لكل حالة اختبار. كل اختبار هو دالة Go تصف
بوضوح ما تتحقق منه.

### اختبار UI (مع نموذج كائن الصفحة)

الهدف: [https://www.saucedemo.com](https://www.saucedemo.com) — موقع تجارة إلكترونية تجريبي للتدرب على أتمتة UI.

```go
func TestLoginAndInventory(t *testing.T) {
    must := require.New(t)
    is   := assert.New(t)
    page := newPage(t)           // يفتح تبويبًا جديدًا في المتصفح
    loginStandardUser(t, page)   // دالة مساعدة تسجّل الدخول ببيانات اعتماد قياسية

    ctx := testCtx(t)
    inv := pom.NewInventoryPage(page)

    count, err := inv.ProductCount(ctx)
    must.NoError(err)
    is.Equal(6, count)   // يجب أن يعرض المخزون 6 منتجات بالضبط
}
```

**نموذج كائن الصفحة (POM)** يغلِّف إجراءات كل شاشة في struct قابل لإعادة الاستخدام
(`pom.InventoryPage`، `pom.CartPage`، إلخ)، مما يجعل الاختبارات نظيفة وسهلة القراءة.

### اختبار API (دورة حياة كاملة)

الهدف: [https://petstore.swagger.io/v2](https://petstore.swagger.io/v2) — REST API عامة تُستخدم كبيئة تجريبية لأتمتة الـ API.

```go
func TestScenario12PetFullLifecycle(t *testing.T) {
    must := require.New(t)
    is   := assert.New(t)
    ctx    := testCtx(t)
    apiCtx := newAPICtx(t)   // عميل HTTP يشير إلى عنوان URL الأساسي للـ API

    // إنشاء حيوان أليف
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

    // التحقق من إمكانية استرجاعه
    getResp, err := apiCtx.Get(ctx, fmt.Sprintf("/pet/%d", created.ID))
    must.NoError(err)
    defer getResp.Dispose(ctx)
    is.Equal(200, getResp.Status())

    // تحديث الحالة إلى "sold"
    created.Status = "sold"
    putResp, err := apiCtx.Put(ctx, "/pet", &playwright.APIFetchOptions{Headers: jsonHeaders(), Data: mustJSON(t, created)})
    must.NoError(err)
    defer putResp.Dispose(ctx)
    is.Equal(200, putResp.Status())

    // حذف الحيوان الأليف
    delResp, err := apiCtx.Delete(ctx, fmt.Sprintf("/pet/%d", created.ID))
    must.NoError(err)
    defer delResp.Dispose(ctx)
    is.Equal(200, delResp.Status())

    // التأكد من عدم وجوده بعد الآن (404)
    confirmResp, err := apiCtx.Get(ctx, fmt.Sprintf("/pet/%d", created.ID))
    must.NoError(err)
    defer confirmResp.Dispose(ctx)
    is.Equal(404, confirmResp.Status())
}
```

---

## الأمثلة المُضمَّنة

### أتمتة UI — SauceDemo (`examples/ui-automation/`)

أتمتة المتصفح لـ [SauceDemo](https://www.saucedemo.com) باستخدام نمط **Page Object Model**.
تغطي السيناريوهات العشرة التدفق الكامل لمتجر إلكتروني.

```bash
# تشغيل جميع السيناريوهات
PLAYWRIGHT_CLI_PATH="..." go test -tags e2e -v -timeout 300s ./examples/ui-automation/...

# تشغيل سيناريو واحد
PLAYWRIGHT_CLI_PATH="..." go test -tags e2e -v -timeout 120s -run TestFullCheckoutFlow ./examples/ui-automation/...
```

| # | السيناريو | ما يتحقق منه |
|---|---|---|
| 01 | تسجيل الدخول والمخزون | اختبار دخاني للمصادقة؛ يؤكد ظهور 6 منتجات |
| 02 | إضافة عناصر إلى السلة | يضيف 3 منتجات؛ يتحقق من عدد شارة السلة وقائمة العناصر |
| 03 | الترتيب بالسعر من الأقل للأعلى | أرخص منتج (Onesie، $7.99) يظهر أولًا |
| 04 | الترتيب بالاسم من Z إلى A | المنتجات مُرتَّبة أبجديًا عكسيًا |
| 05 | تدفق الدفع الكامل | إضافة → سلة → تعبئة البيانات → مراجعة → تأكيد |
| 06 | إزالة عنصر من السلة | يضيف 3، يزيل 1، يؤكد بقاء 2 |
| 07 | صفحة تفاصيل المنتج | يتحقق من الاسم والسعر والوصف |
| 08 | تسجيل الخروج عبر قائمة البرغر | يفتح القائمة الجانبية، ينقر خروج، يؤكد شاشة الدخول |
| 09 | إضافة أغلى منتج | يرتب من الأعلى للأقل؛ يضيف Fleece Jacket ($49.99) |
| 10 | روابط وسائل التواصل في التذييل | يتحقق من ظهور روابط Twitter وFacebook وLinkedIn |

الشاشات مجرَّدة كـ**Page Objects** في `examples/ui-automation/pom/`:

| Page Object | الإجراءات المتاحة |
|---|---|
| `LoginPage` | `Navigate`, `Login`, `IsVisible`, `ErrorMessage` |
| `InventoryPage` | `ProductCount`, `AddToCart`, `CartBadgeCount`, `Sort`, `Logout` |
| `CartPage` | `ItemCount`, `GetItemNames`, `GetItemPrices`, `RemoveItem`, `Checkout` |
| `CheckoutInfoPage` | `FillInfo`, `Continue`, `ErrorMessage` |
| `CheckoutOverviewPage` | `GetItemCount`, `GetSubtotal`, `GetTotal`, `Finish` |
| `CheckoutCompletePage` | `GetHeader`, `BackHome` |
| `ProductDetailPage` | `GetName`, `GetPrice`, `GetDescription`, `AddToCart`, `BackToInventory` |

أوصاف السيناريوهات التفصيلية: [`examples/ui-automation/ui-automation.md`](examples/ui-automation/ui-automation.md)

---

### أتمتة API — Petstore (`examples/api-automation/`)

30 سيناريو اختبار على [Petstore API](https://petstore.swagger.io/v2) العامة.
لا يُشغَّل أي متصفح — طلبات HTTP والتحقق من الاستجابات فقط.

```bash
# تشغيل جميع السيناريوهات
PLAYWRIGHT_CLI_PATH="..." go test -tags e2e -v -timeout 300s ./examples/api-automation/...

# تشغيل سيناريو واحد
PLAYWRIGHT_CLI_PATH="..." go test -tags e2e -v -timeout 60s -run TestScenario01AddPet ./examples/api-automation/...
```

#### الحيوانات الأليفة — نقاط نهاية `/pet` (01–12)

| # | السيناريو | نقطة النهاية | يتحقق من |
|---|---|---|---|
| 01 | إضافة حيوان أليف | `POST /pet` | `id > 0`، الاسم والحالة صحيحان |
| 02 | الحصول بالمعرّف | `GET /pet/{id}` | الحقول تتطابق مع المُنشأ |
| 03 | تحديث حيوان أليف | `PUT /pet` | الحالة محدَّثة إلى `"sold"` |
| 04 | حذف حيوان أليف | `DELETE /pet/{id}` | GET التالي يُعيد 404 |
| 05–07 | البحث بالحالة | `GET /pet/findByStatus` | available / pending / sold |
| 08 | البحث بحالات متعددة | `GET /pet/findByStatus?status=available&status=pending` | هاتان الحالتان فقط |
| 09 | التحديث ببيانات النموذج | `POST /pet/{id}` (form-encoded) | الحالة 200 |
| 10 | البحث بالوسم | `GET /pet/findByTags?tags=goaneco` | مصفوفة JSON صالحة |
| 11 | حيوان أليف غير موجود | `GET /pet/999999999` | الحالة 404 |
| 12 | دورة الحياة الكاملة | `POST → GET → PUT → DELETE → GET` | GET الأخير يُعيد 404 |

#### المتجر — نقاط نهاية `/store` (13–20)

| # | السيناريو | نقطة النهاية | يتحقق من |
|---|---|---|---|
| 13 | الحصول على المخزون | `GET /store/inventory` | خريطة حالة-عدد غير فارغة |
| 14 | تقديم طلب | `POST /store/order` | `id > 0`، `status == "placed"` |
| 15 | الحصول على طلب بالمعرّف | `GET /store/order/{id}` | الحقول تتطابق مع الطلب |
| 16 | حذف طلب | `DELETE /store/order/{id}` | GET التالي يُعيد 404 |
| 17 | طلب مكتمل | `POST /store/order` (`complete: true`) | `complete == true` |
| 18 | طلب غير موجود | `GET /store/order/999999` | الحالة 404 |
| 19 | طلب بكمية محددة | `POST /store/order` (`quantity: 5`) | `quantity == 5` |
| 20 | دورة حياة الطلب الكاملة | `POST → GET → DELETE → GET` | GET الأخير يُعيد 404 |

#### المستخدمون — نقاط نهاية `/user` (21–30)

| # | السيناريو | نقطة النهاية | يتحقق من |
|---|---|---|---|
| 21 | إنشاء مستخدم | `POST /user` + `GET /user/{username}` | الحقول متطابقة |
| 22 | الحصول باسم المستخدم | `GET /user/{username}` | `username` متطابق |
| 23 | تحديث مستخدم | `PUT /user/{username}` + `GET` | البريد الإلكتروني المحدَّث مُعاد |
| 24 | حذف مستخدم | `DELETE /user/{username}` | GET التالي يُعيد 404 |
| 25 | تسجيل دخول صالح | `GET /user/login` | الحالة 200، token في الجسم |
| 26 | تسجيل الخروج | `GET /user/logout` | الحالة 200 |
| 27 | إنشاء مستخدمين بمصفوفة | `POST /user/createWithArray` | كل مستخدم قابل للاسترجاع |
| 28 | إنشاء مستخدمين بقائمة | `POST /user/createWithList` | كل مستخدم قابل للاسترجاع |
| 29 | تسجيل دخول غير صالح | `GET /user/login` (بيانات خاطئة) | الاستجابة ليست 2xx |
| 30 | دورة حياة المستخدم الكاملة | `POST → login → GET → PUT → DELETE → GET` | GET الأخير يُعيد 404 |

أوصاف السيناريوهات التفصيلية: [`examples/api-automation/api-automation.md`](examples/api-automation/api-automation.md)

---

## التشغيل على متصفحات متعددة

يمكن تشغيل اختبارات UI على Chromium (افتراضي) أو Firefox أو WebKit:

```bash
# Firefox
PLAYWRIGHT_CLI_PATH="..." PLAYWRIGHT_BROWSER=firefox go test -tags e2e -v -timeout 300s ./examples/ui-automation/...

# WebKit (محرك Safari)
PLAYWRIGHT_CLI_PATH="..." PLAYWRIGHT_BROWSER=webkit go test -tags e2e -v -timeout 300s ./examples/ui-automation/...
```

---

## تغطية الاختبارات

تتحقق هذه المكتبة عن طريق **3,318 اختبارًا** موزعة على ثلاث طبقات:

| الطبقة | الاختبارات | ما تغطيه |
|---|---|---|
| **مجموعة E2E للمكتبة** (`e2e/`) | 3,278 | قدرات Playwright الأساسية: المتصفح، الصفحة، المحدد، الإطارات، الشبكة، طلبات API، التحقق، التتبع، لقطات الشاشة والمزيد — 124 ملف اختبار تغطي سطح API بالكامل |
| **أمثلة أتمتة UI** (`examples/ui-automation/`) | 10 | تدفقات متصفح كاملة على SauceDemo باستخدام نمط Page Object Model |
| **أمثلة أتمتة API** (`examples/api-automation/`) | 30 | سيناريوهات REST API على Petstore API العامة — دون الحاجة إلى متصفح |

تم بناء مجموعة E2E للمكتبة باستخدام [playwright-java](https://github.com/microsoft/playwright-java) كمرجع، وذلك بنقل حالات اختباره وتكييفها مع Go لضمان تغطية سلوكية مكافئة.

---

## الرخصة

رخصة MIT — راجع [LICENSE](LICENSE) للتفاصيل.

[![Go Reference](https://pkg.go.dev/badge/github.com/jcfrperu/goaneco-playwright.svg)](https://pkg.go.dev/github.com/jcfrperu/goaneco-playwright)
[![Go Report Card](https://goreportcard.com/badge/github.com/jcfrperu/goaneco-playwright)](https://goreportcard.com/report/github.com/jcfrperu/goaneco-playwright)