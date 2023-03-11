//go:build ios
// +build ios

package mobile

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation

#import <stdlib.h>
#import <stdbool.h>

bool iosCanList(const char* url);
bool iosCreateListable(const char* url);
char* iosList(const char* url);
*/
import "C"
import (
	"errors"
	"strings"
	"unsafe"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"
)

func canListURI(uri fyne.URI) bool {
	uriStr := C.CString(uri.String())
	defer C.free(unsafe.Pointer(uriStr))

	return bool(C.iosCanList(uriStr))
}

func createListableURI(uri fyne.URI) error {
	uriStr := C.CString(uri.String())
	defer C.free(unsafe.Pointer(uriStr))

	ok := bool(C.iosCreateListable(uriStr))
	if ok {
		return nil
	}
	return errors.New("failed to create directory")
}

func listURI(uri fyne.URI) ([]fyne.URI, error) {
	uriStr := C.CString(uri.String())
	defer C.free(unsafe.Pointer(uriStr))

	str := C.iosList(uriStr)
	parts := strings.Split(C.GoString(str), "|")
	var list []fyne.URI
	for _, part := range parts {
		if len(part) == 0 {
			continue
		}
		list = append(list, storage.NewURI(part))
	}
	return list, nil
}
