package widget

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"

	"github.com/stretchr/testify/assert"
)

type extendEntry struct {
	Entry
}

func TestEntry_Password_Extended_CreateRenderer(t *testing.T) {
	a := test.NewApp()
	defer test.NewApp()
	w := a.NewWindow("")
	entry := &extendEntry{}
	entry.ExtendBaseWidget(entry)
	entry.Password = true
	entry.Wrapping = fyne.TextTruncate
	assert.NotNil(t, test.WidgetRenderer(entry))
	r := test.WidgetRenderer(entry).(*entryRenderer).scroll.Content.(*entryContent)
	p := test.WidgetRenderer(r).(*entryContentRenderer).provider
	texts := test.WidgetRenderer(p).(*textRenderer).texts

	w.SetContent(entry)

	test.Type(entry, "Pass")
	assert.Equal(t, passwordChar+passwordChar+passwordChar+passwordChar, texts[0].Text)
	assert.NotNil(t, entry.ActionItem)
	test.Tap(entry.ActionItem.(*passwordRevealer))
	assert.Equal(t, "Pass", texts[0].Text)
	assert.Equal(t, entry, w.Canvas().Focused())
}
