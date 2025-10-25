//go:build !noos && !tinygo

package app

import (
	"os"
	"path/filepath"
)

func rootCacheDir() string {
	desktopCache, _ := os.UserCacheDir()
	return filepath.Join(desktopCache, "fyne")
}
