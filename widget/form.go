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

	validationError error
}

// NewFormItem creates a new form item with the specified label text and input widget
func NewFormItem(text string, widget fyne.CanvasObject) *FormItem {
	return &FormItem{Text: text, Widget: widget}
}

// Form widget is two column grid where each row has a label and a widget (usually an input).
// The last row of the grid will contain the appropriate form control buttons if any should be shown.
// Setting OnSubmit will set the submit button to be visible and call back the function when tapped.
// Setting OnCancel will do the same for a cancel button.
// If you change OnSubmit/OnCancel after the form is created and rendered, you need to call
// Refresh() to update the form with the correct buttons.
// Setting OnSubmit/OnCancel to nil will remove the buttons.
type Form struct {
	BaseWidget

	Items      []*FormItem
	OnSubmit   func()
	OnCancel   func()
	SubmitText string
	CancelText string

	itemGrid     *fyne.Container
	buttonBox    *Box
	cancelButton *Button
	submitButton *Button
}

func (f *Form) createLabel(text string) *Label {
	return NewLabelWithStyle(text, fyne.TextAlignTrailing, fyne.TextStyle{Bold: true})
}

// Append adds a new row to the form, using the text as a label next to the specified Widget
func (f *Form) Append(text string, widget fyne.CanvasObject) {
	item := &FormItem{Text: text, Widget: widget}
	f.AppendItem(item)
}

// AppendItem adds the specified row to the end of the Form
func (f *Form) AppendItem(item *FormItem) {
	f.ExtendBaseWidget(f) // could be called before render

	f.Items = append(f.Items, item)
	if f.itemGrid != nil {
		f.itemGrid.Add(f.createLabel(item.Text))
		f.itemGrid.Add(item.Widget)
		f.setUpValidation(item.Widget, len(f.Items)-1)
	}

	f.Refresh()
}

// MinSize returns the size that this widget should not shrink below
func (f *Form) MinSize() fyne.Size {
	f.ExtendBaseWidget(f)
	return f.BaseWidget.MinSize()
}

// Refresh updates the widget state when requested.
func (f *Form) Refresh() {
	cache.Renderer(f.super()) // we are about to make changes to renderer created content... not great!
	f.updateButtons()
	f.updateLabels()
	f.BaseWidget.Refresh()
	canvas.Refresh(f.super()) // refresh ourselves for BG color - the above updates the content
}

func (f *Form) updateButtons() {
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

func (f *Form) checkValidation(err error) {
	if err != nil {
		f.submitButton.Disable()
		return
	}

	for _, item := range f.Items {
		if item.validationError != nil {
			f.submitButton.Disable()
			return
		}
	}

	f.submitButton.Enable()
}

func (f *Form) setUpValidation(widget fyne.CanvasObject, i int) {
	if w, ok := widget.(fyne.Validatable); ok {
		f.Items[i].validationError = w.Validate()
		w.SetOnValidationChanged(func(err error) {
			f.Items[i].validationError = err
			f.checkValidation(err)
		})
	}
}

func (f *Form) updateLabels() {
	for i, item := range f.Items {
		l := f.itemGrid.Objects[i*2].(*Label)
		if l.Text == item.Text {
			continue
		}

		l.SetText(item.Text)
	}
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
func (f *Form) CreateRenderer() fyne.WidgetRenderer {
	f.ExtendBaseWidget(f)

	f.cancelButton = &Button{Icon: theme.CancelIcon(), OnTapped: f.OnCancel}
	f.submitButton = &Button{Icon: theme.ConfirmIcon(), OnTapped: f.OnSubmit, Importance: HighImportance}
	f.buttonBox = NewHBox(layout.NewSpacer(), f.cancelButton, f.submitButton)

	objects := make([]fyne.CanvasObject, len(f.Items)*2)
	for i, item := range f.Items {
		objects[i*2] = f.createLabel(item.Text)
		objects[i*2+1] = item.Widget
		f.setUpValidation(item.Widget, i)
	}
	f.itemGrid = fyne.NewContainerWithLayout(layout.NewFormLayout(), objects...)

	renderer := cache.Renderer(NewVBox(f.itemGrid, f.buttonBox))
	f.updateButtons()      // will set correct visibility on the submit/cancel btns
	f.checkValidation(nil) // will trigger a validation check for correct intial validation status
	return renderer
}

// NewForm creates a new form widget with the specified rows of form items
// and (if any of them should be shown) a form controls row at the bottom
func NewForm(items ...*FormItem) *Form {
	form := &Form{Items: items}
	form.ExtendBaseWidget(form)

	return form
}
