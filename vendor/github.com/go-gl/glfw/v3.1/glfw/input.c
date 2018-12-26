#include "_cgo_export.h"

void glfwMouseButtonCB(GLFWwindow* window, int button, int action, int mods) {
	goMouseButtonCB(window, button, action, mods);
}

void glfwCursorPosCB(GLFWwindow* window, double xpos, double ypos) {
	goCursorPosCB(window, xpos, ypos);
}

void glfwCursorEnterCB(GLFWwindow* window, int entered) {
	goCursorEnterCB(window, entered);
}

void glfwScrollCB(GLFWwindow* window, double xoff, double yoff) {
	goScrollCB(window, xoff, yoff);
}

void glfwKeyCB(GLFWwindow* window, int key, int scancode, int action, int mods) {
	goKeyCB(window, key, scancode, action, mods);
}

void glfwCharCB(GLFWwindow* window, unsigned int character) {
	goCharCB(window, character);
}

void glfwCharModsCB(GLFWwindow* window, unsigned int character, int mods) {
	goCharModsCB(window, character, mods);
}

void glfwDropCB(GLFWwindow* window, int count, const char **names) {
	goDropCB(window, count, (char**)names);
}

void glfwSetKeyCallbackCB(GLFWwindow *window) {
	glfwSetKeyCallback(window, glfwKeyCB);
}

void glfwSetCharCallbackCB(GLFWwindow *window) {
	glfwSetCharCallback(window, glfwCharCB);
}

void glfwSetCharModsCallbackCB(GLFWwindow *window) {
	glfwSetCharModsCallback(window, glfwCharModsCB);
}

void glfwSetMouseButtonCallbackCB(GLFWwindow *window) {
	glfwSetMouseButtonCallback(window, glfwMouseButtonCB);
}

void glfwSetCursorPosCallbackCB(GLFWwindow *window) {
	glfwSetCursorPosCallback(window, glfwCursorPosCB);
}

void glfwSetCursorEnterCallbackCB(GLFWwindow *window) {
	glfwSetCursorEnterCallback(window, glfwCursorEnterCB);
}

void glfwSetScrollCallbackCB(GLFWwindow *window) {
	glfwSetScrollCallback(window, glfwScrollCB);
}

void glfwSetDropCallbackCB(GLFWwindow *window) {
	glfwSetDropCallback(window, glfwDropCB);
}

float GetAxisAtIndex(float *axis, int i) {
	return axis[i];
}

unsigned char GetButtonsAtIndex(unsigned char *buttons, int i) {
	return buttons[i];
}
