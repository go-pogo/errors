package errs

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/roeldev/go-fail"
)

func TestUnwrapAll(t *testing.T) {
	tests := map[string]struct {
		err    error
		wantFn func(e error) []error
	}{
		"nil": {
			err: Trace(nil),
			wantFn: func(e error) []error {
				return nil
			},
		},
		"primitive error": {
			err: errors.New("foo bar"),
			wantFn: func(e error) []error {
				return []error{errors.New("foo bar")}
			},
		},
		"wrapped primitive": {
			err: Trace(errors.New("bar: baz")),
			wantFn: func(e error) []error {
				return []error{e}
			},
		},
		"double wrapped primitive": {
			err: Trace(Trace(errors.New("qux: xoo"))),
			wantFn: func(e error) []error {
				return []error{e}
			},
		},
		"primitive wrap": {
			err: fmt.Errorf("cause: %w", errors.New("foo bar")),
			wantFn: func(e error) []error {
				return []error{e, errors.New("foo bar")}
			},
		},
		// "wrapped primitive wrap": {
		// 	err: Trace(fmt.Errorf("cause: %w", errors.New("foo bar"))),
		// 	wantFn: func(e error) []error {
		// 		return []error{
		// 			errors.Unwrap(e),
		// 			errors.New("foo bar"),
		// 		}
		// 	},
		// },
		"error": {
			err: New("kind", "err msg"),
			wantFn: func(e error) []error {
				return []error{
					&err{Inner{
						frames: *GetFrames(e),
						kind:   "kind",
						msg:    "err msg",
					}},
				}
			},
		},
		"wrapped error": {
			err: Trace(New("kind", "err msg")),
			wantFn: func(e error) []error {
				return []error{
					&err{Inner{
						frames: *GetFrames(e),
						kind:   "kind",
						msg:    "err msg",
					}},
				}
			},
		},
		"double wrapped error": {
			err: Trace(Trace(New("kind", "err msg"))),
			wantFn: func(e error) []error {
				return []error{
					&err{Inner{
						frames: *GetFrames(e),
						kind:   "kind",
						msg:    "err msg",
					}},
				}
			},
		},
		"wrapped error error": {
			err: Wrap(New("baz", "qux"), "foo kind", "bar msg"),
			wantFn: func(e error) []error {
				cause := &err{
					Inner: Inner{
						frames: *GetFrames(errors.Unwrap(e)),
						kind:   "baz",
						msg:    "qux",
					},
				}

				return []error{
					&err{Inner{
						frames: *GetFrames(e),
						cause:  cause,
						kind:   "foo kind",
						msg:    "bar msg",
					}},
					cause,
				}
			},
		},
	}

	for label, tc := range tests {
		t.Run(label, func(t *testing.T) {
			have := UnwrapAll(tc.err)
			want := tc.wantFn(tc.err)

			if !reflect.DeepEqual(have, want) {
				t.Error(fail.Diff{
					Func: "UnwrapAll",
					Msg:  "should unwrap all and return a slice of errors",
					Have: have,
					Want: want,
					Opts: cmp.Options{
						cmp.AllowUnexported(traceErr{}),
					},
				})
			}
		})
	}
}

func TestUnwrapCause(t *testing.T) {
	tests := map[string]struct {
		err    error
		wantFn func(e error) error
	}{
		"primitive error": {
			err: errors.New("foo bar"),
			wantFn: func(e error) error {
				return errors.New("foo bar")
			},
		},
		"wrapped primitive error": {
			err: Trace(errors.New("foo bar")),
			wantFn: func(e error) error {
				return e
			},
		},
		"primitive wrap": {
			err: fmt.Errorf("%w", errors.New("foo bar")),
			wantFn: func(e error) error {
				return errors.New("foo bar")
			},
		},
		"wrapped primitive wrap": {
			err: Trace(fmt.Errorf("cause: %w", errors.New("baz"))),
			wantFn: func(e error) error {
				return errors.New("baz")
			},
		},
		"error": {
			err: New("qux", "xoo"),
			wantFn: func(e error) error {
				return &err{
					Inner: Inner{
						frames: *GetFrames(e),
						kind:   "qux",
						msg:    "xoo",
					},
				}
			},
		},
		"wrapped error": {
			err: Trace(New("qux", "xoo")),
			wantFn: func(e error) error {
				return &err{Inner{
					frames: *GetFrames(e),
					kind:   "qux",
					msg:    "xoo",
				}}
			},
		},
		"double wrapped error": {
			err: Trace(Trace(New("qux", "xoo"))),
			wantFn: func(e error) error {
				return &err{Inner{
					frames: *GetFrames(e),
					kind:   "qux",
					msg:    "xoo",
				}}
			},
		},
	}

	for label, tc := range tests {
		t.Run(label, func(t *testing.T) {
			have := UnwrapCause(tc.err)
			want := tc.wantFn(tc.err)

			if !reflect.DeepEqual(have, want) {
				t.Error(fail.Diff{
					Func: "UnwrapCause",
					Have: have,
					Want: want,
				})
			}
		})
	}
}
