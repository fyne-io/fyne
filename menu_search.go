package fyne

// searchMenuItemMarker is used as a sentinel value for search menu items
var searchMenuItemMarker = struct{}{}

// searchMenuItems tracks which menu items are search items
var searchMenuItems = make(map[*MenuItem]bool)

// IsSearchMenuItem checks if a menu item is a search placeholder
func IsSearchMenuItem(item *MenuItem) bool {
	// Check if this item has been marked as a search item
	return searchMenuItems[item]
}

// AddSearchToMainMenu adds search functionality to a MainMenu
// By default, it adds the search to the File menu
// This allows users to search and trigger any menu item across all menus
func AddSearchToMainMenu(mainMenu *MainMenu) *MainMenu {
	return AddSearchToMenu(mainMenu, "File")
}

// AddSearchToMenu adds search functionality to a specific menu in the MainMenu
// menuLabel specifies which menu should contain the search (e.g., "File", "Edit", "Help")
// The search item will be added at the beginning of the menu
func AddSearchToMenu(mainMenu *MainMenu, menuLabel string) *MainMenu {
	var targetMenu *Menu
	for _, menu := range mainMenu.Items {
		if menu.Label == menuLabel {
			targetMenu = menu
			break
		}
	}

	if targetMenu == nil && menuLabel == "Help" {
		targetMenu = NewMenu("Help")
		mainMenu.Items = append(mainMenu.Items, targetMenu)
	} else if targetMenu == nil {
		return mainMenu
	}

	for _, item := range targetMenu.Items {
		if IsSearchMenuItem(item) {
			return mainMenu
		}
	}

	searchItem := &MenuItem{
		Label:  "Search All Menus...",
		Action: func() {},
	}

	searchMenuItems[searchItem] = true

	newItems := make([]*MenuItem, 0, len(targetMenu.Items)+2)
	newItems = append(newItems, searchItem)

	if len(targetMenu.Items) > 0 && !targetMenu.Items[0].IsSeparator {
		newItems = append(newItems, NewMenuItemSeparator())
	}

	newItems = append(newItems, targetMenu.Items...)
	targetMenu.Items = newItems

	return mainMenu
}

// AddSearchToMenuAtPosition adds search functionality at a specific position in the menu
// position specifies where to insert the search item (0 = beginning, -1 = end)
func AddSearchToMenuAtPosition(mainMenu *MainMenu, menuLabel string, position int) *MainMenu {
	var targetMenu *Menu
	for _, menu := range mainMenu.Items {
		if menu.Label == menuLabel {
			targetMenu = menu
			break
		}
	}

	if targetMenu == nil && menuLabel == "Help" {
		targetMenu = NewMenu("Help")
		mainMenu.Items = append(mainMenu.Items, targetMenu)
	} else if targetMenu == nil {
		return mainMenu
	}

	for _, item := range targetMenu.Items {
		if IsSearchMenuItem(item) {
			return mainMenu
		}
	}

	searchItem := &MenuItem{
		Label:  "Search All Menus...",
		Action: func() {},
	}

	searchMenuItems[searchItem] = true
	if position < 0 || position >= len(targetMenu.Items) {
		targetMenu.Items = append(targetMenu.Items, searchItem)
	} else {
		newItems := make([]*MenuItem, 0, len(targetMenu.Items)+1)
		newItems = append(newItems, targetMenu.Items[:position]...)
		newItems = append(newItems, searchItem)
		newItems = append(newItems, targetMenu.Items[position:]...)
		targetMenu.Items = newItems
	}

	return mainMenu
}

// NewMainMenuWithSearch creates a new MainMenu with search functionality
// It automatically adds search to the File menu
func NewMainMenuWithSearch(items ...*Menu) *MainMenu {
	mainMenu := NewMainMenu(items...)
	return AddSearchToMainMenu(mainMenu)
}
