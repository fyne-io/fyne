package widget

import (
	"net/url"
	"testing"

	"fyne.io/fyne"
	_ "fyne.io/fyne/test"
	"fyne.io/fyne/theme"
	"github.com/stretchr/testify/assert"
)

func TestHyperlink_MinSize(t *testing.T) {
	hyperlink := NewHyperlink("Test", "url1")
	min := hyperlink.MinSize()

	assert.True(t, min.Width > theme.Padding()*2)

	hyperlink.SetText("Longer")
	assert.True(t, hyperlink.MinSize().Width > min.Width)
}

func TestHyperlink_Alignment(t *testing.T) {
	hyperlink := &Hyperlink{Text: "Test", Alignment: fyne.TextAlignTrailing}
	assert.Equal(t, fyne.TextAlignTrailing, Renderer(hyperlink).(*textRenderer).texts[0].Alignment)
}

func TestHyperlink_SetText(t *testing.T) {
	var text string = "dummy text"

	// test constructor
	hyperlink := NewHyperlink(text, "")
	assert.Equal(t, text, hyperlink.Text)

	text = "new dummy text!!"
	hyperlink.SetText(text)
	assert.Equal(t, text, hyperlink.Text)
}

func TestHyperlink_SetUrl(t *testing.T) {
	var sUrl string = "url1"

	// test constructor
	hyperlink := NewHyperlink("Test", sUrl)
	assert.Equal(t, sUrl, hyperlink.Url)

	// test setting functions
	sUrl = "https://fyne.io"
	hyperlink.SetUrl(sUrl)
	assert.Equal(t, sUrl, hyperlink.Url)
	sUrl = "duck.com"
	uUrl, err := url.Parse(sUrl)
	assert.Nil(t, err) // should always be nil since we're testing, but good to check in case url changes in the future
	hyperlink.SetUrlFromUrl(uUrl)
	assert.Equal(t, sUrl, hyperlink.Url)
}
