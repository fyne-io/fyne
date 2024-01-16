//go:generate fyne bundle -o translations.go -package lang ../internal/translations/

// Package lang introduces a translation and localisation API for Fyne applications
//
// Since 2.5
package lang

import (
	"encoding/json"
	"os"

	"fyne.io/fyne/v2"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var (
	// L is a shortcut to localise a string, similar to the gettext "_" function.
	L = Localize

	localizer *i18n.Localizer
)

// Localize asks the translation engine to translate a string, this behaves like the gettext "_" function.
func Localize(in string) string {
	ret, err := localizer.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    in,
			Other: in,
		},
	})

	if err != nil {
		fyne.LogError("Translation failure", err)
		return in
	}
	return ret
}

func init() {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	bundle.MustParseMessageFileBytes(resourceBaseFrJson.Content(), resourceBaseFrJson.Name())
	localizer = i18n.NewLocalizer(bundle, os.Getenv("LANG")) // TODO more sophisticated
}
