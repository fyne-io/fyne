// +build !linux,!darwin,!windows

package fyneapp

import "log"

func (app *eflApp) OpenURL(url string) {
	log.Fatalf("Unable to open url for unknown operating system")
}
