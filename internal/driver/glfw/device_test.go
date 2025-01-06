//go:build !wasm

package glfw

import (
	"testing"

	"fyne.io/fyne/v2"
	"github.com/stretchr/testify/assert"
)

func Test_Device(t *testing.T) {
	dev := &glDevice{}

	assert.False(t, dev.IsMobile())
	assert.Equal(t, fyne.OrientationHorizontalLeft, dev.Orientation())
	assert.True(t, dev.HasKeyboard())
	assert.False(t, dev.IsBrowser())
}
