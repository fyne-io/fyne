// +build android ios mobile

package cache

import (
	"fyne.io/fyne/v2"

	"github.com/fyne-io/mobile/gl"
)

// TextureType represents an uploaded GL texture
type TextureType = gl.Texture

var noTexture = gl.Texture{0}

type textureInfo struct {
	expiringCache
	texture TextureType
	freeFn  func(obj fyne.CanvasObject)
}
