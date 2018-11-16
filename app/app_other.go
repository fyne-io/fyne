// +build ci !linux,!darwin,!windows,!freebsd,!openbsd,!netbsd

package app

import "log"

func (app *fyneApp) OpenURL(url string) {
	log.Printf("Unable to open url for unknown operating system")
}
