// Copyright (c) 2020, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errlist

import (
	"sync"

	"github.com/go-pogo/errors"
)

var DefaultCapacity uint = 8

type ErrorLister interface {
	// ErrorList returns a List of collected non-nil errors.
	ErrorList() *List
}

// List is a thread-safe error list. Its zero value is ready to use.
type List struct {
	sync.RWMutex
	list []error
}

const (
	panicNewListCap  = "errors.New: cap cannot be below 0"
	panicNewListArgs = "errors.New: only one argument is allowed"
)

// New creates a new List with a pre-allocated capacity of cap.
func New(cap int) *List {
	return &List{
		list: make([]error, 0, cap),
	}
}

// All returns the error slice within the list.
func (l *List) All() []error {
	l.RLock()
	defer l.RUnlock()
	if l.list == nil {
		return nil
	}

	var res = make([]error, 0, len(l.list))
	copy(res, l.list)
	return res
}

// Len returns the number of errors within the List.
func (l *List) Len() int { return len(l.list) }

// Empty return true when the list is empty.
func (l *List) Empty() bool { return len(l.list) == 0 }

// Append an error to the list. It guarantees only non-nil errors are added.
// It returns false when a nil error is encountered. And true when the error
// is appended to the list.
func (l *List) Append(err error) bool {
	if err == nil {
		return false
	}

	l.Lock()
	if l.list == nil {
		l.list = make([]error, 0, DefaultCapacity)
	}
	l.list = append(l.list, err)
	l.Unlock()
	return true
}

// Prepend an error to the list. It guarantees only non-nil errors are added.
// It returns false when a nil error is encountered. And true when the error
// is prepended to the list.
func (l *List) Prepend(err error) bool {
	if err == nil {
		return false
	}

	l.Lock()
	if l.list == nil {
		l.list = make([]error, 0, DefaultCapacity)
		l.list = append(l.list, err)
	} else {
		l.list = prepend(l.list, err)
	}
	l.Unlock()
	return true
}

func prepend(errs []error, err error) []error {
	errs = append(errs, err)
	if len(errs) > 1 {
		copy(errs[1:], errs)
		errs[0] = err
	}
	return errs
}

// Join the collected errors. It uses the same rules and logic as the
// Join function.
func (l *List) Join() error {
	l.RLock()
	err := errors.Join(l.list...)
	l.RUnlock()
	return err
}
