package widget

import (
	"image"
	"image/color"
	"net/url"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/internal/cache"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHyperlink_MinSize(t *testing.T) {
	u, err := url.Parse("https://fyne.io/")
	require.NoError(t, err)

	hyperlink := NewHyperlink("Test", u)
	hyperlink.CreateRenderer()
	hyperlink.provider.CreateRenderer()
	minA := hyperlink.MinSize()

	assert.Less(t, theme.InnerPadding(), minA.Width)

	hyperlink.SetText("Longer")
	minB := hyperlink.MinSize()
	assert.Less(t, minA.Width, minB.Width)

	hyperlink.Text = "."
	hyperlink.Refresh()
	minC := hyperlink.MinSize()
	assert.Greater(t, minB.Width, minC.Width)
}

func TestHyperlink_Cursor(t *testing.T) {
	u, err := url.Parse("https://fyne.io/")
	hyperlink := NewHyperlink("Test", u)

	require.NoError(t, err)
	assert.Equal(t, desktop.DefaultCursor, hyperlink.Cursor())

	hyperlink.hovered = true
	assert.Equal(t, desktop.PointerCursor, hyperlink.Cursor())
}

func TestHyperlink_Alignment(t *testing.T) {
	hyperlink := &Hyperlink{Text: "Test", Alignment: fyne.TextAlignTrailing}
	hyperlink.CreateRenderer()

	assert.Equal(t, fyne.TextAlignTrailing, richTextRenderTexts(&hyperlink.provider)[0].Alignment)
}

func TestHyperlink_Hide(t *testing.T) {
	hyperlink := &Hyperlink{Text: "Test"}
	hyperlink.CreateRenderer()
	hyperlink.Hide()
	hyperlink.Refresh()

	assert.True(t, hyperlink.Hidden)
	assert.False(t, hyperlink.provider.Hidden) // we don't propagate hide

	hyperlink.Show()
	assert.False(t, hyperlink.Hidden)
	assert.False(t, hyperlink.provider.Hidden)
}

func TestHyperlink_Focus(t *testing.T) {
	app := test.NewTempApp(t)
	app.Settings().SetTheme(test.Theme())

	hyperlink := &Hyperlink{Text: "Test"}
	w := test.NewWindow(hyperlink)
	defer w.Close()
	w.SetPadded(false)
	w.Resize(hyperlink.MinSize())

	test.AssertImageMatches(t, "hyperlink/initial.png", w.Canvas().Capture())
	hyperlink.FocusGained()
	test.AssertImageMatches(t, "hyperlink/focus.png", w.Canvas().Capture())
	hyperlink.FocusLost()
	test.AssertImageMatches(t, "hyperlink/initial.png", w.Canvas().Capture())
}

func TestHyperlink_OnTapped(t *testing.T) {
	tapped := 0
	link := &Hyperlink{Text: "Test"}
	test.Tap(link)
	assert.Equal(t, 0, tapped)

	link.OnTapped = func() {
		tapped++
	}
	test.Tap(link)
	assert.Equal(t, 1, tapped)
}

func TestHyperlink_TappedOutsideTextBoundary(t *testing.T) {
	tapped := 0
	link := &Hyperlink{Text: "Test"}
	link.OnTapped = func() {
		tapped++
	}
	link.syncSegments()
	link.Tapped(&fyne.PointEvent{
		Position: fyne.NewPos(50 /*past text boundary*/, 2),
	})
	assert.Equal(t, 0, tapped)
}

func TestHyperlink_KeyboardOnTapped(t *testing.T) {
	tapped := 0
	link := &Hyperlink{Text: "Test"}
	link.TypedKey(&fyne.KeyEvent{Name: fyne.KeySpace})
	assert.Equal(t, 0, tapped)

	link.OnTapped = func() {
		tapped++
	}
	link.TypedKey(&fyne.KeyEvent{Name: fyne.KeySpace})
	assert.Equal(t, 1, tapped)
}

