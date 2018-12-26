#include "_cgo_export.h"

GLFWmonitor *GetMonitorAtIndex(GLFWmonitor **monitors, int index) {
	return monitors[index];
}

GLFWvidmode GetVidmodeAtIndex(GLFWvidmode *vidmodes, int index) {
	return vidmodes[index];
}

void glfwMonitorCB(GLFWmonitor* monitor, int event) {
	goMonitorCB(monitor, event);
}

void glfwSetMonitorCallbackCB() {
	glfwSetMonitorCallback(glfwMonitorCB);
}

unsigned int GetGammaAtIndex(unsigned short *color, int i) {
	return color[i];
}

void SetGammaAtIndex(unsigned short *color, int i, unsigned short value) {
	color[i] = value;
}
