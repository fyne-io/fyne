package widget

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/lang"
)

type weekday int

const (
	monday weekday = iota // default
	sunday
	saturday
)

type localeSetting struct {
	dateFormat   string
	weekStartDay weekday
}

const defaultDateFormat = "02/01/2006"

var localeSettings = map[string]localeSetting{
	"": {
		dateFormat:   defaultDateFormat,
		weekStartDay: monday,
	},
	"BR": {
		weekStartDay: sunday,
	},
	"BZ": {
		weekStartDay: sunday,
	},
	"CA": {
		weekStartDay: sunday,
	},
	"CO": {
		weekStartDay: sunday,
	},
	"DE": {
		dateFormat: "02.01.2006",
	},
	"DO": {
		weekStartDay: sunday,
	},
	"GT": {
		weekStartDay: sunday,
	},
	"JP": {
		weekStartDay: sunday,
	},
	"MX": {
		weekStartDay: sunday,
	},
	"NI": {
		weekStartDay: sunday,
	},
	"PE": {
		weekStartDay: sunday,
	},
	"PA": {
		weekStartDay: sunday,
	},
	"PY": {
		weekStartDay: sunday,
	},
	"SE": {
		dateFormat: "2006-01-02",
	},
	"US": {
		dateFormat:   "01/02/2006",
		weekStartDay: sunday,
	},
	"VE": {
		weekStartDay: sunday,
	},
	"ZA": {
		weekStartDay: sunday,
	},
}

func getLocaleDateFormat() string {
	s := lookupLocaleSetting(lang.SystemLocale())
	if d := s.dateFormat; d != "" {
		return d
	}
	return defaultDateFormat
}

func getLocaleWeekStart() weekday {
	return lookupLocaleSetting(lang.SystemLocale()).weekStartDay
}

func lookupLocaleSetting(l fyne.Locale) localeSetting {
	region := ""
	language := l.LanguageString()
	if pos := strings.Index(language, "-"); pos != -1 {
		region = strings.Split(language, "-")[1]
	}

	if setting, ok := localeSettings[region]; ok {
		return setting
	}

	return localeSettings[""]
}
