// Copyright (c) 2020, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errors

import (
	stderrors "errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/go-pogo/errors/internal"
)

func TestWrapWrapf(t *testing.T) {
	internal.DisableCaptureFrames()
	defer internal.EnableCaptureFrames()

	tests := map[string]struct {
		message string
		format  string
		args    []interface{}
		cause   error
	}{
		"with std error": {
			message: "foobar",
			format:  "%s",
			args:    []interface{}{"foobar"},
			cause:   stderrors.New("some err"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			wrap := Wrap(tc.cause, tc.message)
			assert.Exactly(t, tc.cause, Unwrap(wrap))

			wrapf := Wrapf(tc.cause, tc.format, tc.args...)
			assert.Exactly(t, wrap, wrapf)
		})
	}

	t.Run("with nil cause", func(t *testing.T) {
		assert.Nil(t, Wrap(nil, "foobar"))
	})
}

func TestWrap(t *testing.T) {
	t.Run("with nil cause", func(t *testing.T) {
		assert.Nil(t, Wrap(nil, "foobar"))
	})
}

func TestWrapf(t *testing.T) {
	t.Run("with nil cause", func(t *testing.T) {
		assert.Nil(t, Wrapf(nil, "%s", "foobar"))
	})
}

func TestUnwrapAll(t *testing.T) {
	tests := map[string]func(chain *errChainHelper) error{
		"nil": func(chain *errChainHelper) error {
			return Trace(nil)
		},
		"std error": func(chain *errChainHelper) error {
			return chain.prepend(stderrors.New("foo bar"))
		},
		"traced std error": func(chain *errChainHelper) error {
			err := stderrors.New("bar: baz")
			return chain.prepend(Trace(err))
		},
		"double traced std error": func(chain *errChainHelper) error {
			err := stderrors.New("bar: baz")
			return chain.prepend(Trace(Trace(err)))
		},
		"std wrap": func(chain *errChainHelper) error {
			err := chain.prepend(stderrors.New("foo bar"))
			return chain.prepend(fmt.Errorf("cause: %w", err))
		},
		"traced std wrap": func(chain *errChainHelper) error {
			err := chain.prepend(stderrors.New("foo bar"))
			wrap := fmt.Errorf("cause: %w", err)
			return chain.prepend(Trace(wrap))
		},
		"error": func(chain *errChainHelper) error {
			return chain.prepend(New("err msg"))
		},
		"traced error": func(chain *errChainHelper) error {
			err := chain.prepend(New("err msg"))
			return Trace(err)
		},
		"double traced error": func(chain *errChainHelper) error {
			err := chain.prepend(New("err msg"))
			return Trace(Trace(err))
		},
		"wrapped error": func(chain *errChainHelper) error {
			err := chain.prepend(New("qux"))
			return chain.prepend(Wrap(err, "bar msg"))
		},
		"traced wrapped error": func(chain *errChainHelper) error {
			err := chain.prepend(New("qux"))
			err = chain.prepend(Wrap(err, "bar msg"))
			return Trace(err)
		},
	}

	for label, setupFn := range tests {
		t.Run(label, func(t *testing.T) {
			var chain errChainHelper
			err := setupFn(&chain)
			have := UnwrapAll(err)

			assert.Equal(t, len(chain), len(have))
			assert.Exactly(t, []error(chain), have)
		})
	}
}

func TestRootCause(t *testing.T) {
	tests := map[string]struct {
		want  error
		setup func(e error) error
	}{
		"std error": {
			want: stderrors.New("foo bar"),
			setup: func(e error) error {
				return e
			},
		},
		"traced std error": {
			want: stderrors.New("foo bar"),
			setup: func(e error) error {
				return Trace(e)
			},
		},
		"std wrap": {
			want: stderrors.New("foo bar"),
			setup: func(e error) error {
				return fmt.Errorf("%w", e)
			},
		},
		"traced std wrap": {
			want: stderrors.New("baz"),
			setup: func(e error) error {
				return Trace(fmt.Errorf("cause: %w", e))
			},
		},
		"error": {
			want: New("xoo"),
			setup: func(e error) error {
				return e
			},
		},
		"traced error": {
			want: New("xoo"),
			setup: func(e error) error {
				return Trace(e)
			},
		},
		"double traced error": {
			want: New("xoo"),
			setup: func(e error) error {
				return Trace(Trace(e))
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			have := Original(RootCause(tc.setup(tc.want)))
			assert.Same(t, tc.want, have)
		})
	}
}
