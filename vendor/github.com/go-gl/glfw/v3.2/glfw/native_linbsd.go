// +build linux,!wayland freebsd,!wayland

package glfw

//#define GLFW_EXPOSE_NATIVE_X11
//#define GLFW_EXPOSE_NATIVE_GLX
//#include "glfw/include/GLFW/glfw3.h"
//#include "glfw/include/GLFW/glfw3native.h"
import "C"

func GetX11Display() *C.Display {
	ret := C.glfwGetX11Display()
	panicError()
	return ret
}

// GetX11Adapter returns the RRCrtc of the monitor.
func (m *Monitor) GetX11Adapter() C.RRCrtc {
	ret := C.glfwGetX11Adapter(m.data)
	panicError()
	return ret
}

// GetX11Monitor returns the RROutput of the monitor.
func (m *Monitor) GetX11Monitor() C.RROutput {
	ret := C.glfwGetX11Monitor(m.data)
	panicError()
	return ret
}

// GetX11Window returns the Window of the window.
func (w *Window) GetX11Window() C.Window {
	ret := C.glfwGetX11Window(w.data)
	panicError()
	return ret
}

// GetGLXContext returns the GLXContext of the window.
func (w *Window) GetGLXContext() C.GLXContext {
	ret := C.glfwGetGLXContext(w.data)
	panicError()
	return ret
}

// GetGLXWindow returns the GLXWindow of the window.
func (w *Window) GetGLXWindow() C.GLXWindow {
	ret := C.glfwGetGLXWindow(w.data)
	panicError()
	return ret
}
