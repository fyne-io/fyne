package widget

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/lang"
)

const (
	localeDateFormatKey = "dateformat"
	localeWeekStartKey  = "startofweek"
)

var localeSettings = map[string]map[string]string{
	"": {
		localeDateFormatKey: "02/01/2006",
		localeWeekStartKey:  "Monday",
	},
	"BR": {
		localeWeekStartKey: "Sunday",
	},
	"BZ": {
		localeWeekStartKey: "Sunday",
	},
	"CA": {
		localeWeekStartKey: "Sunday",
	},
	"CO": {
		localeWeekStartKey: "Sunday",
	},
	"DO": {
		localeWeekStartKey: "Sunday",
	},
	"GT": {
		localeWeekStartKey: "Sunday",
	},
	"JP": {
		localeWeekStartKey: "Sunday",
	},
	"MX": {
		localeWeekStartKey: "Sunday",
	},
	"NI": {
		localeWeekStartKey: "Sunday",
	},
	"PE": {
		localeWeekStartKey: "Sunday",
	},
	"PA": {
		localeWeekStartKey: "Sunday",
	},
	"PY": {
		localeWeekStartKey: "Sunday",
	},
	"US": {
		localeDateFormatKey: "01/02/2006",
		localeWeekStartKey:  "Sunday",
	},
	"VE": {
		localeWeekStartKey: "Sunday",
	},
	"ZA": {
		localeWeekStartKey: "Sunday",
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
