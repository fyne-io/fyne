// +build android

package gomobile

/*
#cgo LDFLAGS: -landroid -llog

#include <stdlib.h>

char* contentURIGetFileName(uintptr_t jni_env, uintptr_t ctx, char* uriCstr);
*/
import "C"
import (
	"path/filepath"
	"unsafe"

	"github.com/fyne-io/mobile/app"
)

func (m *mobileURI) Name() string {
	if m.Scheme() == "content" {
		result := contentURIGetFileName(m.systemURI)
		if result != "" {
			return result
		}
	}
	return m.URI.Name()
}

func (m *mobileURI) Extension() string {
	return filepath.Ext(m.Name())
}

func contentURIGetFileName(uri string) string {
	uriStr := C.CString(uri)
	defer C.free(unsafe.Pointer(uriStr))

	var filename string
	app.RunOnJVM(func(_, env, ctx uintptr) error {
		fnamePtr := C.contentURIGetFileName(C.uintptr_t(env), C.uintptr_t(ctx), uriStr)
		vPtr := unsafe.Pointer(fnamePtr)
		if vPtr == C.NULL {
			return nil
		}
		filename = C.GoString(fnamePtr)
		C.free(vPtr)

		return nil
	})
	return filename
}
