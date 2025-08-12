package widget

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLocaleWeekStart(t *testing.T) {
	assert.Equal(t, "Monday", lookupLocaleSetting("").weekStartDay.String())
	assert.Equal(t, "Monday", lookupLocaleSetting("en").weekStartDay.String())
	assert.Equal(t, "Monday", lookupLocaleSetting("en-GB").weekStartDay.String())
	assert.Equal(t, "Sunday", lookupLocaleSetting("en-US").weekStartDay.String())
	assert.Equal(t, "Sunday", lookupLocaleSetting("es-US").weekStartDay.String())
	assert.Equal(t, "Monday", lookupLocaleSetting("de-DE").weekStartDay.String())
}

func TestGetLocaleDateFormat(t *testing.T) {
	assert.Equal(t, "02/01/2006", lookupLocaleSetting("").dateFormat)
	assert.Equal(t, "02/01/2006", lookupLocaleSetting("en").dateFormat)
	assert.Equal(t, "02/01/2006", lookupLocaleSetting("en-GB").dateFormat)
	assert.Equal(t, "01/02/2006", lookupLocaleSetting("en-US").dateFormat)
	assert.Equal(t, "01/02/2006", lookupLocaleSetting("es-US").dateFormat)
	assert.Equal(t, "02.01.2006", lookupLocaleSetting("de-DE").dateFormat)
}
