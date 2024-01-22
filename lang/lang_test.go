package lang_test

import (
	"os"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/lang"
	"github.com/stretchr/testify/assert"
)

func TestAddTranslations(t *testing.T) {
	oldLang := os.Getenv("LANG")
	_ = os.Setenv("LANG", "en")

	err := lang.AddTranslations(fyne.NewStaticResource("en.json", []byte(`{
  "Test": "Match"
}`)))
	assert.Nil(t, err)
	assert.Equal(t, "Match", lang.L("Test"))

	err = lang.AddTranslationsForLocale([]byte(`{
  "Test2": "Match2"
}`), "en")
	assert.Nil(t, err)
	assert.Equal(t, "Match2", lang.L("Test2"))

	_ = os.Setenv("LANG", oldLang)
}

func TestLocalize_Fallback(t *testing.T) {
	assert.Equal(t, "Missing", lang.L("Missing"))
}

func TestLocalize_Struct(t *testing.T) {
	type data struct {
		Str string
	}

	assert.Equal(t, "Hello World!", lang.L("Hello {{.Str}}!", data{Str: "World"}))
}

func TestLocalize_Map(t *testing.T) {
	assert.Equal(t, "Hello World!", lang.L("Hello {{.Str}}!", map[string]string{"Str": "World"}))
}

func TestLocalizePlural_Fallback(t *testing.T) {
	_ = lang.AddTranslations(fyne.NewStaticResource("en.json", []byte(`{
  "Apple": {
    "one": "Apple",
    "other": "Apples"
  }
}`)))
	assert.Equal(t, "Missing", lang.N("Missing", 1))
	assert.Equal(t, "Apple", lang.N("Apple", 1))
	assert.Equal(t, "Apples", lang.N("Apple", 2))
}

func TestLocalizeKey_Fallback(t *testing.T) {
	oldLang := os.Getenv("LANG")
	_ = os.Setenv("LANG", "en")

	_ = lang.AddTranslations(fyne.NewStaticResource("en.json", []byte(`{
  "appleID": "Apple Matched"
}`)))

	assert.Equal(t, "Apple", lang.X("appleIDMissing", "Apple"))
	assert.Equal(t, "Apple Matched", lang.X("appleID", "Apple"))

	_ = os.Setenv("LANG", oldLang)
}
