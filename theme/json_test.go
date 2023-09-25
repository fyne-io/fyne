package theme

import (
	"image/color"
	"testing"

	"fyne.io/fyne/v2"
	intRepo "fyne.io/fyne/v2/internal/repository"
	"fyne.io/fyne/v2/storage/repository"

	"github.com/stretchr/testify/assert"
)

func TestFromJSON(t *testing.T) {
	repository.Register("file", intRepo.NewFileRepository()) // file uri resolving (avoid test import loop)
	th, err := FromJSON(`{
"Colors": {"background": "#c0c0c0ff"},
"Colors-light": {"foreground": "#ffffffff"},
"Sizes": {"iconInline": 5.0},
"Fonts": {"monospace": "file://./testdata/NotoMono-Regular.ttf"},
"Icons": {"cancel": "file://./testdata/cancel_Paths.svg"}
}`)

	assert.Nil(t, err)
	assert.Equal(t, &color.NRGBA{R: 0xc0, G: 0xc0, B: 0xc0, A: 0xff}, th.Color(ColorNameBackground, VariantDark))
	assert.Equal(t, &color.NRGBA{R: 0xc0, G: 0xc0, B: 0xc0, A: 0xff}, th.Color(ColorNameBackground, VariantLight))
	assert.Equal(t, &color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}, th.Color(ColorNameForeground, VariantLight))
	assert.Equal(t, float32(5), th.Size(SizeNameInlineIcon))
	assert.Equal(t, "NotoMono-Regular.ttf", th.Font(fyne.TextStyle{Monospace: true}).Name())
	assert.Equal(t, "cancel_Paths.svg", th.Icon(IconNameCancel).Name())

	th, _ = FromJSON("{\"Colors\":{\"foreground\":\"\xb1\"}}")
	c := th.Color(ColorNameForeground, VariantLight)
	assert.NotNil(t, c)
}

func TestFromTOML_Resource(t *testing.T) {
	r, err := fyne.LoadResourceFromPath("./testdata/theme.json")
	assert.Nil(t, err)
	th, err := FromJSON(string(r.Content()))

	assert.Nil(t, err)
	assert.Equal(t, &color.NRGBA{R: 0x30, G: 0x30, B: 0x30, A: 0xff}, th.Color(ColorNameBackground, VariantLight))
	assert.Equal(t, &color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}, th.Color(ColorNameForeground, VariantDark))
	assert.Equal(t, &color.NRGBA{R: 0xc0, G: 0xc0, B: 0xc0, A: 0xff}, th.Color(ColorNameForeground, VariantLight))
	assert.Equal(t, float32(10), th.Size(SizeNameInlineIcon))
}

func TestHexColor(t *testing.T) {
	c, err := hexColor("#abc").color()
	assert.Nil(t, err)
	assert.Equal(t, &color.NRGBA{R: 0xaa, G: 0xbb, B: 0xcc, A: 0xff}, c)
	c, err = hexColor("abc").color()
	assert.Nil(t, err)
	assert.Equal(t, &color.NRGBA{R: 0xaa, G: 0xbb, B: 0xcc, A: 0xff}, c)
	c, err = hexColor("#abcd").color()
	assert.Nil(t, err)
	assert.Equal(t, &color.NRGBA{R: 0xaa, G: 0xbb, B: 0xcc, A: 0xdd}, c)

	c, err = hexColor("#a1b2c3").color()
	assert.Nil(t, err)
	assert.Equal(t, &color.NRGBA{R: 0xa1, G: 0xb2, B: 0xc3, A: 0xff}, c)
	c, err = hexColor("a1b2c3").color()
	assert.Nil(t, err)
	assert.Equal(t, &color.NRGBA{R: 0xa1, G: 0xb2, B: 0xc3, A: 0xff}, c)
	c, err = hexColor("#a1b2c3f4").color()
	assert.Nil(t, err)
	assert.Equal(t, &color.NRGBA{R: 0xa1, G: 0xb2, B: 0xc3, A: 0xf4}, c)
	c, err = hexColor("a1b2c3f4").color()
	assert.Nil(t, err)
	assert.Equal(t, &color.NRGBA{R: 0xa1, G: 0xb2, B: 0xc3, A: 0xf4}, c)
}
