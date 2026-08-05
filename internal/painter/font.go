package painter

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"math"
	"strings"
	"sync"

	"github.com/go-text/render"
	"github.com/go-text/typesetting/di"
	"github.com/go-text/typesetting/font"
	"github.com/go-text/typesetting/fontscan"
	"github.com/go-text/typesetting/language"
	"github.com/go-text/typesetting/shaping"
	"golang.org/x/image/math/fixed"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/async"
	"fyne.io/fyne/v2/internal/cache"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/theme"
)

const (
	// DefaultTabWidth is the default width in spaces
	DefaultTabWidth               = 4
	UnderlineOffsetFromBaseline   = 2
	StrikethroughToBaselineFactor = 0.75

	fontTabSpaceSize = 10
	replacementChar  = 0xfffd // that’s '�'
)

var (
	fm           *fontscan.FontMap
	fontScanLock sync.Mutex
	loaded       bool
)

var shaper = &shaping.HarfbuzzShaper{}

func loadMap() {
	loaded = true

	fm = fontscan.NewFontMap(noopLogger{})
	err := loadSystemFonts(fm)
	if err != nil {
		fm = nil // just don't fallback
	}
}

func lookupLangFont(family string, aspect font.Aspect) *font.Face {
	fontScanLock.Lock()
	defer fontScanLock.Unlock()

	if !loaded {
		loadMap()
	}
	if fm == nil {
		return nil
	}

	fm.SetQuery(fontscan.Query{Families: []string{family}, Aspect: aspect})
	l, _ := language.NewLangID(language.NewLanguage(lang.SystemLocale().LanguageString()))
	return fm.ResolveFaceForLang(l)
}

func lookupRuneFont(r rune, family string, aspect font.Aspect) *font.Face {
	fontScanLock.Lock()
	defer fontScanLock.Unlock()

	if !loaded {
		loadMap()
	}
	if fm == nil {
		return nil
	}

	fm.SetQuery(fontscan.Query{Families: []string{family}, Aspect: aspect})
	fm.SetScript(language.LookupScript(r))
	return fm.ResolveFace(r)
}

func lookupFaces(t, fallback fyne.Resource, additional []fyne.Resource, family string, style fyne.TextStyle) (faces *dynamicFontMap) {
	f1 := loadMeasureFont(t)
	if t == fallback {
		faces = &dynamicFontMap{family: family, faces: []*font.Face{f1}}
	} else {
		f2 := loadMeasureFont(fallback)
		faces = &dynamicFontMap{family: family, faces: []*font.Face{f1, f2}}
	}

	aspect := font.Aspect{Style: font.StyleNormal}
	if style.Italic {
		aspect.Style = font.StyleItalic
	}
	if style.Bold {
		aspect.Weight = font.WeightBold
	}

	for _, added := range additional {
		faces.addFace(loadMeasureFont(added))
	}

	local := lookupLangFont(family, aspect)
	if local != nil {
		faces.addFace(local)
	}

	return faces
}

