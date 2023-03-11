// +build js,wasm

package glfw

import (
	"errors"
	"syscall/js"
)

func newContext(canvas js.Value, ca *contextAttributes) (context js.Value, err error) {
	if js.Global().Get("WebGLRenderingContext").Equal(js.Undefined()) {
		return js.Value{}, errors.New("Your browser doesn't appear to support WebGL.")
	}

	attrs := map[string]interface{}{
		"alpha":                           ca.Alpha,
		"depth":                           ca.Depth,
		"stencil":                         ca.Stencil,
		"antialias":                       ca.Antialias,
		"premultipliedAlpha":              ca.PremultipliedAlpha,
		"preserveDrawingBuffer":           ca.PreserveDrawingBuffer,
		"preferLowPowerToHighPerformance": ca.PreferLowPowerToHighPerformance,
		"failIfMajorPerformanceCaveat":    ca.FailIfMajorPerformanceCaveat,
	}

	if gl := canvas.Call("getContext", "webgl", attrs); !gl.Equal(js.Null()) {
		debug := js.Global().Get("WebGLDebugUtils")
		if debug.Equal(js.Undefined()) {
			return gl, errors.New("No debugging for WebGL.")
		}
		gl = debug.Call("makeDebugContext", gl)
		return gl, nil
	} else if gl := canvas.Call("getContext", "experimental-webgl", attrs); gl.Equal(js.Null()) {
		return gl, nil
	} else {
		return js.Value{}, errors.New("Creating a WebGL context has failed.")
	}
}

type contextAttributes struct {
	Alpha                           bool
	Depth                           bool
	Stencil                         bool
	Antialias                       bool
	PremultipliedAlpha              bool
	PreserveDrawingBuffer           bool
	PreferLowPowerToHighPerformance bool
	FailIfMajorPerformanceCaveat    bool
}

func defaultAttributes() *contextAttributes {
	return &contextAttributes{
		Alpha:                 false,
		Depth:                 true,
		Stencil:               false,
		Antialias:             false,
		PremultipliedAlpha:    false,
		PreserveDrawingBuffer: false,
	}
}
