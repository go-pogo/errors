package internal

var captureFrames = true

func CaptureFrames() bool   { return captureFrames }
func EnableCaptureFrames()  { captureFrames = true }
func DisableCaptureFrames() { captureFrames = false }
