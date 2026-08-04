//go:build windows

// Win32 window class registration and the window procedure: every message the
// driver cares about is dispatched from here to the handlers in window.go,
// mouse.go and keyboard.go.

package directx

import (
	"sync"
	"unicode/utf16"
	"unicode/utf8"
	"unsafe"

	"golang.org/x/sys/windows"

	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/internal/scale"
)

const windowClassName = "FyneDirectXWindow"

var (
	classOnce sync.Once
	classErr  error
	// windowsByHandle resolves a HWND back to its window. Only the main thread
	// touches it, matching the single-threaded message loop.
	windowsByHandle = map[windows.Handle]*window{}
	// pendingWindow holds the window being created, because the first messages
	// arrive from inside CreateWindowEx before it returns a handle.
	pendingWindow *window
)

func registerWindowClass() error {
	classOnce.Do(func() {
		wc := wndClassExW{
			cbSize: uint32(unsafe.Sizeof(wndClassExW{})),
			// CS_OWNDC only. CS_HREDRAW/CS_VREDRAW would invalidate the whole client
			// area on every resize step, so the background erase below would repaint
			// the entire window instead of just the strip the drag exposed - trading
			// a black edge for a full-window flash. Every WM_SIZE repaints
			// everything anyway, so there is nothing for those styles to buy.
			style:         csOwnDC,
			lpfnWndProc:   windows.NewCallback(wndProc),
			hInstance:     getModuleHandle(),
			hCursor:       loadCursor(idcArrow),
			lpszClassName: utf16Ptr(windowClassName),
		}
		_, classErr = registerClassEx(&wc)
	})
	return classErr
}

func wndProc(hwnd windows.Handle, message uint32, wParam, lParam uintptr) uintptr {
	w := windowsByHandle[hwnd]
	if w == nil {
		if pendingWindow != nil {
			w = pendingWindow
			w.hwnd = hwnd
			windowsByHandle[hwnd] = w
		} else {
			return defWindowProc(hwnd, message, wParam, lParam)
		}
	}
	return w.handleMessage(message, wParam, lParam)
}

