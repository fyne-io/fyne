package widget

import (
	"errors"
	"reflect"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/internal/cache"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
)

// errFormItemInitialState defines the error if the initial validation for a FormItem result
// in an error
var errFormItemInitialState = errors.New("widget.FormItem initial state error")

// FormItem provides the details for a row in a form
type FormItem struct {
	Text   string
	Widget fyne.CanvasObject

	// Since: 2.0
	HintText string

	validationError error
	invalid         bool
	helperOutput    *canvas.Text
	wasFocused      bool
}

// NewFormItem creates a new form item with the specified label text and input widget
func NewFormItem(text string, widget fyne.CanvasObject) *FormItem {
	return &FormItem{Text: text, Widget: widget}
}

var _ fyne.Validatable = (*Form)(nil)

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
	OnSubmit   func() `json:"-"`
	OnCancel   func() `json:"-"`
	SubmitText string
	CancelText string

	// Orientation allows a form to be vertical (a single column), horizontal (default, label then input)
	// or to adapt according to the orientation of the mobile device (adaptive).
	//
	// Since: 2.5
	Orientation Orientation

	itemGrid     *fyne.Container
	buttonBox    *fyne.Container
	cancelButton *Button
	submitButton *Button

	disabled bool

	onValidationChanged func(error)
	validationError     error
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
		f.itemGrid.Add(f.createInput(item))
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
	f.ExtendBaseWidget(f)
	cache.Renderer(f.super()) // we are about to make changes to renderer created content... not great!
	f.ensureRenderItems()
	f.updateButtons()
	f.updateLabels()

	if f.isVertical() {
		f.itemGrid.Layout = layout.NewVBoxLayout()
	} else {
		f.itemGrid.Layout = layout.NewFormLayout()
	}

	f.BaseWidget.Refresh()
	canvas.Refresh(f.super()) // refresh ourselves for BG color - the above updates the content
}

// Enable enables submitting this form.
//
// Since: 2.1
func (f *Form) Enable() {
	f.disabled = false
	f.cancelButton.Enable()
	f.checkValidation(nil) // as the form may be invalid
}

// Disable disables submitting this form.
//
// Since: 2.1
func (f *Form) Disable() {
	f.disabled = true
	f.submitButton.Disable()
	f.cancelButton.Disable()
}

// Disabled returns whether submitting the form is disabled.
// Note that, if the form fails validation, the submit button may be
// disabled even if this method returns true.
//
// Since: 2.1
func (f *Form) Disabled() bool {
	return f.disabled
}

// SetOnValidationChanged is intended for parent widgets or containers to hook into the validation.
// The function might be overwritten by a parent that cares about child validation (e.g. widget.Form)
func (f *Form) SetOnValidationChanged(callback func(error)) {
	f.onValidationChanged = callback
}

// Validate validates the entire form and returns the first error that is encountered.
func (f *Form) Validate() error {
	for _, item := range f.Items {
		if w, ok := item.Widget.(fyne.Validatable); ok {
			if err := w.Validate(); err != nil {
				return err
			}
		}
	}
	return nil
}

func (f *Form) createInput(item *FormItem) fyne.CanvasObject {
	_, ok := item.Widget.(fyne.Validatable)
	if item.HintText == "" {
		if !ok {
			return item.Widget
		}
		if !f.itemWidgetHasValidator(item.Widget) { // we don't have validation
			return item.Widget
		}
	}

	th := f.Theme()
	v := fyne.CurrentApp().Settings().ThemeVariant()

	text := canvas.NewText(item.HintText, th.Color(theme.ColorNamePlaceHolder, v))
	text.TextSize = th.Size(theme.SizeNameCaptionText)
	item.helperOutput = text
	f.updateHelperText(item)
	textContainer := &fyne.Container{Objects: []fyne.CanvasObject{text}}
	return &fyne.Container{Layout: formItemLayout{form: f}, Objects: []fyne.CanvasObject{item.Widget, textContainer}}
}

func (f *Form) itemWidgetHasValidator(w fyne.CanvasObject) bool {
	value := reflect.ValueOf(w).Elem()
	validatorField := value.FieldByName("Validator")
	if validatorField == (reflect.Value{}) {
		return false
	}
	validator, ok := validatorField.Interface().(fyne.StringValidator)
	if !ok {
		return false
	}
	return validator != nil
}

