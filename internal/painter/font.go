package painter

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"math"
	"sync"

	"github.com/go-text/typesetting/di"
	gotext "github.com/go-text/typesetting/font"
	"github.com/go-text/typesetting/shaping"
	"github.com/goki/freetype"
	"github.com/goki/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/cache"
	"fyne.io/fyne/v2/theme"
)

const (
	// DefaultTabWidth is the default width in spaces
	DefaultTabWidth = 4

	// TextDPI is a global constant that determines how text scales to interface sizes
	TextDPI = 78

	fontTabSpaceSize = 10
)

// CachedFontFace returns a font face held in memory. These are loaded from the current theme.
func CachedFontFace(style fyne.TextStyle, fontDP float32, texScale float32) (font.Face, gotext.Face) {
	key := faceCacheKey{float32ToFixed266(fontDP), float32ToFixed266(texScale)}
	val, ok := fontCache.Load(style)
	if !ok {
		var f1, f2 *truetype.Font
		switch {
		case style.Monospace:
			f1 = loadFont(theme.TextMonospaceFont())
			f2 = loadFont(theme.DefaultTextMonospaceFont())
		case style.Bold:
			if style.Italic {
				f1 = loadFont(theme.TextBoldItalicFont())
				f2 = loadFont(theme.DefaultTextBoldItalicFont())
			} else {
				f1 = loadFont(theme.TextBoldFont())
				f2 = loadFont(theme.DefaultTextBoldFont())
			}
		case style.Italic:
			f1 = loadFont(theme.TextItalicFont())
			f2 = loadFont(theme.DefaultTextItalicFont())
		case style.Symbol:
			f2 = loadFont(theme.DefaultSymbolFont())
		default:
			f1 = loadFont(theme.TextFont())
			f2 = loadFont(theme.DefaultTextFont())
		}

		if f1 == nil {
			f1 = f2
		}
		val = &fontCacheItem{font: f1, fallback: f2, faces: make(map[faceCacheKey]font.Face),
			measureFaces: make(map[faceCacheKey]gotext.Face)}
		fontCache.Store(style, val)
	}

	comp := val.(*fontCacheItem)
	comp.facesMutex.RLock()
	face := comp.faces[key]
	measureFace := comp.measureFaces[key]
	comp.facesMutex.RUnlock()
	if face == nil {
		var opts truetype.Options
		opts.Size = float64(fontDP)
		opts.DPI = float64(TextDPI * texScale)

		f1 := truetype.NewFace(comp.font, &opts)
		f2 := truetype.NewFace(comp.fallback, &opts)
		face = newFontWithFallback(f1, f2, comp.font, comp.fallback)

		switch {
		case style.Monospace:
			measureFace = loadMeasureFont(theme.TextMonospaceFont())
		case style.Bold:
			if style.Italic {
				measureFace = loadMeasureFont(theme.TextBoldItalicFont())
			} else {
				measureFace = loadMeasureFont(theme.TextBoldFont())
			}
		case style.Italic:
			measureFace = loadMeasureFont(theme.TextItalicFont())
		case style.Symbol:
			measureFace = loadMeasureFont(theme.DefaultSymbolFont())
		default:
			measureFace = loadMeasureFont(theme.TextFont())
		}

		comp.facesMutex.Lock()
		comp.faces[key] = face
		comp.measureFaces[key] = measureFace
		comp.facesMutex.Unlock()
	}

	return face, measureFace
}

// ClearFontCache is used to remove cached fonts in the case that we wish to re-load font faces
func ClearFontCache() {
	fontCache.Range(func(_, val interface{}) bool {
		item := val.(*fontCacheItem)
		for _, face := range item.faces {
			if face == nil {
				continue
			}
			err := face.Close()

			if err != nil {
				fyne.LogError("failed to close font face", err)
				return false
			}
		}
		return true
	})

	fontCache = &sync.Map{}
}

// DrawString draws a string into an image.
func DrawString(dst draw.Image, s string, color color.Color, f font.Face, face gotext.Face, fontSize, scale float32,
	height int, tabWidth int) {
	src := &image.Uniform{C: color}
	dot := freetype.Pt(0, height-f.Metrics().Descent.Ceil())
	walkString(face, s, float32ToFixed266(fontSize), tabWidth, &dot.X, scale, func(g gotext.GID) {
		dr, mask, maskp, _, ok := f.(truetype.IndexableFace).GlyphAtIndex(dot, truetype.Index(g))
		if !ok {
			dr, mask, maskp, _, ok = f.Glyph(dot, 0xfffd)
		}
		if ok {
			draw.DrawMask(dst, dr, src, image.Point{}, mask, maskp, draw.Over)
		}
	})
}

func loadMeasureFont(data fyne.Resource) gotext.Face {
	loaded, err := gotext.ParseTTF(bytes.NewReader(data.Content()))
	if err != nil {
		fyne.LogError("font load error", err)
	}

	return loaded
}

