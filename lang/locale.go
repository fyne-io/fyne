package lang

import (
	"fyne.io/fyne/v2"

	"github.com/fyne-io/go-locale"

	"golang.org/x/text/language"
)

type Language string
type Region string
type Script string

type Locale struct {
	Language Language
	Region   Region
	Script   Script
}

func (l Locale) LanguageString() string {
	return string(l.Language) + "-" + string(l.Region)
}

func (l Locale) String() string {
	return l.LanguageString() + "-" + string(l.Script)
}

func SystemLocale() Locale {
	matcher := language.NewMatcher([]language.Tag{
		language.English,
		language.French}) // TODO look this list up from resources loaded

	locs, err := locale.GetLocales()
	if err != nil {
		fyne.LogError("Failed to look up user locale", err)
		locs = []string{"en"}
	}

	tags := make([]language.Tag, len(locs))
	for i, loc := range locs {
		tags[i], err = language.Parse(loc)
		if err != nil {
			fyne.LogError("Error parsing user locale", err)
		}
	}
	best, _, _ := matcher.Match(tags...)
	return localeFromTag(best)
}

func localeFromTag(in language.Tag) Locale {
	b, _ := in.Base()
	r, _ := in.Region()
	s, _ := in.Script()

	return Locale{
		Language: Language(b.String()),
		Region:   Region(r.String()),
		Script:   Script(s.String()),
	}
}
