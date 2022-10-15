package app

import (
	"errors"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"

	"github.com/stretchr/testify/assert"
)

func TestFyneApp_SetCloudProvider(t *testing.T) {
	a := NewWithID("io.fyne.test")
	p := &mockCloud{}
	a.SetCloudProvider(p)

	assert.Equal(t, p, a.CloudProvider())
	assert.True(t, p.configured)
}

func TestFyneApp_SetCloudProvider_Cleanup(t *testing.T) {
	a := NewWithID("io.fyne.test")
	p1 := &mockCloud{}
	p2 := &mockCloud{}
	a.SetCloudProvider(p1)

	assert.True(t, p1.configured)
	assert.False(t, p1.cleaned)

	a.SetCloudProvider(p2)

	assert.True(t, p1.cleaned)
	assert.True(t, p2.configured)
}

func TestFyneApp_transitionCloud(t *testing.T) {
	a := NewWithID("io.fyne.test")
	p := &mockCloud{}
	preferenceChanged := false
	settingsChan := make(chan fyne.Settings)
	a.Preferences().AddChangeListener(func() {
		preferenceChanged = true
	})
	a.Settings().AddChangeListener(settingsChan)
	a.SetCloudProvider(p)

	<-settingsChan // settings were updated
	assert.True(t, preferenceChanged)
}

func TestFyneApp_transitionCloud_Preferences(t *testing.T) {
	a := NewWithID("io.fyne.test")
	a.Preferences().SetString("key", "blank")

	assert.Equal(t, "blank", a.Preferences().String("key"))

	p := &mockCloud{}
	a.SetCloudProvider(p)

	assert.Equal(t, "", a.Preferences().String("key"))
}

func TestFyneApp_transitionCloud_Storage(t *testing.T) {
	a := NewWithID("io.fyne.test")
	a.Storage().Create("nothere")

	l := a.Storage().List()
	assert.Equal(t, 1, len(l))

	p := &mockCloud{}
	a.SetCloudProvider(p)

	l = a.Storage().List()
	assert.Equal(t, 0, len(l))
}

type mockCloud struct {
	configured, cleaned bool
}

func (c *mockCloud) Cleanup(_ fyne.App) {
	c.cleaned = true
}

func (c *mockCloud) CloudPreferences(fyne.App) fyne.Preferences {
	return &internal.InMemoryPreferences{}
}

func (c *mockCloud) CloudStorage(fyne.App) fyne.Storage {
	return &mockCloudStorage{}
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

func (c *mockCloud) Setup(_ fyne.App) error {
	c.configured = true
	return nil
}

type mockCloudStorage struct {
}

func (s *mockCloudStorage) Create(name string) (fyne.URIWriteCloser, error) {
	return nil, errors.New("not implemented")
}

func (s *mockCloudStorage) List() []string {
	return []string{}
}

func (s *mockCloudStorage) Open(name string) (fyne.URIReadCloser, error) {
	return nil, errors.New("not implemented")
}

func (s *mockCloudStorage) Remove(name string) error {
	return errors.New("not implemented")
}

func (s *mockCloudStorage) RootURI() fyne.URI {
	u, _ := storage.ParseURI("mock://")
	return u
}

func (s *mockCloudStorage) Save(name string) (fyne.URIWriteCloser, error) {
	return nil, errors.New("not implemented")
}
