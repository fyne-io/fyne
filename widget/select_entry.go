package widget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// SelectEntry is an input field which supports selecting from a fixed set of options.
type SelectEntry struct {
	Entry
	popUp   *PopUp
	list    *List
	Options []string
}

// NewSelectEntry creates a SelectEntry.
func NewSelectEntry(options []string) *SelectEntry {
	e := &SelectEntry{Options: options}
	e.ExtendBaseWidget(e)
	e.Wrapping = fyne.TextTruncate
	return e
}

// CreateRenderer returns a new renderer for this select entry.
//
// Implements: fyne.Widget
func (e *SelectEntry) CreateRenderer() fyne.WidgetRenderer {
	e.ExtendBaseWidget(e)

	e.ActionItem = e.setupDropDown()
	if e.Disabled() {
		e.ActionItem.(fyne.Disableable).Disable()
	}

	return e.Entry.CreateRenderer()
}

// Enable this widget, updating any style or features appropriately.
//
// Implements: fyne.DisableableWidget
func (e *SelectEntry) Enable() {
	if e.ActionItem != nil {
		e.ActionItem.(fyne.Disableable).Enable()
	}
	e.Entry.Enable()
}

// Disable this widget so that it cannot be interacted with, updating any style appropriately.
//
// Implements: fyne.DisableableWidget
func (e *SelectEntry) Disable() {
	if e.ActionItem != nil {
		e.ActionItem.(fyne.Disableable).Disable()
	}
	e.Entry.Disable()
}

// MinSize returns the minimal size of the select entry.
//
// Implements: fyne.Widget
func (e *SelectEntry) MinSize() fyne.Size {
	e.ExtendBaseWidget(e)

	min := e.Entry.MinSize()
	padding := fyne.NewSize(4*theme.Padding(), 0)
	for _, item := range e.Options {
		itemMin := fyne.MeasureText(item, theme.TextSize(), fyne.TextStyle{}).Add(padding)
		min = min.Max(itemMin)
	}

	return min
}

// Move changes the relative position of the select entry.
//
// Implements: fyne.Widget
func (e *SelectEntry) Move(pos fyne.Position) {
	e.Entry.Move(pos)
	if e.popUp != nil {
		e.popUp.Move(e.popUpPos())
	}
}

// Resize changes the size of the select entry.
//
// Implements: fyne.Widget
func (e *SelectEntry) Resize(size fyne.Size) {
	e.Entry.Resize(size)
	if e.popUp != nil {
		e.popUp.Resize(e.popUpSize())
	}
}

// SetOptions sets the options the user might select from.
func (e *SelectEntry) SetOptions(options []string) {
	e.Options = options
	e.Refresh()
}

func (e *SelectEntry) popUpSize() fyne.Size {
	length := float32(len(e.Options))
	if length > 10 {
		length = 10
	}

	return fyne.NewSize(e.Size().Width, e.popUp.MinSize().Height*length-3*theme.Padding())
}

func (e *SelectEntry) popUpPos() fyne.Position {
	buttonPos := fyne.CurrentApp().Driver().AbsolutePositionForObject(e.super())
	return buttonPos.Add(fyne.NewPos(0, e.Size().Height-theme.InputBorderSize()))
}

func (e *SelectEntry) onDropDownTapped() {
	c := fyne.CurrentApp().Driver().CanvasForObject(e.super())

	e.popUp = NewPopUp(e.list, c)
	e.popUp.ShowAtPosition(e.popUpPos())
	e.popUp.Resize(e.popUpSize())
}

func (e *SelectEntry) popUpPos() fyne.Position {
	entryPos := fyne.CurrentApp().Driver().AbsolutePositionForObject(e.super())
	return entryPos.Add(fyne.NewPos(0, e.Size().Height-theme.InputBorderSize()))
}

func (e *SelectEntry) setupDropDown() *Button {
	length := func() int {
		return len(e.Options)
	}

	create := func() fyne.CanvasObject {
		return &Button{Importance: LowImportance}
	}

	update := func(item ListItemID, object fyne.CanvasObject) {
		button := object.(*Button)
		button.SetText(e.Options[item])
		button.OnTapped = func() {
			e.SetText(e.Options[item])
			e.popUp.Hide()
		}
	}

	e.list = NewList(length, create, update)
	return &Button{
		Importance: LowImportance,
		Icon:       theme.MenuDropDownIcon(),
		OnTapped:   e.onDropDownTapped,
	}
}
