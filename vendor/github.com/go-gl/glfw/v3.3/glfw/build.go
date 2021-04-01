package glfw

/*
// Windows Build Tags
// ----------------
// GLFW Options:
#cgo windows CFLAGS: -D_GLFW_WIN32 -Iglfw/deps/mingw

// Linker Options:
#cgo windows LDFLAGS: -lgdi32

#cgo !gles2,windows LDFLAGS: -lopengl32
#cgo gles2,windows LDFLAGS: -lGLESv2

// Darwin Build Tags
// ----------------
// GLFW Options:
#cgo darwin CFLAGS: -D_GLFW_COCOA -Wno-deprecated-declarations

// Linker Options:
#cgo darwin LDFLAGS: -framework Cocoa -framework IOKit -framework CoreVideo

#cgo !gles2,darwin LDFLAGS: -framework OpenGL
#cgo gles2,darwin LDFLAGS: -lGLESv2

// Linux Build Tags
// ----------------
// GLFW Options:
#cgo linux,!wayland CFLAGS: -D_GLFW_X11 -D_GNU_SOURCE
#cgo linux,wayland CFLAGS: -D_GLFW_WAYLAND -D_GNU_SOURCE

// Linker Options:
#cgo linux,!gles1,!gles2,!gles3,!vulkan LDFLAGS: -lGL
#cgo linux,gles1 LDFLAGS: -lGLESv1
#cgo linux,gles2 LDFLAGS: -lGLESv2
#cgo linux,gles3 LDFLAGS: -lGLESv3
#cgo linux,vulkan LDFLAGS: -lvulkan
#cgo linux,!wayland LDFLAGS: -lX11 -lXrandr -lXxf86vm -lXi -lXcursor -lm -lXinerama -ldl -lrt
#cgo linux,wayland LDFLAGS: -lwayland-client -lwayland-cursor -lwayland-egl -lxkbcommon -lm -ldl -lrt

// FreeBSD Build Tags
// ----------------
// GLFW Options:
#cgo freebsd pkg-config: glfw3
#cgo freebsd CFLAGS: -D_GLFW_HAS_DLOPEN
#cgo freebsd,!wayland CFLAGS: -D_GLFW_X11 -D_GLFW_HAS_GLXGETPROCADDRESSARB
#cgo freebsd,wayland CFLAGS: -D_GLFW_WAYLAND

// Linker Options:
#cgo freebsd,!wayland LDFLAGS: -lm -lGL -lX11 -lXrandr -lXxf86vm -lXi -lXcursor -lXinerama
#cgo freebsd,wayland LDFLAGS: -lm -lGL -lwayland-client -lwayland-cursor -lwayland-egl -lxkbcommon
*/
import "C"
