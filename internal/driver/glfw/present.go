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
}

type noGate struct{}

func (noGate) ready() bool   { return true }
func (noGate) requestFrame() {}
func (noGate) markReady()    {}
func (noGate) free()         {}
