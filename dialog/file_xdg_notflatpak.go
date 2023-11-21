//go:build !flatpak && (linux || openbsd || freebsd || netbsd) && !android && !wasm && !js
// +build !flatpak
// +build linux openbsd freebsd netbsd
// +build !android
// +build !wasm
// +build !js

package dialog

func fileOpenOSOverride(d *FileDialog) bool {
	return false
}

func fileSaveOSOverride(d *FileDialog) bool {
	return false
}
