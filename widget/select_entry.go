package widget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// SelectEntry is an input field which supports selecting from a fixed set of options.
type SelectEntry struct {
	Entry
	dropDown *fyne.Menu
	popUp    *PopUpMenu
	options  []string
}

// NewSelectEntry creates a SelectEntry.
func NewSelectEntry(options []string) *SelectEntry {
	e := &SelectEntry{options: options}
	e.ExtendBaseWidget(e)
	e.Wrapping = fyne.TextTruncate
	return e
}

// CreateRenderer returns a new renderer for this select entry.
//
// Implements: fyne.Widget
func (e *SelectEntry) CreateRenderer() fyne.WidgetRenderer {
	e.ExtendBaseWidget(e)
	e.SetOptions(e.options)
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
	return e.Entry.MinSize()
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
		e.popUp.Resize(fyne.NewSize(size.Width, e.popUp.Size().Height))
	}
}

// SetOptions sets the options the user might select from.
func (e *SelectEntry) SetOptions(options []string) {
	e.options = options
	items := make([]*fyne.MenuItem, len(options))
	for i, option := range options {
		option := option // capture
		items[i] = fyne.NewMenuItem(option, func() { e.SetText(option) })
	}
	e.dropDown = fyne.NewMenu("", items...)

	if e.ActionItem == nil {
		e.ActionItem = e.setupDropDown()
		if e.Disabled() {
			e.ActionItem.(fyne.Disableable).Disable()
		}
	}
}

func (e *SelectEntry) popUpPos() fyne.Position {
	entryPos := fyne.CurrentApp().Driver().AbsolutePositionForObject(e.super())
	return entryPos.Add(fyne.NewPos(0, e.Size().Height-theme.InputBorderSize()))
}

func (e *SelectEntry) setupDropDown() *Button {
	dropDownButton := NewButton("", func() {
		c := fyne.CurrentApp().Driver().CanvasForObject(e.super())

		e.popUp = NewPopUpMenu(e.dropDown, c)
		e.popUp.ShowAtPosition(e.popUpPos())
		e.popUp.Resize(fyne.NewSize(e.Size().Width, e.popUp.MinSize().Height))
	})
	dropDownButton.Importance = LowImportance
	dropDownButton.SetIcon(theme.MenuDropDownIcon())
	return dropDownButton
}
