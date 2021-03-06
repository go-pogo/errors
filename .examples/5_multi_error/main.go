package main

import (
	"encoding/json"
	stderrors "errors"
	"fmt"
	"os"

	"github.com/go-pogo/errors"
)

func unmarshal() (struct{}, error) {
	dest := struct{}{}
	err := json.Unmarshal([]byte("invalid"), &dest) // this wil result in an error
	return dest, errors.Trace(err)
}

func finish() error {
	return stderrors.New("some error occurred while closing something")
}

func someAction() (err error) {
	defer errors.Append(&err, finish())

	data, unmarshalErr := unmarshal()
	if unmarshalErr != nil {
		errors.Append(&err, unmarshalErr)
		return err
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
