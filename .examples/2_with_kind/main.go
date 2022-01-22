// Copyright (c) 2019, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"

	"github.com/go-pogo/errors"
)

const (
	ErrSomethingWentWrong errors.Msg  = "something went wrong"
	ActionError           errors.Kind = "action error"
)

func doSomething() error {
	return errors.WithKind(ErrSomethingWentWrong, ActionError)
}

func main() {
	err := doSomething()
	if err != nil {
		fmt.Printf("%v\n", err)
		fmt.Println("//////////")
		fmt.Printf("%+v\n", err)
	}
}
