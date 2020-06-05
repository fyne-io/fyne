package test

import (
	"image"
	"image/color"
	"math"
)

// NewCheckedImage returns a new black/white checked image with the specified size
// and the specified amount of horizontal and vertical tiles.
func NewCheckedImage(w, h, hTiles, vTiles int) image.Image {
	img := image.NewNRGBA(image.Rect(0, 0, w, h))
	colors := []color.Color{color.White, color.Black}
	tileWidth := float64(w) / float64(hTiles)
	tileHeight := float64(h) / float64(vTiles)
	for y := 0; y < h; y++ {
		yTile := int(math.Floor(float64(y) / tileHeight))
		for x := 0; x < w; x++ {
			xTile := int(math.Floor(float64(x) / tileWidth))
			img.Set(x, y, colors[(xTile+yTile)%2])
		}
	}
	return img
}
