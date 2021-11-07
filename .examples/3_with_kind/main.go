// Copyright (c) 2019, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"

	"github.com/go-pogo/errors"
)

const someError errors.Kind = "some error"

func someAction() error {
	return errors.WithKind(errors.New("something happened"), someError)
}

func doSomething() error {
	return someAction()
}

func main() {
	err := doSomething()
	if err != nil {
		fmt.Printf("%v\n", err)
		fmt.Println("//////////")
		fmt.Printf("%+v\n", err)
	}
}
