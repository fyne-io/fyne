package tutorials

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

func windowScreen(_ fyne.Window) fyne.CanvasObject {
	var visibilityWindow fyne.Window = nil
	var visibilityState bool = false

	windowGroup := container.NewVBox(
		widget.NewButton("New window", func() {
			w := fyne.CurrentApp().NewWindow("Hello")
			w.SetContent(widget.NewLabel("Hello World!"))
			w.Show()
		}),
		widget.NewButton("Fixed size window", func() {
			w := fyne.CurrentApp().NewWindow("Fixed")
			w.SetContent(container.NewCenter(widget.NewLabel("Hello World!")))

			w.Resize(fyne.NewSize(240, 180))
			w.SetFixedSize(true)
			w.Show()
		}),
		widget.NewButton("Toggle between fixed/not fixed window size", func() {
			w := fyne.CurrentApp().NewWindow("Toggle fixed size")
			w.SetContent(container.NewCenter(widget.NewCheck("Fixed size", func(toggle bool) {
				if toggle {
					w.Resize(fyne.NewSize(240, 180))
				}
				w.SetFixedSize(toggle)
			})))
			w.Show()
		}),
		widget.NewButton("Centered window", func() {
			w := fyne.CurrentApp().NewWindow("Central")
			w.SetContent(container.NewCenter(widget.NewLabel("Hello World!")))

			w.CenterOnScreen()
			w.Show()
		}),
		widget.NewButton("Show/Hide window", func() {
			if visibilityWindow == nil {
				visibilityWindow = fyne.CurrentApp().NewWindow("Hello")
				visibilityWindow.SetContent(widget.NewLabel("Hello World!"))
				visibilityWindow.SetOnClosed(func() {
					visibilityWindow = nil
				})
			}
			if visibilityState {
				visibilityWindow.Hide()
			} else {
				visibilityWindow.Show()
			}
			visibilityState = !visibilityState
		}))

	drv := fyne.CurrentApp().Driver()
	if drv, ok := drv.(desktop.Driver); ok {
		windowGroup.Objects = append(windowGroup.Objects,
			widget.NewButton("Splash Window (only use on start)", func() {
				w := drv.CreateSplashWindow()
				w.SetContent(widget.NewLabelWithStyle("Hello World!\n\nMake a splash!",
					fyne.TextAlignCenter, fyne.TextStyle{Bold: true}))
				w.Show()

				go func() {
					time.Sleep(time.Second * 3)
					w.Close()
				}()
			}))
	}

	otherGroup := widget.NewCard("Other", "",
		widget.NewButton("Notification", func() {
			fyne.CurrentApp().SendNotification(&fyne.Notification{
				Title:   "Fyne Demo",
				Content: "Testing notifications...",
			})
		}))

	return container.NewVBox(widget.NewCard("Windows", "", windowGroup), otherGroup)
}
