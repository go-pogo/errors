package errors

import (
	"reflect"
	"testing"

	"github.com/roeldev/go-fail"
)

// func TestPrint(t *testing.T) {
// 	err1 := Err("foo error")
// 	err2 := Wrap(err1)
// 	err3 := Wrap(err2)
//
// 	have := err3.Error()
// 	want := `foo error
//
// Trace:
// error_test.go:13: errors.TestPrint():
// error_test.go:12: errors.TestPrint():
// error_test.go:11: errors.TestPrint():
// >	foo error
// `
//
// 	if have != want {
// 		t.Error(fail.Diff{
// 			Func: "errs.Print",
// 			Msg:  "should create the same error output message",
// 			Have: have,
// 			Want: want,
// 		})
// 	}
// }

func TestUnwrapAll(t *testing.T) {
	err1 := Err("bar error")
	err2 := Wrap(err1)
	err3 := Wrap(err2)

	have := UnwrapAll(err3)
	want := []error{err1}

	if !reflect.DeepEqual(have, want) {
		t.Error(fail.Diff{
			Func: "UnwrapAll",
			Msg:  "should return a slice of errors",
			Have: have,
			Want: want,
		})
	}
}
