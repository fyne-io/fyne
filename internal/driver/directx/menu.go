//go:build windows

package directx

import (
	"runtime"
	"strings"
	"unsafe"

	"golang.org/x/sys/windows"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/lang"
)

// Native Win32 menus. Fyne's own MenuBar widget lives inside the GLFW driver, so
// rather than porting it this driver uses the real menu bar - which is also what
// Windows users expect, and gets keyboard access and theming for free.

var (
	procCreateMenu      = user32.NewProc("CreateMenu")
	procCreatePopupMenu = user32.NewProc("CreatePopupMenu")
	procAppendMenuW     = user32.NewProc("AppendMenuW")
	procSetMenu         = user32.NewProc("SetMenu")
	procDrawMenuBar     = user32.NewProc("DrawMenuBar")
	procDestroyMenu     = user32.NewProc("DestroyMenu")
)

const (
	wmCommand = 0x0111

	mfString    = 0x0000
	mfGrayed    = 0x0001
	mfChecked   = 0x0008
	mfPopup     = 0x0010
	mfSeparator = 0x0800

	// firstMenuID keeps generated command ids clear of the standard ones.
	firstMenuID = 0x1000
)

// menuActions maps a generated command id back to the item that owns it. Menus
// are rebuilt wholesale on SetMainMenu, so ids are never reused within a window.
type menuState struct {
	handle  windows.Handle
	actions map[uint16]*fyne.MenuItem
	nextID  uint16

	// iconSize is the pixel size of the icon column at this window's DPI, and
	// bitmaps are the item icons built for it. DestroyMenu does not own them, so
	// they are deleted by hand when the menu is replaced or the window closes.
	iconSize int
	bitmaps  []windows.Handle
}

// release frees the menu and the icon bitmaps hanging off it.
func (s *menuState) release() {
	for _, b := range s.bitmaps {
		deleteObject(b)
	}
	s.bitmaps = nil
	if s.handle != 0 {
		procDestroyMenu.Call(uintptr(s.handle))
		s.handle = 0
	}
}

func (w *window) setNativeMenu(main *fyne.MainMenu) {
	if w.hwnd == 0 {
		return
	}
	if w.menu != nil {
		procSetMenu.Call(uintptr(w.hwnd), 0)
		w.menu.release()
		w.menu = nil
	}
	if main == nil || len(main.Items) == 0 {
		procDrawMenuBar.Call(uintptr(w.hwnd))
		return
	}

	bar, _, _ := procCreateMenu.Call()
	if bar == 0 {
		return
	}
	state := &menuState{
		handle:   windows.Handle(bar),
		actions:  map[uint16]*fyne.MenuItem{},
		nextID:   firstMenuID,
		iconSize: menuIconSize(w.hwnd),
	}

	for _, m := range main.Items {
		sub := state.buildMenu(m.Items)
		if sub == 0 {
			continue
		}
		appendMenu(bar, mfPopup, sub, m.Label)
	}

	w.menu = state
	procSetMenu.Call(uintptr(w.hwnd), bar)
	procDrawMenuBar.Call(uintptr(w.hwnd))

	// A menu bar takes its height out of the client area without any WM_SIZE, so
	// the window has to grow by that much or the content loses its bottom rows.
	w.applyClientSize()
}

// buildMenu creates a popup menu for the given items, recursing into submenus.
func (s *menuState) buildMenu(items []*fyne.MenuItem) uintptr {
	h, _, _ := procCreatePopupMenu.Call()
	if h == 0 {
		return 0
	}
	pos := uintptr(0) // position of the next item, for setIcon
	for _, item := range items {
		if item.IsSeparator {
			procAppendMenuW.Call(h, mfSeparator, 0, 0)
			pos++
			continue
		}
		if item.ChildMenu != nil && len(item.ChildMenu.Items) > 0 {
			child := s.buildMenu(item.ChildMenu.Items)
			if child == 0 {
				continue
			}
			appendMenu(h, mfPopup, child, item.Label)
			s.setIcon(h, pos, item)
			pos++
			continue
		}

		flags := uintptr(mfString)
		if item.Disabled {
			flags |= mfGrayed
		}
		if item.Checked {
			flags |= mfChecked
		}
		id := s.nextID
		s.nextID++
		s.actions[id] = item
		appendMenu(h, flags, uintptr(id), menuLabel(item))
		s.setIcon(h, pos, item)
		pos++
	}
	return h
}

