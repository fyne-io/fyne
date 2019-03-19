package glfw

//#include <stdlib.h>
//#include "glfw/include/GLFW/glfw3.h"
import "C"

import (
	"reflect"
	"unsafe"
)

func glfwbool(b C.int) bool {
	if b == C.GL_TRUE {
		return true
	}
	return false
}

func bytes(origin []byte) (pointer *uint8, free func()) {
	n := len(origin)

	if n == 0 {
		return nil, func() {}
	}

	data := C.malloc(C.size_t(n))

	dataSlice := *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(data),
		Len:  n,
		Cap:  n,
	}))

	copy(dataSlice, origin)

	return &dataSlice[0], func() { C.free(data) }
}
