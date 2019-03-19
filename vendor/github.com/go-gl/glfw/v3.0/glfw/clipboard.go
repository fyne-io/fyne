package glfw

//#include <stdlib.h>
//#include <GLFW/glfw3.h>
import "C"

import (
	"errors"
	"unsafe"
)

//SetClipboardString sets the system clipboard to the specified UTF-8 encoded
//string.
//
//This function may only be called from the main thread. See
//https://code.google.com/p/go-wiki/wiki/LockOSThread
func (w *Window) SetClipboardString(str string) {
	cp := C.CString(str)
	defer C.free(unsafe.Pointer(cp))

	C.glfwSetClipboardString(w.data, cp)
}

//GetClipboardString returns the contents of the system clipboard, if it
//contains or is convertible to a UTF-8 encoded string.
//
//This function may only be called from the main thread. See
//https://code.google.com/p/go-wiki/wiki/LockOSThread
func (w *Window) GetClipboardString() (string, error) {
	cs := C.glfwGetClipboardString(w.data)
	if cs == nil {
		return "", errors.New("Can't get clipboard string.")
	}

	return C.GoString(cs), nil
}
