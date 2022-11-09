package test

import (
	"testing"

	"fyne.io/fyne/v2"
	"github.com/stretchr/testify/assert"
)

func TestTestApp_CloudProvider(t *testing.T) {
	a := NewApp()
	c := &mockCloud{}
	a.SetCloudProvider(c)

	assert.Equal(t, c, a.CloudProvider())
}

func TestFyneApp_transitionCloud(t *testing.T) {
	a := NewApp()
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
