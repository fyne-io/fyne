package painter

import (
	"bytes"
	"fmt"
	"image"
	_ "image/jpeg" // avoid users having to import when using image widget
	_ "image/png"  // avoid the same for PNG images
	"io"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/image/draw"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/internal"
	"fyne.io/fyne/v2/internal/cache"
	"fyne.io/fyne/v2/internal/svg"
)

var aspects = make(map[interface{}]float32, 16)

// GetAspect looks up an aspect ratio of an image
func GetAspect(img *canvas.Image) float32 {
	aspect := float32(0.0)
	if img.Resource != nil {
		aspect = aspects[img.Resource.Name()]
	} else if img.File != "" {
		aspect = aspects[img.File]
	}

	if aspect == 0 {
		aspect = aspects[img]
	}

	return aspect
}

// PaintImage renders a given fyne Image to a Go standard image
// If a fyne.Canvas is given and the image’s fill mode is “fill original” the image’s min size has
// to fit its original size. If it doesn’t, PaintImage does not paint the image but adjusts its min size.
// The image will then be painted on the next frame because of the min size change.
func PaintImage(img *canvas.Image, c fyne.Canvas, width, height int) image.Image {
	var wantOrigW, wantOrigH int
	wantOrigSize := false
	if img.FillMode == canvas.ImageFillOriginal && c != nil {
		wantOrigW = internal.ScaleInt(c, img.MinSize().Width)
		wantOrigH = internal.ScaleInt(c, img.MinSize().Height)
		wantOrigSize = true
	}

	dst, origW, origH, err := paintImage(img, width, height, wantOrigSize, wantOrigW, wantOrigH)
	if err != nil {
		fyne.LogError("failed to paint image", err)
		return nil
	}

	if wantOrigSize && dst == nil {
		dpSize := fyne.NewSize(internal.UnscaleInt(c, origW), internal.UnscaleInt(c, origH))
		img.SetMinSize(dpSize)
		canvas.Refresh(img) // force the initial size to be respected
	}
	return dst
}

func paintImage(img *canvas.Image, width, height int, wantOrigSize bool, wantOrigW, wantOrigH int) (dst image.Image, origW, origH int, err error) {
	if width <= 0 || height <= 0 {
		return
	}

	var aspectCacheKey interface{} = img
	checkSize := func(origW, origH int) bool {
		aspect := float32(origW) / float32(origH)
		// this is used by our render code, so let's set it to the file aspect
		aspects[aspectCacheKey] = aspect
		return !wantOrigSize || (wantOrigW == origW && wantOrigH == origH)
	}

	switch {
	case img.File != "" || img.Resource != nil:
		var (
			file  io.Reader
			name  string
			isSVG bool
		)
		if img.Resource != nil {
			name = img.Resource.Name()
			file = bytes.NewReader(img.Resource.Content())
			isSVG = IsResourceSVG(img.Resource)
		} else {
			name = img.File
			var handle *os.File
			handle, err = os.Open(img.File)
			if err != nil {
				err = fmt.Errorf("image load error: %w", err)
				return
			}
			defer handle.Close()
			file = handle
			isSVG = isFileSVG(img.File)
		}
		aspectCacheKey = name

		if isSVG {
			tex := cache.GetSvg(name, width, height)
			if tex == nil {
				// Not in cache, so load the item and add to cache
				tex, err = svg.ToImage(file, width, height, checkSize)
				if err != nil {
					return
				}

				cache.SetSvg(name, tex, width, height)
			}
			dst = tex
		} else {
			var pixels image.Image
			pixels, _, err = image.Decode(file)
			if err != nil {
				err = fmt.Errorf("failed to decode image: %w", err)
				return
			}

			origSize := pixels.Bounds().Size()
			origW, origH = origSize.X, origSize.Y
			if checkSize(origSize.X, origSize.Y) {
				dst = scaleImage(pixels, width, height, img.ScaleMode)
			}
		}
	case img.Image != nil:
		origSize := img.Image.Bounds().Size()
		origW, origH = origSize.X, origSize.Y
		if checkSize(origSize.X, origSize.Y) {
			dst = scaleImage(img.Image, width, height, img.ScaleMode)
		}
	default:
		dst = image.NewNRGBA(image.Rect(0, 0, 1, 1))
	}
	return
}

func scaleImage(pixels image.Image, scaledW, scaledH int, scale canvas.ImageScale) image.Image {
	if scale == canvas.ImageScaleFastest || scale == canvas.ImageScalePixels {
		// do not perform software scaling
		return pixels
	}

	pixW := int(fyne.Min(float32(scaledW), float32(pixels.Bounds().Dx()))) // don't push more pixels than we have to
	pixH := int(fyne.Min(float32(scaledH), float32(pixels.Bounds().Dy()))) // the GL calls will scale this up on GPU.
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

func isFileSVG(path string) bool {
	return strings.ToLower(filepath.Ext(path)) == ".svg"
}

// IsResourceSVG checks if the resource is an SVG or not.
func IsResourceSVG(res fyne.Resource) bool {
	if strings.ToLower(filepath.Ext(res.Name())) == ".svg" {
		return true
	}

	if len(res.Content()) < 5 {
		return false
	}

	switch strings.ToLower(string(res.Content()[:5])) {
	case "<!doc", "<?xml", "<svg ":
		return true
	}
	return false
}
