//go:build windows

// Win32 constants and structure layouts. These must match the C headers
// exactly; the wrappers that use them are in win32.go.

package directx

import (
	"golang.org/x/sys/windows"
)

// Window messages we care about.
const (
	wmDestroy       = 0x0002
	wmSize          = 0x0005
	wmSetFocus      = 0x0007
	wmKillFocus     = 0x0008
	wmClose         = 0x0010
	wmEraseBkgnd    = 0x0014
	wmQuit          = 0x0012
	wmSetIcon       = 0x0080
	wmSetCursor     = 0x0020
	wmGetMinMaxInfo = 0x0024
	wmDrawItem      = 0x002B
	wmMeasureItem   = 0x002C
	wmEnterSizeMove = 0x0231
	wmExitSizeMove  = 0x0232
	wmDropFiles     = 0x0233
	wmKeyDown       = 0x0100
	wmKeyUp         = 0x0101
	wmChar          = 0x0102
	wmSysKeyDown    = 0x0104
	wmSysKeyUp      = 0x0105
	wmSysChar       = 0x0106
	wmMouseMove     = 0x0200
	wmLButtonDown   = 0x0201
	wmLButtonUp     = 0x0202
	wmRButtonDown   = 0x0204
	wmRButtonUp     = 0x0205
	wmMButtonDown   = 0x0207
	wmMButtonUp     = 0x0208
	wmMouseWheel    = 0x020A
	wmMouseLeave    = 0x02A3
	wmMouseHWheel   = 0x020E
	wmDpiChanged    = 0x02E0
	wmUser          = 0x0400

	// wmFyneDo wakes the message loop when another goroutine queues work.
	wmFyneDo = wmUser + 1
)

// Window styles.

// Window styles.
const (
	wsOverlappedWindow = 0x00CF0000
	wsPopup            = 0x80000000
	wsVisible          = 0x10000000
	wsThickFrame       = 0x00040000
	wsMaximizeBox      = 0x00010000
	wsCaption          = 0x00C00000
	wsSysMenu          = 0x00080000
	wsMinimizeBox      = 0x00020000

	swHide       = 0
	swShow       = 5
	swShowNormal = 1
	swMaximize   = 3
	swRestore    = 9

	cwUseDefault = ^0x7fffffff // 0x80000000 as a signed int32

	pmRemove = 0x0001

	swpNoSize       = 0x0001
	swpNoMove       = 0x0002
	swpNoZOrder     = 0x0004
	swpNoOwnerZ     = 0x0200
	swpFrameChanged = 0x0020

	gwlStyle    = -16
	gwlUserData = -21

	idcArrow    = 32512
	idcIBeam    = 32513
	idcCross    = 32515
	idcHand     = 32649
	idcSizeNWSE = 32642
	idcSizeNESW = 32643
	idcSizeWE   = 32644
	idcSizeNS   = 32645

	// csOwnDC gives the window its own device context, which the swap chain wants.
	csOwnDC = 0x0020

	// htClient is the WM_SETCURSOR hit-test result meaning "over the client area".
	htClient = 1

	smCxScreen  = 0
	smCyScreen  = 1
	smCMonitors = 80
	// smCxMenuCheck is the width of a menu item's check/icon column, which is the
	// size Windows draws hbmpItem at.
	smCxMenuCheck = 71

	// colorMenuText is COLOR_MENUTEXT, the ink Windows paints menu item labels in.
	colorMenuText = 7

	// miimBitmap selects MENUITEMINFOW.hbmpItem.
	miimBitmap = 0x00000080
	// miimData selects MENUITEMINFOW.dwItemData, where the icon handle is parked
	// for the WM_DRAWITEM that comes back.
	miimData = 0x00000020

	// hbmMenuCallback is HBMMENU_CALLBACK, the hbmpItem value that asks Windows to
	// send WM_MEASUREITEM and WM_DRAWITEM instead of blitting a bitmap itself.
	hbmMenuCallback = ^uintptr(0) // (HBITMAP)-1

	// odtMenu is ODT_MENU, the CtlType for menu owner-draw messages.
	odtMenu = 1
	// diNormal is DI_NORMAL, DrawIconEx drawing both the image and its mask.
	diNormal = 3

	// dibRGBColors: the DIB has no colour table, the pixels are literal BGRA.
	dibRGBColors = 0

	monitorDefaultToNearest = 0x00000002
	// monitorInfoPrimary marks the primary monitor in MONITORINFO.dwFlags.
	monitorInfoPrimary = 0x00000001

	// SetThreadExecutionState flags: esContinuous makes the setting stick until
	// changed, esDisplayRequired keeps the display from blanking.
	esContinuous      = 0x80000000
	esDisplayRequired = 0x00000002

	cfUnicodeText = 13
	gmemMoveable  = 0x0002

	// Per-monitor-v2 DPI awareness context.
	dpiAwarenessContextPerMonitorAwareV2 = ^uintptr(3) // (HANDLE)-4
	dpiAwarenessPerMonitorLegacy         = 2

	tmeLeave = 0x00000002

	// WM_SETICON sizes: small shows in the title bar and taskbar, big in Alt+Tab.
	iconSmall = 0
	iconBig   = 1

	// DWMWA_USE_IMMERSIVE_DARK_MODE, Windows 10 1809 and later.
	dwmwaUseImmersiveDarkMode = 20

	// dragQueryCount is the DragQueryFile index that asks for the file count
	// rather than a name.
	dragQueryCount = 0xFFFFFFFF
)

