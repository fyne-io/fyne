//go:build !ci && !ios && !wasm && !test_web_driver && !mobile

package app

import (
	"os"
	"path/filepath"
)

func rootConfigDir() string {
	homeDir, _ := os.UserHomeDir()

	desktopConfig := filepath.Join(filepath.Join(homeDir, "Library"), "Preferences")
	return filepath.Join(desktopConfig, "fyne")
}
