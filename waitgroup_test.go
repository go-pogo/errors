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

	want := New("some err")
	var wg WaitGroup
	var i int32

	wg.Go(func() error {
		atomic.AddInt32(&i, 1)
		return nil
	})
	wg.Go(func() error {
		atomic.AddInt32(&i, 1)
		return want
	})
	wg.Go(func() error {
		atomic.AddInt32(&i, 1)
		return nil
	})

	have := wg.Wait()
	assert.Exactly(t, int32(3), i)
	assert.Same(t, want, have)
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

		wantErr := New("some err")
		wg.Go(func() error {
			return wantErr
		})

		assert.Same(t, wantErr, wg.Wait())

		wantList := NewList(1)
		wantList.Append(wantErr)
		assert.Exactly(t, wantList, wg.ErrorList())
	})
}
