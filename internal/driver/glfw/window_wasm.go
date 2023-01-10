//go:build wasm || test_web_driver
// +build wasm test_web_driver

package glfw

func (w *window) scaleInput(in float64) float64 {
	return in
}
