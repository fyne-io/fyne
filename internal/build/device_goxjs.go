//go:build js || wasm
// +build js wasm

package build

import (
	"regexp"
	"syscall/js"
)

var isMobile = regexp.MustCompile("Android|BlackBerry|iPhone|iPad|iPod|Opera Mini|IEMobile")

var userAgent = js.Global().Get("navigator").Get("userAgent").String()
var isMobile = isMobile.MatchString(userAgent)

// IsMobile returns true if running on a mobile device.
func IsMobile() bool {
	return isMobile
}
