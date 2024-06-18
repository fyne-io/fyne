package tutorials

import (
	"fmt"
	"image/color"
	"net/url"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/cmd/fyne_demo/data"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/mobile"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const (
	loremIpsum = `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Pellentesque quis consectetur nisi. Suspendisse id interdum felis.
Sed egestas eget tellus eu pharetra. Praesent pulvinar sed massa id placerat. Etiam sem libero, semper vitae consequat ut, volutpat id mi.
Mauris volutpat pellentesque convallis. Curabitur rutrum venenatis orci nec ornare. Maecenas quis pellentesque neque.
Aliquam consectetur dapibus nulla, id maximus odio ultrices ac. Sed luctus at felis sed faucibus. Cras leo augue, congue in velit ut, mattis rhoncus lectus.

Praesent viverra, mauris ut ullamcorper semper, leo urna auctor lectus, vitae vehicula mi leo quis lorem.
Nullam condimentum, massa at tempor feugiat, metus enim lobortis velit, eget suscipit eros ipsum quis tellus. Aenean fermentum diam vel felis dictum semper.
Duis nisl orci, tincidunt ut leo quis, luctus vehicula diam. Sed velit justo, congue id augue eu, euismod dapibus lacus. Proin sit amet imperdiet sapien.
Mauris erat urna, fermentum et quam rhoncus, fringilla consequat ante. Vivamus consectetur molestie odio, ac rutrum erat finibus a.
Suspendisse id maximus felis. Sed mauris odio, mattis eget mi eu, consequat tempus purus.`
)

var (
	progress    *widget.ProgressBar
	fprogress   *widget.ProgressBar
	infProgress *widget.ProgressBarInfinite
	endProgress chan any
)

func makeAccordionTab(_ fyne.Window) fyne.CanvasObject {
	link, err := url.Parse("https://fyne.io/")
	if err != nil {
		fyne.LogError("Could not parse URL", err)
	}
	ac := widget.NewAccordion(
		widget.NewAccordionItem("A", widget.NewHyperlink("One", link)),
		widget.NewAccordionItem("B", widget.NewLabel("Two")),
		&widget.AccordionItem{
			Title:  "C",
			Detail: widget.NewLabel("Three"),
		},
	)
	ac.MultiOpen = true
	ac.Append(widget.NewAccordionItem("D", &widget.Entry{Text: "Four"}))
	return ac
}

func makeActivityTab(win fyne.Window) fyne.CanvasObject {
	a1 := widget.NewActivity()
	a2 := widget.NewActivity()

	var button *widget.Button
	start := func() {
		button.Disable()
		a1.Start()
		a1.Show()
		a2.Start()
		a2.Show()

		defer func() {
			go func() {
				time.Sleep(time.Second * 10)
				a1.Stop()
				a1.Hide()
				a2.Stop()
				a2.Hide()

				button.Enable()
			}()
		}()
	}

	button = widget.NewButton("Animate", start)
	start()

	return container.NewCenter(container.NewGridWithColumns(1,
		container.NewCenter(container.NewVBox(
			container.NewHBox(widget.NewLabel("Working..."), a1),
			container.NewStack(button, a2))),
		container.NewCenter(widget.NewButton("Show dialog", func() {
			prop := canvas.NewRectangle(color.Transparent)
			prop.SetMinSize(fyne.NewSize(50, 50))

			a3 := widget.NewActivity()
			d := dialog.NewCustomWithoutButtons("Please wait...", container.NewStack(prop, a3), win)
			a3.Start()
			d.Show()

			go func() {
				time.Sleep(time.Second * 5)
				a3.Stop()
				d.Hide()
			}()
		}))))
}

func makeButtonTab(_ fyne.Window) fyne.CanvasObject {
	disabled := widget.NewButton("Disabled", func() {})
	disabled.Disable()

	shareItem := fyne.NewMenuItem("Share via", nil)
	shareItem.ChildMenu = fyne.NewMenu("",
		fyne.NewMenuItem("Twitter", func() { fmt.Println("context menu Share->Twitter") }),
		fyne.NewMenuItem("Reddit", func() { fmt.Println("context menu Share->Reddit") }),
	)
	menuLabel := newContextMenuButton("tap me for pop-up menu with submenus", fyne.NewMenu("",
		fyne.NewMenuItem("Copy", func() { fmt.Println("context menu copy") }),
		shareItem,
	))

	return container.NewVScroll(container.NewVBox(
		widget.NewButton("Button (text only)", func() { fmt.Println("tapped text button") }),
		widget.NewButtonWithIcon("Button (text & leading icon)", theme.ConfirmIcon(), func() { fmt.Println("tapped text & leading icon button") }),
		&widget.Button{
			Alignment: widget.ButtonAlignLeading,
			Text:      "Button (leading-aligned, text only)",
			OnTapped:  func() { fmt.Println("tapped leading-aligned, text only button") },
		},
		&widget.Button{
			Alignment:     widget.ButtonAlignTrailing,
			IconPlacement: widget.ButtonIconTrailingText,
			Text:          "Button (trailing-aligned, text & trailing icon)",
			Icon:          theme.ConfirmIcon(),
			OnTapped:      func() { fmt.Println("tapped trailing-aligned, text & trailing icon button") },
		},
		disabled,
		&widget.Button{
			Text:       "Primary button",
			Importance: widget.HighImportance,
			OnTapped:   func() { fmt.Println("high importance button") },
		},
		&widget.Button{
			Text:       "Success button",
			Importance: widget.SuccessImportance,
			OnTapped:   func() { fmt.Println("success importance button") },
		},
		&widget.Button{
			Text:       "Danger button",
			Importance: widget.DangerImportance,
			OnTapped:   func() { fmt.Println("tapped danger button") },
		},
		&widget.Button{
			Text:       "Warning button",
			Importance: widget.WarningImportance,
			OnTapped:   func() { fmt.Println("tapped warning button") },
		},
		layout.NewSpacer(),
		layout.NewSpacer(),
		menuLabel,
		layout.NewSpacer(),
	))
}

func makeCardTab(_ fyne.Window) fyne.CanvasObject {
	card1 := widget.NewCard("Book a table", "Which time suits?",
		widget.NewRadioGroup([]string{"6:30pm", "7:00pm", "7:45pm"}, func(string) {}))
	card2 := widget.NewCard("With media", "No content, with image", nil)
	card2.Image = canvas.NewImageFromResource(data.FyneLogo)
	card3 := widget.NewCard("Title 3", "Another card", widget.NewLabel("Content"))
	return container.NewGridWithColumns(2, container.NewVBox(card1, card3),
		container.NewVBox(card2))
}

func makeEntryTab(_ fyne.Window) fyne.CanvasObject {
	entry := widget.NewEntry()
	entry.SetPlaceHolder("Entry")
	entryDisabled := widget.NewEntry()
	entryDisabled.SetText("Entry (disabled)")
	entryDisabled.Disable()
	entryValidated := newNumEntry()
	entryValidated.SetPlaceHolder("Must contain a number")
	entryMultiLine := widget.NewMultiLineEntry()
	entryMultiLine.SetPlaceHolder("MultiLine Entry")
	entryMultiLine.SetMinRowsVisible(4)

	return container.NewVBox(
		entry,
		entryDisabled,
		entryValidated,
		entryMultiLine)
}

func makeTextGrid() *widget.TextGrid {
	grid := widget.NewTextGridFromString("TextGrid\n\tContent\nZebra")
	grid.SetStyleRange(0, 4, 0, 7,
		&widget.CustomTextGridStyle{BGColor: &color.NRGBA{R: 64, G: 64, B: 192, A: 128}})
	grid.SetRowStyle(1, &widget.CustomTextGridStyle{BGColor: &color.NRGBA{R: 64, G: 192, B: 64, A: 128}})

	white := &widget.CustomTextGridStyle{FGColor: color.White, BGColor: color.Black}
	black := &widget.CustomTextGridStyle{FGColor: color.Black, BGColor: color.White}
	grid.Rows[2].Cells[0].Style = white
	grid.Rows[2].Cells[1].Style = black
	grid.Rows[2].Cells[2].Style = white
	grid.Rows[2].Cells[3].Style = black
	grid.Rows[2].Cells[4].Style = white

	grid.ShowLineNumbers = true
	grid.ShowWhitespace = true

	return grid
}

func makeTextTab(_ fyne.Window) fyne.CanvasObject {
	label := widget.NewLabel("Label")

	link, err := url.Parse("https://fyne.io/")
	if err != nil {
		fyne.LogError("Could not parse URL", err)
	}
	hyperlink := widget.NewHyperlink("Hyperlink", link)

	entryLoremIpsum := widget.NewMultiLineEntry()
	entryLoremIpsum.SetText(loremIpsum)

	label.Alignment = fyne.TextAlignLeading
	hyperlink.Alignment = fyne.TextAlignLeading

	label.Wrapping = fyne.TextWrapWord
	hyperlink.Wrapping = fyne.TextWrapWord
	entryLoremIpsum.Wrapping = fyne.TextWrapWord

	rich := widget.NewRichTextFromMarkdown(`
# RichText Heading

## A Sub Heading

![title](../../theme/icons/fyne.png)

---

* Item1 in _three_ segments
* Item2
* Item3

Normal **Bold** *Italic* [Link](https://fyne.io/) and some ` + "`Code`" + `.
This styled row should also wrap as expected, but only *when required*.

> An interesting quote here, most likely sharing some very interesting wisdom.`)
	rich.Scroll = container.ScrollBoth
	rich.Segments[2].(*widget.ImageSegment).Alignment = fyne.TextAlignTrailing

	radioAlign := widget.NewRadioGroup([]string{"Leading", "Center", "Trailing"}, func(s string) {
		var align fyne.TextAlign
		switch s {
		case "Leading":
			align = fyne.TextAlignLeading
		case "Center":
			align = fyne.TextAlignCenter
		case "Trailing":
			align = fyne.TextAlignTrailing
		}

		label.Alignment = align
		hyperlink.Alignment = align
		for i := range rich.Segments {
			if seg, ok := rich.Segments[i].(*widget.TextSegment); ok {
				seg.Style.Alignment = align
			}
			if seg, ok := rich.Segments[i].(*widget.HyperlinkSegment); ok {
				seg.Alignment = align
			}
		}

		label.Refresh()
		hyperlink.Refresh()
		rich.Refresh()
	})
	radioAlign.Horizontal = true
	radioAlign.SetSelected("Leading")

	radioWrap := widget.NewRadioGroup([]string{"Off", "Scroll", "Break", "Word"}, func(s string) {
		var wrap fyne.TextWrap
		scroll := container.ScrollBoth
		switch s {
		case "Off":
			wrap = fyne.TextWrapOff
			scroll = container.ScrollNone
		case "Scroll":
			wrap = fyne.TextWrapOff
		case "Break":
			wrap = fyne.TextWrapBreak
		case "Word":
			wrap = fyne.TextWrapWord
		}

		label.Wrapping = wrap
		hyperlink.Wrapping = wrap
		entryLoremIpsum.Wrapping = wrap
		entryLoremIpsum.Scroll = scroll
		rich.Wrapping = wrap

		label.Refresh()
		hyperlink.Refresh()
		entryLoremIpsum.Refresh()
		rich.Refresh()
	})
	radioWrap.Horizontal = true
	radioWrap.SetSelected("Word")

	radioTrunc := widget.NewRadioGroup([]string{"Off", "Clip", "Ellipsis"}, func(s string) {
		var trunc fyne.TextTruncation
		switch s {
		case "Off":
			trunc = fyne.TextTruncateOff
		case "Clip":
			trunc = fyne.TextTruncateClip
		case "Ellipsis":
			trunc = fyne.TextTruncateEllipsis
		}

		label.Truncation = trunc
		rich.Truncation = trunc

		label.Refresh()
		hyperlink.Refresh()
		entryLoremIpsum.Refresh()
		rich.Refresh()
	})
	radioTrunc.Horizontal = true
	radioTrunc.SetSelected("Off")

	fixed := container.NewVBox(
		widget.NewForm(
			widget.NewFormItem("Text Alignment", radioAlign),
			widget.NewFormItem("Wrapping", radioWrap),
			widget.NewFormItem("Truncation", radioTrunc)),
		label,
		hyperlink,
	)

	grid := makeTextGrid()
	return container.NewBorder(fixed, grid, nil, nil,
		container.NewGridWithRows(2, rich, entryLoremIpsum))
}

func makeInputTab(_ fyne.Window) fyne.CanvasObject {
	selectEntry := widget.NewSelectEntry([]string{
		"Option A",
		"Option B",
		"Option C",
		"Option D",
		"Option E",
		"Option F",
		"Option G",
		"Option H",
		"Option I",
		"Option J",
		"Option K",
		"Option L",
		"Option M",
		"Option N",
		"Option O",
		"Option P",
		"Option Q",
		"Option R",
		"Option S",
		"Option T",
		"Option U",
		"Option V",
		"Option W",
		"Option X",
		"Option Y",
		"Option Z",
	})
	selectEntry.PlaceHolder = "Type or select"
	disabledCheck := widget.NewCheck("Disabled check", func(bool) {})
	disabledCheck.Disable()
	checkGroup := widget.NewCheckGroup([]string{"CheckGroup Item 1", "CheckGroup Item 2"}, func(s []string) { fmt.Println("selected", s) })
	checkGroup.Horizontal = true
	radio := widget.NewRadioGroup([]string{"Radio Item 1", "Radio Item 2"}, func(s string) { fmt.Println("selected", s) })
	radio.Horizontal = true
	disabledRadio := widget.NewRadioGroup([]string{"Disabled radio"}, func(string) {})
	disabledRadio.Disable()

	disabledSlider := widget.NewSlider(0, 1000)
	disabledSlider.Disable()
	return container.NewVBox(
		widget.NewSelect([]string{"Option 1", "Option 2", "Option 3"}, func(s string) { fmt.Println("selected", s) }),
		selectEntry,
		widget.NewCheck("Check", func(on bool) { fmt.Println("checked", on) }),
		disabledCheck,
		checkGroup,
		radio,
		disabledRadio,
		container.NewBorder(nil, nil, widget.NewLabel("Slider"), nil, widget.NewSlider(0, 1000)),
		container.NewBorder(nil, nil, widget.NewLabel("Disabled slider"), nil, disabledSlider),
	)
}

func makeProgressTab(_ fyne.Window) fyne.CanvasObject {
	stopProgress()

	progress = widget.NewProgressBar()

	fprogress = widget.NewProgressBar()
	fprogress.TextFormatter = func() string {
		return fmt.Sprintf("%.2f out of %.2f", fprogress.Value, fprogress.Max)
	}

	infProgress = widget.NewProgressBarInfinite()
	endProgress = make(chan any, 1)
	startProgress()

	return container.NewVBox(
		widget.NewLabel("Percent"), progress,
		widget.NewLabel("Formatted"), fprogress,
		widget.NewLabel("Infinite"), infProgress)
}

func makeFormTab(_ fyne.Window) fyne.CanvasObject {
	name := widget.NewEntry()
	name.SetPlaceHolder("John Smith")

	email := widget.NewEntry()
	email.SetPlaceHolder("test@example.com")
	email.Validator = validation.NewRegexp(`\w{1,}@\w{1,}\.\w{1,4}`, "not a valid email")

	password := widget.NewPasswordEntry()
	password.SetPlaceHolder("Password")

	disabled := widget.NewRadioGroup([]string{"Option 1", "Option 2"}, func(string) {})
	disabled.Horizontal = true
	disabled.Disable()
	largeText := widget.NewMultiLineEntry()

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Name", Widget: name, HintText: "Your full name"},
			{Text: "Email", Widget: email, HintText: "A valid email address"},
		},
		OnCancel: func() {
			fmt.Println("Cancelled")
		},
		OnSubmit: func() {
			fmt.Println("Form submitted")
			fyne.CurrentApp().SendNotification(&fyne.Notification{
				Title:   "Form for: " + name.Text,
				Content: largeText.Text,
			})
		},
	}
	form.Append("Password", password)
	form.Append("Disabled", disabled)
	form.Append("Message", largeText)
	return form
}

