//go:build noos || tinygo

package app

import (
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"
)

func rootCacheDir(_ fyne.App) string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "fyne")
}

func rootCacheURI(a fyne.App) fyne.URI {
	return storage.NewFileURI(rootCacheDir(a))
}
