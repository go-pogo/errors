package main

import (
	"fmt"

	"github.com/go-pogo/errors"
)

const SomeError errors.Kind = "some error"

func someAction() error {
	return errors.New(SomeError, "something happened")
}

func doSomething() error {
	err := someAction()
	return errors.Trace(err)
}

func main() {
	err := doSomething()
	if err != nil {
		fmt.Printf("%v\n", err)
		fmt.Println("//////////")
		fmt.Printf("%+v\n", err)
	}
}
