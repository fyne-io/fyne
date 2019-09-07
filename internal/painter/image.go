package painter

import (
	"bytes"
	"errors"
	"image"
	"image/draw"
	_ "image/jpeg" // avoid users having to import when using image widget
	_ "image/png"  // avoid the same for PNG images
	"io"
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/internal"
	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
)

// PaintImage renders a given fyne Image to a Go standard image
func PaintImage(img *canvas.Image, c fyne.Canvas, width, height int) image.Image {
	if width <= 0 || height <= 0 {
		return nil
	}

	switch {
	case img.File != "" || img.Resource != nil:
		var file io.Reader
		var name string
		if img.Resource != nil {
			name = img.Resource.Name()
			file = bytes.NewReader(img.Resource.Content())
		} else {
			name = img.File
			handle, _ := os.Open(img.File)
			defer handle.Close()
			file = handle
		}

		if strings.ToLower(filepath.Ext(name)) == ".svg" {
			tex := svgCacheGet(img.Resource, width, height)
			if tex == nil {
				// Not in cache, so load the item and add to cache

				icon, err := oksvg.ReadIconStream(file)
				if err != nil {
					fyne.LogError("SVG Load error:", err)
					return nil
				}
				icon.SetTarget(0, 0, float64(width), float64(height))

				w, h := int(icon.ViewBox.W), int(icon.ViewBox.H)
				// this is used by our render code, so let's set it to the file aspect
				aspects[img.Resource] = float32(w) / float32(h)
				// if the image specifies it should be original size we need at least that many pixels on screen
				if img.FillMode == canvas.ImageFillOriginal {
					checkImageMinSize(img, c, w, h)
				}

				tex = image.NewRGBA(image.Rect(0, 0, width, height))
				scanner := rasterx.NewScannerGV(w, h, tex, tex.Bounds())
				raster := rasterx.NewDasher(width, height, scanner)

				err = drawSVGSafely(icon, raster)
				if err != nil {
					fyne.LogError("SVG Render error:", err)
					return nil
				}
				svgCachePut(img.Resource, tex, width, height)
			}

			return tex
		}

		pixels, _, err := image.Decode(file)

		if err != nil {
			fyne.LogError("image err", err)

			return nil
		}
		origSize := pixels.Bounds().Size()
		// this is used by our render code, so let's set it to the file aspect
		aspects[img] = float32(origSize.X) / float32(origSize.Y)
		// if the image specifies it should be original size we need at least that many pixels on screen
		if img.FillMode == canvas.ImageFillOriginal {
			checkImageMinSize(img, c, origSize.X, origSize.Y)
		}

		tex := image.NewRGBA(image.Rect(0, 0, pixels.Bounds().Dx(), pixels.Bounds().Dy()))
		draw.Draw(tex, tex.Bounds(), pixels, pixels.Bounds().Min, draw.Src)

		return tex
	case img.Image != nil:
		origSize := img.Image.Bounds().Size()
		// this is used by our render code, so let's set it to the file aspect
		aspects[img] = float32(origSize.X) / float32(origSize.Y)
		// if the image specifies it should be original size we need at least that many pixels on screen
		if img.FillMode == canvas.ImageFillOriginal {
			checkImageMinSize(img, c, origSize.X, origSize.Y)
		}

		tex := image.NewRGBA(image.Rect(0, 0, origSize.X, origSize.Y))
		draw.Draw(tex, tex.Bounds(), img.Image, img.Image.Bounds().Min, draw.Src)

		return tex
	default:
		return image.NewRGBA(image.Rect(0, 0, 1, 1))
	}
}

func drawSVGSafely(icon *oksvg.SvgIcon, raster *rasterx.Dasher) error {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("Crash when rendering SVG")
		}
	}()
	icon.Draw(raster, 1)

	return err
}

func checkImageMinSize(img *canvas.Image, c fyne.Canvas, pixX, pixY int) {
	pixSize := fyne.NewSize(internal.UnscaleInt(c, pixX), internal.UnscaleInt(c, pixY))

	if img.MinSize() != pixSize {
		img.SetMinSize(pixSize)
		canvas.Refresh(img) // force the initial size to be respected
	}
}
