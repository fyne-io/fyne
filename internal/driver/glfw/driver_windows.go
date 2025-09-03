package glfw

import (
	"fmt"
	"runtime"
	"syscall"
	"time"
	"unsafe"
)

type (
	MB uint32
	ES uint
)

const (
	MB_OK        MB = 0x0000_0000
	MB_ICONERROR MB = 0x0000_0010

	ES_CONTINUOUS       ES = 0x80000000
	ES_DISPLAY_REQUIRED ES = 0x00000002
)

var (
	kernel32 = syscall.NewLazyDLL("kernel32.dll")
	user32   = syscall.NewLazyDLL("user32.dll")

	executionState     = kernel32.NewProc("SetThreadExecutionState")
	MessageBox         = user32.NewProc("MessageBoxW")
	getDoubleClickTime = user32.NewProc("GetDoubleClickTime")
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
	uType := MB_OK | MB_ICONERROR

	syscall.SyscallN(MessageBox.Addr(),
		uintptr(unsafe.Pointer(nil)), uintptr(unsafe.Pointer(toNativePtr(text))),
		uintptr(unsafe.Pointer(toNativePtr(caption))), uintptr(uType))
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

func setDisableScreenBlank(disable bool) {
	uType := ES_CONTINUOUS
	if disable {
		uType |= ES_DISPLAY_REQUIRED
	}

	syscall.SyscallN(executionState.Addr(), uintptr(uType))
}

func (d *gLDriver) DoubleTapDelay() time.Duration {
	// https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-getdoubleclicktime
	if getDoubleClickTime == nil {
		return desktopDefaultDoubleTapDelay
	}
	r1, _, _ := syscall.SyscallN(getDoubleClickTime.Addr())
	return time.Duration(uint64(r1) * uint64(time.Millisecond))
}
