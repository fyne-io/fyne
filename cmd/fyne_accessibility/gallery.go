package main

import (
	"errors"
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var scenarios = []struct {
	name  string
	build func(fyne.Window) fyne.CanvasObject
}{
	{"controls", controls},
	{"text", text},
	{"collections", collections},
	{"containers", containers},
	{"dynamic", dynamic},
}

func scenarioNames() []string {
	names := make([]string, len(scenarios))
	for i, scenario := range scenarios {
		names[i] = scenario.name
	}
	return names
}

func validScenario(name string) bool {
	if name == "all" {
		return true
	}
	for _, scenario := range scenarios {
		if scenario.name == name {
			return true
		}
	}
	return false
}

func makeGallery(name string, w fyne.Window) (fyne.CanvasObject, error) {
	var tabs []*container.TabItem
	for _, scenario := range scenarios {
		if scenario.name == name {
			return scenario.build(w), nil
		}
		if name == "all" {
			tabs = append(tabs, container.NewTabItem(scenario.name, scenario.build(w)))
		}
	}
	if name == "all" {
		return container.NewAppTabs(tabs...), nil
	}
	return nil, fmt.Errorf("unknown accessibility scenario %q", name)
}

func controls(fyne.Window) fyne.CanvasObject {
	status := widget.NewLabel("Controls ready")
	presses := 0
	button := widget.NewButton("Count presses", func() {
		presses++
		status.SetText(fmt.Sprintf("Pressed %d times", presses))
	})
	disabled := widget.NewButton("Disabled button", func() { status.SetText("Disabled button activated") })
	disabled.Disable()
	check := widget.NewCheck("Enable notifications", func(checked bool) {
		status.SetText(fmt.Sprintf("Notifications: %t", checked))
	})
	checked := widget.NewCheck("Already checked", nil)
	checked.SetChecked(true)
	radio := widget.NewRadioGroup([]string{"Radio Alpha", "Radio Beta"}, nil)
	radio.SetSelected("Radio Alpha")
	radio.OnChanged = func(value string) {
		status.SetText("Radio: " + value)
	}
	checks := widget.NewCheckGroup([]string{"Group Red", "Group Green"}, func(values []string) {
		status.SetText("Group: " + strings.Join(values, ", "))
	})
	selectWidget := widget.NewSelect([]string{"Choice One", "Choice Two"}, func(value string) {
		status.SetText("Choice: " + value)
	})
	selectWidget.PlaceHolder = "Choose an option"
	slider := widget.NewSlider(0, 10)
	slider.SetValue(4)
	slider.OnChanged = func(value float64) { status.SetText(fmt.Sprintf("Slider: %.0f", value)) }
	progress := widget.NewProgressBar()
	progress.SetValue(0.25)
	infinite := widget.NewProgressBarInfinite()
	infinite.MinSize() // CreateRenderer starts the animation, so stop it after creating the renderer.
	infinite.Stop()

	return container.NewBorder(nil, status, nil, nil, container.NewGridWithColumns(2,
		container.NewVBox(button, disabled, check, checked, radio, checks),
		container.NewVBox(selectWidget, widget.NewLabel("Level"), slider,
			widget.NewLabel("Progress"), progress, infinite, widget.NewActivity(),
			widget.NewIcon(theme.InfoIcon()), widget.NewSeparator()),
	))
}

func text(fyne.Window) fyne.CanvasObject {
	status := widget.NewLabel("Text ready")
	name := widget.NewEntry()
	name.SetPlaceHolder("Display name")
	name.SetText("Ada")
	name.OnChanged = func(value string) { status.SetText("Name: " + value) }
	notes := widget.NewMultiLineEntry()
	notes.SetPlaceHolder("Notes")
	notes.SetText("First line\nSecond line")
	notes.OnChanged = func(value string) { status.SetText("Notes: " + value) }
	password := widget.NewPasswordEntry()
	password.SetPlaceHolder("Demo password")
	password.SetText("gallery-password")
	disabled := widget.NewEntry()
	disabled.SetPlaceHolder("Read only example")
	disabled.SetText("Cannot edit")
	disabled.Disable()
	required := widget.NewEntry()
	required.SetPlaceHolder("Required field")
	required.Validator = func(value string) error {
		if strings.TrimSpace(value) == "" {
			return errors.New("a value is required")
		}
		return nil
	}
	form := widget.NewForm(widget.NewFormItem("Required value", required))
	form.OnSubmit = func() { status.SetText("Form submitted") }
	link := widget.NewHyperlink("Example link", nil)
	link.OnTapped = func() { status.SetText("Link activated") }

	return container.NewBorder(nil, status, nil, nil, container.NewGridWithColumns(2,
		container.NewVBox(widget.NewLabel("Plain text"), name, notes, password, disabled),
		container.NewVBox(widget.NewRichTextFromMarkdown("# Heading\n\nRich **text** example."),
			link, widget.NewTextGridFromString("Text grid\nSecond row"), form),
	))
}

func collections(fyne.Window) fyne.CanvasObject {
	status := widget.NewLabel("Collections ready")
	list := widget.NewList(
		func() int { return 30 },
		func() fyne.CanvasObject { return widget.NewLabel("List row 00") },
		func(id widget.ListItemID, object fyne.CanvasObject) {
			object.(*widget.Label).SetText(fmt.Sprintf("List row %02d", id+1))
		},
	)
	list.OnSelected = func(id widget.ListItemID) { status.SetText(fmt.Sprintf("List selected: %d", id+1)) }
	tree := widget.NewTreeWithStrings(map[string][]string{
		"":       {"Branch"},
		"Branch": {"Leaf One", "Leaf Two"},
	})
	tree.OnSelected = func(id widget.TreeNodeID) { status.SetText("Tree selected: " + id) }
	table := widget.NewTableWithHeaders(
		func() (int, int) { return 12, 2 },
		func() fyne.CanvasObject { return widget.NewLabel("Cell 00,0") },
		func(id widget.TableCellID, object fyne.CanvasObject) {
			object.(*widget.Label).SetText(fmt.Sprintf("Cell %d,%d", id.Row+1, id.Col+1))
		},
	)
	table.OnSelected = func(id widget.TableCellID) {
		status.SetText(fmt.Sprintf("Table selected: %d,%d", id.Row+1, id.Col+1))
	}
	grid := widget.NewGridWrap(
		func() int { return 20 },
		func() fyne.CanvasObject { return widget.NewLabel("Grid item 00") },
		func(id widget.GridWrapItemID, object fyne.CanvasObject) {
			object.(*widget.Label).SetText(fmt.Sprintf("Grid item %02d", id+1))
		},
	)
	grid.OnSelected = func(id widget.GridWrapItemID) {
		status.SetText(fmt.Sprintf("Grid selected: %d", id+1))
	}
	scroll := widget.NewButton("Scroll list to end", list.ScrollToBottom)
	return container.NewBorder(scroll, status, nil, nil, container.NewGridWithColumns(2,
		widget.NewCard("List", "", list), widget.NewCard("Tree", "", tree),
		widget.NewCard("Table", "", table), widget.NewCard("Grid wrap", "", grid),
	))
}

func containers(fyne.Window) fyne.CanvasObject {
	status := widget.NewLabel("Containers ready")
	tabs := container.NewAppTabs(
		container.NewTabItem("First tab", widget.NewLabel("First tab content")),
		container.NewTabItem("Second tab", widget.NewLabel("Second tab content")),
	)
	tabs.OnSelected = func(item *container.TabItem) { status.SetText("Tab: " + item.Text) }
	docs := container.NewDocTabs(
		container.NewTabItem("Document A", widget.NewLabel("Document A content")),
		container.NewTabItem("Document B", widget.NewLabel("Document B content")),
	)
	accordion := widget.NewAccordion(
		widget.NewAccordionItem("Expand details", widget.NewLabel("Accordion details")),
	)
	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.HomeIcon(), func() { status.SetText("Toolbar activated") }),
		widget.NewToolbarSeparator(), widget.NewToolbarSpacer(),
	)
	split := container.NewHSplit(tabs, docs)
	return container.NewBorder(toolbar, container.NewVBox(accordion, status), nil, nil, split)
}

