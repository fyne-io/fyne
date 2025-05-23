package tutorials

import (
	"fmt"
	"image/color"
	"log"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/cmd/fyne_demo/data"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
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
	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon("Tab 1", theme.HomeIcon(), widget.NewLabel("Content of tab 1")),
		container.NewTabItemWithIcon("Tab 2 bigger", theme.ComputerIcon(), widget.NewLabel("Content of tab 2")),
		container.NewTabItemWithIcon("Tab 3", theme.MediaVideoIcon(), widget.NewLabel("Content of tab 3")),
	)
	for i := 4; i <= 12; i++ {
		tabs.Append(container.NewTabItem(fmt.Sprintf("Tab %d", i), widget.NewLabel(fmt.Sprintf("Content of tab %d", i))))
	}
	locations := makeTabLocationSelect(tabs.SetTabLocation)
	return container.NewBorder(locations, nil, nil, nil, tabs)
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
		items = append(items, widget.NewButton("Button "+strconv.Itoa(index), func() {
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

func makeDocTabsTab(_ fyne.Window) fyne.CanvasObject {
	tabs := container.NewDocTabs(
		container.NewTabItem("Doc 1", widget.NewLabel("Content of document 1")),
		container.NewTabItem("Doc 2 bigger", widget.NewLabel("Content of document 2")),
		container.NewTabItem("Doc 3", widget.NewLabel("Content of document 3")),
	)
	i := 3
	tabs.CreateTab = func() *container.TabItem {
		i++
		return container.NewTabItem(fmt.Sprintf("Doc %d", i), widget.NewLabel(fmt.Sprintf("Content of document %d", i)))
	}
	locations := makeTabLocationSelect(tabs.SetTabLocation)
	return container.NewBorder(locations, nil, nil, nil, tabs)
}

func makeGridLayout(_ fyne.Window) fyne.CanvasObject {
	box1 := makeCell()
	box2 := widget.NewLabel("Grid")
	box3 := makeCell()
	box4 := makeCell()

	return container.NewGridWithColumns(2,
		box1, box2, box3, box4)
}

func makeInnerWindowTab(_ fyne.Window) fyne.CanvasObject {
	label := widget.NewLabel("Window content for inner demo")
	win1 := container.NewInnerWindow("Inner Demo", container.NewVBox(
		label,
		widget.NewButton("Tap Me", func() {
			label.SetText("Tapped")
		})))
	win1.Icon = data.FyneLogo
	win1.OnMaximized = func() {
		log.Println("Should maximize here")
	}
	win1.OnMinimized = func() {
		log.Println("Should minimize here")
	}

	win2 := container.NewInnerWindow("Inner2", widget.NewLabel("Win 2"))

	multi := container.NewMultipleWindows()
	multi.Windows = []*container.InnerWindow{win1, win2}
	return multi
}

func makeScrollTab(_ fyne.Window) fyne.CanvasObject {
	hlist := makeButtonList(20)
	vlist := makeButtonList(50)

	horiz := container.NewHScroll(container.NewHBox(hlist...))
	vert := container.NewVScroll(container.NewVBox(vlist...))

	return container.NewAdaptiveGrid(2,
		container.NewBorder(horiz, nil, nil, nil, vert),
		makeScrollBothTab())
}

func makeScrollBothTab() fyne.CanvasObject {
	logo := canvas.NewImageFromResource(data.FyneLogo)
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

func makeTabLocationSelect(callback func(container.TabLocation)) *widget.Select {
	locations := widget.NewSelect([]string{"Top", "Bottom", "Leading", "Trailing"}, func(s string) {
		callback(map[string]container.TabLocation{
			"Top":      container.TabLocationTop,
			"Bottom":   container.TabLocationBottom,
			"Leading":  container.TabLocationLeading,
			"Trailing": container.TabLocationTrailing,
		}[s])
	})
	locations.SetSelected("Top")
	return locations
}
