//go:build e2e

// Additional E2E tests without an equivalent in playwright-java.
package e2e

// overlayHTML returns HTML for a button that shows a dismissable overlay on click.
// Used by locator handler tests.
func overlayHTML() string {
	return `
<button id="action">Do Action</button>
<div id="overlay" style="display:none">
  <p>Overlay!</p>
  <button id="dismiss">Dismiss</button>
</div>
<div id="result"></div>
<script>
  document.getElementById('action').addEventListener('click', function() {
    document.getElementById('overlay').style.display = 'block';
  });
  document.getElementById('dismiss').addEventListener('click', function() {
    document.getElementById('overlay').style.display = 'none';
    document.getElementById('result').textContent = 'done';
  });
</script>`
}
