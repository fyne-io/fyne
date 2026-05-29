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
)

func TestSettingsBuildType(t *testing.T) {
	set := test.NewApp().Settings()
	assert.Equal(t, fyne.BuildStandard, set.BuildType()) // during test we should have a normal build

	set = &settings{}
	assert.Equal(t, build.Mode, set.BuildType()) // when testing this package only it could be debug or release
}

func TestSettingsLoad(t *testing.T) {
	settings := &settings{}

	err := settings.loadFromFile(filepath.Join("testdata", "light-theme.json"))
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, "light", settings.schema.ThemeName)

	err = settings.loadFromFile(filepath.Join("testdata", "dark-theme.json"))
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, "dark", settings.schema.ThemeName)
}

func TestOverrideTheme(t *testing.T) {
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

	err := os.Setenv("FYNE_THEME", "light")
	if err != nil {
		t.Error(err)
	}
	set.setupTheme()
	assert.Equal(t, theme.DefaultTheme(), set.Theme())
	assert.Equal(t, theme.VariantLight, set.ThemeVariant())
	err = os.Setenv("FYNE_THEME", "")
	if err != nil {
		t.Error(err)
	}
}

func TestOverrideTheme_IgnoresSettingsChange(t *testing.T) {
	// check that a file-load does not overwrite our value
	set := &settings{}
	err := os.Setenv("FYNE_THEME", "light")
	if err != nil {
		t.Error(err)
	}
	set.setupTheme()
	assert.Equal(t, theme.DefaultTheme(), set.Theme())
	assert.Equal(t, theme.VariantLight, set.ThemeVariant())

	err = set.loadFromFile(filepath.Join("testdata", "dark-theme.json"))
	if err != nil {
		t.Error(err)
	}
	set.setupTheme()
	assert.Equal(t, theme.DefaultTheme(), set.Theme())
	assert.Equal(t, theme.VariantLight, set.ThemeVariant())
	err = os.Setenv("FYNE_THEME", "")
	if err != nil {
		t.Error(err)
	}
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

	err := set.loadFromFile(filepath.Join("testdata", "light-theme.json"))
	if err != nil {
		t.Error(err)
	}
	set.setupTheme()
	assert.Equal(t, set.Theme(), ctheme)
	assert.Equal(t, theme.VariantLight, set.ThemeVariant())

	err = set.loadFromFile(filepath.Join("testdata", "dark-theme.json"))
	if err != nil {
		t.Error(err)
	}
	set.setupTheme()
	assert.Equal(t, set.Theme(), ctheme)
	assert.Equal(t, theme.VariantDark, set.ThemeVariant())

	err = os.Setenv("FYNE_THEME", "light")
	if err != nil {
		t.Error(err)
	}
	set.setupTheme()
	assert.Equal(t, set.Theme(), ctheme)
	assert.Equal(t, theme.VariantLight, set.ThemeVariant())

	err = os.Setenv("FYNE_THEME", "dark")
	if err != nil {
		t.Error(err)
	}
	set.setupTheme()
	assert.Equal(t, set.Theme(), ctheme)
	assert.Equal(t, theme.VariantDark, set.ThemeVariant())

	err = os.Setenv("FYNE_THEME", "")
	if err != nil {
		t.Error(err)
	}
	set.setupTheme()
	assert.Equal(t, set.Theme(), ctheme)
	assert.Equal(t, theme.VariantDark, set.ThemeVariant())
}

func TestSetThemeVariant_Dark(t *testing.T) {
	set := &settings{}
	set.setupTheme()

	set.SetThemeVariant(theme.VariantDark)
	assert.Equal(t, theme.VariantDark, set.ThemeVariant())
	assert.True(t, set.variantSpecified)
}

func TestSetThemeVariant_Light(t *testing.T) {
	set := &settings{}
	set.setupTheme()

	set.SetThemeVariant(theme.VariantLight)
	assert.Equal(t, theme.VariantLight, set.ThemeVariant())
	assert.True(t, set.variantSpecified)
}

func TestSetThemeVariant_System(t *testing.T) {
	set := &settings{}
	set.setupTheme()

	// First set to a specific variant
	set.SetThemeVariant(theme.VariantDark)
	assert.True(t, set.variantSpecified)

	// Then switch to system variant
	set.SetThemeVariant(theme.VariantSystem)
	assert.False(t, set.variantSpecified)
	assert.Equal(t, internalapp.DefaultVariant(), set.ThemeVariant())
}

