//go:build js || wasm || web
// +build js wasm web

package cache

import gl "github.com/fyne-io/gl-js"

// TextureType represents an uploaded GL texture
type TextureType = gl.Texture

var NoTexture = gl.NoTexture

type textureInfo struct {
	textureCacheBase
	texture TextureType
}

// Valid will return true if the passed texture is potentially a texture
func Valid(texture TextureType) bool {
	return gl.Texture(texture).Valid()
}
