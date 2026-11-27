package painter

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/internal/cache"
	"fyne.io/fyne/v2/theme"
)

type scopeFontTheme struct {
	fyne.Theme
	font fyne.Resource
}

func (f scopeFontTheme) Font(fyne.TextStyle) fyne.Resource { return f.font }

// newResourceTheme builds a new resource on every Font call, as a theme loading its font on demand does.
type newResourceTheme struct {
	fyne.Theme
	calls *int
}

func (n newResourceTheme) Font(fyne.TextStyle) fyne.Resource {
	*n.calls++
	return fyne.NewStaticResource("loaded.ttf", theme.DefaultTextFont().Content())
}

func TestCachedFontFace_OverrideScopesShareFaces(t *testing.T) {
	ClearFontCache()
	CachedFontFace(fyne.TextStyle{}, nil, canvas.NewText("app", nil))
	parsed := mapLength(&parsedFonts)

	a, b := canvas.NewText("a", nil), canvas.NewText("b", nil)
	cache.OverrideTheme(a, scopeFontTheme{theme.DefaultTheme(), theme.DefaultTextFont()})
	cache.OverrideTheme(b, scopeFontTheme{theme.DefaultTheme(), theme.DefaultTextFont()})
	CachedFontFace(fyne.TextStyle{}, nil, a)
	CachedFontFace(fyne.TextStyle{}, nil, b)

	assert.Equal(t, parsed, mapLength(&parsedFonts))
}

func TestCachedFontFace_RefreshedOverrideFreesOldFaces(t *testing.T) {
	ClearFontCache()
	text := canvas.NewText("text", nil)
	th := scopeFontTheme{theme.DefaultTheme(), theme.DefaultTextFont()}
	scope := cache.OverrideTheme(text, th)
	CachedFontFace(fyne.TextStyle{Bold: true}, nil, text)
	entries := mapLength(&fontCache)

	for range 50 {
		scope = rescope(text, th, scope)
		CachedFontFace(fyne.TextStyle{Bold: true}, nil, text)
	}
	assert.Equal(t, entries, mapLength(&fontCache))
}

func TestCachedFontFace_ThemeWithNewResourcePerCall(t *testing.T) {
	ClearFontCache()
	calls := 0
	text := canvas.NewText("text", nil)
	th := newResourceTheme{theme.DefaultTheme(), &calls}
	scope := cache.OverrideTheme(text, th)

	first := CachedFontFace(fyne.TextStyle{}, nil, text)
	for range 100 {
		assert.Same(t, first, CachedFontFace(fyne.TextStyle{}, nil, text))
	}
	assert.Equal(t, 1, calls)

	entries, parsed := mapLength(&fontCache), mapLength(&parsedFonts)
	for range 50 {
		scope = rescope(text, th, scope)
		CachedFontFace(fyne.TextStyle{}, nil, text)
	}
	assert.Equal(t, entries, mapLength(&fontCache))
	assert.Equal(t, parsed, mapLength(&parsedFonts))
}

// rescope does what ThemeOverride does on each Refresh.
func rescope(o fyne.CanvasObject, th fyne.Theme, old string) string {
	scope := cache.OverrideTheme(o, th)
	ReleaseFontScope(old)
	return scope
}

func mapLength[K comparable, V any](m interface{ Range(func(K, V) bool) }) (n int) {
	m.Range(func(K, V) bool { n++; return true })
	return n
}
