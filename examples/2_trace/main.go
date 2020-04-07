package main

import (
	"fmt"

	"github.com/roeldev/go-errs"
)

const SomeError errs.Kind = "some error"

func someAction() error {
	return errs.New(SomeError, "something happened")
}

func doSomething() error {
	err := someAction()
	return errs.Trace(err)
}

func main() {
	err := doSomething()
	if err != nil {
		fmt.Printf("%v\n", err)
		fmt.Println("//////////")
		fmt.Printf("%+v\n", err)
	}
}