func TestHyperlink_Resize(t *testing.T) {
	hyperlink := &Hyperlink{Text: "Test"}
	hyperlink.CreateRenderer()
	size := fyne.NewSize(100, 20)
	hyperlink.Resize(size)

	assert.Equal(t, size, hyperlink.Size())
	assert.Equal(t, size, hyperlink.provider.Size())
}

func TestHyperlink_SetText(t *testing.T) {
	u, err := url.Parse("https://fyne.io/")
	require.NoError(t, err)

	hyperlink := &Hyperlink{Text: "Test", URL: u}
	hyperlink.CreateRenderer()
	hyperlink.SetText("New")

	assert.Equal(t, "New", hyperlink.Text)
	assert.Equal(t, "New", richTextRenderTexts(&hyperlink.provider)[0].Text)
}

func TestHyperlink_SetUrl(t *testing.T) {
	sURL, err := url.Parse("https://github.com/fyne-io/fyne")
	require.NoError(t, err)

	// test constructor
	hyperlink := NewHyperlink("Test", sURL)
	assert.Equal(t, sURL, hyperlink.URL)

	// test setting functions
	sURL, err = url.Parse("https://fyne.io")
	require.NoError(t, err)
	hyperlink.SetURL(sURL)
	assert.Equal(t, sURL, hyperlink.URL)
}

func TestHyperlink_ThemeOverride(t *testing.T) {
	test.NewTempApp(t)
	test.ApplyTheme(t, test.Theme())

	hyperlink := &Hyperlink{Text: "Test"}
	bg := canvas.NewRectangle(color.Gray{Y: 0xc0})
	w := test.NewWindow(&fyne.Container{Layout: layout.NewStackLayout(),
		Objects: []fyne.CanvasObject{bg, hyperlink}})
	defer w.Close()
	w.SetPadded(false)
	w.Resize(hyperlink.MinSize())

	light := w.Canvas().Capture()
	test.ApplyTheme(t, test.NewTheme())
	hyperlink.Refresh()
	ugly := w.Canvas().Capture()
	assertPixelsMatch(t, false, ugly, light)

	cache.OverrideTheme(hyperlink, test.Theme())
	hyperlink.Refresh()
	override := w.Canvas().Capture()
	assertPixelsMatch(t, true, override, light)
}

func TestHyperlink_Truncate(t *testing.T) {
	hyperlink := &Hyperlink{Text: "TestingWithLongText"}
	hyperlink.CreateRenderer()
	hyperlink.Resize(fyne.NewSize(100, 20))

	texts := richTextRenderTexts(&hyperlink.provider)
	assert.Equal(t, "TestingWithLongText", texts[0].Text)

	hyperlink.Truncation = fyne.TextTruncateClip
	hyperlink.Refresh()
	texts = richTextRenderTexts(&hyperlink.provider)
	assert.Equal(t, "TestingWith", texts[0].Text)

	hyperlink.Truncation = fyne.TextTruncateEllipsis
	hyperlink.Refresh()
	texts = richTextRenderTexts(&hyperlink.provider)
	assert.Equal(t, "TestingWit…", texts[0].Text)
}

func TestHyperlink_CreateRendererDoesNotAffectSize(t *testing.T) {
	u, err := url.Parse("https://github.com/fyne-io/fyne")
	require.NoError(t, err)
	link := NewHyperlink("Test", u)
	link.Resize(link.MinSize())
	size := link.Size()
	assert.NotEqual(t, fyne.NewSize(0, 0), size)
	assert.Equal(t, size, link.MinSize())

	r := link.CreateRenderer()
	link.provider.CreateRenderer()
	assert.Equal(t, size, link.Size())
	assert.Equal(t, size, link.MinSize())
	assert.Equal(t, size, r.MinSize())
	r.Layout(size)
	assert.Equal(t, size, link.Size())
	assert.Equal(t, size, link.MinSize())
}

func assertPixelsMatch(t *testing.T, match bool, img1, img2 image.Image) {
	pix1 := img1.(*image.NRGBA).Pix
	pix2 := img2.(*image.NRGBA).Pix
	if match {
		assert.Equal(t, pix1, pix2)
	} else {
		assert.NotEqual(t, pix1, pix2)
	}
}
