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
void        assignDarwinSubmenu(const void*, const void*);
void        completeDarwinMenu(void* menu, bool prepend);
const void* createDarwinMenu(const char* label);
const void* darwinAppMenu();
const void* insertDarwinMenuItem(const void* menu, const char* label, int id, int index, bool isSeparator);

// Used for tests.
const void* test_darwinMainMenu();
const void* test_NSMenu_itemAtIndex(const void*, NSInteger);
NSInteger   test_NSMenu_numberOfItems(const void*);
void        test_NSMenu_performActionForItemAtIndex(const void*, NSInteger);
void        test_NSMenu_removeItemAtIndex(const void* m, NSInteger i);
const char* test_NSMenu_title(const void*);
bool        test_NSMenuItem_isSeparatorItem(const void*);
const void* test_NSMenuItem_submenu(const void*);
const char* test_NSMenuItem_title(const void*);
*/
import "C"

var callbacks []func()
var ecb func(string)

func addNativeMenu(w *window, menu *fyne.Menu, nextItemID int, prepend bool) int {
	menu, nextItemID = handleSpecialItems(w, menu, nextItemID, true)

	containsItems := false
	for _, item := range menu.Items {
		if !item.IsSeparator {
			containsItems = true
			break
		}
	}
	if !containsItems {
		return nextItemID
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

//export exceptionCallback
func exceptionCallback(e *C.char) {
	msg := C.GoString(e)
	if ecb == nil {
		panic("unhandled Obj-C exception: " + msg)
	}
	ecb(msg)
}

func handleSpecialItems(w *window, menu *fyne.Menu, nextItemID int, addSeparator bool) (*fyne.Menu, int) {
	for i, item := range menu.Items {
		if item.Label == "Settings" || item.Label == "Settings…" || item.Label == "Preferences" || item.Label == "Preferences…" {
			items := make([]*fyne.MenuItem, 0, len(menu.Items)-1)
			items = append(items, menu.Items[:i]...)
			items = append(items, menu.Items[i+1:]...)
			menu, nextItemID = handleSpecialItems(w, fyne.NewMenu(menu.Label, items...), nextItemID, false)

			C.insertDarwinMenuItem(
				C.darwinAppMenu(),
				C.CString(item.Label),
				C.int(nextItemID),
				C.int(1),
				C.bool(false),
			)
			if addSeparator {
				C.insertDarwinMenuItem(
					C.darwinAppMenu(),
					C.CString(""),
					C.int(nextItemID),
					C.int(1),
					C.bool(true),
				)
			}
			nextItemID = registerCallback(w, item, nextItemID)
			break
		}
	}
	return menu, nextItemID
}

func registerCallback(w *window, item *fyne.MenuItem, nextItemID int) int {
	if !item.IsSeparator {
		callbacks = append(callbacks, func() {
			if item.Action != nil {
				w.queueEvent(item.Action)
			}
		})
		nextItemID++
	}
	return nextItemID
}

func setExceptionCallback(cb func(string)) {
	ecb = cb
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
	callbacks = []func(){}
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

//
// Test support methods
// These are needed because CGo is not supported inside test files.
//

func testDarwinMainMenu() unsafe.Pointer {
	return C.test_darwinMainMenu()
}

func testNSMenuItemAtIndex(m unsafe.Pointer, i int) unsafe.Pointer {
	return C.test_NSMenu_itemAtIndex(m, C.long(i))
}

func testNSMenuNumberOfItems(m unsafe.Pointer) int {
	return int(C.test_NSMenu_numberOfItems(m))
}

func testNSMenuPerformActionForItemAtIndex(m unsafe.Pointer, i int) {
	C.test_NSMenu_performActionForItemAtIndex(m, C.long(i))
}

func testNSMenuRemoveItemAtIndex(m unsafe.Pointer, i int) {
	C.test_NSMenu_removeItemAtIndex(m, C.long(i))
}

func testNSMenuTitle(m unsafe.Pointer) string {
	return C.GoString(C.test_NSMenu_title(m))
}

func testNSMenuItemIsSeparatorItem(i unsafe.Pointer) bool {
	return bool(C.test_NSMenuItem_isSeparatorItem(i))
}

func testNSMenuItemSubmenu(i unsafe.Pointer) unsafe.Pointer {
	return C.test_NSMenuItem_submenu(i)
}

func testNSMenuItemTitle(i unsafe.Pointer) string {
	return C.GoString(C.test_NSMenuItem_title(i))
}
