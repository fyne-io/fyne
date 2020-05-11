package screens

import (
	"fmt"
	"net/url"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

func scaleString(c fyne.Canvas) string {
	return fmt.Sprintf("%0.2f", c.Scale())
}

func prependTo(g *widget.Group, s string) {
	g.Prepend(widget.NewLabel(s))
}

func setScaleText(obj *widget.Label, win fyne.Window) {
	for obj.Visible() {
		obj.SetText(scaleString(win.Canvas()))

		time.Sleep(time.Second)
	}
}

// AdvancedScreen loads a panel that shows details and settings that are a bit
// more detailed than normally needed.
func AdvancedScreen(win fyne.Window) fyne.CanvasObject {
	scale := widget.NewLabel("")

	screen := widget.NewGroup("Screen", widget.NewForm(
		&widget.FormItem{Text: "Scale", Widget: scale},
	))

	go setScaleText(scale, win)

	label := widget.NewLabel("Just type...")
	generic := widget.NewGroupWithScroller("Generic")
	desk := widget.NewGroupWithScroller("Desktop")

	win.Canvas().SetOnTypedRune(func(r rune) {
		prependTo(generic, "Rune: "+string(r))
	})
	win.Canvas().SetOnTypedKey(func(ev *fyne.KeyEvent) {
		prependTo(generic, "Key : "+string(ev.Name))
	})
	if deskCanvas, ok := win.Canvas().(desktop.Canvas); ok {
		deskCanvas.SetOnKeyDown(func(ev *fyne.KeyEvent) {
			prependTo(desk, "KeyDown: "+string(ev.Name))
		})
		deskCanvas.SetOnKeyUp(func(ev *fyne.KeyEvent) {
			prependTo(desk, "KeyUp  : "+string(ev.Name))
		})
	}

	// if we had chaining here, wouldnt need to create all this guff outside of the HBox
	themeUpdateEntry := widget.NewMultiLineEntry()
	link, _ := url.Parse("https://cimdalli.github.io/mui-theme-generator/")
	notifyError := func(err error) {
		fyne.CurrentApp().SendNotification(&fyne.Notification{
			Title:   "Error",
			Content: err.Error(),
		})
	}
	return widget.NewHBox(
		widget.NewVBox(screen,
			widget.NewHBox(
				widget.NewHyperlink("MaterialUI Themes", link),
				widget.NewButton("Paste", func() {
					if err := theme.Extend(win.Clipboard().Content()); err != nil {
						notifyError(err)
					}
				}).SetStyle(widget.SecondaryButton),
			),
			themeUpdateEntry,
			widget.NewButton("Extend Theme", func() {
				if err := theme.Extend(themeUpdateEntry.Text); err != nil {
					notifyError(err)
				}
				themeUpdateEntry.SetText("")
			}).SetStyle(widget.PrimaryButton),
			widget.NewButton("Custom Theme", func() {
				fyne.CurrentApp().Settings().SetTheme(newCustomTheme())
			}).SetStyle(widget.SecondaryButton),
			widget.NewButton("Fullscreen", func() {
				win.SetFullScreen(!win.FullScreen())
			}),
		),
		fyne.NewContainerWithLayout(layout.NewBorderLayout(label, nil, nil, nil),
			label,
			fyne.NewContainerWithLayout(layout.NewGridLayout(2),
				generic, desk,
			),
		),
	)
}
