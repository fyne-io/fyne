// +build !js

package glfw

import "github.com/go-gl/glfw/v3.3/glfw"

type Hint int

const (
	AlphaBits   = Hint(glfw.AlphaBits)
	DepthBits   = Hint(glfw.DepthBits)
	StencilBits = Hint(glfw.StencilBits)
	Samples     = Hint(glfw.Samples)
	Resizable   = Hint(glfw.Resizable)

	// These hints used for WebGL contexts, ignored on desktop.
	PremultipliedAlpha = noopHint
	PreserveDrawingBuffer
	PreferLowPowerToHighPerformance
	FailIfMajorPerformanceCaveat
)

// noopHint is ignored.
const noopHint Hint = -1

func WindowHint(target Hint, hint int) {
	if target == noopHint {
		return
	}

	glfw.WindowHint(glfw.Hint(target), hint)
}
