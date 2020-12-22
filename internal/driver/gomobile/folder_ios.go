// +build ios

package gomobile

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation

#import <stdlib.h>
#import <stdbool.h>

bool iosCanList(const char* url);
char* iosList(const char* url);
*/
import "C"
import (
	"strings"
	"unsafe"

	"fyne.io/fyne"
	"fyne.io/fyne/storage"
)

func canListURI(uri fyne.URI) bool {
	uriStr := C.CString(uri.String())
	defer C.free(unsafe.Pointer(uriStr))

	return bool(C.iosCanList(uriStr))
}

func listURI(uri fyne.URI) ([]fyne.URI, error) {
	uriStr := C.CString(uri.String())
	defer C.free(unsafe.Pointer(uriStr))

	str := C.iosList(uriStr)
	parts := strings.Split(C.GoString(str), "|")
	var list []fyne.URI
	for _, part := range parts {
		list = append(list, storage.NewURI(part))
	}
	return list, nil
}
