package fyne

import "sync/atomic"

// themeRefreshDepth tracks active ApplyThemeTo passes.
var themeRefreshDepth atomic.Int32

// BeginThemeRefresh starts a theme refresh pass and returns a function that ends it.
// For internal use only.
func BeginThemeRefresh() func() {
	themeRefreshDepth.Add(1)
	return func() {
		themeRefreshDepth.Add(-1)
	}
}

// IsThemeRefreshing reports whether a theme refresh pass is active.
// For internal use only.
func IsThemeRefreshing() bool {
	return themeRefreshDepth.Load() > 0
}
