//go:build android

package mobile

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver"
	"fyne.io/fyne/v2/internal/driver/mobile/app"
)

// Assert we are satisfying the NativeWindow interface
var _ fyne.NativeWindow = (*window)(nil)

func (w *window) RunNative(f func(context any) error) error {
	return app.RunOnJVM(func(vm, env, ctx uintptr) error {
		// TODO: define driver.AndroidWindowContext that also includes the View
		data := &driver.AndroidContext{VM: vm, Env: env, Ctx: ctx}
		return f(data)
	})
}
