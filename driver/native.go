package driver

// AndroidContext is passed to the `driver.RunNative` callback when it is executed on an Android device.
// The VM, Env and Ctx pointers are reqiured to make various calls into JVM methods.
//
// Since: 2.3
type AndroidContext struct {
	VM, Env, Ctx uintptr
}

// UnknownContext is passed to the RunNative callback when it is executed
// on devices without special native context.
//
// Since: 2.3
type UnknownContext struct{}

// WindowsContext is passed to the `(desktop.Window).RunNative` callback
// when it is executed on a Microsoft Windows desktop device.
//
// Since: 2.5
type WindowsContext struct {
	// HWND is the WinAPI HWND for the window.
	HWND uintptr
}
