//go:generate fyne bundle -o translations.go -package lang ../internal/translations/

// Package lang introduces a translation and localisation API for Fyne applications
//
// Since 2.5
package lang

import (
	"encoding/json"

	"fyne.io/fyne/v2"

	"github.com/nicksnyder/go-i18n/v2/i18n"

	"golang.org/x/text/language"
)

var (
	// L is a shortcut to localize a string, similar to the gettext "_" function.
	// More info available on the `Localize` function
	L = Localize

	// N is a shortcut to localize a string with plural forms, similar to the ngettext function.
	// More info available on the `LocalizePlural` function.
	N = LocalizePlural

	bundle    *i18n.Bundle
	localizer *i18n.Localizer
)

// Localize asks the translation engine to translate a string, this behaves like the gettext "_" function.
// The string can be templated and the template data can be passed as a struct with exported fields,
// or as a map of string keys to any suitable value.
func Localize(in string, data ...any) string {
	var d0 any
	if len(data) > 0 {
		d0 = data[0]
	}

	ret, err := localizer.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    in,
			Other: in,
		},
		TemplateData: d0,
	})

	if err != nil {
		fyne.LogError("Translation failure", err)
		return in
	}
	return ret
}

// LocalizePlural asks the translation engine to translate a string from one of a number of plural forms.
// This behaves like the ngettext function, with the `count` parameter determining the plurality looked up.
// The string can be templated and the template data can be passed as a struct with exported fields,
// or as a map of string keys to any suitable value.
func LocalizePlural(in string, count int, data ...any) string {
	var d0 any
	if len(data) > 0 {
		d0 = data[0]
	}

	ret, err := localizer.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    in,
			Other: in,
		},
		PluralCount:  count,
		TemplateData: d0,
	})

	if err != nil {
		fyne.LogError("Translation failure", err)
		return in
	}
	return ret
}

// AddLanguage allows an app to load a bundle of translations.
// The language that this relates to will be inferred from the resource name, for example "fr.json".
func AddLanguage(r fyne.Resource) error {
	return addLanguage(r.Content(), r.Name())
}

// AddLanguageForLocale allows an app to load a bundle of translations for a specified locale.
func AddLanguageForLocale(r fyne.Resource, l Locale) error {
	return addLanguage(r.Content(), l.String())
}

func addLanguage(data []byte, name string) error {
	_, err := bundle.ParseMessageFileBytes(data, name)
	return err
}

func init() {
	bundle = i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	bundle.MustParseMessageFileBytes(resourceBaseFrJson.Content(), resourceBaseFrJson.Name())
	str := SystemLocale().LanguageString()
	localizer = i18n.NewLocalizer(bundle, str)
}
