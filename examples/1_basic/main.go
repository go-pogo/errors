package main

import (
	"fmt"

	"github.com/go-pogo/errors"
)

const SomeError errors.Kind = "some error"

func doSomething() error {
	return errors.New(SomeError, "something happened")
}

func main() {
	err := doSomething()
	if err != nil {
		fmt.Printf("%v\n", err)
		fmt.Println("//////////")
		fmt.Printf("%+v\n", err)
	}
}
