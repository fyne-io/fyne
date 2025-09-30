//go:build !tamago && !noos

package test

import (
	"testing"

	"fyne.io/fyne/v2"
)

// NewTempApp returns a new dummy app and tears it down at the end of the test.
//
// Since: 2.5
func NewTempApp(t testing.TB) fyne.App {
	app := NewApp()
	t.Cleanup(func() { NewApp() })
	return app
}
