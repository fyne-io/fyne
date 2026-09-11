//go:build windows && directx

package directx

import (
	"runtime"
	"strings"
	"unsafe"

	"golang.org/x/sys/windows"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/lang"
)

// Native Win32 menus. The main menu is the real Windows menu bar - which is
// what Windows users expect, and gets keyboard access and theming for free.

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
	// icons are the item icons built for it. DestroyMenu does not own them, so
	// they are destroyed by hand when the menu is replaced or the window closes.
	iconSize int
	icons    []windows.Handle
}

// release frees the menu and the icons hanging off it.
func (s *menuState) release() {
	for _, icon := range s.icons {
		procDestroyIcon.Call(uintptr(icon))
	}
	s.icons = nil
	if s.handle != 0 {
		procDestroyMenu.Call(uintptr(s.handle))
		s.handle = 0
	}
}

func (w *window) setNativeMenu(main *fyne.MainMenu) {
	if w.hwnd == 0 {
		return
	}
	old := w.menu
	w.menu = nil

	if main == nil || len(main.Items) == 0 {
		if old != nil {
			procSetMenu.Call(uintptr(w.hwnd), 0)
			old.release()
			procDrawMenuBar.Call(uintptr(w.hwnd))
		}
		return
	}

	bar, _, _ := procCreateMenu.Call()
	if bar == 0 {
		w.menu = old // keep what is there rather than leaking a dead handle
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

	// One SetMenu call replaces any bar already attached in place. Detaching
	// first (SetMenu 0, then the new bar) resizes the client area twice, and a
	// menu Refresh - which lands here - made every window's content visibly
	// jump down and back up. The old menu is released only once replaced;
	// SetMenu detaches but never destroys.
	w.menu = state
	procSetMenu.Call(uintptr(w.hwnd), bar)
	if old != nil {
		old.release()
	}
	procDrawMenuBar.Call(uintptr(w.hwnd))

	// A menu bar takes its height out of the client area without any WM_SIZE, so
	// the window has to grow by that much or the content loses its bottom rows.
	// On a same-height rebuild setWindowPos is a no-op and Windows sends nothing.
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
//
// The icon is drawn by drawMenuIcon rather than handed over as a bitmap. Windows
// would alpha blend an hbmpItem correctly, but Wine blits it, turning every
// transparent pixel into black and the icon into a black box. Drawing it ourselves
// with DrawIconEx honours the alpha on both.
func (s *menuState) setIcon(menu, pos uintptr, item *fyne.MenuItem) {
	if item.Icon == nil {
		return
	}
	icon := iconFromResource(menuIconResource(item.Icon), s.iconSize)
	if icon == 0 {
		return
	}
	s.icons = append(s.icons, icon)

	info := menuItemInfoW{
		cbSize:     uint32(unsafe.Sizeof(menuItemInfoW{})),
		fMask:      miimBitmap | miimData,
		hbmpItem:   windows.Handle(hbmMenuCallback),
		dwItemData: uintptr(icon),
	}
	if ok, _, err := procSetMenuItemInfoW.Call(menu, pos, 1, /* by position */
		uintptr(unsafe.Pointer(&info))); ok == 0 {
		fyne.LogError("directx: could not set the icon for menu item "+item.Label, err)
	}
}

// measureMenuIcon answers WM_MEASUREITEM for an owner-drawn menu icon, claiming a
// square of the icon column's size. Reporting nothing leaves no room and the icon
// is never asked for.
func (s *menuState) measureMenuIcon(lParam uintptr) bool {
	info := (*measureItemStruct)(pointerFromAddr(lParam))
	if info.ctlType != odtMenu || info.itemData == 0 {
		return false
	}

	info.itemWidth = uint32(s.iconSize)
	info.itemHeight = uint32(s.iconSize)
	return true
}

// drawMenuIcon answers WM_DRAWITEM by painting the icon parked in the item's data
// over whatever Windows has already drawn, so the highlight shows through.
func (s *menuState) drawMenuIcon(lParam uintptr) bool {
	info := (*drawItemStruct)(pointerFromAddr(lParam))
	if info.ctlType != odtMenu || info.itemData == 0 {
		return false
	}

	// Centre the icon in the rect Windows allotted, which is at least the size
	// asked for in measureMenuIcon but can be taller on a roomy menu.
	w := int32(s.iconSize)
	if avail := info.rcItem.Right - info.rcItem.Left; avail < w {
		w = avail
	}
	h := int32(s.iconSize)
	if avail := info.rcItem.Bottom - info.rcItem.Top; avail < h {
		h = avail
	}
	x := info.rcItem.Left + (info.rcItem.Right-info.rcItem.Left-w)/2
	y := info.rcItem.Top + (info.rcItem.Bottom-info.rcItem.Top-h)/2

	procDrawIconEx.Call(uintptr(info.hDC), uintptr(x), uintptr(y), info.itemData,
		uintptr(w), uintptr(h), 0, 0, diNormal)
	return true
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

// addMissingQuitForMainMenu guarantees the first menu can quit the app: when it
// does not already end in a quit item, a separator and a Quit item are appended.
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
