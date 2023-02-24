package glfw

import (
	"runtime"
	"syscall"
	"unsafe"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/scale"

	"golang.org/x/sys/windows/registry"
)

func (w *window) setDarkMode() {
	if runtime.GOOS == "windows" {
		hwnd := w.view().GetWin32Window()
		dark := isDark()

		dwm := syscall.NewLazyDLL("dwmapi.dll")
		setAtt := dwm.NewProc("DwmSetWindowAttribute")
		ret, _, err := setAtt.Call(uintptr(unsafe.Pointer(hwnd)), // window handle
			20,                             // DWMWA_USE_IMMERSIVE_DARK_MODE
			uintptr(unsafe.Pointer(&dark)), // on or off
			8)                              // sizeof(darkMode)

		if ret != 0 && ret != 0x80070057 { // err is always non-nil, we check return value (except erroneous code)
			fyne.LogError("Failed to set dark mode", err)
		}
	}
}

func isDark() bool {
	k, err := registry.OpenKey(registry.CURRENT_USER, `SOFTWARE\Microsoft\Windows\CurrentVersion\Themes\Personalize`, registry.QUERY_VALUE)
	if err != nil { // older version of Windows will not have this key
		return false
	}
	defer k.Close()

	useLight, _, err := k.GetIntegerValue("AppsUseLightTheme")
	if err != nil { // older version of Windows will not have this value
		return false
	}

	return useLight == 0
}

func (w *window) computeCanvasSize(width, height int) fyne.Size {
	if w.fixedSize {
		return fyne.NewSize(scale.ToFyneCoordinate(w.canvas, w.width), scale.ToFyneCoordinate(w.canvas, w.height))
	}
	return fyne.NewSize(scale.ToFyneCoordinate(w.canvas, width), scale.ToFyneCoordinate(w.canvas, height))
}
