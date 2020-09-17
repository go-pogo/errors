package errors

import (
	"sync"
)

var DefaultListCapacity uint = 8

type ErrLister interface {
	// ErrList returns a List of collected non-nil errors that were encountered.
	ErrList() *List
}

type List struct {
	sync.RWMutex
	list []error
}

const panicNewListArgs = "errs.NewList: only one argument is allowed"

// NewList creates a new List with a pre-allocated capacity of `cap`.
func NewList(cap ...uint) *List {
	var c int
	switch len(cap) {
	case 0:
		c = int(DefaultListCapacity)
	case 1:
		c = int(cap[0])
	default:
		panic(panicNewListArgs)
	}

	return &List{
		list: make([]error, 0, c),
	}
}

// All returns the error slice within the list.
func (l *List) All() []error {
	if l.list == nil {
		l.list = make([]error, 0, 0)
	}
	return l.list
}

// Len returns the number of errors within the list.
func (l *List) Len() int { return len(l.list) }

// Append an error to the list. It guarantees only non-nil errors are added.
// It returns `false` when a nil error is encountered. And `true` when the error
// is appended to the list.
func (l *List) Append(err error) bool {
	if err == nil {
		return false
	}

	l.Lock()
	l.list = append(l.list, err)
	l.Unlock()
	return true
}

// Prepend an error to the list. It guarantees only non-nil errors are added.
// It returns `false` when a nil error is encountered. And `true` when the error
// is prepended to the list.
func (l *List) Prepend(err error) bool {
	if err == nil {
		return false
	}

	l.Lock()
	l.list = append([]error{err}, l.list...)
	l.Unlock()
	return true
}

// Combine the collected errors. It uses the same rules and logic as the
// `Combine` function.
func (l *List) Combine() error {
	l.RLock()
	err := combine(l.list)
	l.RUnlock()

	return err
}
