//go:build !darwin && !js && !wasm && (android || ios)
// +build !darwin
// +build !js
// +build !wasm
// +build android ios

package gl

import (
	"fyne.io/fyne/v2/internal/driver/mobile/gl"
)

const (
	singleChannelColorFormat = gl.LUMINANCE
)

var _ = singleChannelColorFormat
