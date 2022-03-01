//go:build !ci && !legacy
// +build !ci,!legacy

package app

/*
#cgo LDFLAGS: -framework Foundation -framework UserNotifications
*/
import "C"
