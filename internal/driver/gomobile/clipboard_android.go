// +build android

package gomobile

/*
#cgo LDFLAGS: -landroid -llog -lEGL -lGLESv2

#include <stdlib.h>

char *getClipboardContent(uintptr_t java_vm, uintptr_t jni_env, uintptr_t ctx);
void setClipboardContent(uintptr_t java_vm, uintptr_t jni_env, uintptr_t ctx, char *content);
*/
import "C"
import (
	"unsafe"

	"github.com/fyne-io/mobile/app"
)

// Content returns the clipboard content for Android
func (c *mobileClipboard) Content() string {
	content := ""
	app.RunOnJVM(func(vm, env, ctx uintptr) error {
		chars := C.getClipboardContent(C.uintptr_t(vm), C.uintptr_t(env), C.uintptr_t(ctx))
		if chars == nil {
			return nil
		}

		content = C.GoString(chars)
		C.free(unsafe.Pointer(chars))
		return nil
	})
	return content
}

// SetContent sets the clipboard content for Android
func (c *mobileClipboard) SetContent(content string) {
	contentStr := C.CString(content)
	defer C.free(unsafe.Pointer(contentStr))

	app.RunOnJVM(func(vm, env, ctx uintptr) error {
		C.setClipboardContent(C.uintptr_t(vm), C.uintptr_t(env), C.uintptr_t(ctx), contentStr)
		return nil
	})
}
