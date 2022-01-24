//go:build js || wasm || test_web_driver
// +build js wasm test_web_driver

package cache

import gl "github.com/fyne-io/gl-js"

// TextureType represents an uploaded GL texture
type TextureType = gl.Texture

var NoTexture = gl.NoTexture

type textureInfo struct {
	textureCacheBase
	texture TextureType
}

// IsValid will return true if the passed texture is potentially a texture
func IsValid(texture TextureType) bool {
	return gl.Texture(texture).IsValid()
}
