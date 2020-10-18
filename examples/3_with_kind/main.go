package main

import (
	"fmt"

	"github.com/go-pogo/errors"
)

const someError errors.Kind = "some error"

func someAction() error {
	return errors.WithKind(errors.New("something happened"), someError)
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
