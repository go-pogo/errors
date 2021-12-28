// Copyright (c) 2021, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errors

import (
	stderrors "errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func embedders() map[string]func(parent error) error {
	return map[string]func(parent error) error{
		"WithFormatter": func(parent error) error {
			return WithFormatter(parent)
		},
		"WithStack": func(parent error) error {
			return WithStack(parent)
		},

		"WithExitCode": func(parent error) error {
			return WithExitCode(parent, 1)
		},
		"WithTime": func(parent error) error {
			return WithTime(parent, time.Now())
		},
	}
}

func TestUnembed(t *testing.T) {
	targets := map[string]error{
		"error":     New("original"),
		"std error": stderrors.New("original std error"),
	}

	for targetName, want := range targets {
		t.Run(targetName, func(t *testing.T) {
			for name, embedFn := range embedders() {
				t.Run(name, func(t *testing.T) {
					have := embedFn(want)
					assert.Same(t, want, Unembed(have))
					assertErrorIs(t, have, want)
				})
			}
		})
	}
}
