package main

import (
	"fmt"

	"github.com/roeldev/go-errs"
)

//
// define custom error
//
type CustomError struct {
	errs.Inner
	Value string
}

func (ce CustomError) Message() string {
	return ce.Inner.Message() + ": " + ce.Value
}

func (ce CustomError) Error() string { return errs.Print(ce) }

func newCustomErr(msg string) *CustomError {
	return &CustomError{Inner: errs.MakeInnerWith(msg)}
}

//
// actual "program"
//
func doSomething() error {
	err := newCustomErr("something is wrong")
	err.Value = "some important value"

	return err
}

func main() {
	err := doSomething()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("All is well")
	}
}
