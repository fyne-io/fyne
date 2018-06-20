// +build !ci

// +build !linux,!darwin,!windows

package desktop

func oSEngineName() string {
	return oSEngineOther
}

func oSWindowInit(w *window) {
}
