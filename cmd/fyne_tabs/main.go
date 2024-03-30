package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func main() {

	a := app.New()
	a.Settings().SetTheme(NewThemeDark())
	w := a.NewWindow("Fyne Tabs")

	evenTabs := container.NewAppTabs(
		container.NewTabItemWithIcon("", theme.HomeIcon(), widget.NewLabel("One!")),
		container.NewTabItemWithIcon("", theme.SearchIcon(), widget.NewLabel("Two!")),
	)
	evenTabs.UseMobileLayout()
	evenTabs.SetTabLocation(container.TabLocationBottom)

	leftTabs := container.NewAppTabs(
		container.NewTabItemWithIcon("", theme.HomeIcon(), widget.NewLabel("Three!")),
		container.NewTabItemWithIcon("", theme.SearchIcon(), widget.NewLabel("Four!")),
	)
	leftTabs.SetTabLocation(container.TabLocationBottom)

	tabs := container.NewAppTabs(
		container.NewTabItem("Even", evenTabs),
		container.NewTabItem("Left", leftTabs),
	)
	tabs.UseMobileLayout()
	tabs.SetTabLocation(container.TabLocationTop)

	w.SetContent(tabs)
	w.Resize(fyne.NewSize(400, 700))
	w.ShowAndRun()
}
