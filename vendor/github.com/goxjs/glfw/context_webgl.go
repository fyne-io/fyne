// +build js,!wasm

package glfw

import (
	"errors"

	"github.com/gopherjs/gopherjs/js"
)

func newContext(canvas *js.Object, ca *contextAttributes) (context *js.Object, err error) {
	if js.Global.Get("WebGLRenderingContext") == js.Undefined {
		return nil, errors.New("Your browser doesn't appear to support WebGL.")
	}

	attrs := map[string]bool{
		"alpha":                           ca.Alpha,
		"depth":                           ca.Depth,
		"stencil":                         ca.Stencil,
		"antialias":                       ca.Antialias,
		"premultipliedAlpha":              ca.PremultipliedAlpha,
		"preserveDrawingBuffer":           ca.PreserveDrawingBuffer,
		"preferLowPowerToHighPerformance": ca.PreferLowPowerToHighPerformance,
		"failIfMajorPerformanceCaveat":    ca.FailIfMajorPerformanceCaveat,
	}

	if gl := canvas.Call("getContext", "webgl", attrs); gl != nil {
		return gl, nil
	} else if gl := canvas.Call("getContext", "experimental-webgl", attrs); gl != nil {
		return gl, nil
	} else {
		return nil, errors.New("Creating a WebGL context has failed.")
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
