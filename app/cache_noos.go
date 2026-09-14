//go:build noos || tinygo

package app

import (
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
)

func rootCacheDir(_ fyne.App) string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "fyne")
}
