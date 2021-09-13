//go:build js || wasm || web
// +build js wasm web

package cache

import "github.com/goxjs/gl"

// TextureType represents an uploaded GL texture
type TextureType = gl.Texture

var NoTexture = gl.NoTexture

type textureInfo struct {
	textureCacheBase
	texture TextureType
}
