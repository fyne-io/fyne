//go:build darwin && !no_native_menus

package glfw

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation -framework AppKit

#include <AppKit/AppKit.h>

// Search-related functions
void        insertDarwinSearchMenuItem(const void* menu, int index);
void        searchDarwinMenus(const char* searchTerm);
void        clearDarwinMenuSearch();
void        enableDarwinHelpMenuSearch();
const void* createDarwinSearchMenuItem();
*/
import "C"
import (
	"unsafe"

	"fyne.io/fyne/v2"
)

var (
	// Store callbacks for search functionality
	onMenuSearch      func(string)
	onMenuItemFound   func(*fyne.MenuItem)
	menuItemCallbacks map[int]*fyne.MenuItem
)

func init() {
	menuItemCallbacks = make(map[int]*fyne.MenuItem)
}

// AddSearchToNativeMenu adds a search field to a native macOS menu
func AddSearchToNativeMenu(menu unsafe.Pointer, index int) {
	C.insertDarwinSearchMenuItem(menu, C.int(index))
}

// SearchNativeMenus searches through all native menus
func SearchNativeMenus(searchTerm string) {
	C.searchDarwinMenus(C.CString(searchTerm))
}

// ClearNativeMenuSearch clears the search highlighting
func ClearNativeMenuSearch() {
	C.clearDarwinMenuSearch()
}

// EnableHelpMenuSearch enables macOS's built-in Help menu search
func EnableHelpMenuSearch() {
	C.enableDarwinHelpMenuSearch()
}

// SetMenuSearchCallback sets the callback for when search text changes
func SetMenuSearchCallback(callback func(string)) {
	onMenuSearch = callback
}

// SetMenuItemFoundCallback sets the callback for when a menu item is found
func SetMenuItemFoundCallback(callback func(*fyne.MenuItem)) {
	onMenuItemFound = callback
}

// RegisterMenuItemForSearch registers a menu item for search callbacks
func RegisterMenuItemForSearch(id int, item *fyne.MenuItem) {
	menuItemCallbacks[id] = item
}

//export menuSearchCallback
func menuSearchCallback(searchTerm *C.char) {
	term := C.GoString(searchTerm)
	if onMenuSearch != nil {
		onMenuSearch(term)
	}

	// Also trigger the native search
	SearchNativeMenus(term)
}

//export menuItemFoundCallback
func menuItemFoundCallback(menuItemId C.int) {
	id := int(menuItemId)
	if item, ok := menuItemCallbacks[id]; ok {
		if onMenuItemFound != nil {
			onMenuItemFound(item)
		}
	}
}

// AddNativeSearchToFileMenu adds a native search field to the File menu
func AddNativeSearchToFileMenu(w *window, menu *fyne.Menu) {
	// This function would be called when setting up the File menu
	// to add the search field at the appropriate position

	// For now, we can use the Help menu search as it's more native to macOS
	EnableHelpMenuSearch()
}

// Alternative approach: Add search to Help menu (more macOS-like)
func SetupNativeHelpMenuSearch(w *window, mainMenu *fyne.MainMenu) {
	// Enable the native Help menu search
	EnableHelpMenuSearch()

	// Set up callbacks to handle search results
	SetMenuSearchCallback(func(term string) {
		fyne.LogError("Searching for: "+term, nil)
	})

	SetMenuItemFoundCallback(func(item *fyne.MenuItem) {
		fyne.LogError("Found menu item: "+item.Label, nil)
		// Could trigger the item's action or highlight it
		if item.Action != nil {
			item.Action()
		}
	})
}
