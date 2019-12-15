package main

import (
	"fmt"

	"github.com/roeldev/go-errs"
)

func doSomething() error {
	return errs.Err("some error", "something happened")
}

func main() {
	err := doSomething()
	if err != nil {
		fmt.Print(err)
	}
}