func dynamic(w fyne.Window) fyne.CanvasObject {
	status := widget.NewLabel("Dynamic ready")
	extra := widget.NewLabel("Extra content")
	extra.Hide()
	toggle := widget.NewButton("Toggle extra content", func() {
		if extra.Visible() {
			extra.Hide()
			status.SetText("Extra content hidden")
		} else {
			extra.Show()
			status.SetText("Extra content shown")
		}
	})
	action := widget.NewButton("Conditional action", func() { status.SetText("Conditional action activated") })
	action.Disable()
	enable := widget.NewCheck("Enable action", func(enabled bool) {
		if enabled {
			action.Enable()
		} else {
			action.Disable()
		}
	})
	openDialog := widget.NewButton("Open dialog", func() {
		dialog.ShowCustom("Example dialog", "Dismiss dialog", widget.NewLabel("Dialog content"), w)
	})
	openPopup := widget.NewButton("Open popup", func() {
		var popup *widget.PopUp
		popup = widget.NewPopUp(container.NewVBox(
			widget.NewLabel("Popup content"),
			widget.NewButton("Dismiss popup", func() { popup.Hide() }),
		), w.Canvas())
		popup.Show()
	})
	return container.NewBorder(nil, status, nil, nil,
		container.NewVBox(toggle, extra, enable, action, openDialog, openPopup))
}
