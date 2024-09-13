package app

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	_ "fyne.io/fyne/v2/test"
)

func TestDummyApp(t *testing.T) {
	app := NewWithID("io.fyne.test")

	app.Quit()
}

func TestCurrentApp(t *testing.T) {
	app := NewWithID("io.fyne.test")

	assert.Equal(t, app, fyne.CurrentApp())
}

func TestFyneApp_UniqueID(t *testing.T) {
	appID := "io.fyne.test"
	app := NewWithID(appID)

	assert.Equal(t, appID, app.UniqueID())

	meta.ID = "fakedInject"
	app = New()

	assert.Equal(t, meta.ID, app.UniqueID())
}

func TestFyneApp_SetIcon(t *testing.T) {
	app := NewWithID("io.fyne.test")

	metaIcon := &fyne.StaticResource{StaticName: "Metadata"}
	SetMetadata(fyne.AppMetadata{
		Icon: metaIcon,
	})

	assert.Equal(t, metaIcon, app.Icon())

	setIcon := &fyne.StaticResource{StaticName: "Test"}
	app.SetIcon(setIcon)

	assert.Equal(t, setIcon, app.Icon())
}

func TestFynaApp_Clipboard(t *testing.T) {
	app := test.NewTempApp(t)
	test.NewTempWindow(t, nil)

	text := "My content from test window"
	cb := app.Clipboard()

	cliboardContent := cb.Content()
	if cliboardContent != "" {
		// Current environment has some content stored in clipboard,
		// set temporary to an empty string to allow test and restore later.
		cb.SetContent("")
	}

	assert.Empty(t, cb.Content())

	cb.SetContent(text)
	assert.Equal(t, text, cb.Content())

	// Restore clipboardContent, if any
	cb.SetContent(cliboardContent)
}
