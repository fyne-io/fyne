package widget

import (
	"net/url"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	_ "fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHyperlink_MinSize(t *testing.T) {
	u, err := url.Parse("https://fyne.io/")
	assert.Nil(t, err)

	hyperlink := NewHyperlink("Test", u)
	minA := hyperlink.MinSize()

	assert.Less(t, theme.Padding()*2, minA.Width)

	hyperlink.SetText("Longer")
	minB := hyperlink.MinSize()
	assert.Less(t, minA.Width, minB.Width)

	hyperlink.Text = "."
	minC := hyperlink.MinSize()
	assert.Greater(t, minB.Width, minC.Width)
}

func TestHyperlink_Cursor(t *testing.T) {
	u, err := url.Parse("https://fyne.io/")
	hyperlink := NewHyperlink("Test", u)

	assert.Nil(t, err)
	assert.Equal(t, desktop.PointerCursor, hyperlink.Cursor())
}

func TestHyperlink_Alignment(t *testing.T) {
	hyperlink := &Hyperlink{Text: "Test", Alignment: fyne.TextAlignTrailing}
	assert.Equal(t, fyne.TextAlignTrailing, textRenderTexts(hyperlink)[0].Alignment)
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
	assert.Nil(t, err)

	hyperlink := &Hyperlink{Text: "Test", URL: u}
	hyperlink.Refresh()
	hyperlink.SetText("New")

	assert.Equal(t, "New", hyperlink.Text)
	assert.Equal(t, "New", textRenderTexts(hyperlink)[0].Text)
}

func TestHyperlink_SetUrl(t *testing.T) {
	sURL, err := url.Parse("https://github.com/fyne-io/fyne")
	assert.Nil(t, err)

	// test constructor
	hyperlink := NewHyperlink("Test", sURL)
	assert.Equal(t, sURL, hyperlink.URL)

	// test setting functions
	sURL, err = url.Parse("https://fyne.io")
	assert.Nil(t, err)
	hyperlink.SetURL(sURL)
	assert.Equal(t, sURL, hyperlink.URL)
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
	assert.Equal(t, size, link.Size())
	assert.Equal(t, size, link.MinSize())
	assert.Equal(t, size, r.MinSize())
	r.Layout(size)
	assert.Equal(t, size, link.Size())
	assert.Equal(t, size, link.MinSize())
}
