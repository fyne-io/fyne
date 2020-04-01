// +build ios

package gomobile

import "io"

func nativeFileOpen(uri string) (io.ReadCloser, error) {
	panic("Please implement me")
}

func nativeFileSave(uri string) (io.WriteCloser, error) {
	panic("Please implement me")
}
