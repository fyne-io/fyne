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
	assert.NotNil(t, test.WidgetRenderer(form))
	assert.Equal(t, 2, len(form.itemGrid.Objects))

	form.Append("test2", NewEntry())
	assert.Equal(t, 4, len(form.itemGrid.Objects))
}

func TestForm_Extended_Append(t *testing.T) {
	form := &extendedForm{}
	form.ExtendBaseWidget(form)
	form.Items = []*FormItem{{Text: "test1", Widget: NewEntry()}}
	assert.Equal(t, 1, len(form.Items))

	form.Append("test2", NewEntry())
	assert.True(t, len(form.Items) == 2)

	item := &FormItem{Text: "test3", Widget: NewEntry()}
	form.AppendItem(item)
	assert.True(t, len(form.Items) == 3)
	assert.Equal(t, item, form.Items[2])
}
