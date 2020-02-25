package widget

import (
	"fmt"
	"image/color"
	"math"
	"strconv"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
)

func sgn(a float64) float64 {
	switch {
	case a < 0:
		return -1
	case a > 0:
		return +1
	}
	return 0
}

// extend entry to detect enter / and scrolling
// and propergate to spinner widget
type customEntry struct {
	Entry
	onEnter    func()
	onScrolled func(s *fyne.ScrollEvent)
}

func newCustomEntry() *customEntry {
	cutomEntry := &customEntry{}
	cutomEntry.ExtendBaseWidget(cutomEntry)
	return cutomEntry
}

func (e *customEntry) KeyDown(key *fyne.KeyEvent) {

	if key.Name == fyne.KeyReturn {
		if e.onEnter != nil {
			e.onEnter()
		}
	} else {
		e.Entry.KeyDown(key)
	}
}

func (e *customEntry) Scrolled(s *fyne.ScrollEvent) {
	if e.onScrolled != nil {
		e.onScrolled(s)
	}
}

type Spinner struct {
	fyne.Container
	value    float64
	Min      float64
	Max      float64
	Step     float64
	Fmt      string
	startVal float64

	buttonUp   *Button
	buttonDown *Button
	entry      *customEntry
	integer    bool
}

func NewSpinner(min, max, step float64) *Spinner {
	return newSpinnerImpl(min, max, step, false)
}

func NewIntegralSpinner(min, max, step int) *Spinner {
	return newSpinnerImpl(float64(min), float64(max), float64(step), true)
}

func newSpinnerImpl(min, max, step float64, integer bool) *Spinner {
	buttonUp := NewButtonWithIcon("", theme.MenuDropUpIcon(), func() {})
	buttonDown := NewButtonWithIcon("", theme.MenuDropDownIcon(), func() {})
	updown := NewHBox(buttonUp, buttonDown)

	entry := newCustomEntry()
	nip := &Spinner{
		buttonUp:   buttonUp,
		buttonDown: buttonDown,
		entry:      entry,
		Min:        min,
		Max:        max,
		value:      math.Min(math.Max(min, 1), max),
		startVal:   math.Min(math.Max(min, 1), max),
		Step:       1,
		integer:    integer,
	}
	buttonDown.OnTapped = nip.onDown
	buttonUp.OnTapped = nip.onUp
	//entry.OnChanged = nip.onTextChanged
	entry.onEnter = nip.onEnter
	entry.onScrolled = nip.onScrolled
	nip.Layout = layout.NewBorderLayout(nil, nil, nil, updown)
	nip.AddObject(updown)
	nip.AddObject(entry)
	nip.updateVal()
	return nip
}

func (spinner *Spinner) GetValue() float64 {
	spinner.parseEntry()
	return spinner.value
}

func (spinner *Spinner) SetValue(value float64) {
	spinner.value = value
	spinner.updateVal()
}

func (spinner *Spinner) onEnter() {
	spinner.parseEntry()
}
func (spinner *Spinner) parseEntry() {
	if f, err := strconv.ParseFloat(spinner.entry.Text, 64); err != nil {
		spinner.value = spinner.startVal
	} else {
		spinner.value = math.Min(math.Max(spinner.Min, f), spinner.Max)
	}
	spinner.updateVal()
}
func (spinner *Spinner) onScrolled(e *fyne.ScrollEvent) {
	if e.DeltaY != 0 {
		spinner.value += spinner.Step * sgn(float64(e.DeltaY))
		spinner.updateVal()
	}
}
func (spinner *Spinner) onUp() {
	spinner.parseEntry()
	spinner.value += spinner.Step
	spinner.updateVal()
}
func (spinner *Spinner) onDown() {
	spinner.parseEntry()
	spinner.value -= spinner.Step
	spinner.updateVal()
}
func (spinner *Spinner) updateVal() {
	spinner.value = math.Min(math.Max(spinner.Min, spinner.value), spinner.Max)
	if spinner.integer {
		spinner.value = math.Round(spinner.value)
		spinner.entry.SetText(fmt.Sprintf("%d", int(spinner.value)))
	} else if spinner.Fmt == "" {
		spinner.entry.SetText(fmt.Sprintf("%.4f", spinner.value))
	} else {
		spinner.entry.SetText(fmt.Sprintf(spinner.Fmt, spinner.value))
	}
	if spinner.value <= spinner.Min {
		spinner.buttonDown.Disable()
	} else {
		spinner.buttonDown.Enable()
	}
	if spinner.value >= spinner.Max {
		spinner.buttonUp.Disable()
	} else {
		spinner.buttonUp.Enable()
	}

}

func (spinner *Spinner) CreateRenderer() fyne.WidgetRenderer {
	return &spinnerRenderer{
		spinner: spinner,
	}
}

type spinnerRenderer struct {
	spinner *Spinner
}

func (renderer *spinnerRenderer) Layout(size fyne.Size) {
	renderer.spinner.Layout.Layout(renderer.spinner.Objects, size)
}
func (renderer *spinnerRenderer) MinSize() fyne.Size {
	return renderer.spinner.MinSize()
}

func (renderer *spinnerRenderer) Refresh() {
	renderer.spinner.Refresh()
}

func (renderer *spinnerRenderer) BackgroundColor() color.Color {
	return theme.BackgroundColor()
}

func (renderer *spinnerRenderer) Objects() []fyne.CanvasObject {
	return renderer.spinner.Objects
}

func (renderer *spinnerRenderer) Destroy() {
}