func (f *Form) createLabel(text string) fyne.CanvasObject {
	th := f.Theme()
	v := fyne.CurrentApp().Settings().ThemeVariant()
	label := &canvas.Text{
		Text:      text,
		Alignment: fyne.TextAlignTrailing,
		Color:     th.Color(theme.ColorNameForeground, v),
		TextSize:  th.Size(theme.SizeNameText),
		TextStyle: fyne.TextStyle{Bold: true},
	}
	if f.isVertical() {
		label.Alignment = fyne.TextAlignLeading
	}

	return &fyne.Container{Layout: &formLabelLayout{form: f}, Objects: []fyne.CanvasObject{label}}
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
		if item.invalid {
			f.submitButton.Disable()
			return
		}
	}

	if !f.disabled {
		f.submitButton.Enable()
	}
}

func (f *Form) ensureRenderItems() {
	done := len(f.itemGrid.Objects) / 2
	if done >= len(f.Items) {
		f.itemGrid.Objects = f.itemGrid.Objects[0 : len(f.Items)*2]
		return
	}

	adding := len(f.Items) - done
	objects := make([]fyne.CanvasObject, adding*2)
	off := 0
	for i, item := range f.Items {
		if i < done {
			continue
		}

		objects[off] = f.createLabel(item.Text)
		off++
		f.setUpValidation(item.Widget, i)
		objects[off] = f.createInput(item)
		off++
	}
	f.itemGrid.Objects = append(f.itemGrid.Objects, objects...)
}

func (f *Form) isVertical() bool {
	if f.Orientation == Vertical {
		return true
	} else if f.Orientation == Horizontal {
		return false
	}

	dev := fyne.CurrentDevice()
	if dev.IsMobile() {
		orient := dev.Orientation()
		return orient == fyne.OrientationVertical || orient == fyne.OrientationVerticalUpsideDown
	}

	return false
}

func (f *Form) setUpValidation(widget fyne.CanvasObject, i int) {
	updateValidation := func(err error) {
		if err == errFormItemInitialState {
			return
		}
		f.Items[i].validationError = err
		f.Items[i].invalid = err != nil
		f.setValidationError(err)
		f.checkValidation(err)
		f.updateHelperText(f.Items[i])
	}
	if w, ok := widget.(fyne.Validatable); ok {
		f.Items[i].invalid = w.Validate() != nil
		if e, ok := w.(*Entry); ok {
			e.onFocusChanged = func(bool) {
				updateValidation(e.validationError)
			}
			if e.Validator != nil && f.Items[i].invalid {
				// set initial state error to guarantee next error (if triggers) is always different
				e.SetValidationError(errFormItemInitialState)
			}
		}
		w.SetOnValidationChanged(updateValidation)
	}
}

func (f *Form) setValidationError(err error) {
	if err == nil && f.validationError == nil {
		return
	}

	if !errors.Is(err, f.validationError) {
		if err == nil {
			for _, item := range f.Items {
				if item.invalid {
					err = item.validationError
					break
				}
			}
		}
		f.validationError = err

		if f.onValidationChanged != nil {
			f.onValidationChanged(err)
		}
	}
}

func (f *Form) updateHelperText(item *FormItem) {
	th := f.Theme()
	v := fyne.CurrentApp().Settings().ThemeVariant()

	if item.helperOutput == nil {
		return // testing probably, either way not rendered yet
	}
	showHintIfError := false
	if e, ok := item.Widget.(*Entry); ok {
		if !e.dirty || (e.focused && !item.wasFocused) {
			showHintIfError = true
		}
		if e.dirty && !e.focused {
			item.wasFocused = true
		}
	}

	if item.validationError == nil || showHintIfError {
		item.helperOutput.Text = item.HintText
		item.helperOutput.Color = th.Color(theme.ColorNamePlaceHolder, v)
	} else {
		item.helperOutput.Text = item.validationError.Error()
		item.helperOutput.Color = th.Color(theme.ColorNameError, v)
	}

	if item.helperOutput.Text == "" {
		item.helperOutput.Hide()
	} else {
		item.helperOutput.Show()
	}
	item.helperOutput.Refresh()
}

