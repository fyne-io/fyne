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
const void* darwinAppMenu();
const void* createDarwinMenu(const char* label);
void insertDarwinMenuItem(const void* menu, const char* label, int id, int index, bool isSeparator);
void completeDarwinMenu(void* menu, bool prepend);
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
	for i := len(main.Items) - 1; i >= 0; i-- {
		menu := main.Items[i]
		if !menu.AfterNativeMenus {
			nextItemID = addNativeMenu(w, menu, nextItemID, true)
		}
	}
	for _, menu := range main.Items {
		if menu.AfterNativeMenus {
			nextItemID = addNativeMenu(w, menu, nextItemID, false)
		}
	}
}

func addNativeMenu(w *window, menu *fyne.Menu, nextItemID int, prepend bool) int {
	createMenu := false
	for _, item := range menu.Items {
		if !item.PlaceInNativeMenu {
			createMenu = true
			break
		}
	}

	var nsMenu unsafe.Pointer
	if createMenu {
		nsMenu = C.createDarwinMenu(C.CString(menu.Label))
	}

	appMenuIndex := 1
	for _, item := range menu.Items {
		if item.PlaceInNativeMenu {
			C.insertDarwinMenuItem(
				C.darwinAppMenu(),
				C.CString(item.Label),
				C.int(nextItemID),
				C.int(appMenuIndex),
				C.bool(item.IsSeparator),
			)
			appMenuIndex++
		} else {
			C.insertDarwinMenuItem(
				nsMenu,
				C.CString(item.Label),
				C.int(nextItemID),
				C.int(-1),
				C.bool(item.IsSeparator),
			)
		}
		action := item.Action // catch
		callbacks = append(callbacks, func() { w.queueEvent(action) })
		nextItemID++
	}

	if nsMenu != nil {
		C.completeDarwinMenu(nsMenu, C.bool(prepend))
	}
	return nextItemID
}
