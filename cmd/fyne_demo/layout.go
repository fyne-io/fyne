package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
)

func makeBorderLayout() *fyne.Container {
	top := widget.NewEntry()
	bottom := widget.NewEntry()
	left := widget.NewEntry()
	right := widget.NewEntry()
	middle := &widget.Label{
		Text:      "BorderLayout",
		Alignment: fyne.TextAlignCenter,
	}

	borderLayout := layout.NewBorderLayout(top, bottom, left, right)
	return fyne.NewContainerWithLayout(borderLayout,
		top, bottom, left, right, middle)
}

func makeBoxLayout() *fyne.Container {
	top := widget.NewEntry()
	bottom := widget.NewEntry()
	middle := widget.NewLabel("BoxLayout")
	center := widget.NewEntry()
	right := widget.NewEntry()

	col := fyne.NewContainerWithLayout(layout.NewVBoxLayout(),
		top, middle, bottom)

	return fyne.NewContainerWithLayout(layout.NewHBoxLayout(),
		col, center, right)
}

func makeFixedGridLayout() *fyne.Container {
	box1 := widget.NewEntry()
	box2 := widget.NewLabel("FixedGrid")
	box3 := widget.NewEntry()
	box4 := widget.NewEntry()

	return fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.NewSize(75, 75)),
		box1, box2, box3, box4)
}

func makeGridLayout() *fyne.Container {
	box1 := widget.NewEntry()
	box2 := widget.NewLabel("Grid")
	box3 := widget.NewEntry()
	box4 := widget.NewEntry()

	return fyne.NewContainerWithLayout(layout.NewGridLayout(2),
		box1, box2, box3, box4)
}

// Layout loads a window that shows the layouts available for a container
func Layout(app fyne.App) {
	w := app.NewWindow("Layout")

	w.SetContent(widget.NewTabContainer(
		widget.NewTabItem("Border", makeBorderLayout()),
		widget.NewTabItem("Box", makeBoxLayout()),
		widget.NewTabItem("Fixed Grid", makeFixedGridLayout()),
		widget.NewTabItem("Grid", makeGridLayout()),
	))

	w.Show()
}
