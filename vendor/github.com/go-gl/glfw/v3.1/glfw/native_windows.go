package glfw

//#define GLFW_EXPOSE_NATIVE_WIN32
//#define GLFW_EXPOSE_NATIVE_WGL
//#include "glfw/include/GLFW/glfw3.h"
//#include "glfw/include/GLFW/glfw3native.h"
import "C"

func (w *Window) GetWin32Window() C.HWND {
	ret := C.glfwGetWin32Window(w.data)
	panicError()
	return ret
}

func (w *Window) GetWGLContext() C.HGLRC {
	ret := C.glfwGetWGLContext(w.data)
	panicError()
	return ret
}
