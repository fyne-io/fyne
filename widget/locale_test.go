package widget

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLocaleWeekStart(t *testing.T) {
	assert.Equal(t, monday, lookupLocaleSetting("").weekStartDay)
	assert.Equal(t, monday, lookupLocaleSetting("en").weekStartDay)
	assert.Equal(t, monday, lookupLocaleSetting("en-GB").weekStartDay)
	assert.Equal(t, sunday, lookupLocaleSetting("en-US").weekStartDay)
	assert.Equal(t, sunday, lookupLocaleSetting("es-US").weekStartDay)
	assert.Equal(t, monday, lookupLocaleSetting("de-DE").weekStartDay)
}

func TestGetLocaleDateFormat(t *testing.T) {
	assert.Equal(t, "02/01/2006", lookupLocaleSetting("").dateFormat)
	assert.Equal(t, "02/01/2006", lookupLocaleSetting("en").dateFormat)
	assert.Equal(t, "02/01/2006", lookupLocaleSetting("en-GB").dateFormat)
	assert.Equal(t, "01/02/2006", lookupLocaleSetting("en-US").dateFormat)
	assert.Equal(t, "01/02/2006", lookupLocaleSetting("es-US").dateFormat)
	assert.Equal(t, "02.01.2006", lookupLocaleSetting("de-DE").dateFormat)
}
