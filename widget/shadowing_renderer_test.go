package widget

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne"
)

func TestShadowingRenderer_Objects(t *testing.T) {
	tests := map[string]struct {
		level                elevationLevel
		wantPrependedObjects []fyne.CanvasObject
	}{
		"with shadow": {
			12,
			[]fyne.CanvasObject{newShadow(shadowAround, 12)},
		},
		"without shadow": {
			0,
			[]fyne.CanvasObject{},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			objects := []fyne.CanvasObject{NewLabel("A"), NewLabel("B")}
			r := newShadowingRenderer(objects, tt.level)
			assert.Equal(t, append(tt.wantPrependedObjects, objects...), r.Objects())

			otherObjects := []fyne.CanvasObject{NewLabel("X"), NewLabel("Y")}
			r.SetObjects(otherObjects)
			assert.Equal(t, append(tt.wantPrependedObjects, otherObjects...), r.Objects())
		})
	}
}
