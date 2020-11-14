package tutorials

import (
	"fmt"

	"fyne.io/fyne"
	"fyne.io/fyne/container"
	internalWidget "fyne.io/fyne/internal/widget"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
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

func makeListTab(_ fyne.Window) fyne.CanvasObject {
	var data []string
	for i := 0; i < 1000; i++ {
		data = append(data, fmt.Sprintf("Test Item %d", i))
	}

	icon := widget.NewIcon(nil)
	label := widget.NewLabel("Select An Item From The List")
	hbox := fyne.NewContainerWithLayout(layout.NewHBoxLayout(), icon, label)

	list := widget.NewList(
		func() int {
			return len(data)
		},
		func() fyne.CanvasObject {
			return fyne.NewContainerWithLayout(layout.NewHBoxLayout(), widget.NewIcon(theme.DocumentIcon()), widget.NewLabel("Template Object"))
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			item.(*fyne.Container).Objects[1].(*widget.Label).SetText(data[id])
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
	return widget.NewHSplitContainer(list, fyne.NewContainerWithLayout(layout.NewCenterLayout(), hbox))
}

func makeTableTab(_ fyne.Window) fyne.CanvasObject {
	t := widget.NewTable(
		func() (int, int) { return 500, 150 },
		func() fyne.CanvasObject {
			return widget.NewLabel("Cell 000, 000")
		},
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			label := cell.(*widget.Label)
			switch id.Col {
			case 0:
				label.SetText(fmt.Sprintf("%d", id.Row+1))
			case 1:
				label.SetText("A longer cell")
			default:
				label.SetText(fmt.Sprintf("Cell %d, %d", id.Row+1, id.Col+1))
			}
		})
	t.SetColumnWidth(0, 34)
	t.SetColumnWidth(1, 102)
	return t
}

func makeTreeTab(_ fyne.Window) fyne.CanvasObject {
	data := make(map[string][]string)
	internalWidget.AddTreePath(data, "A", "B", "C", "abc")
	internalWidget.AddTreePath(data, "A", "D", "E", "F", "adef")
	internalWidget.AddTreePath(data, "A", "D", "E", "G", "adeg")
	internalWidget.AddTreePath(data, "A", "H", "I", "ahi")
	internalWidget.AddTreePath(data, "A", "J", "K", "ajk")
	internalWidget.AddTreePath(data, "A", "L", "M", "N", "almn")
	internalWidget.AddTreePath(data, "A", "O", "ao")
	internalWidget.AddTreePath(data, "A", "P", "Q", "R", "apqr")
	internalWidget.AddTreePath(data, "A", "S", "T", "U", "astu")
	internalWidget.AddTreePath(data, "A", "V", "W", "X", "Y", "Z", "avwxyz")

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
