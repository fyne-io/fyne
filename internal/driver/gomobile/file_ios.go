// +build ios

package gomobile

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation

#import <stdlib.h>
#import <stdbool.h>

void* iosParseUrl(const char* url);
const void* iosReadFromURL(void* url, int* len);
*/
import "C"
import (
	"io"
	"unsafe"
)

type secureReadCloser struct {
	url    unsafe.Pointer
	closer func()

	data   []byte
	offset int
}

// Declare conformity to ReadCloser interface
var _ io.ReadCloser = (*secureReadCloser)(nil)

func (s *secureReadCloser) Read(p []byte) (int, error) {
	if s.data == nil {
		var length C.int
		s.data = C.GoBytes(C.iosReadFromURL(s.url, &length), length)
	}

	count := len(p)
	remain := len(s.data) - s.offset
	var err error
	if count >= remain {
		count = remain
		err = io.EOF
	}

	newOffset := s.offset + count

	o := 0
	for i := s.offset; i < newOffset; i++ {
		p[o] = s.data[i]
		o++
	}
	s.offset = newOffset
	return count, err
}

func (s *secureReadCloser) Close() error {
	if s.closer != nil {
		s.closer()
	}
	s.url = nil
	return nil
}

func nativeFileOpen(f *fileOpen) (io.ReadCloser, error) {
	if f.uri == nil || f.uri.String() == "" {
		return nil, nil
	}

	cStr := C.CString(f.uri.String())
	defer C.free(unsafe.Pointer(cStr))

	url := C.iosParseUrl(cStr)

	fileStruct := &secureReadCloser{url: url, closer: f.done}
	return fileStruct, nil
}
