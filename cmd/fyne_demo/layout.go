package main

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
)

func makeCell() fyne.CanvasObject {
	rect := canvas.NewRectangle(&color.RGBA{128, 128, 128, 255})
	rect.SetMinSize(fyne.NewSize(30, 30))
	return rect
}

func makeBorderLayout(o fyne.CanvasObject) *fyne.Container {
	top := makeCell()
	bottom := makeCell()
	left := makeCell()
	right := makeCell()
	middle := widget.NewVBox(
		widget.NewLabelWithStyle("BorderLayout", fyne.TextAlignCenter, fyne.TextStyle{}),
		o,
	)

	borderLayout := layout.NewBorderLayout(top, bottom, left, right)
	return fyne.NewContainerWithLayout(borderLayout,
		top, bottom, left, right, middle)
}

func makeBoxLayout(o fyne.CanvasObject) *fyne.Container {
	top := makeCell()
	bottom := makeCell()
	middle := widget.NewVBox(widget.NewLabel("BoxLayout"), o)
	center := makeCell()
	right := makeCell()

	col := fyne.NewContainerWithLayout(layout.NewVBoxLayout(),
		top, middle, bottom)

	return fyne.NewContainerWithLayout(layout.NewHBoxLayout(),
		col, center, right)
}

func makeFixedGridLayout(o fyne.CanvasObject) *fyne.Container {
	box1 := makeCell()
	box2 := widget.NewVBox(widget.NewLabel("FixedGrid"), o)
	box3 := makeCell()
	box4 := makeCell()

	return fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.NewSize(75, 75)),
		box1, box2, box3, box4)
}

func makeGridLayout(o fyne.CanvasObject) *fyne.Container {
	box1 := makeCell()
	box2 := widget.NewVBox(widget.NewLabel("Grid"), o)
	box3 := makeCell()
	box4 := makeCell()

	return fyne.NewContainerWithLayout(layout.NewGridLayout(2),
		box1, box2, box3, box4)
}

// Layout loads a window that shows the layouts available for a container
func Layout(app fyne.App) {
	w := app.NewWindow("Layout")

	borderButton := widget.NewButton("Move Tabs", nil)
	boxButton := widget.NewButton("Move Tabs", nil)
	fixedGridButton := widget.NewButton("Move Tabs", nil)
	gridButton := widget.NewButton("Move Tabs", nil)

	t := widget.NewTabContainer(
		widget.NewTabItem("Border", makeBorderLayout(borderButton)),
		widget.NewTabItem("Box", makeBoxLayout(boxButton)),
		widget.NewTabItem("Fixed Grid", makeFixedGridLayout(fixedGridButton)),
		widget.NewTabItem("Grid", makeGridLayout(gridButton)),
	)
	w.SetContent(t)

	l := widget.TabLocationTop
	onTapped := func() {
		switch l {
		case widget.TabLocationTop:
			l = widget.TabLocationLeading
		case widget.TabLocationLeading:
			l = widget.TabLocationBottom
		case widget.TabLocationBottom:
			l = widget.TabLocationTrailing
		default:
			l = widget.TabLocationTop
		}
		t.SetTabLocation(l)
	}
	borderButton.OnTapped = onTapped
	boxButton.OnTapped = onTapped
	fixedGridButton.OnTapped = onTapped
	gridButton.OnTapped = onTapped

	w.Show()
}
