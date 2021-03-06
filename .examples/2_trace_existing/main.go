package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/go-pogo/errors"
)

func unmarshal() (struct{}, error) {
	dest := struct{}{}
	err := json.Unmarshal([]byte("invalid"), &dest) // this wil result in an error
	return dest, errors.Trace(err)
}

func someAction() error {
	data, err := unmarshal()
	if err != nil {
		return errors.Wrapf(err, "something bad happened while performing %s", "someAction")
	}

	// this code never runs
	fmt.Println(data)
	return nil
}

func main() {
	err := someAction()
	if err != nil {
		fmt.Printf("%v\n", err)
		fmt.Println("//////////")
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}
}
