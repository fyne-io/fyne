package examples

import "fmt"

import "github.com/fyne-io/fyne/app"
import "github.com/fyne-io/fyne/ui/widget"

import "github.com/mmcdole/gofeed"

const feedURL = "http://fyne.io/feed.xml"

var parent app.App

func parse(list *widget.List) {
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL(feedURL)

	for i, _ := range feed.Items {
		item := feed.Items[i] // keep a reference to the slices
		list.Append(widget.NewButton(item.Title, func() {
			parent.OpenURL(fmt.Sprintf("%s#about", item.Link))
		}))
	}
}

func Blog(app app.App) {
	parent = app
	w := app.NewWindow("Blog")
	list := widget.NewList(widget.NewLabel(feedURL))
	w.Canvas().SetContent(list)

	go parse(list)

	w.Show()
}
