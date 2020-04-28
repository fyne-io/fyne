package widget_test

import (
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/internal/cache"
	"fyne.io/fyne/internal/painter/software"
	"fyne.io/fyne/internal/widget"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/test"
	"fyne.io/fyne/theme"
)

func TestMenu_ItemHovered(t *testing.T) {
	sub1 := fyne.NewMenuItem("sub1", nil)
	sub1.ChildMenu = fyne.NewMenu("",
		fyne.NewMenuItem("sub1 A", nil),
		fyne.NewMenuItem("sub1 B", nil),
	)
	sub2sub := fyne.NewMenuItem("sub2sub", nil)
	sub2sub.ChildMenu = fyne.NewMenu("",
		fyne.NewMenuItem("sub2sub A", nil),
		fyne.NewMenuItem("sub2sub B", nil),
	)
	sub2 := fyne.NewMenuItem("sub2", nil)
	sub2.ChildMenu = fyne.NewMenu("",
		fyne.NewMenuItem("sub2 A", nil),
		fyne.NewMenuItem("sub2 B", nil),
		sub2sub,
	)
	m := widget.NewMenu(
		fyne.NewMenu("",
			fyne.NewMenuItem("Foo", nil),
			fyne.NewMenuItemSeparator(),
			fyne.NewMenuItem("Bar", nil),
			sub1,
			sub2,
		),
	)
	size := m.MinSize()
	m.Resize(size)

	mi := m.Items[0].(*widget.MenuItem)
	r1 := cache.Renderer(mi)
	assert.Equal(t, color.Transparent, r1.BackgroundColor())
	mi.MouseIn(nil)
	assert.Equal(t, theme.HoverColor(), r1.BackgroundColor())
	mi.MouseOut()
	assert.Equal(t, color.Transparent, r1.BackgroundColor())

	sub1Widget := m.Items[3].(*widget.MenuItem)
	assert.Equal(t, sub1, sub1Widget.Item)
	sub2Widget := m.Items[4].(*widget.MenuItem)
	assert.Equal(t, sub2, sub2Widget.Item)
	assert.False(t, sub1Widget.Child.Visible(), "submenu is invisible initially")
	assert.False(t, sub2Widget.Child.Visible(), "submenu is invisible initially")
	sub1Widget.MouseIn(nil)
	assert.True(t, sub1Widget.Child.Visible(), "hovering item shows submenu")
	assert.False(t, sub2Widget.Child.Visible(), "other Child menu stays hidden")
	sub1Widget.MouseOut()
	assert.True(t, sub1Widget.Child.Visible(), "hover out does not hide submenu")
	assert.False(t, sub2Widget.Child.Visible(), "other Child menu still hidden")
	sub2Widget.MouseIn(nil)
	assert.False(t, sub1Widget.Child.Visible(), "hovering other item hides current submenu")
	assert.True(t, sub2Widget.Child.Visible(), "other Child menu shows up")

	sub2subWidget := sub2Widget.Child.Items[2].(*widget.MenuItem)
	assert.Equal(t, sub2sub, sub2subWidget.Item)
	assert.False(t, sub2subWidget.Child.Visible(), "2nd level submenu is invisible initially")
	sub2Widget.MouseOut()
	sub2subWidget.MouseIn(nil)
	assert.True(t, sub2Widget.Child.Visible(), "1st level submenu stays visible")
	assert.True(t, sub2subWidget.Child.Visible(), "2nd level submenu shows up")
	sub2subWidget.MouseOut()
	assert.True(t, sub2Widget.Child.Visible(), "1st level submenu still visible")
	assert.True(t, sub2subWidget.Child.Visible(), "2nd level submenu still visible")

	sub1Widget.MouseIn(nil)
	assert.False(t, sub2Widget.Child.Visible(), "1st level submenu is hidden by other submenu")
	sub2Widget.MouseIn(nil)
	assert.False(t, sub2subWidget.Child.Visible(), "2nd level submenu is hidden when re-entering its parent")
}

