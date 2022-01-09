// Copyright (c) 2020, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errors

import (
	stderrors "errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/xerrors"
)

func BenchmarkNew(b *testing.B) {
	disableTraceStack()
	defer enableTraceStack()

	str := "some err"
	msg := Msg("some err")

	b.Run("Msg", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = New(msg)
		}
	})
	b.Run("*Msg", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = New(&msg)
		}
	})
	b.Run("string", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = New(str)
		}
	})
	b.Run("*string", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = New(&str)
		}
	})
	b.Run("string to Msg", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = New(Msg(str))
		}
	})
}

func TestNew(t *testing.T) {
	disableTraceStack()
	defer enableTraceStack()

	t.Run("with nil", func(t *testing.T) {
		assert.Nil(t, New(nil))
	})

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
			have := New(input).(*commonError)
			assert.Equal(t, msg, have.error)
			assert.Nil(t, have.cause)
			assert.Nil(t, Unwrap(have))
			assert.Same(t, have, Unembed(have))
		})
	}

	t.Run("with error", func(t *testing.T) {
		assert.PanicsWithValue(t, panicUseWithStackInstead, func() {
			_ = New(stderrors.New(str))
		})
	})

	tests = map[string]interface{}{
		"int":  10,
		"bool": false,
	}

	t.Run("unsupported type", func(t *testing.T) {
		for typ, input := range tests {
			t.Run(typ, func(t *testing.T) {
				assert.PanicsWithValue(t,
					UnsupportedTypeError{Func: "errors.New", Type: typ},
					func() { _ = New(input) },
				)
			})
		}
	})
}

func TestNewf(t *testing.T) {
	disableTraceStack()
	defer enableTraceStack()

	t.Run("without args", func(t *testing.T) {
		assert.Equal(t, New("some err"), Newf("some err"))
	})
	t.Run("with cause", func(t *testing.T) {
		cause := stderrors.New("some err")
		have := Newf("whoops: %w", cause).(*commonError)
		assert.ErrorIs(t, have, cause)
		assert.Equal(t, cause, have.cause)
		assert.Equal(t, cause, Unwrap(have))
		assert.Same(t, have, Unembed(have))
	})
}

func TestMsg(t *testing.T) {
	msg := Msg("some msg")
	assert.Equal(t, msg.String(), msg.Error())
}

func TestMsg_Is(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		msg := Msg("some err")
		tests := map[string]error{
			"Msg":  Msg("some err"),
			"*Msg": &msg,
		}
		for a, err := range tests {
			for b, target := range tests {
				t.Run(a+"/"+b, func(t *testing.T) {
					assert.ErrorIs(t, err, target)
				})
			}
		}
	})

	t.Run("false", func(t *testing.T) {
		msg := "some err"
		targets := map[string]error{
			"different msg string": Msg("blabla"),
			"stderror":             stderrors.New(msg),
			"error":                New(msg),
			"Kind":                 Kind(msg),
		}
		for name, target := range targets {
			t.Run(name, func(t *testing.T) {
				assert.NotErrorIs(t, Msg(msg), target)
			})
		}
	})
}

func TestMsg_As(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		var msg Msg
		assert.True(t, Msg("hi there").As(&msg))
		assert.Exactly(t, Msg("hi there"), msg)
	})
	t.Run("false", func(t *testing.T) {
		var msg Msg
		assert.False(t, Msg("hi there").As(msg))
		assert.Exactly(t, Msg(""), msg)
	})
}

func TestSameErrors(t *testing.T) {
	disableTraceStack()
	defer enableTraceStack()

	cause := xerrors.New("cause of error")
	tests := map[string]map[string][2]error{
		"New&Newf": {
			"empty": {New(""), Newf("")},
			"message only": {
				New("some `foo` happened"),
				Newf("some `%s` happened", "foo"),
			},
		},
		"Wrap&Wrapf": {
			"empty": {
				Wrap(cause, ""),
				Wrapf(cause, ""),
			},
			"message only": {
				Wrap(cause, "some `foo` happened"),
				Wrapf(cause, "some `%s` happened", "foo"),
			},
		},
	}

	for group, ts := range tests {
		t.Run(group, func(t *testing.T) {
			for name, errs := range ts {
				t.Run(name, func(t *testing.T) {
					assert.Equal(t, errs[0].Error(), errs[1].Error())
				})
			}
		})
	}
}
