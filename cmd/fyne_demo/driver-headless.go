// +build ci

package main

import "github.com/fyne-io/fyne"
import "github.com/fyne-io/fyne/test"

// NewApp creates a new headless app to test the examples code
func NewApp() fyne.App {
	return test.NewApp()
}
