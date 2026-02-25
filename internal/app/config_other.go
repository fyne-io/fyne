//go:build ci || (mobile && !android && !ios) || (!linux && !darwin && !windows && !freebsd && !openbsd && !netbsd && !wasm && !test_web_driver && !noos && !tinygo)

package app

import (
	"os"
	"path/filepath"
)

func rootConfigDir() string {
	return filepath.Join(os.TempDir(), "fyne-test")
}
