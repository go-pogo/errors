package errs

import (
	"fmt"
)

// WrapPanic wraps a panicking sequence with the given prefix. It then panics
// again.
func WrapPanic(prefix string) {
	if r := recover(); r != nil {
		panic(fmt.Sprintf("%s: %s", prefix, r))
	}
}
