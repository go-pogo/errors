// Copyright (c) 2020, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errgroup

import (
	"context"
	"github.com/go-pogo/errors"
	"github.com/go-pogo/errors/errlist"
	"sync"
)

var _ errlist.ErrorLister = (*Group)(nil)

// A Group is a collection of goroutines working on subtasks that are part
// of the same overall task. It collects possible errors returned from the
// subtasks.
type Group struct {
	cancel func(error)
	wg     sync.WaitGroup
	errs   errlist.List
}

// WithContext returns a new Group and an associated context.Context derived
// from ctx. Unlike golang.org/x/sync/errgroup.Group, the derived context is not
// canceled the first time a function passed to Group.Go returns a non-nil
// error.
func WithContext(ctx context.Context) (*Group, context.Context) {
	ctx, cancel := context.WithCancelCause(ctx)
	return &Group{cancel: cancel}, ctx
}

// ErrorList returns a List of collected errors from the called goroutines.
func (g *Group) ErrorList() *errlist.List { return &g.errs }

// Wait blocks until all function calls from the Go method have returned, then
// returns all collected errors as a (multi) error.
func (g *Group) Wait() error {
	g.wg.Wait()

	err := g.errs.Join()
	if g.cancel != nil {
		g.cancel(err)
	}
	return err
}

// Go calls the given function in a new goroutine. Errors from all calls are
// collected, combined and returned by Wait.
func (g *Group) Go(fn func() error) {
	g.wg.Add(1)
	go func() {
		if err := fn(); err != nil {
			if errors.IsCause(err) {
				g.errs.AppendUnique(err)
			} else {
				g.errs.Append(err)
			}
		}
		g.wg.Done()
	}()
}
