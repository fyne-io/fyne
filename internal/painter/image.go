package painter

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	_ "image/jpeg" // avoid users having to import when using image widget
	_ "image/png"  // avoid the same for PNG images
	"io"
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/internal"
	"fyne.io/fyne/v2/internal/cache"

	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
	"golang.org/x/image/draw"
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

// PaintImageWithFillModeRespected renders a given fyne Image to a Go standard image.
// If the image’s fill mode is “fill original” but its min size does not match the original image size,
// it does not paint the image but adjusts the image’s min size.
// The image will then be painted on the next frame because of the min size change.
func PaintImageWithFillModeRespected(img *canvas.Image, c fyne.Canvas, width, height int) image.Image {
	var wantOrigW, wantOrigH int
	wantOrigSize := false
	if img.FillMode == canvas.ImageFillOriginal {
		wantOrigW = internal.ScaleInt(c, img.Size().Width)
		wantOrigH = internal.ScaleInt(c, img.Size().Height)
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

// PaintImageWithFillModeIgnored renders a given fyne Image to a Go standard image.
// The image’s fill mode is ignored.
func PaintImageWithFillModeIgnored(img *canvas.Image, width, height int) image.Image {
	dst, _, _, err := paintImage(img, width, height, false, 0, 0)
	if err != nil {
		fyne.LogError("failed to paint image", err)
		return nil
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
				tex, err = svgToImage(file, width, height, checkSize)
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
			if checkSize(origSize.X, origSize.Y) {
				dst = scaleImage(pixels, width, height, img.ScaleMode)
			}
		}
	case img.Image != nil:
		origSize := img.Image.Bounds().Size()
		if checkSize(origSize.X, origSize.Y) {
			dst = scaleImage(img.Image, width, height, img.ScaleMode)
		}
	default:
		dst = image.NewNRGBA(image.Rect(0, 0, 1, 1))
	}
	return
}

func svgToImage(file io.Reader, width, height int, validateSize func(origW, origH int) bool) (*image.NRGBA, error) {
	icon, err := oksvg.ReadIconStream(file)
	if err != nil {
		return nil, fmt.Errorf("SVG Load error: %w", err)
	}

	origW, origH := int(icon.ViewBox.W), int(icon.ViewBox.H)
	if !validateSize(origW, origH) {
		return nil, nil
	}

	aspect := float32(origW) / float32(origH)
	viewAspect := float32(width) / float32(height)
	imgW, imgH := width, height
	if viewAspect > aspect {
		imgW = int(float32(height) * aspect)
	} else if viewAspect < aspect {
		imgH = int(float32(width) / aspect)
	}
	icon.SetTarget(0, 0, float64(imgW), float64(imgH))

	img := image.NewNRGBA(image.Rect(0, 0, imgW, imgH))
	scanner := rasterx.NewScannerGV(origW, origH, img, img.Bounds())
	raster := rasterx.NewDasher(width, height, scanner)

	err = drawSVGSafely(icon, raster)
	if err != nil {
		err = fmt.Errorf("SVG render error: %w", err)
		return nil, err
	}
	return img, nil
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
