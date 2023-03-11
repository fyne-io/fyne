//go:build ios
// +build ios

package mobile

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation

#import <stdlib.h>
#import <stdbool.h>

bool iosExistsPath(const char* path);
void* iosParseUrl(const char* url);
const void* iosReadFromURL(void* url, int* len);
const int iosWriteToURL(void* url, const void* bytes, int len);
*/
import "C"
import (
	"errors"
	"io"
	"unsafe"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage/repository"
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

type secureWriteCloser struct {
	url    unsafe.Pointer
	closer func()

	offset int
}

// Declare conformity to WriteCloser interface
var _ io.WriteCloser = (*secureWriteCloser)(nil)

func (s *secureWriteCloser) Write(p []byte) (int, error) {
	count := int(C.iosWriteToURL(s.url, C.CBytes(p), C.int(len(p))))
	s.offset += count

	return count, nil
}

func (s *secureWriteCloser) Close() error {
	if s.closer != nil {
		s.closer()
	}
	s.url = nil
	return nil
}

func existsURI(u fyne.URI) (bool, error) {
	if u.Scheme() != "file" {
		return true, errors.New("cannot check existence of " + u.Scheme() + " on iOS")
	}

	cStr := C.CString(u.Path())
	defer C.free(unsafe.Pointer(cStr))

	exists := C.iosExistsPath(cStr)
	return bool(exists), nil
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

func nativeFileSave(f *fileSave) (io.WriteCloser, error) {
	if f.uri == nil || f.uri.String() == "" {
		return nil, nil
	}

	cStr := C.CString(f.uri.String())
	defer C.free(unsafe.Pointer(cStr))

	url := C.iosParseUrl(cStr)

	fileStruct := &secureWriteCloser{url: url, closer: f.done}
	return fileStruct, nil
}

func registerRepository(d *mobileDriver) {
	repo := &mobileFileRepo{}
	repository.Register("file", repo)
}
