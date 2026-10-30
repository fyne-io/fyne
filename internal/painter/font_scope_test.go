package painter_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/internal/cache"
	"fyne.io/fyne/v2/internal/painter"
	"fyne.io/fyne/v2/theme"
)

type fontTheme struct {
	fyne.Theme
	font fyne.Resource
}

func (f fontTheme) Font(fyne.TextStyle) fyne.Resource { return f.font }

func TestCachedFontFace_OverrideScopesShareFaces(t *testing.T) {
	painter.ClearFontCache()
	app := painter.CachedFontFace(fyne.TextStyle{}, nil, canvas.NewText("app", nil))

	a, b := canvas.NewText("a", nil), canvas.NewText("b", nil)
	cache.OverrideTheme(a, fontTheme{theme.DefaultTheme(), theme.DefaultTextFont()})
	cache.OverrideTheme(b, fontTheme{theme.DefaultTheme(), theme.DefaultTextFont()})
	assert.NotEqual(t, cache.WidgetScopeID(a), cache.WidgetScopeID(b))

	assert.Same(t, painter.CachedFontFace(fyne.TextStyle{}, nil, a), painter.CachedFontFace(fyne.TextStyle{}, nil, b))
	assert.Same(t, app, painter.CachedFontFace(fyne.TextStyle{}, nil, a))
}

func TestCachedFontFace_OverrideScopeWithOtherFont(t *testing.T) {
	painter.ClearFontCache()
	app := painter.CachedFontFace(fyne.TextStyle{}, nil, nil)

	mono := canvas.NewText("mono", nil)
	cache.OverrideTheme(mono, fontTheme{theme.DefaultTheme(), theme.DefaultTextMonospaceFont()})
	assert.NotSame(t, app, painter.CachedFontFace(fyne.TextStyle{}, nil, mono))
}

func TestCachedFontFace_RefreshedOverrideKeepsFaces(t *testing.T) {
	painter.ClearFontCache()
	text := canvas.NewText("text", nil)
	th := fontTheme{theme.DefaultTheme(), theme.DefaultTextFont()}
	cache.OverrideTheme(text, th)
	before := painter.CachedFontFace(fyne.TextStyle{Bold: true}, nil, text)

	cache.OverrideTheme(text, th) // as ThemeOverride.Refresh does
	assert.Same(t, before, painter.CachedFontFace(fyne.TextStyle{Bold: true}, nil, text))
}

// A theme may return a Resource whose concrete type cannot be a map key.
type sliceResource struct{ content []byte }

func (r sliceResource) Name() string    { return "slice.ttf" }
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
