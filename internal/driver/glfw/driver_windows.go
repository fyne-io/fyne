package glfw

import (
	"fmt"
	"runtime"
	"syscall"
	"unsafe"
)

type MB uint32

const (
	MB_OK        MB = 0x0000_0000
	MB_ICONERROR MB = 0x0000_0010
)

func toNativePtr(s string) *uint16 {
	pstr, err := syscall.UTF16PtrFromString(s)
	if err != nil {
		panic(fmt.Sprintf("toNativePtr() failed \"%s\": %s", s, err))
	}
	return pstr
}

// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-messageboxw
func messageBoxError(text, caption string) {
	user32 := syscall.NewLazyDLL("user32.dll")
	MessageBox := user32.NewProc("MessageBoxW")

	uType := MB_OK | MB_ICONERROR

	syscall.Syscall6(MessageBox.Addr(), 4,
		uintptr(unsafe.Pointer(nil)), uintptr(unsafe.Pointer(toNativePtr(text))),
		uintptr(unsafe.Pointer(toNativePtr(caption))), uintptr(uType),
		0, 0)
}

func logError(msg string, err error) {
	text := fmt.Sprintf("Fyne error: %v", msg)
	if err != nil {
		text = text + fmt.Sprintf("\n  Cause:%v", err)
	}

	_, file, line, ok := runtime.Caller(1)
	if ok {
		text = text + fmt.Sprintf("\n  At: %s:%d", file, line)
	}

	messageBoxError(text, "Fyne Error")
}
