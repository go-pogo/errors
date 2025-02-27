// Copyright (c) 2020, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errgroup

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/go-pogo/errors"
	"github.com/go-pogo/errors/errlist"
	"github.com/go-pogo/errors/internal"
	"github.com/stretchr/testify/assert"
)

func TestGroup(t *testing.T) {
	internal.DisableTraceStack()
	defer internal.EnableTraceStack()

	someErr := errors.New("some err")
	canceledCtx, cancel := context.WithCancel(context.Background())
	cancel()

	t.Run("unique error", func(t *testing.T) {
		wg, ctx := WithContext(canceledCtx)
		wg.Go(func() error { return ctx.Err() })
		wg.Go(func() error { return ctx.Err() })

		err := wg.Wait()
		assert.Equal(t, 1, wg.ErrorList().Len())
		assert.ErrorIs(t, err, context.Canceled)
	})
	t.Run("unique and wrapped error", func(t *testing.T) {
		wg, ctx := WithContext(canceledCtx)
		wg.Go(func() error { return ctx.Err() })
		wg.Go(func() error {
			return errors.Wrap(ctx.Err(), "wrapped")
		})
		wg.Go(func() error { return ctx.Err() })

		_ = wg.Wait()
		assert.Equal(t, 2, wg.ErrorList().Len())
	})

	t.Run("run all goroutines", func(t *testing.T) {
		var wg Group
		var i int32

		for n := 3; n > 0; n-- {
			wg.Go(func() error {
				atomic.AddInt32(&i, 1)
				return someErr
			})
		}

		assert.Same(t, someErr, wg.Wait())
		assert.Exactly(t, int32(3), i)
	})
	t.Run("return on first error", func(t *testing.T) {
		wg, ctx := WithContext(context.Background())
		var i atomic.Int32

		for n := 3; n > 0; n-- {
			wg.Go(func() error {
				time.Sleep(time.Duration(i.Load()) * time.Second)
				if ctx.Err() == nil {
					i.Add(1)
					return someErr
				}
				return nil
			})
		}

		haveErr := wg.Wait()
		assert.Same(t, someErr, haveErr)
		assert.Same(t, context.Cause(ctx), haveErr)
		assert.NotEqual(t, int32(3), i.Load())
	})
}

func TestWaitGroup_ErrorList(t *testing.T) {
	var wg Group
	t.Run("nil", func(t *testing.T) {
		wg.Go(func() error {
			return nil
		})

		assert.Nil(t, wg.Wait())
		assert.Exactly(t, 0, wg.ErrorList().Len())
	})
	t.Run("error", func(t *testing.T) {
		internal.DisableTraceStack()
		defer internal.EnableTraceStack()

		wantErr := errors.New("some err")
		wg.Go(func() error {
			return wantErr
		})

		assert.Same(t, wantErr, wg.Wait())

		wantList := errlist.NewWithCapacity(1)
		wantList.Append(wantErr)
		assert.Exactly(t, wantList, wg.ErrorList())
	})
}
