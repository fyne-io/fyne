package tutorials

import (
	"fmt"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
)

func scaleString(c fyne.Canvas) string {
	return fmt.Sprintf("%0.2f", c.Scale())
}

func prependTo(g *widget.Group, s string) {
	g.Prepend(widget.NewLabel(s))
}

func setScaleText(obj *widget.Label, win fyne.Window) {
	for obj.Visible() {
		obj.SetText(scaleString(win.Canvas()))

		time.Sleep(time.Second)
	}
}

// advancedScreen loads a panel that shows details and settings that are a bit
// more detailed than normally needed.
func advancedScreen(win fyne.Window) fyne.CanvasObject {
	scale := widget.NewLabel("")

	screen := widget.NewCard("Screen info", "", widget.NewForm(
		&widget.FormItem{Text: "Scale", Widget: scale},
	))

	go setScaleText(scale, win)

	label := widget.NewLabel("Just type...")
	generic := widget.NewGroupWithScroller("Generic")
	desk := widget.NewGroupWithScroller("Desktop")

	win.Canvas().SetOnTypedRune(func(r rune) {
		prependTo(generic, "Rune: "+string(r))
	})
	win.Canvas().SetOnTypedKey(func(ev *fyne.KeyEvent) {
		prependTo(generic, "Key : "+string(ev.Name))
	})
	if deskCanvas, ok := win.Canvas().(desktop.Canvas); ok {
		deskCanvas.SetOnKeyDown(func(ev *fyne.KeyEvent) {
			prependTo(desk, "KeyDown: "+string(ev.Name))
		})
		deskCanvas.SetOnKeyUp(func(ev *fyne.KeyEvent) {
			prependTo(desk, "KeyUp  : "+string(ev.Name))
		})
	}

	return widget.NewHBox(
		widget.NewVBox(screen,
			widget.NewButton("Custom Theme", func() {
				fyne.CurrentApp().Settings().SetTheme(newCustomTheme())
			}),
			widget.NewButton("Fullscreen", func() {
				win.SetFullScreen(!win.FullScreen())
			}),
		),
		fyne.NewContainerWithLayout(layout.NewBorderLayout(label, nil, nil, nil),
			label,
			fyne.NewContainerWithLayout(layout.NewGridLayout(2),
				generic, desk,
			),
		),
	)
}
