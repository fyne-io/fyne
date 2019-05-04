// Package main provides various examples of Fyne API capabilities
package main

import (
	"fmt"
	"net/url"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/canvas"
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

	logo := canvas.NewImageFromResource(theme.FyneLogo())
	logo.SetMinSize(fyne.NewSize(64, 64))

	fyneio, err := url.Parse("https://fyne.io/")
	if err != nil {
		fyne.LogError("Could not parse URL", err)
	}

	w.SetContent(widget.NewVBox(
		widget.NewToolbar(widget.NewToolbarAction(theme.MailComposeIcon(), func() { fmt.Println("New") }),
			widget.NewToolbarSeparator(),
			widget.NewToolbarSpacer(),
			widget.NewToolbarAction(theme.ContentCutIcon(), func() { fmt.Println("Cut") }),
			widget.NewToolbarAction(theme.ContentCopyIcon(), func() { fmt.Println("Copy") }),
			widget.NewToolbarAction(theme.ContentPasteIcon(), func() { fmt.Println("Paste") }),
		),

		widget.NewGroup("Demos", widget.NewVBox(
			widget.NewButton("Canvas", func() { Canvas(a) }),
			widget.NewButton("Icons", func() { Icons(a) }),
			widget.NewButton("Layout", func() { Layout(a) }),
			widget.NewButton("Widgets", func() { Widget(a) }),
			widget.NewButton("Form", func() { formApp(a) }),
			widget.NewButton("Dialogs", func() { Dialogs(a) }),
		)),

		widget.NewGroup("Theme",
			fyne.NewContainerWithLayout(layout.NewGridLayout(2),
				widget.NewButton("Dark", func() {
					a.Settings().SetTheme(theme.DarkTheme())
				}),
				widget.NewButton("Light", func() {
					a.Settings().SetTheme(theme.LightTheme())
				}),
			),
		),

		widget.NewHBox(layout.NewSpacer(), logo, layout.NewSpacer()),
		widget.NewHyperlinkWithStyle("fyne.io", fyneio, fyne.TextAlignCenter, fyne.TextStyle{}),
		layout.NewSpacer(),

		widget.NewButton("Input", func() { Input(a) }),
		widget.NewButton("Advanced", func() { Advanced(a) }),
		widget.NewButtonWithIcon("Quit", theme.CancelIcon(), func() {
			a.Quit()
		}),
	))
	w.ShowAndRun()
}
