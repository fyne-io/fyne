#include "_cgo_export.h"

void glfwSetWindowPosCallbackCB(GLFWwindow *window) {
  glfwSetWindowPosCallback(window, (GLFWwindowposfun)goWindowPosCB);
}

void glfwSetWindowSizeCallbackCB(GLFWwindow *window) {
  glfwSetWindowSizeCallback(window, (GLFWwindowsizefun)goWindowSizeCB);
}

void glfwSetWindowCloseCallbackCB(GLFWwindow *window) {
  glfwSetWindowCloseCallback(window, (GLFWwindowclosefun)goWindowCloseCB);
}

void glfwSetWindowRefreshCallbackCB(GLFWwindow *window) {
  glfwSetWindowRefreshCallback(window, (GLFWwindowrefreshfun)goWindowRefreshCB);
}

void glfwSetWindowFocusCallbackCB(GLFWwindow *window) {
  glfwSetWindowFocusCallback(window, (GLFWwindowfocusfun)goWindowFocusCB);
}

void glfwSetWindowIconifyCallbackCB(GLFWwindow *window) {
  glfwSetWindowIconifyCallback(window, (GLFWwindowiconifyfun)goWindowIconifyCB);
}

void glfwSetFramebufferSizeCallbackCB(GLFWwindow *window) {
  glfwSetFramebufferSizeCallback(window,
                                 (GLFWframebuffersizefun)goFramebufferSizeCB);
}

void glfwSetWindowMaximizeCallbackCB(GLFWwindow *window) {
  glfwSetWindowMaximizeCallback(window,
                                (GLFWwindowmaximizefun)goWindowMaximizeCB);
}

void glfwSetWindowContentScaleCallbackCB(GLFWwindow *window) {
  glfwSetWindowContentScaleCallback(
      window, (GLFWwindowcontentscalefun)goWindowContentScaleCB);
}
