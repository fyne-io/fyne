//go:build !ci && !wasm && !test_web_driver && !mobile && !tinygo

package app

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation

#include <stdbool.h>
#include <stdlib.h>

bool isBundled();
void sendNotification(char *title, char *content);
bool scheduleNotification(char *id, char *title, char *content, double seconds);
void cancelScheduledNotification(char *id);
*/
import "C"

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"time"
	"unsafe"

	"fyne.io/fyne/v2"
)

func (a *fyneApp) SendNotification(n *fyne.Notification) {
	if C.isBundled() {
		titleStr := C.CString(n.Title)
		defer C.free(unsafe.Pointer(titleStr))
		contentStr := C.CString(n.Content)
		defer C.free(unsafe.Pointer(contentStr))

		C.sendNotification(titleStr, contentStr)
		return
	}

	fallbackNotification(n.Title, n.Content)
}

func (a *fyneApp) ScheduleNotification(n *fyne.Notification, when time.Time) (*fyne.ScheduledNotification, error) {
	if !C.isBundled() {
		// osascript path has no native scheduler - use the in-process fallback.
		return a.scheduleViaScheduler(n, when)
	}

	delay := time.Until(when).Seconds()
	if delay <= 0 {
		return nil, errors.New("scheduled delivery time must be in the future")
	}

	id, err := newDarwinNotificationID()
	if err != nil {
		return nil, err
	}
	idStr := C.CString(id)
	defer C.free(unsafe.Pointer(idStr))
	titleStr := C.CString(n.Title)
	defer C.free(unsafe.Pointer(titleStr))
	contentStr := C.CString(n.Content)
	defer C.free(unsafe.Pointer(contentStr))

	if !bool(C.scheduleNotification(idStr, titleStr, contentStr, C.double(delay))) {
		// older SDK or runtime refusal - use the in-process scheduler instead
		return a.scheduleViaScheduler(n, when)
	}
	return fyne.NewScheduledNotification(id, n, when), nil
}

func (a *fyneApp) CancelScheduledNotification(id string) error {
	if !C.isBundled() {
		return a.cancelViaScheduler(id)
	}

	idStr := C.CString(id)
	defer C.free(unsafe.Pointer(idStr))
	C.cancelScheduledNotification(idStr)
	return nil
}

func newDarwinNotificationID() (string, error) {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	return "fyne-sched-" + hex.EncodeToString(b[:]), nil
}

func escapeNotificationString(in string) string {
	noSlash := strings.ReplaceAll(in, "\\", "\\\\")
	return strings.ReplaceAll(noSlash, "\"", "\\\"")
}

//export fallbackSend
func fallbackSend(cTitle, cContent *C.char) {
	title := C.GoString(cTitle)
	content := C.GoString(cContent)
	fallbackNotification(title, content)
}

func fallbackNotification(title, content string) {
	template := `display notification "%s" with title "%s"`
	script := fmt.Sprintf(template, escapeNotificationString(content), escapeNotificationString(title))

	err := exec.Command("osascript", "-e", script).Start()
	if err != nil {
		fyne.LogError("Failed to launch darwin notify script", err)
	}
}
