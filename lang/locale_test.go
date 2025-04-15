package lang

import (
	"testing"

	"github.com/jeandeaual/go-locale"
	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
)

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

func TestClosestSupportedLocale(t *testing.T) {
	assert.Equal(t, fyne.Locale("en"), closestSupportedLocale([]string{"en"}))
	assert.Equal(t, fyne.Locale("en"), closestSupportedLocale([]string{"en_GB"}))
	assert.Equal(t, fyne.Locale("pt-BR"), closestSupportedLocale([]string{"pt"}))

	assert.Equal(t, fyne.Locale("uk"), closestSupportedLocale([]string{"uk"}))
	// Don't fall back to uk from ru #5671
	assert.Equal(t, fyne.Locale("en"), closestSupportedLocale([]string{"ru"}))
}
