package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/go-pogo/errors"
)

//
// define custom error
//
type customErr struct {
	Cause error
	Value string
}

func (ce *customErr) Error() string {
	return fmt.Sprintf("just a custom error message with `%s`", ce.Value)
}

func (ce *customErr) Format(s fmt.State, v rune) { errors.FormatError(ce, s, v) }

func newCustomErr(cause error) error {
	return errors.Upgrade(&customErr{
		Cause: cause,
	})
}

//
// actual "program"
//
func unmarshal() (struct{}, error) {
	dest := struct{}{}
	err := json.Unmarshal([]byte("invalid"), &dest) // this wil result in an error
	return dest, errors.Trace(err)
}

func someAction() error {
	data, err := unmarshal()
	if err != nil {
		err := newCustomErr(err)
		errors.Original(err).(*customErr).Value = "some important value"
		return errors.Trace(err)
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
