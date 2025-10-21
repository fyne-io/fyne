package fyne

// DefaultSearchMenuLabel is the default label for search menu items
const DefaultSearchMenuLabel = "Search..."

// searchMenuItemMarker is used as a sentinel value for search menu items
var searchMenuItemMarker = struct{}{}

// searchMenuItems tracks which menu items are search items
var searchMenuItems = make(map[*MenuItem]bool)

// IsSearchMenuItem checks if a menu item is a search placeholder
func IsSearchMenuItem(item *MenuItem) bool {
	return searchMenuItems[item]
}

// AddSearchToMainMenu adds search functionality to a MainMenu
// It automatically adds the search to the Help menu (creating it if it doesn't exist)
// The search item will be added at the top of the Help menu
// If searchLabel is empty, DefaultSearchMenuLabel ("Search...") will be used
func AddSearchToMainMenu(mainMenu *MainMenu, searchLabel string) *MainMenu {
	if searchLabel == "" {
		searchLabel = DefaultSearchMenuLabel
	}

	var helpMenu *Menu
	for _, menu := range mainMenu.Items {
		if menu.Label == "Help" {
			helpMenu = menu
			break
		}
	}

	if helpMenu == nil {
		helpMenu = NewMenu("Help")
		mainMenu.Items = append(mainMenu.Items, helpMenu)
	}

	for _, item := range helpMenu.Items {
		if IsSearchMenuItem(item) {
			return mainMenu
		}
	}

	searchItem := &MenuItem{
		Label:  searchLabel,
		Action: func() {},
	}

	searchMenuItems[searchItem] = true

	newItems := make([]*MenuItem, 0, len(helpMenu.Items)+2)
	newItems = append(newItems, searchItem)

	if len(helpMenu.Items) > 0 && !helpMenu.Items[0].IsSeparator {
		newItems = append(newItems, NewMenuItemSeparator())
	}

	newItems = append(newItems, helpMenu.Items...)
	helpMenu.Items = newItems

	return mainMenu
}

// NewMainMenuWithSearch creates a new MainMenu with search functionality
// It automatically adds search to the Help menu (creating it if it doesn't exist)
// Pass an empty string for searchLabel to use the default "Search..." label
// Example: NewMainMenuWithSearch("", fileMenu, editMenu) or NewMainMenuWithSearch("Find...", fileMenu, editMenu)
func NewMainMenuWithSearch(searchLabel string, items ...*Menu) *MainMenu {
	mainMenu := NewMainMenu(items...)
	return AddSearchToMainMenu(mainMenu, searchLabel)
}
