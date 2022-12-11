// Copyright (c) 2020, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errors

import (
	"sync"
)

// A WaitGroup is a collection of goroutines working on subtasks that are part
// of the same overall task. It collects possible errors returned from the
// subtasks and, unlike golang.org/x/sync/errgroup.Group, does not cancel the
// group when an error is encountered.
type WaitGroup struct {
	errs List
	wg   sync.WaitGroup
}

// ErrorList returns a List of collected errors from the called goroutines.
func (g *WaitGroup) ErrorList() *List { return &g.errs }

// Wait blocks until all function calls from the Go method have returned, then
// returns all collected errors as a (multi) error.
func (g *WaitGroup) Wait() error {
	g.wg.Wait()
	return g.errs.Join()
}

// Go calls the given function in a new goroutine. Errors from all calls are
// collected, combined and returned by Wait.
func (g *WaitGroup) Go(fn func() error) {
	g.wg.Add(1)
	go func() {
		g.errs.Append(fn())
		g.wg.Done()
	}()
}
