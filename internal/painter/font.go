package painter

import (
	"bytes"
	"image/color"
	"image/draw"
	"math"
	"strings"
	"sync"

	"github.com/go-text/render"
	"github.com/go-text/typesetting/di"
	"github.com/go-text/typesetting/font"
	"github.com/go-text/typesetting/shaping"
	"golang.org/x/image/math/fixed"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/cache"
	"fyne.io/fyne/v2/theme"
)

const (
	// DefaultTabWidth is the default width in spaces
	DefaultTabWidth = 4

	fontTabSpaceSize = 10
)

// CachedFontFace returns a Font face held in memory. These are loaded from the current theme.
func CachedFontFace(style fyne.TextStyle, fontDP float32, texScale float32) *FontCacheItem {
	val, ok := fontCache.Load(style)
	if !ok {
		var f1, f2 font.Face
		switch {
		case style.Monospace:
			f1 = loadMeasureFont(theme.TextMonospaceFont())
			f2 = loadMeasureFont(theme.DefaultTextMonospaceFont())
		case style.Bold:
			if style.Italic {
				f1 = loadMeasureFont(theme.TextBoldItalicFont())
				f2 = loadMeasureFont(theme.DefaultTextBoldItalicFont())
			} else {
				f1 = loadMeasureFont(theme.TextBoldFont())
				f2 = loadMeasureFont(theme.DefaultTextBoldFont())
			}
		case style.Italic:
			f1 = loadMeasureFont(theme.TextItalicFont())
			f2 = loadMeasureFont(theme.DefaultTextItalicFont())
		case style.Symbol:
			f1 = loadMeasureFont(theme.SymbolFont())
			f2 = loadMeasureFont(theme.DefaultSymbolFont())
		default:
			f1 = loadMeasureFont(theme.TextFont())
			f2 = loadMeasureFont(theme.DefaultTextFont())
		}

		if f1 == nil {
			f1 = f2
		}
		faces := []font.Face{f1, f2}
		if emoji := theme.DefaultEmojiFont(); emoji != nil {
			faces = append(faces, loadMeasureFont(emoji))
		}
		val = &FontCacheItem{Fonts: faces}
		fontCache.Store(style, val)
	}

	return val.(*FontCacheItem)
}

// ClearFontCache is used to remove cached fonts in the case that we wish to re-load Font faces
func ClearFontCache() {

	fontCache = &sync.Map{}
}

// DrawString draws a string into an image.
func DrawString(dst draw.Image, s string, color color.Color, f []font.Face, fontSize, scale float32, tabWidth int) {
	r := render.Renderer{
		FontSize: fontSize,
		PixScale: scale,
		Color:    color,
	}

	// TODO avoid shaping twice!
	sh := &shaping.HarfbuzzShaper{}
	out := sh.Shape(shaping.Input{
		Text:     []rune(s),
		RunStart: 0,
		RunEnd:   len(s),
		Face:     f[0],
		Size:     fixed.I(int(fontSize * r.PixScale)),
	})

	advance := float32(0)
	y := int(math.Ceil(float64(fixed266ToFloat32(out.LineBounds.Ascent))))
	walkString(f, s, float32ToFixed266(fontSize), tabWidth, &advance, scale, func(run shaping.Output, x float32) {
		if len(run.Glyphs) == 1 {
			if run.Glyphs[0].GlyphID == 0 {
				r.DrawStringAt(string([]rune{0xfffd}), dst, int(x), y, f[0])
				return
			}
		}

		r.DrawShapedRunAt(run, dst, int(x), y)
	})
}

func loadMeasureFont(data fyne.Resource) font.Face {
	loaded, err := font.ParseTTF(bytes.NewReader(data.Content()))
	if err != nil {
		fyne.LogError("font load error", err)
		return nil
	}

	return loaded
}

// MeasureString returns how far dot would advance by drawing s with f.
// Tabs are translated into a dot location change.
func MeasureString(f []font.Face, s string, textSize float32, tabWidth int) (size fyne.Size, advance float32) {
	return walkString(f, s, float32ToFixed266(textSize), tabWidth, &advance, 1, func(shaping.Output, float32) {})
}

