package wrapper_test

import (
	"errors"
	"testing"

	"github.com/oarkflow/wrapper"
)

func add(a, b int) (int, error) {
	if a < 0 || b < 0 {
		return 0, errors.New("inputs must be non-negative")
	}
	return a + b, nil
}

// Benchmark wrapped function
func BenchmarkWrappedAdd(b *testing.B) {
	preHook := func(args ...any) error { return nil }
	postHook := func(results ...any) error { return nil }
	errorHook := func(err error) {}
	wrappedAdd := wrapper.Wrap(add, wrapper.WithPreHook(preHook), wrapper.WithPostHook(postHook), wrapper.WithErrorHook(errorHook))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wrappedAdd(5, 7)
	}
}

// Benchmark non-wrapped function
func BenchmarkNonWrappedAdd(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		add(5, 7)
	}
}
