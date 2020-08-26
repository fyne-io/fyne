package widget_test

import (
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/test"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

func TestCard_Layout(t *testing.T) {
	test.NewApp()
	test.ApplyTheme(t, theme.LightTheme())

	for name, tt := range map[string]struct {
		title, subtitle string
		icon            *canvas.Image
		content         fyne.CanvasObject
	}{
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
			size := card.MinSize().Max(fyne.NewSize(120, 80))
			window.Resize(size.Add(fyne.NewSize(theme.Padding()*2, theme.Padding()*2)))

			test.AssertImageMatches(t, "card/layout_"+name+".png", window.Canvas().Capture())

			window.Close()
		})
	}
}
