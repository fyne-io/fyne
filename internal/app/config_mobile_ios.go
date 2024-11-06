//go:build !ci && ios && !mobile

package app

func rootConfigDir() string {
	root := C.documentsPath()
	return filepath.Join(C.GoString(root), "fyne")
}
