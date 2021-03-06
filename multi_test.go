// Copyright (c) 2020, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errors

import (
	stderrors "errors"
	"fmt"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilter(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		have := Filter(nil)
		assert.Exactly(t, 0, len(have))
		assert.Exactly(t, 0, cap(have))
	})
	t.Run("empty", func(t *testing.T) {
		have := make([]error, 0)
		have = Filter(have)

		assert.Exactly(t, 0, len(have))
		assert.Exactly(t, 0, cap(have))
	})
	t.Run("with nils", func(t *testing.T) {
		input := []error{nil, nil}
		have := Filter(input)

		assert.Exactly(t, 0, len(have))
		assert.Exactly(t, 2, cap(have))
	})
	t.Run("with errors and nils", func(t *testing.T) {
		err1 := stderrors.New("some err")
		err2 := New("")

		input := []error{err1, nil, nil, err2, nil}
		have := Filter(input)

		assert.Equal(t, []error{err1, err2}, have)
		assert.Equal(t, []error{err1, err2, nil, nil, nil}, have[:cap(input)])
	})
}

func BenchmarkFilter(b *testing.B) {
	err1 := stderrors.New("some err")
	err2 := New("")

	tests := map[string]func(errors []error) []error{
		"filterV1": func(errors []error) []error {
			l := len(errors)
			if l == 0 {
				return errors
			}

			res := make([]error, 0, l)
			for _, err := range errors {
				if err != nil {
					res = append(res, err)
				}
			}
			return res
		},

		"filterV2": func(errors []error) []error {
			n := 0
			for i, err := range errors {
				if err == nil {
					continue
				}
				if i != n {
					errors[i] = nil
					errors[n] = err
				}
				n++
			}
			return errors[:n]
		},
	}

	// data sets to run the benchmarks with
	sets := [][]error{
		nil,
		{},
		{nil, nil},
		{err1, nil, nil, err2, nil},
		{err1, err2},
	}

	for name, fn := range tests {
		b.Run(name, func(b *testing.B) {
			b.StopTimer()
			b.ReportAllocs()

			for _, set := range sets {
				input := make([]error, len(set))
				copy(input, set)
				b.StartTimer()

				for n := 0; n < b.N; n++ {
					fn(input)
				}
				b.StopTimer()
			}
		})
	}
}

func TestCombine(t *testing.T) {
	t.Run("with empty and nil", func(t *testing.T) {
		assert.Nil(t, Combine())
		assert.Nil(t, Combine(nil))
	})
	t.Run("with errors", func(t *testing.T) {
		err1 := stderrors.New("first error")
		err2 := Newf("err with trace")
		multi := Combine(err1, err2).(*multiErr)

		assert.Exactly(t, []error{err1, err2}, multi.Errors())
	})
}

func TestAppend(t *testing.T) {
	t.Run("panic on nil dest ptr", func(t *testing.T) {
		assert.PanicsWithValue(t, panicAppendNilPtr, func() {
			Append(nil, New("bar"))
		})
	})
	t.Run("with nil", func(t *testing.T) {
		err := New("err")
		want := err.Error()
		Append(&err, nil)
		assert.Equal(t, want, err.Error())
	})
	t.Run("with error", func(t *testing.T) {
		var have error
		want := stderrors.New("foobar")
		Append(&have, want)

		assert.Same(t, want, have)
	})
	t.Run("with errors", func(t *testing.T) {
		var have error
		errs := []error{
			New("some err"),
			stderrors.New("whoops"),
			fmt.Errorf("another %s", "error"),
		}

		Append(&have, errs[0]) // set value to *have
		Append(&have, errs[1]) // create multi error from errors 0 and 1
		Append(&have, errs[2]) // append error 2 to multi error

		assert.IsType(t, new(multiErr), have)

		multi := have.(*multiErr)
		assert.Exactly(t, errs, multi.Errors())
		assert.Equal(t, len(multi.frames), 1)

		_, file, line, _ := runtime.Caller(0)
		assert.Contains(t, multi.StackFrames().String(), fmt.Sprintf("%s:%d", file, line-9))
	})
}
