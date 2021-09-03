//go:build ios
// +build ios

package mobileinit

import (
	"log"
	"unsafe"
)

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation

#include <asl.h>
#include <stdlib.h>

void log_wrap(const char *logStr);
*/
import "C"

type aslWriter struct{}

func (aslWriter) Write(p []byte) (n int, err error) {
	cstr := C.CString(string(p))
	C.log_wrap(cstr)
	C.free(unsafe.Pointer(cstr))
	return len(p), nil
}

func init() {
	log.SetOutput(aslWriter{})
}
