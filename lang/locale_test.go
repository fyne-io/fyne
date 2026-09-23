package lang

import (
	"testing"

	"github.com/jeandeaual/go-locale"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
)

func TestClosestSupportedLocale_FallsBackForUnrelatedScript(t *testing.T) {
	// Save and restore the translated list.
	orig := translated
	t.Cleanup(func() { translated = orig })

	translated = []language.Tag{
		language.Make("en"),
		language.Make("ru"),
		language.Make("uk"),
	}

	// Serbian Cyrillic shares its script with Russian but is otherwise unrelated;
	// previously the matcher would resolve it to Russian. It should now fall back to English.
	got := closestSupportedLocale([]string{"sr-Cyrl"})
	assert.Equal(t, "en", got.LanguageString())

	// Confirm that exact and high-confidence matches still resolve normally.
	got = closestSupportedLocale([]string{"ru"})
	assert.Equal(t, "ru", got.LanguageString())

	got = closestSupportedLocale([]string{"uk-UA"})
	assert.Equal(t, "uk", got.LanguageString())
}

func TestSystemLocale(t *testing.T) {
	info, err := locale.GetLocale()
	if err != nil {
		// something not testable
		t.Log("Unable to run locale test because", err)
		return
	}

	if len(info) < 2 {
		info = "en_US"
	}

	loc := SystemLocale()
	assert.Equal(t, info[:2], loc.String()[:2])
}
