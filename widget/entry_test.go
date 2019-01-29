package widget

import (
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/test"
	"fyne.io/fyne/theme"

	"github.com/stretchr/testify/assert"
)

func entryRenderTexts(e *Entry) []*canvas.Text {
	textWid := Renderer(e).(*entryRenderer).text
	return Renderer(textWid).(*textRenderer).texts
}

func TestEntry_MinSize(t *testing.T) {
	entry := NewEntry()
	min := entry.MinSize()
	entry.SetPlaceHolder("")
	assert.Equal(t, min, entry.MinSize())
	entry.SetText("")
	assert.Equal(t, min, entry.MinSize())
	entry.SetPlaceHolder("Hi")
	assert.True(t, entry.MinSize().Width > min.Width)
	assert.Equal(t, entry.MinSize().Height, min.Height)

	assert.True(t, min.Width > theme.Padding()*2)
	assert.True(t, min.Height > theme.Padding()*2)
}

func TestMultiLineEntry_MinSize(t *testing.T) {
	entry := NewEntry()
	entry.MinSize()
	singleMin := entry.MinSize()

	multi := NewMultiLineEntry()
	multiMin := multi.MinSize()

	assert.Equal(t, singleMin.Width, multiMin.Width)
	assert.True(t, multiMin.Height > singleMin.Height)

	multi.MultiLine = false
	multiMin = multi.MinSize()
	assert.Equal(t, singleMin.Height, multiMin.Height)
}

func TestEntry_SetPlaceHolder(t *testing.T) {
	entry := NewEntry()

	assert.Equal(t, 0, len(entry.Text))
	assert.Equal(t, 0, entry.textProvider().len())

	entry.SetPlaceHolder("Test")
	assert.Equal(t, 0, len(entry.Text))
	assert.Equal(t, 0, entry.textProvider().len())
	assert.Equal(t, 4, entry.placeholderProvider().len())
	assert.False(t, entry.placeholderProvider().Hidden)

	entry.SetText("Hi")
	assert.Equal(t, 2, len(entry.Text))
	assert.True(t, entry.placeholderProvider().Hidden)

	assert.Equal(t, 2, entry.textProvider().len())
}

func TestEntry_OnKeyDown(t *testing.T) {
	entry := NewEntry()

	key := new(fyne.KeyEvent)
	key.String = "H"
	entry.OnKeyDown(key)
	key.String = "i"
	entry.OnKeyDown(key)

	assert.Equal(t, "Hi", entry.Text)
}

func TestEntry_SetReadOnly_KeyDown(t *testing.T) {
	entry := NewEntry()

	key := new(fyne.KeyEvent)
	key.String = "H"
	entry.OnKeyDown(key)
	entry.SetReadOnly(true)
	key.String = "i"
	entry.OnKeyDown(key)
	assert.Equal(t, "H", entry.Text)

	entry.SetReadOnly(false)
	key.String = "i"
	entry.OnKeyDown(key)
	assert.Equal(t, "Hi", entry.Text)
}

func TestEntry_SetReadOnly_OnFocus(t *testing.T) {
	entry := NewEntry()
	entry.SetReadOnly(true)

	entry.OnFocusGained()
	assert.False(t, entry.Focused())

	entry.SetReadOnly(false)
	entry.OnFocusGained()
	assert.True(t, entry.Focused())
}

func TestEntry_OnKeyDown_Insert(t *testing.T) {
	entry := NewEntry()

	key := new(fyne.KeyEvent)
	key.String = "H"
	entry.OnKeyDown(key)
	key.String = "i"
	entry.OnKeyDown(key)
	assert.Equal(t, "Hi", entry.Text)

	left := &fyne.KeyEvent{Name: fyne.KeyLeft}
	entry.OnKeyDown(left)

	key.String = "o"
	entry.OnKeyDown(key)
	assert.Equal(t, "Hoi", entry.Text)
}

func TestEntry_OnKeyDown_Newline(t *testing.T) {
	entry := &Entry{MultiLine: true}
	entry.SetText("Hi")
	assert.Equal(t, 0, entry.CursorRow)
	assert.Equal(t, 0, entry.CursorColumn)

	right := &fyne.KeyEvent{Name: fyne.KeyRight}
	entry.OnKeyDown(right)
	assert.Equal(t, 0, entry.CursorRow)
	assert.Equal(t, 1, entry.CursorColumn)

	key := &fyne.KeyEvent{Name: fyne.KeyReturn}
	entry.OnKeyDown(key)

	assert.Equal(t, "H\ni", entry.Text)
	assert.Equal(t, 1, entry.CursorRow)
	assert.Equal(t, 0, entry.CursorColumn)

	key = new(fyne.KeyEvent)
	key.String = "o"
	entry.OnKeyDown(key)
	assert.Equal(t, "H\noi", entry.textProvider().String())
	assert.Equal(t, "H", entryRenderTexts(entry)[0].Text)
	assert.Equal(t, "oi", entryRenderTexts(entry)[1].Text)
}

