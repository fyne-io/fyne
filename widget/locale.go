package widget

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/lang"
)

const (
	localeDateFormatKey = "locale.date.format"
	localeWeekStartKey  = "locale.date.startweek"
)

var localeSettings = map[string]map[string]string{
	"": {
		localeDateFormatKey: "02/01/2006",
		localeWeekStartKey:  "Monday",
	},
	"US": {
		localeDateFormatKey: "01/02/2006",
		localeWeekStartKey:  "Sunday",
	},
}

func getLocaleDateFormat() string {
	return lookupLocaleSetting(localeDateFormatKey, lang.SystemLocale())
}

func getLocaleWeekStart() string {
	return lookupLocaleSetting(localeWeekStartKey, lang.SystemLocale())
}

func lookupLocaleSetting(key string, l fyne.Locale) string {
	region := ""
	lang := l.LanguageString()
	if pos := strings.Index(lang, "-"); pos != -1 {
		region = strings.Split(lang, "-")[1]
	}

	if settings, ok := localeSettings[region]; ok {
		if setting, ok := settings[key]; ok {
			return setting
		}
	}

	return localeSettings[""][key]
}
