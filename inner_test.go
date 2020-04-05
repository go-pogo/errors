package errs

import (
	"reflect"
	"testing"

	"github.com/roeldev/go-fail"
)

func TestKind_String(t *testing.T) {
	want := "foo"
	have := Kind(want).String()
	if have != want {
		t.Error(fail.Diff{
			Func: "Kind.String",
			Msg:  "should return the string representation of the Kind",
			Have: have,
			Want: want,
		})
	}
}

func TestMakeInnerWith(t *testing.T) {
	tests := map[string]struct {
		args []interface{}
		want Inner
	}{
		"none": {
			args: []interface{}{},
			want: MakeInner(nil, UnknownKind, ""),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			have := MakeInnerWith(tc.args...)
			have.frames = tc.want.frames

			if !reflect.DeepEqual(have, tc.want) {
				t.Error(fail.Diff{
					Func: "MakeInnerWith",
					Msg:  "",
					Have: have,
					Want: tc.want,
				})
			}
		})
	}
}

func TestMakeInnerWith_panic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error(fail.Msg{
				Func: "MakeInnerWith",
				Msg:  "should panic on invalid argument type",
			})
		}
	}()

	_ = MakeInnerWith(true)
}
