package painter

import (
	"bytes"
	"errors"
	"image"
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
	"golang.org/x/image/draw"
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
			tex := svgCacheGet(name, width, height)
			if tex == nil {
				// Not in cache, so load the item and add to cache

				icon, err := oksvg.ReadIconStream(file)
				if err != nil {
					fyne.LogError("SVG Load error:", err)
					return nil
				}

				origW, origH := int(icon.ViewBox.W), int(icon.ViewBox.H)
				aspect := float32(origW) / float32(origH)
				viewAspect := float32(width) / float32(height)

				texW, texH := width, height
				if viewAspect > aspect {
					texW = int(float32(height) * aspect)
				} else if viewAspect < aspect {
					texH = int(float32(width) / aspect)
				}

				icon.SetTarget(0, 0, float64(texW), float64(texH))
				// this is used by our render code, so let's set it to the file aspect
				aspects[img.Resource] = aspect
				// if the image specifies it should be original size we need at least that many pixels on screen
				if img.FillMode == canvas.ImageFillOriginal {
					if !checkImageMinSize(img, c, origW, origH) {
						return nil
					}
				}

				tex = image.NewNRGBA(image.Rect(0, 0, texW, texH))
				scanner := rasterx.NewScannerGV(origW, origH, tex, tex.Bounds())
				raster := rasterx.NewDasher(width, height, scanner)

				err = drawSVGSafely(icon, raster)
				if err != nil {
					fyne.LogError("SVG Render error:", err)
					return nil
				}

				svgCachePut(name, tex, width, height)
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
			if !checkImageMinSize(img, c, origSize.X, origSize.Y) {
				return nil
			}
		}

		return scaleImage(pixels, width, height, img.ScaleMode)
	case img.Image != nil:
		origSize := img.Image.Bounds().Size()
		// this is used by our render code, so let's set it to the file aspect
		aspects[img] = float32(origSize.X) / float32(origSize.Y)
		// if the image specifies it should be original size we need at least that many pixels on screen
		if img.FillMode == canvas.ImageFillOriginal {
			if !checkImageMinSize(img, c, origSize.X, origSize.Y) {
				return nil
			}
		}

		return scaleImage(img.Image, width, height, img.ScaleMode)
	default:
		return image.NewNRGBA(image.Rect(0, 0, 1, 1))
	}
}

func scaleImage(pixels image.Image, scaledW, scaledH int, scale canvas.ImageScale) image.Image {
	pixW := fyne.Min(scaledW, pixels.Bounds().Dx()) // don't push more pixels than we have to
	pixH := fyne.Min(scaledH, pixels.Bounds().Dy()) // the GL calls will scale this up on GPU.
	scaledBounds := image.Rect(0, 0, pixW, pixH)
	tex := image.NewNRGBA(scaledBounds)
	switch scale {
	case canvas.ImageScalePixels:
		draw.NearestNeighbor.Scale(tex, scaledBounds, pixels, pixels.Bounds(), draw.Over, nil)
	default:
		if scale != canvas.ImageScaleSmooth {
			fyne.LogError("Invalid canvas.ImageScale value, using canvas.ImageScaleSmooth", nil)
		}
		draw.CatmullRom.Scale(tex, scaledBounds, pixels, pixels.Bounds(), draw.Over, nil)
	}
	return tex
}

func drawSVGSafely(icon *oksvg.SvgIcon, raster *rasterx.Dasher) error {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("crash when rendering svg")
		}
	}()
	icon.Draw(raster, 1)

	return err
}

func checkImageMinSize(img *canvas.Image, c fyne.Canvas, pixX, pixY int) bool {
	dpSize := fyne.NewSize(internal.UnscaleInt(c, pixX), internal.UnscaleInt(c, pixY))

	if img.MinSize() != dpSize {
		img.SetMinSize(dpSize)
		canvas.Refresh(img) // force the initial size to be respected
		return false
	}

	return true
}
