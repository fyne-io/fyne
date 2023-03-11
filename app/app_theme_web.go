//go:build !ci && !js && !wasm && test_web_driver
// +build !ci,!js,!wasm,test_web_driver

package app

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

func defaultVariant() fyne.ThemeVariant {
	return theme.VariantDark
}
