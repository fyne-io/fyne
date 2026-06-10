//go:build !ci && android

package app

/*
#cgo LDFLAGS: -landroid -llog

#include <stdbool.h>
#include <stdlib.h>

void openURL(uintptr_t java_vm, uintptr_t jni_env, uintptr_t ctx, char *url);
void sendNotification(uintptr_t java_vm, uintptr_t jni_env, uintptr_t ctx, char *title, char *content);
bool scheduleNotification(uintptr_t java_vm, uintptr_t jni_env, uintptr_t ctx,
	char *id, char *title, char *body, long long deliveryMillis);
void cancelScheduledNotification(uintptr_t java_vm, uintptr_t jni_env, uintptr_t ctx, char *id);
*/
import "C"

import (
	"errors"
	"net/url"
	"time"
	"unsafe"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/driver/mobile/app"
	"fyne.io/fyne/v2/internal/scheduler"
)

func (a *fyneApp) OpenURL(url *url.URL) error {
	urlStr := C.CString(url.String())
	defer C.free(unsafe.Pointer(urlStr))

	app.RunOnJVM(func(vm, env, ctx uintptr) error {
		C.openURL(C.uintptr_t(vm), C.uintptr_t(env), C.uintptr_t(ctx), urlStr)
		return nil
	})
	return nil
}

func (a *fyneApp) SendNotification(n *fyne.Notification) {
	titleStr := C.CString(n.Title)
	defer C.free(unsafe.Pointer(titleStr))
	contentStr := C.CString(n.Content)
	defer C.free(unsafe.Pointer(contentStr))

	app.RunOnJVM(func(vm, env, ctx uintptr) error {
		C.sendNotification(C.uintptr_t(vm), C.uintptr_t(env), C.uintptr_t(ctx), titleStr, contentStr)
		return nil
	})
}

// ScheduleNotification posts a notification via Android's AlarmManager, which
// survives the app process being killed. This requires the Fyne packaging tool
// to register the FyneNotificationReceiver in AndroidManifest.xml; if it is
// missing the Java bridge returns false and we fall back to the in-process
// scheduler with cache persistence.
func (a *fyneApp) ScheduleNotification(n *fyne.Notification, when time.Time) (*fyne.ScheduledNotification, error) {
	if !when.After(time.Now()) {
		return nil, errors.New("scheduled delivery time must be in the future")
	}

	id, err := scheduler.NewID()
	if err != nil {
		return nil, err
	}
	idStr := C.CString(id)
	defer C.free(unsafe.Pointer(idStr))
	titleStr := C.CString(n.Title)
	defer C.free(unsafe.Pointer(titleStr))
	bodyStr := C.CString(n.Content)
	defer C.free(unsafe.Pointer(bodyStr))

	deliveryMillis := C.longlong(when.UnixMilli())
	var ok bool
	app.RunOnJVM(func(vm, env, ctx uintptr) error {
		ok = bool(C.scheduleNotification(C.uintptr_t(vm), C.uintptr_t(env), C.uintptr_t(ctx),
			idStr, titleStr, bodyStr, deliveryMillis))
		return nil
	})

	if !ok {
		return a.scheduleViaScheduler(n, when)
	}
	return fyne.NewScheduledNotification(id, n, when), nil
}

func (a *fyneApp) CancelScheduledNotification(id string) error {
	idStr := C.CString(id)
	defer C.free(unsafe.Pointer(idStr))

	app.RunOnJVM(func(vm, env, ctx uintptr) error {
		C.cancelScheduledNotification(C.uintptr_t(vm), C.uintptr_t(env), C.uintptr_t(ctx), idStr)
		return nil
	})

	// Also cancel any in-process schedule with the same ID, in case this app
	// previously fell back to scheduleViaScheduler before the manifest was wired.
	return a.cancelViaScheduler(id)
}
