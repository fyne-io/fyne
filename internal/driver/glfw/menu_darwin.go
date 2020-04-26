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

//export menuCallback
func menuCallback(id int) {
	callbacks[id]()
}

func hasNativeMenu() bool {
	return true
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

func addNativeMenu(w *window, menu *fyne.Menu, nextItemID int, prepend bool) int {
	createMenu := false
	for _, item := range menu.Items {
		if item.Label != "Settings" {
			createMenu = true
			break
		}
	}

	var nsMenu unsafe.Pointer
	if createMenu {
		nsMenu = C.createDarwinMenu(C.CString(menu.Label))
	}

	for _, item := range menu.Items {
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
		} else {
			nsMenuItem := C.insertDarwinMenuItem(
				nsMenu,
				C.CString(item.Label),
				C.int(nextItemID),
				C.int(-1),
				C.bool(item.IsSeparator),
			)
			if item.ChildMenu != nil {
				nextItemID = addNativeSubMenu(w, nsMenuItem, item.ChildMenu, nextItemID)
			}
		}
		if !item.IsSeparator {
			if action := item.Action; action != nil { // catch action value
				callbacks = append(callbacks, func() { w.queueEvent(action) })
				nextItemID++
			}
		}
	}

	if nsMenu != nil {
		C.completeDarwinMenu(nsMenu, C.bool(prepend))
	}
	return nextItemID
}

func addNativeSubMenu(w *window, nsParentMenuItem unsafe.Pointer, menu *fyne.Menu, nextItemID int) int {
	nsMenu := C.createDarwinMenu(C.CString(menu.Label))
	for _, item := range menu.Items {
		nsMenuItem := C.insertDarwinMenuItem(
			nsMenu,
			C.CString(item.Label),
			C.int(nextItemID),
			C.int(-1),
			C.bool(item.IsSeparator),
		)
		if item.ChildMenu != nil {
			nextItemID = addNativeSubMenu(w, nsMenuItem, item.ChildMenu, nextItemID)
		}
		if !item.IsSeparator {
			if action := item.Action; action != nil { // catch action value
				callbacks = append(callbacks, func() { w.queueEvent(action) })
				nextItemID++
			}
		}
	}
	C.assignDarwinSubmenu(nsParentMenuItem, nsMenu)
	return nextItemID
}