// CachedFontFace returns a Font face held in memory. These are loaded from the current theme.
func CachedFontFace(style fyne.TextStyle, source fyne.Resource, o fyne.CanvasObject) *FontCacheItem {
	if source != nil {
		val, ok := fontCustomCache.Load(source)
		if !ok {
			face := loadMeasureFont(source)
			if face == nil {
				face = loadMeasureFont(theme.TextFont())
			}
			faces := &dynamicFontMap{family: source.Name(), faces: []*font.Face{face}}

			val = &FontCacheItem{Fonts: faces}
			fontCustomCache.Store(source, val)
		}
		return val
	}

	scope := ""
	if o != nil { // for overridden themes get the cache key right
		scope = cache.WidgetScopeID(o)
	}

	val, ok := fontCache.Load(cacheID{style: style, scope: scope})
	if !ok {
		var faces *dynamicFontMap

		th := theme.CurrentForWidget(o)
		font1 := th.Font(style)

		// Skip any nil fallback fonts — they can be nil when built with
		// -tags no_emoji, and the lookupFaces loop expects non-nil entries.
		var fallbacks []fyne.Resource
		if emoji := theme.DefaultEmojiFont(); emoji != nil { // TODO only one emoji - maybe others too
			fallbacks = append(fallbacks, emoji)
		}
		if sym := theme.DefaultSymbolFont(); sym != nil {
			fallbacks = append(fallbacks, sym)
		}
		switch {
		case style.Monospace:
			faces = lookupFaces(font1, theme.DefaultTextMonospaceFont(), fallbacks, fontscan.Monospace, style)
		case style.Bold:
			if style.Italic {
				faces = lookupFaces(font1, theme.DefaultTextBoldItalicFont(), fallbacks, fontscan.SansSerif, style)
			} else {
				faces = lookupFaces(font1, theme.DefaultTextBoldFont(), fallbacks, fontscan.SansSerif, style)
			}
		case style.Italic:
			faces = lookupFaces(font1, theme.DefaultTextItalicFont(), fallbacks, fontscan.SansSerif, style)
		case style.Symbol:
			th := theme.SymbolFont()
			fallback := theme.DefaultSymbolFont()
			f1 := loadMeasureFont(th)

			if th == fallback {
				faces = &dynamicFontMap{family: fontscan.SansSerif, faces: []*font.Face{f1}}
			} else {
				f2 := loadMeasureFont(fallback)
				faces = &dynamicFontMap{family: fontscan.SansSerif, faces: []*font.Face{f1, f2}}
			}
		default:
			faces = lookupFaces(font1, theme.DefaultTextFont(), fallbacks, fontscan.SansSerif, style)
		}

		val = &FontCacheItem{Fonts: faces}
		fontCache.Store(cacheID{style: style, scope: scope}, val)
	}

	return val
}

// ClearFontCache is used to remove cached fonts in the case that we wish to re-load Font faces
func ClearFontCache() {
	fontCache.Clear()
	fontCustomCache.Clear()
}

// DrawString draws a string into an image.
func DrawString(dst draw.Image, s string, c color.Color, f shaping.Fontmap, fontSize, scale float32, style fyne.TextStyle) {
	DrawStringOffset(dst, s, c, f, fontSize, scale, style, 0)
}

// DrawStringOffset draws a string shifted left by the specified pixel offset.
func DrawStringOffset(dst draw.Image, s string, c color.Color, f shaping.Fontmap, fontSize, scale float32, style fyne.TextStyle, offset int) {
	r := render.Renderer{
		FontSize: fontSize,
		PixScale: scale,
		Color:    c,
	}

	advance := float32(0)
	walkString(f, s, float32ToFixed266(fontSize), style, &advance, scale, func(run shaping.Output, x, y float32) {
		yPix := int(math.Ceil(float64(y)))
		if len(run.Glyphs) == 1 && run.Glyphs[0].GlyphID == 0 {
			r.DrawStringAt(string([]rune{replacementChar}), dst, int(x)-offset, yPix, f.ResolveFace(replacementChar))
			return
		}

		r.DrawShapedRunAt(run, dst, int(x)-offset, yPix)
	})
}

// WalkStringGlyphs calls cb once for each glyph in s. The callback receives
// the shaped run containing the glyph, the glyph's index within run.Glyphs,
// the accumulated pen X position in device pixels at the start of that glyph,
// the shared baseline Y for the line in device pixels, and the HarfBuzz X/Y
// positioning offsets (also in device pixels).
//
// baseY is the same value for every run in the string, taken from the largest
// ascent across the runs, so glyphs from differently sized faces sit on one
// baseline rather than each on their own.
//
// Runs that consist solely of a replacement-char glyph (GlyphID==0) are skipped
// so callers do not need to handle them.
func WalkStringGlyphs(f shaping.Fontmap, s string, fontSize float32, style fyne.TextStyle, scale float32,
	cb func(run shaping.Output, idx int, penX, baseY, xOff, yOff float32),
) {
	advance := float32(0)
	walkString(f, s, float32ToFixed266(fontSize), style, &advance, scale, func(run shaping.Output, x, y float32) {
		if len(run.Glyphs) == 1 && run.Glyphs[0].GlyphID == 0 {
			return
		}
		penX := x
		for i, g := range run.Glyphs {
			xOff := float32(math.Round(float64(fixed266ToFloat32(g.XOffset) * scale)))
			yOff := float32(math.Round(float64(fixed266ToFloat32(g.YOffset) * scale)))
			cb(run, i, penX, y, xOff, yOff)
			penX += fixed266ToFloat32(g.Advance) * scale
		}
	})
}

