package widget_test

import (
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/test"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

func TestCard_SetImage(t *testing.T) {
	c := widget.NewCard("Title", "sub", widget.NewLabel("Content"))
	r := test.WidgetRenderer(c)
	assert.Equal(t, 4, len(r.Objects())) // the 3 above plus shadow

	c.SetImage(canvas.NewImageFromResource(theme.FyneLogo()))
	assert.Equal(t, 5, len(r.Objects()))
}

func TestCard_SetContent(t *testing.T) {
	c := widget.NewCard("Title", "sub", widget.NewLabel("Content"))
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
		"titles_image": {
			title:    "Title",
			subtitle: "Subtitle",
			icon:     canvas.NewImageFromResource(theme.FyneLogo()),
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
			content:  newContentRect(),
		},
		"title_content": {
			title:    "Hello",
			subtitle: "",
			icon:     nil,
			content:  newContentRect(),
		},
		"image_content": {
			title:    "",
			subtitle: "",
			icon:     canvas.NewImageFromResource(theme.FyneLogo()),
			content:  newContentRect(),
		},
		"all_items": {
			title:    "Longer title",
			subtitle: "subtitle with length",
			icon:     canvas.NewImageFromResource(theme.FyneLogo()),
			content:  newContentRect(),
		},
	} {
		t.Run(name, func(t *testing.T) {
			card := &widget.Card{
				Title:    tt.title,
				Subtitle: tt.subtitle,
				Image:    tt.icon,
				Content:  tt.content,
			}

			window := test.NewWindow(card)
			size := card.MinSize().Max(fyne.NewSize(80, 0)) // give a little width for image only tests
			window.Resize(size.Add(fyne.NewSize(theme.Padding()*2, theme.Padding()*2)))
			if tt.content != nil {
				assert.Equal(t, 10, tt.content.Size().Height)
			}
			test.AssertImageMatches(t, "card/layout_"+name+".png", window.Canvas().Capture())

			window.Close()
		})
	}
}

func TestCard_MinSize(t *testing.T) {
	content := widget.NewLabel("simple")
	card := &widget.Card{Content: content}

	inner := card.MinSize().Subtract(fyne.NewSize(theme.Padding()*3, theme.Padding()*3)) // shadow + content pad
	assert.Equal(t, content.MinSize(), inner)
}

func newContentRect() *canvas.Rectangle {
	rect := canvas.NewRectangle(color.Gray{0x66})
	rect.StrokeColor = color.Black
	rect.StrokeWidth = 2
	rect.SetMinSize(fyne.NewSize(10, 10))

	return rect
}
