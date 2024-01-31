// Copyright (c) 2021, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build notrace
// +build notrace

package internal

var TraceStack = false

func EnableTraceStack()  {}
func DisableTraceStack() {}
