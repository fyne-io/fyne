package driver

// NativeWindow is an extension interface for `fyne.Window` that gives access
// to platform-native features of application windows.
//
// Since: 2.5
type NativeWindow interface {
	// RunNative  provides a way to execute code within the platform-specific runtime context for a window.
	// The context types are defined in the `driver` package and the specific context passed will differ by platform.
	RunNative(func(context any) error) error
}

// AndroidContext is passed to the RunNative callback when it is executed on an Android device.
// The VM, Env and Ctx pointers are reqiured to make various calls into JVM methods.
//
// Since: 2.3
type AndroidContext struct {
	VM, Env, Ctx uintptr
}

// UnknownContext is passed to the RunNative callback when it is executed
// on devices or windows without special native context.
//
// Since: 2.3
type UnknownContext struct{}

// WindowsWindowContext is passed to the NativeWindow.RunNative callback
// when it is executed on a Microsoft Windows device.
//
// Since: 2.5
type WindowsWindowContext struct {
	// HWND is the window handle for the native window.
	HWND uintptr
}

// MacWindowContext is passed to the NativeWindow.RunNative callback
// when it is executed on a Mac OS device.
//
// Since: 2.5
type MacWindowContext struct {
	// NSWindow is the window handle for the native window.
	NSWindow uintptr
}

// X11WindowContext is passed to the NativeWindow.RunNative callback
// when it is executed on a device with the X11 windowing system.
//
// Since: 2.5
type X11WindowContext struct {
	// WindowHandle is the window handle for the native X11 window.
	WindowHandle string
}

// WaylandWindowContext is passed to the NativeWindow.RunNative callback
// when it is executed on a device with the Wayland windowing system.
//
// Since: 2.5
type WaylandWindowContext struct {
	// WaylandSurface is the handle to the native Wayland surface.
	WaylandSurface uintptr

	// EGLSurface is the handle to the native EGL surface.
	EGLSurface uintptr
}
