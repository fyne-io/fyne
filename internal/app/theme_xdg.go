//go:build !wasm && !test_web_driver && !android && !ios && !mobile && (linux || openbsd || freebsd || netbsd)

package app

import (
	"sync/atomic"

	"fyne.io/fyne/v2"
)

// CurrentVariant contains the system’s theme variant.
// It is intended for internal use, only!
var CurrentVariant atomic.Uint64

// DefaultVariant returns the systems default fyne.ThemeVariant.
// Normally, you should not need this. It is extracted out of the root app package to give the
// settings app access to it.
func DefaultVariant() fyne.ThemeVariant {
	return fyne.ThemeVariant(CurrentVariant.Load())
}
