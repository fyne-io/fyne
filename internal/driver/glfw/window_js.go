//go:build js && !wasm && !test_web_driver
// +build js,!wasm,!test_web_driver

package glfw

import "math"

func (w *window) scaleInput(in float64) float64 {
	return math.Ceil(in * float64(w.canvas.Scale()))

}
