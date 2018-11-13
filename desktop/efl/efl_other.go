// +build !ci,!gl

// +build !linux,!darwin,!windows,!freebsd,!openbsd,!netbsd

package efl

func oSEngineName() string {
	return oSEngineOther
}

func oSWindowInit(w *window) {
}
