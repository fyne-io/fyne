package glfw

//#include <GLFW/glfw3.h>
import "C"

func glfwbool(b C.int) bool {
	if b == C.GL_TRUE {
		return true
	}
	return false
}
