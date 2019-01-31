// +build ci !linux,!darwin,!windows,!freebsd,!openbsd,!netbsd

package app

import (
	"log"
	"net/url"
)

func (app *fyneApp) OpenURL(url *url.URL) {
	log.Printf("Unable to open url for unknown operating system")
}
