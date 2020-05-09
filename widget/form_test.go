package widget

import (
	"testing"
	"time"

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
	form := NewForm().
		OnSubmit(func() {}).
		OnCancel(func() {})
	form.Append("test", NewEntry())
	assert.Equal(t, "Submit", form.SubmitText)
	assert.Equal(t, "Cancel", form.CancelText)

	form = NewForm().
		OnSubmitWithLabel("Apply", func() {}).
		OnCancelWithLabel("Close", func() {})
	assert.Equal(t, "Apply", form.SubmitText)
	assert.Equal(t, "Close", form.CancelText)
}

func TestForm_AddRemoveButton(t *testing.T) {
	scount := 0
	ccount := 0
	sscount := 10
	tapped := make(chan bool)
	form := NewForm().
		Append("test", NewEntry()).
		OnSubmit(func() { scount++; tapped <- true }).
		OnCancel(func() { ccount++; tapped <- true })

	go test.Tap(form.submitButton)
	func() {
		select {
		case <-tapped:
		case <-time.After(1 * time.Second):
			assert.Fail(t, "Timed out waiting for submit button tap")
		}
	}()
	assert.Equal(t, 1, scount)
	go test.Tap(form.cancelButton)
	func() {
		select {
		case <-tapped:
		case <-time.After(1 * time.Second):
			assert.Fail(t, "Timed out waiting for cancel button tap")
		}
	}()
	assert.Equal(t, 1, ccount)

	form.OnSubmit(func() {
		sscount++
		tapped <- true
	})
	go test.Tap(form.submitButton)
	func() {
		select {
		case <-tapped:
		case <-time.After(1 * time.Second):
			assert.Fail(t, "Timed out waiting for cancel button tap")
		}
	}()
	assert.Equal(t, 11, sscount)
}
