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
				<canvas padded size="88x54">
					<content>
						<widget pos="4,4" size="80x46" type="*widget.Card">
							<widget pos="2,2" size="76x42" type="*widget.Shadow">
								<radialGradient centerOffset="0.5,0.5" pos="-1,-1" size="1x1" startColor="shadow"/>
								<linearGradient endColor="shadow" pos="0,-1" size="76x1"/>
								<radialGradient centerOffset="-0.5,0.5" pos="76,-1" size="1x1" startColor="shadow"/>
								<linearGradient angle="270" pos="76,0" size="1x42" startColor="shadow"/>
								<radialGradient centerOffset="-0.5,-0.5" pos="76,42" size="1x1" startColor="shadow"/>
								<linearGradient pos="0,42" size="76x1" startColor="shadow"/>
								<radialGradient centerOffset="0.5,-0.5" pos="-1,42" size="1x1" startColor="shadow"/>
								<linearGradient angle="270" endColor="shadow" pos="-1,0" size="1x42"/>
							</widget>
							<text bold pos="10,6" size="60x34" textSize="23">Title</text>
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
				<canvas padded size="88x41">
					<content>
						<widget pos="4,4" size="80x33" type="*widget.Card">
							<widget pos="2,2" size="76x29" type="*widget.Shadow">
								<radialGradient centerOffset="0.5,0.5" pos="-1,-1" size="1x1" startColor="shadow"/>
								<linearGradient endColor="shadow" pos="0,-1" size="76x1"/>
								<radialGradient centerOffset="-0.5,0.5" pos="76,-1" size="1x1" startColor="shadow"/>
								<linearGradient angle="270" pos="76,0" size="1x29" startColor="shadow"/>
								<radialGradient centerOffset="-0.5,-0.5" pos="76,29" size="1x1" startColor="shadow"/>
								<linearGradient pos="0,29" size="76x1" startColor="shadow"/>
								<radialGradient centerOffset="0.5,-0.5" pos="-1,29" size="1x1" startColor="shadow"/>
								<linearGradient angle="270" endColor="shadow" pos="-1,0" size="1x29"/>
							</widget>
							<text bold size="0x0" textSize="23"></text>
							<text pos="10,6" size="60x21">Subtitle</text>
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
				<canvas padded size="88x83">
					<content>
						<widget pos="4,4" size="80x75" type="*widget.Card">
							<widget pos="2,2" size="76x71" type="*widget.Shadow">
								<radialGradient centerOffset="0.5,0.5" pos="-1,-1" size="1x1" startColor="shadow"/>
								<linearGradient endColor="shadow" pos="0,-1" size="76x1"/>
								<radialGradient centerOffset="-0.5,0.5" pos="76,-1" size="1x1" startColor="shadow"/>
								<linearGradient angle="270" pos="76,0" size="1x71" startColor="shadow"/>
								<radialGradient centerOffset="-0.5,-0.5" pos="76,71" size="1x1" startColor="shadow"/>
								<linearGradient pos="0,71" size="76x1" startColor="shadow"/>
								<radialGradient centerOffset="0.5,-0.5" pos="-1,71" size="1x1" startColor="shadow"/>
								<linearGradient angle="270" endColor="shadow" pos="-1,0" size="1x71"/>
							</widget>
							<text bold pos="10,6" size="60x34" textSize="23">Title</text>
							<text pos="10,44" size="60x21">Subtitle</text>
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
				<canvas padded size="88x211">
					<content>
						<widget pos="4,4" size="80x203" type="*widget.Card">
							<widget pos="2,2" size="76x199" type="*widget.Shadow">
								<radialGradient centerOffset="0.5,0.5" pos="-1,-1" size="1x1" startColor="shadow"/>
								<linearGradient endColor="shadow" pos="0,-1" size="76x1"/>
								<radialGradient centerOffset="-0.5,0.5" pos="76,-1" size="1x1" startColor="shadow"/>
								<linearGradient angle="270" pos="76,0" size="1x199" startColor="shadow"/>
								<radialGradient centerOffset="-0.5,-0.5" pos="76,199" size="1x1" startColor="shadow"/>
								<linearGradient pos="0,199" size="76x1" startColor="shadow"/>
								<radialGradient centerOffset="0.5,-0.5" pos="-1,199" size="1x1" startColor="shadow"/>
								<linearGradient angle="270" endColor="shadow" pos="-1,0" size="1x199"/>
							</widget>
							<text bold pos="10,134" size="60x34" textSize="23">Title</text>
							<text pos="10,172" size="60x21">Subtitle</text>
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
			content:  widget.NewHyperlink("link", nil),
			want: `
				<canvas padded size="88x57">
					<content>
						<widget pos="4,4" size="80x49" type="*widget.Card">
							<widget pos="2,2" size="76x45" type="*widget.Shadow">
								<radialGradient centerOffset="0.5,0.5" pos="-1,-1" size="1x1" startColor="shadow"/>
								<linearGradient endColor="shadow" pos="0,-1" size="76x1"/>
								<radialGradient centerOffset="-0.5,0.5" pos="76,-1" size="1x1" startColor="shadow"/>
								<linearGradient angle="270" pos="76,0" size="1x45" startColor="shadow"/>
								<radialGradient centerOffset="-0.5,-0.5" pos="76,45" size="1x1" startColor="shadow"/>
								<linearGradient pos="0,45" size="76x1" startColor="shadow"/>
								<radialGradient centerOffset="0.5,-0.5" pos="-1,45" size="1x1" startColor="shadow"/>
								<linearGradient angle="270" endColor="shadow" pos="-1,0" size="1x45"/>
							</widget>
							<text bold size="0x0" textSize="23"></text>
							<text size="0x0"></text>
							<widget pos="6,10" size="68x23" type="*widget.Hyperlink">
								<text color="focus" pos="4,4" size="60x21">link</text>
							</widget>
						</widget>
					</content>
				</canvas>
			`,
		},
		"title_content": {
			title:    "Hello",
			subtitle: "",
			icon:     nil,
			content:  widget.NewHyperlink("link", nil),
			want: `
				<canvas padded size="93x91">
					<content>
						<widget pos="4,4" size="85x83" type="*widget.Card">
							<widget pos="2,2" size="81x79" type="*widget.Shadow">
								<radialGradient centerOffset="0.5,0.5" pos="-1,-1" size="1x1" startColor="shadow"/>
								<linearGradient endColor="shadow" pos="0,-1" size="81x1"/>
								<radialGradient centerOffset="-0.5,0.5" pos="81,-1" size="1x1" startColor="shadow"/>
								<linearGradient angle="270" pos="81,0" size="1x79" startColor="shadow"/>
								<radialGradient centerOffset="-0.5,-0.5" pos="81,79" size="1x1" startColor="shadow"/>
								<linearGradient pos="0,79" size="81x1" startColor="shadow"/>
								<radialGradient centerOffset="0.5,-0.5" pos="-1,79" size="1x1" startColor="shadow"/>
								<linearGradient angle="270" endColor="shadow" pos="-1,0" size="1x79"/>
							</widget>
							<text bold pos="10,6" size="65x34" textSize="23">Hello</text>
							<text size="0x0"></text>
							<widget pos="6,48" size="73x19" type="*widget.Hyperlink">
								<text color="focus" pos="4,4" size="65x21">link</text>
							</widget>
						</widget>
					</content>
				</canvas>
			`,
		},
		"image_content": {
			title:    "",
			subtitle: "",
			icon:     canvas.NewImageFromResource(theme.FyneLogo()),
			content:  widget.NewHyperlink("link", nil),
			want: `
				<canvas padded size="88x185">
					<content>
						<widget pos="4,4" size="80x177" type="*widget.Card">
							<widget pos="2,2" size="76x173" type="*widget.Shadow">
								<radialGradient centerOffset="0.5,0.5" pos="-1,-1" size="1x1" startColor="shadow"/>
								<linearGradient endColor="shadow" pos="0,-1" size="76x1"/>
								<radialGradient centerOffset="-0.5,0.5" pos="76,-1" size="1x1" startColor="shadow"/>
								<linearGradient angle="270" pos="76,0" size="1x173" startColor="shadow"/>
								<radialGradient centerOffset="-0.5,-0.5" pos="76,173" size="1x1" startColor="shadow"/>
								<linearGradient pos="0,173" size="76x1" startColor="shadow"/>
								<radialGradient centerOffset="0.5,-0.5" pos="-1,173" size="1x1" startColor="shadow"/>
								<linearGradient angle="270" endColor="shadow" pos="-1,0" size="1x173"/>
							</widget>
							<text bold size="0x0" textSize="23"></text>
							<text size="0x0"></text>
							<image pos="2,2" rsc="fyneLogo" size="76x128"/>
							<widget pos="6,138" size="68x23" type="*widget.Hyperlink">
								<text color="focus" pos="4,4" size="60x21">link</text>
							</widget>
						</widget>
					</content>
				</canvas>
			`,
		},
		"all_items": {
			title:    "Longer title",
			subtitle: "subtitle with length",
			icon:     canvas.NewImageFromResource(theme.FyneLogo()),
			content:  widget.NewHyperlink("link", nil),
			want: `
				<canvas padded size="174x240">
					<content>
						<widget pos="4,4" size="166x232" type="*widget.Card">
							<widget pos="2,2" size="162x228" type="*widget.Shadow">
								<radialGradient centerOffset="0.5,0.5" pos="-1,-1" size="1x1" startColor="shadow"/>
								<linearGradient endColor="shadow" pos="0,-1" size="162x1"/>
								<radialGradient centerOffset="-0.5,0.5" pos="162,-1" size="1x1" startColor="shadow"/>
								<linearGradient angle="270" pos="162,0" size="1x228" startColor="shadow"/>
								<radialGradient centerOffset="-0.5,-0.5" pos="162,228" size="1x1" startColor="shadow"/>
								<linearGradient pos="0,228" size="162x1" startColor="shadow"/>
								<radialGradient centerOffset="0.5,-0.5" pos="-1,228" size="1x1" startColor="shadow"/>
								<linearGradient angle="270" endColor="shadow" pos="-1,0" size="1x228"/>
							</widget>
							<text bold pos="10,134" size="146x34" textSize="23">Longer title</text>
							<text pos="10,172" size="146x21">subtitle with length</text>
							<image pos="2,2" rsc="fyneLogo" size="162x128"/>
							<widget pos="6,201" size="154x15" type="*widget.Hyperlink">
								<text color="focus" pos="4,4" size="146x21">link</text>
							</widget>
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

			test.AssertRendersToMarkup(t, tt.want, window.Canvas())

			window.Close()
		})
	}
}
