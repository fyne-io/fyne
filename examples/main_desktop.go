// +build !ci

package main

import "github.com/fyne-io/fyne/api/app"
import "github.com/fyne-io/fyne/desktop"

func newApp() app.App {
	return fyneapp.NewApp()
}
