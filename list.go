package errs

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
func (l *List) All() []error { return l.list }

// Len returns the number of errors within the list.
func (l *List) Len() int { return len(l.list) }

// Append an error to the list. It guarantees only non-nil errors are added and
// returns `true`. It returns `false` when the error is nil.
func (l *List) Append(err error) bool {
	if err == nil {
		return false
	}

	l.Lock()
	l.list = append(l.list, err)
	l.Unlock()
	return true
}

// Prepend an error to the list. It guarantees only non-nil errors are added and
// returns `true`. It returns `false` when the error is nil.
func (l *List) Prepend(err error) bool {
	if err == nil {
		return false
	}

	l.Lock()
	l.list = append([]error{err}, l.list...)
	l.Unlock()
	return true
}

// Iter iterates over the errors in the list. Each error is sent over the
// returned channel. This way it is possible to iterate over the list using
// the build in `range` keyword.
func (l *List) Iter() <-chan error {
	ch := make(chan error)
	go func() {
		l.RLock()
		for _, v := range l.list {
			ch <- v
		}
		close(ch)
		defer l.RUnlock()
	}()

	return ch
}

func (l *List) Combine() error {
	l.RLock()
	err := combine(l.list)
	l.RUnlock()

	return err
}
