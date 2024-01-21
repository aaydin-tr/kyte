package kyte

import (
	"testing"
)

func TestFilter(t *testing.T) {
	t.Parallel()

	t.Run("without options", func(t *testing.T) {
		filter := Filter()
		if filter == nil {
			t.Error("Filter should not be nil")
		}

		if filter.kyte.source != nil {
			t.Error("Filter.kyte should be nil")
		}

		if filter.kyte.checkField != false {
			t.Error("Filter.kyte.checkField should be false")
		}
	})

	t.Run("with options", func(t *testing.T) {
		type Temp struct {
			Name string `bson:"name"`
		}

		filter := Filter(Source(&Temp{}), ValidateField(true))
		if filter == nil {
			t.Error("Filter should not be nil")
		}

		if filter.kyte.source == nil {
			t.Error("Filter.kyte should not be nil")
		}

		if filter.kyte.checkField != true {
			t.Error("Filter.kyte.checkField should be true")
		}
	})

	t.Run("with validate field is false", func(t *testing.T) {
		filter := Filter(ValidateField(false))
		if filter == nil {
			t.Error("Filter should not be nil")
		}

		if filter.kyte.source != nil {
			t.Error("Filter.kyte should be nil")
		}

		if filter.kyte.checkField != false {
			t.Error("Filter.kyte.checkField should be false")
		}
	})

}
