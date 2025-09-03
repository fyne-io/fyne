package widget

import (
	"math"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
)

// Declare conformity with Layout interface
var _ fyne.Layout = (*calendarLayout)(nil)

const (
	daysPerWeek      = 7
	maxWeeksPerMonth = 6
)

var minCellContent = NewLabel("22")

// Calendar creates a new date time picker which returns a time object
//
// Since: 2.6
type Calendar struct {
	BaseWidget
	currentTime time.Time

	monthPrevious *Button
	monthNext     *Button
	monthLabel    *Label

	dates *fyne.Container

	OnChanged func(time.Time) `json:"-"`
}

// NewCalendar creates a calendar instance
//
// Since: 2.6
func NewCalendar(cT time.Time, changed func(time.Time)) *Calendar {
	c := &Calendar{
		currentTime: cT,
		OnChanged:   changed,
	}

	c.ExtendBaseWidget(c)
	return c
}

// CreateRenderer returns a new WidgetRenderer for this widget.
// This should not be called by regular code, it is used internally to render a widget.
func (c *Calendar) CreateRenderer() fyne.WidgetRenderer {
	c.monthPrevious = NewButtonWithIcon("", theme.NavigateBackIcon(), func() {
		c.currentTime = c.currentTime.AddDate(0, -1, 0)
		// Dates are 'normalised', forcing date to start from the start of the month ensures move from March to February
		c.currentTime = time.Date(c.currentTime.Year(), c.currentTime.Month(), 1, 0, 0, 0, 0, c.currentTime.Location())
		c.monthLabel.SetText(c.monthYear())
		c.dates.Objects = c.calendarObjects()
	})
	c.monthPrevious.Importance = LowImportance

	c.monthNext = NewButtonWithIcon("", theme.NavigateNextIcon(), func() {
		c.currentTime = c.currentTime.AddDate(0, 1, 0)
		c.monthLabel.SetText(c.monthYear())
		c.dates.Objects = c.calendarObjects()
	})
	c.monthNext.Importance = LowImportance

	c.monthLabel = NewLabel(c.monthYear())

	nav := &fyne.Container{
		Layout: layout.NewBorderLayout(nil, nil, c.monthPrevious, c.monthNext),
		Objects: []fyne.CanvasObject{
			c.monthPrevious, c.monthNext,
			&fyne.Container{Layout: layout.NewCenterLayout(), Objects: []fyne.CanvasObject{c.monthLabel}},
		},
	}

	c.dates = &fyne.Container{Layout: newCalendarLayout(), Objects: c.calendarObjects()}

	dateContainer := &fyne.Container{
		Layout:  layout.NewBorderLayout(nav, nil, nil, nil),
		Objects: []fyne.CanvasObject{nav, c.dates},
	}

	return NewSimpleRenderer(dateContainer)
}

func (c *Calendar) calendarObjects() []fyne.CanvasObject {
	offset := 0
	switch getLocaleWeekStart() {
	case "Saturday":
		offset = 6
	case "Sunday":
	default:
		offset = 1
	}

	var columnHeadings []fyne.CanvasObject
	for i := 0; i < daysPerWeek; i++ {
		t := NewLabel(shortDayName(time.Weekday((i + offset) % daysPerWeek).String()))
		t.Alignment = fyne.TextAlignCenter
		columnHeadings = append(columnHeadings, t)
	}
	return append(columnHeadings, c.daysOfMonth()...)
}

func (c *Calendar) dateForButton(dayNum int) time.Time {
	oldName, off := c.currentTime.Zone()
	return time.Date(c.currentTime.Year(), c.currentTime.Month(), dayNum, c.currentTime.Hour(), c.currentTime.Minute(), 0, 0, time.FixedZone(oldName, off)).In(c.currentTime.Location())
}

func (c *Calendar) daysOfMonth() []fyne.CanvasObject {
	start := time.Date(c.currentTime.Year(), c.currentTime.Month(), 1, 0, 0, 0, 0, c.currentTime.Location())
	var buttons []fyne.CanvasObject

	dayIndex := int(start.Weekday())
	// account for Go time pkg starting on sunday at index 0
	switch getLocaleWeekStart() {
	case "Saturday":
		if dayIndex == daysPerWeek-1 {
			dayIndex = 0
		} else {
			dayIndex++
		}
	case "Sunday": // nothing to do
	default:
		if dayIndex == 0 {
			dayIndex += daysPerWeek - 1
		} else {
			dayIndex--
		}
	}

	// add spacers if week doesn't start on Monday
	for i := 0; i < dayIndex; i++ {
		buttons = append(buttons, layout.NewSpacer())
	}

	for d := start; d.Month() == start.Month(); d = d.AddDate(0, 0, 1) {
		dayNum := d.Day()
		s := strconv.Itoa(dayNum)
		b := NewButton(s, func() {
			selectedDate := c.dateForButton(dayNum)

			c.OnChanged(selectedDate)
		})
		b.Importance = LowImportance

		buttons = append(buttons, b)
	}

	return buttons
}

func (c *Calendar) monthYear() string {
	return c.currentTime.Format("January 2006")
}

type calendarLayout struct {
	cellSize fyne.Size
}

func newCalendarLayout() fyne.Layout {
	return &calendarLayout{}
}

// Layout is called to pack all child objects into a specified size.
// For a calendar grid this will pack objects into a table format.
func (g *calendarLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	weeks := 1
	day := 0
	for i, child := range objects {
		if !child.Visible() {
			continue
		}

		if day%daysPerWeek == 0 && i >= daysPerWeek {
			weeks++
		}
		day++
	}

	g.cellSize = fyne.NewSize(size.Width/float32(daysPerWeek),
		size.Height/float32(weeks))
	row, col := 0, 0
	i := 0
	for _, child := range objects {
		if !child.Visible() {
			continue
		}

		lead := g.getLeading(row, col)
		trail := g.getTrailing(row, col)
		child.Move(lead)
		child.Resize(fyne.NewSize(trail.X, trail.Y).Subtract(lead))

		if (i+1)%daysPerWeek == 0 {
			row++
			col = 0
		} else {
			col++
		}
		i++
	}
}

// MinSize sets the minimum size for the calendar
func (g *calendarLayout) MinSize(_ []fyne.CanvasObject) fyne.Size {
	pad := theme.Padding()
	largestMin := minCellContent.MinSize()
	return fyne.NewSize(largestMin.Width*daysPerWeek+pad*(daysPerWeek-1),
		largestMin.Height*maxWeeksPerMonth+pad*(maxWeeksPerMonth-1))
}

// Get the leading edge position of a grid cell.
// The row and col specify where the cell is in the calendar.
func (g *calendarLayout) getLeading(row, col int) fyne.Position {
	x := (g.cellSize.Width) * float32(col)
	y := (g.cellSize.Height) * float32(row)

	return fyne.NewPos(float32(math.Round(float64(x))), float32(math.Round(float64(y))))
}

// Get the trailing edge position of a grid cell.
// The row and col specify where the cell is in the calendar.
func (g *calendarLayout) getTrailing(row, col int) fyne.Position {
	return g.getLeading(row+1, col+1)
}

func shortDayName(in string) string {
	lower := strings.ToLower(in)
	key := lower + ".short"
	long := lang.X(lower, in)
	return strings.ToUpper(lang.X(key, long[:3]))
}
