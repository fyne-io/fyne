// Package main provides various examples of Fyne API capabilities
package main

import "errors"
import "fmt"
import "log"

import "github.com/fyne-io/fyne/examples/apps"

import "github.com/fyne-io/fyne"
import "github.com/fyne-io/fyne/layout"
import "github.com/fyne-io/fyne/theme"
import "github.com/fyne-io/fyne/dialog"
import W "github.com/fyne-io/fyne/widget"

func canvasApp(app fyne.App) {
	apps.Canvas(app)
}

func layoutApp(app fyne.App) {
	apps.Layout(app)
}

func appButton(app fyne.App, label string, onClick func(fyne.App)) *W.Button {
	return &W.Button{Text: label, OnTapped: func() {
		onClick(app)
	}}
}

func confirmCallback(response bool) {
	log.Println("Responded with", response)
}

func main() {
	app := apps.NewApp()

	w := app.NewWindow("Examples")
	w.SetContent(&W.Box{Children: []fyne.CanvasObject{
		&W.Toolbar{Items: []W.ToolbarItem{
			&W.ToolbarAction{Icon: theme.MailComposeIcon(), OnActivated: func() { log.Println("New") }},
			&W.ToolbarSeparator{},
			&W.ToolbarSpacer{},
			&W.ToolbarAction{Icon: theme.CutIcon(), OnActivated: func() { log.Println("Cut") }},
			&W.ToolbarAction{Icon: theme.CopyIcon(), OnActivated: func() { log.Println("Copy") }},
			&W.ToolbarAction{Icon: theme.PasteIcon(), OnActivated: func() { log.Println("Paste") }},
		}},
		&W.Label{Text: "Fyne Examples!"},

		&W.Button{Text: "Apps", OnTapped: func() {
			dialog.ShowInformation("Information", "Example applications have moved to https://github.com/fyne-io/examples", w)
		}},

		W.NewGroup("Demos", []fyne.CanvasObject{
			appButton(app, "Canvas", canvasApp),
			appButton(app, "Layout", layoutApp),
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
	w.Show()
}
