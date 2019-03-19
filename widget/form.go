package widget

import (
	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
)

// FormItem provides the details for a row in a form
type FormItem struct {
	Text   string
	Widget fyne.CanvasObject
}

// Form widget is two column grid where each row has a label and a widget (usually an input).
// The last row of the grid will contain the appropriate form control buttons if any should be shown.
// Setting OnSubmit will set the submit button to be visible and call back the function when tapped.
// Setting OnCancel will do the same for a cancel button.
type Form struct {
	baseWidget

	Items    []*FormItem
	OnSubmit func()
	OnCancel func()

	itemGrid *fyne.Container
}

// Resize sets a new size for a widget.
// Note this should not be used if the widget is being managed by a Layout within a Container.
func (f *Form) Resize(size fyne.Size) {
	f.resize(size, f)
}

// Move the widget to a new position, relative to it's parent.
// Note this should not be used if the widget is being managed by a Layout within a Container.
func (f *Form) Move(pos fyne.Position) {
	f.move(pos, f)
}

// MinSize returns the smallest size this widget can shrink to
func (f *Form) MinSize() fyne.Size {
	return f.minSize(f)
}

// Show this widget, if it was previously hidden
func (f *Form) Show() {
	f.show(f)
}

// Hide this widget, if it was previously visible
func (f *Form) Hide() {
	f.hide(f)
}

func (f *Form) createLabel(text string) *Label {
	return NewLabelWithStyle(text, fyne.TextAlignTrailing, fyne.TextStyle{Bold: true})
}

func (f *Form) ensureGrid() {
	if f.itemGrid != nil {
		return
	}

	f.itemGrid = fyne.NewContainerWithLayout(layout.NewFormLayout(), []fyne.CanvasObject{}...)
}

// Append adds a new row to the form, using the text as a label next to the specified Widget
func (f *Form) Append(text string, widget fyne.CanvasObject) {
	item := &FormItem{Text: text, Widget: widget}
	f.AppendItem(item)
}

// AppendItem adds the specified row to the end of the Form
func (f *Form) AppendItem(item *FormItem) {
	// ensure we have a renderer set up
	Renderer(f)

	f.Items = append(f.Items, item)
	f.itemGrid.AddObject(f.createLabel(item.Text))
	f.itemGrid.AddObject(item.Widget)

	Refresh(f)
}

// CreateRenderer is a private method to Fyne which links this widget to it's renderer
func (f *Form) CreateRenderer() fyne.WidgetRenderer {
	f.ensureGrid()
	for _, item := range f.Items {
		f.itemGrid.AddObject(f.createLabel(item.Text))
		f.itemGrid.AddObject(item.Widget)
	}

	buttons := NewHBox(layout.NewSpacer())
	if f.OnCancel != nil {
		buttons.Append(NewButtonWithIcon("Cancel", theme.CancelIcon(), f.OnCancel))
	}
	if f.OnSubmit != nil {
		submit := NewButtonWithIcon("Submit", theme.ConfirmIcon(), f.OnSubmit)
		submit.Style = PrimaryButton

		buttons.Append(submit)
	}
	return Renderer(NewVBox(f.itemGrid, buttons))
}

// NewForm creates a new form widget with the specified rows of form items
// and (if any of them should be shown) a form controls row at the bottom
func NewForm(items ...*FormItem) *Form {
	form := &Form{baseWidget: baseWidget{}, Items: items}

	Renderer(form).Layout(form.MinSize())
	return form
}
