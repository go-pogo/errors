package errors

import (
	"testing"

	"github.com/roeldev/go-fail"
)

func TestFrame_IsEmpty(t *testing.T) {
	var frame Frame
	if !frame.IsEmpty() {
		t.Error(fail.Msg{
			Func: "errs.Frame.IsEmpty",
			Msg:  "should return true on empty frame",
		})
	}

	frame = Frame{}
	if !frame.IsEmpty() {
		t.Error(fail.Msg{
			Func: "errs.Frame.IsEmpty",
			Msg:  "should return true on empty frame",
		})
	}
}

func TestFrame_String(t *testing.T) {
	var tests = map[string]struct {
		subj Frame
		want string
	}{
		"empty": {
			subj: Frame{},
			want: "",
		},
		"filled": {
			subj: Frame{
				Path: "/test/file.go",
				Line: 123,
				Pkg:  "foo",
				Func: "Bar",
			},
			want: "/test/file.go:123: foo/Bar()",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			have := tc.subj.String()
			if have != tc.want {
				t.Error(fail.Diff{
					Func: "errs.Frame.String",
					Have: have,
					Want: tc.want,
				})
			}
		})
	}
}
