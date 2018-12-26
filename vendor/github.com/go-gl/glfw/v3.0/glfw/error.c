#include "_cgo_export.h"

void glfwErrorCB(int code, const char *desc) {
	goErrorCB(code, (char*)desc);
}

void glfwSetErrorCallbackCB() {
	glfwSetErrorCallback(glfwErrorCB);
}
