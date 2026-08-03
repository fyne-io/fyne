//go:build windows

// The Win32 syscall layer: lazily bound procedures and thin Go wrappers.
// Constants and structures live in win32_types.go.

package directx

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	user32   = windows.NewLazySystemDLL("user32.dll")
	kernel32 = windows.NewLazySystemDLL("kernel32.dll")
	shcore   = windows.NewLazySystemDLL("shcore.dll")

	procRegisterClassExW             = user32.NewProc("RegisterClassExW")
	procCreateWindowExW              = user32.NewProc("CreateWindowExW")
	procDestroyWindow                = user32.NewProc("DestroyWindow")
	procDefWindowProcW               = user32.NewProc("DefWindowProcW")
	procShowWindow                   = user32.NewProc("ShowWindow")
	procUpdateWindow                 = user32.NewProc("UpdateWindow")
	procPeekMessageW                 = user32.NewProc("PeekMessageW")
	procGetMessageW                  = user32.NewProc("GetMessageW")
	procTranslateMessage             = user32.NewProc("TranslateMessage")
	procDispatchMessageW             = user32.NewProc("DispatchMessageW")
	procPostQuitMessage              = user32.NewProc("PostQuitMessage")
	procPostMessageW                 = user32.NewProc("PostMessageW")
	procSetWindowTextW               = user32.NewProc("SetWindowTextW")
	procGetClientRect                = user32.NewProc("GetClientRect")
	procGetMenu                      = user32.NewProc("GetMenu")
	procFillRect                     = user32.NewProc("FillRect")
	procGetWindowRect                = user32.NewProc("GetWindowRect")
	procAdjustWindowRectExForDpi     = user32.NewProc("AdjustWindowRectExForDpi")
	procAdjustWindowRectEx           = user32.NewProc("AdjustWindowRectEx")
	procSetWindowPos                 = user32.NewProc("SetWindowPos")
	procGetWindowLongPtrW            = user32.NewProc("GetWindowLongPtrW")
	procSetWindowLongPtrW            = user32.NewProc("SetWindowLongPtrW")
	procLoadCursorW                  = user32.NewProc("LoadCursorW")
	procSetCursor                    = user32.NewProc("SetCursor")
	procDestroyCursor                = user32.NewProc("DestroyCursor")
	procEnumDisplayMonitors          = user32.NewProc("EnumDisplayMonitors")
	procSetCapture                   = user32.NewProc("SetCapture")
	procReleaseCapture               = user32.NewProc("ReleaseCapture")
	procSetFocus                     = user32.NewProc("SetFocus")
	procSetForegroundWindow          = user32.NewProc("SetForegroundWindow")
	procGetDpiForWindow              = user32.NewProc("GetDpiForWindow")
	procSetProcessDpiAwarenessCtx    = user32.NewProc("SetProcessDpiAwarenessContext")
	procMonitorFromWindow            = user32.NewProc("MonitorFromWindow")
	procGetMonitorInfoW              = user32.NewProc("GetMonitorInfoW")
	procGetKeyState                  = user32.NewProc("GetKeyState")
	procTrackMouseEvent              = user32.NewProc("TrackMouseEvent")
	procOpenClipboard                = user32.NewProc("OpenClipboard")
	procCloseClipboard               = user32.NewProc("CloseClipboard")
	procEmptyClipboard               = user32.NewProc("EmptyClipboard")
	procGetClipboardData             = user32.NewProc("GetClipboardData")
	procSetClipboardData             = user32.NewProc("SetClipboardData")
	procGetDoubleClickTime           = user32.NewProc("GetDoubleClickTime")
	procGetSystemMetrics             = user32.NewProc("GetSystemMetrics")
	procGetSystemMetricsForDpi       = user32.NewProc("GetSystemMetricsForDpi")
	procSetMenuItemInfoW             = user32.NewProc("SetMenuItemInfoW")
	procSendMessageW                 = user32.NewProc("SendMessageW")
	procCreateIconIndirect           = user32.NewProc("CreateIconIndirect")
	procDestroyIcon                  = user32.NewProc("DestroyIcon")
	procDragAcceptFiles              = windows.NewLazySystemDLL("shell32.dll").NewProc("DragAcceptFiles")
	procDragQueryFileW               = windows.NewLazySystemDLL("shell32.dll").NewProc("DragQueryFileW")
	procDragQueryPoint               = windows.NewLazySystemDLL("shell32.dll").NewProc("DragQueryPoint")
	procDragFinish                   = windows.NewLazySystemDLL("shell32.dll").NewProc("DragFinish")
	procGetModuleHandleW             = kernel32.NewProc("GetModuleHandleW")
	procSetThreadExecutionState      = kernel32.NewProc("SetThreadExecutionState")
	procGlobalAlloc                  = kernel32.NewProc("GlobalAlloc")
	procGlobalFree                   = kernel32.NewProc("GlobalFree")
	procGlobalLock                   = kernel32.NewProc("GlobalLock")
	procGlobalUnlock                 = kernel32.NewProc("GlobalUnlock")
	procSetProcessDpiAwarenessLegacy = shcore.NewProc("SetProcessDpiAwareness")

	procCreateSolidBrush = windows.NewLazySystemDLL("gdi32.dll").NewProc("CreateSolidBrush")
	procCreateBitmap     = windows.NewLazySystemDLL("gdi32.dll").NewProc("CreateBitmap")
	procCreateDIBSection = windows.NewLazySystemDLL("gdi32.dll").NewProc("CreateDIBSection")
	procDeleteObject     = windows.NewLazySystemDLL("gdi32.dll").NewProc("DeleteObject")

	procDwmSetWindowAttribute = windows.NewLazySystemDLL("dwmapi.dll").NewProc("DwmSetWindowAttribute")
)

