package kyte

import (
	"fmt"
	"testing"
)

func BenchmarkFilterWithoutGlobal(b *testing.B) {
	sizes := []int{10, 100, 1000, 10000}
	for _, size := range sizes {
		b.Run(fmt.Sprintf("input_size_%d", size), func(b *testing.B) {
			ClearGlobalFilters()
			filter := buildFilterWithSize(size)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _ = filter.Build()
			}
		})
	}
}

func BenchmarkFilterWithGlobal(b *testing.B) {
	sizes := []int{10, 100, 1000, 10000}
	for _, size := range sizes {
		b.Run(fmt.Sprintf("input_size_%d", size), func(b *testing.B) {
			ClearGlobalFilters()
			// Add global filters
			for i := 0; i < 5; i++ {
				AddGlobalFilter(Filter().Equal(fmt.Sprintf("global%c", rune('A'+i)), i))
			}
			filter := buildFilterWithSize(size)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _ = filter.Build()
			}
		})
	}
}

func buildFilterWithSize(size int) *filter {
	f := Filter()
	// Add base conditions
	f.Equal("name", "John").
		GreaterThan("age", 18).
		In("roles", []string{"admin", "user"})

	// Add additional conditions based on size
	for i := 0; i < size; i++ {
		fieldName := fmt.Sprintf("field%c", rune('A'+i%26))
		f.Equal(fieldName, i)
	}

	return f
}
