package cache

import (
	"sync"
	"time"

	"fyne.io/fyne/v2"
)

var texturesLock sync.RWMutex
var textures = make(map[fyne.CanvasObject]*textureInfo, 1024)

// DeleteTexture deletes the texture from the cache map.
func DeleteTexture(obj fyne.CanvasObject) {
	texturesLock.Lock()
	delete(textures, obj)
	texturesLock.Unlock()
}

// RangeExpiredTextures range over the expired textures.
func RangeExpiredTextures(f func(fyne.CanvasObject)) {
	now := time.Now()
	expired := make([]fyne.CanvasObject, 0, 10)
	texturesLock.RLock()
	for obj, tinfo := range textures {
		if tinfo.isExpired(now) {
			expired = append(expired, obj)
		}
	}
	texturesLock.RUnlock()
	for _, obj := range expired {
		f(obj)
	}
}

// GetTexture gets cached texture.
func GetTexture(obj fyne.CanvasObject) (TextureType, bool) {
	texturesLock.RLock()
	texInfo, ok := textures[obj]
	texturesLock.RUnlock()
	if texInfo == nil || !ok {
		return noTexture, false
	}
	texInfo.setAlive()
	return texInfo.texture, true
}

// SetTexture sets cached texture.
func SetTexture(obj fyne.CanvasObject, texture TextureType) {
	texInfo := &textureInfo{texture: texture}
	texInfo.setAlive()
	texturesLock.Lock()
	textures[obj] = texInfo
	texturesLock.Unlock()
}
