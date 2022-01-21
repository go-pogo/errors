// Copyright (c) 2021, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errors

import (
	stderrors "errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func provideErrors(withEmbedders bool) map[string]error {
	res := make(map[string]error, 10)
	errs := map[string]func() error{
		"std error": func() error { return stderrors.New("some err") },
		"error":     func() error { return New("whoopsie") },
		"Msg":       func() error { return Msg("my error message") },
		"Kind":      func() error { return Kind("my error kind") },
	}
	for a, errFn := range errs {
		res[a] = errFn()
		if withEmbedders {
			for b, embedFn := range provideEmbedders() {
				func() {
					defer func() { recover() }()
					res[a+"/"+b] = embedFn(errFn())
				}()
			}
		}
	}
	return res
}

func provideEmbedders() map[string]func(parent error) error {
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
			for name, embedFn := range provideEmbedders() {
				t.Run(name, func(t *testing.T) {
					have := embedFn(want)
					assert.Same(t, want, Unembed(have))
					assert.ErrorIs(t, have, want)
				})
			}
		})
	}
}

func TestEmbedError_Format(t *testing.T) {
	err := stderrors.New("foobar")
	for name, embedFn := range provideEmbedders() {
		t.Run(name, func(t *testing.T) {
			assert.Exactly(t, "foobar", fmt.Sprintf("%v", embedFn(err)))
		})
	}
}
