//+build android

package driver

import "fyne.io/fyne/v2/internal/driver/mobile/app"

func RunNative(fn func(interface{}) error) error {
	return app.RunOnJVM(func(vm, env, ctx uintptr) error {
		data := &AndroidContext{VM: vm, Env: env, Ctx: ctx}
		return fn(data)
	})
}
