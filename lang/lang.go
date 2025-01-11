// Package lang introduces a translation and localisation API for Fyne applications
//
// By default, a few custom localisation keys for a few default languages are included.
// They can be overwritten with custom localisations.
// Without further intervention, the default language will be English. This is overwritten
// with the language of the first translation that is being registered explicitly.
//
// Since 2.5
package lang

import (
	"embed"
	"encoding/json"
	"log"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/jeandeaual/go-locale"
	"github.com/nicksnyder/go-i18n/v2/i18n"

	"fyne.io/fyne/v2"

	"golang.org/x/text/language"
)

var (
	// L is a shortcut to localize a string, similar to the gettext "_" function.
	// More info available on the `Localize` function.
	L = Localize

	// N is a shortcut to localize a string with plural forms, similar to the ngettext function.
	// More info available on the `LocalizePlural` function.
	N = LocalizePlural

	// X is a shortcut to get the localization of a string with specified key, similar to pgettext.
	// More info available on the `LocalizeKey` function.
	X = LocalizeKey

	// XN is a shortcut to get the localization plural form of a string with specified key, similar to npgettext.
	// More info available on the `LocalizePluralKey` function.
	XN = LocalizePluralKey

	bundle          *i18n.Bundle
	localizer       *i18n.Localizer
	bundleIsDefault bool
	unmarshallers   map[string]i18n.UnmarshalFunc

	//go:embed translations
	defaultTranslationFS embed.FS
	defaultTranslations  map[language.Tag]*i18n.MessageFile
	translated           []language.Tag
)

// Localize asks the translation engine to translate a string, this behaves like the gettext "_" function.
// The string can be templated and the template data can be passed as a struct with exported fields,
// or as a map of string keys to any suitable value.
func Localize(in string, data ...any) string {
	return LocalizeKey(in, in, data...)
}

// LocalizeKey asks the translation engine for the translation with specific ID.
// If it cannot be found then the fallback will be used.
// The string can be templated and the template data can be passed as a struct with exported fields,
// or as a map of string keys to any suitable value.
func LocalizeKey(key, fallback string, data ...any) string {
	var d0 any
	if len(data) > 0 {
		d0 = data[0]
	}

	ret, err := localizer.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    key,
			Other: fallback,
		},
		TemplateData: d0,
	})

	if err != nil {
		fyne.LogError("Translation failure", err)
		return fallbackWithData(key, fallback, d0)
	}
	return ret
}

// LocalizePlural asks the translation engine to translate a string from one of a number of plural forms.
// This behaves like the ngettext function, with the `count` parameter determining the plurality looked up.
// The string can be templated and the template data can be passed as a struct with exported fields,
// or as a map of string keys to any suitable value.
func LocalizePlural(in string, count int, data ...any) string {
	return LocalizePluralKey(in, in, count, data...)
}

// LocalizePluralKey asks the translation engine for the translation with specific ID in plural form.
// This behaves like the npgettext function, with the `count` parameter determining the plurality looked up.
// If it cannot be found then the fallback will be used.
// The string can be templated and the template data can be passed as a struct with exported fields,
// or as a map of string keys to any suitable value.
func LocalizePluralKey(key, fallback string, count int, data ...any) string {
	var d0 any
	if len(data) > 0 {
		d0 = data[0]
	}

	ret, err := localizer.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    key,
			Other: fallback,
		},
		PluralCount:  count,
		TemplateData: d0,
	})

	if err != nil {
		fyne.LogError("Translation failure", err)
		return fallbackWithData(key, fallback, d0)
	}
	return ret
}

// AddTranslations allows an app to load a bundle of translations.
// The language that this relates to will be inferred from the resource name, for example "fr.json".
// The data should be in json format.
func AddTranslations(r fyne.Resource) error {
	defer updateLocalizer()
	return addLanguage(r.Content(), r.Name())
}

// AddTranslationsForLocale allows an app to load a bundle of translations for a specified locale.
// The data should be in json format.
func AddTranslationsForLocale(data []byte, l fyne.Locale) error {
	defer updateLocalizer()
	return addLanguage(data, l.String()+".json")
}

