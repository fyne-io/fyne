// +build android

package gomobile

/*
#cgo LDFLAGS: -landroid -llog -lEGL -lGLESv2

#include <stdbool.h>
#include <stdlib.h>

bool canListURI(uintptr_t jni_env, uintptr_t ctx, char* uriCstr);
char *listURI(uintptr_t jni_env, uintptr_t ctx, char* uriCstr);
*/
import "C"
import (
	"strings"
	"unsafe"

	"github.com/fyne-io/mobile/app"

	"fyne.io/fyne"
	"fyne.io/fyne/storage"
)

func canListURI(uri fyne.URI) bool {
	uriStr := C.CString(uri.String())
	defer C.free(unsafe.Pointer(uriStr))
	listable := false

	app.RunOnJVM(func(_, env, ctx uintptr) error {
		listable = bool(C.canListURI(C.uintptr_t(env), C.uintptr_t(ctx), uriStr))
		return nil
	})

	return listable
}

func listURI(uri fyne.URI) ([]fyne.URI, error) {
	uriStr := C.CString(uri.String())
	defer C.free(unsafe.Pointer(uriStr))

	var str *C.char
	app.RunOnJVM(func(_, env, ctx uintptr) error {
		str = C.listURI(C.uintptr_t(env), C.uintptr_t(ctx), uriStr)
		return nil
	})

	parts := strings.Split(C.GoString(str), "|")
	var list []fyne.URI
	for _, part := range parts {
		list = append(list, storage.NewURI(part))
	}
	return list, nil
}
