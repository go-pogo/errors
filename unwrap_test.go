package errs

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type unwrapAllHelper []error

func (h *unwrapAllHelper) add(err error) error {
	x := append(*h, err)
	if len(x) > 1 {
		copy(x[1:], x)
		x[0] = err
	}
	*h = x
	return err
}

func TestUnwrapAll(t *testing.T) {
	tests := map[string]func(want *unwrapAllHelper) error{
		"nil": func(want *unwrapAllHelper) error {
			return Trace(nil)
		},
		"primitive error": func(want *unwrapAllHelper) error {
			return errors.New("foo bar")
		},
		"traced primitive": func(want *unwrapAllHelper) error {
			err := want.add(errors.New("bar: baz"))
			return Trace(err)
		},
		"double traced primitive": func(want *unwrapAllHelper) error {
			err := want.add(errors.New("qux: xoo"))
			return Trace(Trace(err))
		},
		"primitive wrap": func(want *unwrapAllHelper) error {
			err := want.add(errors.New("foo bar"))
			return fmt.Errorf("cause: %w", err)
		},
		"traced primitive wrap": func(want *unwrapAllHelper) error {
			err := want.add(errors.New("foo bar"))
			err = want.add(fmt.Errorf("cause: %w", err))
			return Trace(err)
		},
		"error": func(want *unwrapAllHelper) error {
			return New("kind", "err msg")
		},
		"traced error": func(want *unwrapAllHelper) error {
			return Trace(New("kind", "err msg"))
		},
		"double traced error": func(want *unwrapAllHelper) error {
			return Trace(Trace(New("kind", "err msg")))
		},
		"wrapped error error": func(want *unwrapAllHelper) error {
			err := want.add(New("baz", "qux"))
			return Wrap(err, "foo kind", "bar msg")
		},
	}

	for label, setup := range tests {
		t.Run(label, func(t *testing.T) {
			var h unwrapAllHelper
			err := setup(&h)
			if err != nil {
				_ = h.add(err)
			}

			assert.Exactly(t, []error(h), UnwrapAll(err))
		})
	}
}

func TestUnwrapCause(t *testing.T) {
	tests := map[string]struct {
		want  error
		setup func(e error) error
	}{
		"primitive error": {
			want: errors.New("foo bar"),
			setup: func(e error) error {
				return e
			},
		},
		"traced primitive error": {
			want: errors.New("foo bar"),
			setup: func(e error) error {
				return Trace(e)
			},
		},
		"primitive wrap": {
			want: errors.New("foo bar"),
			setup: func(e error) error {
				return fmt.Errorf("%w", e)
			},
		},
		"traced primitive wrap": {
			want: errors.New("baz"),
			setup: func(e error) error {
				return Trace(fmt.Errorf("cause: %w", e))
			},
		},
		"error": {
			want: New("qux", "xoo"),
			setup: func(e error) error {
				return e
			},
		},
		"traced error": {
			want: New("qux", "xoo"),
			setup: func(e error) error {
				return Trace(e)
			},
		},
		"double traced error": {
			want: New("qux", "xoo"),
			setup: func(e error) error {
				return Trace(Trace(e))
			},
		},
	}

	for label, tc := range tests {
		t.Run(label, func(t *testing.T) {
			have := UnwrapCause(tc.setup(tc.want))
			assert.Same(t, tc.want, have)
		})
	}
}
