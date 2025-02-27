// Copyright (c) 2020, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errlist

import (
	"sync"

	"github.com/go-pogo/errors"
)

// DefaultCapacity of internal errors slice when using zero value List.
var DefaultCapacity uint = 8

type ErrorLister interface {
	// ErrorList returns a [List] of collected non-nil errors.
	ErrorList() *List
}

// List is a thread-safe error list. Its zero value is ready to use.
type List struct {
	mut  sync.RWMutex
	list []error
}

// New creates a new [List] using the provided slice.
func New(slice []error) *List {
	return &List{list: slice}
}

// NewWithCapacity creates a new [List] with a pre-allocated capacity.
func NewWithCapacity(cap uint) *List {
	return &List{list: make([]error, 0, cap)}
}

// Len returns the number of errors within the [List].
func (l *List) Len() int {
	l.mut.RLock()
	defer l.mut.RUnlock()
	return len(l.list)
}

// IsEmpty return true when [List] is empty.
func (l *List) IsEmpty() bool {
	l.mut.RLock()
	defer l.mut.RUnlock()
	return len(l.list) == 0
}

// All returns a copy of the error slice within [List].
func (l *List) All() []error {
	l.mut.RLock()
	defer l.mut.RUnlock()
	if l.list == nil {
		return nil
	}

	res := make([]error, 0, len(l.list))
	res = append(res, l.list...)
	return res
}

// Join the collected errors. It uses the same rules and logic as the
// [errors.Join] function.
func (l *List) Join() error {
	l.mut.RLock()
	defer l.mut.RUnlock()
	return errors.Join(l.list...)
}

// Append an error to the [List]. It guarantees only non-nil errors are added.
// It returns true when the error is appended to [List], false otherwise.
func (l *List) Append(err error) bool {
	if err == nil {
		return false
	}

	l.mut.Lock()
	defer l.mut.Unlock()

	if l.list == nil {
		l.list = make([]error, 0, DefaultCapacity)
	}
	l.list = append(l.list, err)
	return true
}

// AppendUnique appends an error to [List] and guarantees that the error is
// non-nil and unique within the [List].
// It returns true when the error is appended, false otherwise.
func (l *List) AppendUnique(err error) bool {
	if err == nil {
		return false
	}

	l.mut.Lock()
	defer l.mut.Unlock()

	if l.list == nil {
		l.list = make([]error, 0, DefaultCapacity)
		l.list = append(l.list, err)
		return true
	}
	if l.isUnique(err) {
		l.list = append(l.list, err)
		return true
	}
	return false
}

// Prepend an error to the [List]. It guarantees only non-nil errors are added.
// It returns true when the error is appended, false otherwise.
func (l *List) Prepend(err error) bool {
	if err == nil {
		return false
	}

	l.mut.Lock()
	defer l.mut.Unlock()

	if l.list == nil {
		l.list = make([]error, 0, DefaultCapacity)
		l.list = append(l.list, err)
	} else {
		l.list = prepend(l.list, err)
	}
	return true
}

// PrependUnique prepends an error to [List] and guarantees that the error is
// non-nil and unique within the [List].
// It returns true when the error is appended, false otherwise.
func (l *List) PrependUnique(err error) bool {
	if err == nil {
		return false
	}

	l.mut.Lock()
	defer l.mut.Unlock()

	if l.list == nil {
		l.list = make([]error, 0, DefaultCapacity)
		l.list = append(l.list, err)
		return true
	}
	if l.isUnique(err) {
		l.list = prepend(l.list, err)
		return true
	}
	return false
}

func (l *List) isUnique(err error) bool {
	matchCause := errors.IsCause(err)
	for _, e := range l.list {
		if errors.IsCause(e) == matchCause && errors.Is(err, e) {
			return false
		}
	}
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
