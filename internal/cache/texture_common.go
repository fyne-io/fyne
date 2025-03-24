package cache

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/async/migration"
)

var (
	textTextures   = make(map[FontCacheEntry]*textureInfo)
	objectTextures = make(map[fyne.CanvasObject]*textureInfo)
	texturesLock   migration.Mutex
)

// DeleteTexture deletes the texture from the cache map.
func DeleteTexture(obj fyne.CanvasObject) {
	texturesLock.Lock()
	defer texturesLock.Unlock()

	delete(objectTextures, obj)
}

// GetTextTexture gets cached texture for a text run.
func GetTextTexture(ent FontCacheEntry) (TextureType, bool) {
	texturesLock.Lock()
	defer texturesLock.Unlock()

	texInfo, ok := textTextures[ent]
	if texInfo == nil || !ok {
		return NoTexture, false
	}
	texInfo.setAlive()
	return texInfo.texture, true
}

// GetTexture gets cached texture.
func GetTexture(obj fyne.CanvasObject) (TextureType, bool) {
	texturesLock.Lock()
	defer texturesLock.Unlock()

	texInfo, ok := objectTextures[obj]
	if texInfo == nil || !ok {
		return NoTexture, false
	}
	texInfo.setAlive()
	return texInfo.texture, true
}

// RangeExpiredTexturesFor range over the expired textures for the specified canvas.
//
// Note: If this is used to free textures, then it should be called inside a current
// gl context to ensure textures are deleted from gl.
func RangeExpiredTexturesFor(canvas fyne.Canvas, f func(fyne.CanvasObject)) {
	now := timeNow()

	texturesLock.Lock()
	defer texturesLock.Unlock()

	for key, tinfo := range textTextures {
		// Just free text directly when that string/style combo is done.
		if tinfo.isExpired(now) && tinfo.canvas == canvas {
			delete(textTextures, key)
			tinfo.textFree()
		}
	}

	for obj, tinfo := range objectTextures {
		if tinfo.isExpired(now) && tinfo.canvas == canvas {
			f(obj)
		}
	}
}

// RangeTexturesFor range over the textures for the specified canvas.
// It will not return the texture for a `canvas.Text` as their render lifecycle is handled separately.
//
// Note: If this is used to free textures, then it should be called inside a current
// gl context to ensure textures are deleted from gl.
func RangeTexturesFor(canvas fyne.Canvas, f func(fyne.CanvasObject)) {
	texturesLock.Lock()
	defer texturesLock.Unlock()

	// Do nothing for texture cache, it lives outside the scope of an object.
	for obj, tinfo := range objectTextures {
		if tinfo.canvas == canvas {
			f(obj)
		}
	}
}

// DeleteTextTexturesFor deletes all text textures for the given canvas.
func DeleteTextTexturesFor(canvas fyne.Canvas) {
	texturesLock.Lock()
	defer texturesLock.Unlock()

	for key, tinfo := range textTextures {
		if tinfo.canvas == canvas {
			delete(textTextures, key)
			tinfo.textFree()
		}
	}
}

// SetTextTexture sets cached texture for a text run.
func SetTextTexture(ent FontCacheEntry, texture TextureType, canvas fyne.Canvas, free func()) {
	tinfo := prepareTexture(texture, canvas, free)

	texturesLock.Lock()
	defer texturesLock.Unlock()
	textTextures[ent] = tinfo
}

// SetTexture sets cached texture.
func SetTexture(obj fyne.CanvasObject, texture TextureType, canvas fyne.Canvas) {
	tinfo := prepareTexture(texture, canvas, nil)

	texturesLock.Lock()
	defer texturesLock.Unlock()
	objectTextures[obj] = tinfo
}

func prepareTexture(texture TextureType, canvas fyne.Canvas, free func()) *textureInfo {
	tinfo := &textureInfo{texture: texture, textFree: free}
	tinfo.canvas = canvas
	tinfo.setAlive()
	return tinfo
}

// textureCacheBase defines base texture cache object.
type textureCacheBase struct {
	expiringCache
	canvas fyne.Canvas
}