func (f *Form) updateLabels() {
	th := f.Theme()
	v := fyne.CurrentApp().Settings().ThemeVariant()

	for i, item := range f.Items {
		l := f.itemGrid.Objects[i*2].(*fyne.Container).Objects[0].(*canvas.Text)
		l.TextSize = th.Size(theme.SizeNameText)
		if dis, ok := item.Widget.(fyne.Disableable); ok {
			if dis.Disabled() {
				l.Color = th.Color(theme.ColorNameDisabled, v)
			} else {
				l.Color = th.Color(theme.ColorNameForeground, v)
			}
		} else {
			l.Color = th.Color(theme.ColorNameForeground, v)
		}

		l.Text = item.Text
		if f.isVertical() {
			l.Alignment = fyne.TextAlignLeading
		} else {
			l.Alignment = fyne.TextAlignTrailing
		}
		l.Refresh()
		f.updateHelperText(item)
	}
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
func (f *Form) CreateRenderer() fyne.WidgetRenderer {
	f.ExtendBaseWidget(f)
	th := f.Theme()
	f.cancelButton = &Button{Icon: th.Icon(theme.IconNameCancel), OnTapped: f.OnCancel}
	f.submitButton = &Button{Icon: th.Icon(theme.IconNameConfirm), OnTapped: f.OnSubmit, Importance: HighImportance}
	buttons := &fyne.Container{Layout: layout.NewGridLayoutWithRows(1), Objects: []fyne.CanvasObject{f.cancelButton, f.submitButton}}
	f.buttonBox = &fyne.Container{Layout: layout.NewBorderLayout(nil, nil, nil, buttons), Objects: []fyne.CanvasObject{buttons}}
	f.validationError = errFormItemInitialState // set initial state error to guarantee next error (if triggers) is always different

	f.itemGrid = &fyne.Container{Layout: layout.NewFormLayout()}
	if f.isVertical() {
		f.itemGrid.Layout = layout.NewVBoxLayout()
	} else {
		f.itemGrid.Layout = layout.NewFormLayout()
	}
	content := &fyne.Container{Layout: layout.NewVBoxLayout(), Objects: []fyne.CanvasObject{f.itemGrid, f.buttonBox}}
	renderer := NewSimpleRenderer(content)
	f.ensureRenderItems()
	f.updateButtons()
	f.updateLabels()
	f.checkValidation(nil) // will trigger a validation check for correct initial validation status
	return renderer
}

// NewForm creates a new form widget with the specified rows of form items
// and (if any of them should be shown) a form controls row at the bottom
func NewForm(items ...*FormItem) *Form {
	form := &Form{Items: items}
	form.ExtendBaseWidget(form)

	return form
}

type formLabelLayout struct {
	form *Form
}

func (f formLabelLayout) Layout(objs []fyne.CanvasObject, size fyne.Size) {
	innerPad := f.form.Theme().Size(theme.SizeNameInnerPadding)
	xPad := innerPad
	yPos := float32(0)
	if !f.form.isVertical() {
		xPad += innerPad
		yPos = innerPad
	}
	objs[0].Move(fyne.NewPos(innerPad, yPos))
	objs[0].Resize(fyne.NewSize(size.Width-xPad, objs[0].MinSize().Height))
}

func (f formLabelLayout) MinSize(objs []fyne.CanvasObject) fyne.Size {
	innerPad := f.form.Theme().Size(theme.SizeNameInnerPadding)
	min0 := objs[0].MinSize()

	if !f.form.isVertical() {
		min0 = min0.AddWidthHeight(innerPad, 0)
	}

	return min0.AddWidthHeight(innerPad, 0)
}

type formItemLayout struct {
	form *Form
}

func (f formItemLayout) Layout(objs []fyne.CanvasObject, size fyne.Size) {
	innerPad := f.form.Theme().Size(theme.SizeNameInnerPadding)
	itemHeight := objs[0].MinSize().Height
	objs[0].Resize(fyne.NewSize(size.Width, itemHeight))

	objs[1].Move(fyne.NewPos(innerPad, itemHeight+innerPad/2))
	objs[1].Resize(fyne.NewSize(size.Width, objs[1].MinSize().Width))
}

func (f formItemLayout) MinSize(objs []fyne.CanvasObject) fyne.Size {
	innerPad := f.form.Theme().Size(theme.SizeNameInnerPadding)
	min0 := objs[0].MinSize()
	min1 := objs[1].MinSize()

	minWidth := fyne.Max(min0.Width, min1.Width)
	height := min0.Height

	items := objs[1].(*fyne.Container).Objects
	if len(items) > 0 && items[0].(*canvas.Text).Text != "" {
		height += min1.Height + innerPad
	}
	return fyne.NewSize(minWidth, height)
}
