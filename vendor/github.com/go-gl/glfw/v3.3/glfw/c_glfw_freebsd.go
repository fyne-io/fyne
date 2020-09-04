// +build freebsd

package glfw

/*
#ifdef _GLFW_WAYLAND
	#include "glfw/src/wl_init.c"
	#include "glfw/src/wl_monitor.c"
	#include "glfw/src/wl_window.c"
	#include "glfw/src/wl_platform.h"
#endif
#ifdef _GLFW_X11
	#include "glfw/src/x11_init.c"
	#include "glfw/src/x11_monitor.c"
	#include "glfw/src/x11_window.c"
	#include "glfw/src/glx_context.c"
#endif
#include "glfw/src/null_joystick.c"
#include "glfw/src/posix_time.c"
#include "glfw/src/posix_thread.c"
#include "glfw/src/xkb_unicode.c"
#include "glfw/src/egl_context.c"
*/
import "C"
