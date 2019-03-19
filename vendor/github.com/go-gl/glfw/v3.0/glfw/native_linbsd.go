//+build linux freebsd

package glfw

//#define GLFW_EXPOSE_NATIVE_X11
//#define GLFW_EXPOSE_NATIVE_GLX
//#include <X11/extensions/Xrandr.h>
//#include <GLFW/glfw3.h>
//#include <GLFW/glfw3native.h>
import "C"

func (w *Window) GetX11Window() C.Window {
	return C.glfwGetX11Window(w.data)
}

func (w *Window) GetGLXContext() C.GLXContext {
	return C.glfwGetGLXContext(w.data)
}

func GetX11Display() *C.Display {
	return C.glfwGetX11Display()
}
