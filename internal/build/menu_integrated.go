//go:build !darwin || no_native_menus

package build

// HasNativeMenu is true if the app is built with support for native menu.
const HasNativeMenu = false