// AddTranslationsFS supports adding all translations in one calling using an `embed.FS` setup.
// The `dir` parameter specifies the name or path of the directory containing translation files
// inside this embedded filesystem.
// Each file should be a json file with the name following pattern [prefix.]lang.json.
func AddTranslationsFS(fs embed.FS, dir string) (retErr error) {
	files, err := fs.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, f := range files {
		name := f.Name()
		data, err := fs.ReadFile(filepath.Join(dir, name))
		if err != nil {
			if retErr == nil {
				retErr = err
			}
			continue
		}

		err = addLanguage(data, name)
		if err != nil {
			if retErr == nil {
				retErr = err
			}
			continue
		}
	}

	updateLocalizer()

	return retErr
}

// RegisteredLanguages returns the list of locales we have localization for.
func RegisteredLanguages() []fyne.Locale {
	var ret []fyne.Locale
	for _, l := range translated {
		ret = append(ret, localeFromTag(l))
	}
	return ret
}

func addLanguage(data []byte, name string) error {
	mf, err := i18n.ParseMessageFileBytes(data, name, unmarshallers)
	if err != nil {
		return err
	}

	return addLanguageByMessages(mf.Messages, mf.Tag)
}

func addLanguageByMessages(msgs []*i18n.Message, tag language.Tag) error {
	if bundleIsDefault {
		// This is the first time, the user adds a language on their own. This implies, that
		// they are aware of localizations and how to use them.
		// To prevent half-translated applications, we "forget" our base translations and only
		// use languages for resolution, that the user supplied.
		// The bundle will still contain our base translations; if a language overlaps with those,
		// our base strings will be there unless they explicitly get overwritten.
		bundleIsDefault = false
		if err := setupBundle(tag); err != nil {
			return err
		}
		translated = []language.Tag{tag}
	}

	languageKnown := false
	for _, t := range translated {
		if t == tag {
			languageKnown = true
			break
		}
	}
	if err := bundle.AddMessages(tag, msgs...); err != nil {
		return err
	}
	if !languageKnown {
		translated = append(translated, tag)
	}
	return nil
}

func init() {
	unmarshallers = map[string]i18n.UnmarshalFunc{
		"json": json.Unmarshal,
	}

	// Put English first as our preferred fallback language.
	translated = []language.Tag{language.English}
	if err := loadDefaultTranslations(); err != nil {
		fyne.LogError("Error occurred loading built-in translations", err)
	}
	if err := setupBundle(language.English); err != nil {
		fyne.LogError("Error occurred loading built-in translations", err)
	}
	bundleIsDefault = true
}

func loadDefaultTranslations() error {
	defaultTranslations = map[language.Tag]*i18n.MessageFile{}
	const dir = "translations"
	files, err := defaultTranslationFS.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, f := range files {
		name := f.Name()
		data, err := defaultTranslationFS.ReadFile(filepath.Join(dir, name))
		if err != nil {
			return err
		}

		mf, err := i18n.ParseMessageFileBytes(data, name, unmarshallers)
		if err != nil {
			return err
		}
		defaultTranslations[mf.Tag] = mf
	}
	return nil
}

func setupBundle(lng language.Tag) error {
	bundle = i18n.NewBundle(lng)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	translated = []language.Tag{lng} // the first item in this list will be the fallback if none match
	for tag, mf := range defaultTranslations {
		if err := addLanguageByMessages(mf.Messages, tag); err != nil {
			return err
		}
	}
	updateLocalizer()

	return nil
}

func fallbackWithData(key, fallback string, data any) string {
	t, err := template.New(key).Parse(fallback)
	if err != nil {
		log.Println("Could not parse fallback template")
		return fallback
	}
	str := &strings.Builder{}
	_ = t.Execute(str, data)
	return str.String()
}

// A utility for setting up languages - available to unit tests for overriding system
func setupLang(lang string) {
	localizer = i18n.NewLocalizer(bundle, lang)
}

// updateLocalizer Finds the closest translation from the user's locale list and sets it up
func updateLocalizer() {
	var languageStr string
	all, err := locale.GetLocales()
	if err == nil {
		languageStr = closestSupportedLocale(all).LanguageString()
	} else {
		fyne.LogError("Failed to load user locales", err)
		languageStr = translated[0].String()
	}
	setupLang(languageStr)
}
