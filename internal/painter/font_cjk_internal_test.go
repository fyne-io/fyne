package painter

import (
	"os"
	"testing"

	"github.com/go-text/typesetting/font"
	"github.com/go-text/typesetting/fontscan"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/image/math/fixed"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// han_partial.ttf has some Han runes and lacks others that han_full.ttf has,
// the way a Japanese font lacks some simplified Chinese runes. Both were made
// with fontTools, one box glyph for every rune of a format 13 cmap.
// han_other.ttf has only the three that han_partial.ttf lacks and han_full.ttf
// has, so no font of the pair han_partial.ttf and han_other.ttf has them all.
func useHanFonts(t *testing.T, names ...string) {
	t.Helper()

	fontScanLock.Lock()
	oldMap, oldLoaded := fm, loaded
	fm, loaded = fontscan.NewFontMap(noopLogger{}), true
	fontScanLock.Unlock()
	t.Cleanup(func() {
		fontScanLock.Lock()
		fm, loaded = oldMap, oldLoaded
		fontScanLock.Unlock()
	})

	if len(names) == 0 {
		names = []string{"han_partial.ttf", "han_full.ttf"}
	}
	for _, name := range names {
		f, err := os.Open("testdata/font/" + name)
		require.NoError(t, err)
		require.NoError(t, fm.AddFont(f, name, ""))
		require.NoError(t, f.Close())
	}
}

func shapedFamilies(t *testing.T, s string, setUp ...*font.Face) []string {
	t.Helper()

	faces := &dynamicFontMap{family: fontscan.SansSerif, faces: []*font.Face{loadMeasureFont(theme.DefaultTextFont())}, fixed: 1}
	faces.faces = append(faces.faces, setUp...)
	faces.fixed += len(setUp)
	runBufferMut.Lock()
	defer runBufferMut.Unlock()
	shaped := shapeString(faces, s, fixed.I(20))

	var families []string
	for _, step := range shaped.steps {
		families = append(families, step.out.Face.Describe().Family)
	}
	return families
}

func TestShapeString_CJKStaysInOneFont(t *testing.T) {
	useHanFonts(t)

	// 页, 础 and 镜 are not in han_partial.ttf
	assert.Equal(t, []string{"Han Full"}, shapedFamilies(t, "首页基础像镜"))
	// 、 is in both fonts, and has to follow the rest
	assert.Equal(t, []string{"Han Full"}, shapedFamilies(t, "首页、基础镜像"))
	// every rune is in han_partial.ttf, which is found first
	assert.Equal(t, []string{"Han Partial"}, shapedFamilies(t, "首基像"))
	// U+3400 is in neither font, and must not undo the choice for the rest
	assert.Equal(t, []string{"Han Full"}, shapedFamilies(t, "首页\u3400基础"))
}

func TestShapeString_CJKFontOfTheMapIsKept(t *testing.T) {
	useHanFonts(t)

	data, err := os.ReadFile("testdata/font/han_partial.ttf")
	require.NoError(t, err)
	partial := loadMeasureFont(fyne.NewStaticResource("han_partial.ttf", data))

	// the font set up with the map draws 首, so the string is not moved to han_full.ttf
	families := shapedFamilies(t, "首页基础", partial)
	assert.Equal(t, "Han Partial", families[0])
}

func TestShapeString_CJKNoFontHasAll(t *testing.T) {
	useHanFonts(t, "han_partial.ttf", "han_other.ttf")

	// 首 and 基 are only in han_partial.ttf, 页 and 础 only in han_other.ttf, so
	// each rune keeps the font it is found in
	assert.Equal(t, []string{"Han Partial", "Han Other", "Han Partial", "Han Other"}, shapedFamilies(t, "首页基础"))
	// a string that one of them has all of is still drawn by it
	assert.Equal(t, []string{"Han Other"}, shapedFamilies(t, "页础镜"))
}
