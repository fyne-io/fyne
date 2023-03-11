// +build required

package gles2

// This file exists purely to prevent the golang toolchain from stripping
// away the c source directories and files when `go mod vendor` is used
// to populate a `vendor/` directory of a project depending on `go-gl/go`.
//
// How it works:
//  - every directory which only includes c source files receives a dummy.go file.
//  - every directory we want to preserve is included here as a _ import.
//  - this file is given a build tag to exclude it from the regular build.
import (
	// Prevent go tooling from stripping out the c source files.
	_ "github.com/go-gl/gl/v3.1/gles2/KHR"
)
