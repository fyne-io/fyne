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
	u, err := url.Parse("https://fyne.io/")
	assert.Nil(t, err)

	hyperlink := NewHyperlink("Test", u)
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
	u, err := url.Parse("https://fyne.io/")
	assert.Nil(t, err)

	hyperlink := &Hyperlink{Text: "Test", Url: u}
	Refresh(hyperlink)
	hyperlink.SetText("New")

	assert.Equal(t, "New", hyperlink.Text)
	assert.Equal(t, "New", textRenderTexts(hyperlink)[0].Text)
}

func TestHyperlink_SetUrl(t *testing.T) {
	sUrl, err := url.Parse("https://github.com/fyne-io/fyne")
	assert.Nil(t, err)

	// test constructor
	hyperlink := NewHyperlink("Test", sUrl)
	assert.Equal(t, sUrl, hyperlink.Url)

	// test setting functions
	sUrl, err = url.Parse("https://fyne.io")
	assert.Nil(t, err)
	hyperlink.SetUrl(sUrl)
	assert.Equal(t, sUrl, hyperlink.Url)
	sUrl, err = url.Parse("duck.com")
	assert.Nil(t, err)
	uUrl := sUrl //what?
	assert.Nil(t, err)
	hyperlink.SetUrl(uUrl)
	assert.Equal(t, sUrl, hyperlink.Url)
}
