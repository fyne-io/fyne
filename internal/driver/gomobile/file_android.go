// +build android

package gomobile

import "io"

func nativeFileOpen(f *file) (io.ReadCloser, error) {
	panic("Please implement me")
}

func nativeFileSave(f *file) (io.WriteCloser, error) {
	panic("Please implement me")
}
