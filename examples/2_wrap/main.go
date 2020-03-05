package main

import (
	"fmt"

	"github.com/roeldev/go-errs"
)

func someAction() error {
	return errs.New("some error", "something happened")
}

func doSomething() error {
	err := someAction()
	return errs.Wrap(err)
}

func main() {
	err := doSomething()
	if err != nil {
		errs.GetStackTrace(err).Capture(0)
		// err = errs.Wrap(err)
		fmt.Print(err)
	}
}
