//go:build e2e

// Additional E2E tests without an equivalent in playwright-java.
package e2e

// tagSelectorScript is a custom selector engine that maps selectors to HTML tag names.
const tagSelectorScript = `{
  create(root, target) {
    return target.nodeName;
  },
  query(root, selector) {
    return root.querySelector(selector);
  },
  queryAll(root, selector) {
    return Array.from(root.querySelectorAll(selector));
  }
}`
