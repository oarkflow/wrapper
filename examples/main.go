// main.go
package main

import (
	"errors"
	"fmt"

	"github.com/oarkflow/wrapper"
)

func main() {
	add := func(a int, b float32) (int, error) {
		if a < 0 || b < 0 {
			return 0, errors.New("negative input")
		}
		return a + int(b), nil
	}

	// Pre-hook: log inputs
	pre := func(args ...any) error {
		fmt.Printf("About to add %v\n", args)
		return nil
	}

	// Post-hook: ensure result is even
	post := func(args ...any) error {
		fmt.Printf("Result from add %v\n", args)
		return nil
	}

	// Error-hook: log error
	onErr := func(err error) {
		fmt.Printf("❌ error: %v\n", err)
	}

	wrappedAdd := wrapper.Wrap2(add, wrapper.WithPreHook(pre), wrapper.WithPostHook(post), wrapper.WithErrorHook(onErr))

	if res, err := wrappedAdd(1, 2); err != nil {
		fmt.Println("wrappedAdd error:", err)
	} else {
		fmt.Println("wrappedAdd result:", res)
	}
}
