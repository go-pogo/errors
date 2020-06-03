package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/roeldev/go-errs"
)

func unmarshal() (struct{}, error) {
	dest := struct{}{}
	err := json.Unmarshal([]byte("invalid"), &dest) // this wil result in an error
	return dest, errs.Trace(err)
}

func finish() error {
	return errors.New("some error occurred while closing something")
}

func someAction() (err error) {
	defer errs.Append(&err, finish())

	data, unmarshalErr := unmarshal()
	if unmarshalErr != nil {
		return errs.Append(&err, unmarshalErr)
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
