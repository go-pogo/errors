package errors

import (
	stderrors "errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/go-pogo/errors/internal"
)

func assertEmptyList(t *testing.T, list *List) {
	assert.Len(t, list.list, 0)
	assert.Exactly(t, []error{}, list.All())
}

func assertListLen(t *testing.T, list *List, length int) {
	assert.Len(t, list.list, length)
	assert.Exactly(t, len(list.All()), list.Len())
}

func TestNewList(t *testing.T) {
	t.Run("0 args", func(t *testing.T) {
		list := NewList()
		assert.Exactly(t, int(DefaultListCapacity), cap(list.list))
		assertEmptyList(t, list)
	})
	t.Run("1 arg", func(t *testing.T) {
		list := NewList(22)
		assert.Equal(t, 22, cap(list.list))
		assertEmptyList(t, list)
	})
	t.Run("2 args", func(t *testing.T) {
		assert.PanicsWithValue(t, panicNewListArgs, func() {
			_ = NewList(1, 2)
		})
	})
}

func TestList_Append(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		assert.False(t, NewList().Append(nil))
	})
	t.Run("error", func(t *testing.T) {
		assert.True(t, NewList().Append(New(UnknownKind, UnknownError)))
	})
}

func TestList_Prepend(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		assert.False(t, NewList().Prepend(nil))
	})
	t.Run("error", func(t *testing.T) {
		assert.True(t, NewList().Prepend(New(UnknownKind, UnknownError)))
	})
}

func TestList_All(t *testing.T) {
	t.Run("nil error", func(t *testing.T) {
		var list List
		assertEmptyList(t, &list)

		list.Append(nil)
		assertListLen(t, &list, 0)
	})
	t.Run("error", func(t *testing.T) {
		var list List
		assertEmptyList(t, &list)

		err := New(UnknownKind, UnknownError)
		list.Append(err)
		assert.Exactly(t, []error{err}, list.All())
	})
	t.Run("errors", func(t *testing.T) {
		var list List
		assertEmptyList(t, &list)

		err1 := New(UnknownKind, UnknownError)
		err2 := stderrors.New("prepend me")

		list.Append(err1)
		list.Append(nil)
		list.Prepend(err2)

		assert.Exactly(t, []error{err2, err1}, list.All())
	})
}

func TestList_Combine(t *testing.T) {
	internal.DisableCaptureFrames()
	defer internal.EnableCaptureFrames()

	errors := []error{
		New(UnknownKind, UnknownError),
		nil,
		stderrors.New("foobar"),
	}

	list := NewList()
	for _, e := range errors {
		list.Append(e)
	}

	multi := list.Combine().(MultiError)
	assert.Exactly(t, []error{errors[0], errors[2]}, multi.Errors())
	assert.Exactly(t, Combine(errors...), multi)
}
