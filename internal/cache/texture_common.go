package cache

import (
	"sync"

	"fyne.io/fyne/v2"
)

var texturesLock sync.RWMutex
var textures = make(map[fyne.CanvasObject]*textureInfo, 1024)

// DestroyTexture destroys cached texture.
func DestroyTexture(obj fyne.CanvasObject) {
	texturesLock.Lock()
	delete(textures, obj)
	texturesLock.Unlock()
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
func SetTexture(obj fyne.CanvasObject, texture TextureType, freeFn func(obj fyne.CanvasObject)) {
	texInfo := &textureInfo{texture: texture}
	texInfo.freeFn = freeFn
	texInfo.setAlive()
	texturesLock.Lock()
	textures[obj] = texInfo
	texturesLock.Unlock()
}
