//go:build !no_native_menus && !js && !wasm && !test_web_driver
// +build !no_native_menus,!js,!wasm,!test_web_driver

package glfw

import (
	"bytes"
	"fmt"
	"image/color"
	"image/png"
	"strings"
	"unsafe"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/internal/painter"
	"fyne.io/fyne/v2/internal/svg"
	"fyne.io/fyne/v2/theme"
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
void        getTextColorRGBA(int* r, int* g, int* b, int* a);
const void* insertDarwinMenuItem(const void* menu, const char* label, const char* keyEquivalent, unsigned int keyEquivalentModifierMask, int id, int index, bool isSeparator, const void *imageData, unsigned int imageDataLength);
int         menuFontSize();
void        resetDarwinMenu();

// Used for tests.
const void*   test_darwinMainMenu();
const void*   test_NSMenu_itemAtIndex(const void*, NSInteger);
NSInteger     test_NSMenu_numberOfItems(const void*);
void          test_NSMenu_performActionForItemAtIndex(const void*, NSInteger);
void          test_NSMenu_removeItemAtIndex(const void* m, NSInteger i);
const char*   test_NSMenu_title(const void*);
bool          test_NSMenuItem_isSeparatorItem(const void*);
const char*   test_NSMenuItem_keyEquivalent(const void*);
unsigned long test_NSMenuItem_keyEquivalentModifierMask(const void*);
const void*   test_NSMenuItem_submenu(const void*);
const char*   test_NSMenuItem_title(const void*);
*/
import "C"

type menuCallbacks struct {
	action  func()
	enabled func() bool
	checked func() bool
}

var callbacks []*menuCallbacks
var ecb func(string)
var specialKeys = map[fyne.KeyName]string{
	fyne.KeyBackspace: "\x08",
	fyne.KeyDelete:    "\x7f",
	fyne.KeyDown:      "\uf701",
	fyne.KeyEnd:       "\uf72b",
	fyne.KeyEnter:     "\x03",
	fyne.KeyEscape:    "\x1b",
	fyne.KeyF10:       "\uf70d",
	fyne.KeyF11:       "\uf70e",
	fyne.KeyF12:       "\uf70f",
	fyne.KeyF1:        "\uf704",
	fyne.KeyF2:        "\uf705",
	fyne.KeyF3:        "\uf706",
	fyne.KeyF4:        "\uf707",
	fyne.KeyF5:        "\uf708",
	fyne.KeyF6:        "\uf709",
	fyne.KeyF7:        "\uf70a",
	fyne.KeyF8:        "\uf70b",
	fyne.KeyF9:        "\uf70c",
	fyne.KeyHome:      "\uf729",
	fyne.KeyInsert:    "\uf727",
	fyne.KeyLeft:      "\uf702",
	fyne.KeyPageDown:  "\uf72d",
	fyne.KeyPageUp:    "\uf72c",
	fyne.KeyReturn:    "\n",
	fyne.KeyRight:     "\uf703",
	fyne.KeySpace:     " ",
	fyne.KeyTab:       "\t",
	fyne.KeyUp:        "\uf700",
}

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

func addNativeSubmenu(w *window, nsParentMenuItem unsafe.Pointer, menu *fyne.Menu, nextItemID int) int {
	nsMenu, nextItemID := createNativeMenu(w, menu, nextItemID)
	C.assignDarwinSubmenu(nsParentMenuItem, nsMenu)
	return nextItemID
}

func clearNativeMenu() {
	C.resetDarwinMenu()
}

func createNativeMenu(w *window, menu *fyne.Menu, nextItemID int) (unsafe.Pointer, int) {
	nsMenu := C.createDarwinMenu(C.CString(menu.Label))
	for _, item := range menu.Items {
		nsMenuItem := insertNativeMenuItem(nsMenu, item, nextItemID, -1)
		nextItemID = registerCallback(w, item, nextItemID)
		if item.ChildMenu != nil {
			nextItemID = addNativeSubmenu(w, nsMenuItem, item.ChildMenu, nextItemID)
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

			insertNativeMenuItem(C.darwinAppMenu(), item, nextItemID, 1)
			if addSeparator {
				C.insertDarwinMenuItem(
					C.darwinAppMenu(),
					C.CString(""),
					C.CString(""),
					C.uint(0),
					C.int(nextItemID),
					C.int(1),
					C.bool(true),
					unsafe.Pointer(nil),
					C.uint(0),
				)
			}
			nextItemID = registerCallback(w, item, nextItemID)
			break
		}
	}
	return menu, nextItemID
}

// TODO: theme change support, see NSSystemColorsDidChangeNotification
func insertNativeMenuItem(nsMenu unsafe.Pointer, item *fyne.MenuItem, nextItemID, index int) unsafe.Pointer {
	var imgData unsafe.Pointer
	var imgDataLength uint
	if item.Icon != nil {
		if painter.IsResourceSVG(item.Icon) {
			rsc := item.Icon
			if _, isThemed := rsc.(*theme.ThemedResource); isThemed {
				var r, g, b, a C.int
				C.getTextColorRGBA(&r, &g, &b, &a)
				rsc = &fyne.StaticResource{
					StaticName:    rsc.Name(),
					StaticContent: svg.Colorize(rsc.Content(), color.NRGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}),
				}
			}
			size := int(C.menuFontSize())
			img := painter.PaintImage(&canvas.Image{Resource: rsc}, nil, size, size)
			var buf bytes.Buffer
			if err := png.Encode(&buf, img); err != nil {
				fyne.LogError("failed to render menu icon", err)
			} else {
				imgData = unsafe.Pointer(&buf.Bytes()[0])
				imgDataLength = uint(buf.Len())
			}
		} else {
			imgData = unsafe.Pointer(&item.Icon.Content()[0])
			imgDataLength = uint(len(item.Icon.Content()))
		}
	}
	return C.insertDarwinMenuItem(
		nsMenu,
		C.CString(item.Label),
		C.CString(keyEquivalent(item)),
		C.uint(keyEquivalentModifierMask(item)),
		C.int(nextItemID),
		C.int(index),
		C.bool(item.IsSeparator),
		imgData,
		C.uint(imgDataLength),
	)
}

