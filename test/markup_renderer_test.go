package test

import (
	"image"
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func Test_snapshot(t *testing.T) {
	// content elements
	for name, tt := range map[string]struct {
		content  fyne.CanvasObject
		size     fyne.Size
		pos      fyne.Position
		overlays []fyne.CanvasObject
		want     string
	}{
		"circle": {
			content: canvas.NewCircle(color.RGBA{R: 200, G: 100, B: 0, A: 50}),
			size:    fyne.NewSize(42, 42),
			pos:     fyne.NewPos(17, 17),
			want: "<canvas size=\"100x100\">\n" +
				"\t<content>\n" +
				"\t\t<circle fillColor=\"rgba(200,100,0,50)\" pos=\"17,17\" size=\"42x42\"/>\n" +
				"\t</content>\n" +
				"</canvas>\n",
		},
		"circle with theme color": { // we won’t test _all_ valid values … it’s not that important
			content: canvas.NewCircle(theme.ScrollBarColor()),
			want: "<canvas size=\"100x100\">\n" +
				"\t<content>\n" +
				"\t\t<circle fillColor=\"scrollbar\" size=\"100x100\"/>\n" +
				"\t</content>\n" +
				"</canvas>\n",
		},
		"circle with named primary color": { // we won’t test _all_ valid values … it’s not that important
			content: canvas.NewCircle(theme.PrimaryColorNamed(theme.ColorPurple)),
			want: "<canvas size=\"100x100\">\n" +
				"\t<content>\n" +
				"\t\t<circle fillColor=\"purple\" size=\"100x100\"/>\n" +
				"\t</content>\n" +
				"</canvas>\n",
		},
		"circle with stroke": {
			content: func() fyne.CanvasObject {
				c := canvas.NewCircle(color.RGBA{R: 200, G: 100, B: 0, A: 50})
				c.StrokeWidth = 3
				c.StrokeColor = theme.BackgroundColor()
				return c
			}(),
			size: fyne.NewSize(42, 42),
			pos:  fyne.NewPos(17, 17),
			want: "<canvas size=\"100x100\">\n" +
				"\t<content>\n" +
				"\t\t<circle fillColor=\"rgba(200,100,0,50)\" pos=\"17,17\" size=\"42x42\" strokeColor=\"background\" strokeWidth=\"3\"/>\n" +
				"\t</content>\n" +
				"</canvas>\n",
		},
		"image from file": {
			content: canvas.NewImageFromFile("testdata/markup_renderer_image.svg"),
			size:    fyne.NewSize(66, 33),
			pos:     fyne.NewPos(10, 20),
			want: "<canvas size=\"100x100\">\n" +
				"\t<content>\n" +
				"\t\t<image file=\"testdata/markup_renderer_image.svg\" pos=\"10,20\" size=\"66x33\"/>\n" +
				"\t</content>\n" +
				"</canvas>\n",
		},
		"image from image": {
			content: canvas.NewImageFromImage(image.NewUniform(color.Transparent)),
			want: "<canvas size=\"100x100\">\n" +
				"\t<content>\n" +
				"\t\t<image img size=\"100x100\"/>\n" +
				"\t</content>\n" +
				"</canvas>\n",
		},
		"image from resource": {
			content: canvas.NewImageFromResource(fyne.NewStaticResource("resource name", []byte{})),
			want: "<canvas size=\"100x100\">\n" +
				"\t<content>\n" +
				"\t\t<image rsc=\"resource name\" size=\"100x100\"/>\n" +
				"\t</content>\n" +
				"</canvas>\n",
		},
		"image theme resource": { // we won’t test _all_ valid values … it’s not that important
			content: canvas.NewImageFromResource(theme.VolumeDownIcon()),
			want: "<canvas size=\"100x100\">\n" +
				"\t<content>\n" +
				"\t\t<image rsc=\"volumeDownIcon\" size=\"100x100\"/>\n" +
				"\t</content>\n" +
				"</canvas>\n",
		},
		"image themed resource": {
			content: canvas.NewImageFromResource(theme.NewThemedResource(fyne.NewStaticResource("resource name", []byte{}))),
			want: "<canvas size=\"100x100\">\n" +
				"\t<content>\n" +
				"\t\t<image rsc=\"resource name\" size=\"100x100\" themed=\"default\"/>\n" +
				"\t</content>\n" +
				"</canvas>\n",
		},
		"image disabled resource": {
			content: canvas.NewImageFromResource(theme.NewDisabledResource(theme.CancelIcon())),
			want: "<canvas size=\"100x100\">\n" +
				"\t<content>\n" +
				"\t\t<image rsc=\"cancelIcon\" size=\"100x100\" themed=\"disabled\"/>\n" +
				"\t</content>\n" +
				"</canvas>\n",
		},
		"image error themed resource": {
			content: canvas.NewImageFromResource(theme.NewErrorThemedResource(theme.ErrorIcon())),
			want: "<canvas size=\"100x100\">\n" +
				"\t<content>\n" +
				"\t\t<image rsc=\"errorIcon\" size=\"100x100\" themed=\"error\"/>\n" +
				"\t</content>\n" +
				"</canvas>\n",
		},
		"image inverted themed resource": {
			content: canvas.NewImageFromResource(theme.NewInvertedThemedResource(theme.WarningIcon())),
			want: "<canvas size=\"100x100\">\n" +
				"\t<content>\n" +
				"\t\t<image rsc=\"warningIcon\" size=\"100x100\" themed=\"inverted\"/>\n" +
				"\t</content>\n" +
				"</canvas>\n",
		},
		"image primary themed resource": {
			content: canvas.NewImageFromResource(theme.NewPrimaryThemedResource(theme.VolumeDownIcon())),
			want: "<canvas size=\"100x100\">\n" +
				"\t<content>\n" +
				"\t\t<image rsc=\"volumeDownIcon\" size=\"100x100\" themed=\"primary\"/>\n" +
				"\t</content>\n" +
				"</canvas>\n",
		},
		"image with translucency": {
			content: func() fyne.CanvasObject {
				img := canvas.NewImageFromResource(theme.ZoomOutIcon())
				img.Translucency = 1.3
				return img
			}(),
			want: "<canvas size=\"100x100\">\n" +
				"\t<content>\n" +
				"\t\t<image rsc=\"zoomOutIcon\" size=\"100x100\" translucency=\"1.3\"/>\n" +
				"\t</content>\n" +
				"</canvas>\n",
		},
		"image with fill mode contain": {
			content: func() fyne.CanvasObject {
				img := canvas.NewImageFromResource(theme.ZoomOutIcon())
				img.FillMode = canvas.ImageFillContain
				return img
			}(),
			want: "<canvas size=\"100x100\">\n" +
				"\t<content>\n" +
				"\t\t<image fillMode=\"contain\" rsc=\"zoomOutIcon\" size=\"100x100\"/>\n" +
				"\t</content>\n" +
				"</canvas>\n",
		},
		"image with fill mode original": {
			content: func() fyne.CanvasObject {
				img := canvas.NewImageFromResource(theme.ZoomInIcon())
				img.FillMode = canvas.ImageFillOriginal
				return img
			}(),
			want: "<canvas size=\"100x100\">\n" +
				"\t<content>\n" +
				"\t\t<image fillMode=\"original\" rsc=\"zoomInIcon\" size=\"100x100\"/>\n" +
				"\t</content>\n" +
				"</canvas>\n",
		},
		"image with scale mode pixels": {
			content: func() fyne.CanvasObject {
				img := canvas.NewImageFromResource(theme.ViewRestoreIcon())
				img.ScaleMode = canvas.ImageScalePixels
				return img
			}(),
			want: "<canvas size=\"100x100\">\n" +
				"\t<content>\n" +
				"\t\t<image rsc=\"viewRestoreIcon\" scaleMode=\"pixels\" size=\"100x100\"/>\n" +
				"\t</content>\n" +
				"</canvas>\n",
		},
		"image with inline icon size": {
			content: canvas.NewImageFromResource(theme.VisibilityIcon()),
			size:    fyne.NewSize(theme.IconInlineSize(), theme.IconInlineSize()),
			want: "<canvas size=\"100x100\">\n" +
				"\t<content>\n" +
				"\t\t<image rsc=\"visibilityIcon\" size=\"iconInlineSize\"/>\n" +
				"\t</content>\n" +
				"</canvas>\n",
		},
		"line": {
			content: canvas.NewLine(color.RGBA{R: 17, G: 42, B: 128, A: 255}),
			pos:     fyne.NewPos(13, 13),
			size:    fyne.NewSize(10, 50),
			want: "<canvas size=\"100x100\">\n" +
				"\t<content>\n" +
				"\t\t<line pos=\"13,13\" size=\"10x50\" strokeColor=\"rgba(17,42,128,255)\"/>\n" +
				"\t</content>\n" +
				"</canvas>\n",
		},
		"line with theme color": { // we won’t test _all_ valid values … it’s not that important
			content: canvas.NewLine(theme.ShadowColor()),
			want: "<canvas size=\"100x100\">\n" +
				"\t<content>\n" +
				"\t\t<line size=\"100x100\" strokeColor=\"shadow\"/>\n" +
				"\t</content>\n" +
				"</canvas>\n",
		},
		"line with named primary color": { // we won’t test _all_ valid values … it’s not that important
			content: canvas.NewLine(theme.PrimaryColorNamed(theme.ColorBrown)),
			want: "<canvas size=\"100x100\">\n" +
				"\t<content>\n" +
				"\t\t<line size=\"100x100\" strokeColor=\"brown\"/>\n" +
				"\t</content>\n" +
				"</canvas>\n",
		},
		"line with stroke width": {
			content: func() fyne.CanvasObject {
				l := canvas.NewLine(color.NRGBA{R: 17, G: 42, B: 128, A: 255})
				l.StrokeWidth = 1.125
				return l
			}(),
			pos:  fyne.NewPos(13, 13),
			size: fyne.NewSize(10, 50),
			want: "<canvas size=\"100x100\">\n" +
				"\t<content>\n" +
				"\t\t<line pos=\"13,13\" size=\"10x50\" strokeColor=\"rgba(17,42,128,255)\" strokeWidth=\"1.125\"/>\n" +
				"\t</content>\n" +
				"</canvas>\n",
		},
		"linear gradient": {
			content: canvas.NewLinearGradient(color.RGBA{R: 1, G: 2, B: 3, A: 4}, theme.DisabledColor(), 13.25),
			pos:     fyne.NewPos(6, 13),
			size:    fyne.NewSize(50, 10),
			want: "<canvas size=\"100x100\">\n" +
				"\t<content>\n" +
				"\t\t<linearGradient angle=\"13.25\" endColor=\"disabled\" pos=\"6,13\" size=\"50x10\" startColor=\"rgba(1,2,3,4)\"/>\n" +
				"\t</content>\n" +
				"</canvas>\n",
		},
		"radial gradient": {
			content: canvas.NewRadialGradient(color.RGBA{R: 1, G: 2, B: 3, A: 4}, theme.DisabledColor()),
			want: "<canvas size=\"100x100\">\n" +
				"\t<content>\n" +
				"\t\t<radialGradient endColor=\"disabled\" size=\"100x100\" startColor=\"rgba(1,2,3,4)\"/>\n" +
				"\t</content>\n" +
				"</canvas>\n",
		},
		"radial gradient with offset": {
			content: func() fyne.CanvasObject {
				g := canvas.NewRadialGradient(color.RGBA{R: 1, G: 2, B: 3, A: 4}, theme.DisabledColor())
				g.CenterOffsetX = 1.5
				g.CenterOffsetY = -13.7
				return g
			}(),
			want: "<canvas size=\"100x100\">\n" +
				"\t<content>\n" +
				"\t\t<radialGradient centerOffset=\"1.5,-13.7\" endColor=\"disabled\" size=\"100x100\" startColor=\"rgba(1,2,3,4)\"/>\n" +
				"\t</content>\n" +
				"</canvas>\n",
		},
		"raster": {
			content: &canvas.Raster{Translucency: 1},
			want: "<canvas size=\"100x100\">\n" +
				"\t<content>\n" +
				"\t\t<raster size=\"100x100\" translucency=\"1\"/>\n" +
				"\t</content>\n" +
				"</canvas>\n",
		},
		"rectangle": {
			content: canvas.NewRectangle(color.RGBA{R: 100, G: 200, B: 50, A: 0}),
			size:    fyne.NewSize(17, 42),
			pos:     fyne.NewPos(42, 17),
			want: "<canvas size=\"100x100\">\n" +
				"\t<content>\n" +
				"\t\t<rectangle fillColor=\"rgba(100,200,50,0)\" pos=\"42,17\" size=\"17x42\"/>\n" +
				"\t</content>\n" +
				"</canvas>\n",
		},
		"rectangle with theme color": { // we won’t test _all_ valid values … it’s not that important
			content: canvas.NewRectangle(theme.HoverColor()),
			want: "<canvas size=\"100x100\">\n" +
				"\t<content>\n" +
				"\t\t<rectangle fillColor=\"hover\" size=\"100x100\"/>\n" +
				"\t</content>\n" +
				"</canvas>\n",
		},
		"rectangle with named primary color": { // we won’t test _all_ valid values … it’s not that important
			content: canvas.NewRectangle(theme.PrimaryColorNamed(theme.ColorOrange)),
			want: "<canvas size=\"100x100\">\n" +
				"\t<content>\n" +
				"\t\t<rectangle fillColor=\"orange\" size=\"100x100\"/>\n" +
				"\t</content>\n" +
				"</canvas>\n",
		},
		"rectangle with stroke": {
			content: func() fyne.CanvasObject {
				r := canvas.NewRectangle(color.RGBA{R: 200, G: 100, B: 0, A: 50})
				r.StrokeWidth = 6.375
				r.StrokeColor = theme.PlaceHolderColor()
				return r
			}(),
			size: fyne.NewSize(42, 42),
			pos:  fyne.NewPos(17, 17),
			want: "<canvas size=\"100x100\">\n" +
				"\t<content>\n" +
				"\t\t<rectangle fillColor=\"rgba(200,100,0,50)\" pos=\"17,17\" size=\"42x42\" strokeColor=\"placeholder\" strokeWidth=\"6.375\"/>\n" +
				"\t</content>\n" +
				"</canvas>\n",
		},
		"text": {
			content: canvas.NewText("foo", theme.PrimaryColorNamed(theme.ColorYellow)),
			size:    fyne.NewSize(50, 50),
			pos:     fyne.NewPos(20, 20),
			want: "<canvas size=\"100x100\">\n" +
				"\t<content>\n" +
				"\t\t<text color=\"yellow\" pos=\"20,20\" size=\"50x50\">foo</text>\n" +
				"\t</content>\n" +
				"</canvas>\n",
		},
		"text with text color": {
			content: canvas.NewText("bar", theme.ForegroundColor()),
			want: "<canvas size=\"100x100\">\n" +
				"\t<content>\n" +
				"\t\t<text size=\"100x100\">bar</text>\n" +
				"\t</content>\n" +
				"</canvas>\n",
		},
		"text with alignment center": {
			content: func() fyne.CanvasObject {
				t := canvas.NewText("bar", theme.ForegroundColor())
				t.Alignment = fyne.TextAlignCenter
				return t
			}(),
			want: "<canvas size=\"100x100\">\n" +
				"\t<content>\n" +
				"\t\t<text alignment=\"center\" size=\"100x100\">bar</text>\n" +
				"\t</content>\n" +
				"</canvas>\n",
		},
		"text with alignment trailing": {
			content: func() fyne.CanvasObject {
				txt := canvas.NewText("bar", theme.ForegroundColor())
				txt.Alignment = fyne.TextAlignTrailing
				return txt
			}(),
			want: "<canvas size=\"100x100\">\n" +
				"\t<content>\n" +
				"\t\t<text alignment=\"trailing\" size=\"100x100\">bar</text>\n" +
				"\t</content>\n" +
				"</canvas>\n",
		},
		"text with size": {
			content: func() fyne.CanvasObject {
				txt := canvas.NewText("big", theme.ForegroundColor())
				txt.TextSize = 42
				return txt
			}(),
			want: "<canvas size=\"100x100\">\n" +
				"\t<content>\n" +
				"\t\t<text size=\"100x100\" textSize=\"42\">big</text>\n" +
				"\t</content>\n" +
				"</canvas>\n",
		},
		"text bold": {
			content: func() fyne.CanvasObject {
				txt := canvas.NewText("bold", theme.ForegroundColor())
				txt.TextStyle.Bold = true
				return txt
			}(),
			want: "<canvas size=\"100x100\">\n" +
				"\t<content>\n" +
				"\t\t<text bold size=\"100x100\">bold</text>\n" +
				"\t</content>\n" +
				"</canvas>\n",
		},
		"text italic": {
			content: func() fyne.CanvasObject {
				txt := canvas.NewText("italic", theme.ForegroundColor())
				txt.TextStyle.Italic = true
				return txt
			}(),
			want: "<canvas size=\"100x100\">\n" +
				"\t<content>\n" +
				"\t\t<text italic size=\"100x100\">italic</text>\n" +
				"\t</content>\n" +
				"</canvas>\n",
		},
		"text monospace": {
			content: func() fyne.CanvasObject {
				txt := canvas.NewText("mono", theme.ForegroundColor())
				txt.TextStyle.Monospace = true
				return txt
			}(),
			want: "<canvas size=\"100x100\">\n" +
				"\t<content>\n" +
				"\t\t<text monospace size=\"100x100\">mono</text>\n" +
				"\t</content>\n" +
				"</canvas>\n",
		},
		"container": {
			content: container.NewVBox(canvas.NewCircle(color.Black), canvas.NewLine(color.RGBA{R: 250, G: 250, B: 250, A: 250})),
			want: "<canvas size=\"100x100\">\n" +
				"\t<content>\n" +
				"\t\t<container size=\"100x100\">\n" +
				"\t\t\t<circle fillColor=\"rgba(0,0,0,255)\" size=\"100x1\"/>\n" +
				"\t\t\t<line pos=\"0,5\" size=\"100x1\" strokeColor=\"rgba(250,250,250,250)\"/>\n" +
				"\t\t</container>\n" +
				"\t</content>\n" +
				"</canvas>\n",
		},
		"widget": {
			content: &markupRendererTestWidget{},
			want: "<canvas size=\"100x100\">\n" +
				"\t<content>\n" +
				"\t\t<widget size=\"100x100\" type=\"*test.markupRendererTestWidget\">\n" +
				"\t\t</widget>\n" +
				"\t</content>\n" +
				"</canvas>\n",
		},
		"widget with subobjects": {
			content: &markupRendererTestWidget{
				objs: []fyne.CanvasObject{
					canvas.NewCircle(color.Black),
					canvas.NewLine(color.RGBA{R: 250, G: 250, B: 250, A: 250}),
				},
			},
			want: "<canvas size=\"100x100\">\n" +
				"\t<content>\n" +
				"\t\t<widget size=\"100x100\" type=\"*test.markupRendererTestWidget\">\n" +
				"\t\t\t<circle fillColor=\"rgba(0,0,0,255)\" size=\"0x0\"/>\n" +
				"\t\t\t<line size=\"0x0\" strokeColor=\"rgba(250,250,250,250)\"/>\n" +
				"\t\t</widget>\n" +
				"\t</content>\n" +
				"</canvas>\n",
		},
		"spacer": {
			content: layout.NewSpacer(),
			want: "<canvas size=\"100x100\">\n" +
				"\t<content>\n" +
				"\t\t<spacer size=\"100x100\"/>\n" +
				"\t</content>\n" +
				"</canvas>\n",
		},
		"invisible": {
			content: func() fyne.CanvasObject {
				c := canvas.NewCircle(color.Black)
				c.Hide()
				l := canvas.NewLine(color.RGBA{R: 250, G: 250, B: 250, A: 250})
				l.Hide()
				w := widget.NewButton("tap me if you can", nil)
				w.Hide()
				return container.NewVBox(c, l, w)
			}(),
			want: "<canvas size=\"100x100\">\n" +
				"\t<content>\n" +
				"\t\t<container size=\"100x100\">\n" +
				"\t\t</container>\n" +
				"\t</content>\n" +
				"</canvas>\n",
		},
	} {
		t.Run(name, func(t *testing.T) {
			c := NewCanvas()
			c.SetPadded(false)
			c.SetContent(tt.content)
			c.Resize(fyne.NewSize(100, 100))
			if !tt.pos.IsZero() {
				tt.content.Move(tt.pos)
			}
			if !tt.size.IsZero() {
				tt.content.Resize(tt.size)
			}
			assert.Equal(t, tt.want, snapshot(c))
		})
	}

	t.Run("canvas with padding", func(t *testing.T) {
		c := NewCanvas()
		c.SetPadded(true)
		c.SetContent(canvas.NewCircle(color.Black))
		c.Resize(fyne.NewSize(100, 100))
		assert.Equal(
			t,
			"<canvas padded size=\"100x100\">\n"+
				"\t<content>\n"+
				"\t\t<circle fillColor=\"rgba(0,0,0,255)\" pos=\"4,4\" size=\"92x92\"/>\n"+
				"\t</content>\n"+
				"</canvas>\n",
			snapshot(c),
		)
	})

	t.Run("canvas with overlays", func(t *testing.T) {
		c := NewCanvas()
		c.SetPadded(true)
		c.SetContent(canvas.NewCircle(color.Black))
		c.Overlays().Add(canvas.NewRectangle(color.RGBA{R: 250, G: 250, B: 250, A: 250}))
		c.Overlays().Add(canvas.NewRectangle(color.Transparent))
		c.Resize(fyne.NewSize(100, 100))
		assert.Equal(
			t,
			"<canvas padded size=\"100x100\">\n"+
				"\t<content>\n"+
				"\t\t<circle fillColor=\"rgba(0,0,0,255)\" pos=\"4,4\" size=\"92x92\"/>\n"+
				"\t</content>\n"+
				"\t<overlay>\n"+
				"\t\t<rectangle fillColor=\"rgba(250,250,250,250)\" size=\"100x100\"/>\n"+
				"\t</overlay>\n"+
				"\t<overlay>\n"+
				"\t\t<rectangle size=\"100x100\"/>\n"+
				"\t</overlay>\n"+
				"</canvas>\n",
			snapshot(c),
		)
	})
}

type markupRendererTestWidget struct {
	hidden bool
	objs   []fyne.CanvasObject
	pos    fyne.Position
	size   fyne.Size
}

var _ fyne.Widget = (*markupRendererTestWidget)(nil)

func (w *markupRendererTestWidget) CreateRenderer() fyne.WidgetRenderer {
	return &markupRendererTestWidgetRenderer{w: w}
}

func (w *markupRendererTestWidget) Hide() {
	w.hidden = true
}

func (w *markupRendererTestWidget) MinSize() fyne.Size {
	return fyne.Size{}
}

func (w *markupRendererTestWidget) Move(position fyne.Position) {
	w.pos = position
}

func (w *markupRendererTestWidget) Position() fyne.Position {
	return w.pos
}

func (w *markupRendererTestWidget) Refresh() {
}

func (w *markupRendererTestWidget) Resize(size fyne.Size) {
	w.size = size
}

func (w *markupRendererTestWidget) SetObjects(objects ...fyne.CanvasObject) {
	w.objs = objects
}

func (w *markupRendererTestWidget) Show() {
	w.hidden = false
}

func (w *markupRendererTestWidget) Size() fyne.Size {
	return w.size
}

func (w *markupRendererTestWidget) Visible() bool {
	return !w.hidden
}

type markupRendererTestWidgetRenderer struct {
	w *markupRendererTestWidget
}

func (r *markupRendererTestWidgetRenderer) Destroy() {
}

func (r *markupRendererTestWidgetRenderer) Layout(_ fyne.Size) {
}

func (r *markupRendererTestWidgetRenderer) MinSize() fyne.Size {
	return fyne.Size{}
}

func (r *markupRendererTestWidgetRenderer) Objects() []fyne.CanvasObject {
	return r.w.objs
}

func (r *markupRendererTestWidgetRenderer) Refresh() {
}

var _ fyne.WidgetRenderer = (*markupRendererTestWidgetRenderer)(nil)
