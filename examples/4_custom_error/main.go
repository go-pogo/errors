package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/roeldev/go-errs"
)

const SomeError errs.Kind = "some error"

//
// define custom error
//
type CustomErr struct {
	errs.Inner
	Value string
}

func (ce *CustomErr) Error() string {
	return fmt.Sprintf("just a custom error message with `%s`", ce.Value)
}

func (ce *CustomErr) Format(s fmt.State, v rune) { errs.FormatError(ce, s, v) }

func customErr(cause error) *CustomErr {
	err := &CustomErr{Inner: errs.MakeInner(cause, SomeError)}
	err.Frames().Capture(1)
	return err
}

//
// actual "program"
//
func unmarshal() (struct{}, error) {
	dest := struct{}{}
	err := json.Unmarshal([]byte("invalid"), &dest) // this wil result in an error
	return dest, errs.Trace(err)
}

func someAction() error {
	data, err := unmarshal()
	if err != nil {
		err := customErr(err)
		err.Value = "some important value"
		return err
	}

	// this code never runs
	fmt.Println(data)
	return nil
}

func main() {
	err := someAction()
	if err != nil {
		fmt.Printf("%v\n", err)
		fmt.Println("//////////")
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}
}
