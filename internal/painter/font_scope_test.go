package painter_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/internal/cache"
	"fyne.io/fyne/v2/internal/painter"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"
)

type fontTheme struct {
	fyne.Theme
	font fyne.Resource
}

func (f fontTheme) Font(fyne.TextStyle) fyne.Resource { return f.font }

func TestCachedFontFace_OverrideScopesShareFaces(t *testing.T) {
	painter.ClearFontCache()
	painter.CachedFontFace(fyne.TextStyle{}, nil, canvas.NewText("app", nil))
	parsed := painter.ParsedFontsLen()

	a, b := canvas.NewText("a", nil), canvas.NewText("b", nil)
	cache.OverrideTheme(a, fontTheme{theme.DefaultTheme(), theme.DefaultTextFont()})
	cache.OverrideTheme(b, fontTheme{theme.DefaultTheme(), theme.DefaultTextFont()})
	painter.CachedFontFace(fyne.TextStyle{}, nil, a)
	painter.CachedFontFace(fyne.TextStyle{}, nil, b)

	assert.Equal(t, parsed, painter.ParsedFontsLen())
}

func TestCachedFontFace_OverrideScopeWithOtherFont(t *testing.T) {
	painter.ClearFontCache()
	app := painter.CachedFontFace(fyne.TextStyle{}, nil, nil)

	mono := canvas.NewText("mono", nil)
	cache.OverrideTheme(mono, fontTheme{theme.DefaultTheme(), theme.DefaultTextMonospaceFont()})
	assert.NotSame(t, app, painter.CachedFontFace(fyne.TextStyle{}, nil, mono))
}

func TestCachedFontFace_RefreshedOverrideFreesOldFaces(t *testing.T) {
	test.NewTempApp(t)
	painter.ClearFontCache()
	text := canvas.NewText("text", nil)
	override := container.NewThemeOverride(text, fontTheme{theme.DefaultTheme(), theme.DefaultTextFont()})
	painter.CachedFontFace(fyne.TextStyle{Bold: true}, nil, text)
	entries := painter.FontCacheLen()

	for range 50 {
		override.Refresh()
		painter.CachedFontFace(fyne.TextStyle{Bold: true}, nil, text)
	}
	assert.Equal(t, entries, painter.FontCacheLen())
}

// newResourceTheme builds a new resource on every Font call, as a theme loading its font on demand does.
type newResourceTheme struct {
	fyne.Theme
	calls *int
}

func (n newResourceTheme) Font(fyne.TextStyle) fyne.Resource {
	*n.calls++
	return fyne.NewStaticResource("loaded.ttf", theme.DefaultTextFont().Content())
}

func TestCachedFontFace_ThemeWithNewResourcePerCall(t *testing.T) {
	test.NewTempApp(t)
	painter.ClearFontCache()
	calls := 0
	text := canvas.NewText("text", nil)
	override := container.NewThemeOverride(text, newResourceTheme{theme.DefaultTheme(), &calls})

	first := painter.CachedFontFace(fyne.TextStyle{}, nil, text)
	for range 100 {
		assert.Same(t, first, painter.CachedFontFace(fyne.TextStyle{}, nil, text))
	}
	assert.Equal(t, 1, calls)

	entries, parsed := painter.FontCacheLen(), painter.ParsedFontsLen()
	for range 50 {
		override.Refresh()
		painter.CachedFontFace(fyne.TextStyle{}, nil, text)
	}
	assert.Equal(t, entries, painter.FontCacheLen())
	assert.Equal(t, parsed, painter.ParsedFontsLen())
}

// A theme may return a Resource whose concrete type cannot be a map key.
type sliceResource struct{ content []byte }

func (sliceResource) Name() string      { return "slice.ttf" }
func (r sliceResource) Content() []byte { return r.content }

func TestCachedFontFace_NonComparableFontKeepsScope(t *testing.T) {
	painter.ClearFontCache()
	res := sliceResource{theme.DefaultTextFont().Content()}
	a, b := canvas.NewText("a", nil), canvas.NewText("b", nil)
	cache.OverrideTheme(a, fontTheme{theme.DefaultTheme(), res})
	cache.OverrideTheme(b, fontTheme{theme.DefaultTheme(), res})

	fa := painter.CachedFontFace(fyne.TextStyle{}, nil, a)
	assert.NotNil(t, fa)
	assert.Same(t, fa, painter.CachedFontFace(fyne.TextStyle{}, nil, a))
	assert.NotSame(t, fa, painter.CachedFontFace(fyne.TextStyle{}, nil, b))
}

// A struct type is comparable even when the value held in its interface field is not.
type wrappedResource struct{ fyne.Resource }

func TestCachedFontFace_WrappedNonComparableFontKeepsScope(t *testing.T) {
	painter.ClearFontCache()
	res := wrappedResource{sliceResource{theme.DefaultTextFont().Content()}}
	a, b := canvas.NewText("a", nil), canvas.NewText("b", nil)
	cache.OverrideTheme(a, fontTheme{theme.DefaultTheme(), res})
	cache.OverrideTheme(b, fontTheme{theme.DefaultTheme(), res})

	fa := painter.CachedFontFace(fyne.TextStyle{}, nil, a)
	assert.NotNil(t, fa)
	assert.Same(t, fa, painter.CachedFontFace(fyne.TextStyle{}, nil, a))
	assert.NotSame(t, fa, painter.CachedFontFace(fyne.TextStyle{}, nil, b))
}
