// Copyright (c) 2020, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errors

import (
	stderrors "errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func assertEmptyList(t *testing.T, list *List) {
	assert.Exactly(t, []error{}, list.All())
	assert.Exactly(t, 0, list.Len())
	assert.True(t, list.Empty())
}

func TestNewList(t *testing.T) {
	t.Run("default capacity", func(t *testing.T) {
		list := NewList()
		assert.Exactly(t, int(DefaultListCapacity), cap(list.list))
		assertEmptyList(t, list)
	})
	t.Run("set capacity", func(t *testing.T) {
		list := NewList(22)
		assert.Equal(t, 22, cap(list.list))
		assertEmptyList(t, list)
	})
	t.Run("negative capacity", func(t *testing.T) {
		assert.PanicsWithValue(t, panicNewListCap, func() {
			_ = NewList(-1)
		})
	})
	t.Run("too much arguments", func(t *testing.T) {
		assert.PanicsWithValue(t, panicNewListArgs, func() {
			_ = NewList(1, 2)
		})
	})
}

func TestList_Append(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		list := NewList()
		assert.False(t, list.Append(nil))
		assertEmptyList(t, list)
	})
	t.Run("error", func(t *testing.T) {
		list := NewList()
		assert.True(t, list.Append(New("some err")))
		assert.Equal(t, 1, list.Len())
	})
}

func TestList_Prepend(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		assert.False(t, NewList().Prepend(nil))
	})
	t.Run("error", func(t *testing.T) {
		assert.True(t, NewList().Prepend(New("some err")))
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

		err := New("some err")
		list.Append(err)
		assert.Exactly(t, []error{err}, list.All())
	})
	t.Run("errors", func(t *testing.T) {
		var list List
		assertEmptyList(t, &list)

		err1 := New("some err")
		err2 := stderrors.New("prepend me")

		list.Append(err1)
		list.Append(nil)
		list.Prepend(err2)

		assert.Exactly(t, []error{err2, err1}, list.All())
	})
}

func TestList_Join(t *testing.T) {
	disableTraceStack()
	defer enableTraceStack()

	errs := []error{
		New("some err"),
		nil,
		stderrors.New("foobar"),
	}

	list := NewList(3)
	for _, e := range errs {
		list.Append(e)
	}

	multi := list.Join().(MultiError)
	assert.Exactly(t, []error{errs[0], errs[2]}, multi.Unwrap())
	assert.Equal(t, Join(errs...), multi)
}
