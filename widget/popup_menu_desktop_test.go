//go:build !mobile

package widget_test

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
)

func TestPopUpMenu_KeyboardControl(t *testing.T) {
	var lastTriggered string
	m, w := setupPopUpMenuWithSubmenusTest(t, func(triggered string) { lastTriggered = triggered })
	c := w.Canvas()
	m.ShowAtPosition(fyne.NewPos(13, 45))

	focused := c.Focused()
	assert.Equal(t, "*widget.PopUpMenu", reflect.TypeOf(focused).String())
	test.AssertRendersToMarkup(t, "popup_menu/desktop/kbd_ctrl_shown.xml", c)
	assert.Equal(t, "", lastTriggered)

	// Traverse
	c.Focused().TypedKey(&fyne.KeyEvent{Name: fyne.KeyDown})
	assert.Same(t, focused, c.Focused())
	test.AssertRendersToMarkup(t, "popup_menu/desktop/kbd_ctrl_first_active.xml", c)
	c.Focused().TypedKey(&fyne.KeyEvent{Name: fyne.KeyRight})
	assert.Same(t, focused, c.Focused())
	test.AssertRendersToMarkup(t, "popup_menu/desktop/kbd_ctrl_first_active.xml", c, "nothing happens when trying open entry without submenu")
	c.Focused().TypedKey(&fyne.KeyEvent{Name: fyne.KeyDown})
	assert.Same(t, focused, c.Focused())
	test.AssertRendersToMarkup(t, "popup_menu/desktop/kbd_ctrl_second_active.xml", c)
	c.Focused().TypedKey(&fyne.KeyEvent{Name: fyne.KeyDown})
	assert.Same(t, focused, c.Focused())
	test.AssertRendersToMarkup(t, "popup_menu/desktop/kbd_ctrl_second_active.xml", c, "no wrap when reaching end")
	c.Focused().TypedKey(&fyne.KeyEvent{Name: fyne.KeyRight})
	assert.Same(t, focused, c.Focused())
	test.AssertRendersToMarkup(t, "popup_menu/desktop/kbd_ctrl_first_sub_active.xml", c)
	c.Focused().TypedKey(&fyne.KeyEvent{Name: fyne.KeyDown})
	assert.Same(t, focused, c.Focused())
	test.AssertRendersToMarkup(t, "popup_menu/desktop/kbd_ctrl_second_sub_active.xml", c)
	c.Focused().TypedKey(&fyne.KeyEvent{Name: fyne.KeyDown})
	assert.Same(t, focused, c.Focused())
	test.AssertRendersToMarkup(t, "popup_menu/desktop/kbd_ctrl_second_sub_active.xml", c, "no wrap when reaching end")
	c.Focused().TypedKey(&fyne.KeyEvent{Name: fyne.KeyRight})
	assert.Same(t, focused, c.Focused())
	test.AssertRendersToMarkup(t, "popup_menu/desktop/kbd_ctrl_first_sub_sub_active.xml", c)
	c.Focused().TypedKey(&fyne.KeyEvent{Name: fyne.KeyDown})
	assert.Same(t, focused, c.Focused())
	test.AssertRendersToMarkup(t, "popup_menu/desktop/kbd_ctrl_second_sub_sub_active.xml", c)
	c.Focused().TypedKey(&fyne.KeyEvent{Name: fyne.KeyDown})
	assert.Same(t, focused, c.Focused())
	test.AssertRendersToMarkup(t, "popup_menu/desktop/kbd_ctrl_second_sub_sub_active.xml", c, "no wrap when reaching end")
	c.Focused().TypedKey(&fyne.KeyEvent{Name: fyne.KeyRight})
	assert.Same(t, focused, c.Focused())
	test.AssertRendersToMarkup(t, "popup_menu/desktop/kbd_ctrl_second_sub_sub_active.xml", c, "nothing happens when trying open entry without submenu")
	c.Focused().TypedKey(&fyne.KeyEvent{Name: fyne.KeyLeft})
	assert.Same(t, focused, c.Focused())
	test.AssertRendersToMarkup(t, "popup_menu/desktop/kbd_ctrl_second_sub_active.xml", c)
	c.Focused().TypedKey(&fyne.KeyEvent{Name: fyne.KeyUp})
	assert.Same(t, focused, c.Focused())
	test.AssertRendersToMarkup(t, "popup_menu/desktop/kbd_ctrl_first_sub_active.xml", c)
	c.Focused().TypedKey(&fyne.KeyEvent{Name: fyne.KeyUp})
	assert.Same(t, focused, c.Focused())
	test.AssertRendersToMarkup(t, "popup_menu/desktop/kbd_ctrl_first_sub_active.xml", c, "no wrap when reaching start")
	c.Focused().TypedKey(&fyne.KeyEvent{Name: fyne.KeyLeft})
	assert.Same(t, focused, c.Focused())
	test.AssertRendersToMarkup(t, "popup_menu/desktop/kbd_ctrl_second_active.xml", c)
	c.Focused().TypedKey(&fyne.KeyEvent{Name: fyne.KeyLeft})
	assert.Same(t, focused, c.Focused())
	test.AssertRendersToMarkup(t, "popup_menu/desktop/kbd_ctrl_second_active.xml", c, "nothing happens when no sub-menus are open")

	// trigger actions
	c.Focused().TypedKey(&fyne.KeyEvent{Name: fyne.KeyEnter})
	assert.Nil(t, c.Focused())
	assert.Equal(t, "Option B", lastTriggered)
	test.AssertRendersToMarkup(t, "popup_menu/desktop/kbd_ctrl_dismissed.xml", c)

	m.Show()
	assert.Equal(t, "*widget.PopUpMenu", reflect.TypeOf(focused).String())
	test.AssertRendersToMarkup(t, "popup_menu/desktop/kbd_ctrl_shown.xml", c)
	c.Focused().TypedKey(&fyne.KeyEvent{Name: fyne.KeyDown})
	c.Focused().TypedKey(&fyne.KeyEvent{Name: fyne.KeyDown})
	c.Focused().TypedKey(&fyne.KeyEvent{Name: fyne.KeyRight})
	c.Focused().TypedKey(&fyne.KeyEvent{Name: fyne.KeyReturn})
	assert.Nil(t, c.Focused())
	assert.Equal(t, "Sub Option A", lastTriggered)
	test.AssertRendersToMarkup(t, "popup_menu/desktop/kbd_ctrl_dismissed.xml", c)

	m.Show()
	assert.Equal(t, "*widget.PopUpMenu", reflect.TypeOf(focused).String())
	test.AssertRendersToMarkup(t, "popup_menu/desktop/kbd_ctrl_shown.xml", c)
	c.Focused().TypedKey(&fyne.KeyEvent{Name: fyne.KeyDown})
	c.Focused().TypedKey(&fyne.KeyEvent{Name: fyne.KeyDown})
	c.Focused().TypedKey(&fyne.KeyEvent{Name: fyne.KeyRight})
	c.Focused().TypedKey(&fyne.KeyEvent{Name: fyne.KeyDown})
	c.Focused().TypedKey(&fyne.KeyEvent{Name: fyne.KeyRight})
	c.Focused().TypedKey(&fyne.KeyEvent{Name: fyne.KeySpace})
	assert.Nil(t, c.Focused())
	assert.Equal(t, "Sub Sub Option A", lastTriggered)
	test.AssertRendersToMarkup(t, "popup_menu/desktop/kbd_ctrl_dismissed.xml", c)

	// dismiss without triggering action
	lastTriggered = "none"
	m.Show()
	assert.Equal(t, "*widget.PopUpMenu", reflect.TypeOf(focused).String())
	test.AssertRendersToMarkup(t, "popup_menu/desktop/kbd_ctrl_shown.xml", c)
	c.Focused().TypedKey(&fyne.KeyEvent{Name: fyne.KeyDown})
	c.Focused().TypedKey(&fyne.KeyEvent{Name: fyne.KeyDown})
	c.Focused().TypedKey(&fyne.KeyEvent{Name: fyne.KeyEscape})
	assert.Nil(t, c.Focused())
	assert.Equal(t, "none", lastTriggered)
	test.AssertRendersToMarkup(t, "popup_menu/desktop/kbd_ctrl_dismissed.xml", c)
}
