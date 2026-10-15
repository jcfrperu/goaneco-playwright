//go:build e2e

package apiautomation

// Feature: Find pets by tags
// Scenario: Filtering by tag "goaneco" returns a valid array
//
// Feature file: features/10_find_by_tags.feature

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScenario10FindByTags(t *testing.T) {
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	apiCtx := newAPICtx(t)

	resp, err := apiCtx.Get(ctx, "/pet/findByTags?tags=goaneco")
	must.NoError(err)
	defer resp.Dispose(ctx) //nolint:errcheck

	is.Equal(200, resp.Status())
	body, err := resp.Body(ctx)
	must.NoError(err)

	pets := mustUnmarshalPets(t, body)
	_ = pets // result may be empty; just assert 200 and valid JSON array
}
