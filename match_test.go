// Copyright (c) 2020, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errors

import (
	stderrors "errors"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
)

func assertErrorIs(t *testing.T, err, target error) bool {
	return assert.True(t, Is(err, target), fmt.Sprintf("error %T should match with target %T", err, target))
}

type errChainHelper []error

func (h *errChainHelper) append(err error) error {
	*h = append(*h, err)
	return err
}

func (h *errChainHelper) prepend(err error) error {
	*h = prepend(*h, err)
	return err
}

func (h errChainHelper) last() error {
	return h[len(h)-1]
}

// test if the root cause error matches all wrapping errors in the chain
func TestIs(t *testing.T) {
	rootCause := New("root cause")
	stdRootCause := stderrors.New("root cause")

	chains := map[string]func(chain *errChainHelper){
		"base": func(chain *errChainHelper) {
			_ = chain.append(rootCause)
		},
		"std": func(chain *errChainHelper) {
			_ = chain.append(stdRootCause)
		},
		"traced base": func(chain *errChainHelper) {
			err := chain.append(rootCause)
			_ = chain.append(Trace(err))
		},
		"traced std": func(chain *errChainHelper) {
			err := chain.append(stdRootCause)
			_ = chain.append(Trace(err))
		},
		"base with kind": func(chain *errChainHelper) {
			err := chain.append(rootCause)
			_ = chain.append(WithKind(err, "kind"))
		},
		"std with kind": func(chain *errChainHelper) {
			err := chain.append(stdRootCause)
			_ = chain.append(WithKind(err, "kind"))
		},
	}

	tests := map[string]func(parent error) error{
		"WithKind": func(parent error) error {
			return WithKind(parent, "some kind")
		},
		"WithExitCode": func(parent error) error {
			return WithExitCode(parent, 1)
		},
		"WithFormatter": func(parent error) error {
			return WithFormatter(parent)
		},
		"Trace": func(parent error) error {
			return Trace(parent)
		},
	}

	for group, upgradeFn := range tests {
		t.Run(group, func(t *testing.T) {
			for name, setupFn := range chains {
				t.Run(name, func(t *testing.T) {
					var chain errChainHelper
					setupFn(&chain)

					// pass the last error to the upgrade function we'd like
					// to test
					err := upgradeFn(chain.last())
					_ = chain.append(err)

					assert.Same(t, chain[0], Original(RootCause(err)))

					for i, target := range chain {
						for j := i; j < len(chain); j++ {
							err = chain[j]
							assertErrorIs(t, err, target)
						}
					}
				})
			}
		})
	}

	t.Run("manual", func(t *testing.T) {
		rootCause := stderrors.New("root cause")
		withKind := WithKind(rootCause, "some kind")
		withFormatter := WithFormatter(withKind)

		// both upgrades should match with the original error
		assertErrorIs(t, withKind, rootCause)
		assertErrorIs(t, withFormatter, rootCause)

		// both upgrades should match with each other
		assertErrorIs(t, withKind, withFormatter)
		assertErrorIs(t, withFormatter, withKind)
	})
}

type customError struct{}

func (ce *customError) Error() string { return "this is a custom error" }

func TestAs(t *testing.T) {
	disableCaptureFrames()
	defer enableCaptureFrames()

	var kinder Kinder
	var customErr *customError
	var pathErrPtr *os.PathError
	_, pathErr := os.Open("non-existing")

	tests := map[string]struct {
		error  error
		target interface{}
		wantFn func(err error) interface{}
	}{
		"commonErr with kind": {
			error:  WithKind(New("err with kind"), "foobar"),
			target: &kinder,
			wantFn: func(err error) interface{} {
				ce := toCommonErr(stderrors.New("err with kind"), false)
				ce.kind = "foobar"
				return ce
			},
		},
		"upgraded os.PathError": {
			error:  Upgrade(pathErr),
			target: &pathErrPtr,
			wantFn: func(err error) interface{} {
				return pathErr
			},
		},
		"upgraded custom error": {
			error:  Upgrade(&customError{}),
			target: &customErr,
			wantFn: func(err error) interface{} {
				return &customError{}
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			val := reflect.ValueOf(tc.target)
			val.Elem().Set(reflect.Zero(reflect.TypeOf(tc.target).Elem()))

			assert.True(t, As(tc.error, tc.target))

			got := val.Elem().Interface()
			want := tc.wantFn(tc.error)

			assert.Equal(t, want, got, cmp.Diff(got, want, cmpopts.EquateErrors()))
		})
	}
}
