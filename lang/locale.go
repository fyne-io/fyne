package lang

import (
	"github.com/jeandeaual/go-locale"
	"golang.org/x/text/language"

	"fyne.io/fyne/v2"
)

// SystemLocale returns the primary locale on the current system.
// This may refer to a language that Fyne does not have translations for.
func SystemLocale() fyne.Locale {
	loc, err := locale.GetLocale()
	if err != nil {
		fyne.LogError("Failed to look up user locale", err)
	}
	if len(loc) < 2 {
		loc = "en"
	}

	tag, err := language.Parse(loc)
	if err != nil {
		fyne.LogError("Error parsing user locale "+loc, err)
	}
	return localeFromTag(tag)
}

func closestSupportedLocale(locs []string) fyne.Locale {
	matcher := language.NewMatcher(translated)

	tags := make([]language.Tag, len(locs))
	for i, loc := range locs {
		tag, err := language.Parse(loc)
		if err != nil {
			fyne.LogError("Error parsing user locale "+loc, err)
		}
		tags[i] = tag
	}
	best, _, _ := matcher.Match(tags...)
	return localeFromTag(best)
}

func localeFromTag(in language.Tag) fyne.Locale {
	b, s, r := in.Raw()
	ret := b.String()

	if r.String() != "ZZ" {
		ret += "-" + r.String()

		if s.String() != "Zzzz" {
			ret += "-" + s.String()
		}
	}

	return fyne.Locale(ret)
}
