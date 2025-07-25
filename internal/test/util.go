package test

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"path/filepath"
)

// pixCloseEnough reports whether a and b are mostly the same.
func pixCloseEnough(a, b []uint8) bool {
	if len(a) != len(b) {
		return false
	}
	mismatches := 0

	for i, v := range a {
		w := b[i]
		if v == w {
			continue
		}
		// Allow a small delta for rendering variation.
		delta := int(v) - int(w) // use int to avoid overflow
		if delta > 4 || delta < -4 {
			return false
		}
		mismatches++
	}

	// Allow up to 1% of pixels to mismatch.
	return mismatches == 0 || mismatches < len(a)/100
}

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

func writeImage(path string, img image.Image) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	if err = png.Encode(f, img); err != nil {
		f.Close()
		return err
	}
	return f.Close()
}