func TestMenu_ItemTapped(t *testing.T) {
	tapped := false
	item1 := fyne.NewMenuItem("Foo", nil)
	item2 := fyne.NewMenuItem("Bar", func() { tapped = true })
	item3 := fyne.NewMenuItem("Sub", nil)
	subItem := fyne.NewMenuItem("Foo", func() {})
	item3.ChildMenu = fyne.NewMenu("", subItem)
	m := widget.NewMenu(fyne.NewMenu("", item1, item2, item3))
	size := m.MinSize()
	m.Resize(size)
	dismissed := false
	m.DismissAction = func() { dismissed = true }

	mi1 := m.Items[0].(*widget.MenuItem)
	mi2 := m.Items[1].(*widget.MenuItem)
	mi3 := m.Items[2].(*widget.MenuItem)
	assert.Equal(t, item1, mi1.Item)
	assert.Equal(t, item2, mi2.Item)
	assert.Equal(t, item3, mi3.Item)

	// tap on item without action does not panic
	test.Tap(mi1)
	assert.False(t, tapped)
	assert.False(t, dismissed, "tap on item w/o action does not dismiss the menu")
	assert.True(t, m.Visible(), "tap on item w/o action does not hide the menu")

	test.Tap(mi2)
	assert.True(t, tapped)
	assert.True(t, dismissed, "tap on item dismisses the menu")
	assert.True(t, m.Visible(), "tap on item does not hide the menu … the DismissAction is reponsible for that")

	dismissed = false // reset
	mi3.MouseIn(nil)
	sm := mi3.Child
	smi := sm.Items[0].(*widget.MenuItem)
	assert.Equal(t, subItem, smi.Item)
	assert.True(t, sm.Visible(), "sub menu is visible")

	test.Tap(smi)
	assert.True(t, dismissed, "tap on sub item dismisses the root menu")
	assert.True(t, m.Visible(), "tap on item does not hide the menu … the DismissAction is reponsible for that")
	assert.False(t, sm.Visible(), "tap on sub item hides the sub menu")

	newActionTapped := false
	item2.Action = func() { newActionTapped = true }
	test.Tap(mi2)
	assert.True(t, newActionTapped, "tap on item performs its current action")
}

func TestMenu_Layout(t *testing.T) {
	app := test.NewApp()
	defer test.NewApp()
	app.Settings().SetTheme(theme.DarkTheme())

	item1 := fyne.NewMenuItem("A", nil)
	item2 := fyne.NewMenuItem("B (long)", nil)
	sep := fyne.NewMenuItemSeparator()
	item3 := fyne.NewMenuItem("C", nil)
	subItem1 := fyne.NewMenuItem("subitem A", nil)
	subItem2 := fyne.NewMenuItem("subitem B", nil)
	subItem3 := fyne.NewMenuItem("subitem C (long)", nil)
	subsubItem1 := fyne.NewMenuItem("subsubitem A (long)", nil)
	subsubItem2 := fyne.NewMenuItem("subsubitem B", nil)
	subItem3.ChildMenu = fyne.NewMenu("", subsubItem1, subsubItem2)
	item3.ChildMenu = fyne.NewMenu("", subItem1, subItem2, subItem3)
	m := widget.NewMenu(fyne.NewMenu("", item1, sep, item2, item3))
	w := test.NewWindowWithPainter(m, software.NewPainter())
	defer w.Close()
	w.Resize(fyne.NewSize(1000, 1000))
	m.Resize(m.MinSize())
	c := w.Canvas()

	subItem := m.Items[3].(*widget.MenuItem)
	subItem.MouseIn(nil)
	subsubItem := subItem.Child.Items[2].(*widget.MenuItem)
	subsubItem.MouseIn(nil)

	menuWidth := 0
	submenuWidth := 0
	objects := test.LaidOutObjects(m)
	cons := selectContainers(objects)
	shadows := selectShadows(objects)
	submenuIcons := selectImages(objects)

	test.AssertImageMatches(t, "menu_layout.png", c.Capture())

	if assert.Len(t, cons, 3, "one container for each menu") &&
		assert.Len(t, shadows, 3, "one container for each menu") &&
		assert.Len(t, submenuIcons, 2, "one icon for each menu item with submenu") {
		// root menu
		submenuPos := assertMenu(t, cons[0], shadows[0], submenuIcons[0], []*fyne.MenuItem{item1, nil, item2, item3}, "B (long)", false, 3)
		menuWidth = submenuPos.X
		assert.Equal(t, submenuPos, subItem.Child.Position(), "correct submenu position")
		// sub menu
		subsubmenuPos := assertMenu(t, cons[1], shadows[1], submenuIcons[1], []*fyne.MenuItem{subItem1, subItem2, subItem3}, "subitem C (long)", true, 2)
		submenuWidth = subsubmenuPos.X
		assert.Equal(t, subsubmenuPos, subsubItem.Child.Position(), "correct subsubmenu position")
		// sub sub menu
		assertMenu(t, cons[2], shadows[2], nil, []*fyne.MenuItem{subsubItem1, subsubItem2}, "subsubitem A (long)", false, -1)
	}

	// move menu to the far right -> no space left for the submenu
	m.Move(fyne.NewPos(1000-menuWidth-10, 0))
	test.LaidOutObjects(m)
	test.AssertImageMatches(t, "menu_layout_no_space_on_right.png", c.Capture())
	assert.Equal(t, fyne.NewPos(-submenuWidth, 0), subItem.Child.Position(), "submenu is placed to the left if insufficient space to the right")

	// window space too small to place submenu to the left or to the right
	w.Resize(fyne.NewSize(menuWidth/2+submenuWidth, 1000))
	m.Resize(m.MinSize())
	m.Move(fyne.NewPos(0, 0))
	test.LaidOutObjects(m)
	test.AssertImageMatches(t, "menu_layout_no_space_on_both_sides.png", c.Capture())
	assert.Equal(t, fyne.NewPos(menuWidth/2, 0), subItem.Child.Position(), "submenu is placed as far right as possible if space is too tight to both sides")

	// window too short to place submenu
	winHeight := m.Size().Height + subItem.Child.Size().Height/2
	w.Resize(fyne.NewSize(1000, winHeight))
	m.Resize(m.MinSize())
	absSubItemPos := fyne.CurrentApp().Driver().AbsolutePositionForObject(subItem)
	test.LaidOutObjects(m)
	assert.Equal(t, fyne.NewPos(menuWidth, winHeight-(absSubItemPos.Y+subItem.Child.Size().Height)), subItem.Child.Position(), "submenu is placed as far right as possible if space is too tight to both sides")

	test.AssertImageMatches(t, "menu_layout_window_too_short.png", c.Capture())
}

