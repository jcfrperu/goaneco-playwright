//go:build e2e

package apiautomation

// Feature: Find pets by status (available, pending, sold)
// Scenario: All pets returned when filtering by a given status have that status
//
// Feature file: features/05_find_by_status_available.feature

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScenario05FindByStatus(t *testing.T) {
	t.Parallel()
	must := require.New(t)
	is := assert.New(t)

	statuses := []string{"available", "pending", "sold"}
	for _, status := range statuses {
		status := status
		t.Run(status, func(t *testing.T) {
			t.Parallel()
			ctx := testCtx(t)
			apiCtx := newAPICtx(t)

			resp, err := apiCtx.Get(ctx, "/pet/findByStatus?status="+status)
			must.NoError(err)
			defer resp.Dispose(ctx) //nolint:errcheck

			is.Equal(200, resp.Status())
			body, err := resp.Body(ctx)
			must.NoError(err)

			pets := mustUnmarshalPets(t, body)
			is.NotEmpty(pets, "expected at least one %s pet", status)
			for _, p := range pets {
				is.Equal(status, p.Status, "pet %d has unexpected status", p.ID)
			}
		})
	}
}
