//+build windows

package fyne

import (
	"syscall"
	"unsafe"
)

// hideConsoleOnWindows hides the associated console window that gets created
// for Windows applications that are of type console instead of type GUI. When
// building you can pass the ldflag H=windowsgui to suppress this but if you
// just go build or go run, a console window will pop open along with the GUI
// window. hideConsoleOnWindows hides it.
func hideConsoleOnWindows() {
	console := getConsoleWindow()
	if console == 0 {
		return // No console attached.
	}
	// If this application is the process that created the console window, then
	// this program was not compiled with the -H=windowsgui flag and on start-up
	// it created a console along with the main application window. In this case
	// hide the console window. See
	// http://stackoverflow.com/questions/9009333/how-to-check-if-the-program-is-run-from-a-console
	_, consoleProcID := getWindowThreadProcessId(console)
	if getCurrentProcessId() == consoleProcID {
		const SW_HIDE = 0
		showWindowAsync(console, SW_HIDE)
	}
}

var (
	kernel32                = syscall.NewLazyDLL("kernel32.dll")
	procGetConsoleWindow    = kernel32.NewProc("GetConsoleWindow")
	procGetCurrentProcessId = kernel32.NewProc("GetCurrentProcessId")

	user32                       = syscall.NewLazyDLL("user32.dll")
	procGetWindowThreadProcessId = user32.NewProc("GetWindowThreadProcessId")
	procShowWindowAsync          = user32.NewProc("ShowWindowAsync")
)

func getConsoleWindow() uintptr {
	ret, _, _ := procGetConsoleWindow.Call()
	return ret
}

func getCurrentProcessId() uint32 {
	id, _, _ := procGetCurrentProcessId.Call()
	return uint32(id)
}

func getWindowThreadProcessId(hwnd uintptr) (uintptr, uint32) {
	var processId uint32
	ret, _, _ := procGetWindowThreadProcessId.Call(
		hwnd,
		uintptr(unsafe.Pointer(&processId)),
	)
	return ret, processId
}

func showWindowAsync(hwnd uintptr, cmdshow uintptr) {
	procShowWindowAsync.Call(hwnd, cmdshow)
}
