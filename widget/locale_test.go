package widget

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLocaleWeekStart(t *testing.T) {
	assert.Equal(t, "Monday", lookupLocaleSetting(localeWeekStartKey, ""))
	assert.Equal(t, "Monday", lookupLocaleSetting(localeWeekStartKey, "en"))
	assert.Equal(t, "Monday", lookupLocaleSetting(localeWeekStartKey, "en-GB"))
	assert.Equal(t, "Sunday", lookupLocaleSetting(localeWeekStartKey, "en-US"))
	assert.Equal(t, "Sunday", lookupLocaleSetting(localeWeekStartKey, "es-US"))
}

func TestGetLocaleDateFormat(t *testing.T) {
	assert.Equal(t, "02/01/2006", lookupLocaleSetting(localeDateFormatKey, ""))
	assert.Equal(t, "02/01/2006", lookupLocaleSetting(localeDateFormatKey, "en"))
	assert.Equal(t, "02/01/2006", lookupLocaleSetting(localeDateFormatKey, "en-GB"))
	assert.Equal(t, "01/02/2006", lookupLocaleSetting(localeDateFormatKey, "en-US"))
	assert.Equal(t, "01/02/2006", lookupLocaleSetting(localeDateFormatKey, "es-US"))
}
