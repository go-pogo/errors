// Copyright (c) 2020, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errors

import (
	stderrors "errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWrapWrapf(t *testing.T) {
	disableTraceStack()
	defer enableTraceStack()

	tests := map[string]struct {
		cause           error
		message, format string
		args            []interface{}
	}{
		"with std error": {
			cause:   stderrors.New("some err"),
			message: "foobar",
			format:  "%s",
			args:    []interface{}{"foobar"},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			wrap := Wrap(tc.cause, tc.message)
			assert.Exactly(t, tc.cause, Unwrap(wrap))

			wrapf := Wrapf(tc.cause, tc.format, tc.args...)
			assert.Exactly(t, wrap.Error(), wrapf.Error())
		})
	}

	t.Run("with nil cause", func(t *testing.T) {
		assert.Nil(t, Wrap(nil, "foobar"))
	})
}

func TestWrap(t *testing.T) {
	disableTraceStack()
	defer enableTraceStack()

	t.Run("with nil cause", func(t *testing.T) {
		assert.Nil(t, Wrap(nil, "foobar"))
	})

	cause := stderrors.New("the cause")
	str := "my error message"
	msg := Msg(str)

	tests := map[string]interface{}{
		"with Msg":     msg,
		"with *Msg":    &msg,
		"with string":  str,
		"with *string": &str,
	}
	for name, input := range tests {
		t.Run(name, func(t *testing.T) {
			have := Wrap(cause, input).(*commonError)
			assert.Equal(t, msg, have.error)
			assert.Same(t, cause, Unwrap(have))
			assert.Same(t, have, Unembed(have))
		})
	}

	t.Run("with error", func(t *testing.T) {
		assert.PanicsWithValue(t, panicUseWithKindInstead, func() {
			_ = Wrap(cause, Kind(str))
		})
	})

	tests = map[string]interface{}{
		"int":                 10,
		"bool":                false,
		"*errors.errorString": stderrors.New("not supported"),
	}

	t.Run("unsupported type", func(t *testing.T) {
		for typ, input := range tests {
			t.Run(typ, func(t *testing.T) {
				assert.PanicsWithValue(t,
					unsupportedType("errors.Wrap", typ),
					func() { _ = Wrap(cause, input) },
				)
			})
		}
	})
}

func TestWrapf(t *testing.T) {
	t.Run("with nil cause", func(t *testing.T) {
		assert.Nil(t, Wrapf(nil, "%s", "foobar"))
	})
}

func TestUnwrapAll(t *testing.T) {
	tests := map[string]func(chain *errChainHelper) error{
		"std error": func(chain *errChainHelper) error {
			return chain.prepend(stderrors.New("foo bar"))
		},
		"traced std error": func(chain *errChainHelper) error {
			err := chain.prepend(stderrors.New("bar: baz"))
			return chain.prepend(WithStack(err))
		},
		"std wrap": func(chain *errChainHelper) error {
			err := chain.prepend(stderrors.New("foo bar"))
			return chain.prepend(fmt.Errorf("cause: %w", err))
		},
		"traced std wrap": func(chain *errChainHelper) error {
			err := chain.prepend(stderrors.New("foo bar"))
			wrap := chain.prepend(fmt.Errorf("cause: %w", err))
			return chain.prepend(WithStack(wrap))
		},
		"error": func(chain *errChainHelper) error {
			return chain.prepend(New("err msg"))
		},
		"traced error": func(chain *errChainHelper) error {
			err := chain.prepend(New("err msg"))
			return WithStack(err)
		},
		"wrapped error": func(chain *errChainHelper) error {
			err := chain.prepend(New("qux"))
			return chain.prepend(Wrap(err, "bar msg"))
		},
	}

	for name, setupFn := range tests {
		t.Run(name, func(t *testing.T) {
			var chain errChainHelper
			err := setupFn(&chain)
			have := UnwrapAll(err)

			assert.Equal(t, len(chain), len(have))
			assert.Exactly(t, []error(chain), have)
		})
	}

	t.Run("nil", func(t *testing.T) {
		assert.Exactly(t, []error{}, UnwrapAll(nil))
	})
}

func TestCause(t *testing.T) {
	tests := map[string]struct {
		want  error
		setup func(e error) error
	}{
		"std error": {
			want:  stderrors.New("foo bar"),
			setup: func(e error) error { return e },
		},
		"traced std error": {
			want:  stderrors.New("foo bar"),
			setup: func(e error) error { return WithStack(e) },
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
				return WithStack(fmt.Errorf("cause: %w", e))
			},
		},
		"error": {
			want:  New("xoo"),
			setup: func(e error) error { return e },
		},
		"embedded error": {
			want:  New("xoo"),
			setup: func(e error) error { return WithExitCode(e, 1) },
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			have := Cause(tc.setup(tc.want))
			assert.Same(t, tc.want, have)
		})
	}
}
