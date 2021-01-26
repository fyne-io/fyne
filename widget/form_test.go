package widget

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"

	"github.com/stretchr/testify/assert"
)

func TestFormSize(t *testing.T) {
	form := &Form{Items: []*FormItem{
		{Text: "test1", Widget: NewEntry()},
		{Text: "test2", Widget: NewEntry()},
	}}

	assert.Equal(t, 2, len(form.Items))
}

func TestForm_CreateRenderer(t *testing.T) {
	form := &Form{Items: []*FormItem{{Text: "test1", Widget: NewEntry()}}}
	assert.NotNil(t, test.WidgetRenderer(form))
	assert.Equal(t, 2, len(form.itemGrid.Objects))

	form.Append("test2", NewEntry())
	assert.Equal(t, 4, len(form.itemGrid.Objects))
}

func TestForm_Append(t *testing.T) {
	form := &Form{Items: []*FormItem{{Text: "test1", Widget: NewEntry()}}}
	assert.Equal(t, 1, len(form.Items))

	form.Append("test2", NewEntry())
	assert.True(t, len(form.Items) == 2)

	item := &FormItem{Text: "test3", Widget: NewEntry()}
	form.AppendItem(item)
	assert.True(t, len(form.Items) == 3)
	assert.Equal(t, item, form.Items[2])
}

func TestForm_CustomButtonsText(t *testing.T) {
	form := &Form{OnSubmit: func() {}, OnCancel: func() {}}
	form.Append("test", NewEntry())
	assert.Equal(t, "Submit", form.SubmitText)
	assert.Equal(t, "Cancel", form.CancelText)

	form = &Form{OnSubmit: func() {}, SubmitText: "Apply",
		OnCancel: func() {}, CancelText: "Close"}
	assert.Equal(t, "Apply", form.SubmitText)
	assert.Equal(t, "Close", form.CancelText)
}

func TestForm_AddRemoveButton(t *testing.T) {
	scount := 0
	ccount := 0
	sscount := 10
	form := &Form{OnSubmit: func() {}, OnCancel: func() {}}
	form.Append("test", NewEntry())
	form.OnSubmit = func() { scount++ }
	form.OnCancel = func() { ccount++ }
	form.Refresh()

	test.Tap(form.submitButton)
	assert.Equal(t, 1, scount, "tapping submit should incr scount")

	test.Tap(form.cancelButton)
	assert.Equal(t, 1, ccount, "tapping cancel should incr ccount")

	form.OnSubmit = func() { sscount++ }
	form.Refresh()
	test.Tap(form.submitButton)
	assert.Equal(t, 11, sscount, "tapping new submit should incr sscount from 10 to 11")

	form.OnCancel = func() { sscount = sscount - 6 }
	form.Refresh()
	test.Tap(form.cancelButton)
	assert.Equal(t, 5, sscount, "tapping new cancel should decr ssount from 11 down to 5")
}

func TestForm_Renderer(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	form := &Form{
		Items: []*FormItem{
			{Text: "test1", Widget: NewEntry()},
			{Text: "test2", Widget: NewEntry()},
		},
		OnSubmit: func() {}, OnCancel: func() {}}
	w := test.NewWindow(form)
	defer w.Close()

	test.AssertRendersToMarkup(t, "form/layout.xml", w.Canvas())
}

func TestForm_ChangeText(t *testing.T) {
	item := NewFormItem("Test", NewEntry())
	form := NewForm(item)

	renderer := test.WidgetRenderer(form)
	c := renderer.Objects()[0].(*fyne.Container).Objects[0].(*fyne.Container)
	assert.Equal(t, "Test", c.Objects[0].(*Label).Text)

	item.Text = "Changed"
	form.Refresh()
	assert.Equal(t, "Changed", c.Objects[0].(*Label).Text)
}

func TestForm_ChangeTheme(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	form := &Form{
		Items: []*FormItem{
			{Text: "test1", Widget: NewEntry()},
			{Text: "test2", Widget: NewEntry()},
		},
		OnSubmit: func() {}, OnCancel: func() {}}
	w := test.NewWindow(form)
	defer w.Close()

	test.AssertImageMatches(t, "form/theme_initial.png", w.Canvas().Capture())

	test.WithTestTheme(t, func() {
		form.Refresh()
		w.Resize(form.MinSize().Add(fyne.NewSize(theme.Padding()*2, theme.Padding()*2)))
		test.AssertImageMatches(t, "form/theme_changed.png", w.Canvas().Capture())
	})
}

func TestForm_Hints(t *testing.T) {
	app := test.NewApp()
	defer test.NewApp()
	app.Settings().SetTheme(theme.LightTheme())

	entry1 := &Entry{}
	entry2 := &Entry{Validator: validation.NewRegexp(`^\w{3}-\w{5}$`, "Input is not valid"), Text: "wrong"}
	items := []*FormItem{
		{Text: "First", Widget: entry1, HintText: "An entry hint"},
		{Text: "Second", Widget: entry2},
	}

	form := &Form{Items: items, OnSubmit: func() {}, OnCancel: func() {}}
	w := test.NewWindow(form)
	defer w.Close()

	test.AssertImageMatches(t, "form/hint_initial.png", w.Canvas().Capture())

	test.Type(entry2, "n")
	test.AssertImageMatches(t, "form/hint_invalid.png", w.Canvas().Capture())

	test.Type(entry2, "ot-")
	test.AssertImageMatches(t, "form/hint_valid.png", w.Canvas().Capture())
}

func TestForm_Validation(t *testing.T) {
	app := test.NewApp()
	defer test.NewApp()
	app.Settings().SetTheme(theme.LightTheme())

	entry1 := &Entry{Validator: validation.NewRegexp(`^\d{2}-\w{4}$`, "Input is not valid"), Text: "15-true"}
	entry2 := &Entry{Validator: validation.NewRegexp(`^\w{3}-\w{5}$`, "Input is not valid"), Text: "wrong"}
	entry3 := &Entry{}
	items := []*FormItem{
		{Text: "First", Widget: entry1},
		{Text: "Second", Widget: entry2},
		{Text: "Third", Widget: entry3},
	}

	form := &Form{Items: items, OnSubmit: func() {}, OnCancel: func() {}}
	w := test.NewWindow(form)
	defer w.Close()

	test.AssertImageMatches(t, "form/validation_initial.png", w.Canvas().Capture())

	test.Type(entry2, "not-")
	entry1.SetText("incorrect")
	w = test.NewWindow(form)

	test.AssertImageMatches(t, "form/validation_invalid.png", w.Canvas().Capture())

	entry1.SetText("15-true")
	w = test.NewWindow(form)

	test.AssertImageMatches(t, "form/validation_valid.png", w.Canvas().Capture())
}
