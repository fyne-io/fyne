package glfw

//#include <stdlib.h>
//#define GLFW_INCLUDE_NONE
//#include "glfw/include/GLFW/glfw3.h"
import "C"

func glfwbool(b C.int) bool {
	return b == C.int(True)
}

func bytes(origin []byte) (pointer *uint8, free func()) {
	n := len(origin)
	if n == 0 {
		return nil, func() {}
	}

	ptr := C.CBytes(origin)
	return (*uint8)(ptr), func() { C.free(ptr) }
}
