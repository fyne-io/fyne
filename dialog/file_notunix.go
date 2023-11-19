//go:build (!linux && !openbsd && !freebsd && !netbsd) || android || wasm || js

package dialog

func fileOpenOSOverride(*FileDialog) bool {
	return false
}

func fileSaveOSOverride(*FileDialog) bool {
	return false
}
