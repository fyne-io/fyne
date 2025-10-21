package widget

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
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

// MenuSearchItem represents a searchable menu item with its path through the menu hierarchy
type MenuSearchItem struct {
	Item       *fyne.MenuItem
	Path       []string
	Parent     *fyne.Menu
	ParentItem *fyne.MenuItem
}

// SearchableMainMenu wraps a MainMenu to provide search functionality across all menus
type SearchableMainMenu struct {
	MainMenu    *fyne.MainMenu
	searchItems []MenuSearchItem
}

// NewSearchableMainMenu creates a searchable wrapper around a MainMenu
func NewSearchableMainMenu(mainMenu *fyne.MainMenu) *SearchableMainMenu {
	s := &SearchableMainMenu{
		MainMenu: mainMenu,
	}
	s.indexMenuItems()
	return s
}

// indexMenuItems builds an index of all menu items for searching
func (s *SearchableMainMenu) indexMenuItems() {
	s.searchItems = []MenuSearchItem{}

	for _, menu := range s.MainMenu.Items {
		s.indexMenu(menu, []string{}, nil, nil)
	}
}

// indexMenu recursively indexes menu items
func (s *SearchableMainMenu) indexMenu(menu *fyne.Menu, path []string, parentItem *fyne.MenuItem, parentMenu *fyne.Menu) {
	newPath := append(path, menu.Label)

	for _, item := range menu.Items {
		if item.IsSeparator {
			continue
		}

		searchItem := MenuSearchItem{
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
func (s *SearchableMainMenu) Search(query string) []MenuSearchItem {
	if query == "" {
		return []MenuSearchItem{}
	}

	query = strings.ToLower(strings.TrimSpace(query))
	var results []MenuSearchItem

	for _, searchItem := range s.searchItems {
		if s.matchesQuery(searchItem, query) {
			results = append(results, searchItem)
		}
	}

	return results
}

// matchesQuery checks if a menu item matches the search query
func (s *SearchableMainMenu) matchesQuery(item MenuSearchItem, query string) bool {
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

// CreateSearchResultMenuItem creates a menu item that represents a search result
func CreateSearchResultMenuItem(searchItem MenuSearchItem) *fyne.MenuItem {
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

// MenuWithGlobalSearch extends Menu to search across all menus in a MainMenu
type MenuWithGlobalSearch struct {
	*Menu
	searchableMainMenu *SearchableMainMenu
	searchResults      []MenuSearchItem
	originalItems      []fyne.CanvasObject
	minSearchWidth     float32
}

// NewMenuWithGlobalSearch creates a menu that can search across all menus in a MainMenu
// This is automatically used when a menu contains a search menu item (created by AddSearchToMainMenu)
// The search functionality is automatically added to the Help menu, or creates a Help menu if it doesn't exist
func NewMenuWithGlobalSearch(menu *fyne.Menu, mainMenu *fyne.MainMenu) *MenuWithGlobalSearch {
	var searchLabel string
	filteredItems := make([]*fyne.MenuItem, 0, len(menu.Items))
	for _, item := range menu.Items {
		if fyne.IsSearchMenuItem(item) {
			searchLabel = item.Label
			continue
		}
		filteredItems = append(filteredItems, item)
	}

	if searchLabel == "" {
		searchLabel = fyne.DefaultSearchMenuLabel
	}

	filteredMenu := fyne.NewMenu(menu.Label, filteredItems...)

	m := &MenuWithGlobalSearch{
		Menu:               NewMenuWithSearch(filteredMenu),
		searchableMainMenu: NewSearchableMainMenu(mainMenu),
	}

	if m.Menu.searchEntry != nil {
		m.Menu.searchEntry.PlaceHolder = searchLabel
		m.searchEntry = m.Menu.searchEntry

		placeholderText := NewRichTextWithText(searchLabel)
		textSize := placeholderText.MinSize()

		th := m.Theme()
		innerPadding := th.Size(theme.SizeNameInnerPadding)
		inputBorder := th.Size(theme.SizeNameInputBorder)

		minWidth := textSize.Width + (innerPadding * 4) + (inputBorder * 2) + 40
		m.minSearchWidth = minWidth

		if len(m.Menu.Items) > 0 {
			wrappedEntry := newMinWidthContainer(m.Menu.searchEntry, minWidth)
			m.Menu.Items[0] = wrappedEntry
		}
	}

	m.initGlobalSearchHandlers()
	return m
}

// MinSize returns the minimum size for the menu, ensuring it's wide enough for the search field
func (m *MenuWithGlobalSearch) MinSize() fyne.Size {
	baseSize := m.Menu.MinSize()

	if m.minSearchWidth > 0 && baseSize.Width < m.minSearchWidth {
		baseSize.Width = m.minSearchWidth
	}

	return baseSize
}

// initGlobalSearchHandlers sets up the search handlers for global search
func (m *MenuWithGlobalSearch) initGlobalSearchHandlers() {
	if m.searchEntry == nil {
		return
	}

	if len(m.Menu.Items) > 2 {
		m.originalItems = make([]fyne.CanvasObject, len(m.Menu.Items)-2)
		copy(m.originalItems, m.Menu.Items[2:])
	}

	m.searchEntry.OnChanged = func(s string) {
		m.onGlobalSearchChanged(s)
	}
	m.searchEntry.OnSubmitted = func(_ string) {
		m.onGlobalSearchSubmitted()
	}
}

// onGlobalSearchChanged handles search query changes
func (m *MenuWithGlobalSearch) onGlobalSearchChanged(query string) {
	if query == "" {
		m.resetGlobalSearchResults()
		return
	}

	m.searchResults = m.searchableMainMenu.Search(query)
	m.displaySearchResults()
}

// displaySearchResults updates the menu to show search results
func (m *MenuWithGlobalSearch) displaySearchResults() {
	searchEntry := m.Menu.Items[0]
	separator := m.Menu.Items[1]

	resultItems := make([]fyne.CanvasObject, 0, len(m.searchResults)+2)
	resultItems = append(resultItems, searchEntry, separator)

	if len(m.searchResults) == 0 {
		noResultsItem := newMenuItem(&fyne.MenuItem{
			Label:    "No results found",
			Disabled: true,
		}, m.Menu)
		resultItems = append(resultItems, noResultsItem)
	} else {
		for _, result := range m.searchResults {
			resultMenuItem := CreateSearchResultMenuItem(result)
			menuItem := newMenuItem(resultMenuItem, m.Menu)
			resultItems = append(resultItems, menuItem)
		}
	}

	m.Menu.Items = resultItems
	m.Menu.Refresh()

	if m.Menu.Size().Height > 0 {
		newSize := m.Menu.MinSize()
		if newSize.Height > m.Menu.Size().Height {
			m.Menu.Resize(newSize)
		}
	}
}

// resetGlobalSearchResults resets the menu to show original items
func (m *MenuWithGlobalSearch) resetGlobalSearchResults() {
	m.searchResults = nil

	searchEntry := m.Menu.Items[0]
	separator := m.Menu.Items[1]

	newItems := make([]fyne.CanvasObject, 0, len(m.originalItems)+2)
	newItems = append(newItems, searchEntry, separator)
	newItems = append(newItems, m.originalItems...)
	m.Menu.Items = newItems

	m.Menu.Refresh()
}

// onGlobalSearchSubmitted handles Enter key press in search
func (m *MenuWithGlobalSearch) onGlobalSearchSubmitted() {
	if len(m.searchResults) > 0 {
		firstResult := m.searchResults[0]
		if firstResult.Item.Action != nil {
			firstResult.Item.Action()
			if m.OnDismiss != nil {
				m.OnDismiss()
			}
		}
	}
}
