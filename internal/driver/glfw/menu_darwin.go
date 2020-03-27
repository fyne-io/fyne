// +build !no_native_menus

package glfw

import (
	"unsafe"

	"fyne.io/fyne"
)

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation -framework AppKit

// Using void* as type for pointers is a workaround. See https://github.com/golang/go/issues/12065.
const void* createDarwinMenu(const char* label);
void insertDarwinMenuItem(const void* menu, const char* label, int id);
void completeDarwinMenu(void* menu);
*/
import "C"

var callbacks []func()

//export menu_callback
func menu_callback(id int) {
	callbacks[id]()
}

func hasNativeMenu() bool {
	return true
}

func setupNativeMenu(main *fyne.MainMenu) {
	nextItemID := 0
	for _, menu := range main.Items {
		nextItemID = addNativeMenu(menu, nextItemID)
	}
}

func addNativeMenu(menu *fyne.Menu, nextItemID int) int {
	var nsMenu unsafe.Pointer
	nsMenu = C.createDarwinMenu(C.CString(menu.Label))

	for _, item := range menu.Items {
		C.insertDarwinMenuItem(
			nsMenu,
			C.CString(item.Label),
			C.int(nextItemID),
		)
		callbacks = append(callbacks, item.Action)
		nextItemID++
	}

	C.completeDarwinMenu(nsMenu)
	return nextItemID
}
