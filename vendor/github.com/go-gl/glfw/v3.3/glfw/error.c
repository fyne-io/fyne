#include "_cgo_export.h"

void glfwSetErrorCallbackCB() { glfwSetErrorCallback((GLFWerrorfun)goErrorCB); }