// RenderGlyphToImage rasterises a single glyph from run (at index idx) into a
// freshly allocated RGBA image. XOffset and YOffset of the glyph are zeroed so
// the result is independent of kerning context; callers apply the real offsets
// when positioning the glyph on screen.
//
// Image height equals the full line height (ascent + |descent|). The returned
// baseline is the glyph's baseline measured in pixels down from the top of the
// image, which callers need in order to sit the bitmap on the line's shared
// baseline: a run's own ascent is not necessarily the line's ascent when faces
// of different sizes are mixed.
func RenderGlyphToImage(run shaping.Output, idx int, fontSize, scale float32, col color.Color) (img *image.RGBA, baseline int) {
	g := run.Glyphs[idx]
	ren := &render.Renderer{FontSize: fontSize, PixScale: scale, Color: col}

	baseline = int(math.Ceil(float64(fixed266ToFloat32(run.LineBounds.Ascent) * scale)))
	descent := int(math.Ceil(float64(-fixed266ToFloat32(run.LineBounds.Descent) * scale)))
	h := baseline + descent
	if h <= 0 {
		h = 1
	}
	italicPad := int(math.Ceil(float64(fontSize * scale / 5)))
	w := int(math.Ceil(float64(fixed266ToFloat32(g.Advance)*scale))) + italicPad + 2
	if w <= 0 {
		w = 1
	}

	img = image.NewRGBA(image.Rect(0, 0, w, h))
	noOffset := g
	noOffset.XOffset = 0
	noOffset.YOffset = 0
	singleRun := run
	singleRun.Glyphs = []shaping.Glyph{noOffset}
	ren.DrawShapedRunAt(singleRun, img, 0, baseline)
	return img, baseline
}

func loadMeasureFont(data fyne.Resource) *font.Face {
	loaded, err := font.ParseTTF(bytes.NewReader(data.Content()))
	if err != nil {
		fyne.LogError("font load error", err)
		return nil
	}

	return loaded
}

// MeasureString returns how far dot would advance by drawing s with f.
// Tabs are translated into a dot location change.
func MeasureString(f shaping.Fontmap, s string, textSize float32, style fyne.TextStyle) (size fyne.Size, advance float32) {
	return walkString(f, s, float32ToFixed266(textSize), style, &advance, 1, func(shaping.Output, float32, float32) {})
}

// RenderedTextSize looks up how big a string would be if drawn on screen.
// It also returns the distance from top to the text baseline.
func RenderedTextSize(text string, fontSize float32, style fyne.TextStyle, source fyne.Resource) (size fyne.Size, baseline float32) {
	size, base := cache.GetFontMetrics(text, fontSize, style, source)
	if base != 0 {
		return size, base
	}

	size, base = measureText(text, fontSize, style, source)
	cache.SetFontMetrics(text, fontSize, style, source, size, base)
	return size, base
}

func fixed266ToFloat32(i fixed.Int26_6) float32 {
	return float32(float64(i) / (1 << 6))
}

func float32ToFixed266(f float32) fixed.Int26_6 {
	return fixed.Int26_6(float64(f) * (1 << 6))
}

func measureText(text string, fontSize float32, style fyne.TextStyle, source fyne.Resource) (fyne.Size, float32) {
	face := CachedFontFace(style, source, nil)
	return MeasureString(face.Fonts, text, fontSize, style)
}

func tabStop(spacew, x float32, tabWidth int) float32 {
	if tabWidth <= 0 {
		tabWidth = DefaultTabWidth
	}

	tabw := spacew * float32(tabWidth)
	tabs, _ := math.Modf(float64((x + tabw) / tabw))
	return tabw * float32(tabs)
}

type shapedRun struct {
	out shaping.Output
	x   float32
}

var (
	runBuffer    []shapedRun
	runBufferMut async.Mutex
)

