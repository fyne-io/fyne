package widget

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
)

func TestSearchableMainMenu_IndexMenuItems(t *testing.T) {
	fileMenu := fyne.NewMenu("File",
		fyne.NewMenuItem("New", nil),
		fyne.NewMenuItem("Open", nil),
		fyne.NewMenuItem("Save", nil),
	)

	editMenu := fyne.NewMenu("Edit",
		fyne.NewMenuItem("Cut", nil),
		fyne.NewMenuItem("Copy", nil),
		fyne.NewMenuItem("Paste", nil),
	)

	formatTextMenu := fyne.NewMenu("Text",
		fyne.NewMenuItem("Bold", nil),
		fyne.NewMenuItem("Italic", nil),
	)

	formatMenu := fyne.NewMenu("Format",
		fyne.NewMenuItem("Font", nil),
		&fyne.MenuItem{Label: "Text", ChildMenu: formatTextMenu},
	)

	mainMenu := fyne.NewMainMenu(fileMenu, editMenu, formatMenu)
	searchable := NewSearchableMainMenu(mainMenu)

	assert.Equal(t, 10, len(searchable.searchItems)) // 3 + 3 + 2 + 2

	found := false
	for _, item := range searchable.searchItems {
		if item.Item.Label == "Bold" {
			assert.Equal(t, []string{"Format", "Text"}, item.Path)
			found = true
			break
		}
	}
	assert.True(t, found, "Bold menu item should be indexed")
}

func TestSearchableMainMenu_Search(t *testing.T) {
	fileMenu := fyne.NewMenu("File",
		fyne.NewMenuItem("New", nil),
		fyne.NewMenuItem("Open", nil),
		fyne.NewMenuItem("Save", nil),
		fyne.NewMenuItem("Save As", nil),
	)

	editMenu := fyne.NewMenu("Edit",
		fyne.NewMenuItem("Find", nil),
		fyne.NewMenuItem("Replace", nil),
	)

	mainMenu := fyne.NewMainMenu(fileMenu, editMenu)
	searchable := NewSearchableMainMenu(mainMenu)

	results := searchable.Search("save")
	assert.Equal(t, 2, len(results))
	assert.Equal(t, "Save", results[0].Item.Label)
	assert.Equal(t, "Save As", results[1].Item.Label)

	results = searchable.Search("find")
	assert.Equal(t, 1, len(results))
	assert.Equal(t, "Find", results[0].Item.Label)

	results = searchable.Search("SAVE")
	assert.Equal(t, 2, len(results))

	results = searchable.Search("")
	assert.Equal(t, 0, len(results))

	results = searchable.Search("xyz")
	assert.Equal(t, 0, len(results))
}

func TestSearchableMainMenu_SearchWithShortcuts(t *testing.T) {
	fileMenu := fyne.NewMenu("File",
		&fyne.MenuItem{
			Label:    "Save",
			Action:   nil,
			Shortcut: &fyne.ShortcutPaste{},
		},
	)

	mainMenu := fyne.NewMainMenu(fileMenu)
	searchable := NewSearchableMainMenu(mainMenu)

	results := searchable.Search("Paste")
	assert.Equal(t, 1, len(results))
	assert.Equal(t, "Save", results[0].Item.Label)
}

func TestCreateSearchResultMenuItem(t *testing.T) {
	originalItem := fyne.NewMenuItem("Save", func() {})
	searchItem := MenuSearchItem{
		Item: originalItem,
		Path: []string{"File"},
	}

	resultItem := CreateSearchResultMenuItem(searchItem)

	assert.Equal(t, "File → Save", resultItem.Label)
	assert.NotNil(t, resultItem.Action)
	assert.Equal(t, originalItem.Disabled, resultItem.Disabled)
}

func TestMenuWithGlobalSearch_Creation(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	helpMenu := fyne.NewMenu("Help",
		fyne.NewMenuItem("About", nil),
	)

	mainMenu := fyne.NewMainMenu(
		fyne.NewMenu("File", fyne.NewMenuItem("New", nil)),
		helpMenu,
	)

	globalSearchMenu := NewMenuWithGlobalSearch(helpMenu, mainMenu)

	assert.NotNil(t, globalSearchMenu)
	assert.NotNil(t, globalSearchMenu.searchableMainMenu)
	assert.True(t, globalSearchMenu.searchEnabled)
	assert.NotNil(t, globalSearchMenu.searchEntry)

	assert.Greater(t, len(globalSearchMenu.Items), 2)
	_, isEntry := globalSearchMenu.Items[0].(*Entry)
	assert.True(t, isEntry)
}

