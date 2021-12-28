// Copyright (c) 2021, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errors

import (
	stderrors "errors"
)

type Isser interface {
	error
	Is(target error) bool
}

// Is reports whether any error in err's chain matches target. It is fully
// compatible with both errors.Is and xerrors.Is.
//
// An error is considered to match a target if it is equal to that target or if
// it implements a method Is(error) bool such that Is(target) returns true.
func Is(err, target error) bool { return stderrors.Is(err, target) }

type Asser interface {
	error
	As(target interface{}) bool
}

// As is an alias of errors.As. It finds the first error in err's chain that
// matches target, and if so, sets target to that error value and returns true.
func As(err error, target interface{}) bool {
	//goland:noinspection GoErrorsAs
	return err != nil && stderrors.As(err, target)
}
