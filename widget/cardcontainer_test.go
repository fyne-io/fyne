package widget_test

import (
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/test"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"github.com/stretchr/testify/assert"
)

func TestCard_SetImage(t *testing.T) {
	c := widget.NewCardContainer("Title", "sub", widget.NewLabel("Content"))
	r := test.WidgetRenderer(c)
	assert.Equal(t, 4, len(r.Objects())) // the 3 above plus shadow

	c.SetImage(canvas.NewImageFromResource(theme.FyneLogo()))
	assert.Equal(t, 5, len(r.Objects()))
}

func TestCard_SetContent(t *testing.T) {
	c := widget.NewCardContainer("Title", "sub", widget.NewLabel("Content"))
	r := test.WidgetRenderer(c)
	assert.Equal(t, 4, len(r.Objects())) // the 3 above plus shadow

	newContent := widget.NewLabel("New")
	c.SetContent(newContent)
	assert.Equal(t, 4, len(r.Objects()))
	assert.Equal(t, newContent, r.Objects()[3])
}

func TestCard_Layout(t *testing.T) {
	test.NewApp()
	test.ApplyTheme(t, theme.LightTheme())

	for name, tt := range map[string]struct {
		title, subtitle string
		icon            *canvas.Image
		content         fyne.CanvasObject
	}{
		"title": {
			title:    "Title",
			subtitle: "",
			icon:     nil,
			content:  nil,
		},
		"subtitle": {
			title:    "",
			subtitle: "Subtitle",
			icon:     nil,
			content:  nil,
		},
		"titles": {
			title:    "Title",
			subtitle: "Subtitle",
			icon:     nil,
			content:  nil,
		},
		"just_image": {
			title:    "",
			subtitle: "",
			icon:     canvas.NewImageFromResource(theme.FyneLogo()),
			content:  nil,
		},
		"just_content": {
			title:    "",
			subtitle: "",
			icon:     nil,
			content:  widget.NewHyperlink("link", nil),
		},
		"title_content": {
			title:    "Hello",
			subtitle: "",
			icon:     nil,
			content:  widget.NewHyperlink("link", nil),
		},
		"image_content": {
			title:    "",
			subtitle: "",
			icon:     canvas.NewImageFromResource(theme.FyneLogo()),
			content:  widget.NewHyperlink("link", nil),
		},
	} {
		t.Run(name, func(t *testing.T) {
			card := &widget.CardContainer{
				Title:    tt.title,
				SubTitle: tt.subtitle,
				Image:    tt.icon,
				Content:  tt.content,
			}

			window := test.NewWindow(card)
			size := card.MinSize().Max(fyne.NewSize(80, 0)) // give a little width for image only tests
			window.Resize(size.Add(fyne.NewSize(theme.Padding()*2, theme.Padding()*2)))

			test.AssertImageMatches(t, "card/layout_"+name+".png", window.Canvas().Capture())

			window.Close()
		})
	}
}
