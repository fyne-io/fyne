// +build !ci,efl

package efl

func oSEngineName() string {
	return "opengl_cocoa"
}

func oSWindowInit(w *window) {
}
