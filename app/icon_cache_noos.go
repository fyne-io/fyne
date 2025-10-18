//go:build noos || tinygo

package app

import (
	"os"
	"path/filepath"
)

func rootCacheDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "fyne")
}
