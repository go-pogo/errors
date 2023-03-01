// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !go1.20
// +build !go1.20

package errors

func (m *multiErr) Is(target error) bool {
	for _, err := range m.errs {
		if Is(err, target) {
			return true
		}
	}
	return false
}

func (m *multiErr) As(target interface{}) bool {
	for _, err := range m.errs {
		if As(err, target) {
			return true
		}
	}
	return false
}
