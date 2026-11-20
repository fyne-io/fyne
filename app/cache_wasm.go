//go:build wasm

package app

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"
)

func rootCacheDir(_ fyne.App) string {
	// Only used by the (desktop-only) icon cache helpers, which are not called
	// on wasm. Keeping it defined avoids build errors in shared code.
	return "/fyne/cache"
}

func rootCacheURI(_ fyne.App) fyne.URI {
	// IndexedDB is scoped per browser origin, so an app ID is not required here.
	// Using a fixed path also keeps the cache across restarts when app.New() was used.
	root, _ := storage.ParseURI("idbfile:///fyne/cache")
	return root
}
