# goaneco-playwright

[Microsoft Playwright](https://playwright.dev) के साथ टेस्ट ऑटोमेशन के लिए Go लाइब्रेरी।
Go में UI और API टेस्ट एक ही टूल से लिखें, जो **Chromium**, **Firefox** और **WebKit** पर
चलते हैं।

यह प्रोजेक्ट ओपन सोर्स है और [MIT लाइसेंस](LICENSE) के तहत जारी किया गया है। यह मुख्य रूप से
व्यक्तिगत उपयोग के लिए बनाया गया था, लेकिन विचार, सुझाव और सुधार का स्वागत है — जो योगदान
प्रोजेक्ट को अधिक उपयोगी बनाएं, उन्हें हमेशा विचार किया जाएगा। आप इसे fork करके अपनी
दिशा में भी ले जा सकते हैं।

---

## मैं इससे क्या टेस्ट कर सकता हूँ?

| टेस्ट का प्रकार | विवरण | शामिल उदाहरण |
|---|---|---|
| **UI / E2E** | असली ब्राउज़र को ऑटोमेट करता है: क्लिक, फॉर्म, नेविगेशन, स्क्रीनशॉट | [SauceDemo](examples/ui-automation/) — 10 परिदृश्य |
| **REST API** | HTTP अनुरोध भेजता है (GET, POST, PUT, DELETE) और JSON प्रतिक्रियाएं सत्यापित करता है — ब्राउज़र की जरूरत नहीं | [Petstore API](examples/api-automation/) — 30 परिदृश्य |

---

## सेटअप

### 1. पूर्वावश्यकताएं इंस्टॉल करें

| टूल | न्यूनतम संस्करण | इंस्टॉलेशन |
|---|---|---|
| [Go](https://go.dev/dl/) | 1.21 | go.dev से डाउनलोड और इंस्टॉल करें |
| [Node.js](https://nodejs.org/) | 18 | nodejs.org से डाउनलोड और इंस्टॉल करें |
| Playwright | — | `npm install -g playwright` |

### 2. ब्राउज़र इंस्टॉल करें

```bash
npx playwright install chromium
```

Linux पर, OS-स्तरीय निर्भरताएं भी इंस्टॉल करें:

```bash
npx playwright install --with-deps chromium
```

### 3. एनवायरनमेंट वेरिएबल सेट करें

सभी टेस्ट को यह जानना होता है कि Playwright CLI कहाँ है:

```bash
# Linux / macOS
export PLAYWRIGHT_CLI_PATH=$(npm root -g)/playwright/node_modules/playwright-core/cli.js

# Windows PowerShell
$env:PLAYWRIGHT_CLI_PATH = "$env:APPDATA\npm\node_modules\playwright\node_modules\playwright-core\cli.js"
```

इस लाइन को अपने shell प्रोफाइल (`~/.bashrc`, `~/.zshrc` आदि) में जोड़ें ताकि इसे हर सत्र में
दोहराना न पड़े।

---

## टेस्ट कैसे दिखते हैं

टेस्ट **परिदृश्यों** के रूप में व्यवस्थित हैं — प्रति टेस्ट केस एक फ़ाइल। प्रत्येक टेस्ट एक Go
फ़ंक्शन है जो स्पष्ट रूप से बताता है कि वह क्या सत्यापित कर रहा है।

### UI टेस्ट (Page Object Model के साथ)

लक्ष्य: [https://www.saucedemo.com](https://www.saucedemo.com) — UI ऑटोमेशन अभ्यास के लिए एक डेमो ई-कॉमर्स साइट।

```go
func TestLoginAndInventory(t *testing.T) {
    must := require.New(t)
    is   := assert.New(t)
    page := newPage(t)           // नया ब्राउज़र टैब खोलता है
    loginStandardUser(t, page)   // मानक क्रेडेंशियल से लॉगिन करने वाला हेल्पर

    ctx := testCtx(t)
    inv := pom.NewInventoryPage(page)

    count, err := inv.ProductCount(ctx)
    must.NoError(err)
    is.Equal(6, count)   // इन्वेंटरी में बिल्कुल 6 उत्पाद दिखने चाहिए
}
```

**Page Object Model (POM)** प्रत्येक स्क्रीन के एक्शन को एक पुनः उपयोगी struct में लपेटता है
(`pom.InventoryPage`, `pom.CartPage` आदि), जिससे टेस्ट साफ और पढ़ने में आसान रहते हैं।

### API टेस्ट (पूर्ण जीवनचक्र)

लक्ष्य: [https://petstore.swagger.io/v2](https://petstore.swagger.io/v2) — API ऑटोमेशन अभ्यास के लिए एक सार्वजनिक REST API सैंडबॉक्स।

```go
func TestScenario12PetFullLifecycle(t *testing.T) {
    must := require.New(t)
    is   := assert.New(t)
    ctx    := testCtx(t)
    apiCtx := newAPICtx(t)   // API बेस URL की ओर इंगित HTTP क्लाइंट

    // एक पालतू जानवर बनाएं
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

    // सत्यापित करें कि इसे पुनः प्राप्त किया जा सकता है
    getResp, err := apiCtx.Get(ctx, fmt.Sprintf("/pet/%d", created.ID))
    must.NoError(err)
    defer getResp.Dispose(ctx)
    is.Equal(200, getResp.Status())

    // स्टेटस को "sold" में अपडेट करें
    created.Status = "sold"
    putResp, err := apiCtx.Put(ctx, "/pet", &playwright.APIFetchOptions{Headers: jsonHeaders(), Data: mustJSON(t, created)})
    must.NoError(err)
    defer putResp.Dispose(ctx)
    is.Equal(200, putResp.Status())

    // पालतू जानवर को हटाएं
    delResp, err := apiCtx.Delete(ctx, fmt.Sprintf("/pet/%d", created.ID))
    must.NoError(err)
    defer delResp.Dispose(ctx)
    is.Equal(200, delResp.Status())

    // पुष्टि करें कि अब मौजूद नहीं है (404)
    confirmResp, err := apiCtx.Get(ctx, fmt.Sprintf("/pet/%d", created.ID))
    must.NoError(err)
    defer confirmResp.Dispose(ctx)
    is.Equal(404, confirmResp.Status())
}
```

---

## शामिल उदाहरण

### UI ऑटोमेशन — SauceDemo (`examples/ui-automation/`)

**Page Object Model** पैटर्न का उपयोग करके [SauceDemo](https://www.saucedemo.com) के लिए
ब्राउज़र ऑटोमेशन। 10 परिदृश्य एक ऑनलाइन स्टोर के पूरे प्रवाह को कवर करते हैं।

```bash
# सभी परिदृश्य चलाएं
PLAYWRIGHT_CLI_PATH="..." go test -tags e2e -v -timeout 300s ./examples/ui-automation/...

# एक परिदृश्य चलाएं
PLAYWRIGHT_CLI_PATH="..." go test -tags e2e -v -timeout 120s -run TestFullCheckoutFlow ./examples/ui-automation/...
```

| # | परिदृश्य | क्या सत्यापित करता है |
|---|---|---|
| 01 | लॉगिन और इन्वेंटरी | प्रमाणीकरण स्मोक टेस्ट; 6 उत्पाद दृश्यमान हैं |
| 02 | कार्ट में आइटम जोड़ें | 3 उत्पाद जोड़ता है; कार्ट बैज काउंट और सूची सत्यापित करता है |
| 03 | कम से ज्यादा कीमत से क्रमबद्ध | सबसे सस्ता आइटम (Onesie, $7.99) पहले दिखता है |
| 04 | Z से A नाम से क्रमबद्ध | उत्पाद उल्टे वर्णमाला क्रम में सूचीबद्ध हैं |
| 05 | पूर्ण चेकआउट प्रक्रिया | जोड़ें → कार्ट → जानकारी भरें → सारांश → पुष्टि |
| 06 | कार्ट से आइटम हटाएं | 3 जोड़ता है, 1 हटाता है, 2 शेष की पुष्टि करता है |
| 07 | उत्पाद विवरण पृष्ठ | नाम, कीमत और विवरण सत्यापित करता है |
| 08 | बर्गर मेनू से लॉगआउट | साइड नेव खोलता है, लॉगआउट क्लिक करता है, लॉगिन स्क्रीन की पुष्टि करता है |
| 09 | सबसे महंगा आइटम जोड़ें | ऊपर से नीचे क्रमबद्ध; Fleece Jacket ($49.99) जोड़ता है |
| 10 | फुटर सोशल मीडिया लिंक | Twitter, Facebook और LinkedIn लिंक दृश्यमान हैं |

स्क्रीन `examples/ui-automation/pom/` में **Page Objects** के रूप में अमूर्त हैं:

| Page Object | उपलब्ध एक्शन |
|---|---|
| `LoginPage` | `Navigate`, `Login`, `IsVisible`, `ErrorMessage` |
| `InventoryPage` | `ProductCount`, `AddToCart`, `CartBadgeCount`, `Sort`, `Logout` |
| `CartPage` | `ItemCount`, `GetItemNames`, `GetItemPrices`, `RemoveItem`, `Checkout` |
| `CheckoutInfoPage` | `FillInfo`, `Continue`, `ErrorMessage` |
| `CheckoutOverviewPage` | `GetItemCount`, `GetSubtotal`, `GetTotal`, `Finish` |
| `CheckoutCompletePage` | `GetHeader`, `BackHome` |
| `ProductDetailPage` | `GetName`, `GetPrice`, `GetDescription`, `AddToCart`, `BackToInventory` |

पूर्ण परिदृश्य विवरण: [`examples/ui-automation/ui-automation.md`](examples/ui-automation/ui-automation.md)

---

### API ऑटोमेशन — Petstore (`examples/api-automation/`)

सार्वजनिक [Petstore API](https://petstore.swagger.io/v2) के विरुद्ध 30 टेस्ट परिदृश्य।
कोई ब्राउज़र लॉन्च नहीं होता — केवल HTTP अनुरोध और प्रतिक्रिया सत्यापन।

```bash
# सभी परिदृश्य चलाएं
PLAYWRIGHT_CLI_PATH="..." go test -tags e2e -v -timeout 300s ./examples/api-automation/...

# एक परिदृश्य चलाएं
PLAYWRIGHT_CLI_PATH="..." go test -tags e2e -v -timeout 60s -run TestScenario01AddPet ./examples/api-automation/...
```

#### पालतू जानवर — `/pet` एंडपॉइंट (01–12)

| # | परिदृश्य | एंडपॉइंट | सत्यापित करता है |
|---|---|---|---|
| 01 | पालतू जानवर जोड़ें | `POST /pet` | `id > 0`, नाम और स्टेटस सही |
| 02 | ID से प्राप्त करें | `GET /pet/{id}` | फ़ील्ड बनाए गए से मेल खाते हैं |
| 03 | अपडेट करें | `PUT /pet` | स्टेटस `"sold"` में अपडेट |
| 04 | हटाएं | `DELETE /pet/{id}` | बाद का GET 404 लौटाता है |
| 05–07 | स्टेटस से खोजें | `GET /pet/findByStatus` | available / pending / sold |
| 08 | कई स्टेटस से खोजें | `GET /pet/findByStatus?status=available&status=pending` | केवल ये दो स्टेटस |
| 09 | फॉर्म डेटा से अपडेट | `POST /pet/{id}` (form-encoded) | स्टेटस 200 |
| 10 | टैग से खोजें | `GET /pet/findByTags?tags=goaneco` | वैध JSON array |
| 11 | अनुपस्थित पालतू जानवर | `GET /pet/999999999` | स्टेटस 404 |
| 12 | पूर्ण जीवनचक्र | `POST → GET → PUT → DELETE → GET` | अंतिम GET 404 |

#### स्टोर — `/store` एंडपॉइंट (13–20)

| # | परिदृश्य | एंडपॉइंट | सत्यापित करता है |
|---|---|---|---|
| 13 | इन्वेंटरी प्राप्त करें | `GET /store/inventory` | गैर-खाली स्टेटस-काउंट मैप |
| 14 | ऑर्डर दें | `POST /store/order` | `id > 0`, `status == "placed"` |
| 15 | ID से ऑर्डर प्राप्त करें | `GET /store/order/{id}` | फ़ील्ड ऑर्डर से मेल खाते हैं |
| 16 | ऑर्डर हटाएं | `DELETE /store/order/{id}` | बाद का GET 404 लौटाता है |
| 17 | पूर्ण ऑर्डर | `POST /store/order` (`complete: true`) | `complete == true` |
| 18 | अनुपस्थित ऑर्डर | `GET /store/order/999999` | स्टेटस 404 |
| 19 | मात्रा के साथ ऑर्डर | `POST /store/order` (`quantity: 5`) | `quantity == 5` |
| 20 | पूर्ण ऑर्डर जीवनचक्र | `POST → GET → DELETE → GET` | अंतिम GET 404 |

#### उपयोगकर्ता — `/user` एंडपॉइंट (21–30)

| # | परिदृश्य | एंडपॉइंट | सत्यापित करता है |
|---|---|---|---|
| 21 | उपयोगकर्ता बनाएं | `POST /user` + `GET /user/{username}` | फ़ील्ड मेल खाते हैं |
| 22 | उपयोगकर्ता नाम से प्राप्त | `GET /user/{username}` | `username` मेल खाता है |
| 23 | अपडेट करें | `PUT /user/{username}` + `GET` | अपडेट ईमेल लौटाया गया |
| 24 | हटाएं | `DELETE /user/{username}` | बाद का GET 404 लौटाता है |
| 25 | वैध लॉगिन | `GET /user/login` | स्टेटस 200, body में token |
| 26 | लॉगआउट | `GET /user/logout` | स्टेटस 200 |
| 27 | array से उपयोगकर्ता बनाएं | `POST /user/createWithArray` | प्रत्येक GET से प्राप्य |
| 28 | सूची से उपयोगकर्ता बनाएं | `POST /user/createWithList` | प्रत्येक GET से प्राप्य |
| 29 | अमान्य लॉगिन | `GET /user/login` (गलत क्रेडेंशियल) | प्रतिक्रिया 2xx नहीं |
| 30 | पूर्ण उपयोगकर्ता जीवनचक्र | `POST → login → GET → PUT → DELETE → GET` | अंतिम GET 404 |

पूर्ण परिदृश्य विवरण: [`examples/api-automation/api-automation.md`](examples/api-automation/api-automation.md)

---

## कई ब्राउज़रों पर चलाएं

UI टेस्ट Chromium (डिफ़ॉल्ट), Firefox या WebKit पर चल सकते हैं:

```bash
# Firefox
PLAYWRIGHT_CLI_PATH="..." PLAYWRIGHT_BROWSER=firefox go test -tags e2e -v -timeout 300s ./examples/ui-automation/...

# WebKit (Safari इंजन)
PLAYWRIGHT_CLI_PATH="..." PLAYWRIGHT_BROWSER=webkit go test -tags e2e -v -timeout 300s ./examples/ui-automation/...
```

---

## टेस्ट कवरेज

यह लाइब्रेरी तीन स्तरों में फैले **3,318 टेस्ट** द्वारा मान्य है:

| स्तर | टेस्ट | क्या कवर करता है |
|---|---|---|
| **लाइब्रेरी E2E सूट** (`e2e/`) | 3,278 | Playwright की मुख्य क्षमताएं: browser, page, locator, frames, network, API requests, assertions, tracing, screenshots और अधिक — संपूर्ण API सतह को कवर करने वाली 124 टेस्ट फ़ाइलें |
| **UI ऑटोमेशन उदाहरण** (`examples/ui-automation/`) | 10 | Page Object Model पैटर्न का उपयोग करके SauceDemo पर पूर्ण ब्राउज़र फ्लो |
| **API ऑटोमेशन उदाहरण** (`examples/api-automation/`) | 30 | सार्वजनिक Petstore API के विरुद्ध REST API परिदृश्य — कोई ब्राउज़र आवश्यक नहीं |

लाइब्रेरी E2E सूट को [playwright-java](https://github.com/microsoft/playwright-java) को संदर्भ के रूप में उपयोग करके बनाया गया था, इसके टेस्ट केस को Go में पोर्ट और अनुकूलित करके समकक्ष व्यवहार कवरेज सुनिश्चित किया गया।

---

## लाइसेंस

MIT लाइसेंस — विवरण के लिए [LICENSE](LICENSE) देखें।

[![Go Reference](https://pkg.go.dev/badge/github.com/jcfrperu/goaneco-playwright.svg)](https://pkg.go.dev/github.com/jcfrperu/goaneco-playwright)
[![Go Report Card](https://goreportcard.com/badge/github.com/jcfrperu/goaneco-playwright)](https://goreportcard.com/report/github.com/jcfrperu/goaneco-playwright)