//go:build !ci && ios && !mobile

package app

import (
	"path/filepath"
)

/*
#include <stdlib.h>

char *documentsPath(void);
*/
import "C"

func rootConfigDir() string {
	root := C.documentsPath()
	return filepath.Join(C.GoString(root), "fyne")
}
