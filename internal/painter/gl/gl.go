//go:generate fyne bundle -o shaders.go --prefix shader --package gl shaders/

package gl

import (
	"log"
	"runtime"
)

func logGLError(err uint32) {
	if err == 0 {
		return
	}

	log.Printf("Error %x in GL Renderer", err)
	_, file, line, ok := runtime.Caller(2)
	if ok {
		log.Printf("  At: %s:%d", file, line)
	}
}