// walkString shapes s and invokes cb once per shaped run, in left-to-right order.
// All runs share a single ascent (the max ascent of any run in the string), so that
// runs shaped in different fallback fonts (e.g. mixed-script or emoji + text)
// still align on a common baseline
func walkString(faces shaping.Fontmap, s string, textSize fixed.Int26_6, style fyne.TextStyle, advance *float32, scale float32,
	cb func(run shaping.Output, x, y float32),
) (size fyne.Size, base float32) {
	s = strings.ReplaceAll(s, "\r", "")

	runes := []rune(s)
	in := shaping.Input{
		Text:      []rune{' '},
		RunStart:  0,
		RunEnd:    1,
		Direction: di.DirectionLTR,
		Face:      faces.ResolveFace(' '),
		Size:      textSize,
	}
	segmenter := &shaping.Segmenter{}
	out := shaper.Shape(in)

	in.Text = runes
	in.RunStart = 0
	in.RunEnd = len(runes)

	x := float32(0)
	spacew := scale * fontTabSpaceSize
	if style.Monospace {
		spacew = scale * fixed266ToFloat32(out.Advance)
	}

	maxAscent := fixed.Int26_6(0)
	collect := func(run shaping.Output, runX float32) {
		if run.LineBounds.Ascent > maxAscent {
			maxAscent = run.LineBounds.Ascent
		}
		runBuffer = append(runBuffer, shapedRun{out: run, x: runX})
	}

	ins := segmenter.Split(in, faces)
	runBufferMut.Lock()
	for _, in := range ins {
		inEnd := in.RunEnd

		pending := false
		for i, r := range in.Text[in.RunStart:in.RunEnd] {
			if r == '\t' {
				if pending {
					in.RunEnd = i
					x = shapeCallback(in, x, scale, collect)
				}
				x = tabStop(spacew, x, style.TabWidth)

				in.RunStart = i + 1
				in.RunEnd = inEnd
				pending = false
			} else {
				pending = true
			}
		}

		x = shapeCallback(in, x, scale, collect)
	}

	y := fixed266ToFloat32(maxAscent) * scale
	for _, run := range runBuffer {
		cb(run.out, run.x, y)
	}
	clear(runBuffer)
	runBuffer = runBuffer[:0]
	runBufferMut.Unlock()

	*advance = x
	return fyne.NewSize(*advance, fixed266ToFloat32(out.LineBounds.LineThickness())),
		fixed266ToFloat32(out.LineBounds.Ascent)
}

func shapeCallback(in shaping.Input, x, scale float32, cb func(shaping.Output, float32)) float32 {
	out := shaper.Shape(in)
	glyphs := out.Glyphs
	start := 0
	adv := fixed.I(0)
	for i, g := range out.Glyphs {
		if g.GlyphID == 0 {
			if start < i {
				out.Glyphs = glyphs[start:i]
				cb(out, x)
				x += fixed266ToFloat32(adv) * scale
			}

			out.Glyphs = glyphs[i : i+1]
			cb(out, x)
			x += fixed266ToFloat32(glyphs[i].Advance) * scale

			adv = 0
			start = i + 1
		}
		adv += g.Advance
	}

	if start < len(glyphs) {
		out.Glyphs = glyphs[start:]
		cb(out, x)
		x += fixed266ToFloat32(adv) * scale
		adv = 0
	}
	return x + fixed266ToFloat32(adv)*scale
}

type FontCacheItem struct {
	Fonts shaping.Fontmap
}

type cacheID struct {
	style fyne.TextStyle
	scope string
}

var (
	fontCache       async.Map[cacheID, *FontCacheItem]
	fontCustomCache async.Map[fyne.Resource, *FontCacheItem] // for custom resources
)

type noopLogger struct{}

func (noopLogger) Printf(string, ...any) {}

type dynamicFontMap struct {
	faces  []*font.Face
	family string
}

func (d *dynamicFontMap) ResolveFace(r rune) *font.Face {
	for _, f := range d.faces {
		if _, ok := f.NominalGlyph(r); ok {
			return f
		}
	}

	toAdd := lookupRuneFont(r, d.family, font.Aspect{})
	if toAdd != nil {
		d.addFace(toAdd)
		return toAdd
	}

	return d.faces[0]
}

func (d *dynamicFontMap) addFace(f *font.Face) {
	d.faces = append(d.faces, f)
}
