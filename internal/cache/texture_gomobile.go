// +build android ios mobile

package cache

import "github.com/fyne-io/mobile/gl"

// TextureType represents an uploaded GL texture
type TextureType = gl.Texture

var noTexture = gl.Texture{0}

type textureInfo struct {
	expiringCache
	texture TextureType
}
