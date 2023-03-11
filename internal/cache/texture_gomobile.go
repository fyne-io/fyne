//go:build android || ios || mobile
// +build android ios mobile

package cache

import "fyne.io/fyne/v2/internal/driver/mobile/gl"

// TextureType represents an uploaded GL texture
type TextureType = gl.Texture

var NoTexture = gl.Texture{0}

type textureInfo struct {
	textureCacheBase
	texture TextureType
}

// IsValid will return true if the passed texture is potentially a texture
func IsValid(texture TextureType) bool {
	return texture != NoTexture
}
