//go:build !ci && wasm

package app

import (
	"syscall/js"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

func defaultVariant() fyne.ThemeVariant {
	matches := js.Global().Call("matchMedia", "(prefers-color-scheme: dark)")
	if matches.Truthy() {
		if matches.Get("matches").Bool() {
			return theme.VariantDark
		}
		return theme.VariantLight
	}
	return theme.VariantDark
}
