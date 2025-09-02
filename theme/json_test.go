package theme

import (
	"image/color"
	"io"
	"log"
	"os"
	"testing"

	"fyne.io/fyne/v2"
	intRepo "fyne.io/fyne/v2/internal/repository"
	"fyne.io/fyne/v2/internal/theme"
	"fyne.io/fyne/v2/storage/repository"

	"github.com/stretchr/testify/assert"
)

func TestFromJSON(t *testing.T) {
	// Discarding log output for tests
	// The following method logs an error:
	// th.Color(ColorNameForeground, VariantLight)
	log.SetOutput(io.Discard)
	t.Cleanup(func() { log.SetOutput(os.Stderr) })
	repository.Register("file", intRepo.NewFileRepository()) // file uri resolving (avoid test import loop)
	th, err := FromJSON(`{
"Colors": {"background": "#c0c0c0ff"},
"Colors-light": {"foreground": "#ffffffff"},
"Sizes": {"iconInline": 5.0},
"Fonts": {"monospace": "file://./testdata/NotoMono-Regular.ttf"},
"Icons": {"cancel": "file://./testdata/cancel_Paths.svg"}
}`)

	assert.NoError(t, err)
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
	assert.NoError(t, err)
	th, err := FromJSON(string(r.Content()))

	assert.NoError(t, err)
	assert.Equal(t, &color.NRGBA{R: 0x30, G: 0x30, B: 0x30, A: 0xff}, th.Color(ColorNameBackground, VariantLight))
	assert.Equal(t, &color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}, th.Color(ColorNameForeground, VariantDark))
	assert.Equal(t, &color.NRGBA{R: 0xc0, G: 0xc0, B: 0xc0, A: 0xff}, th.Color(ColorNameForeground, VariantLight))
	assert.Equal(t, float32(10), th.Size(SizeNameInlineIcon))
}

func TestHexColor(t *testing.T) {
	c := jsonColor{}
	err := c.parseColor("#abc")
	assert.NoError(t, err)
	assert.Equal(t, &color.NRGBA{R: 0xaa, G: 0xbb, B: 0xcc, A: 0xff}, c.color)
	err = c.parseColor("abc")
	assert.NoError(t, err)
	assert.Equal(t, &color.NRGBA{R: 0xaa, G: 0xbb, B: 0xcc, A: 0xff}, c.color)
	err = c.parseColor("#abcd")
	assert.NoError(t, err)
	assert.Equal(t, &color.NRGBA{R: 0xaa, G: 0xbb, B: 0xcc, A: 0xdd}, c.color)

	err = c.parseColor("#a1b2c3")
	assert.NoError(t, err)
	assert.Equal(t, &color.NRGBA{R: 0xa1, G: 0xb2, B: 0xc3, A: 0xff}, c.color)
	err = c.parseColor("a1b2c3")
	assert.NoError(t, err)
	assert.Equal(t, &color.NRGBA{R: 0xa1, G: 0xb2, B: 0xc3, A: 0xff}, c.color)
	err = c.parseColor("#a1b2c3f4")
	assert.NoError(t, err)
	assert.Equal(t, &color.NRGBA{R: 0xa1, G: 0xb2, B: 0xc3, A: 0xf4}, c.color)
	err = c.parseColor("a1b2c3f4")
	assert.NoError(t, err)
	assert.Equal(t, &color.NRGBA{R: 0xa1, G: 0xb2, B: 0xc3, A: 0xf4}, c.color)
}

var result color.Color

func BenchmarkJsonTheme_Color(b *testing.B) {
	th, err := FromJSON(`
{
    "Colors-light": {
        "pageBackground": "#EFF1F5FF",
        "listHeader": "#DCE0E8FF",
        "pageHeader": "#E6E9EFFF",
        "background": "#DCE0E8FF",
        "button": "#E6E9EFFF",
        "disabledButton": "#CCD0DAFF",
        "disabled": "#4C4F69FF",
        "error": "#D20F39FF",
        "focus": "#EDCFD5FF",
        "foregroundOnError": "#EFF1F5FF",
        "foregroundOnPrimary": "#EFF1F5FF",
        "foregroundOnSuccess": "#EFF1F5FF",
        "foregroundOnWarning": "#EFF1F5FF",
        "headerBackground": "#CCD0DAFF",
        "hover": "#EFE0E5FF",
        "hyperlink": "#E64553FF",
        "inputBackground": "#E6E9EFFF",
        "inputBorder": "#CCD0DAFF",
        "menuBackground": "#E6E9EFFF",
        "overlayBackground": "#E6E9EFFF",
        "placeholder": "#7C7F93FF",
        "pressed": "#CCD0DA66",
        "primary": "#E64553FF",
        "scrollBar": "#E6E9EF99",
        "selection": "#EDC6CDFF",
        "separator": "#DCE0E8FF",
        "success": "#40A02BFF",
        "warning": "#DF8E1DFF"
    },
	    "Colors-dark": {
        "pageBackground": "#1E1E2EFF",
        "listHeader": "#11111BFF",
        "pageHeader": "#181825FF",
        "background": "#11111BFF",
        "button": "#181825FF",
        "disabledButton": "#313244FF",
        "disabled": "#CDD6F4FF",
        "error": "#F38BA8FF",
        "focus": "#473847FF",
        "foreground": "#CDD6F4FF",
        "foregroundOnError": "#1E1E2EFF",
        "foregroundOnPrimary": "#1E1E2EFF",
        "foregroundOnSuccess": "#1E1E2EFF",
        "foregroundOnWarning": "#1E1E2EFF",
        "headerBackground": "#313244FF",
        "hover": "#322B3AFF",
        "hyperlink": "#EBA0ACFF",
        "inputBackground": "#181825FF",
        "inputBorder": "#313244FF",
        "menuBackground": "#181825FF",
        "overlayBackground": "#181825FF",
        "placeholder": "#9399B2FF",
        "pressed": "#31324466",
        "primary": "#EBA0ACFF",
        "scrollBar": "#18182599",
        "selection": "#513E4DFF",
        "separator": "#11111BFF",
        "success": "#A6E3A1FF",
        "warning": "#F9E2AFFF"
    }
}`)

	assert.NoError(b, err)
	var localResult color.Color
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		localResult = th.Color("primary", theme.VariantDark)
	}
	result = localResult
}
