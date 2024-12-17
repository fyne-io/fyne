package tutorials

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var tutorials = map[string]Tutorial2{
	"widgets/button.md": Tutorial2{
		title:   "Widgets : Button",
		content: []string{"The button widget is the basic tappable interaction for an app.\nA user tapping this will run the function passed to the constructor function. \n\n## Basic usage\n\nSimply create a button using the `NewButton` constructor function, passing in a function\nthat should run when the button is tapped.", "btn := widget.NewButton(\"Tap me\", func() {})", "If you want to use an icon in your button that is possible.\nYou can also set the label to \"\" if you want icon only!", "btn := widget.NewButtonWithIcon(\"Home\",\n    theme.HomeIcon(), func() {})", "## Disabled\n\nA button can also be disabled so that it cannot be tapped:", "btn := widget.NewButton(\"Tap me\", func() {})\nbtn.Disable()", "## Importance\n\nYou can change the colour / style of the button by setting its `Importance` value, like this:", "btn := widget.NewButton(\"Danger!\", func() {})\nbtn.Importance = widget.DangerImportance", ""},

		code: []func() fyne.CanvasObject{
			func() fyne.CanvasObject {
				btn := widget.NewButton("Tap me", func() {})
				return btn
			},
			func() fyne.CanvasObject {
				btn := widget.NewButtonWithIcon("Home",
					theme.HomeIcon(), func() {})
				return btn
			},
			func() fyne.CanvasObject {
				btn := widget.NewButton("Tap me", func() {})
				btn.Disable()
				return btn
			},
			func() fyne.CanvasObject {
				btn := widget.NewButton("Danger!", func() {})
				btn.Importance = widget.DangerImportance
				return btn
			},
		},
	},
}
