//go:build android

package mobile

import driverDefs "fyne.io/fyne/v2/driver"

/*
#include <stdbool.h>
#include <stdlib.h>

void keepScreenOn(uintptr_t jni_env, uintptr_t ctx, bool disabled);
*/
import "C"

func setDisableScreenBlank(disable bool) {
	driverDefs.RunNative(func(ctx any) error {
		ac := ctx.(*driverDefs.AndroidContext)

		C.keepScreenOn(C.uintptr_t(ac.Env), C.uintptr_t(ac.Ctx), C.bool(disable))

		return nil
	})
}
