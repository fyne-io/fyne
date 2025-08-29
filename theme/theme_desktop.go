//go:build !android && !ios && !mobile && !wasm && !test_web_driver

package theme

import (
	"bytes"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	internalApp "fyne.io/fyne/v2/internal/app"
)

func setupSystemTheme(fallback fyne.Theme) fyne.Theme {
	root := internalApp.RootConfigDir()

	path := filepath.Join(root, "theme.json")
	data, err := fyne.LoadResourceFromPath(path)
	if err != nil {
		if !os.IsNotExist(err) {
			fyne.LogError("Failed to load user theme file: "+path, err)
		}
		return nil
	}
	if data != nil && data.Content() != nil {
		th, err := fromJSONWithFallback(bytes.NewReader(data.Content()), fallback)
		if err == nil {
			return th
		}
		fyne.LogError("Failed to parse user theme file: "+path, err)
	}
	return nil
}
