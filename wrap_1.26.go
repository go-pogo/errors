// Copyright (c) 2026, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build go1.26

package errors

import (
	stderrors "errors"
)

// AsType is an alias of [errors.AsType] and finds the first error in err's
// tree that matches the type E, and if one is found, returns that error value
// and true. Otherwise, it returns the zero value of E and false.
func AsType[E error](err error) (E, bool) {
	return stderrors.AsType[E](err)
}
