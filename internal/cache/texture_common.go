package cache

import (
	"image"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/async"
)

var (
	textTextures   async.Map[FontCacheEntry, *textureInfo]
	objectTextures async.Map[fyne.CanvasObject, *textureInfo]
)

// DeleteTexture deletes the texture from the cache map.
func DeleteTexture(obj fyne.CanvasObject) {
	objectTextures.Delete(obj)
}

// GetTextTexture gets cached texture for a text run.
func GetTextTexture(ent FontCacheEntry) (TextureType, bool) {
	texInfo, ok := textTextures.Load(ent)
	if texInfo == nil || !ok {
		return NoTexture, false
	}
	texInfo.setAlive()
	return texInfo.texture, true
}

// GetTexture gets cached texture.
func GetTexture(obj fyne.CanvasObject) (TextureType, bool) {
	texInfo, ok := objectTextures.Load(obj)
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

	textTextures.Range(func(key FontCacheEntry, tinfo *textureInfo) bool {
		// Just free text directly when that string/style combo is done.
		if tinfo.isExpired(now) && tinfo.canvas == canvas {
			textTextures.Delete(key)
			tinfo.textFree()
		}
		return true
	})

	objectTextures.Range(func(obj fyne.CanvasObject, tinfo *textureInfo) bool {
		if tinfo.isExpired(now) && tinfo.canvas == canvas {
			f(obj)
		}
		return true
	})
}

// RangeTexturesFor range over the textures for the specified canvas.
// It will not return the texture for a `canvas.Text` as their render lifecycle is handled separately.
//
// Note: If this is used to free textures, then it should be called inside a current
// gl context to ensure textures are deleted from gl.
func RangeTexturesFor(canvas fyne.Canvas, f func(fyne.CanvasObject)) {
	// Do nothing for texture cache, it lives outside the scope of an object.
	objectTextures.Range(func(obj fyne.CanvasObject, tinfo *textureInfo) bool {
		if tinfo.canvas == canvas {
			f(obj)
		}
		return true
	})
}

// DeleteTextTexturesFor deletes all text textures for the given canvas.
func DeleteTextTexturesFor(canvas fyne.Canvas) {
	textTextures.Range(func(key FontCacheEntry, tinfo *textureInfo) bool {
		if tinfo.canvas == canvas {
			textTextures.Delete(key)
			tinfo.textFree()
		}
		return true
	})
}

// SetTextureDirty notes that the pixels in region r of the image that the texture
// of obj was uploaded from have changed, adding to any region noted before.
// An object without a texture needs no note: creating one uploads every pixel.
func SetTextureDirty(obj fyne.CanvasObject, r image.Rectangle) {
	if texInfo, ok := objectTextures.Load(obj); ok && texInfo != nil {
		texInfo.dirty = texInfo.dirty.Union(r)
	}
}

// TakeTextureDirty returns the region noted by SetTextureDirty for the texture
// of obj, and forgets it.
func TakeTextureDirty(obj fyne.CanvasObject) image.Rectangle {
	texInfo, ok := objectTextures.Load(obj)
	if texInfo == nil || !ok {
		return image.Rectangle{}
	}
	dirty := texInfo.dirty
	texInfo.dirty = image.Rectangle{}
	return dirty
}

// SetTextTexture sets cached texture for a text run.
func SetTextTexture(ent FontCacheEntry, texture TextureType, canvas fyne.Canvas, free func()) {
	tinfo := prepareTexture(texture, canvas, free)
	textTextures.Store(ent, tinfo)
}

// SetTexture sets cached texture.
func SetTexture(obj fyne.CanvasObject, texture TextureType, canvas fyne.Canvas) {
	tinfo := prepareTexture(texture, canvas, nil)
	objectTextures.Store(obj, tinfo)
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
	dirty  image.Rectangle // pixels changed since the texture was uploaded
}
