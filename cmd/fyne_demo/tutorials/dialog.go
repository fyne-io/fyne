package tutorials

import (
	"fmt"
	"image/color"
	"io/ioutil"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
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
	openFile := widget.NewButton("File Open With Filter (.jpg or .png)", func() {
		fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, win)
				return
			}
			if reader == nil {
				log.Println("Cancelled")
				return
			}

			imageOpened(reader)
		}, win)
		fd.SetFilter(storage.NewExtensionFileFilter([]string{".png", ".jpg", ".jpeg"}))
		fd.Show()
	})
	saveFile := widget.NewButton("File Save", func() {
		dialog.ShowFileSave(func(writer fyne.URIWriteCloser, err error) {
			if err != nil {
				dialog.ShowError(err, win)
				return
			}
			if writer == nil {
				log.Println("Cancelled")
				return
			}

			fileSaved(writer, win)
		}, win)
	})
	openFolder := widget.NewButton("Folder Open", func() {
		dialog.ShowFolderOpen(func(list fyne.ListableURI, err error) {
			if err != nil {
				dialog.ShowError(err, win)
				return
			}
			if list == nil {
				log.Println("Cancelled")
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
	})

	if fyne.CurrentDevice().IsBrowser() {
		openFile.Disable()
		saveFile.Disable()
		openFolder.Disable()
	}

	advancedPicker := dialog.NewColorPicker("Pick a Color", "What is your favorite color?", func(c color.Color) {
		colorPicked(c, win)
	}, win)
	advancedPicker.Advanced = true
	return container.NewVScroll(container.NewVBox(
		widget.NewButton("Custom Buttons", func() {
			custom := dialog.NewCustom("Test", "Nope", widget.NewLabel("Empty"), win)
			//custom.SetButtons([]fyne.CanvasObject{widget.NewButton("Left", custom.Hide), widget.NewButton("Middle", custom.Hide), widget.NewButton("Right", custom.Hide)})
			custom.SetButtons(nil)
			custom.Show()
		}),
	))
}

func imageOpened(f fyne.URIReadCloser) {
	if f == nil {
		log.Println("Cancelled")
		return
	}
	defer f.Close()

	showImage(f)
}

func fileSaved(f fyne.URIWriteCloser, w fyne.Window) {
	defer f.Close()
	_, err := f.Write([]byte("Written by Fyne demo\n"))
	if err != nil {
		dialog.ShowError(err, w)
	}
	err = f.Close()
	if err != nil {
		dialog.ShowError(err, w)
	}
	log.Println("Saved to...", f.URI())
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
