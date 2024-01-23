// Copyright (c) 2020, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !notrace

package errors

var traceStack = true

func enableTraceStack()  { traceStack = true }
func disableTraceStack() { traceStack = false }
