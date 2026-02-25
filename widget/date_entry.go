package widget

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

var (
	_ fyne.Widget      = (*DateEntry)(nil)
	_ fyne.Tappable    = (*DateEntry)(nil)
	_ fyne.Disableable = (*DateEntry)(nil)
)

// DateEntry is an input field which supports selecting from a fixed set of options.
//
// Since: 2.6
type DateEntry struct {
	Entry
	Date      *time.Time
	OnChanged func(*time.Time) `json:"-"`

	dropDown *Calendar
	popUp    *PopUp
}

// NewDateEntry creates a date input where the date can be selected or typed.
//
// Since: 2.6
func NewDateEntry() *DateEntry {
	e := &DateEntry{}
	e.ExtendBaseWidget(e)
	e.Wrapping = fyne.TextWrap(fyne.TextTruncateClip)
	return e
}

// CreateRenderer returns a new renderer for this select entry.
func (e *DateEntry) CreateRenderer() fyne.WidgetRenderer {
	e.ExtendBaseWidget(e)

	dateFormat := getLocaleDateFormat()
	e.Validator = func(in string) error {
		_, err := time.Parse(dateFormat, in)
		return err
	}
	e.Entry.OnChanged = func(in string) {
		if in == "" {
			e.Date = nil

			if f := e.OnChanged; f != nil {
				f(nil)
			}
		}
		t, err := time.Parse(dateFormat, in)
		if err != nil {
			return
		}

		e.Date = &t

		if f := e.OnChanged; f != nil {
			f(&t)
		}
	}

	if e.ActionItem == nil {
		e.ActionItem = e.setupDropDown()
		if e.Disabled() {
			e.ActionItem.(fyne.Disableable).Disable()
		}
	}

	return e.Entry.CreateRenderer()
}

// Enable this widget, updating any style or features appropriately.
func (e *DateEntry) Enable() {
	if e.ActionItem != nil {
		if d, ok := e.ActionItem.(fyne.Disableable); ok {
			d.Enable()
		}
	}
	e.Entry.Enable()
}

// Disable this widget so that it cannot be interacted with, updating any style appropriately.
func (e *DateEntry) Disable() {
	if e.ActionItem != nil {
		if d, ok := e.ActionItem.(fyne.Disableable); ok {
			d.Disable()
		}
	}
	e.Entry.Disable()
}

// MinSize returns the minimal size of the select entry.
func (e *DateEntry) MinSize() fyne.Size {
	e.ExtendBaseWidget(e)
	return e.Entry.MinSize()
}

// Move changes the relative position of the date entry.
func (e *DateEntry) Move(pos fyne.Position) {
	e.Entry.Move(pos)
	if e.popUp != nil {
		e.popUp.Move(e.popUpPos())
	}
}

// Resize changes the size of the date entry.
func (e *DateEntry) Resize(size fyne.Size) {
	e.Entry.Resize(size)
	if e.popUp != nil {
		e.popUp.Resize(fyne.NewSize(size.Width, e.popUp.Size().Height))
	}
}

// SetDate will update the widget to a specific date.
// You can pass nil to unselect a date.
func (e *DateEntry) SetDate(d *time.Time) {
	if d == nil {
		e.Date = nil
		e.Entry.SetText("")

		return
	}

	e.setDate(*d)
}

func (e *DateEntry) popUpPos() fyne.Position {
	entryPos := fyne.CurrentApp().Driver().AbsolutePositionForObject(e.super())
	return entryPos.Add(fyne.NewPos(0, e.Size().Height-e.Theme().Size(theme.SizeNameInputBorder)))
}

func (e *DateEntry) setDate(d time.Time) {
	e.Date = &d
	if e.popUp != nil {
		e.popUp.Hide()
	}

	e.Entry.SetText(d.Format(getLocaleDateFormat()))
}

func (e *DateEntry) setupDropDown() *Button {
	if e.dropDown == nil {
		e.dropDown = NewCalendar(time.Now(), e.setDate)
	}
	dropDownButton := NewButton("", func() {
		c := fyne.CurrentApp().Driver().CanvasForObject(e.super())

		e.popUp = NewPopUp(e.dropDown, c)
		e.popUp.ShowAtPosition(e.popUpPos())
		e.popUp.Resize(fyne.NewSize(e.Size().Width, e.popUp.MinSize().Height))
	})
	dropDownButton.Importance = LowImportance
	dropDownButton.SetIcon(e.Theme().Icon(theme.IconNameCalendar))
	return dropDownButton
}
