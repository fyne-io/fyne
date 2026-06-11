package app

import (
	"os"
	"path/filepath"
	"testing"

	"fyne.io/fyne/v2"
	internalapp "fyne.io/fyne/v2/internal/app"
	"fyne.io/fyne/v2/internal/build"
	internalTest "fyne.io/fyne/v2/internal/test"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSettingsBuildType(t *testing.T) {
	set := test.NewApp().Settings()
	assert.Equal(t, fyne.BuildStandard, set.BuildType()) // during test we should have a normal build

	set = &settings{}
	assert.Equal(t, build.Mode, set.BuildType()) // when testing this package only it could be debug or release
}

func TestSettingsLoad(t *testing.T) {
	settings := &settings{}

	require.NoError(t, settings.loadFromFile(filepath.Join("testdata", "light-theme.json")))
	assert.Equal(t, "light", settings.schema.ThemeName)

	require.NoError(t, settings.loadFromFile(filepath.Join("testdata", "dark-theme.json")))
	assert.Equal(t, "dark", settings.schema.ThemeName)
}

func TestOverrideTheme(t *testing.T) {
	require.NoError(t, os.Setenv("FYNE_THEME", ""))
	set := &settings{}
	set.setupTheme()
	assert.Equal(t, internalapp.DefaultVariant(), set.ThemeVariant())

	set.schema.ThemeName = "light"
	set.setupTheme()
	assert.Equal(t, theme.DefaultTheme(), set.Theme())
	assert.Equal(t, theme.VariantLight, set.ThemeVariant())

	set.schema.ThemeName = "dark"
	set.setupTheme()
	assert.Equal(t, theme.DefaultTheme(), set.Theme())
	assert.Equal(t, theme.VariantDark, set.ThemeVariant())

	set = &settings{}
	set.setupTheme()
	assert.Equal(t, internalapp.DefaultVariant(), set.ThemeVariant())

	require.NoError(t, os.Setenv("FYNE_THEME", "light"))
	set.setupTheme()
	assert.Equal(t, theme.DefaultTheme(), set.Theme())
	assert.Equal(t, theme.VariantLight, set.ThemeVariant())

	require.NoError(t, os.Setenv("FYNE_THEME", ""))
}

func TestOverrideTheme_IgnoresSettingsChange(t *testing.T) {
	// check that a file-load does not overwrite our value
	set := &settings{}
	require.NoError(t, os.Setenv("FYNE_THEME", "light"))
	set.setupTheme()
	assert.Equal(t, theme.DefaultTheme(), set.Theme())
	assert.Equal(t, theme.VariantLight, set.ThemeVariant())

	require.NoError(t, set.loadFromFile(filepath.Join("testdata", "dark-theme.json")))
	set.setupTheme()
	assert.Equal(t, theme.DefaultTheme(), set.Theme())
	assert.Equal(t, theme.VariantLight, set.ThemeVariant())
	require.NoError(t, os.Setenv("FYNE_THEME", ""))
}

func TestCustomTheme(t *testing.T) {
	type customTheme struct {
		fyne.Theme
	}
	set := &settings{}
	ctheme := &customTheme{internalTest.LightTheme(theme.DefaultTheme())}
	set.SetTheme(ctheme)

	set.setupTheme()
	assert.Equal(t, set.Theme(), ctheme)
	assert.Equal(t, internalapp.DefaultVariant(), set.ThemeVariant())

	require.NoError(t, set.loadFromFile(filepath.Join("testdata", "light-theme.json")))
	set.setupTheme()
	assert.Equal(t, set.Theme(), ctheme)
	assert.Equal(t, theme.VariantLight, set.ThemeVariant())

	require.NoError(t, set.loadFromFile(filepath.Join("testdata", "dark-theme.json")))
	set.setupTheme()
	assert.Equal(t, set.Theme(), ctheme)
	assert.Equal(t, theme.VariantDark, set.ThemeVariant())

	require.NoError(t, os.Setenv("FYNE_THEME", "light"))
	set.setupTheme()
	assert.Equal(t, set.Theme(), ctheme)
	assert.Equal(t, theme.VariantLight, set.ThemeVariant())

	require.NoError(t, os.Setenv("FYNE_THEME", "dark"))
	set.setupTheme()
	assert.Equal(t, set.Theme(), ctheme)
	assert.Equal(t, theme.VariantDark, set.ThemeVariant())

	require.NoError(t, os.Setenv("FYNE_THEME", ""))
	set.setupTheme()
	assert.Equal(t, set.Theme(), ctheme)
	assert.Equal(t, theme.VariantDark, set.ThemeVariant())
}
