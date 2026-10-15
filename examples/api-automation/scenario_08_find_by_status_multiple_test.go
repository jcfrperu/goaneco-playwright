//go:build e2e

package apiautomation

// Feature: Find pets by multiple statuses
// Scenario: Filtering by "available" and "pending" returns only pets with those statuses
//
// Feature file: features/08_find_by_status_multiple.feature

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScenario08FindByStatusMultiple(t *testing.T) {
	must := require.New(t)
	is := assert.New(t)

	ctx := testCtx(t)
	apiCtx := newAPICtx(t)

	resp, err := apiCtx.Get(ctx, "/pet/findByStatus?status=available&status=pending")
	must.NoError(err)
	defer resp.Dispose(ctx) //nolint:errcheck

	is.Equal(200, resp.Status())
	body, err := resp.Body(ctx)
	must.NoError(err)

	pets := mustUnmarshalPets(t, body)
	is.NotEmpty(pets, "expected at least one pet")
	for _, p := range pets {
		is.Contains([]string{"available", "pending"}, p.Status,
			"pet %d has unexpected status %q", p.ID, p.Status)
	}
}
