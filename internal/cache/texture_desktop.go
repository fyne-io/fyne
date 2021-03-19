// +build !android,!ios,!mobile

package cache

import (
	"fyne.io/fyne/v2"
)

// TextureType represents an uploaded GL texture
type TextureType = uint32

var noTexture = TextureType(0)

type textureInfo struct {
	expiringCache
	texture TextureType
	freeFn  func(obj fyne.CanvasObject)
}
