package fyne

import (
	"reflect"
	"testing"
)

func testFunc1() {}
func testFunc2() {}
func testFunc3() {}
func testFunc4() {}

var menuItemsTest = []struct {
	Label  string
	Action func()
	Item   *MenuItem
}{
	{"item1", testFunc1, &MenuItem{Label: "item1", Action: testFunc1}},
	{"item2", testFunc2, &MenuItem{Label: "item2", Action: testFunc2}},
	{"item3", testFunc3, &MenuItem{Label: "item3", Action: testFunc3}},
	{"item4", testFunc4, &MenuItem{Label: "item4", Action: testFunc4}},
}

var menuTest = []struct {
	Label string
}{
	{"menu1"},
	{"menu2"},
}

func TestNewMainMenu(t *testing.T) {
	var items []*MenuItem
	var menus []*Menu

	for _, tt := range menuItemsTest {
		t.Run(tt.Label, func(t *testing.T) {
			item := NewMenuItem(tt.Label, tt.Action)
			// Compare sprinted address
			if reflect.ValueOf(item.Action) != reflect.ValueOf(tt.Action) {
				t.Errorf("Expected %v but got %v", reflect.ValueOf(tt.Action), reflect.ValueOf(item.Action))
			}
			if item.Label != tt.Label {
				t.Errorf("Expected %v but got %v", tt.Label, item.Label)
			}
			items = append(items, item)
		})
	}

	if len(items) < 4 {
		t.Errorf("Expected %d menu items but got %d", len(menuItemsTest), len(items))
	}

	for _, tt := range menuTest {
		t.Run(tt.Label, func(t *testing.T) {
			menu := NewMenu(tt.Label, items...)

			if menu.Label != tt.Label {
				t.Errorf("Expected menu label %s but got %s", tt.Label, menu.Label)
			}

			if !reflect.DeepEqual(menu.Items, items) {
				t.Errorf("Expected items to resemble what was inputted, but got %v", menu.Items)
			}

			menus = append(menus, menu)
		})
	}

	mm := NewMainMenu(menus...)

	if !reflect.DeepEqual(mm.Items, menus) {
		t.Errorf("Expected main menu to contain all submenus but got %v", mm.Items)
	}

}
