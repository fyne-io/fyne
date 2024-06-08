//go:build android

package mobile

import (
	"fyne.io/fyne/v2/driver"
	"fyne.io/fyne/v2/internal/driver/mobile/app"
)

// Assert we are satisfying the driver.NativeWindow interface
var _ driver.NativeWindow = (*window)(nil)

func (w *window) RunNative(f func(context any) error) error {
	return app.RunOnJVM(func(vm, env, ctx uintptr) error {
		data := &driver.AndroidWindowContext{
			NativeWindow: w.handle,
			AndroidContext: driver.AndroidContext{
				VM:  vm,
				Env: env,
				Ctx: ctx,
			},
		}
		return f(data)
	})
}
