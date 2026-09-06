//go:build !ios

package app

// BeginPaint records that a frame's drawing commands are about to be issued.
// Only iOS needs this, so no-op here for everyone else.
func BeginPaint() {}
