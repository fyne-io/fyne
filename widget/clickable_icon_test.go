package widget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestClickableIcon_Layout(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	for name, tt := range map[string]struct {
		resource fyne.Resource
	}{
		"empty": {},
		"resource": {
			resource: theme.CancelIcon(),
		},
	} {
		t.Run(name, func(t *testing.T) {
			icon := NewClickableIcon(tt.resource, func() {})

			window := test.NewWindow(&fyne.Container{Layout: layout.NewCenterLayout(), Objects: []fyne.CanvasObject{icon}})
			window.Resize(icon.MinSize().Max(fyne.NewSize(150, 200)))

			test.AssertRendersToMarkup(t, "clickable_icon/layout_"+name+".xml", window.Canvas())

			window.Close()
		})
	}
}

func TestClickableIcon_Tapped(t *testing.T) {
	tapped := make(chan bool)
	clickableIcon := NewClickableIcon(theme.CancelIcon(), func() {
		tapped <- true
	})

	go test.Tap(clickableIcon)
	func() {
		select {
		case <-tapped:
		case <-time.After(1 * time.Second):
			assert.Fail(t, "Timed out waiting for button tap")
		}
	}()
}
