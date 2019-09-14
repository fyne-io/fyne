package main

import (
	"fmt"
	"os"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
)

var counter = 5
var lv *widget.ListView

func main() {
	a := app.New()
	w := a.NewWindow("title")
	w.Resize(fyne.NewSize(500, 600))
	w.SetContent(getcontent())
	w.ShowAndRun()
}

func getcontent() fyne.CanvasObject {
	lv = widget.NewVListView(onCreate, onBind, onGetCount)
	return lv
}

func onCreate(vh *widget.ViewHolder) fyne.CanvasObject {
	b := widget.NewButton("label", nil)
	vh.Add(b, "name")
	return b
}

func onBind(vh *widget.ViewHolder, pos int) {
	b := vh.GetButton("name")
	b.SetText(fmt.Sprint(pos, counter))
	b.SetOnTapped(func() {
		counter++
		fmt.Println(pos)
		lv.NotifyDataChange()
	})
}

func onGetCount() int {
	return counter
}
