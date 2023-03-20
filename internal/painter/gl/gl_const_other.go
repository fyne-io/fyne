//go:build !darwin && !js && !wasm
// +build !darwin,!js,!wasm

package gl

import (
	"fyne.io/fyne/v2/internal/driver/mobile/gl"
)

const (
	singleChannelColorFormat = gl.LUMINANCE
)

var _ = singleChannelColorFormat
