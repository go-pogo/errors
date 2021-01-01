package main

import (
	"fmt"

	"github.com/go-pogo/errors"
)

func doSomething() error {
	return errors.New("something happened")
}

func main() {
	err := doSomething()
	if err != nil {
		fmt.Printf("%v\n", err)
		fmt.Println("//////////")
		fmt.Printf("%+v\n", err)
	}
}
