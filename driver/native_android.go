//go:build android
// +build android

package driver

import "fyne.io/fyne/v2/internal/driver/mobile/app"

// RunNative provides a way to execute code within the platform-specific runtime context for various runtimes.
// On Android this provides the JVM pointers required to execute various NDK calls or use JNI APIs.
//
// Since: 2.3
func RunNative(fn func(interface{}) error) error {
	return app.RunOnJVM(func(vm, env, ctx uintptr) error {
		data := &AndroidContext{VM: vm, Env: env, Ctx: ctx}
		return fn(data)
	})
}
