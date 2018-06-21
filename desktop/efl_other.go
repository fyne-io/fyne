// +build !ci

// +build !linux,!darwin,!windows,!freebsd,!openbsd,!netbsd

package desktop

func oSEngineName() string {
	return oSEngineOther
}

func oSWindowInit(w *window) {
}