// Window messages we care about.

func getModuleHandle() windows.Handle {
	h, _, _ := procGetModuleHandleW.Call(0)
	return windows.Handle(h)
}

func registerClassEx(wc *wndClassExW) (uint16, error) {
	r, _, err := procRegisterClassExW.Call(uintptr(unsafe.Pointer(wc)))
	if r == 0 {
		return 0, err
	}
	return uint16(r), nil
}

func createWindowEx(exStyle uint32, class, title *uint16, style uint32, x, y, w, h int32,
	parent, menu, inst windows.Handle, param unsafe.Pointer) (windows.Handle, error) {
	r, _, err := procCreateWindowExW.Call(uintptr(exStyle), uintptr(unsafe.Pointer(class)),
		uintptr(unsafe.Pointer(title)), uintptr(style), uintptr(x), uintptr(y), uintptr(w), uintptr(h),
		uintptr(parent), uintptr(menu), uintptr(inst), uintptr(param))
	if r == 0 {
		return 0, err
	}
	return windows.Handle(r), nil
}

func defWindowProc(hwnd windows.Handle, m uint32, w, l uintptr) uintptr {
	r, _, _ := procDefWindowProcW.Call(uintptr(hwnd), uintptr(m), w, l)
	return r
}

func destroyWindow(hwnd windows.Handle) {
	procDestroyWindow.Call(uintptr(hwnd))
}

func showWindow(hwnd windows.Handle, cmd int32) {
	procShowWindow.Call(uintptr(hwnd), uintptr(cmd))
	procUpdateWindow.Call(uintptr(hwnd))
}

func peekMessage(m *msg) bool {
	r, _, _ := procPeekMessageW.Call(uintptr(unsafe.Pointer(m)), 0, 0, 0, pmRemove)
	return r != 0
}

func getMessage(m *msg) int32 {
	r, _, _ := procGetMessageW.Call(uintptr(unsafe.Pointer(m)), 0, 0, 0)
	return int32(r)
}

func translateDispatch(m *msg) {
	procTranslateMessage.Call(uintptr(unsafe.Pointer(m)))
	procDispatchMessageW.Call(uintptr(unsafe.Pointer(m)))
}

func postQuitMessage(code int32) {
	procPostQuitMessage.Call(uintptr(code))
}

