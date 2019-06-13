package main

import (
	"errors"
	"fmt"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/widget"
)

// Dialogs loads a window that lists the dialog windows that can be tested.
func Dialogs(app fyne.App) {
	win := app.NewWindow("Dialogs")
	win.Resize(fyne.NewSize(400, 300))

	win.SetContent(widget.NewVBox(
		widget.NewButton("Info", func() {
			dialog.ShowInformation("Information", "You should know this thing...", win)
		}),
		widget.NewButton("Error", func() {
			err := errors.New("A dummy error message")
			dialog.ShowError(err, win)
		}),
		widget.NewButton("Confirm", func() {
			cnf := dialog.NewConfirm("Confirmation", "Are you enjoying this demo?", confirmCallback, win)
			cnf.SetDismissText("Nah")
			cnf.SetConfirmText("Oh Yes!")
			cnf.Show()
		}),
		widget.NewButton("Progress", func() {
			prog := dialog.NewProgress("MyProgress", "Nearly there...", win)

			go func() {
				num := 0.0
				for num < 1.0 {
					time.Sleep(50 * time.Millisecond)
					prog.SetValue(num)
					num += 0.01
				}

				prog.SetValue(1)
				prog.Hide()
			}()

			prog.Show()
		}),
		widget.NewButton("Custom", func() {
			content := widget.NewEntry()
			content.SetPlaceHolder("Type something here")
			content.OnChanged = func(text string) {
				fmt.Println("Entered", text)
			}
			dialog.ShowCustom("Custom dialog", "Done", content, win)
		}),
	))
	win.Show()
}
