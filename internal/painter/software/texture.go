package software

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/cache"
)

var noTexture = Texture(cache.NoTexture)

// Texture represents a cached image texture
type Texture cache.TextureType

func (p *Painter) freeTexture(obj fyne.CanvasObject) {
	i, ok := cache.GetTexture(obj)
	if !ok {
		return
	}
	// TODO: this doesn't work as planned, if needed, we could probably pull it from the draw calls
	p.dirtyRects = append(p.dirtyRects, i.Bounds())
	cache.DeleteTexture(obj)
}

func (p *Painter) getTexture(object fyne.CanvasObject, creator func(canvasObject fyne.CanvasObject) Texture) (Texture, error) {
	texture, ok := cache.GetTexture(object)

	if !ok {
		texture = cache.TextureType(creator(object))
		cache.SetTexture(object, texture, p.canvas)
	}
	// TODO: get this to work
	// if !cache.IsValid(texture) {
	// 	return noTexture, fmt.Errorf("no texture available")
	// }
	return Texture(texture), nil
}
