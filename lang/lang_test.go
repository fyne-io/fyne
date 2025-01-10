package lang

import (
	"testing"

	"fyne.io/fyne/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddTranslations(t *testing.T) {
	err := AddTranslations(fyne.NewStaticResource("en.json", []byte(`{
		"Test": "Match"
	}`)))
	if assert.NoError(t, err) {
		setupLang("en")
		assert.Equal(t, "Match", L("Test"))
	}

	err = AddTranslationsForLocale([]byte(`{
		"Test2": "Match2"
	}`), "fr")
	if assert.NoError(t, err) {
		setupLang("fr")
		assert.Equal(t, "Match2", L("Test2"))
	}
}

func TestLocalize_Default(t *testing.T) {
	fallback := closestSupportedLocale([]string{"xx"})
	assert.Equal(t, fyne.Locale("en"), fallback[0:2])
}

func TestLocalize_Fallback(t *testing.T) {
	assert.Equal(t, "Missing", L("Missing"))
}

func TestLocalize_Struct(t *testing.T) {
	type data struct {
		Str string
	}

	assert.Equal(t, "Hello World!", L("Hello {{.Str}}!", data{Str: "World"}))
}

func TestLocalize_Map(t *testing.T) {
	assert.Equal(t, "Hello World!", L("Hello {{.Str}}!", map[string]string{"Str": "World"}))
}

func TestLocalizePlural_Fallback(t *testing.T) {
	require.NoError(t, AddTranslations(fyne.NewStaticResource("en.json", []byte(`{
		"Apple": {
			"one": "Apple",
			"other": "Apples"
		}
	}`))))

	setupLang("en")
	assert.Equal(t, "Missing", N("Missing", 1))
	assert.Equal(t, "Apple", N("Apple", 1))
	assert.Equal(t, "Apples", N("Apple", 2))
}

func TestLocalizeKey_Fallback(t *testing.T) {
	require.NoError(t, AddTranslations(fyne.NewStaticResource("en.json", []byte(`{
		"appleID": "Apple Matched"
	}`))))

	setupLang("en")
	assert.Equal(t, "Apple", X("appleIDMissing", "Apple"))
	assert.Equal(t, "Apple Matched", X("appleID", "Apple"))
}

func TestLocalizePluralKey_Fallback(t *testing.T) {
	require.NoError(t, AddTranslations(fyne.NewStaticResource("en.json", []byte(`{
		"appleID": {
			"one": "Apple",
			"other": "Apples"
		}
	}`))))

	setupLang("en")
	assert.Equal(t, "Missing", XN("appleIDMissing", "Missing", 1))
	assert.Equal(t, "Apple", XN("appleID", "Apple", 1))
	assert.Equal(t, "Apples", XN("appleID", "Apple", 2))
}

func TestDefaultLocalizations(t *testing.T) {
	t.Run("base localizations are loaded by default", func(t *testing.T) {
		languages := RegisteredLanguages()

		translationFiles, err := translations.ReadDir("translations")
		require.NoError(t, err)
		assert.Len(t, languages, len(translationFiles))

		// Two samples for sanity check.
		assert.Contains(t, languages, fyne.Locale("en-US-Latn"))
		assert.Contains(t, languages, fyne.Locale("de-DE-Latn"))
	})

	t.Run("registering custom localizations should wipe default localizations first", func(t *testing.T) {
		err := AddTranslations(fyne.NewStaticResource("de.json", []byte(`{
		  "Test": "Passt"
		}`)))
		require.NoError(t, err)

		err = AddTranslations(fyne.NewStaticResource("en_GB.json", []byte(`{
		  "Test": "Match"
		}`)))
		require.NoError(t, err)

		languages := RegisteredLanguages()
		assert.Equal(t, []fyne.Locale{"de-DE-Latn", "en-GB-Latn"}, languages)

		setupLang("de")
		assert.Equal(t, "Passt", L("Test"))

		setupLang("en-GB")
		assert.Equal(t, "Match", L("Test"))

		t.Run("first registered language becomes default", func(t *testing.T) {
			setupLang("it")
			assert.Equal(t, "Passt", L("Test"))
		})

		t.Run("base translations are still present", func(t *testing.T) {
			setupLang("de")
			assert.Equal(t, "Beenden", L("Quit"))

			setupLang("en-GB")
			assert.Equal(t, "Name", X("file.name", "nope"))
		})
	})
}
