// +build !android,!ios,!mobile

package cache

// TextureType represents an uploaded GL texture
type TextureType = uint32

var noTexture = TextureType(0)

type textureInfo struct {
	textureCacheBase
	texture TextureType
}
