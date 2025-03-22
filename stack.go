package wrapper

import (
	"reflect"
	"sync"
)

type PreHook func(args ...any) error
type PostHook func(results ...any) error
type WrapOption func(*wrapOptions)

type wrapOptions struct {
	preHook   PreHook
	postHook  PostHook
	errorHook func(error)
}

var cache sync.Map // key: uintptr, value: *wrapMetadata

type wrapMetadata struct {
	fnType   reflect.Type
	numOut   int
	errIndex int // -1 if no error return
}

// Wrap wraps a function fn with optional preHook, postHook and errorHook logic.
// It still uses generics while caching reflection metadata for repeated calls.
func Wrap[T any](fn T, opts ...WrapOption) T {
	options := &wrapOptions{}
	for _, opt := range opts {
		opt(options)
	}

	fnValue := reflect.ValueOf(fn)
	fnType := fnValue.Type()
	if fnType.Kind() != reflect.Func {
		panic("Wrap expects a function")
	}

	// Try to get cached metadata for this function.
	metaIface, ok := cache.Load(fnValue.Pointer())
	var metadata *wrapMetadata
	if ok {
		metadata = metaIface.(*wrapMetadata)
	} else {
		errIndex := -1
		if fnType.NumOut() > 0 {
			lastType := fnType.Out(fnType.NumOut() - 1)
			errorType := reflect.TypeOf((*error)(nil)).Elem()
			if lastType.Implements(errorType) {
				errIndex = fnType.NumOut() - 1
			}
		}
		metadata = &wrapMetadata{
			fnType:   fnType,
			numOut:   fnType.NumOut(),
			errIndex: errIndex,
		}
		cache.Store(fnValue.Pointer(), metadata)
	}

	// Create the wrapped function.
	wrappedFn := reflect.MakeFunc(fnType, func(args []reflect.Value) (results []reflect.Value) {
		// Pre-hook: convert arguments to []any inline.
		if options.preHook != nil {
			argInterfaces := make([]any, len(args))
			for i, arg := range args {
				argInterfaces[i] = arg.Interface()
			}
			if err := options.preHook(argInterfaces...); err != nil {
				if options.errorHook != nil {
					options.errorHook(err)
				}
				return createErrorResults(metadata, err)
			}
		}

		results = fnValue.Call(args)

		// If the function returns an error, check it once.
		if metadata.errIndex != -1 {
			lastResult := results[metadata.errIndex]
			if !lastResult.IsNil() {
				if err, ok := lastResult.Interface().(error); ok && err != nil {
					if options.errorHook != nil {
						options.errorHook(err)
					}
					return createErrorResults(metadata, err)
				}
			}
		}

		// Post-hook: convert results to []any inline.
		if options.postHook != nil {
			resInterfaces := make([]any, len(results))
			for i, res := range results {
				resInterfaces[i] = res.Interface()
			}
			if err := options.postHook(resInterfaces...); err != nil {
				if options.errorHook != nil {
					options.errorHook(err)
				}
				return createErrorResults(metadata, err)
			}
		}

		return results
	}).Interface()

	return wrappedFn.(T)
}

// createErrorResults creates a slice of reflect.Values with zero values
// for non-error return types and sets the error value where appropriate.
func createErrorResults(meta *wrapMetadata, err error) []reflect.Value {
	results := make([]reflect.Value, meta.numOut)
	errorType := reflect.TypeOf((*error)(nil)).Elem()
	for i := 0; i < meta.numOut; i++ {
		if meta.fnType.Out(i) == errorType {
			results[i] = reflect.ValueOf(err)
		} else {
			results[i] = reflect.Zero(meta.fnType.Out(i))
		}
	}
	return results
}

func WithPreHook(hook PreHook) WrapOption {
	return func(opts *wrapOptions) {
		opts.preHook = hook
	}
}

func WithPostHook(hook PostHook) WrapOption {
	return func(opts *wrapOptions) {
		opts.postHook = hook
	}
}

func WithErrorHook(hook func(err error)) WrapOption {
	return func(opts *wrapOptions) {
		opts.errorHook = hook
	}
}
