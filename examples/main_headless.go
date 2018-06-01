// +build ci

package main

import "github.com/fyne-io/fyne/api/app"
import "github.com/fyne-io/fyne/test"

func newApp() app.App {
	return test.NewTestApp()
}
