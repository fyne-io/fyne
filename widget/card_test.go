package widget_test

import (
	"image/color"
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/test"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"

	"github.com/stretchr/testify/assert"
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

	for name, tt := range map[string]struct {
		title, subtitle string
		icon            *canvas.Image
		content         fyne.CanvasObject
		want            string
	}{
		"title": {
			title:    "Title",
			subtitle: "",
			icon:     nil,
			content:  nil,
			want: `
				<canvas padded size="88x64">
					<content>
						<widget pos="4,4" size="80x56" type="*widget.Card">
							<widget pos="2,2" size="76x52" type="*widget.Shadow">
								<radialGradient centerOffset="0.5,0.5" pos="-1,-1" size="1x1" startColor="shadow"/>
								<linearGradient endColor="shadow" pos="0,-1" size="76x1"/>
								<radialGradient centerOffset="-0.5,0.5" pos="76,-1" size="1x1" startColor="shadow"/>
								<linearGradient angle="270" pos="76,0" size="1x52" startColor="shadow"/>
								<radialGradient centerOffset="-0.5,-0.5" pos="76,52" size="1x1" startColor="shadow"/>
								<linearGradient pos="0,52" size="76x1" startColor="shadow"/>
								<radialGradient centerOffset="0.5,-0.5" pos="-1,52" size="1x1" startColor="shadow"/>
								<linearGradient angle="270" endColor="shadow" pos="-1,0" size="1x52"/>
							</widget>
							<text bold pos="10,10" size="60x36" textSize="23">Title</text>
							<text size="0x0"></text>
						</widget>
					</content>
				</canvas>
			`,
		},
		"subtitle": {
			title:    "",
			subtitle: "Subtitle",
			icon:     nil,
			content:  nil,
			want: `
				<canvas padded size="88x49">
					<content>
						<widget pos="4,4" size="80x41" type="*widget.Card">
							<widget pos="2,2" size="76x37" type="*widget.Shadow">
								<radialGradient centerOffset="0.5,0.5" pos="-1,-1" size="1x1" startColor="shadow"/>
								<linearGradient endColor="shadow" pos="0,-1" size="76x1"/>
								<radialGradient centerOffset="-0.5,0.5" pos="76,-1" size="1x1" startColor="shadow"/>
								<linearGradient angle="270" pos="76,0" size="1x37" startColor="shadow"/>
								<radialGradient centerOffset="-0.5,-0.5" pos="76,37" size="1x1" startColor="shadow"/>
								<linearGradient pos="0,37" size="76x1" startColor="shadow"/>
								<radialGradient centerOffset="0.5,-0.5" pos="-1,37" size="1x1" startColor="shadow"/>
								<linearGradient angle="270" endColor="shadow" pos="-1,0" size="1x37"/>
							</widget>
							<text bold size="0x0" textSize="23"></text>
							<text pos="10,10" size="60x21">Subtitle</text>
						</widget>
					</content>
				</canvas>
			`,
		},
		"titles": {
			title:    "Title",
			subtitle: "Subtitle",
			icon:     nil,
			content:  nil,
			want: `
				<canvas padded size="88x89">
					<content>
						<widget pos="4,4" size="80x81" type="*widget.Card">
							<widget pos="2,2" size="76x77" type="*widget.Shadow">
								<radialGradient centerOffset="0.5,0.5" pos="-1,-1" size="1x1" startColor="shadow"/>
								<linearGradient endColor="shadow" pos="0,-1" size="76x1"/>
								<radialGradient centerOffset="-0.5,0.5" pos="76,-1" size="1x1" startColor="shadow"/>
								<linearGradient angle="270" pos="76,0" size="1x77" startColor="shadow"/>
								<radialGradient centerOffset="-0.5,-0.5" pos="76,77" size="1x1" startColor="shadow"/>
								<linearGradient pos="0,77" size="76x1" startColor="shadow"/>
								<radialGradient centerOffset="0.5,-0.5" pos="-1,77" size="1x1" startColor="shadow"/>
								<linearGradient angle="270" endColor="shadow" pos="-1,0" size="1x77"/>
							</widget>
							<text bold pos="10,10" size="60x36" textSize="23">Title</text>
							<text pos="10,50" size="60x21">Subtitle</text>
						</widget>
					</content>
				</canvas>
			`,
		},
		"titles_image": {
			title:    "Title",
			subtitle: "Subtitle",
			icon:     canvas.NewImageFromResource(theme.FyneLogo()),
			content:  nil,
			want: `
				<canvas padded size="88x217">
					<content>
						<widget pos="4,4" size="80x209" type="*widget.Card">
							<widget pos="2,2" size="76x205" type="*widget.Shadow">
								<radialGradient centerOffset="0.5,0.5" pos="-1,-1" size="1x1" startColor="shadow"/>
								<linearGradient endColor="shadow" pos="0,-1" size="76x1"/>
								<radialGradient centerOffset="-0.5,0.5" pos="76,-1" size="1x1" startColor="shadow"/>
								<linearGradient angle="270" pos="76,0" size="1x205" startColor="shadow"/>
								<radialGradient centerOffset="-0.5,-0.5" pos="76,205" size="1x1" startColor="shadow"/>
								<linearGradient pos="0,205" size="76x1" startColor="shadow"/>
								<radialGradient centerOffset="0.5,-0.5" pos="-1,205" size="1x1" startColor="shadow"/>
								<linearGradient angle="270" endColor="shadow" pos="-1,0" size="1x205"/>
							</widget>
							<text bold pos="10,138" size="60x36" textSize="23">Title</text>
							<text pos="10,178" size="60x21">Subtitle</text>
							<image pos="2,2" rsc="fyneLogo" size="76x128"/>
						</widget>
					</content>
				</canvas>
			`,
		},
		"just_image": {
			title:    "",
			subtitle: "",
			icon:     canvas.NewImageFromResource(theme.FyneLogo()),
			content:  nil,
			want: `
				<canvas padded size="88x140">
					<content>
						<widget pos="4,4" size="80x132" type="*widget.Card">
							<widget pos="2,2" size="76x128" type="*widget.Shadow">
								<radialGradient centerOffset="0.5,0.5" pos="-1,-1" size="1x1" startColor="shadow"/>
								<linearGradient endColor="shadow" pos="0,-1" size="76x1"/>
								<radialGradient centerOffset="-0.5,0.5" pos="76,-1" size="1x1" startColor="shadow"/>
								<linearGradient angle="270" pos="76,0" size="1x128" startColor="shadow"/>
								<radialGradient centerOffset="-0.5,-0.5" pos="76,128" size="1x1" startColor="shadow"/>
								<linearGradient pos="0,128" size="76x1" startColor="shadow"/>
								<radialGradient centerOffset="0.5,-0.5" pos="-1,128" size="1x1" startColor="shadow"/>
								<linearGradient angle="270" endColor="shadow" pos="-1,0" size="1x128"/>
							</widget>
							<text bold size="0x0" textSize="23"></text>
							<text size="0x0"></text>
							<image pos="2,2" rsc="fyneLogo" size="76x128"/>
						</widget>
					</content>
				</canvas>
			`,
		},
		"just_content": {
			title:    "",
			subtitle: "",
			icon:     nil,
			content:  newContentRect(),
			want: `
				<canvas padded size="88x30">
					<content>
						<widget pos="4,4" size="80x22" type="*widget.Card">
							<widget pos="2,2" size="76x18" type="*widget.Shadow">
								<radialGradient centerOffset="0.5,0.5" pos="-1,-1" size="1x1" startColor="shadow"/>
								<linearGradient endColor="shadow" pos="0,-1" size="76x1"/>
								<radialGradient centerOffset="-0.5,0.5" pos="76,-1" size="1x1" startColor="shadow"/>
								<linearGradient angle="270" pos="76,0" size="1x18" startColor="shadow"/>
								<radialGradient centerOffset="-0.5,-0.5" pos="76,18" size="1x1" startColor="shadow"/>
								<linearGradient pos="0,18" size="76x1" startColor="shadow"/>
								<radialGradient centerOffset="0.5,-0.5" pos="-1,18" size="1x1" startColor="shadow"/>
								<linearGradient angle="270" endColor="shadow" pos="-1,0" size="1x18"/>
							</widget>
							<text bold size="0x0" textSize="23"></text>
							<text size="0x0"></text>
							<rectangle fillColor="rgba(102,102,102,255)" pos="6,6" size="68x10" strokeColor="rgba(0,0,0,255)" strokeWidth="2"/>
						</widget>
					</content>
				</canvas>
			`,
		},
		"title_content": {
			title:    "Hello",
			subtitle: "",
			icon:     nil,
			content:  newContentRect(),
			want: `
				<canvas padded size="95x82">
					<content>
						<widget pos="4,4" size="87x74" type="*widget.Card">
							<widget pos="2,2" size="83x70" type="*widget.Shadow">
								<radialGradient centerOffset="0.5,0.5" pos="-1,-1" size="1x1" startColor="shadow"/>
								<linearGradient endColor="shadow" pos="0,-1" size="83x1"/>
								<radialGradient centerOffset="-0.5,0.5" pos="83,-1" size="1x1" startColor="shadow"/>
								<linearGradient angle="270" pos="83,0" size="1x70" startColor="shadow"/>
								<radialGradient centerOffset="-0.5,-0.5" pos="83,70" size="1x1" startColor="shadow"/>
								<linearGradient pos="0,70" size="83x1" startColor="shadow"/>
								<radialGradient centerOffset="0.5,-0.5" pos="-1,70" size="1x1" startColor="shadow"/>
								<linearGradient angle="270" endColor="shadow" pos="-1,0" size="1x70"/>
							</widget>
							<text bold pos="10,10" size="67x36" textSize="23">Hello</text>
							<text size="0x0"></text>
							<rectangle fillColor="rgba(102,102,102,255)" pos="6,58" size="75x10" strokeColor="rgba(0,0,0,255)" strokeWidth="2"/>
						</widget>
					</content>
				</canvas>
			`,
		},
		"image_content": {
			title:    "",
			subtitle: "",
			icon:     canvas.NewImageFromResource(theme.FyneLogo()),
			content:  newContentRect(),
			want: `
				<canvas padded size="88x158">
					<content>
						<widget pos="4,4" size="80x150" type="*widget.Card">
							<widget pos="2,2" size="76x146" type="*widget.Shadow">
								<radialGradient centerOffset="0.5,0.5" pos="-1,-1" size="1x1" startColor="shadow"/>
								<linearGradient endColor="shadow" pos="0,-1" size="76x1"/>
								<radialGradient centerOffset="-0.5,0.5" pos="76,-1" size="1x1" startColor="shadow"/>
								<linearGradient angle="270" pos="76,0" size="1x146" startColor="shadow"/>
								<radialGradient centerOffset="-0.5,-0.5" pos="76,146" size="1x1" startColor="shadow"/>
								<linearGradient pos="0,146" size="76x1" startColor="shadow"/>
								<radialGradient centerOffset="0.5,-0.5" pos="-1,146" size="1x1" startColor="shadow"/>
								<linearGradient angle="270" endColor="shadow" pos="-1,0" size="1x146"/>
							</widget>
							<text bold size="0x0" textSize="23"></text>
							<text size="0x0"></text>
							<image pos="2,2" rsc="fyneLogo" size="76x128"/>
							<rectangle fillColor="rgba(102,102,102,255)" pos="6,134" size="68x10" strokeColor="rgba(0,0,0,255)" strokeWidth="2"/>
						</widget>
					</content>
				</canvas>
			`,
		},
		"all_items": {
			title:    "Longer title",
			subtitle: "subtitle with length",
			icon:     canvas.NewImageFromResource(theme.FyneLogo()),
			content:  newContentRect(),
			want: `
				<canvas padded size="179x235">
					<content>
						<widget pos="4,4" size="171x227" type="*widget.Card">
							<widget pos="2,2" size="167x223" type="*widget.Shadow">
								<radialGradient centerOffset="0.5,0.5" pos="-1,-1" size="1x1" startColor="shadow"/>
								<linearGradient endColor="shadow" pos="0,-1" size="167x1"/>
								<radialGradient centerOffset="-0.5,0.5" pos="167,-1" size="1x1" startColor="shadow"/>
								<linearGradient angle="270" pos="167,0" size="1x223" startColor="shadow"/>
								<radialGradient centerOffset="-0.5,-0.5" pos="167,223" size="1x1" startColor="shadow"/>
								<linearGradient pos="0,223" size="167x1" startColor="shadow"/>
								<radialGradient centerOffset="0.5,-0.5" pos="-1,223" size="1x1" startColor="shadow"/>
								<linearGradient angle="270" endColor="shadow" pos="-1,0" size="1x223"/>
							</widget>
							<text bold pos="10,138" size="151x36" textSize="23">Longer title</text>
							<text pos="10,178" size="151x21">subtitle with length</text>
							<image pos="2,2" rsc="fyneLogo" size="167x128"/>
							<rectangle fillColor="rgba(102,102,102,255)" pos="6,211" size="159x10" strokeColor="rgba(0,0,0,255)" strokeWidth="2"/>
						</widget>
					</content>
				</canvas>
			`,
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
				assert.Equal(t, float32(10), tt.content.Size().Height)
			}
			test.AssertRendersToMarkup(t, tt.want, window.Canvas())

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
