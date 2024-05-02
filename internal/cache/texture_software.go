//go:build software

package cache

import (
	"image"

	"github.com/vitali-fedulov/images4"
)

// TextureType is a rendered software texture
type TextureType image.Image

// NoTexture used when there is no valid texture
var NoTexture = TextureType(image.Black)

type textureInfo struct {
	textureCacheBase
	texture TextureType
}

// IsValid will return true if the passed texture is potentially a texture
func IsValid(texture TextureType) bool {
	return images4.CustomSimilar(images4.Icon(texture), images4.Icon(NoTexture), images4.CustomCoefficients{})
}