func (w *window) handleMessage(message uint32, wParam, lParam uintptr) uintptr {
	switch message {
	case wmClose:
		if w.onCloseInter != nil {
			w.onCloseInter()
			return 0
		}
		w.Close()
		return 0

	case wmDestroy:
		w.closing = true
		return 0

	case wmEraseBkgnd:
		w.eraseBackground(wParam) // wParam is the device context
		return 1                  // handled; do not let DefWindowProc erase again

	case wmSize:
		w.resized(loWord(lParam), hiWord(lParam))
		return 0

	case wmGetMinMaxInfo:
		// Let Windows enforce the size limits during the drag itself. Correcting an
		// undersized window afterwards with SetWindowPos fights the pointer and the
		// content visibly oscillates.
		info := (*minMaxInfo)(pointerFromAddr(lParam))
		if w.fixedSize {
			cw, ch := w.frameSizeFor(w.width, w.height)
			info.ptMinTrackSize = point{X: cw, Y: ch}
			info.ptMaxTrackSize = point{X: cw, Y: ch}
			return 0
		}
		if w.canvas != nil && w.canvas.Content() != nil {
			minSize := w.canvas.MinSize()
			cw, ch := w.frameSizeFor(
				scale.ToScreenCoordinate(w.canvas, minSize.Width),
				scale.ToScreenCoordinate(w.canvas, minSize.Height),
			)
			info.ptMinTrackSize = point{X: cw, Y: ch}
			return 0
		}

	// Menu icons are owner-drawn, see menuState.setIcon. Both messages must report
	// back that they handled it, or Windows falls back to leaving a blank column.
	case wmMeasureItem:
		if w.menu != nil && w.menu.measureMenuIcon(lParam) {
			return 1
		}

	case wmDrawItem:
		if w.menu != nil && w.menu.drawMenuIcon(lParam) {
			return 1
		}

	case wmEnterSizeMove:
		w.resizing = true
		return 0

	case wmExitSizeMove:
		w.resizing = false
		w.fitContent()
		return 0

	case wmSetFocus:
		w.processFocused(true)
		return 0

	case wmKillFocus:
		w.processFocused(false)
		return 0

	case wmMouseMove:
		if !w.mouseTracking {
			trackMouseLeave(w.hwnd)
			w.mouseTracking = true
		}
		w.processMouseMoved(float64(loWord(lParam)), float64(hiWord(lParam)))
		return 0

	case wmMouseLeave:
		w.mouseTracking = false
		w.mouseOut()
		return 0

	case wmLButtonDown:
		w.mouseClicked(desktop.MouseButtonPrimary, press)
		return 0
	case wmLButtonUp:
		w.mouseClicked(desktop.MouseButtonPrimary, release)
		return 0
	case wmRButtonDown:
		w.mouseClicked(desktop.MouseButtonSecondary, press)
		return 0
	case wmRButtonUp:
		w.mouseClicked(desktop.MouseButtonSecondary, release)
		return 0
	case wmMButtonDown:
		w.mouseClicked(desktop.MouseButtonTertiary, press)
		return 0
	case wmMButtonUp:
		w.mouseClicked(desktop.MouseButtonTertiary, release)
		return 0

	case wmMouseWheel:
		delta := float64(hiWord(wParam)) / wheelDelta
		if keyDown(vkShift) {
			// Shift turns vertical wheel motion into horizontal scrolling, the way
			// the GLFW driver and the rest of the desktop do.
			w.processMouseScrolled(delta, 0)
		} else {
			w.processMouseScrolled(0, delta)
		}
		return 0
	case wmMouseHWheel:
		w.processMouseScrolled(-float64(hiWord(wParam))/wheelDelta, 0)
		return 0

	case wmKeyDown, wmSysKeyDown:
		repeat := lParam&(1<<30) != 0
		w.processKeyPressed(keyNameFor(wParam, lParam), scanCodeFor(lParam), press, repeat)
		if message == wmSysKeyDown {
			// Let the system handle Alt+F4 and friends.
			break
		}
		return 0

	case wmKeyUp, wmSysKeyUp:
		w.processKeyPressed(keyNameFor(wParam, lParam), scanCodeFor(lParam), release, false)
		if message == wmSysKeyUp {
			break
		}
		return 0

	case wmChar:
		r := rune(wParam)
		// Non-BMP characters arrive as two WM_CHAR messages carrying a UTF-16
		// surrogate pair; recombine them before delivering a rune.
		if utf16.IsSurrogate(r) {
			if r < 0xdc00 { // high half
				w.highSurrogate = r
			} else if w.highSurrogate != 0 {
				if dr := utf16.DecodeRune(w.highSurrogate, r); dr != utf8.RuneError {
					w.processCharInput(dr)
				}
				w.highSurrogate = 0
			}
			return 0
		}
		w.highSurrogate = 0
		// Filter control characters; Return and Tab arrive as key events.
		if r >= ' ' && r != 0x7f {
			w.processCharInput(r)
		}
		return 0

	case wmSetCursor:
		// LOWORD(lParam) == HTCLIENT means the pointer is over our content, so the
		// canvas owns the cursor. Anywhere else is the frame: falling through to
		// DefWindowProc is what gives the borders and corners their resize cursors.
		if loWord(lParam) == htClient {
			w.applyCursor(w.cursor)
			return 1
		}

	case wmDropFiles:
		w.processDropped(wParam) // wParam is the HDROP
		return 0

	case wmDpiChanged:
		// The suggested frame rect in lParam is deliberately ignored. It scales the
		// current frame by the DPI ratio, but rescale already derives the new frame
		// from the canvas size times the new scale - applying both compounds, and
		// the window ends up a further 25% too big on a 96 to 120 DPI step.
		w.canvas.reloadScale()
		return 0

	case wmCommand:
		if w.invokeMenu(uint16(wParam & 0xffff)) {
			return 0
		}

	case wmFyneDo:
		// The queue is drained by the run loop; this message only wakes it.
		return 0
	}

	return defWindowProc(w.hwnd, message, wParam, lParam)
}
