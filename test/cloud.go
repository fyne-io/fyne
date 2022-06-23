package test

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type mockCloud struct {
	configured bool
}

func (c *mockCloud) ProviderDescription() string {
	return "Mock cloud implementation"
}

func (c *mockCloud) ProviderIcon() fyne.Resource {
	return theme.FyneLogo()
}

func (c *mockCloud) ProviderName() string {
	return "mock"
}

func (c *mockCloud) Setup(fyne.App) error {
	c.configured = true
	return nil
}
