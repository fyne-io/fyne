package test

import (
	"testing"

	"fyne.io/fyne"

	"github.com/stretchr/testify/assert"
)

func TestAssertNotificationSent(t *testing.T) {
	n := fyne.NewNotification("Test Title", "Some content")
	myApp := fyne.CurrentApp()

	AssertNotificationSent(t, n, func() {
		fyne.CurrentApp().SendNotification(n)
	})
	assert.Equal(t, myApp, fyne.CurrentApp())
}

func TestAssertNotificationSent_NotSent(t *testing.T) {
	myApp := fyne.CurrentApp()

	AssertNotificationSent(t, nil, func() {
		// don't send anything
	})
	assert.Equal(t, myApp, fyne.CurrentApp())
}