func TestSetThemeVariant_WithCustomTheme(t *testing.T) {
	type customTheme struct {
		fyne.Theme
	}
	set := &settings{}
	ctheme := &customTheme{internalTest.LightTheme(theme.DefaultTheme())}
	set.SetTheme(ctheme)

	set.SetThemeVariant(theme.VariantDark)
	assert.Equal(t, ctheme, set.Theme())
	assert.Equal(t, theme.VariantDark, set.ThemeVariant())
	assert.True(t, set.variantSpecified)

	set.SetThemeVariant(theme.VariantLight)
	assert.Equal(t, ctheme, set.Theme())
	assert.Equal(t, theme.VariantLight, set.ThemeVariant())
	assert.True(t, set.variantSpecified)
}

func TestSetThemeVariant_TriggersListeners(t *testing.T) {
	set := &settings{}
	set.setupTheme()

	listenerCalled := false
	set.AddListener(func(s fyne.Settings) {
		listenerCalled = true
	})

	set.SetThemeVariant(theme.VariantDark)
	assert.True(t, listenerCalled)
}

func TestSetThemeVariant_VariantPersistenceAfterSetupTheme(t *testing.T) {
	set := &settings{}
	set.setupTheme()

	// Set a specific variant
	set.SetThemeVariant(theme.VariantLight)
	assert.Equal(t, theme.VariantLight, set.ThemeVariant())

	// Call setupTheme again - should preserve the set variant
	set.setupTheme()
	assert.Equal(t, theme.VariantLight, set.ThemeVariant())
	assert.True(t, set.variantSpecified)
}

func TestSetThemeVariant_OverridesEnvironmentVariable(t *testing.T) {
	set := &settings{}
	set.schema.ThemeName = "light"

	// Set environment variable to dark
	err := os.Setenv("FYNE_THEME", "dark")
	if err != nil {
		t.Error(err)
	}
	set.setupTheme()
	assert.Equal(t, theme.VariantDark, set.ThemeVariant())

	// SetThemeVariant should override the environment variable
	set.SetThemeVariant(theme.VariantLight)
	assert.Equal(t, theme.VariantLight, set.ThemeVariant())

	// setupTheme should still respect the set variant over environment
	set.setupTheme()
	assert.Equal(t, theme.VariantLight, set.ThemeVariant())
}

func TestSetTheme_AndSetThemeVariant_Interaction(t *testing.T) {
	type customTheme struct {
		fyne.Theme
	}
	set := &settings{}
	ctheme := &customTheme{internalTest.LightTheme(theme.DefaultTheme())}

	// Set variant first
	set.SetThemeVariant(theme.VariantDark)
	assert.Equal(t, theme.VariantDark, set.ThemeVariant())
	assert.True(t, set.variantSpecified)

	// Then set custom theme - variant should be preserved
	set.SetTheme(ctheme)
	assert.Equal(t, ctheme, set.Theme())
	assert.Equal(t, theme.VariantDark, set.ThemeVariant())
	assert.True(t, set.variantSpecified)

	// Change variant after theme is set
	set.SetThemeVariant(theme.VariantLight)
	assert.Equal(t, ctheme, set.Theme())
	assert.Equal(t, theme.VariantLight, set.ThemeVariant())
}

func TestSetThemeVariant_BeforeThemeSpecified(t *testing.T) {
	set := &settings{}

	// Set variant before calling setupTheme
	set.SetThemeVariant(theme.VariantDark)
	assert.Equal(t, theme.VariantDark, set.ThemeVariant())
	assert.True(t, set.variantSpecified)

	// setupTheme should respect the previously set variant
	set.setupTheme()
	assert.Equal(t, theme.VariantDark, set.ThemeVariant())
	assert.True(t, set.variantSpecified)
}

func TestSetThemeVariant_PreservesThroughFileLoad(t *testing.T) {
	set := &settings{}
	set.setupTheme()

	// Set a specific variant
	set.SetThemeVariant(theme.VariantLight)
	assert.Equal(t, theme.VariantLight, set.ThemeVariant())
	assert.True(t, set.variantSpecified)

	// Load settings from file - variant should still be preserved
	err := set.loadFromFile(filepath.Join("testdata", "dark-theme.json"))
	if err != nil {
		t.Error(err)
	}
	set.setupTheme()
	// The set variant should take precedence over file settings
	assert.Equal(t, theme.VariantLight, set.ThemeVariant())
	assert.True(t, set.variantSpecified)
}
