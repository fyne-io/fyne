// +build required

package glfw

// This file exists purely to prevent the golang toolchain from stripping
// away the c source directories and files when `go mod vendor` is used
// to populate a `vendor/` directory of a project depending on `go-gl/glfw`.
//
// How it works:
//  - every directory which only includes c source files receives a dummy.go file.
//  - every directory we want to preserve is included here as a _ import.
//  - this file is given a build to exclude it from the regular build.
import (
	// Prevent go tooling from stripping out the c source files.
	_ "github.com/go-gl/glfw/v3.3/glfw/glfw/deps"
	_ "github.com/go-gl/glfw/v3.3/glfw/glfw/include/GLFW"
	_ "github.com/go-gl/glfw/v3.3/glfw/glfw/src"
)
