package gl

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/internal/cache"
)

func FontCacheEntryForText(t *canvas.Text, canvas fyne.Canvas) cache.FontCacheEntry {
	custom := ""
	if t.FontSource != nil {
		custom = t.FontSource.Name()
	}
	ent := cache.FontCacheEntry{Color: t.Color, Canvas: canvas}
	ent.Text = t.Text
	ent.Size = t.TextSize
	ent.Style = t.TextStyle
	ent.Source = custom
	return ent
}
