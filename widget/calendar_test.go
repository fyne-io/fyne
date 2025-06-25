package widget

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
)

func TestNewCalendar(t *testing.T) {
	now := time.Now()
	c := NewCalendar(now, func(time.Time) {})
	assert.Equal(t, now.Day(), c.displayedMonth.Day())
	assert.Equal(t, int(now.Month()), int(c.displayedMonth.Month()))
	assert.Equal(t, now.Year(), c.displayedMonth.Year())

	_ = test.WidgetRenderer(c) // and render
	assert.Equal(t, now.Format("January 2006"), c.monthLabel.Text)
}

func TestNewCalendar_ButtonDate(t *testing.T) {
	date := time.Now()
	c := NewCalendar(date, func(time.Time) {})
	_ = test.WidgetRenderer(c) // and render

	endNextMonth := date.AddDate(0, 1, 0).AddDate(0, 0, -(date.Day() - 1))
	last := endNextMonth.AddDate(0, 0, -1)

	firstDate := firstDateButton(c.dates)
	assert.Equal(t, "1", firstDate.Text)
	lastDate := c.dates.Objects[len(c.dates.Objects)-1].(*Button)
	assert.Equal(t, strconv.Itoa(last.Day()), lastDate.Text)
}

func TestNewCalendar_Next(t *testing.T) {
	date := time.Now()
	c := NewCalendar(date, func(time.Time) {})
	_ = test.WidgetRenderer(c) // and render

	assert.Equal(t, date.Format("January 2006"), c.monthLabel.Text)

	test.Tap(c.monthNext)
	date = date.AddDate(0, 1, 0)
	assert.Equal(t, date.Format("January 2006"), c.monthLabel.Text)
}

func TestNewCalendar_Previous(t *testing.T) {
	date := time.Now()
	c := NewCalendar(date, func(time.Time) {})
	_ = test.WidgetRenderer(c) // and render

	assert.Equal(t, date.Format("January 2006"), c.monthLabel.Text)

	test.Tap(c.monthPrevious)
	date = date.AddDate(0, -1, 0)
	assert.Equal(t, date.Format("January 2006"), c.monthLabel.Text)
}

func TestNewCalendar_Resize(t *testing.T) {
	date := time.Now()
	c := NewCalendar(date, func(time.Time) {})
	r := test.WidgetRenderer(c) // and render
	layout := c.dates.Layout.(*calendarLayout)

	baseSize := c.MinSize()
	r.Layout(baseSize)
	min := layout.cellSize

	r.Layout(baseSize.AddWidthHeight(100, 0))
	assert.Greater(t, layout.cellSize.Width, min.Width)
	assert.Equal(t, layout.cellSize.Height, min.Height)

	r.Layout(baseSize.AddWidthHeight(0, 100))
	assert.Equal(t, layout.cellSize.Width, min.Width)
	assert.Greater(t, layout.cellSize.Height, min.Height)

	r.Layout(baseSize.AddWidthHeight(100, 100))
	assert.Greater(t, layout.cellSize.Width, min.Width)
	assert.Greater(t, layout.cellSize.Height, min.Height)
}

func firstDateButton(c *fyne.Container) *Button {
	for _, b := range c.Objects {
		if nonBlank, ok := b.(*Button); ok {
			return nonBlank
		}
	}

	return nil
}
