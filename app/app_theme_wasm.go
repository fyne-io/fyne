//go:build !ci && wasm
// +build !ci,wasm

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

func init() {
	if matchMedia := js.Global().Call("matchMedia", "(prefers-color-scheme: dark)"); matchMedia.Truthy() {
		matchMedia.Call("addEventListener", "change", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			fyne.CurrentApp().Settings().(*settings).setupTheme()
			return nil
		}))
	}
}