func TestMenuWithGlobalSearch_SearchAndDisplay(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	helpMenu := fyne.NewMenu("Help")
	fileMenu := fyne.NewMenu("File",
		fyne.NewMenuItem("Save", nil),
		fyne.NewMenuItem("Open", nil),
	)

	mainMenu := fyne.NewMainMenu(fileMenu, helpMenu)
	globalSearchMenu := NewMenuWithGlobalSearch(helpMenu, mainMenu)

	globalSearchMenu.onGlobalSearchChanged("save")

	assert.NotNil(t, globalSearchMenu.searchResults)
	assert.Equal(t, 1, len(globalSearchMenu.searchResults))

	assert.GreaterOrEqual(t, len(globalSearchMenu.Items), 3)
}

func TestMenuWithGlobalSearch_ResetSearch(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	helpMenu := fyne.NewMenu("Help",
		fyne.NewMenuItem("About", nil),
	)

	mainMenu := fyne.NewMainMenu(
		fyne.NewMenu("File", fyne.NewMenuItem("Save", nil)),
		helpMenu,
	)

	globalSearchMenu := NewMenuWithGlobalSearch(helpMenu, mainMenu)
	originalItemCount := len(globalSearchMenu.Items)

	globalSearchMenu.onGlobalSearchChanged("save")
	globalSearchMenu.onGlobalSearchChanged("")

	assert.Nil(t, globalSearchMenu.searchResults)
	assert.Equal(t, originalItemCount, len(globalSearchMenu.Items))
}

func TestMenuWithGlobalSearch_NoResults(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	helpMenu := fyne.NewMenu("Help")
	mainMenu := fyne.NewMainMenu(
		fyne.NewMenu("File", fyne.NewMenuItem("Save", nil)),
		helpMenu,
	)

	globalSearchMenu := NewMenuWithGlobalSearch(helpMenu, mainMenu)
	globalSearchMenu.onGlobalSearchChanged("xyz")

	hasNoResults := false
	for _, item := range globalSearchMenu.Items {
		if mi, ok := item.(*menuItem); ok {
			if mi.Item.Label == "No results found" && mi.Item.Disabled {
				hasNoResults = true
				break
			}
		}
	}
	assert.True(t, hasNoResults, "Should show 'No results found' message")
}

func TestAddSearchToMainMenu(t *testing.T) {
	mainMenu := fyne.NewMainMenu(
		fyne.NewMenu("File", fyne.NewMenuItem("New", nil)),
	)

	result := fyne.AddSearchToMainMenu(mainMenu)

	assert.Equal(t, 1, len(result.Items))
	assert.Equal(t, "File", result.Items[0].Label)
	assert.Greater(t, len(result.Items[0].Items), 0)
	assert.Equal(t, "Search Menu...", result.Items[0].Items[0].Label)

	mainMenu2 := fyne.NewMainMenu(
		fyne.NewMenu("File", fyne.NewMenuItem("New", nil)),
		fyne.NewMenu("Help", fyne.NewMenuItem("About", nil)),
	)

	result2 := fyne.AddSearchToMainMenu(mainMenu2)

	assert.Equal(t, 2, len(result2.Items))
	assert.Equal(t, "File", result2.Items[0].Label)
	assert.GreaterOrEqual(t, len(result2.Items[0].Items), 3)
	assert.Equal(t, "Search Menu...", result2.Items[0].Items[0].Label)
}

func TestNewMainMenuWithSearch(t *testing.T) {
	fileMenu := fyne.NewMenu("File", fyne.NewMenuItem("New", nil))
	mainMenu := fyne.NewMainMenuWithSearch(fileMenu)

	assert.NotNil(t, mainMenu)
	assert.Equal(t, 1, len(mainMenu.Items)) // File and Help
	assert.Equal(t, "File", mainMenu.Items[0].Label)
	assert.Greater(t, len(mainMenu.Items[0].Items), 0)
	assert.Equal(t, "Search Menu...", mainMenu.Items[0].Items[0].Label)
}
