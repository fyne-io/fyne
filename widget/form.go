package widget

import (
	"github.com/fyne-io/fyne"
	"github.com/fyne-io/fyne/layout"
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

func (f *Form) createLabel(text string) *Label {
	label := &Label{Text: text,
		Alignment: fyne.TextAlignTrailing,
		TextStyle: fyne.TextStyle{Bold: true},
	}
	return label
}

func (f *Form) ensureGrid() {
	if f.itemGrid != nil {
		return
	}

	f.itemGrid = fyne.NewContainerWithLayout(layout.NewGridLayout(2), []fyne.CanvasObject{}...)
}

// Append adds a new row to the form, using the text as a label next to the specified Widget
func (f *Form) Append(text string, widget fyne.CanvasObject) {
	item := &FormItem{Text: text, Widget: widget}
	f.AppendItem(item)
}

// AppendItem adds the specified row to the end of the Form
func (f *Form) AppendItem(item *FormItem) {
	f.Items = append(f.Items, item)

	f.ensureGrid()
	f.itemGrid.AddObject(f.createLabel(item.Text))
	f.itemGrid.AddObject(item.Widget)

	f.Renderer().Refresh()
}

// Renderer is a private method to Fyne which links this widget to it's renderer
func (f *Form) Renderer() fyne.WidgetRenderer {
	if f.renderer == nil {
		f.ensureGrid()

		buttons := NewHBox(layout.NewSpacer())
		if f.OnCancel != nil {
			buttons.Append(NewButton("Cancel", f.OnCancel))
		}
		if f.OnSubmit != nil {
			submit := NewButton("Submit", f.OnSubmit)
			submit.Style = PrimaryButton

			buttons.Append(submit)
		}
		f.renderer = NewVBox(f.itemGrid, buttons).Renderer()
	}

	return f.renderer
}

// NewForm creates a new form widget with the specified rows of form items
// and (if any of them should be shown) a form controls row at the bottom
func NewForm(items ...*FormItem) *Form {
	form := &Form{baseWidget: baseWidget{}, Items: items}

	form.Renderer().Layout(form.MinSize())
	return form
}
