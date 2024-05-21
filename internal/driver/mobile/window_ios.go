//go:build ios

package mobile

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver"
)

// Assert we are satisfying the NativeWindow interface
var _ fyne.NativeWindow = (*window)(nil)

func (w *window) RunNative(fn func(context any) error) error {
	return fn(&driver.UnknownContext{})
}
