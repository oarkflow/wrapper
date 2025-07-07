package wrapper_test

import (
	"errors"
	"testing"

	"github.com/oarkflow/wrapper"
)

// add is the function under test: returns an error on negative inputs.
func add(a, b int) (int, error) {
	if a < 0 || b < 0 {
		return 0, errors.New("inputs must be non‑negative")
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

// BenchmarkWrap2WithHooks measures Wrap2(add) with no‑op hooks.
func BenchmarkWrap2WithHooks(b *testing.B) {
	wrapped := wrapper.Wrap2(
		add,
		wrapper.WithPreHook(func(args ...any) error { return nil }),
		wrapper.WithPostHook(func(results ...any) error { return nil }),
		wrapper.WithErrorHook(func(err error) {}),
	)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wrapped(5, 7)
	}
}

// BenchmarkWrap2NoHooks measures Wrap2(add) with zero options.
func BenchmarkWrap2NoHooks(b *testing.B) {
	wrapped := wrapper.Wrap2(add)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wrapped(5, 7)
	}
}

// BenchmarkRawAdd measures the unwrapped add(a,b) call.
func BenchmarkRawAdd(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		add(5, 7)
	}
}