// setIcon puts an item's icon in the menu's icon column, the same column the
// check mark uses. Addressing by position rather than by command id is what lets
// submenu headers, which have no id, carry an icon too.
func (s *menuState) setIcon(menu, pos uintptr, item *fyne.MenuItem) {
	if item.Icon == nil {
		return
	}
	bmp := menuBitmapFromResource(item.Icon, s.iconSize)
	if bmp == 0 {
		return
	}
	s.bitmaps = append(s.bitmaps, bmp)

	info := menuItemInfoW{
		cbSize:   uint32(unsafe.Sizeof(menuItemInfoW{})),
		fMask:    miimBitmap,
		hbmpItem: bmp,
	}
	procSetMenuItemInfoW.Call(menu, pos, 1 /* by position */, uintptr(unsafe.Pointer(&info)))
}

// menuLabel renders an item's label with its accelerator after a tab, which is how
// Win32 right-aligns shortcut text in a menu. This is display only: a real Win32
// accelerator needs an HACCEL table, whereas Fyne shortcuts are dispatched from
// the key handler through triggerMainMenuShortcut.
func menuLabel(item *fyne.MenuItem) string {
	sh, ok := item.Shortcut.(fyne.KeyboardShortcut)
	if !ok {
		return item.Label
	}

	var b strings.Builder
	b.WriteString(item.Label)
	b.WriteByte('\t')
	mod := sh.Mod()
	for _, m := range [...]struct {
		bit  fyne.KeyModifier
		name string
	}{
		{fyne.KeyModifierControl, "Ctrl+"},
		{fyne.KeyModifierAlt, "Alt+"},
		{fyne.KeyModifierShift, "Shift+"},
		{fyne.KeyModifierSuper, "Win+"},
	} {
		if mod&m.bit != 0 {
			b.WriteString(m.name)
		}
	}
	b.WriteString(string(sh.Key()))
	return b.String()
}

// triggerMainMenuShortcut runs the first main menu item bound to this shortcut,
// reporting whether one claimed it.
func (w *window) triggerMainMenuShortcut(sh fyne.Shortcut) bool {
	if w.mainMenu == nil {
		return false
	}
	for _, m := range w.mainMenu.Items {
		if triggerMenuShortcut(sh, m) {
			return true
		}
	}
	return false
}

func triggerMenuShortcut(sh fyne.Shortcut, m *fyne.Menu) bool {
	for _, i := range m.Items {
		if i.Shortcut != nil && i.Shortcut.ShortcutName() == sh.ShortcutName() && i.Action != nil {
			i.Action()
			return true
		}
		if i.ChildMenu != nil && triggerMenuShortcut(sh, i.ChildMenu) {
			return true
		}
	}
	return false
}

// addMissingQuitForMainMenu guarantees the first menu can quit the app, matching
// the GLFW driver so the same menu definition behaves the same on both backends.
func addMissingQuitForMainMenu(menus *fyne.MainMenu, w *window) {
	if len(menus.Items) == 0 {
		return
	}
	first := menus.Items[0]

	localQuit := lang.L("Quit")
	var lastItem *fyne.MenuItem
	if len(first.Items) > 0 {
		lastItem = first.Items[len(first.Items)-1]
		if lastItem.Label == localQuit {
			lastItem.IsQuit = true
		}
	}
	if lastItem == nil || !lastItem.IsQuit {
		quitItem := fyne.NewMenuItem(localQuit, nil)
		quitItem.IsQuit = true
		first.Items = append(first.Items, fyne.NewMenuItemSeparator(), quitItem)
	}

	for _, item := range first.Items {
		if item.IsQuit && item.Action == nil {
			item.Action = func() {
				for _, win := range w.driver.AllWindows() {
					win.Close()
				}
			}
		}
	}
}

// invokeMenu runs the action bound to a WM_COMMAND id, reporting whether the id
// belonged to this window's menu.
func (w *window) invokeMenu(id uint16) bool {
	if w.menu == nil {
		return false
	}
	item, ok := w.menu.actions[id]
	if !ok {
		return false
	}
	if item.Action != nil {
		item.Action()
	}
	return true
}

// appendMenu adds one entry, keeping the encoded label alive across the call.
func appendMenu(menu, flags, idOrSubmenu uintptr, label string) {
	p := utf16Ptr(label)
	procAppendMenuW.Call(menu, flags, idOrSubmenu, uintptr(unsafe.Pointer(p)))
	runtime.KeepAlive(p)
}
