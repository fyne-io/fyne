package canvas_test

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
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
	"fyne.io/fyne/v2/internal/cache"
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

func TestImage_SharedResourceIsDecodedOnce(t *testing.T) {
	cache.ResetThemeCaches()
	pwd, _ := os.Getwd()
	path := filepath.Join(filepath.Dir(pwd), "theme", "icons", "fyne.png")
	data, err := os.ReadFile(path)
	assert.NoError(t, err)

	res := fyne.NewStaticResource("shared.png", data)
	first := canvas.NewImageFromResource(res)
	second := canvas.NewImageFromResource(res)
	first.Resize(fyne.NewSize(64, 64))
	second.Resize(fyne.NewSize(64, 64))

	assert.NotNil(t, first.Image)
	assert.Same(t, first.Image, second.Image)
}

func TestImage_RefreshKeepsDecodedImageWhenSourceUnchanged(t *testing.T) {
	cache.ResetThemeCaches()
	res := fyne.NewStaticResource("shared.png", makeTestPNG(t, color.White))
	img := canvas.NewImageFromResource(res)
	img.Resize(fyne.NewSize(16, 16))
	decoded := img.Image

	cache.ResetThemeCaches() // simulate the shared cache having expired

	img.Refresh()
	assert.Same(t, decoded, img.Image)
}

func TestImage_RefreshRedecodesWhenContentChanges(t *testing.T) {
	cache.ResetThemeCaches()
	res := fyne.NewStaticResource("shared.png", makeTestPNG(t, color.White))
	img := canvas.NewImageFromResource(res)
	img.Resize(fyne.NewSize(16, 16))
	decoded := img.Image

	res.StaticContent = makeTestPNG(t, color.Black)
	img.Refresh()
	assert.NotSame(t, decoded, img.Image)
}

func makeTestPNG(t *testing.T, clr color.Color) []byte {
	t.Helper()

	img := image.NewNRGBA(image.Rect(0, 0, 2, 2))
	for y := 0; y < 2; y++ {
		for x := 0; x < 2; x++ {
			img.Set(x, y, clr)
		}
	}

	var buf bytes.Buffer
	assert.NoError(t, png.Encode(&buf, img))
	return buf.Bytes()
}
