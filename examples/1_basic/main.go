package main

import (
	"fmt"

	"github.com/roeldev/go-errs"
)

const SomeError errs.Kind = "some error"

func doSomething() error {
	return errs.New(SomeError, "something happened")
}

func main() {
	err := doSomething()
	if err != nil {
		fmt.Printf("%v\n", err)
		fmt.Println("//////////")
		fmt.Printf("%+v\n", err)
	}
}
