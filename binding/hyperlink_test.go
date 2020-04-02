package binding_test

import (
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
	"time"

	"fyne.io/fyne/binding"
	"fyne.io/fyne/test"
	"fyne.io/fyne/widget"
)

func TestBindHyperlinkText(t *testing.T) {
	a := test.NewApp()
	defer a.Quit()
	u, err := url.Parse("https://fyne.io")
	assert.Nil(t, err)
	hyperlink := widget.NewHyperlink("hyperlink", u)
	data := &binding.StringBinding{}
	binding.BindHyperlinkText(hyperlink, data)
	data.Set("foobar")
	time.Sleep(time.Second)
	assert.Equal(t, "foobar", hyperlink.Text)
}

func TestBindHyperlinkURL(t *testing.T) {
	a := test.NewApp()
	defer a.Quit()
	u, err := url.Parse("https://fyne.io")
	assert.Nil(t, err)
	hyperlink := widget.NewHyperlink("hyperlink", u)
	data := &binding.URLBinding{}
	binding.BindHyperlinkURL(hyperlink, data)
	u, err = url.Parse("https://github.com/fyne-io/fyne")
	assert.Nil(t, err)
	data.Set(u)
	time.Sleep(time.Second)
	assert.Equal(t, u, hyperlink.URL)
}
