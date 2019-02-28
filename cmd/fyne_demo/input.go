package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
)

// Input loads a window that shows how input events are handled.
func Input(app fyne.App) {
	win := app.NewWindow("Input")
	label := widget.NewLabel("Just type...")

	generic := widget.NewVBox()
	desk := widget.NewVBox()

	win.SetContent(fyne.NewContainerWithLayout(layout.NewBorderLayout(label, nil, nil, nil),
		label,
		fyne.NewContainerWithLayout(layout.NewGridLayout(2),
			widget.NewGroup("Generic", widget.NewScroller(generic)),
			widget.NewGroup("Desktop", widget.NewScroller(desk)),
		),
	))

	win.Canvas().SetOnTypedRune(func(r rune) {
		prependTo(generic, "Rune: "+string(r))
	})
	win.Canvas().SetOnTypedKey(func(ev *fyne.KeyEvent) {
		prependTo(generic, "Key : "+string(ev.Name))
	})
	win.Canvas().(desktop.Canvas).SetOnKeyDown(func(ev *fyne.KeyEvent) {
		prependTo(desk, "KeyDown: "+string(ev.Name))
	})
	win.Canvas().(desktop.Canvas).SetOnKeyUp(func(ev *fyne.KeyEvent) {
		prependTo(desk, "KeyUp  : "+string(ev.Name))
	})
	win.Resize(fyne.NewSize(300, 300))
	win.Show()
}

func prependTo(g *widget.Box, s string) {
	g.Prepend(widget.NewLabel(s))
}