// MeasureString returns how far dot would advance by drawing s with f.
// Tabs are translated into a dot location change.
func MeasureString(f gotext.Face, s string, textSize float32, tabWidth int) (size fyne.Size, advance fixed.Int26_6) {
	return walkString(f, s, float32ToFixed266(textSize), tabWidth, &advance, 1, func(gotext.GID) {})
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

func loadFont(data fyne.Resource) *truetype.Font {
	loaded, err := truetype.Parse(data.Content())
	if err != nil {
		fyne.LogError("font load error", err)
	}

	return loaded
}

func measureText(text string, fontSize float32, style fyne.TextStyle) (fyne.Size, float32) {
	_, face := CachedFontFace(style, fontSize, 1)
	size, base := MeasureString(face, text, fontSize, style.TabWidth)
	return size, fixed266ToFloat32(base)
}

func newFontWithFallback(chosen, fallback font.Face, chosenFont, fallbackFont ttfFont) font.Face {
	return &compositeFace{chosen: chosen, fallback: fallback, chosenFont: chosenFont, fallbackFont: fallbackFont}
}

func tabStop(spacew, x fixed.Int26_6, tabWidth int) fixed.Int26_6 {
	if tabWidth <= 0 {
		tabWidth = DefaultTabWidth
	}

	tabw := spacew * fixed.Int26_6(tabWidth)
	tabs, _ := math.Modf(float64((x + tabw) / tabw))
	return tabw * fixed.Int26_6(tabs)
}

func walkString(f gotext.Face, s string, textSize fixed.Int26_6, tabWidth int, advance *fixed.Int26_6, scale float32, cb func(g gotext.GID)) (size fyne.Size, base fixed.Int26_6) {
	runes := []rune(s)
	in := shaping.Input{
		Text:      []rune{' '},
		RunStart:  0,
		RunEnd:    1,
		Direction: di.DirectionLTR,
		Face:      f,
		Size:      textSize,
	}
	shaper := &shaping.HarfbuzzShaper{}
	out := shaper.Shape(in)
	spacew := float32ToFixed266(scale) * fontTabSpaceSize

	in.Text = runes
	in.RunStart = 0
	in.RunEnd = len(runes)

	ins := shaping.SplitByFontGlyphs(in, []gotext.Face{f}) // TODO provide fallback...
	for _, in := range ins {
		out = shaper.Shape(in)

		var c rune
		nextRuneIndex := 0
		last := -1
		for _, g := range out.Glyphs {
			if g.ClusterIndex != last {
				c = in.Text[nextRuneIndex]
				nextRuneIndex += g.RuneCount
				last = g.ClusterIndex
			}

			if c == '\r' {
				continue
			}
			if c == '\t' {
				*advance = tabStop(spacew, *advance, tabWidth)
			} else {
				cb(g.GlyphID)
				*advance += float32ToFixed266(fixed266ToFloat32(g.XAdvance) * scale)
			}
		}
	}

	return fyne.NewSize(fixed266ToFloat32(*advance), fixed266ToFloat32(out.LineBounds.LineHeight())),
		out.LineBounds.Ascent
}

var _ truetype.IndexableFace = (*compositeFace)(nil)

type compositeFace struct {
	sync.Mutex

	chosen, fallback         font.Face
	chosenFont, fallbackFont ttfFont
}

func (c *compositeFace) Close() (err error) {
	c.Lock()
	defer c.Unlock()

	if c.chosen != nil {
		err = c.chosen.Close()
	}

	err2 := c.fallback.Close()
	if err2 != nil {
		return err2
	}

	return
}

func (c *compositeFace) Glyph(dot fixed.Point26_6, r rune) (
	dr image.Rectangle, mask image.Image, maskp image.Point, advance fixed.Int26_6, ok bool) {
	c.Lock()
	defer c.Unlock()

	if c.containsGlyph(c.chosenFont, r) {
		return c.chosen.Glyph(dot, r)
	}

	if c.containsGlyph(c.fallbackFont, r) {
		return c.fallback.Glyph(dot, r)
	}

	return
}

func (c *compositeFace) GlyphAdvance(r rune) (advance fixed.Int26_6, ok bool) {
	c.Lock()
	defer c.Unlock()

	if c.containsGlyph(c.chosenFont, r) {
		return c.chosen.GlyphAdvance(r)
	}

	if c.containsGlyph(c.fallbackFont, r) {
		return c.fallback.GlyphAdvance(r)
	}

	return
}

func (c *compositeFace) GlyphAtIndex(dot fixed.Point26_6, g truetype.Index) (dr image.Rectangle, mask image.Image, maskp image.Point,
	advance fixed.Int26_6, ok bool) {
	if g == 0 {
		return image.Rectangle{}, nil, image.Point{}, 0, false
	}

	c.Lock()
	defer c.Unlock()

	dr, mask, maskp, advance, ok = c.chosen.(truetype.IndexableFace).GlyphAtIndex(dot, g)
	if ok {
		return
	}

	return c.fallback.(truetype.IndexableFace).GlyphAtIndex(dot, g)
}

func (c *compositeFace) GlyphBounds(r rune) (bounds fixed.Rectangle26_6, advance fixed.Int26_6, ok bool) {
	c.Lock()
	defer c.Unlock()

	if c.containsGlyph(c.chosenFont, r) {
		return c.chosen.GlyphBounds(r)
	}

	if c.containsGlyph(c.fallbackFont, r) {
		return c.fallback.GlyphBounds(r)
	}

	return
}

func (c *compositeFace) Kern(r0, r1 rune) fixed.Int26_6 {
	c.Lock()
	defer c.Unlock()

	if c.containsGlyph(c.chosenFont, r0) && c.containsGlyph(c.chosenFont, r1) {
		return c.chosen.Kern(r0, r1)
	}

	return c.fallback.Kern(r0, r1)
}

func (c *compositeFace) Metrics() font.Metrics {
	c.Lock()
	defer c.Unlock()

	return c.chosen.Metrics()
}

func (c *compositeFace) containsGlyph(font ttfFont, r rune) bool {
	return font != nil && font.Index(r) != 0
}

type ttfFont interface {
	Index(rune) truetype.Index
}

type faceCacheKey struct {
	size, scale fixed.Int26_6
}

type fontCacheItem struct {
	font, fallback *truetype.Font
	faces          map[faceCacheKey]font.Face
	measureFaces   map[faceCacheKey]gotext.Face
	facesMutex     sync.RWMutex
}

var fontCache = &sync.Map{} // map[fyne.TextStyle]*fontCacheItem
