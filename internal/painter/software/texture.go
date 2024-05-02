package software

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/cache"
)

var noTexture = Texture(cache.NoTexture)

// Texture represents an uploaded GL texture
type Texture cache.TextureType

func (p *Painter) freeTexture(obj fyne.CanvasObject) {
	_, ok := cache.GetTexture(obj)
	if !ok {
		return
	}
	cache.DeleteTexture(obj)
}

func (p *Painter) getTexture(object fyne.CanvasObject, creator func(canvasObject fyne.CanvasObject) Texture) (Texture, error) {
	texture, ok := cache.GetTexture(object)

	if !ok {
		fmt.Println("cache miss")
		texture = cache.TextureType(creator(object))
		cache.SetTexture(object, texture, p.canvas)
	} else {
		fmt.Println("cache hit")
	}
	// if !cache.IsValid(texture) {
	// 	return noTexture, fmt.Errorf("no texture available")
	// }
	return Texture(texture), nil
}
