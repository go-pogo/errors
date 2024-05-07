// Copyright (c) 2020, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errgroup

import (
	"errors"
	"sync/atomic"
	"testing"

	"github.com/go-pogo/errors/errlist"
	"github.com/go-pogo/errors/internal"
	"github.com/stretchr/testify/assert"
)

func TestWaitGroup_Go(t *testing.T) {
	internal.DisableTraceStack()
	defer internal.EnableTraceStack()

	want := errors.New("some err")
	var wg Group
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
	var wg Group
	t.Run("nil", func(t *testing.T) {
		wg.Go(func() error {
			return nil
		})

		assert.Nil(t, wg.Wait())
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