// RenderedTextSize looks up how big a string would be if drawn on screen.
// It also returns the distance from top to the text baseline.
func RenderedTextSize(text string, fontSize float32, style fyne.TextStyle) (size fyne.Size, baseline float32) {
	size, base := cache.GetFontMetrics(text, fontSize, style)
	if base != 0 {
		return size, base
	}

	size, base = measureText(text, fontSize, style)
	cache.SetFontMetrics(text, fontSize, style, size, base)
	return size, base
}

func fixed266ToFloat32(i fixed.Int26_6) float32 {
	return float32(float64(i) / (1 << 6))
}

func float32ToFixed266(f float32) fixed.Int26_6 {
	return fixed.Int26_6(float64(f) * (1 << 6))
}

func measureText(text string, fontSize float32, style fyne.TextStyle) (fyne.Size, float32) {
	face := CachedFontFace(style, fontSize, 1)
	return MeasureString(face.Fonts, text, fontSize, style.TabWidth)
}

func tabStop(spacew, x float32, tabWidth int) float32 {
	if tabWidth <= 0 {
		tabWidth = DefaultTabWidth
	}

	tabw := spacew * float32(tabWidth)
	tabs, _ := math.Modf(float64((x + tabw) / tabw))
	return tabw * float32(tabs)
}

func walkString(faces []font.Face, s string, textSize fixed.Int26_6, tabWidth int, advance *float32, scale float32,
	cb func(run shaping.Output, x float32)) (size fyne.Size, base float32) {
	s = strings.ReplaceAll(s, "\r", "")

	runes := []rune(s)
	in := shaping.Input{
		Text:      []rune{' '},
		RunStart:  0,
		RunEnd:    1,
		Direction: di.DirectionLTR,
		Face:      faces[0],
		Size:      textSize,
	}
	shaper := &shaping.HarfbuzzShaper{}
	out := shaper.Shape(in)

	in.Text = runes
	in.RunStart = 0
	in.RunEnd = len(runes)

	x := float32(0)
	spacew := scale * fontTabSpaceSize
	ins := shaping.SplitByFontGlyphs(in, faces)
	for _, in := range ins {
		inEnd := in.RunEnd

		pending := false
		for i, r := range in.Text[in.RunStart:in.RunEnd] {
			if r == '\t' {
				if pending {
					in.RunEnd = i
					out = shaper.Shape(in)
					x = shapeCallback(shaper, in, out, x, scale, cb)
				}
				x = tabStop(spacew, x, tabWidth)

				in.RunStart = i + 1
				in.RunEnd = inEnd
				pending = false
			} else {
				pending = true
			}
		}

		x = shapeCallback(shaper, in, out, x, scale, cb)
	}

	*advance = x
	return fyne.NewSize(*advance, fixed266ToFloat32(out.LineBounds.LineThickness())),
		fixed266ToFloat32(out.LineBounds.Ascent)
}

func shapeCallback(shaper shaping.Shaper, in shaping.Input, out shaping.Output, x, scale float32, cb func(shaping.Output, float32)) float32 {
	out = shaper.Shape(in)
	glyphs := out.Glyphs
	start := 0
	pending := false
	adv := fixed.I(0)
	for i, g := range out.Glyphs {
		if g.GlyphID == 0 {
			if pending {
				out.Glyphs = glyphs[start:i]
				cb(out, x)
				x += fixed266ToFloat32(adv) * scale
				adv = 0
			}

			out.Glyphs = glyphs[i : i+1]
			cb(out, x)
			x += fixed266ToFloat32(glyphs[i].XAdvance) * scale
			adv = 0

			start = i + 1
			pending = false
		} else {
			pending = true
		}
		adv += g.XAdvance
	}

	if pending {
		out.Glyphs = glyphs[start:]
		cb(out, x)
		x += fixed266ToFloat32(adv) * scale
		adv = 0
	}
	return x + fixed266ToFloat32(adv)*scale
}

type FontCacheItem struct {
	Fonts []font.Face
}

var fontCache = &sync.Map{} // map[fyne.TextStyle]*FontCacheItem
