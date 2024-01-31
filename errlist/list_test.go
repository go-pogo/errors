// Copyright (c) 2020, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errlist

import (
	stderrors "errors"
	"testing"

	"github.com/go-pogo/errors"
	"github.com/go-pogo/errors/internal"
	"github.com/stretchr/testify/assert"
)

func assertEmptyList(t *testing.T, list *List) {
	assert.Exactly(t, []error{}, list.All())
	assert.Exactly(t, 0, list.Len())
	assert.True(t, list.Empty())
}

func TestNewList(t *testing.T) {
	t.Run("default capacity", func(t *testing.T) {
		list := New()
		assert.Exactly(t, int(DefaultCapacity), cap(list.list))
		assertEmptyList(t, list)
	})
	t.Run("set capacity", func(t *testing.T) {
		list := New(22)
		assert.Equal(t, 22, cap(list.list))
		assertEmptyList(t, list)
	})
	t.Run("negative capacity", func(t *testing.T) {
		assert.PanicsWithValue(t, panicNewListCap, func() {
			_ = New(-1)
		})
	})
	t.Run("too much arguments", func(t *testing.T) {
		assert.PanicsWithValue(t, panicNewListArgs, func() {
			_ = New(1, 2)
		})
	})
}

func TestList_Append(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		list := New()
		assert.False(t, list.Append(nil))
		assertEmptyList(t, list)
	})
	t.Run("error", func(t *testing.T) {
		list := New()
		assert.True(t, list.Append(errors.New("some err")))
		assert.Equal(t, 1, list.Len())
	})
}

func TestList_Prepend(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		assert.False(t, New().Prepend(nil))
	})
	t.Run("error", func(t *testing.T) {
		assert.True(t, New().Prepend(errors.New("some err")))
	})
}

func TestList_All(t *testing.T) {
	t.Run("nil error", func(t *testing.T) {
		var list List
		assertEmptyList(t, &list)

		list.Append(nil)
		assertEmptyList(t, &list)
	})
	t.Run("error", func(t *testing.T) {
		var list List
		assertEmptyList(t, &list)

		err := errors.New("some err")
		list.Append(err)
		assert.Exactly(t, []error{err}, list.All())
	})
	t.Run("errors", func(t *testing.T) {
		var list List
		assertEmptyList(t, &list)

		err1 := errors.New("some err")
		err2 := stderrors.New("prepend me")

		list.Append(err1)
		list.Append(nil)
		list.Prepend(err2)

		assert.Exactly(t, []error{err2, err1}, list.All())
	})
}

func TestList_Join(t *testing.T) {
	internal.DisableTraceStack()
	defer internal.EnableTraceStack()

	errs := []error{
		errors.New("some err"),
		nil,
		stderrors.New("foobar"),
	}

	list := New(3)
	for _, e := range errs {
		list.Append(e)
	}

	multi := list.Join().(errors.MultiError)
	assert.Exactly(t, []error{errs[0], errs[2]}, multi.Unwrap())
	assert.Equal(t, errors.Join(errs...), multi)
}
