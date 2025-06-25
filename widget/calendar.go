package widget

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
	"unicode"

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
	displayedMonth time.Time

	monthPrevious *Button
	monthNext     *Button
	monthLabel    *Label

	dates *fyne.Container

	SelectedDate time.Time // currently selected date

	// TODO discuss a config variable WeekStart time.Weekday to override week's first day

	OnChanged func(time.Time) `json:"-"`
}

// NewCalendar creates a calendar instance
//
// Since: 2.6
func NewCalendar(cT time.Time, changed func(time.Time)) *Calendar {
	c := &Calendar{
		displayedMonth: cT,
		OnChanged:      changed,
	}

	c.ExtendBaseWidget(c)
	return c
}

// SetDisplayedDate sets the currently displayed year and month
func (c *Calendar) SetDisplayedDate(date time.Time) {
	if date.IsZero() {
		date = time.Now()
	}

	// Dates are 'normalised', forcing date to start from the start of the month ensures move from March to February
	c.displayedMonth = time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, time.Local)

	// check if label has been created to avoid nil pointer access panic if function is called before
	// widget instanciation
	if c.monthLabel != nil {
		c.monthLabel.SetText(c.monthYear())
		c.dates.Objects = c.calendarObjects()
	}
}

// SetSelectedDate sets the currently selected date
//
// Date selection works only if .Selectable = true
func (c *Calendar) SetSelectedDate(date time.Time) {
	c.SelectedDate = date
	c.updateSelection()
}

// CreateRenderer returns a new WidgetRenderer for this widget.
// This should not be called by regular code, it is used internally to render a widget.
func (c *Calendar) CreateRenderer() fyne.WidgetRenderer {
	c.monthPrevious = NewButtonWithIcon("", theme.NavigateBackIcon(), func() {
		c.SetDisplayedDate(c.displayedMonth.AddDate(0, -1, 0))
	})
	c.monthPrevious.Importance = LowImportance

	c.monthNext = NewButtonWithIcon("", theme.NavigateNextIcon(), func() {
		c.SetDisplayedDate(c.displayedMonth.AddDate(0, 1, 0))
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
	oldName, off := c.displayedMonth.Zone()
	return time.Date(c.displayedMonth.Year(), c.displayedMonth.Month(), dayNum, c.displayedMonth.Hour(), c.displayedMonth.Minute(), 0, 0, time.FixedZone(oldName, off)).In(c.displayedMonth.Location())
}

func (c *Calendar) daysOfMonth() []fyne.CanvasObject {
	start := time.Date(c.displayedMonth.Year(), c.displayedMonth.Month(), 1, 0, 0, 0, 0, c.displayedMonth.Location())
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
			// TODO discuss if it is a good idea to unset selection when tapping the selected date
			/*if selDate := c.dateForButton(dayNum); selDate.Format("02012006") != c.SelectedDate.Format("02012006") {
				c.SelectedDate = selDate
			} else {
				c.SelectedDate = time.Time{}
			}*/

			c.SelectedDate = c.dateForButton(dayNum)
			c.updateSelection()

			if c.OnChanged != nil {
				c.OnChanged(c.SelectedDate)
			}
		})

		if !c.SelectedDate.IsZero() && c.SelectedDate.Year() == c.displayedMonth.Year() && c.SelectedDate.Month() == c.displayedMonth.Month() && c.SelectedDate.Day() == dayNum {
			b.Importance = HighImportance
		} else {
			b.Importance = LowImportance
		}

		buttons = append(buttons, b)
	}

	return buttons
}

func (c *Calendar) updateSelection() {
	if c.dates == nil || len(c.dates.Objects) <= 0 {
		return
	}

	defer c.dates.Refresh()

	if c.SelectedDate.IsZero() || c.SelectedDate.Month() != c.displayedMonth.Month() || c.SelectedDate.Year() != c.displayedMonth.Year() {
		for i := 0; i < len(c.dates.Objects); i++ {
			if b, ok := c.dates.Objects[i].(*Button); ok {
				b.Importance = LowImportance
				b.Refresh()
			}
		}
		return
	}

	for i := 0; i < len(c.dates.Objects); i++ {
		if b, ok := c.dates.Objects[i].(*Button); ok {
			dayNum, _ := strconv.Atoi(b.Text)
			if dayNum == c.SelectedDate.Day() {
				b.Importance = HighImportance
			} else {
				b.Importance = LowImportance
			}
			b.Refresh()
		}
	}
}

func (c *Calendar) monthYear() string {
	return monthName(c.displayedMonth.Format("January")) + fmt.Sprintf(" %04d", c.displayedMonth.Year())
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

func monthName(in string) string {
	r := []rune(lang.X(strings.ToLower(in), in))
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}
