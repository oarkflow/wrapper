package wrapper

/*
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

type Func0[R any] func() (R, error)
type Func1[A any, R any] func(A) (R, error)
type Func2[A, B any, R any] func(A, B) (R, error)
type Func3[A, B, C any, R any] func(A, B, C) (R, error)
type Func4[A, B, C, D any, R any] func(A, B, C, D) (R, error)
type Func5[A, B, C, D, E any, R any] func(A, B, C, D, E) (R, error)

func preHook(preHook PreHook, errHook func(err error), args ...any) error {
	if preHook != nil {
		err := preHook(args...)
		if err != nil && errHook != nil {
			errHook(err)
		}
		return err
	}
	return nil
}

func postHook(postHook PostHook, errHook func(err error), results ...any) error {
	if postHook != nil {
		err := postHook(results...)
		if err != nil && errHook != nil {
			errHook(err)
		}
		return err
	}
	return nil
}

func errHook(errHook func(err error), err error) error {
	if errHook != nil && err != nil {
		errHook(err)
		return err
	}
	return nil
}

func defaultWrapOptions(opts ...WrapOption) *wrapOptions {
	o := &wrapOptions{}
	for _, opt := range opts {
		opt(o)
	}
	return o
}

func Wrap0[R any](fn Func0[R], opts ...WrapOption) Func0[R] {
	o := defaultWrapOptions(opts...)
	return func() (R, error) {
		var zr R
		err := preHook(o.preHook, o.errorHook)
		if err != nil {
			return zr, err
		}
		zr, err = fn()
		if err != nil {
			return zr, errHook(o.errorHook, err)
		}
		err = postHook(o.postHook, o.errorHook, zr)
		return zr, err
	}
}

func Wrap1[A any, R any](fn Func1[A, R], opts ...WrapOption) Func1[A, R] {
	o := defaultWrapOptions(opts...)
	return func(a A) (R, error) {
		var zr R
		err := preHook(o.preHook, o.errorHook, a)
		if err != nil {
			return zr, err
		}
		zr, err = fn(a)
		if err != nil {
			return zr, errHook(o.errorHook, err)
		}
		err = postHook(o.postHook, o.errorHook, zr)
		return zr, err
	}
}

func Wrap2[A, B any, R any](fn Func2[A, B, R], opts ...WrapOption) Func2[A, B, R] {
	o := defaultWrapOptions(opts...)
	return func(a A, b B) (R, error) {
		var zr R
		err := preHook(o.preHook, o.errorHook, a, b)
		if err != nil {
			return zr, err
		}
		zr, err = fn(a, b)
		if err != nil {
			return zr, errHook(o.errorHook, err)
		}
		err = postHook(o.postHook, o.errorHook, zr)
		return zr, err
	}
}

func Wrap3[A, B, C any, R any](fn Func3[A, B, C, R], opts ...WrapOption) Func3[A, B, C, R] {
	o := defaultWrapOptions(opts...)
	return func(a A, b B, c C) (R, error) {
		var zr R
		err := preHook(o.preHook, o.errorHook, a, b, c)
		if err != nil {
			return zr, err
		}
		zr, err = fn(a, b, c)
		if err != nil {
			return zr, errHook(o.errorHook, err)
		}
		err = postHook(o.postHook, o.errorHook, zr)
		return zr, err
	}
}

func Wrap4[A, B, C, D any, R any](fn Func4[A, B, C, D, R], opts ...WrapOption) Func4[A, B, C, D, R] {
	o := defaultWrapOptions(opts...)
	return func(a A, b B, c C, d D) (R, error) {
		var zr R
		err := preHook(o.preHook, o.errorHook, a, b, c, d)
		if err != nil {
			return zr, err
		}
		zr, err = fn(a, b, c, d)
		if err != nil {
			return zr, errHook(o.errorHook, err)
		}
		err = postHook(o.postHook, o.errorHook, zr)
		return zr, err
	}
}

func Wrap5[A, B, C, D, E any, R any](fn Func5[A, B, C, D, E, R], opts ...WrapOption) Func5[A, B, C, D, E, R] {
	o := defaultWrapOptions(opts...)
	return func(a A, b B, c C, d D, e E) (R, error) {
		var zr R
		err := preHook(o.preHook, o.errorHook, a, b, c, d, e)
		if err != nil {
			return zr, err
		}
		zr, err = fn(a, b, c, d, e)
		if err != nil {
			return zr, errHook(o.errorHook, err)
		}
		err = postHook(o.postHook, o.errorHook, zr)
		return zr, err
	}
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
*/
