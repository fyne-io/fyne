package glfw

/*
// Standard OpenGL client is used on 386 and amd64 architectures, except when
// explicitly asked for gles2 or wayland.
#cgo 386,!gles2,!wayland CFLAGS: -D_GLFW_USE_OPENGL
#cgo amd64,!gles2,!wayland CFLAGS: -D_GLFW_USE_OPENGL

// Choose OpenGL ES V2 on arm, or when explicitly asked for gles2/wayland.
#cgo arm gles2 wayland CFLAGS: -D_GLFW_USE_GLESV2


// Windows Build Tags
// ----------------
// GLFW Options:
#cgo windows CFLAGS: -D_GLFW_WIN32 -D_GLFW_WGL

// Linker Options:
#cgo windows LDFLAGS: -lopengl32 -lgdi32


// Darwin Build Tags
// ----------------
// GLFW Options:
#cgo darwin CFLAGS: -D_GLFW_COCOA -D_GLFW_NSGL -D_GLFW_USE_CHDIR -D_GLFW_USE_MENUBAR -D_GLFW_USE_RETINA -Wno-deprecated-declarations

// Linker Options:
#cgo darwin LDFLAGS: -framework Cocoa -framework OpenGL -framework IOKit -framework CoreVideo


// Linux Build Tags
// ----------------
// GLFW Options:
#cgo linux,!wayland CFLAGS: -D_GLFW_X11 -D_GLFW_GLX
#cgo linux,wayland CFLAGS: -D_GLFW_WAYLAND -D_GLFW_EGL

// Linker Options:
#cgo linux,!wayland LDFLAGS: -lGL -lX11 -lXrandr -lXxf86vm -lXi -lXcursor -lm -lXinerama -ldl -lrt
#cgo linux,wayland LDFLAGS: -lGL -lX11 -lXrandr -lXxf86vm -lXi -lXcursor -lm -lXinerama -ldl -lrt


// FreeBSD Build Tags
// ----------------
// GLFW Options:
#cgo freebsd,!wayland CFLAGS: -D_GLFW_X11 -D_GLFW_GLX -D_GLFW_HAS_GLXGETPROCADDRESSARB -D_GLFW_HAS_DLOPEN
#cgo freebsd,wayland CFLAGS: -D_GLFW_WAYLAND -D_GLFW_EGL -D_GLFW_HAS_DLOPEN

// Linker Options:
#cgo freebsd,!wayland LDFLAGS: -lGL -lX11 -lXrandr -lXxf86vm -lXi -lXcursor -lm -lXinerama
#cgo freebsd,wayland LDFLAGS: -lGL -lX11 -lXrandr -lXxf86vm -lXi -lXcursor -lm -lXinerama
*/
import "C"
