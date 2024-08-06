package lang

import (
	"testing"

	"fyne.io/fyne/v2"
	"github.com/stretchr/testify/assert"
)

func TestAddTranslations(t *testing.T) {
	setupLang("en")
	err := AddTranslations(fyne.NewStaticResource("en.json", []byte(`{
  "Test": "Match"
}`)))
	assert.Nil(t, err)
	assert.Equal(t, "Match", L("Test"))

	err = AddTranslationsForLocale([]byte(`{
  "Test2": "Match2"
}`), "fr")
	setupLang("fr")
	assert.Nil(t, err)
	assert.Equal(t, "Match2", L("Test2"))
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
	_ = AddTranslations(fyne.NewStaticResource("en.json", []byte(`{
  "Apple": {
    "one": "Apple",
    "other": "Apples"
  }
}`)))

	setupLang("en")
	assert.Equal(t, "Missing", N("Missing", 1))
	assert.Equal(t, "Apple", N("Apple", 1))
	assert.Equal(t, "Apples", N("Apple", 2))
}

func TestLocalizeKey_Fallback(t *testing.T) {
	_ = AddTranslations(fyne.NewStaticResource("en.json", []byte(`{
  "appleID": "Apple Matched"
}`)))

	setupLang("en")
	assert.Equal(t, "Apple", X("appleIDMissing", "Apple"))
	assert.Equal(t, "Apple Matched", X("appleID", "Apple"))
}

func TestLocalizePluralKey_Fallback(t *testing.T) {
	_ = AddTranslations(fyne.NewStaticResource("en.json", []byte(`{
  "appleID": {
    "one": "Apple",
    "other": "Apples"
  }
}`)))

	setupLang("en")
	assert.Equal(t, "Missing", XN("appleIDMissing", "Missing", 1))
	assert.Equal(t, "Apple", XN("appleID", "Apple", 1))
	assert.Equal(t, "Apples", XN("appleID", "Apple", 2))
}
