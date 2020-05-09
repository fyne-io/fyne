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
	onSubmit   func()
	onCancel   func()
	SubmitText string
	CancelText string

	itemGrid     *fyne.Container
	rendered     bool
	buttonsBox   *Box
	vbox         *Box
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
	canvas.Refresh(f) // refresh ourselves for BG color - the above updates the content
}

// OnSubmit sets the onSubmit func ptr
func (f *Form) OnSubmit(fn func()) *Form {
	f.onSubmit = fn
	f.setButtons()
	return f
}

// OnSubmitWithLabel sets the onSubmit func ptr and sets the name of the button
func (f *Form) OnSubmitWithLabel(lbl string, fn func()) *Form {
	f.SubmitText = lbl
	return f.OnSubmit(fn)
}

// OnCancel sets the onSubmit func ptr
func (f *Form) OnCancel(fn func()) *Form {
	f.onCancel = fn
	f.setButtons()
	return f
}

// OnCancelWithLabel sets the onCancel func ptr and sets the name of the button
func (f *Form) OnCancelWithLabel(lbl string, fn func()) *Form {
	f.CancelText = lbl
	return f.OnCancel(fn)
}

func (f *Form) setButtons() {
	if f.CancelText == "" {
		f.CancelText = "Cancel"
	}
	if f.SubmitText == "" {
		f.SubmitText = "Submit"
	}

	// if there is no renderer yet, exit
	if !f.rendered {
		return
	}

	// remove the buttonBox from the forms vbox
	if len(f.vbox.Children) > 1 {
		f.vbox.Children = f.vbox.Children[:1]
	}

	// create the buttons
	if f.onCancel != nil {
		f.cancelButton = NewButtonWithIcon(f.CancelText, theme.CancelIcon(), f.onCancel)
	}
	if f.onSubmit != nil {
		f.submitButton = NewButtonWithIcon(f.SubmitText, theme.ConfirmIcon(), f.onSubmit)
	}

	// fill in the button box if needed
	switch {
	case f.onCancel != nil && f.onSubmit != nil:
		f.buttonsBox = NewHBox(layout.NewSpacer(), f.cancelButton, f.submitButton)
	case f.onCancel != nil && f.onSubmit == nil:
		f.buttonsBox = NewHBox(layout.NewSpacer(), f.cancelButton)
	case f.onCancel == nil && f.onSubmit != nil:
		f.buttonsBox = NewHBox(layout.NewSpacer(), f.submitButton)
	case f.onCancel == nil && f.onSubmit == nil:
		// we are done here
		f.buttonsBox = nil
		f.vbox.Refresh()
		return
	}

	// add the button box to the itemsGrid and we are done
	f.vbox.Append(f.buttonsBox)
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
func (f *Form) CreateRenderer() fyne.WidgetRenderer {
	f.ExtendBaseWidget(f)
	f.ensureGrid()
	for _, item := range f.Items {
		f.itemGrid.AddObject(f.createLabel(item.Text))
		f.itemGrid.AddObject(item.Widget)
	}
	f.rendered = true
	f.vbox = NewVBox(f.itemGrid)
	f.setButtons()
	renderer := cache.Renderer(f.vbox)
	return renderer
}

// NewForm creates a new form widget with the specified rows of form items
// and (if any of them should be shown) a form controls row at the bottom
func NewForm(items ...*FormItem) *Form {
	form := &Form{BaseWidget: BaseWidget{}, Items: items}
	form.ExtendBaseWidget(form)

	return form
}
