package widget

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

// minWidthContainer is a container that enforces a minimum width for its content
type minWidthContainer struct {
	BaseWidget
	content  fyne.CanvasObject
	minWidth float32
}

func newMinWidthContainer(content fyne.CanvasObject, minWidth float32) *minWidthContainer {
	c := &minWidthContainer{
		content:  content,
		minWidth: minWidth,
	}
	c.ExtendBaseWidget(c)
	return c
}

func (c *minWidthContainer) CreateRenderer() fyne.WidgetRenderer {
	return &minWidthRenderer{
		container: c,
		objects:   []fyne.CanvasObject{c.content},
	}
}

type minWidthRenderer struct {
	container *minWidthContainer
	objects   []fyne.CanvasObject
}

func (r *minWidthRenderer) Destroy() {}

func (r *minWidthRenderer) Layout(size fyne.Size) {
	r.container.content.Resize(size)
	r.container.content.Move(fyne.NewPos(0, 0))
}

func (r *minWidthRenderer) MinSize() fyne.Size {
	minSize := r.container.content.MinSize()
	if minSize.Width < r.container.minWidth {
		minSize.Width = r.container.minWidth
	}
	return minSize
}

func (r *minWidthRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *minWidthRenderer) Refresh() {
	canvas.Refresh(r.container.content)
}

// menuSearchItem represents a searchable menu item with its path through the menu hierarchy
type menuSearchItem struct {
	Item       *fyne.MenuItem
	Path       []string
	Parent     *fyne.Menu
	ParentItem *fyne.MenuItem
}

// searchableMainMenu wraps a MainMenu to provide search functionality across all menus
type searchableMainMenu struct {
	MainMenu    *fyne.MainMenu
	searchItems []menuSearchItem
}

// newSearchableMainMenu creates a searchable wrapper around a MainMenu
func newSearchableMainMenu(mainMenu *fyne.MainMenu) *searchableMainMenu {
	s := &searchableMainMenu{
		MainMenu: mainMenu,
	}
	s.indexMenuItems()
	return s
}

// indexMenuItems builds an index of all menu items for searching
func (s *searchableMainMenu) indexMenuItems() {
	s.searchItems = []menuSearchItem{}

	for _, menu := range s.MainMenu.Items {
		s.indexMenu(menu, []string{}, nil, nil)
	}
}

// indexMenu recursively indexes menu items
func (s *searchableMainMenu) indexMenu(menu *fyne.Menu, path []string, parentItem *fyne.MenuItem, parentMenu *fyne.Menu) {
	newPath := append(path, menu.Label)

	for _, item := range menu.Items {
		if item.IsSeparator {
			continue
		}

		searchItem := menuSearchItem{
			Item:       item,
			Path:       newPath,
			Parent:     menu,
			ParentItem: parentItem,
		}
		s.searchItems = append(s.searchItems, searchItem)

		if item.ChildMenu != nil {
			s.indexMenu(item.ChildMenu, newPath, item, menu)
		}
	}
}

// Search finds menu items matching the query
func (s *searchableMainMenu) Search(query string) []menuSearchItem {
	if query == "" {
		return []menuSearchItem{}
	}

	query = strings.ToLower(strings.TrimSpace(query))
	var results []menuSearchItem

	for _, searchItem := range s.searchItems {
		if s.matchesQuery(searchItem, query) {
			results = append(results, searchItem)
		}
	}

	return results
}

// matchesQuery checks if a menu item matches the search query
func (s *searchableMainMenu) matchesQuery(item menuSearchItem, query string) bool {
	if strings.Contains(strings.ToLower(item.Item.Label), query) {
		return true
	}

	if item.Item.Shortcut != nil {
		shortcutStr := item.Item.Shortcut.ShortcutName()
		if strings.Contains(strings.ToLower(shortcutStr), query) {
			return true
		}
	}

	for _, pathComponent := range item.Path {
		if strings.Contains(strings.ToLower(pathComponent), query) {
			return true
		}
	}

	return false
}

// createSearchResultMenuItem creates a menu item that represents a search result
func createSearchResultMenuItem(searchItem menuSearchItem) *fyne.MenuItem {
	pathStr := strings.Join(searchItem.Path, " → ")
	if searchItem.Item.Label != "" {
		pathStr += " → " + searchItem.Item.Label
	}

	resultItem := &fyne.MenuItem{
		Label:    pathStr,
		Action:   searchItem.Item.Action,
		Disabled: searchItem.Item.Disabled,
		Icon:     searchItem.Item.Icon,
		Shortcut: searchItem.Item.Shortcut,
	}

	if searchItem.Item.ChildMenu != nil && searchItem.Item.Action == nil {
		resultItem.Action = func() {
			if leaf := findFirstActionableInMenu(searchItem.Item.ChildMenu); leaf != nil && leaf.Action != nil {
				leaf.Action()
			}
		}
	}

	return resultItem
}

// findFirstActionableInMenu finds the first actionable item in a menu
func findFirstActionableInMenu(menu *fyne.Menu) *fyne.MenuItem {
	for _, item := range menu.Items {
		if item.IsSeparator || item.Disabled {
			continue
		}
		if item.Action != nil && item.ChildMenu == nil {
			return item
		}
		if item.ChildMenu != nil {
			if leaf := findFirstActionableInMenu(item.ChildMenu); leaf != nil {
				return leaf
			}
		}
	}
	return nil
}