func postMessage(hwnd windows.Handle, m uint32, w, l uintptr) {
	procPostMessageW.Call(uintptr(hwnd), uintptr(m), w, l)
}

func setWindowText(hwnd windows.Handle, text string) {
	p, err := windows.UTF16PtrFromString(text)
	if err != nil {
		return
	}
	procSetWindowTextW.Call(uintptr(hwnd), uintptr(unsafe.Pointer(p)))
}

func getClientRect(hwnd windows.Handle) rect {
	var r rect
	procGetClientRect.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&r)))
	return r
}

func getWindowRect(hwnd windows.Handle) rect {
	var r rect
	procGetWindowRect.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&r)))
	return r
}

// adjustWindowRect grows a desired client rect into the full window rect for the
// given style. AdjustWindowRectExForDpi is Windows 10 1607+; fall back when absent
// so the driver still runs (slightly wrong frame size) on older builds.
func adjustWindowRect(r *rect, style uint32, dpi uint32, hasMenu bool) {
	menu := uintptr(0)
	if hasMenu {
		menu = 1
	}
	if procAdjustWindowRectExForDpi.Find() == nil {
		procAdjustWindowRectExForDpi.Call(uintptr(unsafe.Pointer(r)), uintptr(style), menu, 0, uintptr(dpi))
		return
	}
	procAdjustWindowRectEx.Call(uintptr(unsafe.Pointer(r)), uintptr(style), menu, 0)
}

// hasMenu reports whether a menu bar is attached, which changes where the client
// area starts and so every frame size derived from it.
func hasMenu(hwnd windows.Handle) bool {
	h, _, _ := procGetMenu.Call(uintptr(hwnd))
	return h != 0
}

func setWindowPos(hwnd windows.Handle, x, y, w, h int32, flags uint32) {
	procSetWindowPos.Call(uintptr(hwnd), 0, uintptr(x), uintptr(y), uintptr(w), uintptr(h), uintptr(flags))
}

func getWindowLongPtr(hwnd windows.Handle, index int32) uintptr {
	r, _, _ := procGetWindowLongPtrW.Call(uintptr(hwnd), uintptr(index))
	return r
}

func setWindowLongPtr(hwnd windows.Handle, index int32, value uintptr) {
	procSetWindowLongPtrW.Call(uintptr(hwnd), uintptr(index), value)
}

func loadCursor(id uint16) windows.Handle {
	h, _, _ := procLoadCursorW.Call(0, uintptr(id))
	return windows.Handle(h)
}

func setCursor(h windows.Handle) {
	procSetCursor.Call(uintptr(h))
}

func getKeyState(vk int32) int16 {
	r, _, _ := procGetKeyState.Call(uintptr(vk))
	return int16(r)
}

// getDpiForWindow reports the window's DPI, defaulting to 96 (scale 1.0) when the
// Windows 10 API is unavailable.
func getDpiForWindow(hwnd windows.Handle) uint32 {
	if procGetDpiForWindow.Find() != nil {
		return 96
	}
	r, _, _ := procGetDpiForWindow.Call(uintptr(hwnd))
	if r == 0 {
		return 96
	}
	return uint32(r)
}

// enableDpiAwareness opts the process into per-monitor DPI so Windows hands us real
// pixels instead of a blurry stretched bitmap.
func enableDpiAwareness() {
	if procSetProcessDpiAwarenessCtx.Find() == nil {
		if r, _, _ := procSetProcessDpiAwarenessCtx.Call(dpiAwarenessContextPerMonitorAwareV2); r != 0 {
			return
		}
	}
	if procSetProcessDpiAwarenessLegacy.Find() == nil {
		procSetProcessDpiAwarenessLegacy.Call(dpiAwarenessPerMonitorLegacy)
	}
}

func monitorRectFor(hwnd windows.Handle) rect {
	h, _, _ := procMonitorFromWindow.Call(uintptr(hwnd), monitorDefaultToNearest)
	mi := monitorInfo{cbSize: uint32(unsafe.Sizeof(monitorInfo{}))}
	if r, _, _ := procGetMonitorInfoW.Call(h, uintptr(unsafe.Pointer(&mi))); r == 0 {
		w, _, _ := procGetSystemMetrics.Call(smCxScreen)
		ht, _, _ := procGetSystemMetrics.Call(smCyScreen)
		return rect{0, 0, int32(w), int32(ht)}
	}
	return mi.rcMonitor
}

