package errs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewList(t *testing.T) {
	t.Run("0 args", func(t *testing.T) {
		list := NewList()
		assert.Exactly(t, int(DefaultListCapacity), cap(list.list))
		assert.Len(t, list.list, 0)
	})
	t.Run("1 arg", func(t *testing.T) {
		list := NewList(22)
		assert.Equal(t, 22, cap(list.list))
		assert.Len(t, list.list, 0)
	})
	t.Run("2 args", func(t *testing.T) {
		assert.PanicsWithValue(t, panicNewListArgs, func() {
			_ = NewList(1, 2)
		})
	})
}

func assertListLen(t *testing.T, list *List, length int) {
	assert.Len(t, list.list, length)
	assert.Exactly(t, len(list.list), list.Len())
}

func TestList_Append(t *testing.T) {
	t.Run("nil error", func(t *testing.T) {
		list := NewList()
		assertListLen(t, list, 0)

		assert.False(t, list.Append(nil))
		assertListLen(t, list, 0)
	})
	t.Run("error", func(t *testing.T) {
		list := NewList()
		assertListLen(t, list, 0)

		err := New(UnknownKind, UnknownError)
		assert.True(t, list.Append(err))
		assertListLen(t, list, 1)
		assert.Exactly(t, []error{err}, list.All())
	})
	t.Run("errors", func(t *testing.T) {
		list := NewList()
		assertListLen(t, list, 0)
	})
}
