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
	return errs.Wrap(err)
}

func main() {
	err := doSomething()
	if err != nil {
		fmt.Print(err)
	}
}
