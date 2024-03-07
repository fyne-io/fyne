package widget_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
	w "fyne.io/fyne/v2/internal/widget"
	"fyne.io/fyne/v2/widget"
)

func TestShadowingRenderer_Objects(t *testing.T) {
	tests := map[string]struct {
		level                w.ElevationLevel
		wantPrependedObjects []fyne.CanvasObject
	}{
		"with shadow": {
			12,
			[]fyne.CanvasObject{},
		},
		"without shadow": {
			0,
			[]fyne.CanvasObject{},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			shadowIndex := 0
			if tt.level > 0 {
				shadowIndex = 1 // Shadow pointers are not the same. Avoid comparing.
			}

			objects := []fyne.CanvasObject{widget.NewLabel("A"), widget.NewLabel("B")}
			r := w.NewShadowingRenderer(objects, tt.level)
			assert.Equal(t, append(tt.wantPrependedObjects, objects...), r.Objects()[shadowIndex:])

			otherObjects := []fyne.CanvasObject{widget.NewLabel("X"), widget.NewLabel("Y")}
			r.SetObjects(otherObjects)
			assert.Equal(t, append(tt.wantPrependedObjects, otherObjects...), r.Objects()[shadowIndex:])
		})
	}
}
