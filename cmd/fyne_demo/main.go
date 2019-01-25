// Package main provides various examples of Fyne API capabilities
package main

import (
	"errors"
	"fmt"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

func formApp(app fyne.App) {
	w := app.NewWindow("Form")

	name := widget.NewEntry()
	name.SetPlaceHolder("John Smith")
	email := widget.NewEntry()
	email.SetPlaceHolder("test@example.com")
	password := widget.NewPasswordEntry()
	largeText := widget.NewMultiLineEntry()

	form := &widget.Form{
		OnCancel: func() {
			w.Close()
		},
		OnSubmit: func() {
			fmt.Println("Form submitted")
			fmt.Println("Name:", name.Text)
			fmt.Println("Email:", email.Text)
			fmt.Println("Password:", password.Text)
			fmt.Println("Message:", largeText.Text)
		},
	}
	form.Append("Name", name)
	form.Append("Email", email)
	form.Append("Password", password)
	form.Append("Message", largeText)
	w.SetContent(form)
	w.Show()
}

func confirmCallback(response bool) {
	fmt.Println("Responded with", response)
}

func main() {
	a := app.New()
	w := a.NewWindow("Fyne Demo")

	cv := canvas.NewImageFromResource(theme.FyneLogo())
	cv.SetMinSize(fyne.NewSize(64, 64))

	w.SetContent(widget.NewVBox(
		widget.NewToolbar(widget.NewToolbarAction(theme.MailComposeIcon(), func() { fmt.Println("New") }),
			widget.NewToolbarSeparator(),
			widget.NewToolbarSpacer(),
			widget.NewToolbarAction(theme.CutIcon(), func() { fmt.Println("Cut") }),
			widget.NewToolbarAction(theme.CopyIcon(), func() { fmt.Println("Copy") }),
			widget.NewToolbarAction(theme.PasteIcon(), func() { fmt.Println("Paste") }),
		),

		widget.NewButton("Apps", func() {
			dialog.ShowInformation("Information", "Example applications have moved to https://github.com/fyne-io/examples", w)
		}),

		widget.NewGroup("Demos",
			widget.NewButton("Canvas", func() { Canvas(a) }),
			widget.NewButton("Icons", func() { Icons(a) }),
			widget.NewButton("Layout", func() { Layout(a) }),
			widget.NewButton("Widgets", func() { Widget(a) }),
			widget.NewButton("Form", func() { formApp(a) }),
		),

		widget.NewHBox(layout.NewSpacer(), cv, layout.NewSpacer()),

		widget.NewGroup("Dialogs",
			widget.NewButton("Info", func() {
				dialog.ShowInformation("Information", "You should know this thing...", w)
			}),
			widget.NewButton("Error", func() {
				err := errors.New("A dummy error message")
				dialog.ShowError(err, w)
			}),
			widget.NewButton("Confirm", func() {
				cnf := dialog.NewConfirm("Confirmation", "Are you enjoying this demo?", confirmCallback, w)
				cnf.SetDismissText("Nah")
				cnf.SetConfirmText("Oh Yes!")
				cnf.Show()
			}),
			widget.NewButton("Custom", func() {
				dialog.ShowCustom("MyDialog", "Nice", widget.NewCheck("Inside a dialog", func(bool) {}), w)
			}),
		),

		layout.NewSpacer(),

		fyne.NewContainerWithLayout(layout.NewGridLayout(2),
			widget.NewButton("Dark", func() {
				a.Settings().SetTheme(theme.DarkTheme())
			}),
			widget.NewButton("Light", func() {
				a.Settings().SetTheme(theme.LightTheme())
			}),
		),
		widget.NewButtonWithIcon("Quit", theme.CancelIcon(), func() {
			a.Quit()
		}),
	))
	w.ShowAndRun()
}
