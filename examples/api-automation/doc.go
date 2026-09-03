// Package apiautomation demonstrates API test automation for the Petstore public API
// (https://petstore.swagger.io/v2) using the goaneco-playwright Go client with
// Playwright's APIRequestContext for HTTP requests.
//
// Run all scenarios:
//
//	go test -tags e2e -v -timeout 300s ./examples/api-automation/...
//
// Run a single scenario:
//
//	go test -tags e2e -v -run TestScenario01AddPet ./examples/api-automation/...
package apiautomation
