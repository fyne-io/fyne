//go:build !debug
// +build !debug

package internal

// BuildTypeIsDebug is true if built with -tags debug. This value exists as a constant
// so the Go compiler can deadcode eliminate reflect from non-debug builds.
const BuildTypeIsDebug = false
