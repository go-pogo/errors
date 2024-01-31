// Copyright (c) 2020, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !notrace

package internal

var TraceStack = true

func EnableTraceStack()  { TraceStack = true }
func DisableTraceStack() { TraceStack = false }
