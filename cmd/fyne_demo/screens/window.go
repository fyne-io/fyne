package screens

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/storage"
	"fyne.io/fyne/widget"
)

func confirmCallback(response bool) {
	fmt.Println("Responded with", response)
}

func fileOpened(f fyne.URIReadCloser) {
	if f == nil {
		log.Println("Cancelled")
		return
	}

	ext := f.URI().Extension()
	if ext == ".png" {
		showImage(f)
	} else if ext == ".txt" {
		showText(f)
	}
	err := f.Close()
	if err != nil {
		fyne.LogError("Failed to close stream", err)
	}
}

func fileSaved(f fyne.URIWriteCloser) {
	if f == nil {
		log.Println("Cancelled")
		return
	}

	log.Println("Save to...", f.URI())
}

func loadImage(f fyne.URIReadCloser) *canvas.Image {
	data, err := ioutil.ReadAll(f)
	if err != nil {
		fyne.LogError("Failed to load image data", err)
		return nil
	}
	res := fyne.NewStaticResource(f.Name(), data)

	return canvas.NewImageFromResource(res)
}

func showImage(f fyne.URIReadCloser) {
	img := loadImage(f)
	if img == nil {
		return
	}
	img.FillMode = canvas.ImageFillOriginal

	w := fyne.CurrentApp().NewWindow(f.Name())
	w.SetContent(widget.NewScrollContainer(img))
	w.Resize(fyne.NewSize(320, 240))
	w.Show()
}

func loadText(f fyne.URIReadCloser) string {
	data, err := ioutil.ReadAll(f)
	if err != nil {
		fyne.LogError("Failed to load text data", err)
		return ""
	}
	if data == nil {
		return ""
	}

	return string(data)
}

func showText(f fyne.URIReadCloser) {
	text := widget.NewLabel(loadText(f))
	text.Wrapping = fyne.TextWrapWord

	w := fyne.CurrentApp().NewWindow(f.Name())
	w.SetContent(widget.NewScrollContainer(text))
	w.Resize(fyne.NewSize(320, 240))
	w.Show()
}

func loadDialogGroup(win fyne.Window) *widget.Group {
	return widget.NewGroup("Dialogs",
		widget.NewButton("Info", func() {
			dialog.ShowInformation("Information", "You should know this thing...", win)
		}),
		widget.NewButton("Error", func() {
			err := errors.New("a dummy error message")
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
		widget.NewButton("ProgressInfinite", func() {
			prog := dialog.NewProgressInfinite("MyProgress", "Closes after 5 seconds...", win)

			go func() {
				time.Sleep(time.Second * 5)
				prog.Hide()
			}()

			prog.Show()
		}),
		widget.NewButton("File Open With Filter (.txt or .png)", func() {
			fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
				if err == nil && reader == nil {
					return
				}
				if err != nil {
					dialog.ShowError(err, win)
					return
				}

				fileOpened(reader)
			}, win)
			fd.SetFilter(storage.NewExtensionFileFilter([]string{".png", ".txt"}))
			fd.Show()
		}),
		widget.NewButton("File Save", func() {
			dialog.ShowFileSave(func(writer fyne.URIWriteCloser, err error) {
				if err != nil {
					dialog.ShowError(err, win)
					return
				}

				fileSaved(writer)
			}, win)
		}),
		widget.NewButton("Custom Dialog (Login Form)", func() {
			username := widget.NewEntry()
			password := widget.NewPasswordEntry()
			content := widget.NewForm(widget.NewFormItem("Username", username),
				widget.NewFormItem("Password", password))

			dialog.ShowCustomConfirm("Login...", "Log In", "Cancel", content, func(b bool) {
				if !b {
					return
				}

				log.Println("Please Authenticate", username.Text, password.Text)
			}, win)
		}),
	)
}

func loadWindowGroup() fyne.Widget {
	windowGroup := widget.NewGroup("Windows",
		widget.NewButton("New window", func() {
			w := fyne.CurrentApp().NewWindow("Hello")
			w.SetContent(widget.NewLabel("Hello World!"))
			w.Show()
		}),
		widget.NewButton("Fixed size window", func() {
			w := fyne.CurrentApp().NewWindow("Fixed")
			w.SetContent(fyne.NewContainerWithLayout(layout.NewCenterLayout(), widget.NewLabel("Hello World!")))

			w.Resize(fyne.NewSize(240, 180))
			w.SetFixedSize(true)
			w.Show()
		}),
		widget.NewButton("Centered window", func() {
			w := fyne.CurrentApp().NewWindow("Central")
			w.SetContent(fyne.NewContainerWithLayout(layout.NewCenterLayout(), widget.NewLabel("Hello World!")))

			w.CenterOnScreen()
			w.Show()
		}))

	drv := fyne.CurrentApp().Driver()
	if drv, ok := drv.(desktop.Driver); ok {
		windowGroup.Append(
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

	otherGroup := widget.NewGroup("Other",
		widget.NewButton("Notification", func() {
			fyne.CurrentApp().SendNotification(&fyne.Notification{
				Title:   "Fyne Demo",
				Content: "Testing notifications...",
			})
		}))

	return widget.NewVBox(windowGroup, otherGroup)
}

// DialogScreen loads a panel that lists the dialog windows that can be tested.
func DialogScreen(win fyne.Window) fyne.CanvasObject {
	return fyne.NewContainerWithLayout(layout.NewAdaptiveGridLayout(2),
		widget.NewScrollContainer(loadDialogGroup(win)),
		widget.NewScrollContainer(loadWindowGroup()))
}
