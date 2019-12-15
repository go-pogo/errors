package main

import (
	"encoding/json"
	"fmt"

	"github.com/roeldev/go-errs"
)

func someAction() error {
	dest := new(struct{})
	return json.Unmarshal([]byte("invalid json"), dest)
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
