// Copyright (c) 2020, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"

	"github.com/go-pogo/errors"
)

type customErr struct {
	cause error
	value string
}

func (ce *customErr) Unwrap() error { return ce.cause }

func (ce *customErr) Error() string {
	return fmt.Sprintf("just a custom error message with `%s`", ce.value)
}

//
// actual "program"
//
func unmarshal() (struct{}, error) {
	dest := struct{}{}
	err := json.Unmarshal([]byte("invalid"), &dest) // this wil result in an error
	return dest, errors.WithStack(err)
}

func doSomething() error {
	data, err := unmarshal()
	if err != nil {
		return errors.WithStack(&customErr{
			cause: err,
			value: "some important value",
		})
	}

	// this code never runs
	fmt.Println(data)
	return nil
}

func main() {
	err := doSomething()
	if err != nil {
		fmt.Printf("%v\n", err)
		fmt.Println("//////////")
		fmt.Printf("%+v\n", err)
	}
}
