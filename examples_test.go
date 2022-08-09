// Copyright (c) 2019, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package errors

import (
	"encoding/json"
	"errors"
	"fmt"
)

func ExampleNew() {
	doSomething := func() error {
		return New("something happened")
	}

	err := doSomething()
	fmt.Printf("%v\n", err)
	// Output: something happened
}

func ExampleNewWithMsg() {
	const ErrSomethingHappened Msg = "something happened"

	doSomething := func() error {
		return New(ErrSomethingHappened)
	}

	err := doSomething()
	fmt.Printf("%v\n", err)
	// Output: something happened
}

func ExampleWithKind() {
	const (
		ErrSomethingWentWrong Msg  = "something went wrong"
		SomeKindOfError       Kind = "some action error"
	)

	doSomethingElse := func() error {
		return New(ErrSomethingWentWrong)
	}
	doSomething := func() error {
		err := doSomethingElse()
		return WithKind(err, SomeKindOfError)
	}

	err := doSomething()
	fmt.Printf("%v\n", err)
	// Output: some action error: something went wrong
}

func ExampleAppend() {
	type Result struct{}

	unmarshal := func() (*Result, error) {
		dest := new(Result)
		err := json.Unmarshal([]byte("invalid"), &dest) // this wil result in an error
		return dest, WithStack(err)
	}

	close := func() error {
		return errors.New("some error occurred while closing something")
	}

	doSomething := func() (err error) {
		defer AppendFunc(&err, close)

		_, unmarshalErr := unmarshal()
		if unmarshalErr != nil {
			Append(&err, unmarshalErr)
			return err
		}
		return nil
	}

	err := doSomething()
	fmt.Printf("%v\n", err)
	// Output:
	// multiple errors occurred:
	// [1/2] invalid character 'i' looking for beginning of value;
	// [2/2] some error occurred while closing something
}

func ExampleCatchPanic() {
	var err error

	done := make(chan struct{})
	go func() {
		defer close(done)
		defer CatchPanic(&err)

		err = New("first error")
		panic("something bad happened")
	}()

	<-done
	fmt.Printf("%v\n", err)
	// Output:
	// multiple errors occurred:
	// [1/2] first error;
	// [2/2] panic: something bad happened
}
