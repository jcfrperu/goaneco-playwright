//go:build e2e

// Additional E2E tests without an equivalent in playwright-java.
package e2e

const selectHTML = `<select id="sel">
	<option value="foo">Foo</option>
	<option value="bar">Bar</option>
	<option value="baz">Baz</option>
</select>`

const selectHTMLMulti = `<select id="sel" multiple>
	<option value="foo">Foo</option>
	<option value="bar">Bar</option>
	<option value="baz">Baz</option>
</select>`
