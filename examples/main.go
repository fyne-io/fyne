// Package main provides various examples of Fyne API capabilities
package main

import "errors"
import "fmt"

import "github.com/fyne-io/fyne/examples/apps"

import "github.com/fyne-io/fyne"
import "github.com/fyne-io/fyne/layout"
import "github.com/fyne-io/fyne/theme"
import "github.com/fyne-io/fyne/dialog"
import W "github.com/fyne-io/fyne/widget"

func canvasApp(app fyne.App) {
	apps.Canvas(app)
}

func iconsApp(app fyne.App) {
	apps.Icons(app)
}

func layoutApp(app fyne.App) {
	apps.Layout(app)
}

func formApp(app fyne.App) {
	w := app.NewWindow("Form")
	largeText := W.NewEntry()
	//	largeText.Text = "\n\n\n"

	form := &W.Form{
		OnCancel: func() {
			w.Close()
		},
		OnSubmit: func() {
			fmt.Println("Form submitted")
		},
	}
	form.Append("Name", W.NewEntry())
	form.Append("Email", W.NewEntry())
	form.Append("Message", largeText)
	w.SetContent(form)
	w.Show()
}

func appButton(app fyne.App, label string, onClick func(fyne.App)) *W.Button {
	return &W.Button{Text: label, OnTapped: func() {
		onClick(app)
	}}
}

func confirmCallback(response bool) {
	fmt.Println("Responded with", response)
}

func main() {
	app := apps.NewApp()

	w := app.NewWindow("Examples")
	w.SetContent(&W.Box{Children: []fyne.CanvasObject{
		&W.Toolbar{Items: []W.ToolbarItem{
			&W.ToolbarAction{Icon: theme.MailComposeIcon(), OnActivated: func() { fmt.Println("New") }},
			&W.ToolbarSeparator{},
			&W.ToolbarSpacer{},
			&W.ToolbarAction{Icon: theme.CutIcon(), OnActivated: func() { fmt.Println("Cut") }},
			&W.ToolbarAction{Icon: theme.CopyIcon(), OnActivated: func() { fmt.Println("Copy") }},
			&W.ToolbarAction{Icon: theme.PasteIcon(), OnActivated: func() { fmt.Println("Paste") }},
		}},
		&W.Label{Text: "Fyne Examples!"},

		&W.Button{Text: "Apps", OnTapped: func() {
			dialog.ShowInformation("Information", "Example applications have moved to https://github.com/fyne-io/examples", w)
		}},

		W.NewGroup("Demos", []fyne.CanvasObject{
			appButton(app, "Canvas", canvasApp),
			appButton(app, "Icons", iconsApp),
			appButton(app, "Layout", layoutApp),
			appButton(app, "Form", formApp),
			&W.Entry{Text: "Entry"},
			&W.Check{Text: "Check", OnChanged: func(on bool) { fmt.Println("checked", on) }},
		}...),

		W.NewGroup("Dialogs", []fyne.CanvasObject{
			&W.Button{Text: "Info", OnTapped: func() {
				dialog.ShowInformation("Information", "You should know this thing...", w)
			}},
			&W.Button{Text: "Error", OnTapped: func() {
				err := errors.New("A dummy error message")
				dialog.ShowError(err, w)
			}},
			&W.Button{Text: "Confirm", OnTapped: func() {
				cnf := dialog.NewConfirm("Confirmation", "Are you enjoying this demo?", confirmCallback, w)
				cnf.SetDismissText("Nah")
				cnf.SetConfirmText("Oh Yes!")
				cnf.Show()
			}},
			&W.Button{Text: "Custom", OnTapped: func() {
				dialog.ShowCustom("MyDialog", "Nice", &W.Check{Text: "Inside a dialog"}, w)
			}},
		}...),
		&layout.Spacer{},

		fyne.NewContainerWithLayout(layout.NewGridLayout(2),
			&W.Button{Text: "Dark", OnTapped: func() {
				fyne.GetSettings().SetTheme("dark")
			}},
			&W.Button{Text: "Light", OnTapped: func() {
				fyne.GetSettings().SetTheme("light")
			}},
		),
		&W.Button{Text: "Quit", Icon: theme.CancelIcon(), OnTapped: func() {
			app.Quit()
		}},
	}})
	w.ShowAndRun()
}
