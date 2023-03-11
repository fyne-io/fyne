//go:build android
// +build android

package mobile

/*
#cgo LDFLAGS: -landroid -llog

#include <stdlib.h>
#include <stdbool.h>

bool existsURI(uintptr_t jni_env, uintptr_t ctx, char* uriCstr);
void* openStream(uintptr_t jni_env, uintptr_t ctx, char* uriCstr);
char* readStream(uintptr_t jni_env, uintptr_t ctx, void* stream, int len, int* total);
void* saveStream(uintptr_t jni_env, uintptr_t ctx, char* uriCstr);
void writeStream(uintptr_t jni_env, uintptr_t ctx, void* stream, char* data, int len);
void closeStream(uintptr_t jni_env, uintptr_t ctx, void* stream);
*/
import "C"
import (
	"errors"
	"io"
	"os"
	"unsafe"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/driver/mobile/app"
	"fyne.io/fyne/v2/storage/repository"
)

type javaStream struct {
	stream unsafe.Pointer // java.io.InputStream
}

// Declare conformity to ReadCloser interface
var _ io.ReadCloser = (*javaStream)(nil)

func (s *javaStream) Read(p []byte) (int, error) {
	count := 0
	err := app.RunOnJVM(func(_, env, ctx uintptr) error {
		cCount := C.int(0)
		cBytes := unsafe.Pointer(C.readStream(C.uintptr_t(env), C.uintptr_t(ctx), s.stream, C.int(len(p)), &cCount))
		if cCount == -1 {
			return io.EOF
		}
		defer C.free(cBytes)
		count = int(cCount) // avoid sending -1 instead of 0 on completion

		bytes := C.GoBytes(cBytes, cCount)
		for i := 0; i < int(count); i++ {
			p[i] = bytes[i]
		}
		return nil
	})

	return int(count), err
}

func (s *javaStream) Close() error {
	app.RunOnJVM(func(_, env, ctx uintptr) error {
		C.closeStream(C.uintptr_t(env), C.uintptr_t(ctx), s.stream)

		return nil
	})

	return nil
}

func openStream(uri string) unsafe.Pointer {
	uriStr := C.CString(uri)
	defer C.free(unsafe.Pointer(uriStr))

	var stream unsafe.Pointer
	app.RunOnJVM(func(_, env, ctx uintptr) error {
		streamPtr := C.openStream(C.uintptr_t(env), C.uintptr_t(ctx), uriStr)
		if streamPtr == C.NULL {
			return os.ErrNotExist
		}

		stream = unsafe.Pointer(streamPtr)
		return nil
	})
	return stream
}

func nativeFileOpen(f *fileOpen) (io.ReadCloser, error) {
	if f.uri == nil || f.uri.String() == "" {
		return nil, nil
	}

	ret := openStream(f.uri.String())
	if ret == nil {
		return nil, errors.New("resource not found at URI")
	}

	stream := &javaStream{}
	stream.stream = ret
	return stream, nil
}

func saveStream(uri string) unsafe.Pointer {
	uriStr := C.CString(uri)
	defer C.free(unsafe.Pointer(uriStr))

	var stream unsafe.Pointer
	app.RunOnJVM(func(_, env, ctx uintptr) error {
		streamPtr := C.saveStream(C.uintptr_t(env), C.uintptr_t(ctx), uriStr)
		if streamPtr == C.NULL {
			return os.ErrNotExist
		}

		stream = unsafe.Pointer(streamPtr)
		return nil
	})
	return stream
}

func nativeFileSave(f *fileSave) (io.WriteCloser, error) {
	if f.uri == nil || f.uri.String() == "" {
		return nil, nil
	}

	ret := saveStream(f.uri.String())
	if ret == nil {
		return nil, errors.New("resource not found at URI")
	}

	stream := &javaStream{}
	stream.stream = ret
	return stream, nil
}

// Declare conformity to WriteCloser interface
var _ io.WriteCloser = (*javaStream)(nil)

func (s *javaStream) Write(p []byte) (int, error) {
	err := app.RunOnJVM(func(_, env, ctx uintptr) error {
		C.writeStream(C.uintptr_t(env), C.uintptr_t(ctx), s.stream, (*C.char)(C.CBytes(p)), C.int(len(p)))
		return nil
	})

	return len(p), err
}

func existsURI(uri fyne.URI) (bool, error) {
	uriStr := C.CString(uri.String())
	defer C.free(unsafe.Pointer(uriStr))

	ok := false
	app.RunOnJVM(func(_, env, ctx uintptr) error {
		ok = bool(C.existsURI(C.uintptr_t(env), C.uintptr_t(ctx), uriStr))
		return nil
	})

	return ok, nil
}

func registerRepository(d *mobileDriver) {
	repo := &mobileFileRepo{}
	repository.Register("file", repo)
	repository.Register("content", repo)
}
