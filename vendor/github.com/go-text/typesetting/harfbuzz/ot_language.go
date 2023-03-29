package harfbuzz

import (
	"strings"

	"github.com/go-text/typesetting/language"
	"github.com/go-text/typesetting/opentype/tables"
)

type langTag struct {
	language string
	tag      tables.Tag
}

// return -1 if `a` < `l`
func (l *langTag) compare(a string) int {
	b := l.language

	p := strings.IndexByte(a, '-')
	// da := len(a)
	if p != -1 {
		// da = p
		a = a[:p]
	}

	p = strings.IndexByte(b, '-')
	// db := len(b)
	if p != -1 {
		// db = p
		b = b[:p]
	}
	// L := min(min(len(a), len(b)), max(da, db))
	return strings.Compare(a, b)
}

func bfindLanguage(lang string) int {
	low, high := 0, len(otLanguages)
	for low <= high {
		mid := (low + high) / 2
		p := &otLanguages[mid]
		cmp := p.compare(lang)
		if cmp < 0 {
			high = mid - 1
		} else if cmp > 0 {
			low = mid + 1
		} else {
			return mid
		}
	}
	return -1
}

func subtagMatches(langStr string, limit int, subtag string) bool {
	LS := len(subtag)
	for {
		s := strings.Index(langStr, subtag)
		if s == -1 || s >= limit {
			return false
		}
		if s+LS >= len(langStr) || !isAlnum(langStr[s+LS]) {
			return true
		}
		langStr = langStr[s+LS:]
	}
}

func langMatches(langStr, spec string) bool {
	l := len(spec)
	return strings.HasPrefix(langStr, spec) && (len(langStr) == l || langStr[l] == '-')
}

func languageToString(l language.Language) string { return string(l) }
