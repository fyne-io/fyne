package test

import (
	"testing"

	"fyne.io/fyne/v2"

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

func TestAssertNotificationSent_Nil(t *testing.T) {
	AssertNotificationSent(t, nil, func() {
		// don't send anything
	})
}

func TestAssertNotificationSent_NotSent(t *testing.T) {
	tt := &testing.T{}

	AssertNotificationSent(tt, &fyne.Notification{}, func() {
		// don't send anything
	})
	assert.True(t, tt.Failed(), "notification assert should fail if no notification was sent")
}
