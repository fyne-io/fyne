package test

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type mockCloud struct {
	configured bool
}

func (c *mockCloud) Cleanup(_ fyne.App) {
	c.configured = false
}

func (*mockCloud) ProviderDescription() string {
	return "Mock cloud implementation"
}

func (*mockCloud) ProviderIcon() fyne.Resource {
	return theme.ComputerIcon()
}

func (*mockCloud) ProviderName() string {
	return "mock"
}

func (c *mockCloud) Setup(_ fyne.App) error {
	c.configured = true
	return nil
}
