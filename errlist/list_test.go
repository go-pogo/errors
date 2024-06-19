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
	assert.Exactly(t, 0, list.Len())
	assert.True(t, list.IsEmpty())
}

func TestNew(t *testing.T) {
	list := New(make([]error, 0, DefaultCapacity))
	assert.Exactly(t, int(DefaultCapacity), cap(list.list))
	assertEmptyList(t, list)
}

func TestNewWithCapacity(t *testing.T) {
	list := NewWithCapacity(22)
	assert.Equal(t, 22, cap(list.list))
	assertEmptyList(t, list)
}

func TestList(t *testing.T) {
	t.Run("zero value", func(t *testing.T) {
		var list List
		assert.Nil(t, list.list)
		assertEmptyList(t, &list)
	})
}

func TestList_Append(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var list List
		assert.False(t, list.Append(nil))
		assertEmptyList(t, &list)
	})
	t.Run("error", func(t *testing.T) {
		var list List
		assert.True(t, list.Append(errors.New("some err")))
		assert.Equal(t, 1, list.Len())
	})
}

func TestList_Prepend(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var list List
		assert.False(t, list.Prepend(nil))
		assertEmptyList(t, &list)
	})
	t.Run("error", func(t *testing.T) {
		var list List
		assert.True(t, list.Prepend(errors.New("some err")))
		assert.Equal(t, 1, list.Len())
	})
}

func TestList_All(t *testing.T) {
	t.Run("nil list", func(t *testing.T) {
		var list List
		assert.Nil(t, list.All())
	})
	t.Run("nil error", func(t *testing.T) {
		var list List
		assertEmptyList(t, &list)

		list.Append(nil)
		assertEmptyList(t, &list)
	})
	t.Run("error", func(t *testing.T) {
		var list List
		err := errors.New("some err")
		list.Append(err)
		assert.Exactly(t, []error{err}, list.All())
	})
	t.Run("errors", func(t *testing.T) {
		err1 := errors.New("some err")
		err2 := stderrors.New("prepend me")

		var list List
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

	var list List
	for _, e := range errs {
		list.Append(e)
	}

	//goland:noinspection GoTypeAssertionOnErrors
	multi := list.Join().(errors.MultiError)
	assert.Exactly(t, []error{errs[0], errs[2]}, multi.Unwrap())
	assert.Equal(t, errors.Join(errs...), multi)
}