// enumMonitors state is package-level because the callback must be created once
// (NewCallback slots are never freed) and Win32 gives it no closure. Only the
// main thread enumerates, matching the rest of the window API.
var (
	enumMonitorsResult rect
	enumMonitorsFound  bool
	enumMonitorsCB     = windows.NewCallback(func(hMon, hdc uintptr, r *rect, lParam uintptr) uintptr {
		mi := monitorInfo{cbSize: uint32(unsafe.Sizeof(monitorInfo{}))}
		if res, _, _ := procGetMonitorInfoW.Call(hMon, uintptr(unsafe.Pointer(&mi))); res != 0 &&
			mi.dwFlags&monitorInfoPrimary == 0 {
			enumMonitorsResult, enumMonitorsFound = mi.rcMonitor, true
			return 0 // found one; stop enumerating
		}
		return 1
	})
)

// secondaryMonitorRect returns the bounds of the first non-primary monitor.
func secondaryMonitorRect() (rect, bool) {
	enumMonitorsFound = false
	procEnumDisplayMonitors.Call(0, 0, enumMonitorsCB, 0)
	return enumMonitorsResult, enumMonitorsFound
}

func monitorWorkRectFor(hwnd windows.Handle) rect {
	h, _, _ := procMonitorFromWindow.Call(uintptr(hwnd), monitorDefaultToNearest)
	mi := monitorInfo{cbSize: uint32(unsafe.Sizeof(monitorInfo{}))}
	if r, _, _ := procGetMonitorInfoW.Call(h, uintptr(unsafe.Pointer(&mi))); r == 0 {
		return monitorRectFor(hwnd)
	}
	return mi.rcWork
}

func trackMouseLeave(hwnd windows.Handle) {
	t := trackMouseEventStruct{
		cbSize:    uint32(unsafe.Sizeof(trackMouseEventStruct{})),
		dwFlags:   tmeLeave,
		hwndTrack: hwnd,
	}
	procTrackMouseEvent.Call(uintptr(unsafe.Pointer(&t)))
}

func doubleClickTime() uint32 {
	r, _, _ := procGetDoubleClickTime.Call()
	if r == 0 {
		return 500
	}
	return uint32(r)
}

// loWord/hiWord unpack the packed coordinate and delta words in mouse messages.
// The cast through int16 keeps negative coordinates (multi-monitor) intact.
func loWord(v uintptr) int32 { return int32(int16(uint32(v) & 0xffff)) }

func hiWord(v uintptr) int32 { return int32(int16((uint32(v) >> 16) & 0xffff)) }

func utf16Ptr(s string) *uint16 {
	p, err := windows.UTF16PtrFromString(s)
	if err != nil {
		return nil
	}
	return p
}

// pointerFromAddr reinterprets an address Win32 handed us as a pointer.
//
// go vet's unsafeptr check rejects a direct uintptr->unsafe.Pointer conversion,
// because for Go heap memory an address can go stale the moment the collector
// stops seeing it as a pointer. That hazard does not apply here: the memory
// belongs to the OS or to a GlobalAlloc block, lives outside the Go heap, and
// never moves.
func pointerFromAddr(addr uintptr) unsafe.Pointer {
	return *(*unsafe.Pointer)(unsafe.Pointer(&addr))
}

// createSolidBrush makes a GDI brush of the given COLORREF (0x00BBGGRR).
func createSolidBrush(colour uint32) windows.Handle {
	h, _, _ := procCreateSolidBrush.Call(uintptr(colour))
	return windows.Handle(h)
}

func fillRect(hdc uintptr, r *rect, brush windows.Handle) {
	procFillRect.Call(hdc, uintptr(unsafe.Pointer(r)), uintptr(brush))
}

func deleteObject(h windows.Handle) {
	procDeleteObject.Call(uintptr(h))
}
