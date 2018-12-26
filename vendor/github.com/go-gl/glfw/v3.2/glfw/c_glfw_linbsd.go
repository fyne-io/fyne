// +build linux freebsd

package glfw

/*
#ifdef _GLFW_MIR
	#include "glfw/src/mir_init.c"
	#include "glfw/src/mir_monitor.c"
	#include "glfw/src/mir_window.c"
#endif
#ifdef _GLFW_WAYLAND
	#include "glfw/src/wl_init.c"
	#include "glfw/src/wl_monitor.c"
	#include "glfw/src/wl_window.c"
	#include "glfw/src/wayland-pointer-constraints-unstable-v1-client-protocol.c"
	#include "glfw/src/wayland-relative-pointer-unstable-v1-client-protocol.c"
#endif
#ifdef _GLFW_X11
	#include "glfw/src/x11_init.c"
	#include "glfw/src/x11_monitor.c"
	#include "glfw/src/x11_window.c"
	#include "glfw/src/glx_context.c"
#endif
#include "glfw/src/linux_joystick.c"
#include "glfw/src/posix_time.c"
#include "glfw/src/posix_tls.c"
#include "glfw/src/xkb_unicode.c"
#include "glfw/src/egl_context.c"
*/
import "C"
