//go:build !wasm

package glfw

import (
	"testing"

	"fyne.io/fyne/v2"
	"github.com/stretchr/testify/assert"
)

func Test_Device(t *testing.T) {
	dev := &glDevice{}

	assert.Equal(t, false, dev.IsMobile())
	assert.Equal(t, fyne.OrientationHorizontalLeft, dev.Orientation())
	assert.Equal(t, true, dev.HasKeyboard())
	assert.Equal(t, false, dev.IsBrowser())
}
