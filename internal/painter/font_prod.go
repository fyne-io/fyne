//go:build !test

package painter

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/go-text/typesetting/fontscan"
)

func loadSystemFonts(fm *fontscan.FontMap) error {
	cacheDir := ""
	if runtime.GOOS == "android" {
		parent := os.Getenv("FILESDIR")
		cacheDir = filepath.Join(parent, "fontcache")
	}

	return fm.UseSystemFonts(cacheDir)
}
