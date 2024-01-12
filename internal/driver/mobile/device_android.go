//go:build android

package mobile

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver"
)

/*
#include <stdbool.h>
#include <stdlib.h>

void keepScreenOn(uintptr_t jni_env, uintptr_t ctx, bool disabled);
*/
import "C"

const tapYOffset = -12.0 // to compensate for how we hold our fingers on the device

func (*device) SystemScaleForWindow(_ fyne.Window) float32 {
	if currentDPI >= 600 {
		return 4
	} else if currentDPI >= 405 {
		return 3
	} else if currentDPI >= 270 {
		return 2
	} else if currentDPI >= 180 {
		return 1.5
	}
	return 1
}

func setDisableScreenBlank(disable bool) {
	driver.RunNative(func(ctx any) error {
		ac := ctx.(*driver.AndroidContext)

		C.keepScreenOn(C.uintptr_t(ac.Env), C.uintptr_t(ac.Ctx), C.bool(disable))

		return nil
	})
}
