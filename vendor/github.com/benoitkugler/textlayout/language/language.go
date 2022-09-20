package language

import (
	"os"
	"strings"
)

var canonMap = [256]byte{
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, '-', 0, 0,
	'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 0, 0, 0, 0, 0, 0,
	'-', 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o',
	'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z', 0, 0, 0, 0, '-',
	0, 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o',
	'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z', 0, 0, 0, 0, 0,
}

// Language store the canonicalized BCP 47 tag,
// which has the generic form <lang>-<country>-<other tags>...
type Language string

// NewLanguage canonicalizes the language input (as a BCP 47 language tag), by converting it to
// lowercase, mapping '_' to '-', and stripping all characters other
// than letters, numbers and '-'.
func NewLanguage(language string) Language {
	out := make([]byte, 0, len(language))
	for _, b := range language {
		can := canonMap[b]
		if can != 0 {
			out = append(out, can)
		}
	}
	return Language(out)
}

func (l Language) primary() Language {
	if index := strings.IndexByte(string(l), '-'); index != -1 {
		l = l[:index]
	}
	return l
}

// SimpleInheritance returns the list of matching language, using simple truncation inheritance.
// The resulting slice starts with the given whole language.
// See http://www.unicode.org/reports/tr35/#Locale_Inheritance for more information.
func (l Language) SimpleInheritance() []Language {
	tags := strings.Split(string(l), "-")
	out := make([]Language, 0, len(tags))
	for len(tags) != 0 {
		out = append(out, Language(strings.Join(tags, "-")))
		tags = tags[:len(tags)-1]
	}
	return out
}

// IsDerivedFrom returns `true` if `l` has
// the `root` as primary language.
func (l Language) IsDerivedFrom(root Language) bool { return l.primary() == root }

// IsUndetermined returns `true` if its primary language is "und".
// It is a shortcut for IsDerivedFrom("und").
func (l Language) IsUndetermined() bool { return l.IsDerivedFrom("und") }

// LanguageComparison is a three state enum resulting from comparing two languages
type LanguageComparison uint8

const (
	LanguagesDiffer      LanguageComparison = iota // the two languages are totally differents
	LanguagesExactMatch                            // the two languages are exactly the same
	LanguagePrimaryMatch                           // the two languages have the same primary language, but differs.
)

// Compare compares `other` and `l`.
// Undetermined languages are only compared using the remaining tags,
// meaning that "und-fr" and "und-be" are compared as LanguagesDiffer, not
// LanguagePrimaryMatch.
func (l Language) Compare(other Language) LanguageComparison {
	if l == other {
		return LanguagesExactMatch
	}

	primary1, primary2 := l.primary(), other.primary()
	if primary1 != primary2 {
		return LanguagesDiffer
	}

	// check for the undetermined special case
	if primary1 == "und" {
		return LanguagesDiffer
	}
	return LanguagePrimaryMatch
}

func languageFromLocale(locale string) Language {
	if i := strings.IndexByte(locale, '.'); i >= 0 {
		locale = locale[:i]
	}
	return NewLanguage(locale)
}

// DefaultLanguage returns the language found in environment variables LC_ALL, LC_CTYPE or
// LANG (in that order), or the zero value if not found.
func DefaultLanguage() Language {
	p, ok := os.LookupEnv("LC_ALL")
	if ok {
		return languageFromLocale(p)
	}

	p, ok = os.LookupEnv("LC_CTYPE")
	if ok {
		return languageFromLocale(p)
	}

	p, ok = os.LookupEnv("LANG")
	if ok {
		return languageFromLocale(p)
	}

	return ""
}
