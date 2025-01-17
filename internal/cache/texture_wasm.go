//go:build wasm || test_web_driver

package cache

import "github.com/fyne-io/gl-js"

// TextureType represents an uploaded GL texture
type TextureType = gl.Texture

var NoTexture = gl.NoTexture

type textureInfo struct {
	textureCacheBase

	texture  TextureType
	textFree func()
}

// IsValid will return true if the passed texture is potentially a texture
func IsValid(texture TextureType) bool {
	return gl.Texture(texture).IsValid()
}
