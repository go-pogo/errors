// Copyright (c) 2020, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package internal

var captureFrames = true

func CaptureFrames() bool   { return captureFrames }
func EnableCaptureFrames()  { captureFrames = true }
func DisableCaptureFrames() { captureFrames = false }
