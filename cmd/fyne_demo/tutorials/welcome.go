package tutorials

import (
	"image/color"
	"net/url"
	"strings"

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

	authors := widget.NewRichTextFromMarkdown(formatAuthors(string(data.Authors.Content())))
	content := container.NewVBox(
		widget.NewLabelWithStyle("\n\nWelcome to the Fyne toolkit demo app", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		logo,
		container.NewCenter(authors),
		widget.NewLabelWithStyle("\nWith great thanks to our many kind sponsors\n", fyne.TextAlignCenter, fyne.TextStyle{Italic: true}))
	scroll := container.NewScroll(content)

	bgColor := withAlpha(theme.Color(theme.ColorNameBackground), 0xe0)
	shadowColor := withAlpha(theme.Color(theme.ColorNameBackground), 0x33)

	underlay := canvas.NewImageFromResource(data.FyneLogo)
	bg := canvas.NewRectangle(bgColor)
	underlayer := underLayout{}
	slideBG := container.New(underlayer, underlay)
	footerBG := canvas.NewRectangle(shadowColor)

	fyne.CurrentApp().Settings().AddListener(func(fyne.Settings) {
		bgColor = withAlpha(theme.Color(theme.ColorNameBackground), 0xe0)
		bg.FillColor = bgColor
		bg.Refresh()

		shadowColor = withAlpha(theme.Color(theme.ColorNameBackground), 0x33)
		footerBG.FillColor = bgColor
		footer.Refresh()
	})

	underlay.Resize(fyne.NewSize(1024, 1024))
	scroll.OnScrolled = func(p fyne.Position) {
		underlayer.offset = -p.Y / 3
		underlayer.Layout(slideBG.Objects, slideBG.Size())
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

type underLayout struct {
	offset float32
}

func (u underLayout) Layout(objs []fyne.CanvasObject, size fyne.Size) {
	under := objs[0]
	left := size.Width/2 - under.Size().Width/2
	under.Move(fyne.NewPos(left, u.offset-50))
}

func (u underLayout) MinSize(_ []fyne.CanvasObject) fyne.Size {
	return fyne.Size{}
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

func formatAuthors(lines string) string {
	markdown := &strings.Builder{}
	markdown.WriteString("### Authors\n\n")

	for _, line := range strings.Split(lines, "\n") {
		if len(line) == 0 {
			continue
		}

		markdown.WriteString("* ")
		markdown.WriteString(line)
		markdown.WriteByte('\n')
	}

	return markdown.String()
}
