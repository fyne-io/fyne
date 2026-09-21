package glfw

// presentGate reports whether a window's surface is currently presentable and
// lets the render loop register interest in the next presentable moment. On
// Wayland this is backed by wl_surface.frame callbacks; elsewhere it is a
// no-op that always reports ready.
type presentGate interface {
	ready() bool
	requestFrame()
	markReady()
	free()
	// resetAfterDisplayLoss drops a pending wl_callback without destroying
	// it (the Wayland display is already gone) and marks the gate ready so
	// the first frame after glfw.Terminate/Init can present.
	resetAfterDisplayLoss()
}

type noGate struct{}

func (noGate) ready() bool            { return true }
func (noGate) requestFrame()          {}
func (noGate) markReady()             {}
func (noGate) free()                  {}
func (noGate) resetAfterDisplayLoss() {}
