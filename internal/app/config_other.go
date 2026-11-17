//go:build ci || (mobile && !android && !ios && !windows) || (!linux && !darwin && !windows && !freebsd && !openbsd && !netbsd && !wasm && !test_web_driver && !noos && !tinygo)

package app

import (
	"os"
	"path/filepath"
)

func RootConfigDir() string {
	return filepath.Join(os.TempDir(), "fyne-test")
}
