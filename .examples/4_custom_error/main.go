// Copyright (c) 2020, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"

	"github.com/go-pogo/errors"
)

//
// define custom error
//
type customErr struct {
	Cause error
	Value string
}

func (ce *customErr) Error() string {
	return fmt.Sprintf("just a custom error message with `%s`", ce.Value)
}

func (ce *customErr) Format(s fmt.State, v rune) { errors.FormatError(ce, s, v) }

//
// actual "program"
//
func unmarshal() (struct{}, error) {
	dest := struct{}{}
	err := json.Unmarshal([]byte("invalid"), &dest) // this wil result in an error
	return dest, errors.Trace(err)
}

func someAction() error {
	data, err := unmarshal()
	if err != nil {
		return errors.Trace(&customErr{
			Cause: err,
			Value: "some important value",
		})
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
	}
}
