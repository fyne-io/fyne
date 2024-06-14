//go:build !android && !ios && !wasm && !test_web_driver

package app

import (
	"golang.org/x/sys/windows/registry"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// DefaultVariant returns the systems default fyne.ThemeVariant.
// Normally, you should not need this. It is extracted out of the root app package to give the
// settings app access to it.
func DefaultVariant() fyne.ThemeVariant {
	if isDark() {
		return theme.VariantDark
	}
	return theme.VariantLight
}

func isDark() bool {
	k, err := registry.OpenKey(registry.CURRENT_USER, `SOFTWARE\Microsoft\Windows\CurrentVersion\Themes\Personalize`, registry.QUERY_VALUE)
	if err != nil { // older version of Windows will not have this key
		return false
	}
	defer k.Close()

	useLight, _, err := k.GetIntegerValue("AppsUseLightTheme")
	if err != nil { // older version of Windows will not have this value
		return false
	}

	return useLight == 0
}
