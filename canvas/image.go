package canvas

import (
	"bytes"
	"errors"
	"image"
	_ "image/jpeg" // avoid users having to import when using image widget
	_ "image/png"  // avoid the same for PNG images
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/cache"
	"fyne.io/fyne/v2/internal/scale"
	"fyne.io/fyne/v2/internal/svg"
	"fyne.io/fyne/v2/storage"
)

// ImageFill defines the different type of ways an image can stretch to fill its space.
type ImageFill int

const (
	// ImageFillStretch will scale the image to match the Size() values.
	// This is the default and does not maintain aspect ratio.
	ImageFillStretch ImageFill = iota
	// ImageFillContain makes the image fit within the object Size(),
	// centrally and maintaining aspect ratio.
	// There may be transparent sections top and bottom or left and right.
	ImageFillContain // (Fit)
	// ImageFillOriginal ensures that the container grows to the pixel dimensions
	// required to fit the original image. The aspect of the image will be maintained so,
	// as with ImageFillContain there may be transparent areas around the image.
	// Note that the minSize may be smaller than the image dimensions if scale > 1.
	ImageFillOriginal
)

// ImageScale defines the different scaling filters used to scaling images
type ImageScale int32

const (
	// ImageScaleSmooth will scale the image using ApproxBiLinear filter (or GL equivalent)
	ImageScaleSmooth ImageScale = iota
	// ImageScalePixels will scale the image using NearestNeighbor filter (or GL equivalent)
	ImageScalePixels
	// ImageScaleFastest will scale the image using hardware GPU if available
	//
	// Since: 2.0
	ImageScaleFastest
)

// Declare conformity with CanvasObject interface
var _ fyne.CanvasObject = (*Image)(nil)

// Image describes a drawable image area that can render in a Fyne canvas
// The image may be a vector or a bitmap representation and it will fill the area.
// The fill mode can be changed by setting FillMode to a different ImageFill.
type Image struct {
	baseObject

	aspect float32
	reader io.Reader
	icon   *svg.Decoder
	isSVG  bool
	lock   sync.Mutex

	// one of the following sources will provide our image data
	File     string        // Load the image from a file
	Resource fyne.Resource // Load the image from an in-memory resource
	Image    image.Image   // Specify a loaded image to use in this canvas object

	Translucency float64    // Set a translucency value > 0.0 to fade the image
	FillMode     ImageFill  // Specify how the image should expand to fill or fit the available space
	ScaleMode    ImageScale // Specify the type of scaling interpolation applied to the image
}

// Alpha is a convenience function that returns the alpha value for an image
// based on its Translucency value. The result is 1.0 - Translucency.
func (i *Image) Alpha() float64 {
	return 1.0 - i.Translucency
}

// Resize on an image will scale the content or reposition it according to FillMode.
// It will normally cause a Refresh to ensure the pixels are recalculated.
func (i *Image) Resize(s fyne.Size) {
	if s == i.Size() {
		return
	}
	if i.FillMode == ImageFillOriginal && i.size.Height > 2 { // don't refresh original scale images after first draw
		return
	}

	i.baseObject.Resize(s)
	i.Refresh()
}

// Aspect will return the original content aspect after it was last Refresh
func (i *Image) Aspect() float32 {
	if i.aspect == 0 {
		i.Refresh()
	}
	return i.aspect
}

func (i *Image) MinSize() fyne.Size {
	if i.Image == nil || i.aspect == 0 {
		i.Refresh()
	}
	return i.baseObject.MinSize()
}

// Refresh causes this object to be redrawn in it's current state
func (i *Image) Refresh() {
	i.lock.Lock()
	defer i.lock.Unlock()

	err := i.updateReader()
	if err != nil {
		fyne.LogError("Failed to load image", err)
		return
	}
	err = i.calculateMinSize()
	if err != nil {
		fyne.LogError("Failed to load image", err)
		return
	}

	if i.File != "" || i.Resource != nil {
		size := i.Size()
		width := size.Width
		height := size.Height

		if width == 0 || height == 0 {
			return
		}

		if i.isSVG {
			tex, err := i.renderSVG(width, height)
			if err != nil {
				fyne.LogError("Failed to render SVG", err)
				return
			}
			i.Image = tex
		} else {
			if i.reader == nil {
				return
			}

			image, _, err := image.Decode(i.reader)
			if err != nil {
				fyne.LogError("Failed to render image", err)
				return
			}
			i.Image = image
		}
	}

	Refresh(i)
}

// NewImageFromFile creates a new image from a local file.
// Images returned from this method will scale to fit the canvas object.
// The method for scaling can be set using the Fill field.
func NewImageFromFile(file string) *Image {
	return &Image{File: file}
}

