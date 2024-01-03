//go:build !flatpak && !windows && !android && !ios && !wasm && !js

package dialog

func fileOpenOSOverride(d *FileDialog) bool {
	return false
}

func fileSaveOSOverride(d *FileDialog) bool {
	return false
}
