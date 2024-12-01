//go:build !ci && ios && !mobile

package app

import (
	"path/filepath"
)

func rootConfigDir() string {
	root := C.documentsPath()
	return filepath.Join(C.GoString(root), "fyne")
}
