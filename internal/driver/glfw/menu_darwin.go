// +build !no_native_menus

package glfw

import (
	"unsafe"

	"fyne.io/fyne"
)

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation -framework AppKit

#include <AppKit/AppKit.h>

// Using void* as type for pointers is a workaround. See https://github.com/golang/go/issues/12065.
void assignDarwinSubmenu(const void*, const void*);
void completeDarwinMenu(void* menu, bool prepend);
const void* createDarwinMenu(const char* label);
const void* darwinAppMenu();
const void* insertDarwinMenuItem(const void* menu, const char* label, int id, int index, bool isSeparator);
*/
import "C"

var callbacks []func()

func addNativeMenu(w *window, menu *fyne.Menu, nextItemID int, prepend bool) int {
	for i, item := range menu.Items {
		if item.Label == "Settings" {
			C.insertDarwinMenuItem(
				C.darwinAppMenu(),
				C.CString(""),
				C.int(nextItemID),
				C.int(1),
				C.bool(true),
			)
			C.insertDarwinMenuItem(
				C.darwinAppMenu(),
				C.CString(item.Label),
				C.int(nextItemID),
				C.int(2),
				C.bool(false),
			)
			nextItemID = registerCallback(w, item, nextItemID)
			if len(menu.Items) == 1 {
				return nextItemID
			}

			items := make([]*fyne.MenuItem, 0, len(menu.Items)-1)
			items = append(items, menu.Items[:i]...)
			items = append(items, menu.Items[i+1:]...)
			menu = fyne.NewMenu(menu.Label, items...)
			break
		}
	}

	nsMenu, nextItemID := createNativeMenu(w, menu, nextItemID)
	C.completeDarwinMenu(nsMenu, C.bool(prepend))
	return nextItemID
}

func addNativeSubMenu(w *window, nsParentMenuItem unsafe.Pointer, menu *fyne.Menu, nextItemID int) int {
	nsMenu, nextItemID := createNativeMenu(w, menu, nextItemID)
	C.assignDarwinSubmenu(nsParentMenuItem, nsMenu)
	return nextItemID
}

func createNativeMenu(w *window, menu *fyne.Menu, nextItemID int) (unsafe.Pointer, int) {
	nsMenu := C.createDarwinMenu(C.CString(menu.Label))
	for _, item := range menu.Items {
		nsMenuItem := C.insertDarwinMenuItem(
			nsMenu,
			C.CString(item.Label),
			C.int(nextItemID),
			C.int(-1),
			C.bool(item.IsSeparator),
		)
		nextItemID = registerCallback(w, item, nextItemID)
		if item.ChildMenu != nil {
			nextItemID = addNativeSubMenu(w, nsMenuItem, item.ChildMenu, nextItemID)
		}
	}
	return nsMenu, nextItemID
}

func registerCallback(w *window, item *fyne.MenuItem, nextItemID int) int {
	if !item.IsSeparator {
		if action := item.Action; action != nil { // catch action value
			callbacks = append(callbacks, func() { w.queueEvent(action) })
			nextItemID++
		}
	}
	return nextItemID
}

func hasNativeMenu() bool {
	return true
}

//export menuCallback
func menuCallback(id int) {
	callbacks[id]()
}

func setupNativeMenu(w *window, main *fyne.MainMenu) {
	nextItemID := 0
	var helpMenu *fyne.Menu
	for i := len(main.Items) - 1; i >= 0; i-- {
		menu := main.Items[i]
		if menu.Label == "Help" {
			helpMenu = menu
			continue
		}
		nextItemID = addNativeMenu(w, menu, nextItemID, true)
	}
	if helpMenu != nil {
		addNativeMenu(w, helpMenu, nextItemID, false)
	}
}
