package widget

import (
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
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
// Calling OnSubmit will set the submit button to be visible and call back the function when tapped.
// Calling OnCancel will do the same for a cancel button.
type Form struct {
	BaseWidget

	Items      []*FormItem
	OnSubmit   func()
	OnCancel   func()
	SubmitText string
	CancelText string

	itemGrid     *fyne.Container
	vbox         *Box
	buttonBox    *Box
	cancelButton *Button
	submitButton *Button
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
func (f *Form) Append(text string, widget fyne.CanvasObject) *Form {
	item := &FormItem{Text: text, Widget: widget}
	return f.AppendItem(item)
}

// AppendItem adds the specified row to the end of the Form
func (f *Form) AppendItem(item *FormItem) *Form {
	f.ExtendBaseWidget(f) // could be called before render

	// ensure we have a renderer set up (that creates itemGrid)...
	cache.Renderer(f.super())

	f.Items = append(f.Items, item)
	f.itemGrid.AddObject(f.createLabel(item.Text))
	f.itemGrid.AddObject(item.Widget)

	f.Refresh()
	return f
}

// MinSize returns the size that this widget should not shrink below
func (f *Form) MinSize() fyne.Size {
	f.ExtendBaseWidget(f)
	return f.BaseWidget.MinSize()
}

// Refresh updates the widget state when requested.
func (f *Form) Refresh() {
	f.BaseWidget.Refresh()
	f.setButtons()
	canvas.Refresh(f) // refresh ourselves for BG color - the above updates the content
}

func (f *Form) setButtons() {
	if f.CancelText == "" {
		f.CancelText = "Cancel"
	}
	if f.SubmitText == "" {
		f.SubmitText = "Submit"
	}

	// set visibility on the buttons
	if f.OnCancel == nil {
		f.cancelButton.Hide()
	} else {
		f.cancelButton.SetText(f.CancelText)
		f.cancelButton.OnTapped = f.OnCancel
		f.cancelButton.Show()
	}
	if f.OnSubmit == nil {
		f.submitButton.Hide()
	} else {
		f.submitButton.SetText(f.SubmitText)
		f.submitButton.OnTapped = f.OnSubmit
		f.submitButton.Show()
	}
	if f.OnCancel == nil && f.OnSubmit == nil {
		f.buttonBox.Hide()
	} else {
		f.buttonBox.Show()
	}
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
func (f *Form) CreateRenderer() fyne.WidgetRenderer {
	f.ExtendBaseWidget(f)
	f.ensureGrid()
	for _, item := range f.Items {
		f.itemGrid.AddObject(f.createLabel(item.Text))
		f.itemGrid.AddObject(item.Widget)
	}

	// Create the buttons and hide the box till they are needed
	f.cancelButton = NewButtonWithIcon("", theme.CancelIcon(), f.OnCancel)
	f.submitButton = NewButtonWithIcon("", theme.ConfirmIcon(), f.OnSubmit)
	f.submitButton.Hide()
	f.cancelButton.Hide()
	f.buttonBox = NewHBox(layout.NewSpacer(), f.cancelButton, f.submitButton)
	f.buttonBox.Hide()

	f.vbox = NewVBox(f.itemGrid, f.buttonBox)
	renderer := cache.Renderer(f.vbox)
	f.setButtons()
	return renderer
}

// NewForm creates a new form widget with the specified rows of form items
// and (if any of them should be shown) a form controls row at the bottom
func NewForm(items ...*FormItem) *Form {
	form := &Form{BaseWidget: BaseWidget{}, Items: items}
	form.ExtendBaseWidget(form)

	return form
}
