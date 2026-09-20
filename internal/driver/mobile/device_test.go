package mobile

import (
	"io"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_device(t *testing.T) {
	dev := &device{}

	assert.Equal(t, true, dev.IsMobile())
	assert.Equal(t, false, dev.IsBrowser())
	assert.Equal(t, false, dev.HasKeyboard())
}

func Test_device_HideVirtualKeyboard_BeforeRun(t *testing.T) {
	// Discarding log output for tests
	// The following method logs an error:
	// (&device{}).HideVirtualKeyboard() // should not crash!
	log.SetOutput(io.Discard)
	t.Cleanup(func() { log.SetOutput(os.Stderr) })
	(&device{}).HideVirtualKeyboard() // should not crash!
}
