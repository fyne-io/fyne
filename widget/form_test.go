package widget

import (
	"errors"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
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

func TestForm_Append_Items(t *testing.T) {
	form := &Form{Items: []*FormItem{{Text: "test1", Widget: NewEntry()}}}
	assert.Equal(t, 1, len(form.Items))
	renderer := test.WidgetRenderer(form)

	form.Items = append(form.Items, NewFormItem("test2", NewEntry()))
	assert.True(t, len(form.Items) == 2)

	form.Refresh()
	c := renderer.Objects()[0].(*fyne.Container).Objects[0].(*fyne.Container)
	assert.Equal(t, "test2", c.Objects[2].(*canvas.Text).Text)
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
	assert.Equal(t, "Test", c.Objects[0].(*canvas.Text).Text)

	item.Text = "Changed"
	form.Refresh()
	assert.Equal(t, "Changed", c.Objects[0].(*canvas.Text).Text)
}

func TestForm_ChangeTheme(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	form := &Form{
		Items: []*FormItem{
			{Text: "test1", Widget: NewEntry()},
			{Text: "test2", Widget: NewLabel("static")},
		},
		OnSubmit: func() {}, OnCancel: func() {}}
	w := test.NewWindow(form)
	defer w.Close()

	test.AssertImageMatches(t, "form/theme_initial.png", w.Canvas().Capture())

	test.WithTestTheme(t, func() {
		form.Refresh()
		w.Resize(form.MinSize().Add(fyne.NewSize(theme.InnerPadding(), theme.InnerPadding())))
		test.AssertImageMatches(t, "form/theme_changed.png", w.Canvas().Capture())
	})
}

func TestForm_Disabled(t *testing.T) {
	test.NewApp()
	defer test.NewApp()
	test.ApplyTheme(t, theme.LightTheme())

	disabled := NewEntry()
	disabled.Disable()
	f := NewForm(
		NewFormItem("Form Item 1", NewEntry()),
		NewFormItem("Form Item 2", disabled))

	w := test.NewWindow(f)
	defer w.Close()

	test.AssertImageMatches(t, "form/disabled.png", w.Canvas().Capture())
}

func TestForm_Hints(t *testing.T) {
	test.NewApp()
	defer test.NewApp()
	test.ApplyTheme(t, theme.LightTheme())

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
	test.NewApp()
	defer test.NewApp()
	test.ApplyTheme(t, theme.LightTheme())

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

func TestForm_EntryValidation_FirstTypeValid(t *testing.T) {
	test.NewApp()
	defer test.NewApp()
	test.ApplyTheme(t, theme.LightTheme())

	notEmptyValidator := func(s string) error {
		if s == "" {
			return errors.New("can't be empty")
		}
		return nil
	}

	entry1 := &Entry{Validator: notEmptyValidator, Text: ""}
	entry2 := &Entry{Validator: notEmptyValidator, Text: ""}
	items := []*FormItem{
		{Text: "First", Widget: entry1},
		{Text: "Second", Widget: entry2},
	}

	form := &Form{Items: items, OnSubmit: func() {}, OnCancel: func() {}}
	w := test.NewWindow(form)
	defer w.Close()

	test.AssertImageMatches(t, "form/validation_entry_first_type_initial.png", w.Canvas().Capture())

	test.Type(entry1, "H")
	test.Type(entry2, "L")
	entry1.focused = false
	entry1.Refresh()
	w = test.NewWindow(form)

	test.AssertImageMatches(t, "form/validation_entry_first_type_valid.png", w.Canvas().Capture())

	entry1.SetText("")
	entry2.SetText("")
	w = test.NewWindow(form)

	test.AssertImageMatches(t, "form/validation_entry_first_type_invalid.png", w.Canvas().Capture())
}

func TestForm_DisableEnable(t *testing.T) {
	test.NewApp()
	defer test.NewApp()
	test.ApplyTheme(t, theme.LightTheme())

	form := &Form{
		Items: []*FormItem{
			{Text: "test1", Widget: NewEntry()},
		},
		OnSubmit: func() {}, OnCancel: func() {}}
	w := test.NewWindow(form)
	defer w.Close()

	if form.Disabled() {
		t.Error("form.Disabled() returned true when it should have been false")
	}

	test.AssertImageMatches(t, "form/disable_initial.png", w.Canvas().Capture())

	form.Disable()

	if !form.Disabled() {
		t.Error("form.Disabled() returned false when it should have been true")
	}

	test.AssertImageMatches(t, "form/disable_disabled.png", w.Canvas().Capture())

	form.Enable()

	if form.Disabled() {
		t.Error("form.Disabled() returned true when it should have been false")
	}

	test.AssertImageMatches(t, "form/disable_re_enabled.png", w.Canvas().Capture())
}

func TestForm_Disable_Validation(t *testing.T) {
	test.NewApp()
	defer test.NewApp()
	test.ApplyTheme(t, theme.LightTheme())

	entry := &Entry{Validator: validation.NewRegexp(`^\d{2}-\w{4}$`, "Input is not valid"), Text: "wrong"}

	form := &Form{Items: []*FormItem{{Text: "test", Widget: entry}}, OnSubmit: func() {}, OnCancel: func() {}}
	w := test.NewWindow(form)
	defer w.Close()

	test.AssertImageMatches(t, "form/disable_validation_initial.png", w.Canvas().Capture())

	form.Disable()

	test.AssertImageMatches(t, "form/disable_validation_disabled_invalid.png", w.Canvas().Capture())

	form.Enable()

	test.AssertImageMatches(t, "form/disable_validation_enabled_invalid.png", w.Canvas().Capture())

	entry.SetText("15-true")
	test.AssertImageMatches(t, "form/disable_validation_enabled_valid.png", w.Canvas().Capture())

	// ensure we don't re-enable the form when entering something valid
	entry.SetText("invalid")
	form.Disable()
	entry.SetText("15-true")

	test.AssertImageMatches(t, "form/disable_validation_disabled_valid.png", w.Canvas().Capture())
}

func TestForm_HintsRendered(t *testing.T) {
	test.NewApp()
	defer test.NewApp()
	test.ApplyTheme(t, theme.LightTheme())

	f := NewForm()

	fi1 := NewFormItem("Form Item 1", NewEntry())
	fi1.HintText = "HT1"
	f.AppendItem(fi1)

	fi2 := NewFormItem("Form Item 2", NewEntry())
	fi2.HintText = "HT2"

	f.AppendItem(fi2)

	fi3 := NewFormItem("Form Item 3", NewEntry())
	fi3.HintText = "HT3"

	f.AppendItem(fi3)

	w := test.NewWindow(f)
	defer w.Close()

	test.AssertImageMatches(t, "form/hints_rendered.png", w.Canvas().Capture())
}

func TestForm_Validate(t *testing.T) {
	entry1 := &Entry{Validator: validation.NewRegexp(`^\d{2}-\w{4}$`, "Input is not valid 1"), Text: "15-true"}
	entry2 := &Entry{Validator: validation.NewRegexp(`^\w{3}-\w{5}$`, "Input is not valid 2"), Text: "wrong"}

	form := &Form{
		Items: []*FormItem{
			{Text: "First", Widget: entry1},
			{Text: "Second", Widget: entry2},
		},
	}

	err := form.Validate()
	if assert.Error(t, err) {
		assert.Equal(t, "Input is not valid 2", err.Error())
	}

	entry1.SetText("incorrect")
	err = form.Validate()
	if assert.Error(t, err) {
		assert.Equal(t, "Input is not valid 1", err.Error())
	}

	entry1.SetText("15-true")
	err = form.Validate()
	if assert.Error(t, err) {
		assert.Equal(t, "Input is not valid 2", err.Error())
	}

	entry2.SetText("not-wrong")
	err = form.Validate()
	assert.NoError(t, err)

}

func TestForm_SetOnValidationChanged(t *testing.T) {
	entry1 := &Entry{Validator: validation.NewRegexp(`^\d{2}-\w{4}$`, "Input is not valid"), Text: "15-true"}

	form := &Form{
		Items: []*FormItem{
			{Text: "First", Widget: entry1},
		},
	}

	validationError := false

	form.SetOnValidationChanged(func(err error) {
		validationError = err != nil
	})

	form.CreateRenderer()

	entry1.SetText("incorrect")
	assert.Error(t, form.Validate())
	assert.True(t, validationError)

	entry1.SetText("15-true")
	assert.NoError(t, form.Validate())
	assert.False(t, validationError)

}

func TestForm_ExtendedEntry(t *testing.T) {
	extendedEntry := NewSelectEntry([]string{""})

	test.NewApp()
	defer test.NewApp()

	form := &Form{
		Items: []*FormItem{
			{Text: "Extended entry", Widget: extendedEntry},
		},
	}
	w := test.NewWindow(form)
	defer w.Close()

	test.AssertRendersToMarkup(t, "form/extended_entry.xml", w.Canvas())
}

func TestForm_RefreshFromStructInit(t *testing.T) {
	form := &Form{
		Items: []*FormItem{
			{Text: "Entry", Widget: NewEntry()},
		},
	}

	assert.NotPanics(t, func() {
		form.Refresh()
	})

}
