// +build android

package gomobile

import "io"

func nativeFileOpen(f *fileOpen) (io.ReadCloser, error) {
	panic("Please implement me")
}

func nativeFileSave(f *fileSave) (io.WriteCloser, error) {
	panic("Please implement me")
}
