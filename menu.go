package fyne

// Menu stores the information required for a standard menu.
// A menu can pop down from a MainMenu or could be a pop out menu.
type Menu struct {
	Label string
	Items []*MenuItem
}

// NewMenu creates a new menu given the specified label (to show in a MainMenu) and list of items to display.
func NewMenu(label string, items ...*MenuItem) *Menu {
	return &Menu{Label: label, Items: items}
}

// MenuItem is a single item within any menu, it contains a display Label and Action function that is called when tapped.
type MenuItem struct {
	ChildMenu   *Menu
	IsSeparator bool
	Label       string
	Action      func()
}

// NewMenuItem creates a new menu item from the passed label and action parameters.
func NewMenuItem(label string, action func()) *MenuItem {
	return &MenuItem{Label: label, Action: action}
}

// NewMenuItemSeparator creates a menu item that is to be used as a separator.
func NewMenuItemSeparator() *MenuItem {
	return &MenuItem{IsSeparator: true, Action: func() {}}
}

// MainMenu defines the data required to show a menu bar (desktop) or other appropriate top level menu.
type MainMenu struct {
	Items []*Menu
}

// NewMainMenu creates a top level menu structure used by fyne.Window for displaying a menubar
// (or appropriate equivalent).
func NewMainMenu(items ...*Menu) *MainMenu {
	return &MainMenu{items}
}
