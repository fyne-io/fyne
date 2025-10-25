//go:build !android && !ios && !mobile && !wasm && !test_web_driver

package theme

import (
	"bufio"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	internalApp "fyne.io/fyne/v2/internal/app"
)

func setupSystemTheme(fallback fyne.Theme) fyne.Theme {
	path := filepath.Join(internalApp.RootConfigDir(), "theme.json")
	f, err := os.Open(path)
	if err != nil {
		if !os.IsNotExist(err) {
			fyne.LogError("Failed to load user theme file: "+path, err)
		}
		return nil
	}
	defer f.Close()

	th, err := fromJSONWithFallback(bufio.NewReader(f), fallback)
	if err != nil {
		fyne.LogError("Failed to parse user theme file: "+path, err)
		return nil
	}

	return th
}