func makeToolbarTab(_ fyne.Window) fyne.CanvasObject {
	t := widget.NewToolbar(widget.NewToolbarAction(theme.MailComposeIcon(), func() { fmt.Println("New") }),
		widget.NewToolbarSeparator(),
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.ContentCutIcon(), func() { fmt.Println("Cut") }),
		widget.NewToolbarAction(theme.ContentCopyIcon(), func() { fmt.Println("Copy") }),
		widget.NewToolbarAction(theme.ContentPasteIcon(), func() { fmt.Println("Paste") }),
	)

	return container.NewBorder(t, nil, nil, nil)
}

func startProgress() {
	progress.SetValue(0)
	fprogress.SetValue(0)
	select { // ignore stale end message
	case <-endProgress:
	default:
	}

	go func() {
		end := endProgress
		num := 0.0
		for num < 1.0 {
			time.Sleep(16 * time.Millisecond)
			select {
			case <-end:
				return
			default:
			}

			progress.SetValue(num)
			fprogress.SetValue(num)
			num += 0.002
		}

		progress.SetValue(1)
		fprogress.SetValue(1)

		// TODO make sure this resets when we hide etc...
		stopProgress()
	}()
	infProgress.Start()
}

func stopProgress() {
	if infProgress == nil {
		return
	}

	if !infProgress.Running() {
		return
	}

	infProgress.Stop()
	endProgress <- struct{}{}
}

// widgetScreen shows a panel containing widget demos
func widgetScreen(_ fyne.Window) fyne.CanvasObject {
	content := container.NewVBox(
		widget.NewLabel("Labels"),
		widget.NewButtonWithIcon("Icons", theme.HomeIcon(), func() {}),
		widget.NewSlider(0, 1))
	return container.NewCenter(content)
}

type contextMenuButton struct {
	widget.Button
	menu *fyne.Menu
}

func (b *contextMenuButton) Tapped(e *fyne.PointEvent) {
	widget.ShowPopUpMenuAtPosition(b.menu, fyne.CurrentApp().Driver().CanvasForObject(b), e.AbsolutePosition)
}

func newContextMenuButton(label string, menu *fyne.Menu) *contextMenuButton {
	b := &contextMenuButton{menu: menu}
	b.Text = label

	b.ExtendBaseWidget(b)
	return b
}

type numEntry struct {
	widget.Entry
}

func (n *numEntry) Keyboard() mobile.KeyboardType {
	return mobile.NumberKeyboard
}

func newNumEntry() *numEntry {
	e := &numEntry{}
	e.ExtendBaseWidget(e)
	e.Validator = validation.NewRegexp(`\d`, "Must contain a number")
	return e
}
