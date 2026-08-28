package widget

import (
	"testing"
	"time"

	"fyne.io/fyne/v2/test"
	"github.com/stretchr/testify/assert"
)

func TestDateEntry_SetDate(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	e := NewDateEntry()
	date := time.Now()

	e.SetDate(&date)
	assert.NotNil(t, e.Date)
	assert.Equal(t, date.Format(getLocaleDateFormat()), e.Text)

	e.SetDate(nil)
	assert.Nil(t, e.Date)
	assert.Equal(t, "", e.Text)
}

func TestDateEntry_EnableDisable(t *testing.T) {
	e := NewDateEntry()

	e.Disable()
	assert.True(t, e.Disabled())

	e.Enable()
	assert.False(t, e.Disabled())
}

func TestDateEntry_MinSize(t *testing.T) {
	e := NewDateEntry()
	size := e.MinSize()

	assert.True(t, size.Width > 0)
	assert.True(t, size.Height > 0)
}