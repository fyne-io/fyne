package canvas_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/software"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"

	"github.com/stretchr/testify/assert"
)

func TestImage_AlphaDefault(t *testing.T) {
	img := &canvas.Image{}

	assert.Equal(t, 1.0, img.Alpha())
}

func TestImage_TranslucencyDefault(t *testing.T) {
	img := &canvas.Image{}

	assert.Equal(t, 0.0, img.Translucency)
}

func TestImage_RefreshBlank(t *testing.T) {
	img := &canvas.Image{}
	img.Resize(fyne.NewSize(64, 64))
	img.Refresh()
	assert.Nil(t, img.Image)

	img.Resource = theme.HomeIcon()
	img.Refresh()
	assert.NotNil(t, img.Image)

	img.Image = nil
	img.Resource = nil
	img.Refresh()
	assert.Nil(t, img.Image)
}

func TestImage_RefreshSVGChanged(t *testing.T) {
	square := []byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 10 10"><rect width="10" height="10"/></svg>`)
	wide := []byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 10"><rect width="20" height="10"/></svg>`)
	tall := []byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 10 20"><rect width="10" height="20"/></svg>`)

	res := &fyne.StaticResource{StaticName: "shape.svg", StaticContent: square}
	img := canvas.NewImageFromResource(res)
	img.FillMode = canvas.ImageFillOriginal
	img.Refresh()
	assert.Equal(t, float32(1), img.Aspect())
	assert.Equal(t, fyne.NewSize(10, 10), img.MinSize())

	copy(res.StaticContent, wide) // same length, changed in place
	img.Refresh()
	assert.Equal(t, float32(2), img.Aspect())
	assert.Equal(t, fyne.NewSize(20, 10), img.MinSize())

	path := filepath.Join(t.TempDir(), "tall.svg")
	assert.NoError(t, os.WriteFile(path, tall, 0o644))
	img.Resource = nil
	img.File = path
	img.Refresh()
	assert.Equal(t, float32(0.5), img.Aspect())
	assert.Equal(t, fyne.NewSize(10, 20), img.MinSize())
}

func TestNewImageFromFile(t *testing.T) {
	pwd, _ := os.Getwd()
	path := filepath.Join(filepath.Dir(pwd), "theme", "icons", "fyne.png")

	img := canvas.NewImageFromFile(path)
	assert.NotNil(t, img)
	assert.Equal(t, path, img.File)
}

func TestNewImageFromReader(t *testing.T) {
	pwd, _ := os.Getwd()
	path := filepath.Join(filepath.Dir(pwd), "theme", "icons", "fyne.png")
	read, err := os.Open(path)
	assert.NoError(t, err)
	defer read.Close()

	img := canvas.NewImageFromReader(read, "fyne.png")
	assert.NotNil(t, img)
	assert.Equal(t, "", img.File)
	assert.NotNil(t, img.Resource)
	assert.Equal(t, "fyne.png", img.Resource.Name())

	img.FillMode = canvas.ImageFillOriginal

	size := img.MinSize()
	assert.Equal(t, float32(512), size.Width)
	assert.Equal(t, float32(512), size.Height)
}

func TestNewImageFromURI_File(t *testing.T) {
	pwd, _ := os.Getwd()
	path := filepath.Join(filepath.Dir(pwd), "theme", "icons", "fyne.png")

	if runtime.GOOS == "windows" {
		path = strings.ReplaceAll(path, "\\", "/")
	}

	img := canvas.NewImageFromURI(storage.NewFileURI(path))
	assert.NotNil(t, img)
	assert.Equal(t, path, img.File)

	img.FillMode = canvas.ImageFillOriginal

	size := img.MinSize()
	assert.Equal(t, float32(512), size.Width)
	assert.Equal(t, float32(512), size.Height)
}

func TestNewImageFromURI_HTTP(t *testing.T) {
	pwd, _ := os.Getwd()
	path := filepath.Join(filepath.Dir(pwd), "theme", "icons", "fyne.png")
	f, _ := os.ReadFile(path)

	// start a test server to test http calls
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		f, err := os.ReadFile(path)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(f)
	}))
	defer ts.Close()

	// http is mounted in Fyne test handlers by default
	url, _ := storage.ParseURI(ts.URL)
	img := canvas.NewImageFromURI(url)
	assert.NotNil(t, img)
	assert.Equal(t, "", img.File)
	assert.NotNil(t, img.Resource)
	assert.Equal(t, url.Authority(), img.Resource.Name())
	assert.Equal(t, f, img.Resource.Content())

	img.FillMode = canvas.ImageFillOriginal

	size := img.MinSize()
	assert.Equal(t, float32(512), size.Width)
	assert.Equal(t, float32(512), size.Height)
}

func TestImage_CornerRadius(t *testing.T) {
	pwd, _ := os.Getwd()
	path := filepath.Join(filepath.Dir(pwd), "theme", "icons", "fyne.png")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatal("source image not found")
	}

	i := &canvas.Image{
		File:         path,
		CornerRadius: 25,
	}
	c := software.NewCanvas()
	c.SetContent(i)
	c.Resize(fyne.NewSize(120, 120))

	test.AssertRendersToImage(t, "image_rounded_corners.png", c)

	i.CornerRadius = canvas.RadiusMaximum
	c.Resize(fyne.NewSize(60, 60))
	test.AssertRendersToImage(t, "image_fully_rounded_corners.png", c)

	i.FillMode = canvas.ImageFillCover
	i.CornerRadius = 20
	c.Resize(fyne.NewSize(60, 100))
	test.AssertRendersToImage(t, "image_cover_rounded_corners.png", c)
}
