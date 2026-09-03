//go:build e2e

// PDF generation E2E tests. Chromium headless only.
// Migration of: TestPdf.java
package e2e

import (
	"net/http"
	"os"
	"testing"

	playwright "github.com/jcfrperu/goaneco-playwright"
	"github.com/jcfrperu/goaneco-playwright/e2e/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestPdfSaveFile verifies that Pdf returns a non-empty byte slice and that
// the bytes can be written to a file.
// Ref: TestPdf.java#shouldBeAbleToSaveFile
func TestPdfSaveFile(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	if globalBTName != "chromium" {
		t.Skip("PDF generation is Chromium headless only")
	}
	ctx := testCtx(t)
	page := newPage(t)

	pdfBytes, err := page.PDF(ctx)
	must.NoError(err)
	is.Greater(len(pdfBytes), 0)

	f, err := os.CreateTemp(t.TempDir(), "*.pdf")
	must.NoError(err)
	_, err = f.Write(pdfBytes)
	must.NoError(err)
	must.NoError(f.Close())

	info, err := os.Stat(f.Name())
	must.NoError(err)
	is.Greater(info.Size(), int64(0))
}

// TestPdfFractionalScale verifies that a fractional scale value (0.5) is accepted
// and produces a non-empty PDF.
// Ref: TestPdf.java#shouldSupportFractionalScaleValue
func TestPdfFractionalScale(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	if globalBTName != "chromium" {
		t.Skip("PDF generation is Chromium headless only")
	}
	ctx := testCtx(t)
	page := newPage(t)

	scale := 0.5
	pdfBytes, err := page.PDF(ctx, &playwright.PagePdfOptions{Scale: &scale})
	must.NoError(err)
	is.Greater(len(pdfBytes), 0)
}

// TestPdfOutlineLargerThanNoOutline verifies that generating a PDF with tagged=true
// and outline=true produces a larger file than a plain PDF, because the outline
// data is embedded.
// Ref: TestPdf.java#shouldBeAbleToGenerateOutline
func TestPdfOutlineLargerThanNoOutline(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	if globalBTName != "chromium" {
		t.Skip("PDF generation is Chromium headless only")
	}
	ctx := testCtx(t)
	page := newPage(t)

	srv := testserver.New(t)
	srv.SetRoute("/headings.html", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(`<!DOCTYPE html><html><body>
<h1>Heading 1</h1>
<h2>Heading 2</h2>
<h3>Heading 3</h3>
<p>Some content under each heading.</p>
</body></html>`))
	})

	must.NoError(page.Goto(ctx, srv.Prefix()+"/headings.html"))

	plainBytes, err := page.PDF(ctx)
	must.NoError(err)
	is.Greater(len(plainBytes), 0)

	tagged := true
	outline := true
	outlineBytes, err := page.PDF(ctx, &playwright.PagePdfOptions{
		Tagged:  &tagged,
		Outline: &outline,
	})
	must.NoError(err)
	is.Greater(len(outlineBytes), 0)

	assert.Greater(
		t,
		len(outlineBytes),
		len(plainBytes),
		"PDF with outline (%d bytes) should be larger than plain PDF (%d bytes)",
		len(outlineBytes),
		len(plainBytes),
	)
}
