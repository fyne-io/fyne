package widget

import (
	"testing"

	"fyne.io/fyne/v2/test"

	"github.com/stretchr/testify/assert"
)

type extendedForm struct {
	Form
}

func TestForm_Extended_CreateRenderer(t *testing.T) {
	form := &extendedForm{}
	form.ExtendBaseWidget(form)
	form.Items = []*FormItem{{Text: "test1", Widget: NewEntry()}}
	assert.NotNil(t, test.TempWidgetRenderer(t, form))
	assert.Len(t, form.itemGrid.Objects, 2)

	form.Append("test2", NewEntry())
	assert.Len(t, form.itemGrid.Objects, 4)
}

func TestForm_Extended_Append(t *testing.T) {
	form := &extendedForm{}
	form.ExtendBaseWidget(form)
	form.Items = []*FormItem{{Text: "test1", Widget: NewEntry()}}
	assert.Len(t, form.Items, 1)

	form.Append("test2", NewEntry())
	assert.Len(t, form.Items, 2)

	item := &FormItem{Text: "test3", Widget: NewEntry()}
	form.AppendItem(item)
	assert.Len(t, form.Items, 3)
	assert.Equal(t, item, form.Items[2])
}
