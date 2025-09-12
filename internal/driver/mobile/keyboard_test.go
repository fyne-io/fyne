package mobile

import (
	"io"
	"log"
	"os"
	"testing"

	_ "fyne.io/fyne/v2/test"
)

func TestDevice_HideVirtualKeyboard_BeforeRun(t *testing.T) {
	// Discarding log output for tests
	// The following method logs an error:
	// (&device{}).hideVirtualKeyboard() // should not crash!
	log.SetOutput(io.Discard)
	t.Cleanup(func() { log.SetOutput(os.Stderr) })
	(&device{}).hideVirtualKeyboard() // should not crash!
}