type wndClassExW struct {
	cbSize        uint32
	style         uint32
	lpfnWndProc   uintptr
	cbClsExtra    int32
	cbWndExtra    int32
	hInstance     windows.Handle
	hIcon         windows.Handle
	hCursor       windows.Handle
	hbrBackground windows.Handle
	lpszMenuName  *uint16
	lpszClassName *uint16
	hIconSm       windows.Handle
}

type point struct {
	X, Y int32
}

type rect struct {
	Left, Top, Right, Bottom int32
}

type msg struct {
	hwnd    windows.Handle
	message uint32
	wParam  uintptr
	lParam  uintptr
	time    uint32
	pt      point
}

type minMaxInfo struct {
	ptReserved     point
	ptMaxSize      point
	ptMaxPosition  point
	ptMinTrackSize point
	ptMaxTrackSize point
}

type monitorInfo struct {
	cbSize    uint32
	rcMonitor rect
	rcWork    rect
	dwFlags   uint32
}

type trackMouseEventStruct struct {
	cbSize      uint32
	dwFlags     uint32
	hwndTrack   windows.Handle
	dwHoverTime uint32
}

type iconInfo struct {
	fIcon    int32
	xHotspot uint32
	yHotspot uint32
	hbmMask  windows.Handle
	hbmColor windows.Handle
}

type menuItemInfoW struct {
	cbSize        uint32
	fMask         uint32
	fType         uint32
	fState        uint32
	wID           uint32
	hSubMenu      windows.Handle
	hbmpChecked   windows.Handle
	hbmpUnchecked windows.Handle
	dwItemData    uintptr
	dwTypeData    *uint16
	cch           uint32
	hbmpItem      windows.Handle
}

type bitmapInfoHeader struct {
	biSize          uint32
	biWidth         int32
	biHeight        int32
	biPlanes        uint16
	biBitCount      uint16
	biCompression   uint32
	biSizeImage     uint32
	biXPelsPerMeter int32
	biYPelsPerMeter int32
	biClrUsed       uint32
	biClrImportant  uint32
}

// measureItemStruct is MEASUREITEMSTRUCT, the reply buffer for WM_MEASUREITEM.
type measureItemStruct struct {
	ctlType    uint32
	ctlID      uint32
	itemID     uint32
	itemWidth  uint32
	itemHeight uint32
	itemData   uintptr
}

// drawItemStruct is DRAWITEMSTRUCT, the request for WM_DRAWITEM.
type drawItemStruct struct {
	ctlType    uint32
	ctlID      uint32
	itemID     uint32
	itemAction uint32
	itemState  uint32
	hwndItem   windows.Handle
	hDC        windows.Handle
	rcItem     rect
	itemData   uintptr
}
