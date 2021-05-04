package glfw

import (
	"runtime"
	"syscall"
	"unsafe"

	"fyne.io/fyne/v2"
)

func (w *window) platformResize(canvasSize fyne.Size) {
	standardResize(w, canvasSize)
}

func (w *window) setDarkMode(dark bool) {
	if runtime.GOOS == "windows" {
		hwnd := w.view().GetWin32Window()

		dwm := syscall.NewLazyDLL("dwmapi.dll")
		setAtt := dwm.NewProc("DwmSetWindowAttribute")
		ret, _, err := setAtt.Call(uintptr(unsafe.Pointer(hwnd)), // window handle
			20,                             // DWMWA_USE_IMMERSIVE_DARK_MODE
			uintptr(unsafe.Pointer(&dark)), // on or off
			8)                              // sizeof(darkMode)

		if ret != 0 { // err is always non-nil, we check return value
			fyne.LogError("Failed to set dark mode", err)
		}
	}
}
