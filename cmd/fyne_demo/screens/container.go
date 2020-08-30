package screens

import (
	"fmt"
	"image/color"
	"log"
	"net/url"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

// ContainerScreen loads a tab panel for containers and layouts
func ContainerScreen() fyne.CanvasObject {
	return widget.NewTabContainer(
		widget.NewTabItem("Accordion", makeAccordionTab()),
		widget.NewTabItem("Split", makeSplitTab()),
		widget.NewTabItem("Scroll", makeScrollTab()),
		widget.NewTabItem("Tree", makeTreeTab()),
		// layouts
		widget.NewTabItem("Border", makeBorderLayout()),
		widget.NewTabItem("Box", makeBoxLayout()),
		widget.NewTabItem("Center", makeCenterLayout()),
		widget.NewTabItem("Grid", makeGridLayout()),
	)
}

func makeAccordionTab() fyne.CanvasObject {
	link, err := url.Parse("https://fyne.io/")
	if err != nil {
		fyne.LogError("Could not parse URL", err)
	}
	ac := widget.NewAccordionContainer(
		widget.NewAccordionItem("A", widget.NewHyperlink("One", link)),
		widget.NewAccordionItem("B", widget.NewLabel("Two")),
		&widget.AccordionItem{
			Title:  "C",
			Detail: widget.NewLabel("Three"),
		},
	)
	ac.Append(widget.NewAccordionItem("D", &widget.Entry{Text: "Four"}))
	return ac
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

func makeCenterLayout() *fyne.Container {
	middle := widget.NewButton("CenterLayout", func() {})

	return fyne.NewContainerWithLayout(layout.NewCenterLayout(),
		middle)
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
	hlist := makeButtonList(20)
	vlist := makeButtonList(50)

	horiz := widget.NewHScrollContainer(widget.NewHBox(hlist...))
	vert := widget.NewVScrollContainer(widget.NewVBox(vlist...))

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

func makeTreeTab() fyne.CanvasObject {
	left := widget.NewTree()
	left.UseFileSystemIcons()
	left.OnLeafSelected = func(path []string) {
		log.Println("TreeLeafSelected:", path)
	}
	left.Add("A", "B", "C", "abc")
	left.Add("A", "D", "E", "F", "adef")
	left.Add("A", "D", "E", "G", "adeg")
	left.Add("A", "H", "I", "ahi")
	left.Add("A", "J", "K", "ajk")
	left.Add("A", "L", "M", "N", "almn")
	left.Add("A", "O", "ao")
	left.Add("A", "P", "Q", "R", "apqr")
	left.Add("A", "S", "T", "U", "astu")
	left.Add("A", "V", "W", "X", "Y", "Z", "avwxyz")

	right := widget.NewTree()
	right.UseArrowIcons()
	right.OnLeafSelected = func(path []string) {
		log.Println("TreeLeafSelected:", path)
	}
	right.Add("1", "2", "3", "1bc")
	right.Add("1", "4", "5", "6", "1456")
	right.Add("1", "4", "5", "7", "1457")
	right.Add("1", "8", "9", "189")
	right.Add("1", "10", "11", "11011")
	return widget.NewHBox(
		widget.NewVScrollContainer(left),
		widget.NewVScrollContainer(right),
	)
}
