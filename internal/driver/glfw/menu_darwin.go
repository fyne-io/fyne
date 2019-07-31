// +build !no_native_menus

package glfw

import (
	"fyne.io/fyne"
)

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation -framework AppKit

void createDarwinMenu(const char* label);
void addDarwinMenuItem(const char* label, int id);
void completeDarwinMenu();
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

func setupNativeMenu(menu *fyne.MainMenu) {
	id := 0
	for _, menu := range menu.Items {
		C.createDarwinMenu(C.CString(menu.Label))

		for _, item := range menu.Items {
			C.addDarwinMenuItem(C.CString(item.Label), C.int(id))

			callbacks = append(callbacks, item.Action)
			id++
		}

		C.completeDarwinMenu()
	}
}