// NewImageFromURI creates a new image from named resource.
// File URIs will read the file path and other schemes will download the data into a resource.
// HTTP and HTTPs URIs will use the GET method by default to request the resource.
// Images returned from this method will scale to fit the canvas object.
// The method for scaling can be set using the Fill field.
//
// Since: 2.0
func NewImageFromURI(uri fyne.URI) *Image {
	if uri.Scheme() == "file" && len(uri.String()) > 7 {
		return NewImageFromFile(uri.Path())
	}

	var read io.ReadCloser

	read, err := storage.Reader(uri) // attempt unknown / http file type
	if err != nil {
		fyne.LogError("Failed to open image URI", err)
		return &Image{}
	}

	defer read.Close()
	return NewImageFromReader(read, filepath.Base(uri.String()))
}

// NewImageFromReader creates a new image from a data stream.
// The name parameter is required to uniquely identify this image (for caching etc).
// If the image in this io.Reader is an SVG, the name should end ".svg".
// Images returned from this method will scale to fit the canvas object.
// The method for scaling can be set using the Fill field.
//
// Since: 2.0
func NewImageFromReader(read io.Reader, name string) *Image {
	data, err := ioutil.ReadAll(read)
	if err != nil {
		fyne.LogError("Unable to read image data", err)
		return nil
	}

	res := &fyne.StaticResource{
		StaticName:    name,
		StaticContent: data,
	}

	return NewImageFromResource(res)
}

// NewImageFromResource creates a new image by loading the specified resource.
// Images returned from this method will scale to fit the canvas object.
// The method for scaling can be set using the Fill field.
func NewImageFromResource(res fyne.Resource) *Image {
	return &Image{Resource: res}
}

// NewImageFromImage returns a new Image instance that is rendered from the Go
// image.Image passed in.
// Images returned from this method will scale to fit the canvas object.
// The method for scaling can be set using the Fill field.
func NewImageFromImage(img image.Image) *Image {
	return &Image{Image: img}
}

func (i *Image) name() string {
	if i.Resource != nil {
		return i.Resource.Name()
	} else if i.File != "" {
		return i.File
	}
	return ""
}

func (i *Image) updateReader() error {
	i.reader = nil

	var fd io.Reader
	if i.Resource != nil {
		fd = bytes.NewReader(i.Resource.Content())
		i.isSVG = svg.IsResourceSVG(i.Resource)
	} else if i.File != "" {
		var err error

		fd, err = os.Open(i.File)
		if err != nil {
			return err
		}
		i.isSVG = svg.IsFileSVG(i.File)
	}
	i.reader = fd
	return nil
}

func (i *Image) calculateMinSize() error {
	var size fyne.Size

	if i.Resource != nil || i.File != "" {
		var err error

		size, err = i.minSizeFromReader(i.reader)
		if err != nil {
			return err
		}
	} else if i.Image != nil {
		original := i.Image.Bounds().Size()
		size = fyne.NewSize(float32(original.X), float32(original.Y))
		i.aspect = size.Width / size.Height
	} else {
		return errors.New("no matching image source")
	}

	if i.FillMode == ImageFillOriginal {
		i.SetMinSize(size)
	}
	return nil
}

func (img *Image) minSizeFromReader(source io.Reader) (fyne.Size, error) {
	var width int
	var height int

	if source == nil {
		return fyne.NewSize(0, 0), errors.New("no matching reading reader")
	}

	if img.isSVG {
		var err error

		img.icon, err = svg.NewDecoder(source)
		if err != nil {
			return fyne.NewSize(0, 0), err
		}
		config := img.icon.Config()
		width, height = config.Width, config.Height
		img.aspect = config.Aspect
	} else {
		var buf bytes.Buffer
		tee := io.TeeReader(source, &buf)
		img.reader = io.MultiReader(&buf, source)

		config, _, err := image.DecodeConfig(tee)
		if err != nil {
			return fyne.NewSize(0, 0), err
		}
		width, height = config.Width, config.Height
		img.aspect = float32(width) / float32(height)
	}

	app := fyne.CurrentApp()
	if app == nil {
		return fyne.NewSize(0, 0), errors.New("no current app")
	}
	driver := app.Driver()
	if driver == nil {
		return fyne.NewSize(0, 0), errors.New("no current driver")
	}
	c := driver.CanvasForObject(img)
	if c == nil {
		return fyne.NewSize(0, 0), errors.New("object is not attached to a canvas yet")
	}
	dpSize := fyne.NewSize(scale.UnscaleInt(c, width), scale.UnscaleInt(c, height))

	return dpSize, nil
}

func (img *Image) renderSVG(width, height float32) (image.Image, error) {
	c := fyne.CurrentApp().Driver().CanvasForObject(img)
	if c == nil {
		return nil, errors.New("object is not attached to a canvas yet")
	}

	screenWidth := scale.ScaleInt(c, width)
	screenHeight := scale.ScaleInt(c, height)

	tex := cache.GetSvg(img.name(), screenWidth, screenHeight)
	if tex != nil {
		return tex, nil
	}

	var err error
	tex, err = img.icon.Draw(screenWidth, screenHeight)
	if err != nil {
		return nil, err
	}
	cache.SetSvg(img.name(), tex, screenWidth, screenHeight)
	return tex, nil
}