func TestEntry_OnKeyDown_Backspace(t *testing.T) {
	entry := NewEntry()
	entry.SetText("Hi")
	right := &fyne.KeyEvent{Name: fyne.KeyRight}
	entry.OnKeyDown(right)
	entry.OnKeyDown(right)
	assert.Equal(t, 0, entry.CursorRow)
	assert.Equal(t, 2, entry.CursorColumn)

	key := &fyne.KeyEvent{Name: fyne.KeyBackspace}
	entry.OnKeyDown(key)

	assert.Equal(t, "H", entry.Text)
	assert.Equal(t, 0, entry.CursorRow)
	assert.Equal(t, 1, entry.CursorColumn)
}

func TestEntry_OnKeyDown_BackspaceBeyondText(t *testing.T) {
	entry := NewEntry()
	entry.SetText("Hi")
	right := &fyne.KeyEvent{Name: fyne.KeyRight}
	entry.OnKeyDown(right)
	entry.OnKeyDown(right)

	key := &fyne.KeyEvent{Name: fyne.KeyBackspace}
	entry.OnKeyDown(key)
	entry.OnKeyDown(key)
	entry.OnKeyDown(key)

	assert.Equal(t, "", entry.Text)
}

func TestEntry_OnKeyDown_BackspaceNewline(t *testing.T) {
	entry := NewMultiLineEntry()
	entry.SetText("H\ni")

	down := &fyne.KeyEvent{Name: fyne.KeyDown}
	entry.OnKeyDown(down)

	key := &fyne.KeyEvent{Name: fyne.KeyBackspace}
	entry.OnKeyDown(key)

	assert.Equal(t, "Hi", entry.Text)
}

func TestEntry_OnKeyDown_Backspace_Unicode(t *testing.T) {
	entry := NewEntry()

	key := new(fyne.KeyEvent)
	key.String = "è"
	entry.OnKeyDown(key)
	assert.Equal(t, 0, entry.CursorRow)
	assert.Equal(t, 1, entry.CursorColumn)

	bs := &fyne.KeyEvent{Name: fyne.KeyBackspace}
	entry.OnKeyDown(bs)
	assert.Equal(t, "", entry.Text)
	assert.Equal(t, 0, entry.CursorRow)
	assert.Equal(t, 0, entry.CursorColumn)
}

func TestEntry_OnKeyDown_Delete(t *testing.T) {
	entry := NewEntry()
	entry.SetText("Hi")
	right := &fyne.KeyEvent{Name: fyne.KeyRight}
	entry.OnKeyDown(right)
	assert.Equal(t, 0, entry.CursorRow)
	assert.Equal(t, 1, entry.CursorColumn)

	key := &fyne.KeyEvent{Name: fyne.KeyDelete}
	entry.OnKeyDown(key)

	assert.Equal(t, "H", entry.Text)
	assert.Equal(t, 0, entry.CursorRow)
	assert.Equal(t, 1, entry.CursorColumn)
}

func TestEntry_OnKeyDown_DeleteBeyondText(t *testing.T) {
	entry := NewEntry()
	entry.SetText("Hi")

	key := &fyne.KeyEvent{Name: fyne.KeyDelete}
	entry.OnKeyDown(key)
	entry.OnKeyDown(key)
	entry.OnKeyDown(key)

	assert.Equal(t, "", entry.Text)
}

func TestEntry_OnKeyDown_DeleteNewline(t *testing.T) {
	entry := NewEntry()
	entry.SetText("H\ni")

	right := &fyne.KeyEvent{Name: fyne.KeyRight}
	entry.OnKeyDown(right)

	key := &fyne.KeyEvent{Name: fyne.KeyDelete}
	entry.OnKeyDown(key)

	assert.Equal(t, "Hi", entry.Text)
}

func TestEntry_OnKeyDown_Home_End(t *testing.T) {
	entry := &Entry{}
	entry.SetText("Hi")
	assert.Equal(t, 0, entry.CursorRow)
	assert.Equal(t, 0, entry.CursorColumn)

	end := &fyne.KeyEvent{Name: fyne.KeyEnd}
	entry.OnKeyDown(end)
	assert.Equal(t, 0, entry.CursorRow)
	assert.Equal(t, 2, entry.CursorColumn)

	home := &fyne.KeyEvent{Name: fyne.KeyHome}
	entry.OnKeyDown(home)
	assert.Equal(t, 0, entry.CursorRow)
	assert.Equal(t, 0, entry.CursorColumn)
}

func TestEntryNotify(t *testing.T) {
	entry := NewEntry()
	changed := false

	entry.OnChanged = func(string) {
		changed = true
	}
	entry.SetText("Test")

	assert.True(t, changed)
}

func TestEntryFocus(t *testing.T) {
	entry := NewEntry()

	entry.OnFocusGained()
	assert.True(t, entry.Focused())

	entry.OnFocusLost()
	assert.False(t, entry.Focused())
}

func TestEntryWindowFocus(t *testing.T) {
	entry := NewEntry()
	canvas := test.Canvas()

	canvas.Focus(entry)
	assert.True(t, entry.Focused())
}

