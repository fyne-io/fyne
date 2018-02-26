// +build !linux,!darwin

package driver

func oSEngineName() string {
	return oSEngineOther
}

func oSWindowInit(w *window) {
}
