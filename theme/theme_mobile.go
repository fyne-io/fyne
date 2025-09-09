//go:build android || ios || mobile

package theme

import "fyne.io/fyne/v2"

func setupSystemTheme(fallback fyne.Theme) fyne.Theme {
	return fallback
}
