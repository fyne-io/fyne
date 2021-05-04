package glfw

import (
	"fyne.io/fyne/v2"
	"runtime"
	"syscall"
	"unsafe"
)

func (w *window) platformResize(canvasSize fyne.Size) {
	standardResize(w, canvasSize)
}

func (w *window) setDarkMode(dark bool) {
	if runtime.GOOS == "windows" {
		hwnd := w.view().GetWin32Window()

		dwm := syscall.NewLazyDLL("dwmapi.dll")
		setAtt := dwm.NewProc("DwmSetWindowAttribute")
		_, _, err := setAtt.Call(uintptr(unsafe.Pointer(hwnd)), // window handle
			20,                             // DWMWA_USE_IMMERSIVE_DARK_MODE
			uintptr(unsafe.Pointer(&dark)), // on or off
			8)                              // sizeof(darkMode)

		if err != nil {
			fyne.LogError("Failed to set dark mode", err)
		}
	}
}
