package screens

import (
	"fmt"
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

// ContainerScreen loads a tab panel for containers and layouts
func ContainerScreen() fyne.CanvasObject {
	return widget.NewTabContainer(
		widget.NewTabItem("Split", makeSplitTab()),
		widget.NewTabItem("Scroll", makeScrollTab()),
		// layouts
		widget.NewTabItem("Border", makeBorderLayout()),
		widget.NewTabItem("Box", makeBoxLayout()),
		widget.NewTabItem("Grid", makeGridLayout()),
	)
}

func makeBorderLayout() *fyne.Container {
	top := makeCell()
	bottom := makeCell()
	left := makeCell()
	right := makeCell()
	middle := widget.NewLabelWithStyle("BorderLayout", fyne.TextAlignCenter, fyne.TextStyle{})

	borderLayout := layout.NewBorderLayout(top, bottom, left, right)
	return fyne.NewContainerWithLayout(borderLayout,
		top, bottom, left, right, middle)
}

func makeBoxLayout() *fyne.Container {
	top := makeCell()
	bottom := makeCell()
	middle := widget.NewLabel("BoxLayout")
	center := makeCell()
	right := makeCell()

	col := fyne.NewContainerWithLayout(layout.NewVBoxLayout(),
		top, middle, bottom)

	return fyne.NewContainerWithLayout(layout.NewHBoxLayout(),
		col, center, right)
}

func makeCell() fyne.CanvasObject {
	rect := canvas.NewRectangle(&color.RGBA{128, 128, 128, 255})
	rect.SetMinSize(fyne.NewSize(30, 30))
	return rect
}

func makeGridLayout() *fyne.Container {
	box1 := makeCell()
	box2 := widget.NewLabel("Grid")
	box3 := makeCell()
	box4 := makeCell()

	return fyne.NewContainerWithLayout(layout.NewGridLayout(2),
		box1, box2, box3, box4)
}

func makeScrollTab() fyne.CanvasObject {
	list := widget.NewHBox()
	list2 := widget.NewVBox()

	for i := 1; i <= 20; i++ {
		index := i // capture
		list.Append(widget.NewButton(fmt.Sprintf("Button %d", index), func() {
			fmt.Println("Tapped", index)
		}))
		list2.Append(widget.NewButton(fmt.Sprintf("Button %d", index), func() {
			fmt.Println("Tapped", index)
		}))
	}

	horiz := widget.NewHScrollContainer(list)
	vert := widget.NewVScrollContainer(list2)

	return fyne.NewContainerWithLayout(layout.NewAdaptiveGridLayout(2),
		fyne.NewContainerWithLayout(layout.NewBorderLayout(horiz, nil, nil, nil), horiz, vert),
		makeScrollBothTab())
}

func makeScrollBothTab() fyne.CanvasObject {
	logo := canvas.NewImageFromResource(theme.FyneLogo())
	logo.SetMinSize(fyne.NewSize(800, 800))

	scroll := widget.NewScrollContainer(logo)
	scroll.Resize(fyne.NewSize(400, 400))

	return scroll
}

func makeSplitTab() fyne.CanvasObject {
	left := widget.NewMultiLineEntry()
	left.Wrapping = fyne.TextWrapWord
	left.SetText("Long text is looooooooooooooong")
	right := widget.NewVSplitContainer(
		widget.NewLabel("Label"),
		widget.NewButton("Button", func() { fmt.Println("button tapped!") }),
	)
	return widget.NewHSplitContainer(widget.NewVScrollContainer(left), right)
}
