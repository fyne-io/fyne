// +build !ios,!android

package gomobile

import (
	"errors"
	"io"
	"os"
)

func nativeFileOpen(uri string) (io.ReadCloser, error) {
	if len(uri) < 8 || uri[:7] != "file://" {
		return nil, errors.New("Mobile simulator mode only supports file:// URIs")
	}

	return os.Open(uri[7:])
}

func nativeFileSave(uri string) (io.WriteCloser, error) {
	if len(uri) < 8 || uri[:7] != "file://" {
		return nil, errors.New("Mobile simulator mode only supports file:// URIs")
	}

	return os.Open(uri[7:])
}
