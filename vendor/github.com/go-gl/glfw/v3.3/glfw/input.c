#include "_cgo_export.h"

void glfwSetJoystickCallbackCB() {
  glfwSetJoystickCallback((GLFWjoystickfun)goJoystickCB);
}

void glfwSetKeyCallbackCB(GLFWwindow *window) {
  glfwSetKeyCallback(window, (GLFWkeyfun)goKeyCB);
}

void glfwSetCharCallbackCB(GLFWwindow *window) {
  glfwSetCharCallback(window, (GLFWcharfun)goCharCB);
}

void glfwSetCharModsCallbackCB(GLFWwindow *window) {
  glfwSetCharModsCallback(window, (GLFWcharmodsfun)goCharModsCB);
}

void glfwSetMouseButtonCallbackCB(GLFWwindow *window) {
  glfwSetMouseButtonCallback(window, (GLFWmousebuttonfun)goMouseButtonCB);
}

void glfwSetCursorPosCallbackCB(GLFWwindow *window) {
  glfwSetCursorPosCallback(window, (GLFWcursorposfun)goCursorPosCB);
}

void glfwSetCursorEnterCallbackCB(GLFWwindow *window) {
  glfwSetCursorEnterCallback(window, (GLFWcursorenterfun)goCursorEnterCB);
}

void glfwSetScrollCallbackCB(GLFWwindow *window) {
  glfwSetScrollCallback(window, (GLFWscrollfun)goScrollCB);
}

void glfwSetDropCallbackCB(GLFWwindow *window) {
  glfwSetDropCallback(window, (GLFWdropfun)goDropCB);
}

float GetAxisAtIndex(float *axis, int i) { return axis[i]; }

unsigned char GetButtonsAtIndex(unsigned char *buttons, int i) {
  return buttons[i];
}

unsigned char GetGamepadButtonAtIndex(GLFWgamepadstate *gp, int i) {
  return gp->buttons[i];
}

float GetGamepadAxisAtIndex(GLFWgamepadstate *gp, int i) { return gp->axes[i]; }
