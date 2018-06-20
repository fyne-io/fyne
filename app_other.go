// +build ci !linux,!darwin,!windows

package fyne

import "log"

func (app *fyneApp) OpenURL(url string) {
	log.Printf("Unable to open url for unknown operating system")
}
