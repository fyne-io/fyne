package widget_test

import (
	"testing"
	"time"

	"fyne.io/fyne/v2/widget"
	"github.com/stretchr/testify/assert"
)

func TestDatePicker(t *testing.T) {

	dp := widget.NewDatePicker()

	// Simulate selecting a date
	dp.YearSelect.SetSelected("2025")
	dp.MonthSelect.SetSelected("02")
	dp.DaySelect.SetSelected("14")

	// Get the formatted date string
	dateStr, err := dp.GetDate()
	assert.NoError(t, err)

	// Parse the formatted string back into a time.Time object
	date, err := time.Parse("2006-01-02", dateStr)
	assert.NoError(t, err)

	// Validate the parsed date
	assert.Equal(t, 2025, date.Year())
	assert.Equal(t, time.February, date.Month())
	assert.Equal(t, 14, date.Day())
}
