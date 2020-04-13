package widget_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne"
	w "fyne.io/fyne/internal/widget"
	"fyne.io/fyne/widget"
)

func TestShadowingRenderer_Objects(t *testing.T) {
	tests := map[string]struct {
		level                w.ElevationLevel
		wantPrependedObjects []fyne.CanvasObject
	}{
		"with shadow": {
			12,
			[]fyne.CanvasObject{w.NewShadow(w.ShadowAround, 12)},
		},
		"without shadow": {
			0,
			[]fyne.CanvasObject{},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			objects := []fyne.CanvasObject{widget.NewLabel("A"), widget.NewLabel("B")}
			r := w.NewShadowingRenderer(objects, tt.level)
			assert.Equal(t, append(tt.wantPrependedObjects, objects...), r.Objects())

			otherObjects := []fyne.CanvasObject{widget.NewLabel("X"), widget.NewLabel("Y")}
			r.SetObjects(otherObjects)
			assert.Equal(t, append(tt.wantPrependedObjects, otherObjects...), r.Objects())
		})
	}
}
