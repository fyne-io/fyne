// DatePicker is a widget for selecting a date. It consists of three dropdowns
// for year, month, and day. The days dropdown is dynamically updated based
// on the selected month and year.
//
// Usage:
//
//	dp := widgets.NewDatePicker()
//	date, err := dp.GetDate()
//	if err != nil {
//	    log.Println("No date selected")
//	} else {
//	    log.Println("Selected date:", date)
//	}

package widget

import (
	"fmt"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
)

// DatePicker is a custom widget for selecting a date.
type DatePicker struct {
	BaseWidget
	YearSelect  *Select
	MonthSelect *Select
	DaySelect   *Select
}

// NewDatePicker creates a new DatePicker widget.
func NewDatePicker() *DatePicker {
	dp := &DatePicker{}
	dp.ExtendBaseWidget(dp)

	// Generate lists of years and months
	years := generateYears()
	months := generateMonths()

	// Get Current Time to use as placeholders
	currentTime := time.Now()

	// Extract the year, month, and day and convert them to strings
	year := currentTime.Year()
	month := currentTime.Month()
	day := currentTime.Day()

	// Convert year and day to strings
	yearStr := strconv.Itoa(year)
	dayStr := strconv.Itoa(day)

	// Get the numeric value of the month (integer)
	monthNum := int(month)

	// Convert the month to a string
	monthStr := strconv.Itoa(monthNum)

	// Create Select widgets for years and months
	dp.YearSelect = NewSelect(years, nil)
	dp.YearSelect.PlaceHolder = yearStr

	dp.MonthSelect = NewSelect(months, nil)
	dp.MonthSelect.PlaceHolder = monthStr

	dp.DaySelect = NewSelect([]string{}, nil) // Initialize with empty days
	dp.DaySelect.PlaceHolder = dayStr

	// Function to update the days based on the selected month and year
	updateDays := func() {
		selectedYear := dp.YearSelect.Selected
		selectedMonth := dp.MonthSelect.Selected

		if selectedYear == "" || selectedMonth == "" {
			return // No selection yet
		}

		// Convert selected year and month to integers
		year := 0
		month := 0
		fmt.Sscanf(selectedYear, "%d", &year)
		fmt.Sscanf(selectedMonth, "%d", &month)

		// Generate the list of days for the selected month and year
		days := generateDays(month, year)

		// Update the DaySelect widget
		dp.DaySelect.Options = days
		dp.DaySelect.Refresh()
	}

	// Set up callbacks for year and month selection
	dp.YearSelect.OnChanged = func(_ string) {
		updateDays()
	}
	dp.MonthSelect.OnChanged = func(_ string) {
		updateDays()
	}

	return dp
}

// CreateRenderer implements fyne.Widget.
func (dp *DatePicker) CreateRenderer() fyne.WidgetRenderer {
	return &datePickerRenderer{
		dp: dp,
	}
}

// datePickerRenderer is the renderer for the DatePicker widget.
type datePickerRenderer struct {
	dp *DatePicker
}

// Layout implements fyne.WidgetRenderer.
func (r *datePickerRenderer) Layout(size fyne.Size) {
	// Define the layout for the YearSelect, MonthSelect, and DaySelect widgets
	yearWidth := size.Width / 3
	monthWidth := size.Width / 3
	dayWidth := size.Width / 3

	r.dp.YearSelect.Resize(fyne.NewSize(yearWidth, size.Height))
	r.dp.MonthSelect.Resize(fyne.NewSize(monthWidth, size.Height))
	r.dp.DaySelect.Resize(fyne.NewSize(dayWidth, size.Height))

	r.dp.YearSelect.Move(fyne.NewPos(0, 0))
	r.dp.MonthSelect.Move(fyne.NewPos(yearWidth, 0))
	r.dp.DaySelect.Move(fyne.NewPos(yearWidth+monthWidth, 0))
}

// MinSize implements fyne.WidgetRenderer.
func (r *datePickerRenderer) MinSize() fyne.Size {
	// Calculate the minimum size required for the DatePicker
	yearMin := r.dp.YearSelect.MinSize()
	monthMin := r.dp.MonthSelect.MinSize()
	dayMin := r.dp.DaySelect.MinSize()

	totalWidth := yearMin.Width + monthMin.Width + dayMin.Width
	maxHeight := fyne.Max(yearMin.Height, fyne.Max(monthMin.Height, dayMin.Height))

	return fyne.NewSize(totalWidth, maxHeight)
}

// Refresh implements fyne.WidgetRenderer.
func (r *datePickerRenderer) Refresh() {
	// Refresh the child widgets
	r.dp.YearSelect.Refresh()
	r.dp.MonthSelect.Refresh()
	r.dp.DaySelect.Refresh()
}

// Objects implements fyne.WidgetRenderer.
func (r *datePickerRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{
		r.dp.YearSelect,
		r.dp.MonthSelect,
		r.dp.DaySelect,
	}
}

// Destroy implements fyne.WidgetRenderer.
func (r *datePickerRenderer) Destroy() {
}

// GetDate returns the selected date as a formatted string (YYYY-MM-DD).
func (dp *DatePicker) GetDate() (string, error) {
	selectedYear := dp.YearSelect.Selected
	selectedMonth := dp.MonthSelect.Selected
	selectedDay := dp.DaySelect.Selected

	if selectedYear == "" || selectedMonth == "" || selectedDay == "" {
		return "", fmt.Errorf("no date selected")
	}

	year := 0
	month := 0
	day := 0
	fmt.Sscanf(selectedYear, "%d", &year)
	fmt.Sscanf(selectedMonth, "%d", &month)
	fmt.Sscanf(selectedDay, "%d", &day)

	date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
	return date.Format("2006-01-02"), nil
}

// Helper functions
func generateYears() []string {
	var years []string
	currentYear := time.Now().Year()
	for year := currentYear; year >= 1960; year-- {
		years = append(years, fmt.Sprintf("%d", year))
	}
	return years
}

func generateMonths() []string {
	var months []string
	for month := 1; month <= 12; month++ {
		months = append(months, fmt.Sprintf("%02d", month)) // Format as two digits
	}
	return months
}

func generateDays(month int, year int) []string {
	var days []string
	daysInMonth := time.Date(year, time.Month(month+1), 0, 0, 0, 0, 0, time.Local).Day()
	for day := 1; day <= daysInMonth; day++ {
		days = append(days, fmt.Sprintf("%02d", day)) // Format as two digits
	}
	return days
}
