//go:generate go run ../../../cmd/fyne bundle -o shaders.go --prefix shader --package gl shaders/

package gl

import (
	"log"
	"runtime"

	"fyne.io/fyne/v2"
)

const floatSize = 4
const max16bit = float32(255 * 255)

func logGLError(err uint32) {
	if fyne.CurrentApp().Settings().BuildType() != fyne.BuildDebug {
		return
	}

	if err == 0 {
		return
	}

	log.Printf("Error %x in GL Renderer", err)
	_, file, line, ok := runtime.Caller(2)
	if ok {
		log.Printf("  At: %s:%d", file, line)
	}
}
