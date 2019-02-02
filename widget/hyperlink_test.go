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
	assert.Equal(t, fyne.TextAlignTrailing, textRenderTexts(hyperlink)[0].Alignment)
}

func TestHyperlink_SetText(t *testing.T) {
	hyperlink := &Hyperlink{Text: "Test", URL: "TestUrl"}
	Refresh(hyperlink)
	hyperlink.SetText("New")

	assert.Equal(t, "New", hyperlink.Text)
	assert.Equal(t, "New", textRenderTexts(hyperlink)[0].Text)
}

func TestHyperlink_SetUrl(t *testing.T) {
	var input = "url1"

	// test constructor
	hyperlink := NewHyperlink("Test", input)
	assert.Equal(t, input, hyperlink.URL)

	// test setting functions
	input = "https://fyne.io"
	hyperlink.SetURL(input)
	assert.Equal(t, input, hyperlink.URL)
	input = "duck.com"
	URL, err := url.Parse(input)
	assert.Nil(t, err)
	hyperlink.SetURLFromURL(URL)
	assert.Equal(t, input, hyperlink.URL)
}
