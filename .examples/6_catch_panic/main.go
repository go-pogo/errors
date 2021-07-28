// Copyright (c) 2021, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"

	"github.com/go-pogo/errors"
)

func main() {
	var err error

	done := make(chan struct{})
	go func() {
		defer close(done)
		defer errors.CatchPanic(&err)

		panic("something really bad happened")
	}()

	<-done
	if err != nil {
		fmt.Printf("%v\n", err)
		fmt.Println("//////////")
		fmt.Printf("%+v\n", err)
	}
}
