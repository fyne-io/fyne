// +build !linux,!darwin,!windows

package fyne

import "log"

func (app *fyneApp) OpenURL(url string) {
	log.Fatalf("Unable to open url for unknown operating system")
}
