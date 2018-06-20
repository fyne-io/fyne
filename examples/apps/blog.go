package apps

import "log"
import "fmt"

import "github.com/fyne-io/fyne"
import "github.com/fyne-io/fyne/widget"

import "github.com/mmcdole/gofeed"

const feedURL = "http://fyne.io/feed.xml"

var parent fyne.App

func parse(list *widget.List) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(feedURL)

	if err != nil {
		log.Println("Unable to load feed!")
		return
	}

	for i := range feed.Items {
		item := feed.Items[i] // keep a reference to the slices
		list.Append(widget.NewButton(item.Title, func() {
			parent.OpenURL(fmt.Sprintf("%s#about", item.Link))
		}))
	}
}

// Blog loads a blog example window for the specified app context
func Blog(app fyne.App) {
	parent = app
	w := app.NewWindow("Blog")
	list := widget.NewList(widget.NewLabel(feedURL))
	w.SetContent(list)

	go parse(list)

	w.Show()
}
