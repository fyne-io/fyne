//go:build linux
// +build linux

package theme

import (
	"io/ioutil"
	"os"
	"testing"

	"fyne.io/fyne/v2"
)

func setup() (tmp, home string) {
	// create a false home directory
	var err error
	tmp, err = ioutil.TempDir("", "fyne-test-")
	if err != nil {
		panic(err)
	}
	home = os.Getenv("HOME")
	os.Setenv("HOME", tmp)

	// creat a false KDE configuration
	if err = os.MkdirAll(tmp+"/.config", 0755); err != nil {
		panic(err)
	}
	content := []byte("[General]\nwidgetStyle=GTK")
	ioutil.WriteFile(tmp+"/.config/kdeglobals", content, 0644)

	return
}

func teardown(tmp, home string) {
	os.Unsetenv("XDG_CURRENT_DESKTOP")
	os.RemoveAll(tmp)
	os.Setenv("HOME", home)
}

// Test to load from desktop environment.
func TestLoadFromEnvironment(t *testing.T) {
	tmp, home := setup()
	defer teardown(tmp, home)

	// Set XDG_CURRENT_DESKTOP to "GNOME"
	envs := []string{"GNOME", "KDE", "FAKE"}
	app := fyne.CurrentApp()
	for _, env := range envs {
		// chante desktop environment
		os.Setenv("XDG_CURRENT_DESKTOP", env)
		app.Settings().SetTheme(FromDesktopEnvironment())

		// check if the theme is loaded
		current := app.Settings().Theme()
		// Check if the type of the theme is correct
		if current == nil {
			t.Error("Theme is nil")
		}
		switch env {
		case "GNOME":
			switch v := current.(type) {
			case *GnomeTheme:
				// OK
			default:
				t.Error("Theme is not GnomeTheme")
				t.Logf("Theme is %T\n", v)
			}
		case "KDE":
			switch v := current.(type) {
			case *KDETheme:
				// OK
			default:
				t.Error("Theme is not KDETheme")
				t.Logf("Theme is %T\n", v)
			}
		case "FAKE":
			if current != DefaultTheme() {
				t.Error("Theme is not DefaultTheme")
			}
		}

	}
}
