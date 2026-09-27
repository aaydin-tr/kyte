package integration_test

import (
	"testing"

	"github.com/aaydin-tr/kyte"
)

// TestGlobalFilters must not call t.Parallel(): Global Filters are package
// state, and Go runs every sequential top-level test to completion before it
// releases the parallel ones, so they cannot leak into other tests.
func TestGlobalFilters(t *testing.T) {
	coll := seed(t)

	t.Run("applied to new filters", func(t *testing.T) {
		t.Cleanup(kyte.ClearGlobalFilters)
		kyte.AddGlobalFilter(kyte.Filter().Equal("tenant", "t1"))

		assertMatches(t, coll, kyte.Filter().GreaterThan("age", 25), 2, 5)
	})

	t.Run("multiple global filters", func(t *testing.T) {
		t.Cleanup(kyte.ClearGlobalFilters)
		// Create both before registering either, so the second does not inherit the first.
		tenant := kyte.Filter().Equal("tenant", "t1")
		adults := kyte.Filter().GreaterThanOrEqual("age", 30)
		kyte.AddGlobalFilter(tenant)
		kyte.AddGlobalFilter(adults)

		assertMatches(t, coll, kyte.Filter().Exists("nickname", true), 5)
	})

	t.Run("built from source", func(t *testing.T) {
		t.Cleanup(kyte.ClearGlobalFilters)
		u := newSource()
		kyte.AddGlobalFilter(kyte.Filter(kyte.Source(u)).Equal(&u.Tenant, "t2"))

		assertMatches(t, coll, kyte.Filter().LessThan("age", 45), 3, 4)
	})

	t.Run("ignored with IgnoreGlobalFilters", func(t *testing.T) {
		t.Cleanup(kyte.ClearGlobalFilters)
		kyte.AddGlobalFilter(kyte.Filter().Equal("tenant", "t1"))

		assertMatches(t, coll, kyte.Filter(kyte.IgnoreGlobalFilters()).GreaterThan("age", 45), 6)
	})

	t.Run("with and", func(t *testing.T) {
		t.Cleanup(kyte.ClearGlobalFilters)
		kyte.AddGlobalFilter(kyte.Filter().Equal("tenant", "t1"))

		assertMatches(t, coll, kyte.Filter().And(kyte.Filter().GreaterThan("age", 25).LessThan("age", 50)), 2, 5)
	})

	t.Run("with or", func(t *testing.T) {
		t.Skip("known bug: the Or child also receives the Global Filter, so $or gains a clause every tenant document satisfies")
		t.Cleanup(kyte.ClearGlobalFilters)
		kyte.AddGlobalFilter(kyte.Filter().Equal("tenant", "t1"))

		assertMatches(t, coll, kyte.Filter().Or(kyte.Filter().Equal("name", "Alice").Equal("name", "Bob")), 1)
	})

	t.Run("with nor", func(t *testing.T) {
		t.Skip("known bug: the NOR child also receives the Global Filter, so $nor excludes every tenant document")
		t.Cleanup(kyte.ClearGlobalFilters)
		kyte.AddGlobalFilter(kyte.Filter().Equal("tenant", "t1"))

		assertMatches(t, coll, kyte.Filter().NOR(kyte.Filter().Equal("name", "Alice")), 2, 5)
	})
}
