package app

import (
	"testing"

	"fyne.io/fyne"
	_ "fyne.io/fyne/test"
	"github.com/stretchr/testify/assert"
)

func TestDummyApp(t *testing.T) {
	app := New()

	app.Quit()
}

func TestCurrentApp(t *testing.T) {
	app := New()

	assert.Equal(t, app, fyne.CurrentApp())
}
