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
func useHanFonts(t *testing.T) {
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

	for _, name := range []string{"han_partial.ttf", "han_full.ttf"} {
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
