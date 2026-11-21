//go:build !noos && !tinygo && !android && !ios && !mobile && !wasm

package app

import (
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"
)

func rootCacheDir(a fyne.App) string {
	desktopCache, _ := os.UserCacheDir()
	return filepath.Join(desktopCache, "fyne", a.UniqueID())
}

func rootCacheURI(a fyne.App) fyne.URI {
	return storage.NewFileURI(rootCacheDir(a))
}
