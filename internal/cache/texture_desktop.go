//go:build !android && !ios && !mobile && !js && !wasm && !test_web_driver
// +build !android,!ios,!mobile,!js,!wasm,!test_web_driver

package cache

// TextureType represents an uploaded GL texture
type TextureType = uint32

// NoTexture used when there is no valid texture
var NoTexture = TextureType(0)

type textureInfo struct {
	textureCacheBase
	texture TextureType
}

// IsValid will return true if the passed texture is potentially a texture
func IsValid(texture TextureType) bool {
	return texture != NoTexture
}
