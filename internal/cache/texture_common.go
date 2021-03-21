package cache

import (
	"time"

	"fyne.io/fyne/v2"
)

// NOTE: Texture cache functions should be called always in
// the same go routine.

var textures = make(map[fyne.CanvasObject]*textureInfo, 1024)

// DeleteTexture deletes the texture from the cache map.
func DeleteTexture(obj fyne.CanvasObject) {
	delete(textures, obj)
}

// GetTexture gets cached texture.
func GetTexture(obj fyne.CanvasObject) (TextureType, bool) {
	texInfo, ok := textures[obj]
	if texInfo == nil || !ok {
		return noTexture, false
	}
	texInfo.setAlive()
	return texInfo.texture, true
}

// RangeExpiredTextures range over the expired textures.
func RangeExpiredTextures(f func(fyne.CanvasObject)) {
	now := time.Now()
	for obj, tinfo := range textures {
		if tinfo.isExpired(now) {
			f(obj)
		}
	}
}

// SetTexture sets cached texture.
func SetTexture(obj fyne.CanvasObject, texture TextureType) {
	texInfo := &textureInfo{texture: texture}
	texInfo.setAlive()
	textures[obj] = texInfo
}