func assertMenu(t *testing.T, c *fyne.Container, shadow *widget.Shadow, icon *canvas.Image, items []*fyne.MenuItem, longestLabel string, longestIsSub bool, subitem int) (subPos fyne.Position) {
	itemSize := canvas.NewText(longestLabel, color.Black).MinSize().Add(fyne.NewSize(theme.Padding()*4, theme.Padding()*2))
	if longestIsSub {
		itemSize.Width += theme.IconInlineSize()
	}
	yOff := 0
	itemCount := 0
	sepCount := 0
	for i, item := range items {
		o := c.Objects[i]
		assert.Equal(t, fyne.NewPos(0, yOff), o.Position())
		if item != nil {
			assert.Equal(t, item, o.(*widget.MenuItem).Item)
			assert.Equal(t, itemSize, o.Size())
			if i == subitem {
				assert.Equal(t, fyne.NewPos(itemSize.Width-theme.IconInlineSize(), (itemSize.Height-theme.IconInlineSize())/2), icon.Position())
				assert.Equal(t, fyne.NewSize(theme.IconInlineSize(), theme.IconInlineSize()), icon.Size())
				subPos = fyne.NewPos(itemSize.Width, 0)
			}
			yOff += itemSize.Height + theme.Padding()
			itemCount++
		} else { // separator
			assert.IsType(t, (*canvas.Rectangle)(nil), o)
			assert.Equal(t, fyne.NewSize(itemSize.Width, 2), o.Size())
			yOff += 2 + theme.Padding()
			sepCount++
		}
	}

	// height = item heights + sep heights + padding between items/seps & at start/end of menu
	size := fyne.NewSize(itemSize.Width, itemCount*itemSize.Height+sepCount*2+(2+itemCount+sepCount-1)*theme.Padding())
	menu := c.Objects[0].(*widget.MenuItem).Parent
	assert.Equal(t, size, menu.MinSize())

	assert.Equal(t, fyne.NewPos(0, 0), shadow.Position())
	assert.Equal(t, size, shadow.Size())

	assert.Equal(t, layout.NewVBoxLayout(), c.Layout)
	assert.Len(t, c.Objects, itemCount+sepCount, "container children size is equal to item + sep count")
	assert.Equal(t, fyne.NewPos(0, theme.Padding()), c.Position())
	assert.Equal(t, size.Subtract(fyne.NewSize(0, theme.Padding()*2)), c.Size(), "container size does not include leading & trailing padding")

	return
}

func selectContainers(objects []fyne.CanvasObject) (containers []*fyne.Container) {
	for _, object := range objects {
		if c, ok := object.(*fyne.Container); ok {
			containers = append(containers, c)
		}
	}
	return
}

func selectShadows(objects []fyne.CanvasObject) (shadows []*widget.Shadow) {
	for _, object := range objects {
		if c, ok := object.(*widget.Shadow); ok {
			shadows = append(shadows, c)
		}
	}
	return
}

func selectImages(objects []fyne.CanvasObject) (images []*canvas.Image) {
	for _, object := range objects {
		if c, ok := object.(*canvas.Image); ok {
			images = append(images, c)
		}
	}
	return
}
