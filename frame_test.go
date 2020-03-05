package errs

import (
	"testing"

	"github.com/roeldev/go-fail"
)

func TestFrame_IsEmpty(t *testing.T) {
	var frame Frame
	if !frame.IsEmpty() {
		t.Error(fail.Msg{
			Func: "Frame.IsEmpty",
			Msg:  "should return true on empty frame",
		})
	}

	frame = Frame{}
	if !frame.IsEmpty() {
		t.Error(fail.Msg{
			Func: "Frame.IsEmpty",
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
				Func: "foo.Bar",
			},
			want: "/test/file.go:123: foo.Bar()",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			have := tc.subj.String()
			if have != tc.want {
				t.Error(fail.Diff{
					Func: "Frame.String",
					Have: have,
					Want: tc.want,
				})
			}
		})
	}
}

func TestGetFrame(t *testing.T) {
	frame, ok := GetFrame(0)
	if frame.IsEmpty() || !ok {
		t.Error(fail.RetVal{
			Func: "GetFrame",
			Msg:  "should return a Frame with ok = true",
			Have: []interface{}{frame, ok},
			Want: []interface{}{Frame{
				Path: frame.Path,
				Line: 60,
				Func: "github.com/roeldev/go-errs.TestGetFrame",
			}, true},
		})
	}
}

func TestGetFrame__invalid_skip(t *testing.T) {
	frame, ok := GetFrame(9999)
	if !frame.IsEmpty() || ok {
		t.Error(fail.RetVal{
			Func: "GetFrame",
			Msg:  "should return an empty Frame with ok = false on invalid skip value",
			Have: []interface{}{frame, ok},
			Want: []interface{}{Frame{}, false},
		})
	}
}
