package tutorials

import (
	"errors"
	"fmt"
	"image/color"
	"io/ioutil"
	"log"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/container"
	"fyne.io/fyne/data/validation"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/storage"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

func confirmCallback(response bool) {
	fmt.Println("Responded with", response)
}

func colorPicked(c color.Color, w fyne.Window) {
	log.Println("Color picked:", c)
	rectangle := canvas.NewRectangle(c)
	size := 2 * theme.IconInlineSize()
	rectangle.SetMinSize(fyne.NewSize(size, size))
	dialog.ShowCustom("Color Picked", "Ok", rectangle, w)
}

// dialogScreen loads demos of the dialogs we support
func dialogScreen(win fyne.Window) fyne.CanvasObject {
	return container.NewVScroll(container.NewVBox(
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
		widget.NewButton("Progress (deprecated)", func() {
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
		widget.NewButton("ProgressInfinite (deprecated)", func() {
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
		widget.NewButton("Folder Open", func() {
			dialog.ShowFolderOpen(func(list fyne.ListableURI, err error) {
				if err != nil {
					dialog.ShowError(err, win)
					return
				}
				if list == nil {
					return
				}

				children, err := list.List()
				if err != nil {
					dialog.ShowError(err, win)
					return
				}
				out := fmt.Sprintf("Folder %s (%d children):\n%s", list.Name(), len(children), list.String())
				dialog.ShowInformation("Folder Open", out, win)
			}, win)
		}),
		widget.NewButton("Color Picker", func() {
			picker := dialog.NewColorPicker("Pick a Color", "What is your favorite color?", func(c color.Color) {
				colorPicked(c, win)
			}, win)
			picker.Show()
		}),
		widget.NewButton("Advanced Color Picker", func() {
			picker := dialog.NewColorPicker("Pick a Color", "What is your favorite color?", func(c color.Color) {
				colorPicked(c, win)
			}, win)
			picker.Advanced = true
			picker.Show()
		}),
		widget.NewButton("Form Dialog (Login Form)", func() {
			username := widget.NewEntry()
			username.Validator = validation.NewRegexp(`^[A-Za-z0-9_-]+$`, "username can only contain letters, numbers, '_', and '-'")
			password := widget.NewPasswordEntry()
			password.Validator = validation.NewRegexp(`^[A-Za-z0-9_-]+$`, "password can only contain letters, numbers, '_', and '-'")
			remember := false
			items := []*widget.FormItem{
				widget.NewFormItem("Username", username),
				widget.NewFormItem("Password", password),
				widget.NewFormItem("Remember me", widget.NewCheck("", func(checked bool) {
					remember = checked
				})),
			}

			dialog.ShowForm("Login...", "Log In", "Cancel", items, func(b bool) {
				if !b {
					return
				}
				var rememberText string
				if remember {
					rememberText = "and remember this login"
				}

				log.Println("Please Authenticate", username.Text, password.Text, rememberText)
			}, win)
		}),
	))
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
	res := fyne.NewStaticResource(f.URI().Name(), data)

	return canvas.NewImageFromResource(res)
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

func showImage(f fyne.URIReadCloser) {
	img := loadImage(f)
	if img == nil {
		return
	}
	img.FillMode = canvas.ImageFillOriginal

	w := fyne.CurrentApp().NewWindow(f.URI().Name())
	w.SetContent(container.NewScroll(img))
	w.Resize(fyne.NewSize(320, 240))
	w.Show()
}

func showText(f fyne.URIReadCloser) {
	text := widget.NewLabel(loadText(f))
	text.Wrapping = fyne.TextWrapWord

	w := fyne.CurrentApp().NewWindow(f.URI().Name())
	w.SetContent(container.NewScroll(text))
	w.Resize(fyne.NewSize(320, 240))
	w.Show()
}
