package tutorials

import (
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// collectionScreen loads a tab panel for collection widgets
func collectionScreen(_ fyne.Window) fyne.CanvasObject {
	content := container.NewVBox(
		widget.NewLabelWithStyle("func Length() int", fyne.TextAlignLeading, fyne.TextStyle{Monospace: true}),
		widget.NewLabelWithStyle("func CreateItem() fyne.CanvasObject", fyne.TextAlignLeading, fyne.TextStyle{Monospace: true}),
		widget.NewLabelWithStyle("func UpdateItem(ListItemID, fyne.CanvasObject)", fyne.TextAlignLeading, fyne.TextStyle{Monospace: true}),
		widget.NewLabelWithStyle("func OnSelected(ListItemID)", fyne.TextAlignLeading, fyne.TextStyle{Monospace: true}),
		widget.NewLabelWithStyle("func OnUnselected(ListItemID)", fyne.TextAlignLeading, fyne.TextStyle{Monospace: true}))
	return container.NewCenter(content)
}

func makeGridWrapTab(_ fyne.Window) fyne.CanvasObject {
	data := make([]string, 1000)
	for i := range data {
		data[i] = "Test Item " + strconv.Itoa(i)
	}

	icon := widget.NewIcon(nil)
	label := widget.NewLabel("Select An Item From The List")
	vbox := container.NewVBox(icon, label)

	grid := widget.NewGridWrap(
		func() int {
			return len(data)
		},
		func() fyne.CanvasObject {
			text := widget.NewLabel("Template Object")
			text.Alignment = fyne.TextAlignCenter
			return container.NewVBox(container.NewPadded(widget.NewIcon(theme.DocumentIcon())), text)
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			item.(*fyne.Container).Objects[1].(*widget.Label).SetText(data[id])
		},
	)
	grid.OnSelected = func(id widget.ListItemID) {
		label.SetText(data[id])
		icon.SetResource(theme.DocumentIcon())
	}
	grid.OnUnselected = func(id widget.ListItemID) {
		label.SetText("Select An Item From The List")
		icon.SetResource(nil)
	}
	grid.Select(15)

	split := container.NewHSplit(grid, container.NewCenter(vbox))
	split.Offset = 0.6
	return split
}

func makeListTab(_ fyne.Window) fyne.CanvasObject {
	data := make([]string, 1000)
	for i := range data {
		data[i] = "Test Item " + strconv.Itoa(i)
	}

	icon := widget.NewIcon(nil)
	label := widget.NewLabel("Select An Item From The List")
	hbox := container.NewHBox(icon, label)

	list := widget.NewList(
		func() int {
			return len(data)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewIcon(theme.DocumentIcon()), widget.NewLabel("Template Object"))
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			if id == 5 || id == 6 {
				item.(*fyne.Container).Objects[1].(*widget.Label).SetText(data[id] + "\ntaller")
			} else {
				item.(*fyne.Container).Objects[1].(*widget.Label).SetText(data[id])
			}
		},
	)
	list.OnSelected = func(id widget.ListItemID) {
		label.SetText(data[id])
		icon.SetResource(theme.DocumentIcon())
	}
	list.OnUnselected = func(id widget.ListItemID) {
		label.SetText("Select An Item From The List")
		icon.SetResource(nil)
	}
	list.Select(125)
	list.SetItemHeight(5, 50)
	list.SetItemHeight(6, 50)

	return container.NewHSplit(list, container.NewCenter(hbox))
}

func makeTableTab(_ fyne.Window) fyne.CanvasObject {
	t := widget.NewTableWithHeaders(
		func() (int, int) { return 500, 150 },
		func() fyne.CanvasObject {
			return widget.NewLabel("Cell 000, 000")
		},
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			label := cell.(*widget.Label)
			switch id.Col {
			case 0:
				label.SetText("A longer cell")
			default:
				label.SetText(fmt.Sprintf("Cell %d, %d", id.Row+1, id.Col+1))
			}
		})
	t.SetColumnWidth(0, 102)
	t.SetRowHeight(2, 50)
	return t
}

func makeTreeTab(_ fyne.Window) fyne.CanvasObject {
	data := map[string][]string{
		"":  {"A"},
		"A": {"B", "D", "H", "J", "L", "O", "P", "S", "V", "Z"},
		"B": {"C"},
		"C": {"abc"},
		"D": {"E"},
		"E": {"F", "G"},
		"F": {"adef"},
		"G": {"adeg"},
		"H": {"I"},
		"I": {"ahi"},
		"O": {"ao"},
		"P": {"Q"},
		"Q": {"R"},
		"R": {"apqr"},
		"S": {"T"},
		"T": {"U"},
		"U": {"astu"},
		"V": {"W"},
		"W": {"X"},
		"X": {"Y"},
		"Y": {"zzz"},
		"Z": {},
	}

	childlen := 10000
	z := make([]string, childlen)
	for i := 0; i < childlen; i++ {
		z[i] = strconv.Itoa(i)
	}
	data["Z"] = z

	tree := widget.NewTreeWithStrings(data)
	tree.OnSelected = func(id string) {
		fmt.Println("Tree node selected:", id)
	}
	tree.OnUnselected = func(id string) {
		fmt.Println("Tree node unselected:", id)
	}
	tree.OpenBranch("A")
	tree.OpenBranch("D")
	tree.OpenBranch("E")
	tree.OpenBranch("L")
	tree.OpenBranch("M")
	return tree
}
