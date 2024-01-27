package tutorials

import (
	"image/color"
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/cmd/fyne_demo/data"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func parseURL(urlStr string) *url.URL {
	link, err := url.Parse(urlStr)
	if err != nil {
		fyne.LogError("Could not parse URL", err)
	}

	return link
}

func welcomeScreen(_ fyne.Window) fyne.CanvasObject {
	logo := canvas.NewImageFromResource(data.FyneLogoTransparent)
	logo.FillMode = canvas.ImageFillContain
	if fyne.CurrentDevice().IsMobile() {
		logo.SetMinSize(fyne.NewSize(192, 192))
	} else {
		logo.SetMinSize(fyne.NewSize(256, 256))
	}

	footer := container.NewHBox(
		layout.NewSpacer(),
		widget.NewHyperlink("fyne.io", parseURL("https://fyne.io/")),
		widget.NewLabel("-"),
		widget.NewHyperlink("documentation", parseURL("https://docs.fyne.io/")),
		widget.NewLabel("-"),
		widget.NewHyperlink("sponsor", parseURL("https://fyne.io/sponsor/")),
		layout.NewSpacer(),
	)

	// TODO load the AUTHORS file somehow
	authors := widget.NewRichTextFromMarkdown(`### Authors

* Andy Williams <andy@andy.xyz>
* Steve OConnor <steveoc64@gmail.com>
* Luca Corbo <lu.corbo@gmail.com>
* Paul Hovey <paul@paulhovey.org>
* Charles Corbett <nafredy@gmail.com>
* Tilo Prütz <tilo@pruetz.net>
* Stephen Houston <smhouston88@gmail.com>
* Storm Hess <stormhess@gloryskulls.com>
* Stuart Scott <stuart.murray.scott@gmail.com>
* Jacob Alzén <jacalz@tutanota.com>
* Charles A. Daniels <charles@cdaniels.net>
* Pablo Fuentes <f.pablo1@hotmail.com>
* Changkun Ou <hi@changkun.de>
`)
	content := container.NewVBox(
		widget.NewLabelWithStyle("\n\nWelcome to the Fyne toolkit demo app", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		logo,
		container.NewCenter(authors))
	scroll := container.NewScroll(content)

	bgColor := withAlpha(theme.BackgroundColor(), 0xe0)
	shadowColor := withAlpha(theme.BackgroundColor(), 0x33)

	underlay := canvas.NewImageFromResource(data.FyneLogo)
	bg := canvas.NewRectangle(bgColor)
	slideBG := container.NewWithoutLayout(underlay)
	footerBG := canvas.NewRectangle(shadowColor)

	listen := make(chan fyne.Settings)
	fyne.CurrentApp().Settings().AddChangeListener(listen)
	go func() {
		for range listen {
			bgColor = withAlpha(theme.BackgroundColor(), 0xe0)
			bg.FillColor = bgColor
			bg.Refresh()

			shadowColor = withAlpha(theme.BackgroundColor(), 0x33)
			footerBG.FillColor = bgColor
			footer.Refresh()
		}
	}()

	underlay.Resize(fyne.NewSize(1024, 1024))
	underlay.FillMode = canvas.ImageFillStretch
	underlay.Move(fyne.NewPos(-250, -50)) // TODO move it into center of this space, requires layout

	scroll.OnScrolled = func(p fyne.Position) {
		underlay.Move(fyne.NewPos(-250, -50-p.Y/3))
	}

	bgClip := container.NewScroll(slideBG)
	bgClip.Direction = container.ScrollNone
	return container.NewStack(container.New(unpad{top: true}, bgClip, bg),
		container.NewBorder(nil,
			container.NewStack(footerBG, footer), nil, nil,
			container.New(unpad{top: true, bottom: true}, scroll)))
}

func withAlpha(c color.Color, alpha uint8) color.Color {
	r, g, b, _ := c.RGBA()
	return color.NRGBA{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8), A: alpha}
}

type unpad struct {
	top, bottom bool
}

func (u unpad) Layout(objs []fyne.CanvasObject, s fyne.Size) {
	pad := theme.Padding()
	var pos fyne.Position
	if u.top {
		pos = fyne.NewPos(0, -pad)
	}
	size := s
	if u.top {
		size = size.AddWidthHeight(0, pad)
	}
	if u.bottom {
		size = size.AddWidthHeight(0, pad)
	}
	for _, o := range objs {
		o.Move(pos)
		o.Resize(size)
	}
}

func (u unpad) MinSize(_ []fyne.CanvasObject) fyne.Size {
	return fyne.NewSize(100, 100)
}
