//go:build !flatpak && !windows && !android && !ios && !wasm && !js

package dialog

func fileOpenOSOverride(_ *FileDialog) bool {
	return false
}

func fileSaveOSOverride(_ *FileDialog) bool {
	return false
}
