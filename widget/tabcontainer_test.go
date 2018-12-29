package widget

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTabContainer_CurrentTabIndex(t *testing.T) {
	tabs := NewTabContainer(&TabItem{Text: "Test", Content: NewLabel("Test")})

	assert.Equal(t, 1, len(tabs.Items))
	assert.Equal(t, 0, tabs.CurrentTabIndex())
}

func TestTabContainer_CurrentTab(t *testing.T) {
	tab1 := &TabItem{Text: "Test1", Content: NewLabel("Test1")}
	tab2 := &TabItem{Text: "Test2", Content: NewLabel("Test2")}
	tabs := NewTabContainer(tab1, tab2)

	assert.Equal(t, 2, len(tabs.Items))
	assert.Equal(t, tab1, tabs.CurrentTab())
}

func TestTabContainer_SelectTab(t *testing.T) {
	tab1 := &TabItem{Text: "Test1", Content: NewLabel("Test1")}
	tab2 := &TabItem{Text: "Test2", Content: NewLabel("Test2")}
	tabs := NewTabContainer(tab1, tab2)

	assert.Equal(t, 2, len(tabs.Items))
	assert.Equal(t, tab1, tabs.CurrentTab())

	tabs.SelectTab(tab2)
	assert.Equal(t, tab2, tabs.CurrentTab())
}

func TestTabContainer_SelectTabIndex(t *testing.T) {
	tabs := NewTabContainer(&TabItem{Text: "Test1", Content: NewLabel("Test1")},
		&TabItem{Text: "Test2", Content: NewLabel("Test2")})

	assert.Equal(t, 2, len(tabs.Items))
	assert.Equal(t, 0, tabs.CurrentTabIndex())

	tabs.SelectTabIndex(1)
	assert.Equal(t, 1, tabs.CurrentTabIndex())
}
