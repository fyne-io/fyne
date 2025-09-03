//go:build !ci && ios && !mobile

package app

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation -framework UIKit -framework UserNotifications

#include <stdlib.h>

void openURL(char *urlStr);
void sendNotification(char *title, char *content);
*/
import "C"

import (
	"net/url"
	"unsafe"
)

func (a *fyneApp) OpenURL(url *url.URL) error {
	urlStr := C.CString(url.String())
	C.openURL(urlStr)
	C.free(unsafe.Pointer(urlStr))

	return nil
}
