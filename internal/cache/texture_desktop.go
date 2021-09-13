//go:build !android && !ios && !mobile && !js && !wasm && !web
// +build !android,!ios,!mobile,!js,!wasm,!web

package cache

// TextureType represents an uploaded GL texture
type TextureType = uint32

// NoTexture used when there is no valid texture
var NoTexture = TextureType(0)

type textureInfo struct {
	textureCacheBase
	texture TextureType
}