func TestEntryFocusHighlight(t *testing.T) {
	entry := NewEntry()

	entry.OnFocusGained()
	assert.True(t, entry.focused)

	entry.OnFocusLost()
	assert.False(t, entry.focused)
}

func TestEntry_CursorRow(t *testing.T) {
	entry := NewMultiLineEntry()
	entry.SetText("test")
	assert.Equal(t, 0, entry.CursorRow)

	// only 1 line, do nothing
	down := &fyne.KeyEvent{Name: fyne.KeyDown}
	entry.OnKeyDown(down)
	assert.Equal(t, 0, entry.CursorRow)

	// 2 lines, this should increment
	entry.SetText("test\nrows")
	entry.OnKeyDown(down)
	assert.Equal(t, 1, entry.CursorRow)

	up := &fyne.KeyEvent{Name: fyne.KeyUp}
	entry.OnKeyDown(up)
	assert.Equal(t, 0, entry.CursorRow)

	// don't go beyond top
	entry.OnKeyDown(up)
	assert.Equal(t, 0, entry.CursorRow)
}

func TestEntry_CursorColumn(t *testing.T) {
	entry := NewEntry()
	entry.SetText("")
	assert.Equal(t, 0, entry.CursorColumn)

	// only 0 columns, do nothing
	right := &fyne.KeyEvent{Name: fyne.KeyRight}
	entry.OnKeyDown(right)
	assert.Equal(t, 0, entry.CursorColumn)

	// 1, this should increment
	entry.SetText("a")
	entry.OnKeyDown(right)
	assert.Equal(t, 1, entry.CursorColumn)

	left := &fyne.KeyEvent{Name: fyne.KeyLeft}
	entry.OnKeyDown(left)
	assert.Equal(t, 0, entry.CursorColumn)

	// don't go beyond left
	entry.OnKeyDown(left)
	assert.Equal(t, 0, entry.CursorColumn)
}

func TestEntry_CursorColumn_Wrap(t *testing.T) {
	entry := NewMultiLineEntry()
	entry.SetText("a\nb")
	assert.Equal(t, 0, entry.CursorRow)
	assert.Equal(t, 0, entry.CursorColumn)

	// go to end of line
	right := &fyne.KeyEvent{Name: fyne.KeyRight}
	entry.OnKeyDown(right)
	assert.Equal(t, 0, entry.CursorRow)
	assert.Equal(t, 1, entry.CursorColumn)

	// wrap to new line
	entry.OnKeyDown(right)
	assert.Equal(t, 1, entry.CursorRow)
	assert.Equal(t, 0, entry.CursorColumn)

	// and back
	left := &fyne.KeyEvent{Name: fyne.KeyLeft}
	entry.OnKeyDown(left)
	assert.Equal(t, 0, entry.CursorRow)
	assert.Equal(t, 1, entry.CursorColumn)
}

func TestEntry_CursorColumn_Jump(t *testing.T) {
	entry := NewMultiLineEntry()
	entry.SetText("a\nbc")

	// go to end of text
	right := &fyne.KeyEvent{Name: fyne.KeyRight}
	entry.OnKeyDown(right)
	entry.OnKeyDown(right)
	entry.OnKeyDown(right)
	entry.OnKeyDown(right)
	assert.Equal(t, 1, entry.CursorRow)
	assert.Equal(t, 2, entry.CursorColumn)

	// go up, to a shorter line
	up := &fyne.KeyEvent{Name: fyne.KeyUp}
	entry.OnKeyDown(up)
	assert.Equal(t, 0, entry.CursorRow)
	assert.Equal(t, 1, entry.CursorColumn)
}

func checkNewlineIgnored(t *testing.T, entry *Entry) {
	assert.Equal(t, 0, entry.CursorRow)

	// only 1 line, do nothing
	down := &fyne.KeyEvent{Name: fyne.KeyDown}
	entry.OnKeyDown(down)
	assert.Equal(t, 0, entry.CursorRow)

	// return is ignored, do nothing
	ret := &fyne.KeyEvent{Name: fyne.KeyReturn}
	entry.OnKeyDown(ret)
	assert.Equal(t, 0, entry.CursorRow)

	up := &fyne.KeyEvent{Name: fyne.KeyUp}
	entry.OnKeyDown(up)
	assert.Equal(t, 0, entry.CursorRow)

	// don't go beyond top
	entry.OnKeyDown(up)
	assert.Equal(t, 0, entry.CursorRow)
}

func TestSinglelineEntry_NewlineIgnored(t *testing.T) {
	entry := &Entry{MultiLine: false}
	entry.SetText("test")

	checkNewlineIgnored(t, entry)
}

func TestPasswordEntry_NewlineIgnored(t *testing.T) {
	entry := NewPasswordEntry()
	entry.SetText("test")

	checkNewlineIgnored(t, entry)
}

func TestPasswordEntry_Obfuscation(t *testing.T) {
	entry := NewPasswordEntry()

	key := new(fyne.KeyEvent)
	key.String = "Hié™שרה"
	entry.OnKeyDown(key)
	assert.Equal(t, "Hié™שרה", entry.Text)
	assert.Equal(t, "*******", entryRenderTexts(entry)[0].Text)
}
