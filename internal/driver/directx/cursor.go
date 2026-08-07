//go:build windows && directx

// Mapping from Fyne cursors to Win32 cursors: stock cursors by id, custom
// desktop.Cursor images through CreateIconIndirect.

package directx

import (
	"fyne.io/fyne/v2/driver/desktop"
)

// cursorFor maps a Fyne cursor onto a stock Win32 cursor. The second result is
// false for HiddenCursor, which Win32 expresses as a null cursor handle rather
// than a stock id.
func cursorFor(c desktop.Cursor) (id uint16, visible bool) {
	switch c {
	case desktop.TextCursor:
		return idcIBeam, true
	case desktop.CrosshairCursor:
		return idcCross, true
	case desktop.PointerCursor:
		return idcHand, true
	case desktop.VResizeCursor:
		return idcSizeNS, true
	case desktop.HResizeCursor:
		return idcSizeWE, true
	case desktop.NESWResizeCursor:
		return idcSizeNESW, true
	case desktop.NWSEResizeCursor:
		return idcSizeNWSE, true
	case desktop.HiddenCursor:
		return 0, false
	default:
		return idcArrow, true
	}
}

// applyCursor makes the given Fyne cursor current. A non-standard cursor
// providing an image becomes a real HCURSOR, cached on the window until the
// cursor changes; everything else maps to a stock cursor, and Win32 hides the
// pointer when SetCursor is passed a null handle.
func (w *window) applyCursor(c desktop.Cursor) {
	if c == nil { // WM_SETCURSOR arrives before any mouse move has set a cursor
		c = desktop.DefaultCursor
	}
	if _, standard := c.(desktop.StandardCursor); !standard {
		if img, hotX, hotY := c.Image(); img != nil {
			if w.customCursorFor != c {
				if h := hCursorFor(img, hotX, hotY); h != 0 {
					w.freeCustomCursor()
					w.customCursor, w.customCursorFor = h, c
				}
			}
			if w.customCursor != 0 {
				setCursor(w.customCursor)
				return
			}
		}
	}

	w.freeCustomCursor()
	id, visible := cursorFor(c)
	if !visible {
		setCursor(0)
		return
	}
	setCursor(loadCursor(id))
}

func (w *window) freeCustomCursor() {
	if w.customCursor == 0 {
		return
	}
	procDestroyCursor.Call(uintptr(w.customCursor))
	w.customCursor = 0
	w.customCursorFor = nil
}
