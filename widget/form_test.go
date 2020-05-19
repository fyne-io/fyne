package widget

import (
	"testing"

	"fyne.io/fyne/test"
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