func keyEquivalent(item *fyne.MenuItem) (key string) {
	if s, ok := item.Shortcut.(fyne.KeyboardShortcut); ok {
		if key = specialKeys[s.Key()]; key == "" {
			if len(s.Key()) > 1 {
				fyne.LogError(fmt.Sprintf("unsupported key “%s” for menu shortcut", s.Key()), nil)
			}
			key = strings.ToLower(string(s.Key()))
		}
	}
	return
}

func keyEquivalentModifierMask(item *fyne.MenuItem) (mask uint) {
	if s, ok := item.Shortcut.(fyne.KeyboardShortcut); ok {
		if (s.Mod() & fyne.KeyModifierShift) != 0 {
			mask |= 1 << 17 // NSEventModifierFlagShift
		}
		if (s.Mod() & fyne.KeyModifierAlt) != 0 {
			mask |= 1 << 19 // NSEventModifierFlagOption
		}
		if (s.Mod() & fyne.KeyModifierControl) != 0 {
			mask |= 1 << 18 // NSEventModifierFlagControl
		}
		if (s.Mod() & fyne.KeyModifierSuper) != 0 {
			mask |= 1 << 20 // NSEventModifierFlagCommand
		}
	}
	return
}

func registerCallback(w *window, item *fyne.MenuItem, nextItemID int) int {
	if !item.IsSeparator {
		callbacks = append(callbacks, &menuCallbacks{
			action: func() {
				if item.Action != nil {
					w.QueueEvent(item.Action)
				}
			},
			enabled: func() bool {
				return !item.Disabled
			},
			checked: func() bool {
				return item.Checked
			},
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
	callbacks[id].action()
}

//export menuEnabled
func menuEnabled(id int) bool {
	return callbacks[id].enabled()
}

//export menuChecked
func menuChecked(id int) bool {
	return callbacks[id].checked()
}

func setupNativeMenu(w *window, main *fyne.MainMenu) {
	clearNativeMenu()
	nextItemID := 0
	callbacks = []*menuCallbacks{}
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

func testNSMenuItemKeyEquivalent(i unsafe.Pointer) string {
	return C.GoString(C.test_NSMenuItem_keyEquivalent(i))
}

func testNSMenuItemKeyEquivalentModifierMask(i unsafe.Pointer) uint64 {
	return uint64(C.ulong(C.test_NSMenuItem_keyEquivalentModifierMask(i)))
}

func testNSMenuItemSubmenu(i unsafe.Pointer) unsafe.Pointer {
	return C.test_NSMenuItem_submenu(i)
}

func testNSMenuItemTitle(i unsafe.Pointer) string {
	return C.GoString(C.test_NSMenuItem_title(i))
}
