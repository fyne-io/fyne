package glfw

//#define GLFW_EXPOSE_NATIVE_COCOA
//#define GLFW_EXPOSE_NATIVE_NSGL
//#include <GLFW/glfw3.h>
//#include <GLFW/glfw3native.h>
import "C"

func (w *Window) GetCocoaWindow() C.id {
	return C.glfwGetCocoaWindow(w.data)
}

func (w *Window) GetNSGLContext() C.id {
	return C.glfwGetNSGLContext(w.data)
}
