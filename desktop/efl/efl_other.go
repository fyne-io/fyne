// +build !ci

// +build !linux,!darwin,!windows

package efl

func oSEngineName() string {
	return oSEngineOther
}

func oSWindowInit(w *window) {
}
