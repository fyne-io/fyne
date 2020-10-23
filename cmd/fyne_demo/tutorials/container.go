package tutorials

import (
	"fmt"
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/container"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

// containerScreen loads a tab panel for containers
func containerScreen(_ fyne.Window) fyne.CanvasObject {
	content := container.NewBorder(
		widget.NewLabelWithStyle("Top", fyne.TextAlignCenter, fyne.TextStyle{}),
		widget.NewLabelWithStyle("Bottom", fyne.TextAlignCenter, fyne.TextStyle{}),
		widget.NewLabel("Left"),
		widget.NewLabel("Right"),
		widget.NewLabel("Border Container"))
	return container.NewCenter(content)
}

func makeAppTabsTab(_ fyne.Window) fyne.CanvasObject {
	return container.NewAppTabs(
		container.NewTabItem("Tab 1", widget.NewLabel("Content of tab 1")),
		container.NewTabItem("Tab 2", widget.NewLabel("Content of tab 2")),
		container.NewTabItem("Tab 3", widget.NewLabel("Content of tab 3")),
	)
}

func makeBorderLayout(_ fyne.Window) fyne.CanvasObject {
	top := makeCell()
	bottom := makeCell()
	left := makeCell()
	right := makeCell()
	middle := widget.NewLabelWithStyle("BorderLayout", fyne.TextAlignCenter, fyne.TextStyle{})

	return container.NewBorder(top, bottom, left, right, middle)
}

func makeBoxLayout(_ fyne.Window) fyne.CanvasObject {
	top := makeCell()
	bottom := makeCell()
	middle := widget.NewLabel("BoxLayout")
	center := makeCell()
	right := makeCell()

	col := container.NewVBox(top, middle, bottom)

	return container.NewHBox(col, center, right)
}

func makeButtonList(count int) []fyne.CanvasObject {
	var items []fyne.CanvasObject
	for i := 1; i <= count; i++ {
		index := i // capture
		items = append(items, widget.NewButton(fmt.Sprintf("Button %d", index), func() {
			fmt.Println("Tapped", index)
		}))
	}

	return items
}

func makeCell() fyne.CanvasObject {
	rect := canvas.NewRectangle(&color.NRGBA{128, 128, 128, 255})
	rect.SetMinSize(fyne.NewSize(30, 30))
	return rect
}

func makeCenterLayout(_ fyne.Window) fyne.CanvasObject {
	middle := widget.NewButton("CenterLayout", func() {})

	return container.NewCenter(middle)
}

func makeGridLayout(_ fyne.Window) fyne.CanvasObject {
	box1 := makeCell()
	box2 := widget.NewLabel("Grid")
	box3 := makeCell()
	box4 := makeCell()

	return container.NewGridWithColumns(2,
		box1, box2, box3, box4)
}

func makeScrollTab(_ fyne.Window) fyne.CanvasObject {
	hlist := makeButtonList(20)
	vlist := makeButtonList(50)

	horiz := container.NewHScroll(widget.NewHBox(hlist...))
	vert := container.NewVScroll(widget.NewVBox(vlist...))

	return container.NewAdaptiveGrid(2,
		container.NewBorder(horiz, nil, nil, nil, vert),
		makeScrollBothTab())
}

func makeScrollBothTab() fyne.CanvasObject {
	logo := canvas.NewImageFromResource(theme.FyneLogo())
	logo.SetMinSize(fyne.NewSize(800, 800))

	scroll := container.NewScroll(logo)
	scroll.Resize(fyne.NewSize(400, 400))

	return scroll
}

func makeSplitTab(_ fyne.Window) fyne.CanvasObject {
	left := widget.NewMultiLineEntry()
	left.Wrapping = fyne.TextWrapWord
	left.SetText("Long text is looooooooooooooong")
	right := container.NewVSplit(
		widget.NewLabel("Label"),
		widget.NewButton("Button", func() { fmt.Println("button tapped!") }),
	)
	return container.NewHSplit(container.NewVScroll(left), right)
}
