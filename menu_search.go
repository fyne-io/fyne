package fyne

// IsHelpMenu checks if a menu is the Help menu by its label.
// This is used internally to automatically add search functionality to the Help menu.
func IsHelpMenu(menu *Menu) bool {
	return menu.Label == "Help"
}
