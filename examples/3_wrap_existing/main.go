package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/roeldev/go-errs"
)

func unmarshal() (struct{}, error) {
	dest := struct{}{}
	err := json.Unmarshal([]byte("invalid"), &dest) // this wil result in an error
	return dest, err
}

func run() error {
	data, err := unmarshal()
	if err != nil {
		return errs.Wrap(err)
	}

	fmt.Println(data)
	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
}
