//go:build !noos && !tinygo && !android && !ios && !mobile

package app

import (
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
)

func rootCacheDir(a fyne.App) string {
	desktopCache, _ := os.UserCacheDir()
	return a.(*fyneApp).appDir(filepath.Join(desktopCache, "fyne"))
}
