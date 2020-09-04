package screens

import (
	"fmt"
	"image/color"
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
		widget.NewTabItem("Card", makeCardTab()),
		widget.NewTabItem("Split", makeSplitTab()),
		widget.NewTabItem("Scroll", makeScrollTab()),
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

func makeCardTab() fyne.CanvasObject {
	card1 := widget.NewCardContainer("Book a table", "Which time suits?",
		widget.NewRadio([]string{"6:30pm", "7:00pm", "7:45pm"}, func(string) {}))
	card2 := widget.NewCardContainer("With media", "No content, with image", nil)
	card2.Image = canvas.NewImageFromResource(theme.FyneLogo())
	card3 := widget.NewCardContainer("Title 3", "Subtitle", widget.NewCheck("Check me", func(bool) {}))
	card4 := widget.NewCardContainer("Title 4", "Another card", widget.NewLabel("Content"))
	return fyne.NewContainerWithLayout(layout.NewGridLayout(3), widget.NewVBox(card1, card4),
		widget.NewVBox(card2), widget.NewVBox(card3))
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
