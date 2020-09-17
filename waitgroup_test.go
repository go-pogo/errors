package errors

import (
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/go-pogo/errors/internal"
)

func TestWaitGroup_Go(t *testing.T) {
	internal.DisableCaptureFrames()
	defer internal.EnableCaptureFrames()

	var wg WaitGroup
	var i int32

	wg.Go(func() error {
		atomic.AddInt32(&i, 1)
		return nil
	})
	wg.Go(func() error {
		atomic.AddInt32(&i, 1)
		return New(UnknownKind, UnknownError)
	})
	wg.Go(func() error {
		atomic.AddInt32(&i, 1)
		return nil
	})

	assert.Exactly(t, New(UnknownKind, UnknownError), wg.Wait())
	assert.Exactly(t, int32(3), i)
}

func TestWaitGroup_Wait(t *testing.T) {
	var wg WaitGroup
	t.Run("nil", func(t *testing.T) {
		wg.Go(func() error {
			return nil
		})

		assert.Nil(t, wg.Wait())
	})
	t.Run("error", func(t *testing.T) {
		internal.DisableCaptureFrames()
		defer internal.EnableCaptureFrames()

		wg.Go(func() error {
			return New(UnknownKind, UnknownError)
		})

		wantErr := New(UnknownKind, UnknownError)
		assert.Exactly(t, wantErr, wg.Wait())

		wantList := NewList(1)
		wantList.Append(wantErr)
		assert.Exactly(t, wantList, wg.ErrList())
	})
}
