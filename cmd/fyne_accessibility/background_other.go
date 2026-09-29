//go:build !darwin || ci || ios || mobile || test_web_driver || tinygo

package main

import "errors"

func configureBackground() error {
	return errors.New("-background requires the native macOS desktop driver (without ci, mobile, or web build tags)")
}
