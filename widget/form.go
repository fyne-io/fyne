package widget

import (
	"fyne.io/fyne"
	"fyne.io/fyne/internal/cache"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
)

// FormItem provides the details for a row in a form
type FormItem struct {
	Text   string
	Widget fyne.CanvasObject
}

// NewFormItem creates a new form item with the specified label text and input widget
func NewFormItem(text string, widget fyne.CanvasObject) *FormItem {
	return &FormItem{text, widget}
}

// Form widget is two column grid where each row has a label and a widget (usually an input).
// The last row of the grid will contain the appropriate form control buttons if any should be shown.
// Setting OnSubmit will set the submit button to be visible and call back the function when tapped.
// Setting OnCancel will do the same for a cancel button.
type Form struct {
	BaseWidget

	Items    []*FormItem
	OnSubmit func()
	OnCancel func()

	itemGrid *fyne.Container
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
	f.ExtendBaseWidget(f) // could be called before render

	// ensure we have a renderer set up (that creates itemGrid)...
	cache.Renderer(f.super())

	f.Items = append(f.Items, item)
	f.itemGrid.AddObject(f.createLabel(item.Text))
	f.itemGrid.AddObject(item.Widget)

	f.Refresh()
}

// MinSize returns the size that this widget should not shrink below
func (f *Form) MinSize() fyne.Size {
	f.ExtendBaseWidget(f)
	return f.BaseWidget.MinSize()
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
func (f *Form) CreateRenderer() fyne.WidgetRenderer {
	f.ExtendBaseWidget(f)
	f.ensureGrid()
	for _, item := range f.Items {
		f.itemGrid.AddObject(f.createLabel(item.Text))
		f.itemGrid.AddObject(item.Widget)
	}

	if f.OnCancel == nil && f.OnSubmit == nil {
		return cache.Renderer(NewVBox(f.itemGrid))
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
	return cache.Renderer(NewVBox(f.itemGrid, buttons))
}

// NewForm creates a new form widget with the specified rows of form items
// and (if any of them should be shown) a form controls row at the bottom
func NewForm(items ...*FormItem) *Form {
	form := &Form{BaseWidget: BaseWidget{}, Items: items}
	form.ExtendBaseWidget(form)

	return form
}
