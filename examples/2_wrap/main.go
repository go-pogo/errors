package main

import (
	"fmt"

	"github.com/roeldev/go-errs"
)

func someAction() error {
	return errs.Err("some error", "something happened")
}

func doSomething() error {
	err := someAction()
	if err != nil {
		return errs.Wrap(err)
	}
	return nil
}

func main() {
	err := doSomething()
	if err != nil {
		fmt.Print(err)
	}
}
