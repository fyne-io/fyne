//go:build ios

package mobile

import (
	"fyne.io/fyne/v2/driver"
)

// Assert we are satisfying the driver.NativeWindow interface
var _ driver.NativeWindow = (*window)(nil)

func (w *window) RunNative(fn func(context any)) {
	fn(&driver.UnknownContext{})
}
