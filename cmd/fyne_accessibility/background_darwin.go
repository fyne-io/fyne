//go:build darwin && !ci && !ios && !mobile && !test_web_driver && !tinygo

package main

import "github.com/go-gl/glfw/v3.4/glfw"

func configureBackground() error {
	// NewWindow initializes GLFW; Show creates the native window with these hints.
	glfw.WindowHint(glfw.Focused, glfw.False)
	glfw.WindowHint(glfw.FocusOnShow, glfw.False)
	glfw.WindowHint(glfw.MousePassthrough, glfw.True)
	return nil
}
