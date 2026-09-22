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

// CurrentLocale returns the currently set locale. If no locale has been set previously,
// it returns an empty string and the package uses the system locale.
func CurrentLocale() fyne.Locale {
	return currentLocale
}

// SetLocale sets the closest supported locale to the given locale as the current locale.
// If an empty string is given, the package sets the system locale as the current locale.
func SetLocale(loc fyne.Locale) {
	currentLocale = ""
	if loc != "" {
		currentLocale = closestSupportedLocale([]string{loc.String()})
	}
	updateLocalizer()
}

func closestSupportedLocale(locs []string) fyne.Locale {
	matcher := language.NewMatcher(translated)

	tags := make([]language.Tag, len(locs))
	for i, loc := range locs {
		tag, err := language.Parse(loc)
		if err != nil && loc != "C" {
			fyne.LogError("Error parsing user locale "+loc, err)
		}
		tags[i] = tag
	}
	best, _, conf := matcher.Match(tags...)
	// When confidence is No the matcher may pick a language that shares only a script
	// (for example Serbian Cyrillic resolving to Russian). Prefer the default fallback
	// in that case rather than presenting an unrelated language.
	if conf == language.No && len(translated) > 0 {
		best = translated[0]
	}
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
