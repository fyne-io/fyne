package app

import (
	"errors"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"

	"github.com/stretchr/testify/assert"
)

func TestFyneApp_SetCloudProvider(t *testing.T) {
	a := test.NewTempApp(t)
	p := &mockCloud{}
	a.SetCloudProvider(p)

	assert.Equal(t, p, a.CloudProvider())
	assert.True(t, p.configured)
}

func TestFyneApp_SetCloudProvider_Cleanup(t *testing.T) {
	a := test.NewTempApp(t)
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
	a := test.NewTempApp(t)
	p := &mockCloud{}
	preferenceChanged := false
	settingsChan := make(chan fyne.Settings)
	a.Preferences().AddChangeListener(func() {
		preferenceChanged = true
	})
	a.Settings().AddListener(func(s fyne.Settings) {
		go func() { settingsChan <- s }()
	})

	done := make(chan struct{})
	go func() {
		<-settingsChan // settings were updated
		assert.True(t, preferenceChanged)
		close(done)
	}()

	a.SetCloudProvider(p)
	<-done
}

func TestFyneApp_transitionCloud_Preferences(t *testing.T) {
	a := test.NewTempApp(t)
	a.Preferences().SetString("key", "blank")

	assert.Equal(t, "blank", a.Preferences().String("key"))

	p := &mockCloud{}
	a.SetCloudProvider(p)

	assert.Equal(t, "", a.Preferences().String("key"))
}

func TestFyneApp_transitionCloud_Storage(t *testing.T) {
	a := test.NewTempApp(t)
	a.Storage().Create("nothere")

	l := a.Storage().List()
	assert.Len(t, l, 1)

	p := &mockCloud{}
	a.SetCloudProvider(p)

	l = a.Storage().List()
	assert.Empty(t, l)
}

type mockCloud struct {
	configured, cleaned bool
}

func (c *mockCloud) Cleanup(_ fyne.App) {
	c.cleaned = true
}

func (*mockCloud) CloudPreferences(fyne.App) fyne.Preferences {
	return &internal.InMemoryPreferences{}
}

func (*mockCloud) CloudStorage(fyne.App) fyne.Storage {
	return &mockCloudStorage{}
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

type mockCloudStorage struct{}

func (*mockCloudStorage) Create(string) (fyne.URIWriteCloser, error) {
	return nil, errors.New("not implemented")
}

func (*mockCloudStorage) List() []string {
	return []string{}
}

func (*mockCloudStorage) Open(string) (fyne.URIReadCloser, error) {
	return nil, errors.New("not implemented")
}

func (*mockCloudStorage) Remove(string) error {
	return errors.New("not implemented")
}

func (*mockCloudStorage) RootURI() fyne.URI {
	u, _ := storage.ParseURI("mock://")
	return u
}

func (*mockCloudStorage) Save(string) (fyne.URIWriteCloser, error) {
	return nil, errors.New("not implemented")
}
